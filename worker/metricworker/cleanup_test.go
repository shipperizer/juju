// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package metricworker_test

import (
	stdtesting "testing"
	"time"

	gc "launchpad.net/gocheck"

	"github.com/juju/juju/juju/testing"
	"github.com/juju/juju/state/api/metricsmanager"
	coretesting "github.com/juju/juju/testing"
	"github.com/juju/juju/testing/factory"
	"github.com/juju/juju/worker/metricworker"
)

func TestPackage(t *stdtesting.T) {
	coretesting.MgoTestPackage(t)
}

type CleanupSuite struct {
	testing.JujuConnSuite
}

var (
	_                = gc.Suite(&CleanupSuite{})
	unnaceptableWait = time.Second * 5
)

func (s *CleanupSuite) SetUpTest(c *gc.C) {
	s.JujuConnSuite.SetUpTest(c)
}

// TestCleaner create 2 metrics, one old and one new.
// After a single run of the cleanup worker it expects the
// old one to be deleted
func (s *CleanupSuite) TestCleaner(c *gc.C) {
	unit := s.Factory.MakeUnit(c, nil)
	oldTime := time.Now().Add(-(time.Hour * 25))
	now := time.Now()
	oldMetric := s.Factory.MakeMetric(c, &factory.MetricParams{Unit: unit, Sent: true, Time: &oldTime})
	newMetric := s.Factory.MakeMetric(c, &factory.MetricParams{Unit: unit, Sent: true, Time: &now})

	notify := make(chan struct{})
	client := metricsmanager.NewClient(s.APIState)
	worker := metricworker.NewCleanup(client, notify)
	defer worker.Kill()
	select {
	case <-notify:
	case <-time.After(unnaceptableWait):
		c.Fatalf("the cleanup function should have fired by now")
	}
	_, err := s.State.MetricBatch(newMetric.UUID())
	c.Assert(err, gc.IsNil)

	_, err = s.State.MetricBatch(oldMetric.UUID())
	c.Assert(err, gc.ErrorMatches, "not found")
}
