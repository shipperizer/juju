// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/worker/objectstore (interfaces: TrackedObjectStore)

// Package objectstore is a generated GoMock package.
package objectstore

import (
	context "context"
	io "io"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockTrackedObjectStore is a mock of TrackedObjectStore interface.
type MockTrackedObjectStore struct {
	ctrl     *gomock.Controller
	recorder *MockTrackedObjectStoreMockRecorder
}

// MockTrackedObjectStoreMockRecorder is the mock recorder for MockTrackedObjectStore.
type MockTrackedObjectStoreMockRecorder struct {
	mock *MockTrackedObjectStore
}

// NewMockTrackedObjectStore creates a new mock instance.
func NewMockTrackedObjectStore(ctrl *gomock.Controller) *MockTrackedObjectStore {
	mock := &MockTrackedObjectStore{ctrl: ctrl}
	mock.recorder = &MockTrackedObjectStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTrackedObjectStore) EXPECT() *MockTrackedObjectStoreMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockTrackedObjectStore) Get(arg0 context.Context, arg1 string) (io.ReadCloser, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockTrackedObjectStoreMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTrackedObjectStore)(nil).Get), arg0, arg1)
}

// Kill mocks base method.
func (m *MockTrackedObjectStore) Kill() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Kill")
}

// Kill indicates an expected call of Kill.
func (mr *MockTrackedObjectStoreMockRecorder) Kill() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Kill", reflect.TypeOf((*MockTrackedObjectStore)(nil).Kill))
}

// Put mocks base method.
func (m *MockTrackedObjectStore) Put(arg0 context.Context, arg1 string, arg2 io.Reader, arg3 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockTrackedObjectStoreMockRecorder) Put(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockTrackedObjectStore)(nil).Put), arg0, arg1, arg2, arg3)
}

// Wait mocks base method.
func (m *MockTrackedObjectStore) Wait() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait")
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait.
func (mr *MockTrackedObjectStoreMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockTrackedObjectStore)(nil).Wait))
}
