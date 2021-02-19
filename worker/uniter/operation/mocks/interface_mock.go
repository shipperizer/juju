// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/worker/uniter/operation (interfaces: Operation,Factory)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	charm "github.com/juju/charm/v9"
	hook "github.com/juju/juju/worker/uniter/hook"
	operation "github.com/juju/juju/worker/uniter/operation"
	remotestate "github.com/juju/juju/worker/uniter/remotestate"
	reflect "reflect"
)

// MockOperation is a mock of Operation interface
type MockOperation struct {
	ctrl     *gomock.Controller
	recorder *MockOperationMockRecorder
}

// MockOperationMockRecorder is the mock recorder for MockOperation
type MockOperationMockRecorder struct {
	mock *MockOperation
}

// NewMockOperation creates a new mock instance
func NewMockOperation(ctrl *gomock.Controller) *MockOperation {
	mock := &MockOperation{ctrl: ctrl}
	mock.recorder = &MockOperationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOperation) EXPECT() *MockOperationMockRecorder {
	return m.recorder
}

// Commit mocks base method
func (m *MockOperation) Commit(arg0 operation.State) (*operation.State, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", arg0)
	ret0, _ := ret[0].(*operation.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Commit indicates an expected call of Commit
func (mr *MockOperationMockRecorder) Commit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockOperation)(nil).Commit), arg0)
}

// Execute mocks base method
func (m *MockOperation) Execute(arg0 operation.State) (*operation.State, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0)
	ret0, _ := ret[0].(*operation.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockOperationMockRecorder) Execute(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockOperation)(nil).Execute), arg0)
}

// ExecutionGroup mocks base method
func (m *MockOperation) ExecutionGroup() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecutionGroup")
	ret0, _ := ret[0].(string)
	return ret0
}

// ExecutionGroup indicates an expected call of ExecutionGroup
func (mr *MockOperationMockRecorder) ExecutionGroup() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecutionGroup", reflect.TypeOf((*MockOperation)(nil).ExecutionGroup))
}

// NeedsGlobalMachineLock mocks base method
func (m *MockOperation) NeedsGlobalMachineLock() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NeedsGlobalMachineLock")
	ret0, _ := ret[0].(bool)
	return ret0
}

// NeedsGlobalMachineLock indicates an expected call of NeedsGlobalMachineLock
func (mr *MockOperationMockRecorder) NeedsGlobalMachineLock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NeedsGlobalMachineLock", reflect.TypeOf((*MockOperation)(nil).NeedsGlobalMachineLock))
}

// Prepare mocks base method
func (m *MockOperation) Prepare(arg0 operation.State) (*operation.State, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prepare", arg0)
	ret0, _ := ret[0].(*operation.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Prepare indicates an expected call of Prepare
func (mr *MockOperationMockRecorder) Prepare(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockOperation)(nil).Prepare), arg0)
}

// RemoteStateChanged mocks base method
func (m *MockOperation) RemoteStateChanged(arg0 remotestate.Snapshot) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoteStateChanged", arg0)
}

// RemoteStateChanged indicates an expected call of RemoteStateChanged
func (mr *MockOperationMockRecorder) RemoteStateChanged(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteStateChanged", reflect.TypeOf((*MockOperation)(nil).RemoteStateChanged), arg0)
}

// String mocks base method
func (m *MockOperation) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String
func (mr *MockOperationMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockOperation)(nil).String))
}

// MockFactory is a mock of Factory interface
type MockFactory struct {
	ctrl     *gomock.Controller
	recorder *MockFactoryMockRecorder
}

// MockFactoryMockRecorder is the mock recorder for MockFactory
type MockFactoryMockRecorder struct {
	mock *MockFactory
}

// NewMockFactory creates a new mock instance
func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	mock := &MockFactory{ctrl: ctrl}
	mock.recorder = &MockFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFactory) EXPECT() *MockFactoryMockRecorder {
	return m.recorder
}

// NewAcceptLeadership mocks base method
func (m *MockFactory) NewAcceptLeadership() (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewAcceptLeadership")
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewAcceptLeadership indicates an expected call of NewAcceptLeadership
func (mr *MockFactoryMockRecorder) NewAcceptLeadership() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewAcceptLeadership", reflect.TypeOf((*MockFactory)(nil).NewAcceptLeadership))
}

// NewAction mocks base method
func (m *MockFactory) NewAction(arg0 string) (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewAction", arg0)
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewAction indicates an expected call of NewAction
func (mr *MockFactoryMockRecorder) NewAction(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewAction", reflect.TypeOf((*MockFactory)(nil).NewAction), arg0)
}

// NewCommands mocks base method
func (m *MockFactory) NewCommands(arg0 operation.CommandArgs, arg1 operation.CommandResponseFunc) (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCommands", arg0, arg1)
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewCommands indicates an expected call of NewCommands
func (mr *MockFactoryMockRecorder) NewCommands(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCommands", reflect.TypeOf((*MockFactory)(nil).NewCommands), arg0, arg1)
}

// NewFailAction mocks base method
func (m *MockFactory) NewFailAction(arg0 string) (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewFailAction", arg0)
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewFailAction indicates an expected call of NewFailAction
func (mr *MockFactoryMockRecorder) NewFailAction(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewFailAction", reflect.TypeOf((*MockFactory)(nil).NewFailAction), arg0)
}

// NewInstall mocks base method
func (m *MockFactory) NewInstall(arg0 *charm.URL) (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewInstall", arg0)
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewInstall indicates an expected call of NewInstall
func (mr *MockFactoryMockRecorder) NewInstall(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewInstall", reflect.TypeOf((*MockFactory)(nil).NewInstall), arg0)
}

// NewNoOpFinishUpgradeSeries mocks base method
func (m *MockFactory) NewNoOpFinishUpgradeSeries() (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewNoOpFinishUpgradeSeries")
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewNoOpFinishUpgradeSeries indicates an expected call of NewNoOpFinishUpgradeSeries
func (mr *MockFactoryMockRecorder) NewNoOpFinishUpgradeSeries() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewNoOpFinishUpgradeSeries", reflect.TypeOf((*MockFactory)(nil).NewNoOpFinishUpgradeSeries))
}

// NewRemoteInit mocks base method
func (m *MockFactory) NewRemoteInit(arg0 remotestate.ContainerRunningStatus) (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRemoteInit", arg0)
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRemoteInit indicates an expected call of NewRemoteInit
func (mr *MockFactoryMockRecorder) NewRemoteInit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRemoteInit", reflect.TypeOf((*MockFactory)(nil).NewRemoteInit), arg0)
}

// NewResignLeadership mocks base method
func (m *MockFactory) NewResignLeadership() (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewResignLeadership")
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewResignLeadership indicates an expected call of NewResignLeadership
func (mr *MockFactoryMockRecorder) NewResignLeadership() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewResignLeadership", reflect.TypeOf((*MockFactory)(nil).NewResignLeadership))
}

// NewResolvedUpgrade mocks base method
func (m *MockFactory) NewResolvedUpgrade(arg0 *charm.URL) (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewResolvedUpgrade", arg0)
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewResolvedUpgrade indicates an expected call of NewResolvedUpgrade
func (mr *MockFactoryMockRecorder) NewResolvedUpgrade(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewResolvedUpgrade", reflect.TypeOf((*MockFactory)(nil).NewResolvedUpgrade), arg0)
}

// NewRevertUpgrade mocks base method
func (m *MockFactory) NewRevertUpgrade(arg0 *charm.URL) (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRevertUpgrade", arg0)
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRevertUpgrade indicates an expected call of NewRevertUpgrade
func (mr *MockFactoryMockRecorder) NewRevertUpgrade(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRevertUpgrade", reflect.TypeOf((*MockFactory)(nil).NewRevertUpgrade), arg0)
}

// NewRunHook mocks base method
func (m *MockFactory) NewRunHook(arg0 hook.Info) (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRunHook", arg0)
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRunHook indicates an expected call of NewRunHook
func (mr *MockFactoryMockRecorder) NewRunHook(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRunHook", reflect.TypeOf((*MockFactory)(nil).NewRunHook), arg0)
}

// NewSkipHook mocks base method
func (m *MockFactory) NewSkipHook(arg0 hook.Info) (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSkipHook", arg0)
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewSkipHook indicates an expected call of NewSkipHook
func (mr *MockFactoryMockRecorder) NewSkipHook(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSkipHook", reflect.TypeOf((*MockFactory)(nil).NewSkipHook), arg0)
}

// NewSkipRemoteInit mocks base method
func (m *MockFactory) NewSkipRemoteInit(arg0 bool) (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSkipRemoteInit", arg0)
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewSkipRemoteInit indicates an expected call of NewSkipRemoteInit
func (mr *MockFactoryMockRecorder) NewSkipRemoteInit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSkipRemoteInit", reflect.TypeOf((*MockFactory)(nil).NewSkipRemoteInit), arg0)
}

// NewUpgrade mocks base method
func (m *MockFactory) NewUpgrade(arg0 *charm.URL) (operation.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewUpgrade", arg0)
	ret0, _ := ret[0].(operation.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewUpgrade indicates an expected call of NewUpgrade
func (mr *MockFactoryMockRecorder) NewUpgrade(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewUpgrade", reflect.TypeOf((*MockFactory)(nil).NewUpgrade), arg0)
}
