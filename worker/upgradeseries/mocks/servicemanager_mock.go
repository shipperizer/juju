// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/service (interfaces: SystemdServiceManager)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	common "github.com/juju/juju/internal/service/common"
	gomock "go.uber.org/mock/gomock"
)

// MockSystemdServiceManager is a mock of SystemdServiceManager interface.
type MockSystemdServiceManager struct {
	ctrl     *gomock.Controller
	recorder *MockSystemdServiceManagerMockRecorder
}

// MockSystemdServiceManagerMockRecorder is the mock recorder for MockSystemdServiceManager.
type MockSystemdServiceManagerMockRecorder struct {
	mock *MockSystemdServiceManager
}

// NewMockSystemdServiceManager creates a new mock instance.
func NewMockSystemdServiceManager(ctrl *gomock.Controller) *MockSystemdServiceManager {
	mock := &MockSystemdServiceManager{ctrl: ctrl}
	mock.recorder = &MockSystemdServiceManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSystemdServiceManager) EXPECT() *MockSystemdServiceManagerMockRecorder {
	return m.recorder
}

// CreateAgentConf mocks base method.
func (m *MockSystemdServiceManager) CreateAgentConf(arg0, arg1 string) (common.Conf, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAgentConf", arg0, arg1)
	ret0, _ := ret[0].(common.Conf)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAgentConf indicates an expected call of CreateAgentConf.
func (mr *MockSystemdServiceManagerMockRecorder) CreateAgentConf(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAgentConf", reflect.TypeOf((*MockSystemdServiceManager)(nil).CreateAgentConf), arg0, arg1)
}

// FindAgents mocks base method.
func (m *MockSystemdServiceManager) FindAgents(arg0 string) (string, []string, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAgents", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].([]string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// FindAgents indicates an expected call of FindAgents.
func (mr *MockSystemdServiceManagerMockRecorder) FindAgents(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAgents", reflect.TypeOf((*MockSystemdServiceManager)(nil).FindAgents), arg0)
}

// WriteServiceFile mocks base method.
func (m *MockSystemdServiceManager) WriteServiceFile() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteServiceFile")
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteServiceFile indicates an expected call of WriteServiceFile.
func (mr *MockSystemdServiceManagerMockRecorder) WriteServiceFile() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteServiceFile", reflect.TypeOf((*MockSystemdServiceManager)(nil).WriteServiceFile))
}

// WriteSystemdAgent mocks base method.
func (m *MockSystemdServiceManager) WriteSystemdAgent(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteSystemdAgent", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteSystemdAgent indicates an expected call of WriteSystemdAgent.
func (mr *MockSystemdServiceManagerMockRecorder) WriteSystemdAgent(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteSystemdAgent", reflect.TypeOf((*MockSystemdServiceManager)(nil).WriteSystemdAgent), arg0, arg1, arg2)
}
