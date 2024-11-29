// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package secret_test

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/juju/clock"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/changestream"
	corecharm "github.com/juju/juju/core/charm"
	coresecrets "github.com/juju/juju/core/secrets"
	corestorage "github.com/juju/juju/core/storage"
	"github.com/juju/juju/core/unit"
	jujuversion "github.com/juju/juju/core/version"
	corewatcher "github.com/juju/juju/core/watcher"
	"github.com/juju/juju/core/watcher/watchertest"
	"github.com/juju/juju/domain"
	applicationservice "github.com/juju/juju/domain/application/service"
	applicationstate "github.com/juju/juju/domain/application/state"
	"github.com/juju/juju/domain/secret"
	"github.com/juju/juju/domain/secret/service"
	"github.com/juju/juju/domain/secret/state"
	"github.com/juju/juju/internal/changestream/testing"
	"github.com/juju/juju/internal/charm"
	loggertesting "github.com/juju/juju/internal/logger/testing"
	"github.com/juju/juju/internal/storage"
	coretesting "github.com/juju/juju/internal/testing"
	"github.com/juju/juju/internal/uuid"
)

type watcherSuite struct {
	testing.ModelSuite
}

var _ = gc.Suite(&watcherSuite{})

func (s *watcherSuite) SetUpTest(c *gc.C) {
	s.ModelSuite.SetUpTest(c)

	err := s.TxnRunner().StdTxn(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
INSERT INTO model (uuid, controller_uuid, target_agent_version, name, type, cloud, cloud_type)
VALUES (?, ?, ?, "test", "iaas", "fluffy", "ec2")
		`, s.ModelUUID(), coretesting.ControllerTag.Id(), jujuversion.Current.String())
		return err
	})
	c.Assert(err, jc.ErrorIsNil)
}

func (s *watcherSuite) setupUnits(c *gc.C, appName string) {
	logger := loggertesting.WrapCheckLog(c)
	st := applicationstate.NewState(s.TxnRunnerFactory(), logger)
	svc := applicationservice.NewService(st, nil,
		corestorage.ConstModelStorageRegistry(func() storage.ProviderRegistry {
			return storage.NotImplementedProviderRegistry{}
		}),
		nil,
		clock.WallClock,
		logger,
	)

	unitName, err := unit.NewNameFromParts(appName, 0)
	c.Assert(err, jc.ErrorIsNil)
	_, err = svc.CreateApplication(context.Background(),
		appName,
		&stubCharm{},
		corecharm.Origin{
			Source: corecharm.CharmHub,
			Platform: corecharm.Platform{
				Channel:      "24.04",
				OS:           "ubuntu",
				Architecture: "amd64",
			},
		},
		applicationservice.AddApplicationArgs{
			ReferenceName: appName,
		},
		applicationservice.AddUnitArg{UnitName: unitName},
	)
	c.Assert(err, jc.ErrorIsNil)
}

func revID(uri *coresecrets.URI, rev int) string {
	return fmt.Sprintf("%s/%d", uri.ID, rev)
}

func createNewRevision(c *gc.C, st *state.State, uri *coresecrets.URI) {
	sp := secret.UpsertSecretParams{
		Data:       coresecrets.SecretData{"foo-new": "bar-new"},
		RevisionID: ptr(uuid.MustNewUUID().String()),
	}
	err := st.RunAtomic(context.Background(), func(ctx domain.AtomicContext) error {
		return st.UpdateSecret(ctx, uri, sp)
	})
	c.Assert(err, jc.ErrorIsNil)
}

func (s *watcherSuite) setupServiceAndState(c *gc.C) (*service.WatchableService, *state.State) {
	logger := loggertesting.WrapCheckLog(c)
	st := state.NewState(s.TxnRunnerFactory(), logger)
	factory := domain.NewWatcherFactory(
		changestream.NewWatchableDBFactoryForNamespace(s.GetWatchableDB, "secret_revision"),
		logger,
	)
	return service.NewWatchableService(st, nil, logger, factory, service.SecretServiceParams{}), st
}

func createUserSecret(ctx context.Context, st *state.State, version int, uri *coresecrets.URI, secret secret.UpsertSecretParams) error {
	return st.RunAtomic(ctx, func(ctx domain.AtomicContext) error {
		return st.CreateUserSecret(ctx, version, uri, secret)
	})
}

func createCharmApplicationSecret(ctx context.Context, st *state.State, version int, uri *coresecrets.URI, appName string, secret secret.UpsertSecretParams) error {
	return st.RunAtomic(ctx, func(ctx domain.AtomicContext) error {
		appUUID, err := st.GetApplicationUUID(ctx, appName)
		if err != nil {
			return err
		}
		return st.CreateCharmApplicationSecret(ctx, version, uri, appUUID, secret)
	})
}

func createCharmUnitSecret(ctx context.Context, st *state.State, version int, uri *coresecrets.URI, unitName string, secret secret.UpsertSecretParams) error {
	return st.RunAtomic(ctx, func(ctx domain.AtomicContext) error {
		unitUUID, err := st.GetUnitUUID(ctx, unitName)
		if err != nil {
			return err
		}
		return st.CreateCharmUnitSecret(ctx, version, uri, unitUUID, secret)
	})
}

func (s *watcherSuite) TestWatchObsoleteForAppsAndUnitsOwned(c *gc.C) {
	s.setupUnits(c, "mysql")
	s.setupUnits(c, "mediawiki")

	ctx := context.Background()
	svc, st := s.setupServiceAndState(c)

	uri1 := coresecrets.NewURI()
	uri2 := coresecrets.NewURI()
	uri3 := coresecrets.NewURI()
	uri4 := coresecrets.NewURI()

	w, err := svc.WatchObsolete(ctx,
		service.CharmSecretOwner{
			Kind: service.ApplicationOwner,
			ID:   "mysql",
		},
		service.CharmSecretOwner{
			Kind: service.UnitOwner,
			ID:   "mysql/0",
		},

		service.CharmSecretOwner{
			Kind: service.ApplicationOwner,
			ID:   "mediawiki",
		},
		service.CharmSecretOwner{
			Kind: service.UnitOwner,
			ID:   "mediawiki/0",
		},
	)
	c.Assert(err, gc.IsNil)
	c.Assert(w, gc.NotNil)
	defer watchertest.CleanKill(c, w)

	harness := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w))
	harness.AddTest(func(c *gc.C) {
		sp := secret.UpsertSecretParams{
			Data: coresecrets.SecretData{"foo": "bar", "hello": "world"},
		}
		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err := createCharmApplicationSecret(ctx, st, 1, uri1, "mysql", sp)
		c.Assert(err, jc.ErrorIsNil)

		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err = createCharmUnitSecret(ctx, st, 1, uri2, "mysql/0", sp)
		c.Assert(err, jc.ErrorIsNil)

		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err = createCharmApplicationSecret(ctx, st, 1, uri3, "mediawiki", sp)
		c.Assert(err, jc.ErrorIsNil)

		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err = createCharmUnitSecret(ctx, st, 1, uri4, "mediawiki/0", sp)
		c.Assert(err, jc.ErrorIsNil)
	}, func(w watchertest.WatcherC[[]string]) {
		w.AssertNoChange()
	})

	harness.AddTest(func(c *gc.C) {
		// create revision 2, and obsolete revision 1.
		createNewRevision(c, st, uri1)
		createNewRevision(c, st, uri2)
		createNewRevision(c, st, uri3)
		createNewRevision(c, st, uri4)
	}, func(w watchertest.WatcherC[[]string]) {
		w.Check(
			watchertest.StringSliceAssert(
				revID(uri1, 1),
				revID(uri2, 1),
				revID(uri3, 1),
				revID(uri4, 1),
			),
		)
	})

	harness.AddTest(func(c *gc.C) {
		// create revision 3, and obsolete revision 2.
		createNewRevision(c, st, uri1)
		createNewRevision(c, st, uri2)
		createNewRevision(c, st, uri3)
	}, func(w watchertest.WatcherC[[]string]) {
		w.Check(
			watchertest.StringSliceAssert(
				revID(uri1, 2),
				revID(uri2, 2),
				revID(uri3, 2),
			),
		)
	})

	harness.Run(c, []string(nil))
}

func (s *watcherSuite) TestWatchObsoleteForAppsOwned(c *gc.C) {
	s.setupUnits(c, "mysql")

	ctx := context.Background()
	svc, st := s.setupServiceAndState(c)

	uri1 := coresecrets.NewURI()
	uri2 := coresecrets.NewURI()

	w, err := svc.WatchObsolete(ctx,
		service.CharmSecretOwner{
			Kind: service.ApplicationOwner,
			ID:   "mysql",
		},
	)
	c.Assert(err, gc.IsNil)
	c.Assert(w, gc.NotNil)
	defer watchertest.CleanKill(c, w)

	harness := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w))
	harness.AddTest(func(c *gc.C) {
		sp := secret.UpsertSecretParams{
			Data: coresecrets.SecretData{"foo": "bar", "hello": "world"},
		}
		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err := createCharmApplicationSecret(ctx, st, 1, uri1, "mysql", sp)
		c.Assert(err, jc.ErrorIsNil)

		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err = createCharmUnitSecret(ctx, st, 1, uri2, "mysql/0", sp)
		c.Assert(err, jc.ErrorIsNil)
	}, func(w watchertest.WatcherC[[]string]) {
		w.AssertNoChange()
	})

	harness.AddTest(func(c *gc.C) {
		// create revision 2, and obsolete revision 1.
		createNewRevision(c, st, uri1)
		createNewRevision(c, st, uri2)
	}, func(w watchertest.WatcherC[[]string]) {
		w.Check(
			watchertest.StringSliceAssert(
				revID(uri1, 1),
			),
		)
	})

	harness.AddTest(func(c *gc.C) {
		// create revision 3, and obsolete revision 2.
		createNewRevision(c, st, uri1)
		createNewRevision(c, st, uri2)
	}, func(w watchertest.WatcherC[[]string]) {
		w.Check(
			watchertest.StringSliceAssert(
				revID(uri1, 2),
			),
		)
	})
	harness.Run(c, []string(nil))
}

func (s *watcherSuite) TestWatchObsoleteForUnitsOwned(c *gc.C) {
	s.setupUnits(c, "mysql")

	ctx := context.Background()
	svc, st := s.setupServiceAndState(c)

	uri1 := coresecrets.NewURI()
	uri2 := coresecrets.NewURI()

	w, err := svc.WatchObsolete(ctx,
		service.CharmSecretOwner{
			Kind: service.UnitOwner,
			ID:   "mysql/0",
		},
	)
	c.Assert(err, gc.IsNil)
	c.Assert(w, gc.NotNil)
	defer watchertest.CleanKill(c, w)

	harness := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w))
	harness.AddTest(func(c *gc.C) {
		sp := secret.UpsertSecretParams{
			Data: coresecrets.SecretData{"foo": "bar", "hello": "world"},
		}
		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err := createCharmApplicationSecret(ctx, st, 1, uri1, "mysql", sp)
		c.Assert(err, jc.ErrorIsNil)

		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err = createCharmUnitSecret(ctx, st, 1, uri2, "mysql/0", sp)
		c.Assert(err, jc.ErrorIsNil)
	}, func(w watchertest.WatcherC[[]string]) {
		w.AssertNoChange()
	})

	harness.AddTest(func(c *gc.C) {
		// create revision 2, and obsolete revision 1.
		createNewRevision(c, st, uri1)
		createNewRevision(c, st, uri2)
	}, func(w watchertest.WatcherC[[]string]) {
		w.Check(
			watchertest.StringSliceAssert(
				revID(uri2, 1),
			),
		)
	})

	harness.AddTest(func(c *gc.C) {
		// create revision 3, and obsolete revision 2.
		createNewRevision(c, st, uri1)
		createNewRevision(c, st, uri2)
	}, func(w watchertest.WatcherC[[]string]) {
		w.Check(
			watchertest.StringSliceAssert(
				revID(uri2, 2),
			),
		)
	})
	harness.Run(c, []string(nil))
}

func (s *watcherSuite) TestWatchObsoleteUserSecretsToPrune(c *gc.C) {
	ctx := context.Background()
	svc, st := s.setupServiceAndState(c)

	data := coresecrets.SecretData{"foo": "bar", "hello": "world"}
	uri1 := coresecrets.NewURI()
	uri2 := coresecrets.NewURI()
	c.Logf("uri1: %v, uri2: %v", uri1, uri2)

	w, err := svc.WatchObsoleteUserSecretsToPrune(ctx)
	c.Assert(err, gc.IsNil)
	c.Assert(w, gc.NotNil)
	defer watchertest.CleanKill(c, w)

	harness := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w))
	harness.AddTest(func(c *gc.C) {
		err := createUserSecret(ctx, st, 1, uri1, secret.UpsertSecretParams{
			Data:       data,
			RevisionID: ptr(uuid.MustNewUUID().String()),
		})
		c.Assert(err, jc.ErrorIsNil)
		err = createUserSecret(ctx, st, 1, uri2, secret.UpsertSecretParams{
			Data:       data,
			AutoPrune:  ptr(true),
			RevisionID: ptr(uuid.MustNewUUID().String()),
		})
		c.Assert(err, jc.ErrorIsNil)
	}, func(w watchertest.WatcherC[struct{}]) {
		w.AssertNoChange()
	})

	harness.AddTest(func(c *gc.C) {
		// create revision 2, no event is fired because the auto prune is not turned on for uri1.
		createNewRevision(c, st, uri1)
	}, func(w watchertest.WatcherC[struct{}]) {
		w.AssertNoChange()
	})

	harness.AddTest(func(c *gc.C) {
		// create revision 2, and obsolete revision 1. An event is fired because the auto prune is turned on for uri2.
		createNewRevision(c, st, uri2)
	}, func(w watchertest.WatcherC[struct{}]) {
		w.AssertNChanges(2)
	})

	harness.AddTest(func(c *gc.C) {
		err = st.RunAtomic(context.Background(), func(ctx domain.AtomicContext) error {
			return st.UpdateSecret(ctx, uri1, secret.UpsertSecretParams{
				AutoPrune: ptr(true),
			})
		})
		c.Assert(err, jc.ErrorIsNil)
	}, func(w watchertest.WatcherC[struct{}]) {
		w.AssertChange()
	})

	harness.Run(c, struct{}{})

	// Pretend that the agent restarted and the watcher is re-created.
	w1, err := svc.WatchObsoleteUserSecretsToPrune(ctx)
	c.Assert(err, gc.IsNil)
	c.Assert(w1, gc.NotNil)
	defer watchertest.CleanKill(c, w1)

	harness1 := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w1))
	harness1.AddTest(func(c *gc.C) {}, func(w watchertest.WatcherC[struct{}]) {
		w.AssertNoChange()
	})

	harness1.Run(c, struct{}{})
}

func (s *watcherSuite) TestWatchConsumedSecretsChanges(c *gc.C) {
	s.setupUnits(c, "mysql")
	s.setupUnits(c, "mediawiki")

	ctx := context.Background()
	svc, st := s.setupServiceAndState(c)

	saveConsumer := func(uri *coresecrets.URI, revision int, consumerID string) {
		consumer := &coresecrets.SecretConsumerMetadata{
			CurrentRevision: revision,
		}
		err := st.SaveSecretConsumer(ctx, uri, consumerID, consumer)
		c.Assert(err, jc.ErrorIsNil)
	}

	uri1 := coresecrets.NewURI()
	uri2 := coresecrets.NewURI()

	w, err := svc.WatchConsumedSecretsChanges(ctx, "mediawiki/0")
	c.Assert(err, gc.IsNil)
	c.Assert(w, gc.NotNil)
	defer watchertest.CleanKill(c, w)

	harness := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w))
	harness.AddTest(func(c *gc.C) {
		sp := secret.UpsertSecretParams{
			Data: coresecrets.SecretData{"foo": "bar", "hello": "world"},
		}

		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err := createCharmApplicationSecret(ctx, st, 1, uri1, "mysql", sp)
		c.Assert(err, jc.ErrorIsNil)

		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err = createCharmApplicationSecret(ctx, st, 1, uri2, "mysql", sp)
		c.Assert(err, jc.ErrorIsNil)

		// The consumed revision 1 is the initial revision - will be ignored.
		saveConsumer(uri1, 1, "mediawiki/0")
		// The consumed revision 1 is the initial revision - will be ignored.
		saveConsumer(uri2, 1, "mediawiki/0")
	}, func(w watchertest.WatcherC[[]string]) {
		w.AssertNoChange()
	})

	harness.AddTest(func(c *gc.C) {
		// create revision 2.
		createNewRevision(c, st, uri1)
	}, func(w watchertest.WatcherC[[]string]) {
		w.Check(
			watchertest.StringSliceAssert(
				uri1.String(),
			),
		)
	})
	harness.Run(c, []string(nil))

	// Pretend that the agent restarted and the watcher is re-created.
	w1, err := svc.WatchConsumedSecretsChanges(ctx, "mediawiki/0")
	c.Assert(err, gc.IsNil)
	c.Assert(w1, gc.NotNil)
	defer watchertest.CleanKill(c, w1)

	harness1 := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w1))
	harness1.AddTest(func(c *gc.C) {}, func(w watchertest.WatcherC[[]string]) {
		w.Check(
			watchertest.StringSliceAssert(
				uri1.String(),
			),
		)
	})

	harness1.AddTest(func(c *gc.C) {
		// The consumed revision 2 is the updated current_revision.
		saveConsumer(uri1, 2, "mediawiki/0")
	}, func(w watchertest.WatcherC[[]string]) {
		w.AssertNoChange()
	})

	harness1.Run(c, []string(nil))

	// Pretend that the agent restarted and the watcher is re-created again.
	// Since we comsume the latest revision already, so there should be no change.
	w2, err := svc.WatchConsumedSecretsChanges(ctx, "mediawiki/0")
	c.Assert(err, gc.IsNil)
	c.Assert(w2, gc.NotNil)
	defer watchertest.CleanKill(c, w2)
	harness2 := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w2))
	harness2.AddTest(func(c *gc.C) {}, func(w watchertest.WatcherC[[]string]) {
		w.AssertNoChange()
	})
	harness2.Run(c, []string(nil))

}

func (s *watcherSuite) TestWatchConsumedRemoteSecretsChanges(c *gc.C) {
	s.setupUnits(c, "mediawiki")

	ctx := context.Background()
	svc, st := s.setupServiceAndState(c)

	saveConsumer := func(uri *coresecrets.URI, revision int, consumerID string) {
		consumer := &coresecrets.SecretConsumerMetadata{
			CurrentRevision: revision,
		}
		err := st.SaveSecretConsumer(ctx, uri, consumerID, consumer)
		c.Assert(err, jc.ErrorIsNil)
	}

	sourceModelUUID := uuid.MustNewUUID()
	uri1 := coresecrets.NewURI()
	uri1.SourceUUID = sourceModelUUID.String()

	uri2 := coresecrets.NewURI()
	uri2.SourceUUID = sourceModelUUID.String()

	w, err := svc.WatchConsumedSecretsChanges(ctx, "mediawiki/0")
	c.Assert(err, jc.ErrorIsNil)
	defer watchertest.CleanKill(c, w)

	harness := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w))
	harness.AddTest(func(c *gc.C) {
		// The consumed revision 1 is the initial revision - will be ignored.
		saveConsumer(uri1, 1, "mediawiki/0")
		// The consumed revision 1 is the initial revision - will be ignored.
		saveConsumer(uri2, 1, "mediawiki/0")
	}, func(w watchertest.WatcherC[[]string]) {
		w.AssertNoChange()
	})

	harness.AddTest(func(c *gc.C) {
		err = st.UpdateRemoteSecretRevision(ctx, uri1, 2)
		c.Assert(err, jc.ErrorIsNil)
	}, func(w watchertest.WatcherC[[]string]) {
		w.Check(
			watchertest.StringSliceAssert(
				uri1.String(),
			),
		)
	})

	harness.Run(c, []string(nil))

	// Pretend that the agent restarted and the watcher is re-created.
	w1, err := svc.WatchConsumedSecretsChanges(ctx, "mediawiki/0")
	c.Assert(err, jc.ErrorIsNil)
	defer watchertest.CleanKill(c, w1)

	harness1 := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w1))
	harness1.AddTest(func(c *gc.C) {}, func(w watchertest.WatcherC[[]string]) {
		w.Check(
			watchertest.StringSliceAssert(
				uri1.String(),
			),
		)
	})

	harness1.AddTest(func(c *gc.C) {
		// The consumed revision 2 is the updated current_revision.
		saveConsumer(uri1, 2, "mediawiki/0")
	}, func(w watchertest.WatcherC[[]string]) {
		w.AssertNoChange()
	})

	harness1.Run(c, []string(nil))

	// Pretend that the agent restarted and the watcher is re-created again.
	// Since we consume the latest revision already, so there should be no
	// change.
	w2, err := svc.WatchConsumedSecretsChanges(ctx, "mediawiki/0")
	c.Assert(err, jc.ErrorIsNil)
	defer watchertest.CleanKill(c, w2)

	harness2 := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w2))
	harness2.AddTest(func(c *gc.C) {}, func(w watchertest.WatcherC[[]string]) {
		w.AssertNoChange()
	})
	harness2.Run(c, []string(nil))
}

func (s *watcherSuite) TestWatchRemoteConsumedSecretsChanges(c *gc.C) {
	s.setupUnits(c, "mysql")

	ctx := context.Background()
	svc, st := s.setupServiceAndState(c)

	saveRemoteConsumer := func(uri *coresecrets.URI, revision int, consumerID string) {
		consumer := &coresecrets.SecretConsumerMetadata{
			CurrentRevision: revision,
		}
		err := st.SaveSecretRemoteConsumer(ctx, uri, consumerID, consumer)
		c.Assert(err, jc.ErrorIsNil)
	}

	uri1 := coresecrets.NewURI()
	uri1.SourceUUID = s.ModelUUID()
	uri2 := coresecrets.NewURI()
	uri2.SourceUUID = s.ModelUUID()

	w, err := svc.WatchRemoteConsumedSecretsChanges(ctx, "mediawiki")
	c.Assert(err, gc.IsNil)
	c.Assert(w, gc.NotNil)
	defer watchertest.CleanKill(c, w)

	harness := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w))
	harness.AddTest(func(c *gc.C) {
		sp := secret.UpsertSecretParams{
			Data: coresecrets.SecretData{"foo": "bar", "hello": "world"},
		}

		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err := createCharmApplicationSecret(ctx, st, 1, uri1, "mysql", sp)
		c.Assert(err, jc.ErrorIsNil)

		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err = createCharmApplicationSecret(ctx, st, 1, uri2, "mysql", sp)
		c.Assert(err, jc.ErrorIsNil)

		// The consumed revision 1 is the initial revision - will be ignored.
		saveRemoteConsumer(uri1, 1, "mediawiki/0")
		// The consumed revision 1 is the initial revision - will be ignored.
		saveRemoteConsumer(uri2, 1, "mediawiki/0")

	}, func(w watchertest.WatcherC[[]string]) {
		w.AssertNoChange()
	})

	harness.AddTest(func(c *gc.C) {
		// create revision 2.
		createNewRevision(c, st, uri1)
		err = st.UpdateRemoteSecretRevision(ctx, uri1, 2)
		c.Assert(err, jc.ErrorIsNil)
	}, func(w watchertest.WatcherC[[]string]) {
		w.Check(
			watchertest.StringSliceAssert(
				uri1.String(),
			),
		)
	})

	harness.Run(c, []string(nil))

	// Pretend that the agent restarted and the watcher is re-created.
	w1, err := svc.WatchRemoteConsumedSecretsChanges(ctx, "mediawiki")
	c.Assert(err, gc.IsNil)
	c.Assert(w1, gc.NotNil)
	defer watchertest.CleanKill(c, w1)

	harness1 := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w1))
	harness1.AddTest(func(c *gc.C) {}, func(w watchertest.WatcherC[[]string]) {
		w.Check(
			watchertest.StringSliceAssert(
				uri1.String(),
			),
		)
	})

	harness1.AddTest(func(c *gc.C) {
		// The consumed revision 2 is the updated current_revision.
		saveRemoteConsumer(uri1, 2, "mediawiki/0")
	}, func(w watchertest.WatcherC[[]string]) {
		w.AssertNoChange()
	})
	harness1.Run(c, []string(nil))

	// Pretend that the agent restarted and the watcher is re-created again.
	// Since we comsume the latest revision already, so there should be no change.
	w2, err := svc.WatchRemoteConsumedSecretsChanges(ctx, "mediawiki")
	c.Assert(err, gc.IsNil)
	c.Assert(w2, gc.NotNil)
	defer watchertest.CleanKill(c, w2)

	harness2 := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w2))
	harness2.AddTest(func(c *gc.C) {}, func(w watchertest.WatcherC[[]string]) {
		w.AssertNoChange()
	})

	harness2.Run(c, []string(nil))
}

func (s *watcherSuite) TestWatchSecretsRotationChanges(c *gc.C) {
	s.setupUnits(c, "mysql")
	s.setupUnits(c, "mediawiki")

	ctx := context.Background()
	svc, st := s.setupServiceAndState(c)

	uri1 := coresecrets.NewURI()
	uri2 := coresecrets.NewURI()

	w, err := svc.WatchSecretsRotationChanges(context.Background(),
		service.CharmSecretOwner{
			Kind: service.ApplicationOwner,
			ID:   "mysql",
		},
		service.CharmSecretOwner{
			Kind: service.UnitOwner,
			ID:   "mediawiki/0",
		},
	)
	c.Assert(err, gc.IsNil)
	c.Assert(w, gc.NotNil)
	defer watchertest.CleanKill(c, w)

	harness := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w))
	harness.AddTest(func(c *gc.C) {
		sp := secret.UpsertSecretParams{
			Data: coresecrets.SecretData{"foo": "bar", "hello": "world"},
		}

		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err := createCharmApplicationSecret(ctx, st, 1, uri1, "mysql", sp)
		c.Assert(err, jc.ErrorIsNil)

		sp.RevisionID = ptr(uuid.MustNewUUID().String())
		err = createCharmUnitSecret(ctx, st, 1, uri2, "mediawiki/0", sp)
		c.Assert(err, jc.ErrorIsNil)
		createNewRevision(c, st, uri2)
	}, func(w watchertest.WatcherC[[]corewatcher.SecretTriggerChange]) {
		w.AssertNoChange()
	})

	now := time.Now()
	harness.AddTest(func(c *gc.C) {
		err = st.SecretRotated(ctx, uri1, now.Add(1*time.Hour))
		c.Assert(err, jc.ErrorIsNil)
		err = st.SecretRotated(ctx, uri2, now.Add(2*time.Hour))
		c.Assert(err, jc.ErrorIsNil)
	}, func(w watchertest.WatcherC[[]corewatcher.SecretTriggerChange]) {
		w.Check(
			watchertest.SecretTriggerSliceAssert(
				corewatcher.SecretTriggerChange{
					URI:             uri1,
					Revision:        1,
					NextTriggerTime: now.Add(1 * time.Hour),
				},
				corewatcher.SecretTriggerChange{
					URI:             uri2,
					Revision:        2,
					NextTriggerTime: now.Add(2 * time.Hour),
				},
			),
		)
	})

	harness.Run(c, []corewatcher.SecretTriggerChange(nil))

	// Pretend that the agent restarted and the watcher is re-created.
	w1, err := svc.WatchSecretsRotationChanges(context.Background(),
		service.CharmSecretOwner{
			Kind: service.ApplicationOwner,
			ID:   "mysql",
		},
		service.CharmSecretOwner{
			Kind: service.UnitOwner,
			ID:   "mediawiki/0",
		},
	)
	c.Assert(err, gc.IsNil)
	c.Assert(w1, gc.NotNil)
	defer watchertest.CleanKill(c, w1)

	harness1 := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w1))

	harness1.AddTest(func(c *gc.C) {}, func(w watchertest.WatcherC[[]corewatcher.SecretTriggerChange]) {
		w.Check(
			watchertest.SecretTriggerSliceAssert(
				corewatcher.SecretTriggerChange{
					URI:             uri1,
					Revision:        1,
					NextTriggerTime: now.Add(1 * time.Hour),
				},
				corewatcher.SecretTriggerChange{
					URI:             uri2,
					Revision:        2,
					NextTriggerTime: now.Add(2 * time.Hour),
				},
			),
		)
	})

	harness1.Run(c, []corewatcher.SecretTriggerChange(nil))
}

func ptr[T any](v T) *T {
	return &v
}

func (s *watcherSuite) TestWatchSecretsRevisionExpiryChanges(c *gc.C) {
	s.setupUnits(c, "mysql")
	s.setupUnits(c, "mediawiki")

	ctx := context.Background()
	svc, st := s.setupServiceAndState(c)

	uri1 := coresecrets.NewURI()
	uri2 := coresecrets.NewURI()
	c.Logf("uri1: %v, uri2: %v", uri1, uri2)

	w, err := svc.WatchSecretRevisionsExpiryChanges(context.Background(),
		service.CharmSecretOwner{
			Kind: service.ApplicationOwner,
			ID:   "mysql",
		},
		service.CharmSecretOwner{
			Kind: service.UnitOwner,
			ID:   "mediawiki/0",
		},
	)
	c.Assert(err, gc.IsNil)
	c.Assert(w, gc.NotNil)
	defer watchertest.CleanKill(c, w)

	harness := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w))
	harness.AddTest(func(c *gc.C) {
		err := createCharmUnitSecret(ctx, st, 1, uri2, "mediawiki/0", secret.UpsertSecretParams{
			Data:       coresecrets.SecretData{"foo": "bar", "hello": "world"},
			RevisionID: ptr(uuid.MustNewUUID().String()),
		})
		c.Assert(err, jc.ErrorIsNil)
	}, func(w watchertest.WatcherC[[]corewatcher.SecretTriggerChange]) {
		w.AssertNoChange()
	})

	now := time.Now()
	harness.AddTest(func(c *gc.C) {
		err = createCharmApplicationSecret(ctx, st, 1, uri1, "mysql", secret.UpsertSecretParams{
			Data:       coresecrets.SecretData{"foo": "bar", "hello": "world"},
			ExpireTime: ptr(now.Add(1 * time.Hour)),
			RevisionID: ptr(uuid.MustNewUUID().String()),
		})
		c.Assert(err, jc.ErrorIsNil)

		err = st.RunAtomic(context.Background(), func(ctx domain.AtomicContext) error {
			return st.UpdateSecret(ctx, uri2, secret.UpsertSecretParams{
				Data:       coresecrets.SecretData{"foo-new": "bar-new"},
				ExpireTime: ptr(now.Add(2 * time.Hour)),
				RevisionID: ptr(uuid.MustNewUUID().String()),
			})
		})
		c.Assert(err, jc.ErrorIsNil)
	}, func(w watchertest.WatcherC[[]corewatcher.SecretTriggerChange]) {
		w.Check(
			watchertest.SecretTriggerSliceAssert(
				corewatcher.SecretTriggerChange{
					URI:             uri1,
					Revision:        1,
					NextTriggerTime: now.Add(1 * time.Hour),
				},
				corewatcher.SecretTriggerChange{
					URI:             uri2,
					Revision:        2,
					NextTriggerTime: now.Add(2 * time.Hour),
				},
			),
		)
	})

	harness.Run(c, []corewatcher.SecretTriggerChange(nil))

	// Pretend that the agent restarted and the watcher is re-created.
	w1, err := svc.WatchSecretRevisionsExpiryChanges(context.Background(),
		service.CharmSecretOwner{
			Kind: service.ApplicationOwner,
			ID:   "mysql",
		},
		service.CharmSecretOwner{
			Kind: service.UnitOwner,
			ID:   "mediawiki/0",
		},
	)
	c.Assert(err, gc.IsNil)
	c.Assert(w1, gc.NotNil)
	defer watchertest.CleanKill(c, w1)

	harness1 := watchertest.NewHarness(s, watchertest.NewWatcherC(c, w1))
	harness1.AddTest(func(c *gc.C) {}, func(w watchertest.WatcherC[[]corewatcher.SecretTriggerChange]) {
		w.Check(
			watchertest.SecretTriggerSliceAssert(
				corewatcher.SecretTriggerChange{
					URI:             uri1,
					Revision:        1,
					NextTriggerTime: now.Add(1 * time.Hour),
				},
				corewatcher.SecretTriggerChange{
					URI:             uri2,
					Revision:        2,
					NextTriggerTime: now.Add(2 * time.Hour),
				},
			),
		)
	})

	harness1.Run(c, []corewatcher.SecretTriggerChange(nil))
}

type stubCharm struct{}

var _ charm.Charm = (*stubCharm)(nil)

func (m *stubCharm) Meta() *charm.Meta {
	return &charm.Meta{
		Name: "foo",
	}
}

func (m *stubCharm) Manifest() *charm.Manifest {
	return &charm.Manifest{}
}

func (m *stubCharm) Config() *charm.Config {
	return &charm.Config{}
}

func (m *stubCharm) Actions() *charm.Actions {
	return &charm.Actions{}
}

func (m *stubCharm) Revision() int {
	return 1
}
