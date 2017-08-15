// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"fmt"

	"github.com/juju/errors"
	jujutxn "github.com/juju/txn"
	"github.com/juju/version"
	"gopkg.in/juju/names.v2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"

	jujucloud "github.com/juju/juju/cloud"
	"github.com/juju/juju/constraints"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/permission"
	"github.com/juju/juju/status"
	"github.com/juju/juju/storage"
)

// modelGlobalKey is the key for the model, its
// settings and constraints.
const modelGlobalKey = "e"

// modelKey will create the kei for a given model using the modelGlobalKey.
func modelKey(modelUUID string) string {
	return fmt.Sprintf("%s#%s", modelGlobalKey, modelUUID)
}

// MigrationMode specifies where the Model is with respect to migration.
type MigrationMode string

const (
	// MigrationModeNone is the default mode for a model and reflects
	// that it isn't involved with a model migration.
	MigrationModeNone MigrationMode = ""

	// MigrationModeExporting reflects a model that is in the process of being
	// exported from one controller to another.
	MigrationModeExporting MigrationMode = "exporting"

	// MigrationModeImporting reflects a model that is being imported into a
	// controller, but is not yet fully active.
	MigrationModeImporting MigrationMode = "importing"
)

// Model represents the state of a model.
type Model struct {
	// globalState is a State that is only safe for accessing non-model
	// specific data, e.g. (non-model) users, models.
	globalState *State
	doc         modelDoc
}

// modelDoc represents the internal state of the model in MongoDB.
type modelDoc struct {
	UUID           string `bson:"_id"`
	Name           string
	Life           Life
	Owner          string        `bson:"owner"`
	ControllerUUID string        `bson:"controller-uuid"`
	MigrationMode  MigrationMode `bson:"migration-mode"`

	// EnvironVersion is the version of the Environ. As providers
	// evolve, cloud resource representations may change; the environ
	// version tracks the current version of that.
	EnvironVersion int `bson:"environ-version"`

	// Cloud is the name of the cloud to which the model is deployed.
	Cloud string `bson:"cloud"`

	// CloudRegion is the name of the cloud region to which the model is
	// deployed. This will be empty for clouds that do not support regions.
	CloudRegion string `bson:"cloud-region,omitempty"`

	// CloudCredential is the ID of the cloud credential that is used
	// for managing cloud resources for this model. This will be empty
	// for clouds that do not require credentials.
	CloudCredential string `bson:"cloud-credential,omitempty"`

	// LatestAvailableTools is a string representing the newest version
	// found while checking streams for new versions.
	LatestAvailableTools string `bson:"available-tools,omitempty"`

	// SLA is the current support level of the model.
	SLA slaDoc `bson:"sla"`

	// MeterStatus is the current meter status of the model.
	MeterStatus modelMeterStatusdoc `bson:"meter-status"`
}

// slaLevel enumerates the support levels available to a model.
type slaLevel string

const (
	slaNone        = slaLevel("")
	SLAUnsupported = slaLevel("unsupported")
	SLAEssential   = slaLevel("essential")
	SLAStandard    = slaLevel("standard")
	SLAAdvanced    = slaLevel("advanced")
)

// String implements fmt.Stringer returning the string representation of an
// SLALevel.
func (l slaLevel) String() string {
	if l == slaNone {
		l = SLAUnsupported
	}
	return string(l)
}

// newSLALevel returns a new SLA level from a string representation.
func newSLALevel(level string) (slaLevel, error) {
	l := slaLevel(level)
	if l == slaNone {
		l = SLAUnsupported
	}
	switch l {
	case SLAUnsupported, SLAEssential, SLAStandard, SLAAdvanced:
		return l, nil
	}
	return l, errors.NotValidf("SLA level %q", level)
}

// slaDoc represents the state of the SLA on the model.
type slaDoc struct {
	// Level is the current support level set on the model.
	Level slaLevel `bson:"level"`

	// Owner is the SLA owner of the model.
	Owner string `bson:"owner,omitempty"`

	// Credentials authenticates the support level setting.
	Credentials []byte `bson:"credentials"`
}

type modelMeterStatusdoc struct {
	Code string `bson:"code"`
	Info string `bson:"info"`
}

// modelEntityRefsDoc records references to the top-level entities
// in the model.
// (anastasiamac 2017-04-10) This is also used to determine if a model can be destroyed.
// Consequently, any changes, especially additions of entities, here,
// would need to be reflected, at least, in Model.checkEmpty(...) as well as
// Model.destroyOps(...)
type modelEntityRefsDoc struct {
	UUID string `bson:"_id"`

	// Machines contains the names of the top-level machines in the model.
	Machines []string `bson:"machines"`

	// Applicatons contains the names of the applications in the model.
	Applications []string `bson:"applications"`

	// Volumes contains the IDs of the volumes in the model.
	Volumes []string `bson:"volumes"`

	// Filesystems contains the IDs of the filesystems in the model.
	Filesystems []string `bson:"filesystems"`
}

// Model returns the model entity.
func (st *State) Model() (*Model, error) {
	model := &Model{
		globalState: st,
	}
	if err := model.refresh(st.modelTag.Id()); err != nil {
		return nil, errors.Trace(err)
	}
	return model, nil
}

// AllModelUUIDs returns the UUIDs for all models in the controller.
// Results are sorted by (name, owner).
func (st *State) AllModelUUIDs() ([]string, error) {
	models, closer := st.db().GetCollection(modelsC)
	defer closer()

	var docs []bson.M
	err := models.Find(nil).Sort("name", "owner").Select(bson.M{"_id": 1}).All(&docs)
	if err != nil {
		return nil, err
	}

	out := make([]string, len(docs))
	for i, doc := range docs {
		out[i] = doc["_id"].(string)
	}
	return out, nil
}

// ModelExists returns true if a model with the supplied UUID exists.
func (st *State) ModelExists(uuid string) (bool, error) {
	models, closer := st.db().GetCollection(modelsC)
	defer closer()

	count, err := models.FindId(uuid).Count()
	if err != nil {
		return false, errors.Annotate(err, "querying model")
	}
	return count > 0, nil
}

// ModelActive returns true if a model with the supplied UUID exists
// and is not being imported as part of a migration.
func (st *State) ModelActive(uuid string) (bool, error) {
	models, closer := st.db().GetCollection(modelsC)
	defer closer()

	var doc modelDoc
	err := models.FindId(uuid).One(&doc)
	if err == mgo.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, errors.Annotate(err, "querying model")
	}
	return doc.MigrationMode != MigrationModeImporting, nil
}

// ModelArgs is a params struct for creating a new model.
type ModelArgs struct {
	// CloudName is the name of the cloud to which the model is deployed.
	CloudName string

	// CloudRegion is the name of the cloud region to which the model is
	// deployed. This will be empty for clouds that do not support regions.
	CloudRegion string

	// CloudCredential is the tag of the cloud credential that will be
	// used for managing cloud resources for this model. This will be
	// empty for clouds that do not require credentials.
	CloudCredential names.CloudCredentialTag

	// Config is the model config.
	Config *config.Config

	// Constraints contains the initial constraints for the model.
	Constraints constraints.Value

	// StorageProviderRegistry is used to determine and store the
	// details of the default storage pools.
	StorageProviderRegistry storage.ProviderRegistry

	// Owner is the user that owns the model.
	Owner names.UserTag

	// MigrationMode is the initial migration mode of the model.
	MigrationMode MigrationMode

	// EnvironVersion is the initial version of the Environ for the model.
	EnvironVersion int
}

// Validate validates the ModelArgs.
func (m ModelArgs) Validate() error {
	if m.Config == nil {
		return errors.NotValidf("nil Config")
	}
	if !names.IsValidCloud(m.CloudName) {
		return errors.NotValidf("Cloud Name %q", m.CloudName)
	}
	if m.Owner == (names.UserTag{}) {
		return errors.NotValidf("empty Owner")
	}
	if m.StorageProviderRegistry == nil {
		return errors.NotValidf("nil StorageProviderRegistry")
	}
	switch m.MigrationMode {
	case MigrationModeNone, MigrationModeImporting:
	default:
		return errors.NotValidf("initial migration mode %q", m.MigrationMode)
	}
	return nil
}

// NewModel creates a new model with its own UUID and
// prepares it for use. Model and State instances for the new
// model are returned.
//
// The controller model's UUID is attached to the new
// model's document. Having the server UUIDs stored with each
// model document means that we have a way to represent external
// models, perhaps for future use around cross model
// relations.
func (st *State) NewModel(args ModelArgs) (_ *Model, _ *State, err error) {
	if err := args.Validate(); err != nil {
		return nil, nil, errors.Trace(err)
	}
	// For now, the model cloud must be the same as the controller cloud.
	controllerInfo, err := st.ControllerInfo()
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	if controllerInfo.CloudName != args.CloudName {
		return nil, nil, errors.NewNotValid(
			nil, fmt.Sprintf("controller cloud %s does not match model cloud %s", controllerInfo.CloudName, args.CloudName))
	}

	// Ensure that the cloud region is valid, or if one is not specified,
	// that the cloud does not support regions.
	controllerCloud, err := st.Cloud(args.CloudName)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	assertCloudRegionOp, err := validateCloudRegion(controllerCloud, args.CloudRegion)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	// Ensure that the cloud credential is valid, or if one is not
	// specified, that the cloud supports the "empty" authentication
	// type.
	owner := args.Owner
	cloudCredentials, err := st.CloudCredentials(owner, args.CloudName)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	assertCloudCredentialOp, err := validateCloudCredential(
		controllerCloud, cloudCredentials, args.CloudCredential,
	)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	if owner.IsLocal() {
		if _, err := st.User(owner); err != nil {
			return nil, nil, errors.Annotate(err, "cannot create model")
		}
	}

	uuid := args.Config.UUID()
	session := st.session.Copy()
	newSt, err := newState(
		names.NewModelTag(uuid),
		controllerInfo.ModelTag,
		session,
		st.mongoInfo,
		st.newPolicy,
		st.clock(),
		st.runTransactionObserver,
	)
	if err != nil {
		return nil, nil, errors.Annotate(err, "could not create state for new model")
	}
	defer func() {
		if err != nil {
			newSt.Close()
		}
	}()
	newSt.controllerModelTag = st.controllerModelTag

	modelOps, modelStatusDoc, err := newSt.modelSetupOps(st.controllerTag.Id(), args, nil)
	if err != nil {
		return nil, nil, errors.Annotate(err, "failed to create new model")
	}

	prereqOps := []txn.Op{
		assertCloudRegionOp,
		assertCloudCredentialOp,
	}
	ops := append(prereqOps, modelOps...)
	err = newSt.db().RunTransaction(ops)
	if err == txn.ErrAborted {

		// We have a  unique key restriction on the "owner" and "name" fields,
		// which will cause the insert to fail if there is another record with
		// the same "owner" and "name" in the collection. If the txn is
		// aborted, check if it is due to the unique key restriction.
		name := args.Config.Name()
		models, closer := st.db().GetCollection(modelsC)
		defer closer()
		envCount, countErr := models.Find(bson.D{
			{"owner", owner.Id()},
			{"name", name}},
		).Count()
		if countErr != nil {
			err = errors.Trace(countErr)
		} else if envCount > 0 {
			err = errors.AlreadyExistsf("model %q for %s", name, owner.Id())
		} else {
			err = errors.New("model already exists")
		}
	}
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	err = newSt.start(st.controllerTag)
	if err != nil {
		return nil, nil, errors.Annotate(err, "could not start state for new model")
	}

	newModel, err := newSt.Model()
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	if args.MigrationMode != MigrationModeImporting {
		probablyUpdateStatusHistory(newSt.db(), modelGlobalKey, modelStatusDoc)
	}

	_, err = newSt.SetUserAccess(newModel.Owner(), newModel.ModelTag(), permission.AdminAccess)
	if err != nil {
		return nil, nil, errors.Annotate(err, "granting admin permission to the owner")
	}

	if err := InitDbLogs(session, uuid); err != nil {
		return nil, nil, errors.Annotate(err, "initialising model logs collection")
	}
	return newModel, newSt, nil
}

// validateCloudRegion validates the given region name against the
// provided Cloud definition, and returns a txn.Op to include in a
// transaction to assert the same.
func validateCloudRegion(cloud jujucloud.Cloud, regionName string) (txn.Op, error) {
	// Ensure that the cloud region is valid, or if one is not specified,
	// that the cloud does not support regions.
	assertCloudRegionOp := txn.Op{
		C:  cloudsC,
		Id: cloud.Name,
	}
	if regionName != "" {
		region, err := jujucloud.RegionByName(cloud.Regions, regionName)
		if err != nil {
			return txn.Op{}, errors.Trace(err)
		}
		assertCloudRegionOp.Assert = bson.D{
			{"regions." + region.Name, bson.D{{"$exists", true}}},
		}
	} else {
		if len(cloud.Regions) > 0 {
			return txn.Op{}, errors.NotValidf("missing CloudRegion")
		}
		assertCloudRegionOp.Assert = bson.D{
			{"regions", bson.D{{"$exists", false}}},
		}
	}
	return assertCloudRegionOp, nil
}

// validateCloudCredential validates the given cloud credential
// name against the provided cloud definition and credentials,
// and returns a txn.Op to include in a transaction to assert the
// same. A user is supplied, for which access to the credential
// will be asserted.
func validateCloudCredential(
	cloud jujucloud.Cloud,
	cloudCredentials map[string]jujucloud.Credential,
	cloudCredential names.CloudCredentialTag,
) (txn.Op, error) {
	if cloudCredential != (names.CloudCredentialTag{}) {
		if cloudCredential.Cloud().Id() != cloud.Name {
			return txn.Op{}, errors.NotValidf("credential %q", cloudCredential.Id())
		}
		var found bool
		for tag := range cloudCredentials {
			if tag == cloudCredential.Id() {
				found = true
				break
			}
		}
		if !found {
			return txn.Op{}, errors.NotFoundf("credential %q", cloudCredential.Id())
		}
		// NOTE(axw) if we add ACLs for credentials,
		// we'll need to check access here. The map
		// we check above contains only the credentials
		// that the model owner has access to.
		return txn.Op{
			C:      cloudCredentialsC,
			Id:     cloudCredentialDocID(cloudCredential),
			Assert: txn.DocExists,
		}, nil
	}
	var hasEmptyAuth bool
	for _, authType := range cloud.AuthTypes {
		if authType != jujucloud.EmptyAuthType {
			continue
		}
		hasEmptyAuth = true
		break
	}
	if !hasEmptyAuth {
		return txn.Op{}, errors.NotValidf("missing CloudCredential")
	}
	return txn.Op{
		C:      cloudsC,
		Id:     cloud.Name,
		Assert: bson.D{{"auth-types", string(jujucloud.EmptyAuthType)}},
	}, nil
}

// Tag returns a name identifying the model.
// The returned name will be different from other Tag values returned
// by any other entities from the same state.
func (m *Model) Tag() names.Tag {
	return m.ModelTag()
}

// ModelTag is the concrete model tag for this model.
func (m *Model) ModelTag() names.ModelTag {
	return names.NewModelTag(m.doc.UUID)
}

// ControllerTag is the tag for the controller that the model is
// running within.
func (m *Model) ControllerTag() names.ControllerTag {
	return names.NewControllerTag(m.doc.ControllerUUID)
}

// UUID returns the universally unique identifier of the model.
func (m *Model) UUID() string {
	return m.doc.UUID
}

// ControllerUUID returns the universally unique identifier of the controller
// in which the model is running.
func (m *Model) ControllerUUID() string {
	return m.doc.ControllerUUID
}

// Name returns the human friendly name of the model.
func (m *Model) Name() string {
	return m.doc.Name
}

// Cloud returns the name of the cloud to which the model is deployed.
func (m *Model) Cloud() string {
	return m.doc.Cloud
}

// CloudRegion returns the name of the cloud region to which the model is deployed.
func (m *Model) CloudRegion() string {
	return m.doc.CloudRegion
}

// CloudCredential returns the tag of the cloud credential used for managing the
// model's cloud resources, and a boolean indicating whether a credential is set.
func (m *Model) CloudCredential() (names.CloudCredentialTag, bool) {
	if names.IsValidCloudCredential(m.doc.CloudCredential) {
		return names.NewCloudCredentialTag(m.doc.CloudCredential), true
	}
	return names.CloudCredentialTag{}, false
}

// MigrationMode returns whether the model is active or being migrated.
func (m *Model) MigrationMode() MigrationMode {
	return m.doc.MigrationMode
}

// SetMigrationMode updates the migration mode of the model.
func (m *Model) SetMigrationMode(mode MigrationMode) error {
	ops := []txn.Op{{
		C:      modelsC,
		Id:     m.doc.UUID,
		Assert: txn.DocExists,
		Update: bson.D{{"$set", bson.D{{"migration-mode", mode}}}},
	}}
	if err := m.globalState.db().RunTransaction(ops); err != nil {
		return errors.Trace(err)
	}
	return m.Refresh()
}

// Life returns whether the model is Alive, Dying or Dead.
func (m *Model) Life() Life {
	return m.doc.Life
}

// Owner returns tag representing the owner of the model.
// The owner is the user that created the model.
func (m *Model) Owner() names.UserTag {
	return names.NewUserTag(m.doc.Owner)
}

// Status returns the status of the model.
func (m *Model) Status() (status.StatusInfo, error) {
	db, closer := m.modelDatabase()
	defer closer()
	status, err := getStatus(db, m.globalKey(), "model")
	if err != nil {
		return status, err
	}
	return status, nil
}

// SetStatus sets the status of the model.
func (m *Model) SetStatus(sInfo status.StatusInfo) error {
	if !status.ValidModelStatus(sInfo.Status) {
		return errors.Errorf("cannot set invalid status %q", sInfo.Status)
	}
	db, closer := m.modelDatabase()
	defer closer()
	return setStatus(db, setStatusParams{
		badge:     "model",
		globalKey: m.globalKey(),
		status:    sInfo.Status,
		message:   sInfo.Message,
		rawData:   sInfo.Data,
		updated:   timeOrNow(sInfo.Since, m.globalState.clock()),
	})
}

// StatusHistory returns a slice of at most filter.Size StatusInfo items
// or items as old as filter.Date or items newer than now - filter.Delta time
// representing past statuses for this application.
func (m *Model) StatusHistory(filter status.StatusHistoryFilter) ([]status.StatusInfo, error) {
	db, closer := m.modelDatabase()
	defer closer()
	args := &statusHistoryArgs{
		db:        db,
		globalKey: m.globalKey(),
		filter:    filter,
	}
	return statusHistory(args)
}

// Config returns the config for the model.
func (m *Model) Config() (*config.Config, error) {
	db, closer := m.modelDatabase()
	defer closer()
	return getModelConfig(db)
}

// UpdateLatestToolsVersion looks up for the latest available version of
// juju tools and updates environementDoc with it.
func (m *Model) UpdateLatestToolsVersion(ver version.Number) error {
	v := ver.String()
	// TODO(perrito666): I need to assert here that there isn't a newer
	// version in place.
	ops := []txn.Op{{
		C:      modelsC,
		Id:     m.doc.UUID,
		Update: bson.D{{"$set", bson.D{{"available-tools", v}}}},
	}}
	err := m.globalState.db().RunTransaction(ops)
	if err != nil {
		return errors.Trace(err)
	}
	return m.Refresh()
}

// LatestToolsVersion returns the newest version found in the last
// check in the streams.
// Bear in mind that the check was performed filtering only
// new patches for the current major.minor. (major.minor.patch)
func (m *Model) LatestToolsVersion() version.Number {
	ver := m.doc.LatestAvailableTools
	if ver == "" {
		return version.Zero
	}
	v, err := version.Parse(ver)
	if err != nil {
		// This is being stored from a valid version but
		// in case this data would beacame corrupt It is not
		// worth to fail because of it.
		return version.Zero
	}
	return v
}

// SLALevel returns the SLA level as a string.
func (m *Model) SLALevel() string {
	return m.doc.SLA.Level.String()
}

// SLAOwner returns the SLA owner as a string. Note that this may differ from
// the model owner.
func (m *Model) SLAOwner() string {
	return m.doc.SLA.Owner
}

// SLACredential returns the SLA credential.
func (m *Model) SLACredential() []byte {
	return m.doc.SLA.Credentials
}

// SetSLA sets the SLA on the model.
func (m *Model) SetSLA(level, owner string, credentials []byte) error {
	l, err := newSLALevel(level)
	if err != nil {
		return errors.Trace(err)
	}
	ops := []txn.Op{{
		C:  modelsC,
		Id: m.doc.UUID,
		Update: bson.D{{"$set", bson.D{{"sla", slaDoc{
			Level:       l,
			Owner:       owner,
			Credentials: credentials,
		}}}}},
	}}
	err = m.globalState.db().RunTransaction(ops)
	if err != nil {
		return errors.Trace(err)
	}
	return m.Refresh()
}

// SetMeterStatus sets the current meter status for this model.
func (m *Model) SetMeterStatus(status, info string) error {
	if _, err := isValidMeterStatusCode(status); err != nil {
		return errors.Trace(err)
	}
	ops := []txn.Op{{
		C:  modelsC,
		Id: m.doc.UUID,
		Update: bson.D{{"$set", bson.D{{"meter-status", modelMeterStatusdoc{
			Code: status,
			Info: info,
		}}}}},
	}}
	err := m.globalState.db().RunTransaction(ops)
	if err != nil {
		return errors.Trace(err)
	}
	return m.Refresh()
}

// MeterStatus returns the current meter status for this model.
func (m *Model) MeterStatus() MeterStatus {
	ms := m.doc.MeterStatus
	return MeterStatus{
		Code: MeterStatusFromString(ms.Code),
		Info: ms.Info,
	}
}

// EnvironVersion is the version of the model's environ -- the related
// cloud provider resources. The environ version is used by the controller
// to identify environ/provider upgrade steps to run for a model's environ
// after the controller is upgraded, or the model is migrated to another
// controller.
func (m *Model) EnvironVersion() int {
	return m.doc.EnvironVersion
}

// SetEnvironVersion sets the model's current environ version. The value
// must be monotonically increasing.
func (m *Model) SetEnvironVersion(v int) error {
	mOrig := m
	mCopy := *m
	m = &mCopy // copy so we can refresh without affecting the original m
	buildTxn := func(attempt int) ([]txn.Op, error) {
		if attempt > 0 {
			if err := m.Refresh(); err != nil {
				return nil, errors.Trace(err)
			}
		}
		if v < m.doc.EnvironVersion {
			return nil, errors.Errorf(
				"cannot set environ version to %v, which is less than the current version %v",
				v, m.doc.EnvironVersion,
			)
		}
		if v == m.doc.EnvironVersion {
			return nil, jujutxn.ErrNoOperations
		}
		return []txn.Op{{
			C:      modelsC,
			Id:     m.doc.UUID,
			Assert: bson.D{{"environ-version", m.doc.EnvironVersion}},
			Update: bson.D{{"$set", bson.D{{"environ-version", v}}}},
		}}, nil
	}
	if err := m.globalState.db().Run(buildTxn); err != nil {
		return errors.Trace(err)
	}
	mOrig.doc.EnvironVersion = v
	return nil
}

// globalKey returns the global database key for the model.
func (m *Model) globalKey() string {
	return modelGlobalKey
}

func (m *Model) Refresh() error {
	return m.refresh(m.UUID())
}

func (m *Model) refresh(uuid string) error {
	models, closer := m.globalState.db().GetCollection(modelsC)
	defer closer()
	err := models.FindId(uuid).One(&m.doc)
	if err == mgo.ErrNotFound {
		return errors.NotFoundf("model")
	}
	return err
}

// Users returns a slice of all users for this model.
func (m *Model) Users() ([]permission.UserAccess, error) {
	db, dbCloser := m.modelDatabase()
	defer dbCloser()
	coll, closer := db.GetCollection(modelUsersC)
	defer closer()

	var userDocs []userAccessDoc
	err := coll.Find(nil).All(&userDocs)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var modelUsers []permission.UserAccess
	for _, doc := range userDocs {
		// check if the User belonging to this model user has
		// been deleted, in this case we should not return it.
		userTag := names.NewUserTag(doc.UserName)
		if userTag.IsLocal() {
			_, err := m.globalState.User(userTag)
			if err != nil {
				if _, ok := err.(DeletedUserError); !ok {
					// We ignore deleted users for now. So if it is not a
					// DeletedUserError we return the error.
					return nil, errors.Trace(err)
				}
				continue
			}
		}
		mu, err := NewModelUserAccess(m.globalState, doc)
		if err != nil {
			return nil, errors.Trace(err)
		}
		modelUsers = append(modelUsers, mu)
	}

	return modelUsers, nil
}

func (m *Model) isControllerModel() bool {
	return m.globalState.controllerModelTag.Id() == m.doc.UUID
}

// DestroyModelParams contains parameters for destroy a model.
type DestroyModelParams struct {
	// DestroyHostedModels controls whether or not hosted models
	// are destroyed also. This only applies to the controller
	// model.
	//
	// If this is false when destroying the controller model,
	// there must be no hosted models, or an error satisfying
	// IsHasHostedModelsError will be returned.
	//
	// TODO(axw) this should be moved to the Controller type.
	DestroyHostedModels bool

	// DestroyStorage controls whether or not storage in the
	// model (and hosted models, if DestroyHostedModels is true)
	// should be destroyed.
	//
	// This is ternary: nil, false, or true. If nil and
	// there is persistent storage in the model (or hosted
	// models), an error satisfying IsHasPersistentStorageError
	// will be returned.
	DestroyStorage *bool
}

// Destroy sets the models's lifecycle to Dying, preventing
// addition of services or machines to state. If called on
// an empty hosted model, the lifecycle will be advanced
// straight to Dead.
func (m *Model) Destroy(args DestroyModelParams) (err error) {
	defer errors.DeferredAnnotatef(&err, "failed to destroy model")

	buildTxn := func(attempt int) ([]txn.Op, error) {
		// On the first attempt, we assume memory state is recent
		// enough to try using...
		if attempt != 0 {
			// ...but on subsequent attempts, we read fresh environ
			// state from the DB. Note that we do *not* refresh the
			// original `m` itself, as detailed in doc/hacking-state.txt
			if attempt == 1 {
				mCopy := *m
				m = &mCopy
			}
			if err := m.Refresh(); err != nil {
				return nil, errors.Trace(err)
			}
		}

		ops, err := m.destroyOps(args, false, false)
		if err == errModelNotAlive {
			return nil, jujutxn.ErrNoOperations
		} else if err != nil {
			return nil, errors.Trace(err)
		}

		return ops, nil
	}

	return m.globalState.db().RunFor(m.UUID(), buildTxn)
}

// errModelNotAlive is a signal emitted from destroyOps to indicate
// that model destruction is already underway.
var errModelNotAlive = errors.New("model is no longer alive")

type hasHostedModelsError int

func (e hasHostedModelsError) Error() string {
	return fmt.Sprintf("hosting %d other models", e)
}

// IsHasHostedModelsError reports whether or not the given error
// was caused by an attempt to destroy the controller model while
// it contained non-empty hosted models, without specifying that
// they should also be destroyed.
func IsHasHostedModelsError(err error) bool {
	_, ok := errors.Cause(err).(hasHostedModelsError)
	return ok
}

type hasPersistentStorageError struct{}

func (hasPersistentStorageError) Error() string {
	return "model contains persistent storage"
}

// IsHasPersistentStorageError reports whether or not the given
// error was caused by an attempt to destroy a model while it
// contained persistent storage, without specifying how the
// storage should be removed (destroyed or released).
func IsHasPersistentStorageError(err error) bool {
	_, ok := errors.Cause(err).(hasPersistentStorageError)
	return ok
}

type modelNotEmptyError struct {
	error
}

// destroyOps returns the txn operations necessary to begin model
// destruction, or an error indicating why it can't.
//
// If ensureEmpty is true, then destroyOps will return an error
// if the model is non-empty.
//
// If destroyingController is true, then destroyOps will progress
// empty models to Dead, but otherwise will return only non-mutating
// ops to assert the current state of the model, and will leave it
// to the "models" cleanup to destroy the model.
func (m *Model) destroyOps(
	args DestroyModelParams,
	ensureEmpty bool,
	destroyingController bool,
) ([]txn.Op, error) {
	if m.Life() != Alive {
		return nil, errModelNotAlive
	}

	// Check if the model is empty. If it is, we can advance the model's
	// lifecycle state directly to Dead.
	modelEntityRefs, err := m.getEntityRefs()
	if err != nil {
		return nil, errors.Annotate(err, "getting model entity refs")
	}
	modelUUID := m.UUID()
	isEmpty := true
	nextLife := Dying
	prereqOps, err := checkModelEntityRefsEmpty(modelEntityRefs)
	if err != nil {
		if ensureEmpty {
			return nil, modelNotEmptyError{err}
		}
		isEmpty = false
		prereqOps = nil
		if args.DestroyStorage == nil {
			// The model is non-empty, and the user has not specified
			// whether storage should be destroyed or released. Make
			// sure there are no filesystems or volumes in the model.
			db, closer := m.modelDatabase()
			defer closer()
			storageOps, err := checkModelEntityRefsNoPersistentStorage(
				db, modelEntityRefs,
			)
			if err != nil {
				return nil, err
			}
			prereqOps = storageOps
		}
	} else {
		if !m.isControllerModel() {
			// The model is empty, and is not the controller
			// model, so we can move it straight to Dead.
			nextLife = Dead
		}
	}

	if m.isControllerModel() && (!args.DestroyHostedModels || args.DestroyStorage == nil) {
		// This is the controller model, and we've not been instructed
		// to destroy hosted models, or we've not been instructed how
		// to remove storage.
		//
		// Check for any Dying or alive but non-empty models. If there
		// are any and we have not been instructed to destroy them, we
		// return an error indicating that there are hosted models.

		// We need access State instances for hosted models and it's
		// too hard to thread an external StatePool to here, so create
		// a fresh one. Creating new States is relatively slow but
		// this is ok because this is an infrequently used code path.
		pool := NewStatePool(m.globalState)
		defer pool.Close()

		modelUUIDs, err := m.globalState.AllModelUUIDs()
		if err != nil {
			return nil, errors.Trace(err)
		}
		var aliveEmpty, aliveNonEmpty, dying, dead int
		for _, modelUUID := range modelUUIDs {
			if modelUUID == m.UUID() {
				// Ignore the controller model.
				continue
			}

			st, release, err := pool.Get(modelUUID)
			if err != nil {
				return nil, errors.Trace(err)
			}
			defer release()

			model, err := st.Model()
			if err != nil {
				return nil, errors.Trace(err)
			}

			if model.Life() == Dead {
				// Dead hosted models don't affect
				// whether the controller can be
				// destroyed or not, but they are
				// still counted in the hosted models.
				dead++
				continue
			}
			// See if the model is empty, and if it is,
			// get the ops required to ensure it can be
			// destroyed, but without effecting the
			// destruction. The destruction is carried
			// out by the cleanup.
			ops, err := model.destroyOps(args, !args.DestroyHostedModels, true)
			switch err {
			case errModelNotAlive:
				dying++
			case nil:
				prereqOps = append(prereqOps, ops...)
				aliveEmpty++
			default:
				if _, ok := err.(modelNotEmptyError); !ok {
					return nil, errors.Trace(err)
				}
				aliveNonEmpty++
			}
		}
		if !args.DestroyHostedModels && (dying > 0 || aliveNonEmpty > 0) {
			// There are Dying, or Alive but non-empty models.
			// We cannot destroy the controller without first
			// destroying the models and waiting for them to
			// become Dead.
			return nil, errors.Trace(
				hasHostedModelsError(dying + aliveNonEmpty + aliveEmpty),
			)
		}
		// Ensure that the number of active models has not changed
		// between the query and when the transaction is applied.
		//
		// Note that we assert that each empty model that we intend
		// move to Dead is still Alive, so we're protected from an
		// ABA style problem where an empty model is concurrently
		// removed, and replaced with a non-empty model.
		prereqOps = append(prereqOps,
			assertHostedModelsOp(aliveEmpty+aliveNonEmpty+dying+dead),
		)
	}

	ops := []txn.Op{{
		C:      modelsC,
		Id:     modelUUID,
		Assert: isAliveDoc,
	}}
	if !destroyingController || nextLife == Dead {
		timeOfDying := m.globalState.nowToTheSecond()
		modelUpdateValues := bson.D{
			{"life", nextLife},
			{"time-of-dying", timeOfDying},
		}
		if nextLife == Dead {
			modelUpdateValues = append(modelUpdateValues, bson.DocElem{
				"time-of-death", timeOfDying,
			})
		}
		ops[0].Update = bson.D{{"$set", modelUpdateValues}}
	} else {
		// We're destroying the controller, and we're not
		// progressing the model directly to Dead. We leave
		// it to the "models" cleanup to destroy the model.
	}
	if destroyingController {
		// We're destroying the controller model, and being asked
		// to check the validity of destroying hosted models,
		// progressing empty models to Dead. We don't want to
		// include the cleanups below, as they assume we're
		// destroying the model.
		return append(prereqOps, ops...), nil
	}

	// Because txn operations execute in order, and may encounter
	// arbitrarily long delays, we need to make sure every op
	// causes a state change that's still consistent; so we make
	// sure the cleanup ops are the last thing that will execute.
	if m.isControllerModel() {
		ops = append(ops, newCleanupOp(
			cleanupModelsForDyingController, modelUUID,
			// pass through the DestroyModelArgs to the cleanup,
			// so the models can be destroyed according to the
			// same rules.
			args,
		))
	}
	if !isEmpty {
		// We only need to destroy resources if the model is non-empty.
		// It wouldn't normally be harmful to enqueue the cleanups
		// otherwise, except for when we're destroying an empty
		// hosted model in the course of destroying the controller. In
		// that case we'll get errors if we try to enqueue hosted-model
		// cleanups, because the cleanups collection is non-global.
		ops = append(ops,
			newCleanupOp(cleanupMachinesForDyingModel, modelUUID),
			newCleanupOp(cleanupApplicationsForDyingModel, modelUUID),
		)
		if args.DestroyStorage != nil {
			// The user has specified that the storage should be destroyed
			// or released, which we can do in a cleanup. If the user did
			// not specify either, then we have already added prereq ops
			// to assert that there is no storage in the model.
			ops = append(ops, newCleanupOp(
				cleanupStorageForDyingModel, modelUUID,
				// pass through DestroyModelArgs.DestroyStorage to the
				// cleanup, so the storage can be destroyed/released
				// according to the parameters.
				*args.DestroyStorage,
			))
		}
	}
	return append(prereqOps, ops...), nil
}

// getEntityRefs reads the current model entity refs document for the model.
func (m *Model) getEntityRefs() (*modelEntityRefsDoc, error) {
	modelEntityRefs, closer := m.globalState.db().GetCollection(modelEntityRefsC)
	defer closer()

	var doc modelEntityRefsDoc
	if err := modelEntityRefs.FindId(m.UUID()).One(&doc); err != nil {
		if err == mgo.ErrNotFound {
			return nil, errors.NotFoundf("entity references doc for model %s", m.UUID())
		}
		return nil, errors.Annotatef(err, "getting entity references for model %s", m.UUID())
	}
	return &doc, nil
}

// checkModelEntityRefsEmpty checks that the model is empty of any entities
// that may require external resource cleanup. If the model is not empty,
// then an error will be returned; otherwise txn.Ops are returned to assert
// the continued emptiness.
func checkModelEntityRefsEmpty(doc *modelEntityRefsDoc) ([]txn.Op, error) {
	// These errors could be potentially swallowed as we re-try to destroy model.
	// Let's, at least, log them for observation.
	if n := len(doc.Machines); n > 0 {
		logger.Infof("model is still not empty, has machines: %v", doc.Machines)
		return nil, errors.Errorf("model not empty, found %d machine(s)", n)
	}
	if n := len(doc.Applications); n > 0 {
		logger.Infof("model is still not empty, has applications: %v", doc.Applications)
		return nil, errors.Errorf("model not empty, found %d application(s)", n)
	}
	if n := len(doc.Volumes); n > 0 {
		logger.Infof("model is still not empty, has volumes: %v", doc.Volumes)
		return nil, errors.Errorf("model not empty, found %d volume(s)", n)
	}
	if n := len(doc.Filesystems); n > 0 {
		logger.Infof("model is still not empty, has file systems: %v", doc.Filesystems)
		return nil, errors.Errorf("model not empty, found %d filesystem(s)", n)
	}
	return []txn.Op{{
		C:  modelEntityRefsC,
		Id: doc.UUID,
		Assert: bson.D{
			{"machines", bson.D{{"$size", 0}}},
			{"applications", bson.D{{"$size", 0}}},
			{"volumes", bson.D{{"$size", 0}}},
			{"filesystems", bson.D{{"$size", 0}}},
		},
	}}, nil
}

// checkModelEntityRefsNoPersistentStorage checks that there is no
// persistent storage in the model. If there is, then an error of
// type hasPersistentStorageError is returned. If there is not,
// txn.Ops are returned to assert the same.
func checkModelEntityRefsNoPersistentStorage(
	db Database, doc *modelEntityRefsDoc,
) ([]txn.Op, error) {
	for _, volumeId := range doc.Volumes {
		volumeTag := names.NewVolumeTag(volumeId)
		detachable, err := isDetachableVolumeTag(db, volumeTag)
		if err != nil {
			return nil, errors.Trace(err)
		}
		if detachable {
			return nil, hasPersistentStorageError{}
		}
	}
	for _, filesystemId := range doc.Filesystems {
		filesystemTag := names.NewFilesystemTag(filesystemId)
		detachable, err := isDetachableFilesystemTag(db, filesystemTag)
		if err != nil {
			return nil, errors.Trace(err)
		}
		if detachable {
			return nil, hasPersistentStorageError{}
		}
	}
	noNewVolumes := bson.DocElem{
		"volumes", bson.D{{
			"$not", bson.D{{
				"$elemMatch", bson.D{{
					"$nin", doc.Volumes,
				}},
			}},
		}},
		// There are no volumes that are not in
		// the set of volumes we previously knew
		// about => the current set of volumes
		// is a subset of the previously known set.
	}
	noNewFilesystems := bson.DocElem{
		"filesystems", bson.D{{
			"$not", bson.D{{
				"$elemMatch", bson.D{{
					"$nin", doc.Filesystems,
				}},
			}},
		}},
	}
	return []txn.Op{{
		C:  modelEntityRefsC,
		Id: doc.UUID,
		Assert: bson.D{
			noNewVolumes,
			noNewFilesystems,
		},
	}}, nil
}

func addModelMachineRefOp(mb modelBackend, machineId string) txn.Op {
	return addModelEntityRefOp(mb, "machines", machineId)
}

func removeModelMachineRefOp(mb modelBackend, machineId string) txn.Op {
	return removeModelEntityRefOp(mb, "machines", machineId)
}

func addModelApplicationRefOp(mb modelBackend, applicationname string) txn.Op {
	return addModelEntityRefOp(mb, "applications", applicationname)
}

func removeModelApplicationRefOp(mb modelBackend, applicationname string) txn.Op {
	return removeModelEntityRefOp(mb, "applications", applicationname)
}

func addModelVolumeRefOp(mb modelBackend, volumeId string) txn.Op {
	return addModelEntityRefOp(mb, "volumes", volumeId)
}

func removeModelVolumeRefOp(mb modelBackend, volumeId string) txn.Op {
	return removeModelEntityRefOp(mb, "volumes", volumeId)
}

func addModelFilesystemRefOp(mb modelBackend, filesystemId string) txn.Op {
	return addModelEntityRefOp(mb, "filesystems", filesystemId)
}

func removeModelFilesystemRefOp(mb modelBackend, filesystemId string) txn.Op {
	return removeModelEntityRefOp(mb, "filesystems", filesystemId)
}

func addModelEntityRefOp(mb modelBackend, entityField, entityId string) txn.Op {
	return txn.Op{
		C:      modelEntityRefsC,
		Id:     mb.modelUUID(),
		Assert: txn.DocExists,
		Update: bson.D{{"$addToSet", bson.D{{entityField, entityId}}}},
	}
}

func removeModelEntityRefOp(mb modelBackend, entityField, entityId string) txn.Op {
	return txn.Op{
		C:      modelEntityRefsC,
		Id:     mb.modelUUID(),
		Update: bson.D{{"$pull", bson.D{{entityField, entityId}}}},
	}
}

// createModelOp returns the operation needed to create
// an model document with the given name and UUID.
func createModelOp(
	owner names.UserTag,
	name, uuid, controllerUUID, cloudName, cloudRegion string,
	cloudCredential names.CloudCredentialTag,
	migrationMode MigrationMode,
	environVersion int,
) txn.Op {
	doc := &modelDoc{
		UUID:            uuid,
		Name:            name,
		Life:            Alive,
		Owner:           owner.Id(),
		ControllerUUID:  controllerUUID,
		MigrationMode:   migrationMode,
		EnvironVersion:  environVersion,
		Cloud:           cloudName,
		CloudRegion:     cloudRegion,
		CloudCredential: cloudCredential.Id(),
	}
	return txn.Op{
		C:      modelsC,
		Id:     uuid,
		Assert: txn.DocMissing,
		Insert: doc,
	}
}

func createModelEntityRefsOp(uuid string) txn.Op {
	return txn.Op{
		C:      modelEntityRefsC,
		Id:     uuid,
		Assert: txn.DocMissing,
		Insert: &modelEntityRefsDoc{UUID: uuid},
	}
}

const hostedModelCountKey = "hostedModelCount"

type hostedModelCountDoc struct {
	// RefCount is the number of models in the Juju system.
	// We do not count the system model.
	RefCount int `bson:"refcount"`
}

func assertHostedModelsOp(n int) txn.Op {
	return txn.Op{
		C:      controllersC,
		Id:     hostedModelCountKey,
		Assert: bson.D{{"refcount", n}},
	}
}

func incHostedModelCountOp() txn.Op {
	return HostedModelCountOp(1)
}

func decHostedModelCountOp() txn.Op {
	return HostedModelCountOp(-1)
}

func HostedModelCountOp(amount int) txn.Op {
	return txn.Op{
		C:      controllersC,
		Id:     hostedModelCountKey,
		Assert: txn.DocExists,
		Update: bson.M{
			"$inc": bson.M{"refcount": amount},
		},
	}
}

func hostedModelCount(st *State) (int, error) {
	var doc hostedModelCountDoc
	controllers, closer := st.db().GetCollection(controllersC)
	defer closer()

	if err := controllers.Find(bson.D{{"_id", hostedModelCountKey}}).One(&doc); err != nil {
		return 0, errors.Trace(err)
	}
	return doc.RefCount, nil
}

// createUniqueOwnerModelNameOp returns the operation needed to create
// an usermodelnameC document with the given owner and model name.
func createUniqueOwnerModelNameOp(owner names.UserTag, envName string) txn.Op {
	return txn.Op{
		C:      usermodelnameC,
		Id:     userModelNameIndex(owner.Id(), envName),
		Assert: txn.DocMissing,
		Insert: bson.M{},
	}
}

// assertAliveOp returns a txn.Op that asserts the model is alive.
func (m *Model) assertActiveOp() txn.Op {
	return assertModelActiveOp(m.UUID())
}

// assertModelActiveOp returns a txn.Op that asserts the given
// model UUID refers to an Alive model.
func assertModelActiveOp(modelUUID string) txn.Op {
	return txn.Op{
		C:      modelsC,
		Id:     modelUUID,
		Assert: append(isAliveDoc, bson.DocElem{"migration-mode", MigrationModeNone}),
	}
}

func checkModelActive(st *State) error {
	model, err := st.Model()
	if (err == nil && model.Life() != Alive) || errors.IsNotFound(err) {
		return errors.Errorf("model %q is no longer alive", model.Name())
	} else if err != nil {
		return errors.Annotate(err, "unable to read model")
	} else if mode := model.MigrationMode(); mode != MigrationModeNone {
		return errors.Errorf("model %q is being migrated", model.Name())
	}
	return nil
}

// modelDatabase returns a Database scoped to the model's UUID,
// and a function that will close the database when called.
func (m *Model) modelDatabase() (Database, func()) {
	return m.globalState.db().CopyForModel(m.UUID())
}
