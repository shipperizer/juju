// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasapplicationprovisioner

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/names/v5"
	"github.com/juju/retry"

	"github.com/juju/juju/caas"
	"github.com/juju/juju/core/life"
	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/core/status"
	"github.com/juju/juju/internal/charm"
	"github.com/juju/juju/internal/cloudconfig/podcfg"
	"github.com/juju/juju/rpc/params"
)

// ApplicationOps defines all the operations the application worker can perform.
// This is exported for testing only.
type ApplicationOps interface {
	AppAlive(ctx context.Context, appName string, app caas.Application,
		password string, lastApplied *caas.ApplicationConfig,
		facade CAASProvisionerFacade, clk clock.Clock, logger logger.Logger) error

	AppDying(ctx context.Context, appName string, app caas.Application, appLife life.Value,
		facade CAASProvisionerFacade, unitFacade CAASUnitProvisionerFacade, logger logger.Logger) error

	AppDead(ctx context.Context, appName string, app caas.Application,
		broker CAASBroker, facade CAASProvisionerFacade, unitFacade CAASUnitProvisionerFacade,
		clk clock.Clock, logger logger.Logger) error

	CheckCharmFormat(ctx context.Context, appName string,
		facade CAASProvisionerFacade, logger logger.Logger) (isOk bool, err error)

	EnsureTrust(ctx context.Context, appName string, app caas.Application,
		unitFacade CAASUnitProvisionerFacade, logger logger.Logger) error

	UpdateState(ctx context.Context, appName string, app caas.Application, lastReportedStatus map[string]status.StatusInfo,
		broker CAASBroker, facade CAASProvisionerFacade, unitFacade CAASUnitProvisionerFacade, logger logger.Logger) (map[string]status.StatusInfo, error)

	RefreshApplicationStatus(ctx context.Context, appName string, app caas.Application, appLife life.Value,
		facade CAASProvisionerFacade, logger logger.Logger) error

	WaitForTerminated(appName string, app caas.Application,
		clk clock.Clock) error

	ReconcileDeadUnitScale(ctx context.Context, appName string, app caas.Application,
		facade CAASProvisionerFacade, logger logger.Logger) error

	EnsureScale(ctx context.Context, appName string, app caas.Application, appLife life.Value,
		facade CAASProvisionerFacade, unitFacade CAASUnitProvisionerFacade, logger logger.Logger) error
}

type applicationOps struct{}

func (applicationOps) AppAlive(
	ctx context.Context,
	appName string, app caas.Application, password string,
	lastApplied *caas.ApplicationConfig, facade CAASProvisionerFacade,
	clk clock.Clock, logger logger.Logger,
) error {
	return appAlive(ctx, appName, app, password, lastApplied, facade, clk, logger)
}

func (applicationOps) AppDying(
	ctx context.Context,
	appName string, app caas.Application, appLife life.Value,
	facade CAASProvisionerFacade, unitFacade CAASUnitProvisionerFacade,
	logger logger.Logger,
) error {
	return appDying(ctx, appName, app, appLife, facade, unitFacade, logger)
}

func (applicationOps) AppDead(ctx context.Context,
	appName string, app caas.Application,
	broker CAASBroker, facade CAASProvisionerFacade, unitFacade CAASUnitProvisionerFacade,
	clk clock.Clock, logger logger.Logger,
) error {
	return appDead(ctx, appName, app, broker, facade, unitFacade, clk, logger)
}

func (applicationOps) CheckCharmFormat(
	ctx context.Context, appName string,
	facade CAASProvisionerFacade, logger logger.Logger) (isOk bool, err error) {
	return checkCharmFormat(ctx, appName, facade, logger)
}

func (applicationOps) EnsureTrust(
	ctx context.Context,
	appName string, app caas.Application,
	unitFacade CAASUnitProvisionerFacade,
	logger logger.Logger,
) error {
	return ensureTrust(ctx, appName, app, unitFacade, logger)
}

func (applicationOps) UpdateState(
	ctx context.Context,
	appName string, app caas.Application, lastReportedStatus map[string]status.StatusInfo,
	broker CAASBroker, facade CAASProvisionerFacade, unitFacade CAASUnitProvisionerFacade,
	logger logger.Logger,
) (map[string]status.StatusInfo, error) {
	return updateState(ctx, appName, app, lastReportedStatus, broker, facade, unitFacade, logger)
}

func (applicationOps) RefreshApplicationStatus(
	ctx context.Context,
	appName string, app caas.Application, appLife life.Value,
	facade CAASProvisionerFacade,
	logger logger.Logger,
) error {
	return refreshApplicationStatus(ctx, appName, app, appLife, facade, logger)
}

func (applicationOps) WaitForTerminated(
	appName string, app caas.Application,
	clk clock.Clock,
) error {
	return waitForTerminated(appName, app, clk)
}

func (applicationOps) ReconcileDeadUnitScale(
	ctx context.Context,
	appName string, app caas.Application,
	facade CAASProvisionerFacade,
	logger logger.Logger,
) error {
	return reconcileDeadUnitScale(ctx, appName, app, facade, logger)
}

func (applicationOps) EnsureScale(
	ctx context.Context,
	appName string, app caas.Application, appLife life.Value,
	facade CAASProvisionerFacade, unitFacade CAASUnitProvisionerFacade,
	logger logger.Logger,
) error {
	return ensureScale(ctx, appName, app, appLife, facade, unitFacade, logger)
}

type Tomb interface {
	Dying() <-chan struct{}
	ErrDying() error
}

// appAlive handles the life.Alive state for the CAAS application. It handles invoking the
// CAAS broker to create the resources in the k8s cluster for this application.
func appAlive(ctx context.Context, appName string, app caas.Application,
	password string, lastApplied *caas.ApplicationConfig,
	facade CAASProvisionerFacade, clk clock.Clock, logger logger.Logger) error {
	logger.Debugf("ensuring application %q exists", appName)

	provisionInfo, err := facade.ProvisioningInfo(ctx, appName)
	if err != nil {
		return errors.Annotate(err, "retrieving provisioning info")
	}
	if provisionInfo.CharmURL == nil {
		return errors.Errorf("missing charm url in provision info")
	}

	charmInfo, err := facade.CharmInfo(ctx, provisionInfo.CharmURL.String())
	if err != nil {
		return errors.Annotatef(err, "retrieving charm deployment info for %q", appName)
	}

	appState, err := app.Exists()
	if err != nil {
		return errors.Annotatef(err, "retrieving application state for %q", appName)
	}

	if appState.Exists && appState.Terminating {
		if err := waitForTerminated(appName, app, clk); err != nil {
			return errors.Annotatef(err, "%q was terminating and there was an error waiting for it to stop", appName)
		}
	}

	images, err := facade.ApplicationOCIResources(ctx, appName)
	if err != nil {
		return errors.Annotate(err, "getting OCI image resources")
	}

	ch := charmInfo.Charm()
	charmBaseImage, err := podcfg.ImageForBase(provisionInfo.ImageDetails.Repository, charm.Base{
		Name: provisionInfo.Base.OS,
		Channel: charm.Channel{
			Track: provisionInfo.Base.Channel.Track,
			Risk:  charm.Risk(provisionInfo.Base.Channel.Risk),
		},
	})
	if err != nil {
		return errors.Annotate(err, "getting image for base")
	}

	containers := make(map[string]caas.ContainerConfig)
	for k, v := range ch.Meta().Containers {
		container := caas.ContainerConfig{
			Name: k,
			Uid:  v.Uid,
			Gid:  v.Gid,
		}
		if v.Resource == "" {
			return errors.NotValidf("empty container resource reference")
		}
		image, ok := images[v.Resource]
		if !ok {
			return errors.NotFoundf("referenced charm base image resource %s", v.Resource)
		}
		container.Image = image
		for _, m := range v.Mounts {
			container.Mounts = append(container.Mounts, caas.MountConfig{
				StorageName: m.Storage,
				Path:        m.Location,
			})
		}
		containers[k] = container
	}

	// TODO(sidecar): container.Mounts[*].Path <= consolidate? => provisionInfo.Filesystems[*].Attachment.Path
	config := caas.ApplicationConfig{
		IsPrivateImageRepo:   provisionInfo.ImageDetails.IsPrivate(),
		IntroductionSecret:   password,
		AgentVersion:         provisionInfo.Version,
		AgentImagePath:       provisionInfo.ImageDetails.RegistryPath,
		ControllerAddresses:  strings.Join(provisionInfo.APIAddresses, ","),
		ControllerCertBundle: provisionInfo.CACert,
		ResourceTags:         provisionInfo.Tags,
		Constraints:          provisionInfo.Constraints,
		Filesystems:          provisionInfo.Filesystems,
		Devices:              provisionInfo.Devices,
		CharmBaseImagePath:   charmBaseImage,
		Containers:           containers,
		CharmModifiedVersion: provisionInfo.CharmModifiedVersion,
		Trust:                provisionInfo.Trust,
		InitialScale:         provisionInfo.Scale,
	}
	switch ch.Meta().CharmUser {
	case charm.RunAsDefault:
		config.CharmUser = caas.RunAsDefault
	case charm.RunAsRoot:
		config.CharmUser = caas.RunAsRoot
	case charm.RunAsSudoer:
		config.CharmUser = caas.RunAsSudoer
	case charm.RunAsNonRoot:
		config.CharmUser = caas.RunAsNonRoot
	default:
		return errors.NotValidf("unknown RunAs for CharmUser: %q", ch.Meta().CharmUser)
	}
	reason := "unchanged"
	// TODO(sidecar): implement Equals method for caas.ApplicationConfig
	if !reflect.DeepEqual(config, *lastApplied) {
		if err = app.Ensure(config); err != nil {
			_ = setApplicationStatus(ctx, appName, status.Error, err.Error(), nil, facade, logger)
			return errors.Annotatef(err, "ensuring application %q", appName)
		}
		*lastApplied = config
		reason = "deployed"
		if appState.Exists {
			reason = "updated"
		}
	}
	logger.Debugf("application %q was %q", appName, reason)
	return nil
}

// appDying handles the life.Dying state for the CAAS application. It deals with scaling down
// the application and removing units.
func appDying(
	ctx context.Context,
	appName string, app caas.Application, appLife life.Value,
	facade CAASProvisionerFacade, unitFacade CAASUnitProvisionerFacade,
	logger logger.Logger,
) error {
	logger.Debugf("application %q dying", appName)
	err := ensureScale(ctx, appName, app, appLife, facade, unitFacade, logger)
	if err != nil {
		return errors.Annotate(err, "cannot scale dying application to 0")
	}
	err = reconcileDeadUnitScale(ctx, appName, app, facade, logger)
	if err != nil {
		return errors.Annotate(err, "cannot reconcile dead units in dying application")
	}
	return nil
}

// appDead handles the life.Dead state for the CAAS application. It ensures the application
// is removed from the k8s cluster and unblocks the cleanup of the application in state.
func appDead(
	ctx context.Context,
	appName string, app caas.Application,
	broker CAASBroker, facade CAASProvisionerFacade, unitFacade CAASUnitProvisionerFacade,
	clk clock.Clock, logger logger.Logger,
) error {
	logger.Debugf("application %q dead", appName)
	err := app.Delete()
	if err != nil {
		return errors.Trace(err)
	}
	err = waitForTerminated(appName, app, clk)
	if err != nil {
		return errors.Trace(err)
	}
	_, err = updateState(ctx, appName, app, nil, broker, facade, unitFacade, logger)
	if err != nil {
		return errors.Trace(err)
	}
	// Clear "has-resources" flag so state knows it can now remove the application.
	err = facade.ClearApplicationResources(ctx, appName)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// checkCharmFormat checks that the charm is a v2 charm.
func checkCharmFormat(
	ctx context.Context,
	appName string,
	facade CAASProvisionerFacade,
	logger logger.Logger,
) (isOk bool, err error) {
	charmInfo, err := facade.ApplicationCharmInfo(ctx, appName)
	if errors.Is(err, errors.NotFound) {
		logger.Debugf("application %q no longer exists", appName)
		return false, nil
	} else if err != nil {
		return false, errors.Annotatef(err, "failed to get charm info for application %q", appName)
	}
	format := charm.MetaFormat(charmInfo.Charm())
	if format >= charm.FormatV2 {
		return true, nil
	}
	return false, nil
}

// ensureTrust updates the applications Trust status on the CAAS broker, giving it
// access to the k8s api via a service account.
func ensureTrust(
	ctx context.Context,
	appName string, app caas.Application,
	unitFacade CAASUnitProvisionerFacade,
	logger logger.Logger,
) error {
	desiredTrust, err := unitFacade.ApplicationTrust(ctx, appName)
	if err != nil {
		return errors.Annotatef(err, "fetching application %q desired trust", appName)
	}

	logger.Debugf("updating application %q trust to %v", appName, desiredTrust)
	err = app.Trust(desiredTrust)
	if err != nil {
		return errors.Annotatef(
			err,
			"updating application %q to desired trust %v",
			appName,
			desiredTrust)
	}
	return nil
}

// updateState reports back information about the CAAS application into state, such as
// status, IP addresses and volume info.
func updateState(
	ctx context.Context,
	appName string, app caas.Application, lastReportedStatus map[string]status.StatusInfo,
	broker CAASBroker, facade CAASProvisionerFacade, unitFacade CAASUnitProvisionerFacade,
	logger logger.Logger,
) (map[string]status.StatusInfo, error) {
	appTag := names.NewApplicationTag(appName).String()
	appStatus := params.EntityStatus{}
	svc, err := app.Service()
	if err != nil && !errors.Is(err, errors.NotFound) {
		return nil, errors.Trace(err)
	}
	if svc != nil {
		appStatus = params.EntityStatus{
			Status: svc.Status.Status,
			Info:   svc.Status.Message,
			Data:   svc.Status.Data,
		}
		err = unitFacade.UpdateApplicationService(ctx, params.UpdateApplicationServiceArg{
			ApplicationTag: appTag,
			ProviderId:     svc.Id,
			Addresses:      params.FromProviderAddresses(svc.Addresses...),
		})
		if errors.Is(err, errors.NotFound) {
			// Do nothing
		} else if err != nil {
			return nil, errors.Trace(err)
		}
	}

	units, err := app.Units()
	if err != nil {
		return nil, errors.Trace(err)
	}

	reportedStatus := make(map[string]status.StatusInfo)
	args := params.UpdateApplicationUnits{
		ApplicationTag: appTag,
		Status:         appStatus,
	}
	for _, u := range units {
		// For pods managed by the substrate, any marked as dying
		// are treated as non-existing.
		if u.Dying {
			continue
		}
		unitStatus := u.Status
		lastStatus, ok := lastReportedStatus[u.Id]
		reportedStatus[u.Id] = unitStatus
		// TODO: Determine a better way to propagate status
		// without constantly overriding the juju state value.
		if ok {
			// If we've seen the same status value previously,
			// report as unknown as this value is ignored.
			if reflect.DeepEqual(lastStatus, unitStatus) {
				unitStatus = status.StatusInfo{
					Status: status.Unknown,
				}
			}
		}
		unitParams := params.ApplicationUnitParams{
			ProviderId: u.Id,
			Address:    u.Address,
			Ports:      u.Ports,
			Stateful:   u.Stateful,
			Status:     unitStatus.Status.String(),
			Info:       unitStatus.Message,
			Data:       unitStatus.Data,
		}
		// Fill in any filesystem info for volumes attached to the unit.
		// A unit will not become active until all required volumes are
		// provisioned, so it makes sense to send this information along
		// with the units to which they are attached.
		for _, info := range u.FilesystemInfo {
			unitParams.FilesystemInfo = append(unitParams.FilesystemInfo, params.KubernetesFilesystemInfo{
				StorageName:  info.StorageName,
				FilesystemId: info.FilesystemId,
				Size:         info.Size,
				MountPoint:   info.MountPoint,
				ReadOnly:     info.ReadOnly,
				Status:       info.Status.Status.String(),
				Info:         info.Status.Message,
				Data:         info.Status.Data,
				Volume: params.KubernetesVolumeInfo{
					VolumeId:   info.Volume.VolumeId,
					Size:       info.Volume.Size,
					Persistent: info.Volume.Persistent,
					Status:     info.Volume.Status.Status.String(),
					Info:       info.Volume.Status.Message,
					Data:       info.Volume.Status.Data,
				},
			})
		}
		args.Units = append(args.Units, unitParams)
	}

	appUnitInfo, err := facade.UpdateUnits(ctx, args)
	if err != nil {
		// We can ignore not found errors as the worker will get stopped anyway.
		// We can also ignore Forbidden errors raised from SetScale because disordered events could happen often.
		if !errors.Is(err, errors.Forbidden) && !errors.Is(err, errors.NotFound) {
			return nil, errors.Trace(err)
		}
		logger.Warningf("update units %v", err)
	}

	if appUnitInfo != nil {
		for _, unitInfo := range appUnitInfo.Units {
			unit, err := names.ParseUnitTag(unitInfo.UnitTag)
			if err != nil {
				return nil, errors.Trace(err)
			}
			err = broker.AnnotateUnit(context.TODO(), appName, unitInfo.ProviderId, unit)
			if errors.Is(err, errors.NotFound) {
				continue
			} else if err != nil {
				return nil, errors.Trace(err)
			}
		}
	}
	return reportedStatus, nil
}

func refreshApplicationStatus(
	ctx context.Context,
	appName string, app caas.Application, appLife life.Value,
	facade CAASProvisionerFacade,
	logger logger.Logger,
) error {
	if appLife != life.Alive {
		return nil
	}
	st, err := app.State()
	if errors.Is(err, errors.NotFound) {
		// Do nothing.
		return nil
	} else if err != nil {
		return errors.Trace(err)
	}

	// refresh the units information.
	units, err := facade.Units(ctx, appName)
	if errors.Is(err, errors.NotFound) {
		return nil
	} else if err != nil {
		return errors.Trace(err)
	}
	readyUnitsCount := 0
	for _, unit := range units {
		if unit.UnitStatus.AgentStatus.Status == string(status.Active) {
			readyUnitsCount++
		}
	}
	if st.DesiredReplicas > 0 && st.DesiredReplicas > readyUnitsCount {
		// Only set status to waiting for scale up.
		// When the application gets scaled down, the desired units will be kept running and
		// the application should be active always.
		return setApplicationStatus(ctx, appName, status.Waiting, "waiting for units to settle down", nil, facade, logger)
	}
	return setApplicationStatus(ctx, appName, status.Active, "", nil, facade, logger)
}

func waitForTerminated(appName string, app caas.Application,
	clk clock.Clock) error {
	existsFunc := func() error {
		appState, err := app.Exists()
		if err != nil {
			return errors.Trace(err)
		}
		if !appState.Exists {
			return nil
		}
		if appState.Exists && !appState.Terminating {
			return errors.Errorf("application %q should be terminating but is now running", appName)
		}
		return tryAgain
	}
	retryCallArgs := retry.CallArgs{
		Attempts:    60,
		Delay:       3 * time.Second,
		MaxDuration: 3 * time.Minute,
		Clock:       clk,
		Func:        existsFunc,
		IsFatalError: func(err error) bool {
			return !errors.Is(err, tryAgain)
		},
	}
	return errors.Trace(retry.Call(retryCallArgs))
}

// reconcileDeadUnitScale is setup to respond to CAAS sidecar units that become
// dead. It takes stock of what the current desired scale is for the application
// and the number of dead units in the application. Once the number of dead units
// has reached the a point where the desired scale has been achieved this func
// can go ahead and removed the units from CAAS provider.
func reconcileDeadUnitScale(
	ctx context.Context,
	appName string, app caas.Application,
	facade CAASProvisionerFacade,
	logger logger.Logger,
) error {
	units, err := facade.Units(ctx, appName)
	if err != nil {
		return fmt.Errorf("getting units for application %s: %w", appName, err)
	}

	ps, err := facade.ProvisioningState(ctx, appName)
	if err != nil {
		return errors.Trace(err)
	}
	if ps == nil || !ps.Scaling {
		return nil
	}

	desiredScale := ps.ScaleTarget
	unitsToRemove := len(units) - desiredScale

	var deadUnits []params.CAASUnit
	for _, unit := range units {
		unitLife, err := facade.Life(ctx, unit.Tag.Id())
		if err != nil {
			return fmt.Errorf("getting life for unit %q: %w", unit.Tag, err)
		}
		if unitLife == life.Dead {
			deadUnits = append(deadUnits, unit)
		}
	}

	if unitsToRemove <= 0 {
		unitsToRemove = len(deadUnits)
	}

	// We haven't met the threshold to initiate scale down in the CAAS provider
	// yet.
	if unitsToRemove != len(deadUnits) {
		return nil
	}

	logger.Infof("scaling application %q to desired scale %d", appName, desiredScale)
	if err := app.Scale(desiredScale); err != nil && !errors.Is(err, errors.NotFound) {
		return fmt.Errorf(
			"scaling application %q to scale %d: %w",
			appName,
			desiredScale,
			err,
		)
	}

	appState, err := app.State()
	if err != nil && !errors.Is(err, errors.NotFound) {
		return err
	}
	// TODO: stop k8s things from mutating the statefulset.
	if len(appState.Replicas) > desiredScale {
		return tryAgain
	}

	for _, deadUnit := range deadUnits {
		logger.Infof("removing dead unit %s", deadUnit.Tag.Id())
		if err := facade.RemoveUnit(ctx, deadUnit.Tag.Id()); err != nil && !errors.Is(err, errors.NotFound) {
			return fmt.Errorf("removing dead unit %q: %w", deadUnit.Tag.Id(), err)
		}
	}

	return updateProvisioningState(ctx, appName, false, 0, facade)
}

// ensureScale determines how and when to scale up or down based on
// current scale targets that have yet to be met.
func ensureScale(
	ctx context.Context,
	appName string, app caas.Application, appLife life.Value,
	facade CAASProvisionerFacade, unitFacade CAASUnitProvisionerFacade,
	logger logger.Logger,
) error {
	var err error
	var desiredScale int
	switch appLife {
	case life.Alive:
		desiredScale, err = unitFacade.ApplicationScale(ctx, appName)
		if err != nil {
			return errors.Annotatef(err, "fetching application %q desired scale", appName)
		}
	case life.Dying, life.Dead:
		desiredScale = 0
	default:
		return errors.NotImplementedf("unknown life %q", appLife)
	}

	ps, err := facade.ProvisioningState(ctx, appName)
	if err != nil {
		return errors.Trace(err)
	}

	logger.Debugf("updating application %q scale to %d", appName, desiredScale)
	if !ps.Scaling || appLife != life.Alive {
		err := updateProvisioningState(ctx, appName, true, desiredScale, facade)
		if err != nil {
			return err
		}
		ps.Scaling = true
		ps.ScaleTarget = desiredScale
	}

	units, err := facade.Units(ctx, appName)
	if err != nil {
		return err
	}
	if ps.ScaleTarget >= len(units) {
		logger.Infof("scaling application %q to desired scale %d", appName, ps.ScaleTarget)
		err = app.Scale(ps.ScaleTarget)
		if appLife != life.Alive && errors.Is(err, errors.NotFound) {
			logger.Infof("dying application %q is already removed", appName)
		} else if err != nil {
			return err
		}
		return updateProvisioningState(ctx, appName, false, 0, facade)
	}

	unitsToDestroy, err := app.UnitsToRemove(context.TODO(), ps.ScaleTarget)
	if err != nil && errors.Is(err, errors.NotFound) {
		return nil
	} else if err != nil {
		return fmt.Errorf("scaling application %q to desired scale %d: %w",
			appName, ps.ScaleTarget, err)
	}

	if len(unitsToDestroy) > 0 {
		if err := facade.DestroyUnits(ctx, unitsToDestroy); err != nil {
			return errors.Trace(err)
		}
	}

	if ps.ScaleTarget != desiredScale {
		// if the current scale target doesn't equal the desired scale
		// we need to rerun this.
		logger.Debugf("application %q currently scaling to %d but desired scale is %d", appName, ps.ScaleTarget, desiredScale)
		return tryAgain
	}

	return nil
}

func setApplicationStatus(
	ctx context.Context,
	appName string, s status.Status, reason string, data map[string]interface{},
	facade CAASProvisionerFacade,
	logger logger.Logger,
) error {
	logger.Tracef("updating application %q status to %q, %q, %v", appName, s, reason, data)
	return facade.SetOperatorStatus(ctx, appName, s, reason, data)
}

func updateProvisioningState(
	ctx context.Context,
	appName string, scaling bool, scaleTarget int,
	facade CAASProvisionerFacade,
) error {
	newPs := params.CAASApplicationProvisioningState{
		Scaling:     scaling,
		ScaleTarget: scaleTarget,
	}
	err := facade.SetProvisioningState(ctx, appName, newPs)
	if params.IsCodeTryAgain(err) {
		return tryAgain
	} else if err != nil {
		return errors.Annotatef(err, "setting provisiong state for application %q", appName)
	}
	return nil
}
