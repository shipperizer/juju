// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/controller/charmdownloader (interfaces: StateBackend,Application,Charm,Downloader,AuthChecker,ResourcesBackend)
//
// Generated by this command:
//
//	mockgen -typed -package mocks -destination mocks/mocks.go github.com/juju/juju/apiserver/facades/controller/charmdownloader StateBackend,Application,Charm,Downloader,AuthChecker,ResourcesBackend
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	charmdownloader "github.com/juju/juju/apiserver/facades/controller/charmdownloader"
	charm "github.com/juju/juju/core/charm"
	status "github.com/juju/juju/core/status"
	charm0 "github.com/juju/juju/internal/charm"
	services "github.com/juju/juju/internal/charm/services"
	state "github.com/juju/juju/state"
	worker "github.com/juju/worker/v4"
	gomock "go.uber.org/mock/gomock"
)

// MockStateBackend is a mock of StateBackend interface.
type MockStateBackend struct {
	ctrl     *gomock.Controller
	recorder *MockStateBackendMockRecorder
}

// MockStateBackendMockRecorder is the mock recorder for MockStateBackend.
type MockStateBackendMockRecorder struct {
	mock *MockStateBackend
}

// NewMockStateBackend creates a new mock instance.
func NewMockStateBackend(ctrl *gomock.Controller) *MockStateBackend {
	mock := &MockStateBackend{ctrl: ctrl}
	mock.recorder = &MockStateBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStateBackend) EXPECT() *MockStateBackendMockRecorder {
	return m.recorder
}

// Application mocks base method.
func (m *MockStateBackend) Application(arg0 string) (charmdownloader.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Application", arg0)
	ret0, _ := ret[0].(charmdownloader.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Application indicates an expected call of Application.
func (mr *MockStateBackendMockRecorder) Application(arg0 any) *MockStateBackendApplicationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Application", reflect.TypeOf((*MockStateBackend)(nil).Application), arg0)
	return &MockStateBackendApplicationCall{Call: call}
}

// MockStateBackendApplicationCall wrap *gomock.Call
type MockStateBackendApplicationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateBackendApplicationCall) Return(arg0 charmdownloader.Application, arg1 error) *MockStateBackendApplicationCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateBackendApplicationCall) Do(f func(string) (charmdownloader.Application, error)) *MockStateBackendApplicationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateBackendApplicationCall) DoAndReturn(f func(string) (charmdownloader.Application, error)) *MockStateBackendApplicationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ModelUUID mocks base method.
func (m *MockStateBackend) ModelUUID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelUUID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ModelUUID indicates an expected call of ModelUUID.
func (mr *MockStateBackendMockRecorder) ModelUUID() *MockStateBackendModelUUIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelUUID", reflect.TypeOf((*MockStateBackend)(nil).ModelUUID))
	return &MockStateBackendModelUUIDCall{Call: call}
}

// MockStateBackendModelUUIDCall wrap *gomock.Call
type MockStateBackendModelUUIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateBackendModelUUIDCall) Return(arg0 string) *MockStateBackendModelUUIDCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateBackendModelUUIDCall) Do(f func() string) *MockStateBackendModelUUIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateBackendModelUUIDCall) DoAndReturn(f func() string) *MockStateBackendModelUUIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// PrepareCharmUpload mocks base method.
func (m *MockStateBackend) PrepareCharmUpload(arg0 string) (services.UploadedCharm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareCharmUpload", arg0)
	ret0, _ := ret[0].(services.UploadedCharm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareCharmUpload indicates an expected call of PrepareCharmUpload.
func (mr *MockStateBackendMockRecorder) PrepareCharmUpload(arg0 any) *MockStateBackendPrepareCharmUploadCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareCharmUpload", reflect.TypeOf((*MockStateBackend)(nil).PrepareCharmUpload), arg0)
	return &MockStateBackendPrepareCharmUploadCall{Call: call}
}

// MockStateBackendPrepareCharmUploadCall wrap *gomock.Call
type MockStateBackendPrepareCharmUploadCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateBackendPrepareCharmUploadCall) Return(arg0 services.UploadedCharm, arg1 error) *MockStateBackendPrepareCharmUploadCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateBackendPrepareCharmUploadCall) Do(f func(string) (services.UploadedCharm, error)) *MockStateBackendPrepareCharmUploadCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateBackendPrepareCharmUploadCall) DoAndReturn(f func(string) (services.UploadedCharm, error)) *MockStateBackendPrepareCharmUploadCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdateUploadedCharm mocks base method.
func (m *MockStateBackend) UpdateUploadedCharm(arg0 state.CharmInfo) (services.UploadedCharm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUploadedCharm", arg0)
	ret0, _ := ret[0].(services.UploadedCharm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUploadedCharm indicates an expected call of UpdateUploadedCharm.
func (mr *MockStateBackendMockRecorder) UpdateUploadedCharm(arg0 any) *MockStateBackendUpdateUploadedCharmCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUploadedCharm", reflect.TypeOf((*MockStateBackend)(nil).UpdateUploadedCharm), arg0)
	return &MockStateBackendUpdateUploadedCharmCall{Call: call}
}

// MockStateBackendUpdateUploadedCharmCall wrap *gomock.Call
type MockStateBackendUpdateUploadedCharmCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateBackendUpdateUploadedCharmCall) Return(arg0 services.UploadedCharm, arg1 error) *MockStateBackendUpdateUploadedCharmCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateBackendUpdateUploadedCharmCall) Do(f func(state.CharmInfo) (services.UploadedCharm, error)) *MockStateBackendUpdateUploadedCharmCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateBackendUpdateUploadedCharmCall) DoAndReturn(f func(state.CharmInfo) (services.UploadedCharm, error)) *MockStateBackendUpdateUploadedCharmCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchApplicationsWithPendingCharms mocks base method.
func (m *MockStateBackend) WatchApplicationsWithPendingCharms() state.StringsWatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchApplicationsWithPendingCharms")
	ret0, _ := ret[0].(state.StringsWatcher)
	return ret0
}

// WatchApplicationsWithPendingCharms indicates an expected call of WatchApplicationsWithPendingCharms.
func (mr *MockStateBackendMockRecorder) WatchApplicationsWithPendingCharms() *MockStateBackendWatchApplicationsWithPendingCharmsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchApplicationsWithPendingCharms", reflect.TypeOf((*MockStateBackend)(nil).WatchApplicationsWithPendingCharms))
	return &MockStateBackendWatchApplicationsWithPendingCharmsCall{Call: call}
}

// MockStateBackendWatchApplicationsWithPendingCharmsCall wrap *gomock.Call
type MockStateBackendWatchApplicationsWithPendingCharmsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStateBackendWatchApplicationsWithPendingCharmsCall) Return(arg0 state.StringsWatcher) *MockStateBackendWatchApplicationsWithPendingCharmsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStateBackendWatchApplicationsWithPendingCharmsCall) Do(f func() state.StringsWatcher) *MockStateBackendWatchApplicationsWithPendingCharmsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStateBackendWatchApplicationsWithPendingCharmsCall) DoAndReturn(f func() state.StringsWatcher) *MockStateBackendWatchApplicationsWithPendingCharmsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockApplication is a mock of Application interface.
type MockApplication struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationMockRecorder
}

// MockApplicationMockRecorder is the mock recorder for MockApplication.
type MockApplicationMockRecorder struct {
	mock *MockApplication
}

// NewMockApplication creates a new mock instance.
func NewMockApplication(ctrl *gomock.Controller) *MockApplication {
	mock := &MockApplication{ctrl: ctrl}
	mock.recorder = &MockApplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder {
	return m.recorder
}

// Charm mocks base method.
func (m *MockApplication) Charm() (charmdownloader.Charm, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Charm")
	ret0, _ := ret[0].(charmdownloader.Charm)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Charm indicates an expected call of Charm.
func (mr *MockApplicationMockRecorder) Charm() *MockApplicationCharmCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Charm", reflect.TypeOf((*MockApplication)(nil).Charm))
	return &MockApplicationCharmCall{Call: call}
}

// MockApplicationCharmCall wrap *gomock.Call
type MockApplicationCharmCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockApplicationCharmCall) Return(arg0 charmdownloader.Charm, arg1 bool, arg2 error) *MockApplicationCharmCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockApplicationCharmCall) Do(f func() (charmdownloader.Charm, bool, error)) *MockApplicationCharmCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockApplicationCharmCall) DoAndReturn(f func() (charmdownloader.Charm, bool, error)) *MockApplicationCharmCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CharmOrigin mocks base method.
func (m *MockApplication) CharmOrigin() *charm.Origin {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CharmOrigin")
	ret0, _ := ret[0].(*charm.Origin)
	return ret0
}

// CharmOrigin indicates an expected call of CharmOrigin.
func (mr *MockApplicationMockRecorder) CharmOrigin() *MockApplicationCharmOriginCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CharmOrigin", reflect.TypeOf((*MockApplication)(nil).CharmOrigin))
	return &MockApplicationCharmOriginCall{Call: call}
}

// MockApplicationCharmOriginCall wrap *gomock.Call
type MockApplicationCharmOriginCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockApplicationCharmOriginCall) Return(arg0 *charm.Origin) *MockApplicationCharmOriginCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockApplicationCharmOriginCall) Do(f func() *charm.Origin) *MockApplicationCharmOriginCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockApplicationCharmOriginCall) DoAndReturn(f func() *charm.Origin) *MockApplicationCharmOriginCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CharmPendingToBeDownloaded mocks base method.
func (m *MockApplication) CharmPendingToBeDownloaded() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CharmPendingToBeDownloaded")
	ret0, _ := ret[0].(bool)
	return ret0
}

// CharmPendingToBeDownloaded indicates an expected call of CharmPendingToBeDownloaded.
func (mr *MockApplicationMockRecorder) CharmPendingToBeDownloaded() *MockApplicationCharmPendingToBeDownloadedCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CharmPendingToBeDownloaded", reflect.TypeOf((*MockApplication)(nil).CharmPendingToBeDownloaded))
	return &MockApplicationCharmPendingToBeDownloadedCall{Call: call}
}

// MockApplicationCharmPendingToBeDownloadedCall wrap *gomock.Call
type MockApplicationCharmPendingToBeDownloadedCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockApplicationCharmPendingToBeDownloadedCall) Return(arg0 bool) *MockApplicationCharmPendingToBeDownloadedCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockApplicationCharmPendingToBeDownloadedCall) Do(f func() bool) *MockApplicationCharmPendingToBeDownloadedCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockApplicationCharmPendingToBeDownloadedCall) DoAndReturn(f func() bool) *MockApplicationCharmPendingToBeDownloadedCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetDownloadedIDAndHash mocks base method.
func (m *MockApplication) SetDownloadedIDAndHash(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDownloadedIDAndHash", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDownloadedIDAndHash indicates an expected call of SetDownloadedIDAndHash.
func (mr *MockApplicationMockRecorder) SetDownloadedIDAndHash(arg0, arg1 any) *MockApplicationSetDownloadedIDAndHashCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDownloadedIDAndHash", reflect.TypeOf((*MockApplication)(nil).SetDownloadedIDAndHash), arg0, arg1)
	return &MockApplicationSetDownloadedIDAndHashCall{Call: call}
}

// MockApplicationSetDownloadedIDAndHashCall wrap *gomock.Call
type MockApplicationSetDownloadedIDAndHashCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockApplicationSetDownloadedIDAndHashCall) Return(arg0 error) *MockApplicationSetDownloadedIDAndHashCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockApplicationSetDownloadedIDAndHashCall) Do(f func(string, string) error) *MockApplicationSetDownloadedIDAndHashCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockApplicationSetDownloadedIDAndHashCall) DoAndReturn(f func(string, string) error) *MockApplicationSetDownloadedIDAndHashCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetStatus mocks base method.
func (m *MockApplication) SetStatus(arg0 status.StatusInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStatus", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStatus indicates an expected call of SetStatus.
func (mr *MockApplicationMockRecorder) SetStatus(arg0 any) *MockApplicationSetStatusCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStatus", reflect.TypeOf((*MockApplication)(nil).SetStatus), arg0)
	return &MockApplicationSetStatusCall{Call: call}
}

// MockApplicationSetStatusCall wrap *gomock.Call
type MockApplicationSetStatusCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockApplicationSetStatusCall) Return(arg0 error) *MockApplicationSetStatusCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockApplicationSetStatusCall) Do(f func(status.StatusInfo) error) *MockApplicationSetStatusCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockApplicationSetStatusCall) DoAndReturn(f func(status.StatusInfo) error) *MockApplicationSetStatusCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockCharm is a mock of Charm interface.
type MockCharm struct {
	ctrl     *gomock.Controller
	recorder *MockCharmMockRecorder
}

// MockCharmMockRecorder is the mock recorder for MockCharm.
type MockCharmMockRecorder struct {
	mock *MockCharm
}

// NewMockCharm creates a new mock instance.
func NewMockCharm(ctrl *gomock.Controller) *MockCharm {
	mock := &MockCharm{ctrl: ctrl}
	mock.recorder = &MockCharmMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCharm) EXPECT() *MockCharmMockRecorder {
	return m.recorder
}

// URL mocks base method.
func (m *MockCharm) URL() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "URL")
	ret0, _ := ret[0].(string)
	return ret0
}

// URL indicates an expected call of URL.
func (mr *MockCharmMockRecorder) URL() *MockCharmURLCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "URL", reflect.TypeOf((*MockCharm)(nil).URL))
	return &MockCharmURLCall{Call: call}
}

// MockCharmURLCall wrap *gomock.Call
type MockCharmURLCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCharmURLCall) Return(arg0 string) *MockCharmURLCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCharmURLCall) Do(f func() string) *MockCharmURLCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCharmURLCall) DoAndReturn(f func() string) *MockCharmURLCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockDownloader is a mock of Downloader interface.
type MockDownloader struct {
	ctrl     *gomock.Controller
	recorder *MockDownloaderMockRecorder
}

// MockDownloaderMockRecorder is the mock recorder for MockDownloader.
type MockDownloaderMockRecorder struct {
	mock *MockDownloader
}

// NewMockDownloader creates a new mock instance.
func NewMockDownloader(ctrl *gomock.Controller) *MockDownloader {
	mock := &MockDownloader{ctrl: ctrl}
	mock.recorder = &MockDownloaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDownloader) EXPECT() *MockDownloaderMockRecorder {
	return m.recorder
}

// DownloadAndStore mocks base method.
func (m *MockDownloader) DownloadAndStore(arg0 context.Context, arg1 *charm0.URL, arg2 charm.Origin, arg3 bool) (charm.Origin, charm0.Charm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadAndStore", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(charm.Origin)
	ret1, _ := ret[1].(charm0.Charm)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DownloadAndStore indicates an expected call of DownloadAndStore.
func (mr *MockDownloaderMockRecorder) DownloadAndStore(arg0, arg1, arg2, arg3 any) *MockDownloaderDownloadAndStoreCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadAndStore", reflect.TypeOf((*MockDownloader)(nil).DownloadAndStore), arg0, arg1, arg2, arg3)
	return &MockDownloaderDownloadAndStoreCall{Call: call}
}

// MockDownloaderDownloadAndStoreCall wrap *gomock.Call
type MockDownloaderDownloadAndStoreCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockDownloaderDownloadAndStoreCall) Return(arg0 charm.Origin, arg1 charm0.Charm, arg2 error) *MockDownloaderDownloadAndStoreCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockDownloaderDownloadAndStoreCall) Do(f func(context.Context, *charm0.URL, charm.Origin, bool) (charm.Origin, charm0.Charm, error)) *MockDownloaderDownloadAndStoreCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockDownloaderDownloadAndStoreCall) DoAndReturn(f func(context.Context, *charm0.URL, charm.Origin, bool) (charm.Origin, charm0.Charm, error)) *MockDownloaderDownloadAndStoreCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockAuthChecker is a mock of AuthChecker interface.
type MockAuthChecker struct {
	ctrl     *gomock.Controller
	recorder *MockAuthCheckerMockRecorder
}

// MockAuthCheckerMockRecorder is the mock recorder for MockAuthChecker.
type MockAuthCheckerMockRecorder struct {
	mock *MockAuthChecker
}

// NewMockAuthChecker creates a new mock instance.
func NewMockAuthChecker(ctrl *gomock.Controller) *MockAuthChecker {
	mock := &MockAuthChecker{ctrl: ctrl}
	mock.recorder = &MockAuthCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthChecker) EXPECT() *MockAuthCheckerMockRecorder {
	return m.recorder
}

// AuthController mocks base method.
func (m *MockAuthChecker) AuthController() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthController")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthController indicates an expected call of AuthController.
func (mr *MockAuthCheckerMockRecorder) AuthController() *MockAuthCheckerAuthControllerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthController", reflect.TypeOf((*MockAuthChecker)(nil).AuthController))
	return &MockAuthCheckerAuthControllerCall{Call: call}
}

// MockAuthCheckerAuthControllerCall wrap *gomock.Call
type MockAuthCheckerAuthControllerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAuthCheckerAuthControllerCall) Return(arg0 bool) *MockAuthCheckerAuthControllerCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAuthCheckerAuthControllerCall) Do(f func() bool) *MockAuthCheckerAuthControllerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAuthCheckerAuthControllerCall) DoAndReturn(f func() bool) *MockAuthCheckerAuthControllerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockResourcesBackend is a mock of ResourcesBackend interface.
type MockResourcesBackend struct {
	ctrl     *gomock.Controller
	recorder *MockResourcesBackendMockRecorder
}

// MockResourcesBackendMockRecorder is the mock recorder for MockResourcesBackend.
type MockResourcesBackendMockRecorder struct {
	mock *MockResourcesBackend
}

// NewMockResourcesBackend creates a new mock instance.
func NewMockResourcesBackend(ctrl *gomock.Controller) *MockResourcesBackend {
	mock := &MockResourcesBackend{ctrl: ctrl}
	mock.recorder = &MockResourcesBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourcesBackend) EXPECT() *MockResourcesBackendMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *MockResourcesBackend) Register(arg0 worker.Worker) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockResourcesBackendMockRecorder) Register(arg0 any) *MockResourcesBackendRegisterCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockResourcesBackend)(nil).Register), arg0)
	return &MockResourcesBackendRegisterCall{Call: call}
}

// MockResourcesBackendRegisterCall wrap *gomock.Call
type MockResourcesBackendRegisterCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockResourcesBackendRegisterCall) Return(arg0 string) *MockResourcesBackendRegisterCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockResourcesBackendRegisterCall) Do(f func(worker.Worker) string) *MockResourcesBackendRegisterCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockResourcesBackendRegisterCall) DoAndReturn(f func(worker.Worker) string) *MockResourcesBackendRegisterCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
