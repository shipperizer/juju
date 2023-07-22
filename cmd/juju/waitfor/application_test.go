// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package waitfor

import (
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/cmd/juju/waitfor/query"
	"github.com/juju/juju/core/life"
	"github.com/juju/juju/core/status"
	"github.com/juju/juju/rpc/params"
)

type applicationScopeSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&applicationScopeSuite{})

func (s *applicationScopeSuite) TestGetIdentValue(c *gc.C) {
	tests := []struct {
		Field           string
		ApplicationInfo *params.ApplicationInfo
		Expected        query.Box
	}{{
		Field:           "name",
		ApplicationInfo: &params.ApplicationInfo{Name: "application name"},
		Expected:        query.NewString("application name"),
	}, {
		Field:           "life",
		ApplicationInfo: &params.ApplicationInfo{Life: life.Alive},
		Expected:        query.NewString("alive"),
	}, {
		Field:           "charm-url",
		ApplicationInfo: &params.ApplicationInfo{CharmURL: "ch:charm"},
		Expected:        query.NewString("ch:charm"),
	}, {
		Field:           "subordinate",
		ApplicationInfo: &params.ApplicationInfo{Subordinate: true},
		Expected:        query.NewBool(true),
	}, {
		Field: "status",
		ApplicationInfo: &params.ApplicationInfo{Status: params.StatusInfo{
			Current: status.Active,
		}},
		Expected: query.NewString("active"),
	}, {
		Field:           "workload-version",
		ApplicationInfo: &params.ApplicationInfo{WorkloadVersion: "1.2.3"},
		Expected:        query.NewString("1.2.3"),
	}}
	for i, test := range tests {
		c.Logf("%d: GetIdentValue %q", i, test.Field)
		scope := ApplicationScope{
			ctx:             MakeScopeContext(),
			ApplicationInfo: test.ApplicationInfo,
		}
		result, err := scope.GetIdentValue(test.Field)
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(result, gc.DeepEquals, test.Expected)
	}
}

func (s *applicationScopeSuite) TestGetIdentValueError(c *gc.C) {
	scope := ApplicationScope{
		ctx:             MakeScopeContext(),
		ApplicationInfo: &params.ApplicationInfo{},
	}
	result, err := scope.GetIdentValue("bad")
	c.Assert(err, gc.ErrorMatches, `.*"bad" on ApplicationInfo.*`)
	c.Assert(result, gc.IsNil)
}

func (s *applicationScopeSuite) TestDeriveApplicationStatus(c *gc.C) {
	tests := []struct {
		status   status.Status
		units    map[string]*params.UnitInfo
		expected string
	}{{
		status:   status.Unset,
		units:    nil,
		expected: "unknown",
	}, {
		status: status.Unset,
		units: map[string]*params.UnitInfo{
			"foo": {WorkloadStatus: params.StatusInfo{Current: status.Active}},
		},
		expected: "active",
	}, {
		status: status.Unset,
		units: map[string]*params.UnitInfo{
			"foo1": {WorkloadStatus: params.StatusInfo{Current: status.Error}},
			"foo2": {WorkloadStatus: params.StatusInfo{Current: status.Active}},
			"foo3": {WorkloadStatus: params.StatusInfo{Current: status.Active}},
		},
		expected: "error",
	}, {
		status: status.Unset,
		units: map[string]*params.UnitInfo{
			"foo1": {WorkloadStatus: params.StatusInfo{Current: status.Terminated}},
			"foo2": {WorkloadStatus: params.StatusInfo{Current: status.Active}},
			"foo3": {WorkloadStatus: params.StatusInfo{Current: status.Active}},
		},
		expected: "active",
	}, {
		status:   status.Unknown,
		units:    nil,
		expected: "unknown",
	}, {
		status: status.Error,
		units: map[string]*params.UnitInfo{
			"foo": {WorkloadStatus: params.StatusInfo{Current: status.Active}},
		},
		expected: "error",
	}, {
		status: status.Active,
		units: map[string]*params.UnitInfo{
			"foo1": {WorkloadStatus: params.StatusInfo{Current: status.Error}},
			"foo2": {WorkloadStatus: params.StatusInfo{Current: status.Active}},
			"foo3": {WorkloadStatus: params.StatusInfo{Current: status.Active}},
		},
		expected: "active",
	}}
	for _, test := range tests {
		status := deriveApplicationStatus(test.status, test.units)
		c.Check(status.String(), gc.Equals, test.expected)
	}
}
