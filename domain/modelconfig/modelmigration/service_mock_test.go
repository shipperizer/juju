// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/domain/modelconfig/service (interfaces: ModelDefaultsProvider)
//
// Generated by this command:
//
//	mockgen -package modelmigration -destination service_mock_test.go github.com/juju/juju/domain/modelconfig/service ModelDefaultsProvider
//

// Package modelmigration is a generated GoMock package.
package modelmigration

import (
	context "context"
	reflect "reflect"

	modeldefaults "github.com/juju/juju/domain/modeldefaults"
	gomock "go.uber.org/mock/gomock"
)

// MockModelDefaultsProvider is a mock of ModelDefaultsProvider interface.
type MockModelDefaultsProvider struct {
	ctrl     *gomock.Controller
	recorder *MockModelDefaultsProviderMockRecorder
}

// MockModelDefaultsProviderMockRecorder is the mock recorder for MockModelDefaultsProvider.
type MockModelDefaultsProviderMockRecorder struct {
	mock *MockModelDefaultsProvider
}

// NewMockModelDefaultsProvider creates a new mock instance.
func NewMockModelDefaultsProvider(ctrl *gomock.Controller) *MockModelDefaultsProvider {
	mock := &MockModelDefaultsProvider{ctrl: ctrl}
	mock.recorder = &MockModelDefaultsProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModelDefaultsProvider) EXPECT() *MockModelDefaultsProviderMockRecorder {
	return m.recorder
}

// ModelDefaults mocks base method.
func (m *MockModelDefaultsProvider) ModelDefaults(arg0 context.Context) (modeldefaults.Defaults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelDefaults", arg0)
	ret0, _ := ret[0].(modeldefaults.Defaults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelDefaults indicates an expected call of ModelDefaults.
func (mr *MockModelDefaultsProviderMockRecorder) ModelDefaults(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelDefaults", reflect.TypeOf((*MockModelDefaultsProvider)(nil).ModelDefaults), arg0)
}
