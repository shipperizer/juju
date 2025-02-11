// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application_test

import (
	"context"

	"github.com/juju/clock"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/arch"
	corecharm "github.com/juju/juju/core/charm"
	"github.com/juju/juju/core/database"
	corestorage "github.com/juju/juju/core/storage"
	"github.com/juju/juju/domain/application/architecture"
	"github.com/juju/juju/domain/application/charm"
	"github.com/juju/juju/domain/application/service"
	"github.com/juju/juju/domain/application/state"
	domaintesting "github.com/juju/juju/domain/schema/testing"
	secretstate "github.com/juju/juju/domain/secret/state"
	internalcharm "github.com/juju/juju/internal/charm"
	loggertesting "github.com/juju/juju/internal/logger/testing"
	"github.com/juju/juju/internal/storage"
	"github.com/juju/juju/internal/storage/provider"
)

type charmSuite struct {
	domaintesting.ModelSuite
}

var _ = gc.Suite(&charmSuite{})

func (s *charmSuite) TestSetCharmWithArchitecture(c *gc.C) {
	service := s.setupService(c)

	// We can't use the architecture from the manifest, as there may not be one.
	// So we should default to Unknown.

	metadata := internalcharm.Meta{
		Name: "foo",
	}
	manifest := internalcharm.Manifest{
		Bases: []internalcharm.Base{{
			Name: "ubuntu",
			Channel: internalcharm.Channel{
				Risk: internalcharm.Stable,
			},
			Architectures: []string{"amd64"},
		}},
	}

	id, _, err := service.SetCharm(context.Background(), charm.SetCharmArgs{
		Charm:         internalcharm.NewCharmBase(&metadata, &manifest, nil, nil, nil),
		Source:        corecharm.Local,
		ReferenceName: "foo",
		Revision:      1,
		Hash:          "hash",
		ArchivePath:   "archive",
		Version:       "1.0",
		Architecture:  arch.ARM64,
	})
	c.Assert(err, jc.ErrorIsNil)

	_, locator, err := service.GetCharm(context.Background(), id)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(locator.Architecture, gc.Equals, architecture.ARM64)
}

func (s *charmSuite) TestSetCharmWithoutArchitecture(c *gc.C) {
	service := s.setupService(c)

	// We can't use the architecture from the manifest, as there may not be one.
	// So we should default to Unknown.

	metadata := internalcharm.Meta{
		Name: "foo",
	}
	manifest := internalcharm.Manifest{
		Bases: []internalcharm.Base{{
			Name: "ubuntu",
			Channel: internalcharm.Channel{
				Risk: internalcharm.Stable,
			},
			Architectures: []string{"amd64"},
		}},
	}

	id, _, err := service.SetCharm(context.Background(), charm.SetCharmArgs{
		Charm:         internalcharm.NewCharmBase(&metadata, &manifest, nil, nil, nil),
		Source:        corecharm.Local,
		ReferenceName: "foo",
		Revision:      1,
		Hash:          "hash",
		ArchivePath:   "archive",
		Version:       "1.0",
	})
	c.Assert(err, jc.ErrorIsNil)

	_, locator, err := service.GetCharm(context.Background(), id)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(locator.Architecture, gc.Equals, architecture.Unknown)
}

func (s *charmSuite) setupService(c *gc.C) *service.Service {
	modelDB := func() (database.TxnRunner, error) {
		return s.ModelTxnRunner(), nil
	}

	return service.NewService(
		state.NewState(modelDB, clock.WallClock, loggertesting.WrapCheckLog(c)),
		secretstate.NewState(modelDB, loggertesting.WrapCheckLog(c)),
		corestorage.ConstModelStorageRegistry(func() storage.ProviderRegistry {
			return provider.CommonStorageProviders()
		}),
		nil,
		clock.WallClock,
		loggertesting.WrapCheckLog(c),
	)
}
