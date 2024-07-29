// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/state/migrations (interfaces: MigrationFirewallRule,FirewallRuleSource,FirewallRulesModel)
//
// Generated by this command:
//
//	mockgen -typed -package migrations -destination firewallrules_mock_test.go github.com/juju/juju/state/migrations MigrationFirewallRule,FirewallRuleSource,FirewallRulesModel
//

// Package migrations is a generated GoMock package.
package migrations

import (
	reflect "reflect"

	description "github.com/juju/description/v8"
	firewall "github.com/juju/juju/core/network/firewall"
	gomock "go.uber.org/mock/gomock"
)

// MockMigrationFirewallRule is a mock of MigrationFirewallRule interface.
type MockMigrationFirewallRule struct {
	ctrl     *gomock.Controller
	recorder *MockMigrationFirewallRuleMockRecorder
}

// MockMigrationFirewallRuleMockRecorder is the mock recorder for MockMigrationFirewallRule.
type MockMigrationFirewallRuleMockRecorder struct {
	mock *MockMigrationFirewallRule
}

// NewMockMigrationFirewallRule creates a new mock instance.
func NewMockMigrationFirewallRule(ctrl *gomock.Controller) *MockMigrationFirewallRule {
	mock := &MockMigrationFirewallRule{ctrl: ctrl}
	mock.recorder = &MockMigrationFirewallRuleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMigrationFirewallRule) EXPECT() *MockMigrationFirewallRuleMockRecorder {
	return m.recorder
}

// ID mocks base method.
func (m *MockMigrationFirewallRule) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockMigrationFirewallRuleMockRecorder) ID() *MockMigrationFirewallRuleIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockMigrationFirewallRule)(nil).ID))
	return &MockMigrationFirewallRuleIDCall{Call: call}
}

// MockMigrationFirewallRuleIDCall wrap *gomock.Call
type MockMigrationFirewallRuleIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMigrationFirewallRuleIDCall) Return(arg0 string) *MockMigrationFirewallRuleIDCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMigrationFirewallRuleIDCall) Do(f func() string) *MockMigrationFirewallRuleIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMigrationFirewallRuleIDCall) DoAndReturn(f func() string) *MockMigrationFirewallRuleIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WellKnownService mocks base method.
func (m *MockMigrationFirewallRule) WellKnownService() firewall.WellKnownServiceType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WellKnownService")
	ret0, _ := ret[0].(firewall.WellKnownServiceType)
	return ret0
}

// WellKnownService indicates an expected call of WellKnownService.
func (mr *MockMigrationFirewallRuleMockRecorder) WellKnownService() *MockMigrationFirewallRuleWellKnownServiceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WellKnownService", reflect.TypeOf((*MockMigrationFirewallRule)(nil).WellKnownService))
	return &MockMigrationFirewallRuleWellKnownServiceCall{Call: call}
}

// MockMigrationFirewallRuleWellKnownServiceCall wrap *gomock.Call
type MockMigrationFirewallRuleWellKnownServiceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMigrationFirewallRuleWellKnownServiceCall) Return(arg0 firewall.WellKnownServiceType) *MockMigrationFirewallRuleWellKnownServiceCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMigrationFirewallRuleWellKnownServiceCall) Do(f func() firewall.WellKnownServiceType) *MockMigrationFirewallRuleWellKnownServiceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMigrationFirewallRuleWellKnownServiceCall) DoAndReturn(f func() firewall.WellKnownServiceType) *MockMigrationFirewallRuleWellKnownServiceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WhitelistCIDRs mocks base method.
func (m *MockMigrationFirewallRule) WhitelistCIDRs() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WhitelistCIDRs")
	ret0, _ := ret[0].([]string)
	return ret0
}

// WhitelistCIDRs indicates an expected call of WhitelistCIDRs.
func (mr *MockMigrationFirewallRuleMockRecorder) WhitelistCIDRs() *MockMigrationFirewallRuleWhitelistCIDRsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WhitelistCIDRs", reflect.TypeOf((*MockMigrationFirewallRule)(nil).WhitelistCIDRs))
	return &MockMigrationFirewallRuleWhitelistCIDRsCall{Call: call}
}

// MockMigrationFirewallRuleWhitelistCIDRsCall wrap *gomock.Call
type MockMigrationFirewallRuleWhitelistCIDRsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMigrationFirewallRuleWhitelistCIDRsCall) Return(arg0 []string) *MockMigrationFirewallRuleWhitelistCIDRsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMigrationFirewallRuleWhitelistCIDRsCall) Do(f func() []string) *MockMigrationFirewallRuleWhitelistCIDRsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMigrationFirewallRuleWhitelistCIDRsCall) DoAndReturn(f func() []string) *MockMigrationFirewallRuleWhitelistCIDRsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockFirewallRuleSource is a mock of FirewallRuleSource interface.
type MockFirewallRuleSource struct {
	ctrl     *gomock.Controller
	recorder *MockFirewallRuleSourceMockRecorder
}

// MockFirewallRuleSourceMockRecorder is the mock recorder for MockFirewallRuleSource.
type MockFirewallRuleSourceMockRecorder struct {
	mock *MockFirewallRuleSource
}

// NewMockFirewallRuleSource creates a new mock instance.
func NewMockFirewallRuleSource(ctrl *gomock.Controller) *MockFirewallRuleSource {
	mock := &MockFirewallRuleSource{ctrl: ctrl}
	mock.recorder = &MockFirewallRuleSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFirewallRuleSource) EXPECT() *MockFirewallRuleSourceMockRecorder {
	return m.recorder
}

// AllFirewallRules mocks base method.
func (m *MockFirewallRuleSource) AllFirewallRules() ([]MigrationFirewallRule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllFirewallRules")
	ret0, _ := ret[0].([]MigrationFirewallRule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllFirewallRules indicates an expected call of AllFirewallRules.
func (mr *MockFirewallRuleSourceMockRecorder) AllFirewallRules() *MockFirewallRuleSourceAllFirewallRulesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllFirewallRules", reflect.TypeOf((*MockFirewallRuleSource)(nil).AllFirewallRules))
	return &MockFirewallRuleSourceAllFirewallRulesCall{Call: call}
}

// MockFirewallRuleSourceAllFirewallRulesCall wrap *gomock.Call
type MockFirewallRuleSourceAllFirewallRulesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockFirewallRuleSourceAllFirewallRulesCall) Return(arg0 []MigrationFirewallRule, arg1 error) *MockFirewallRuleSourceAllFirewallRulesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockFirewallRuleSourceAllFirewallRulesCall) Do(f func() ([]MigrationFirewallRule, error)) *MockFirewallRuleSourceAllFirewallRulesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockFirewallRuleSourceAllFirewallRulesCall) DoAndReturn(f func() ([]MigrationFirewallRule, error)) *MockFirewallRuleSourceAllFirewallRulesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockFirewallRulesModel is a mock of FirewallRulesModel interface.
type MockFirewallRulesModel struct {
	ctrl     *gomock.Controller
	recorder *MockFirewallRulesModelMockRecorder
}

// MockFirewallRulesModelMockRecorder is the mock recorder for MockFirewallRulesModel.
type MockFirewallRulesModelMockRecorder struct {
	mock *MockFirewallRulesModel
}

// NewMockFirewallRulesModel creates a new mock instance.
func NewMockFirewallRulesModel(ctrl *gomock.Controller) *MockFirewallRulesModel {
	mock := &MockFirewallRulesModel{ctrl: ctrl}
	mock.recorder = &MockFirewallRulesModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFirewallRulesModel) EXPECT() *MockFirewallRulesModelMockRecorder {
	return m.recorder
}

// AddFirewallRule mocks base method.
func (m *MockFirewallRulesModel) AddFirewallRule(arg0 description.FirewallRuleArgs) description.FirewallRule {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFirewallRule", arg0)
	ret0, _ := ret[0].(description.FirewallRule)
	return ret0
}

// AddFirewallRule indicates an expected call of AddFirewallRule.
func (mr *MockFirewallRulesModelMockRecorder) AddFirewallRule(arg0 any) *MockFirewallRulesModelAddFirewallRuleCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFirewallRule", reflect.TypeOf((*MockFirewallRulesModel)(nil).AddFirewallRule), arg0)
	return &MockFirewallRulesModelAddFirewallRuleCall{Call: call}
}

// MockFirewallRulesModelAddFirewallRuleCall wrap *gomock.Call
type MockFirewallRulesModelAddFirewallRuleCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockFirewallRulesModelAddFirewallRuleCall) Return(arg0 description.FirewallRule) *MockFirewallRulesModelAddFirewallRuleCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockFirewallRulesModelAddFirewallRuleCall) Do(f func(description.FirewallRuleArgs) description.FirewallRule) *MockFirewallRulesModelAddFirewallRuleCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockFirewallRulesModelAddFirewallRuleCall) DoAndReturn(f func(description.FirewallRuleArgs) description.FirewallRule) *MockFirewallRulesModelAddFirewallRuleCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
