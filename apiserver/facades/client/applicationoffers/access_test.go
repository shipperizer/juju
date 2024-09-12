// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package applicationoffers_test

import (
	"context"
	"regexp"

	"github.com/go-macaroon-bakery/macaroon-bakery/v3/bakery"
	"github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/names/v5"
	jc "github.com/juju/testing/checkers"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/apiserver/common/crossmodel"
	"github.com/juju/juju/apiserver/facades/client/applicationoffers"
	jujucrossmodel "github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/core/model"
	modeltesting "github.com/juju/juju/core/model/testing"
	"github.com/juju/juju/core/permission"
	coreuser "github.com/juju/juju/core/user"
	access "github.com/juju/juju/domain/access"
	accesserrors "github.com/juju/juju/domain/access/errors"
	loggertesting "github.com/juju/juju/internal/logger/testing"
	"github.com/juju/juju/internal/testing"
	"github.com/juju/juju/internal/uuid"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
)

type offerAccessSuite struct {
	baseSuite
	api *applicationoffers.OffersAPIv5
}

var _ = gc.Suite(&offerAccessSuite{})

func (s *offerAccessSuite) SetUpTest(c *gc.C) {
	s.baseSuite.SetUpTest(c)
	s.authorizer.Tag = names.NewUserTag("admin")

	modelID := modeltesting.GenModelUUID(c)
	var err error
	thirdPartyKey := bakery.MustGenerateKey()
	s.authContext, err = crossmodel.NewAuthContext(
		s.mockState, nil, names.NewModelTag(modelID.String()), thirdPartyKey,
		crossmodel.NewOfferBakeryForTest(s.bakery, clock.WallClock),
	)
	c.Assert(err, jc.ErrorIsNil)
}

// Creates the API to use in testing.
// Call baseSuite.setupMocks before this.
func (s *offerAccessSuite) setupAPI(c *gc.C) {
	getApplicationOffers := func(interface{}) jujucrossmodel.ApplicationOffers {
		return &stubApplicationOffers{}
	}
	api, err := applicationoffers.CreateOffersAPI(
		getApplicationOffers, getFakeControllerInfo,
		s.mockState, s.mockStatePool,
		s.mockAccessService,
		s.mockModelServiceFactoryGetter,
		s.authorizer, s.authContext,
		c.MkDir(), loggertesting.WrapCheckLog(c),
		testing.ControllerTag.Id(), model.UUID(testing.ModelTag.Id()),
	)
	c.Assert(err, jc.ErrorIsNil)
	s.api = api
}

func (s *offerAccessSuite) modifyAccess(
	c *gc.C, user names.UserTag,
	action params.OfferAction,
	access params.OfferAccessPermission,
	offerURL string,
) error {
	args := params.ModifyOfferAccessRequest{
		Changes: []params.ModifyOfferAccess{{
			UserTag:  user.String(),
			Action:   action,
			Access:   access,
			OfferURL: offerURL,
		}}}

	result, err := s.api.ModifyOfferAccess(context.Background(), args)
	if err != nil {
		return err
	}
	return result.OneError()
}

func (s *offerAccessSuite) grant(c *gc.C, user names.UserTag, access params.OfferAccessPermission, offerURL string) error {
	return s.modifyAccess(c, user, params.GrantOfferAccess, access, offerURL)
}

func (s *offerAccessSuite) revoke(c *gc.C, user names.UserTag, access params.OfferAccessPermission, offerURL string) error {
	return s.modifyAccess(c, user, params.RevokeOfferAccess, access, offerURL)
}

func (s *offerAccessSuite) setupOffer(modelUUID, modelName, owner, offerName string) string {
	model := &mockModel{uuid: modelUUID, name: modelName, owner: owner, modelType: state.ModelTypeIAAS}
	s.mockState.allmodels = []applicationoffers.Model{model}
	st := &mockState{
		modelUUID:         modelUUID,
		applicationOffers: make(map[string]jujucrossmodel.ApplicationOffer),
		model:             model,
	}
	s.mockStatePool.st[modelUUID] = st
	uuid := uuid.MustNewUUID().String()
	st.applicationOffers[offerName] = jujucrossmodel.ApplicationOffer{OfferUUID: uuid}
	return uuid
}

func (s *offerAccessSuite) TestGrantMissingUserFails(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.setupAPI(c)

	offerUUID := s.setupOffer("uuid", "test", "admin", "someoffer")
	user := names.NewUserTag("foobar")

	s.mockAccessService.EXPECT().UpdatePermission(gomock.Any(), access.UpdatePermissionArgs{
		AccessSpec: offerAccessSpec(offerUUID, permission.ReadAccess),
		Subject:    coreuser.NameFromTag(user),
		Change:     permission.Grant,
	}).Return(accesserrors.UserNotFound)

	err := s.grant(c, user, params.OfferReadAccess, "test.someoffer")
	expectedErr := `could not grant offer access for "foobar": user not found`
	c.Assert(err, gc.ErrorMatches, expectedErr)
}

func (s *offerAccessSuite) TestGrantMissingOfferFails(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.setupAPI(c)

	s.setupOffer("uuid", "test", "admin", "differentoffer")
	user := names.NewUserTag("foobar")
	err := s.grant(c, user, params.OfferReadAccess, "test.someoffer")
	expectedErr := `.*application offer "someoffer" not found`
	c.Assert(err, gc.ErrorMatches, expectedErr)
}

func (s *offerAccessSuite) TestRevokePermission(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.setupAPI(c)

	offerUUID := s.setupOffer("uuid", "test", "admin", "someoffer")
	user := names.NewUserTag("foobar")
	userName := coreuser.NameFromTag(user)
	s.mockAccessService.EXPECT().UpdatePermission(gomock.Any(), access.UpdatePermissionArgs{
		AccessSpec: offerAccessSpec(offerUUID, permission.ReadAccess),
		Subject:    userName,
		Change:     permission.Revoke,
	})

	err := s.revoke(c, user, params.OfferReadAccess, "test.someoffer")
	c.Assert(err, gc.IsNil)
}

func (s *offerAccessSuite) TestGrantPermission(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.setupAPI(c)

	offerUUID := s.setupOffer("uuid", "test", "admin", "someoffer")

	user := names.NewUserTag("foobar")
	userName := coreuser.NameFromTag(user)
	s.mockAccessService.EXPECT().UpdatePermission(gomock.Any(), access.UpdatePermissionArgs{
		AccessSpec: offerAccessSpec(offerUUID, permission.ReadAccess),
		Subject:    userName,
		Change:     permission.Grant,
	}).Return(accesserrors.PermissionAccessGreater)

	err := s.grant(c, user, params.OfferReadAccess, "test.someoffer")

	c.Assert(errors.Cause(err), gc.ErrorMatches, `could not grant offer access for .*: access or greater`)
}

func (s *offerAccessSuite) TestGrantPermissionAddRemoteUser(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.setupAPI(c)

	offerUUID := s.setupOffer("uuid", "test", "superuser-bob", "someoffer")

	apiUser := names.NewUserTag("superuser-bob")
	s.authorizer.Tag = apiUser
	user := names.NewUserTag("bob@remote")
	userName := coreuser.NameFromTag(user)

	s.mockAccessService.EXPECT().UpdatePermission(gomock.Any(), access.UpdatePermissionArgs{
		AccessSpec: offerAccessSpec(offerUUID, permission.ReadAccess),
		Subject:    userName,
		Change:     permission.Grant,
	})

	err := s.grant(c, user, params.OfferReadAccess, "superuser-bob/test.someoffer")
	c.Assert(err, jc.ErrorIsNil)
}

func (s *offerAccessSuite) assertGrantToOffer(c *gc.C, userAccess permission.Access) {
	offerUUID := s.setupOffer("uuid", "test", "bob@remote", "someoffer")

	user := names.NewUserTag("bob@remote")
	s.authorizer.Tag = user
	other := names.NewUserTag("other@remote")

	s.mockAccessService.EXPECT().ReadUserAccessLevelForTarget(gomock.Any(), coreuser.NameFromTag(user), permission.ID{
		ObjectType: permission.Offer,
		Key:        offerUUID,
	}).Return(userAccess, nil)

	err := s.grant(c, other, params.OfferReadAccess, "bob@remote/test.someoffer")
	c.Assert(err, gc.ErrorMatches, "permission denied")
}

func (s *offerAccessSuite) TestGrantToOfferNoAccess(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.setupAPI(c)

	s.assertGrantToOffer(c, permission.NoAccess)
}

func (s *offerAccessSuite) TestGrantToOfferReadAccess(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.setupAPI(c)

	s.assertGrantToOffer(c, permission.ReadAccess)
}

func (s *offerAccessSuite) TestGrantToOfferConsumeAccess(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.setupAPI(c)

	s.assertGrantToOffer(c, permission.ConsumeAccess)
}

func (s *offerAccessSuite) TestGrantToOfferAdminAccess(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.setupAPI(c)

	offerUUID := s.setupOffer("uuid", "test", "foobar", "someoffer")

	user := names.NewUserTag("foobar")
	s.authorizer.Tag = user
	other := names.NewUserTag("other")
	otherName := coreuser.NameFromTag(other)

	s.mockAccessService.EXPECT().ReadUserAccessLevelForTarget(gomock.Any(), coreuser.NameFromTag(user), permission.ID{
		ObjectType: permission.Offer,
		Key:        offerUUID,
	}).Return(permission.AdminAccess, nil)

	s.mockAccessService.EXPECT().UpdatePermission(gomock.Any(), access.UpdatePermissionArgs{
		AccessSpec: offerAccessSpec(offerUUID, permission.ReadAccess),
		Subject:    otherName,
		Change:     permission.Grant,
	})

	err := s.grant(c, other, params.OfferReadAccess, "foobar/test.someoffer")
	c.Assert(err, jc.ErrorIsNil)
}

func (s *offerAccessSuite) TestGrantOfferInvalidUserTag(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.setupAPI(c)

	s.setupOffer("uuid", "test", "admin", "someoffer")
	for _, testParam := range []struct {
		tag      string
		validTag bool
	}{{
		tag:      "unit-foo/0",
		validTag: true,
	}, {
		tag:      "application-foo",
		validTag: true,
	}, {
		tag:      "relation-wordpress:db mysql:db",
		validTag: true,
	}, {
		tag:      "machine-0",
		validTag: true,
	}, {
		tag:      "user",
		validTag: false,
	}, {
		tag:      "user-Mua^h^h^h^arh",
		validTag: true,
	}, {
		tag:      "user@",
		validTag: false,
	}, {
		tag:      "user@ubuntuone",
		validTag: false,
	}, {
		tag:      "user@ubuntuone",
		validTag: false,
	}, {
		tag:      "@ubuntuone",
		validTag: false,
	}, {
		tag:      "in^valid.",
		validTag: false,
	}, {
		tag:      "",
		validTag: false,
	},
	} {
		var expectedErr string
		errPart := `could not modify offer access: "` + regexp.QuoteMeta(testParam.tag) + `" is not a valid `

		if testParam.validTag {
			// The string is a valid tag, but not a user tag.
			expectedErr = errPart + `user tag`
		} else {
			// The string is not a valid tag of any kind.
			expectedErr = errPart + `tag`
		}

		args := params.ModifyOfferAccessRequest{
			Changes: []params.ModifyOfferAccess{{
				UserTag:  testParam.tag,
				Action:   params.GrantOfferAccess,
				Access:   params.OfferReadAccess,
				OfferURL: "test.someoffer",
			}}}

		result, err := s.api.ModifyOfferAccess(context.Background(), args)
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(result.OneError(), gc.ErrorMatches, expectedErr)
	}
}

func (s *offerAccessSuite) TestModifyOfferAccessEmptyArgs(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.setupAPI(c)

	s.setupOffer("uuid", "test", "admin", "someoffer")
	args := params.ModifyOfferAccessRequest{
		Changes: []params.ModifyOfferAccess{{OfferURL: "test.someoffer"}}}

	result, err := s.api.ModifyOfferAccess(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)
	expectedErr := `could not modify offer access: "" offer access not valid`
	c.Assert(result.OneError(), gc.ErrorMatches, expectedErr)
}

func (s *offerAccessSuite) TestModifyOfferAccessInvalidAction(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.setupAPI(c)

	s.setupOffer("uuid", "test", "admin", "someoffer")

	var dance params.OfferAction = "dance"
	args := params.ModifyOfferAccessRequest{
		Changes: []params.ModifyOfferAccess{{
			UserTag:  "user-user",
			Action:   dance,
			Access:   params.OfferReadAccess,
			OfferURL: "test.someoffer",
		}}}

	result, err := s.api.ModifyOfferAccess(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)
	expectedErr := `unknown action "dance"`
	c.Assert(result.OneError(), gc.ErrorMatches, expectedErr)
}

func offerAccessSpec(offerUUID string, accessLevel permission.Access) permission.AccessSpec {
	return permission.AccessSpec{
		Target: permission.ID{
			ObjectType: permission.Offer,
			Key:        offerUUID,
		},
		Access: accessLevel,
	}
}
