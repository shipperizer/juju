// Copyright 2021 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package dbaccessor

import (
	"github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/dependency"

	coreagent "github.com/juju/juju/agent"
	coredatabase "github.com/juju/juju/core/database"
	"github.com/juju/juju/database"
	"github.com/juju/juju/database/app"
	"github.com/juju/juju/worker/common"
)

// Logger represents the logging methods called.
type Logger interface {
	Errorf(message string, args ...interface{})
	Warningf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Debugf(message string, args ...interface{})
	Tracef(message string, args ...interface{})

	// Logf is used to proxy Dqlite logs via this logger.
	Logf(level loggo.Level, msg string, args ...interface{})

	IsTraceEnabled() bool
}

// Hub defines the methods of the API server central hub
// that the DB accessor requires.
type Hub interface {
	Subscribe(topic string, handler interface{}) (func(), error)
	Publish(topic string, data interface{}) (func(), error)
}

// ManifoldConfig contains:
// - The names of other manifolds on which the DB accessor depends.
// - Other dependencies from ManifoldsConfig required by the worker.
type ManifoldConfig struct {
	AgentName   string
	Clock       clock.Clock
	Hub         Hub
	Logger      Logger
	NewApp      func(string, ...app.Option) (DBApp, error)
	NewDBWorker func(DBApp, string, ...TrackedDBWorkerOption) (TrackedDB, error)
}

func (cfg ManifoldConfig) Validate() error {
	if cfg.AgentName == "" {
		return errors.NotValidf("empty AgentName")
	}
	if cfg.Clock == nil {
		return errors.NotValidf("nil Clock")
	}
	if cfg.Hub == nil {
		return errors.NotValidf("nil Hub")
	}
	if cfg.Logger == nil {
		return errors.NotValidf("nil Logger")
	}
	if cfg.NewApp == nil {
		return errors.NotValidf("nil NewApp")
	}
	if cfg.NewDBWorker == nil {
		return errors.NotValidf("nil NewDBWorker")
	}
	return nil
}

// Manifold returns a dependency manifold that runs the dbaccessor
// worker, using the resource names defined in the supplied config.
func Manifold(config ManifoldConfig) dependency.Manifold {
	return dependency.Manifold{
		Inputs: []string{
			config.AgentName,
		},
		Output: dbAccessorOutput,
		Start: func(context dependency.Context) (worker.Worker, error) {
			if err := config.Validate(); err != nil {
				return nil, errors.Trace(err)
			}

			var agent coreagent.Agent
			if err := context.Get(config.AgentName, &agent); err != nil {
				return nil, err
			}
			agentConfig := agent.CurrentConfig()

			cfg := WorkerConfig{
				NodeManager:  database.NewNodeManager(agentConfig, config.Logger),
				Clock:        config.Clock,
				Hub:          config.Hub,
				ControllerID: agentConfig.Tag().Id(),
				Logger:       config.Logger,
				NewApp:       config.NewApp,
				NewDBWorker:  config.NewDBWorker,
			}

			w, err := newWorker(cfg)
			if err != nil {
				return nil, errors.Trace(err)
			}
			return w, nil
		},
	}
}

func dbAccessorOutput(in worker.Worker, out interface{}) error {
	if w, ok := in.(*common.CleanupWorker); ok {
		in = w.Worker
	}
	w, ok := in.(*dbWorker)
	if !ok {
		return errors.Errorf("expected input of type dbWorker, got %T", in)
	}

	switch out := out.(type) {
	case *coredatabase.DBGetter:
		var target coredatabase.DBGetter = w
		*out = target
	default:
		return errors.Errorf("expected output of *dbaccessor.DBGetter, got %T", out)
	}
	return nil
}
