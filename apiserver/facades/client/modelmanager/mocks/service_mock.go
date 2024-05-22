// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/modelmanager (interfaces: AccessService,SecretBackendService,ModelService)
//
// Generated by this command:
//
//	mockgen -typed -package mocks -destination mocks/service_mock.go github.com/juju/juju/apiserver/facades/client/modelmanager AccessService,SecretBackendService,ModelService
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	credential "github.com/juju/juju/core/credential"
	model "github.com/juju/juju/core/model"
	permission "github.com/juju/juju/core/permission"
	user "github.com/juju/juju/core/user"
	access "github.com/juju/juju/domain/access"
	model0 "github.com/juju/juju/domain/model"
	service "github.com/juju/juju/domain/secretbackend/service"
	gomock "go.uber.org/mock/gomock"
)

// MockAccessService is a mock of AccessService interface.
type MockAccessService struct {
	ctrl     *gomock.Controller
	recorder *MockAccessServiceMockRecorder
}

// MockAccessServiceMockRecorder is the mock recorder for MockAccessService.
type MockAccessServiceMockRecorder struct {
	mock *MockAccessService
}

// NewMockAccessService creates a new mock instance.
func NewMockAccessService(ctrl *gomock.Controller) *MockAccessService {
	mock := &MockAccessService{ctrl: ctrl}
	mock.recorder = &MockAccessServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccessService) EXPECT() *MockAccessServiceMockRecorder {
	return m.recorder
}

// GetUserByName mocks base method.
func (m *MockAccessService) GetUserByName(arg0 context.Context, arg1 string) (user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByName", arg0, arg1)
	ret0, _ := ret[0].(user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByName indicates an expected call of GetUserByName.
func (mr *MockAccessServiceMockRecorder) GetUserByName(arg0, arg1 any) *MockAccessServiceGetUserByNameCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByName", reflect.TypeOf((*MockAccessService)(nil).GetUserByName), arg0, arg1)
	return &MockAccessServiceGetUserByNameCall{Call: call}
}

// MockAccessServiceGetUserByNameCall wrap *gomock.Call
type MockAccessServiceGetUserByNameCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceGetUserByNameCall) Return(arg0 user.User, arg1 error) *MockAccessServiceGetUserByNameCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceGetUserByNameCall) Do(f func(context.Context, string) (user.User, error)) *MockAccessServiceGetUserByNameCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceGetUserByNameCall) DoAndReturn(f func(context.Context, string) (user.User, error)) *MockAccessServiceGetUserByNameCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ReadUserAccessLevelForTarget mocks base method.
func (m *MockAccessService) ReadUserAccessLevelForTarget(arg0 context.Context, arg1 string, arg2 permission.ID) (permission.Access, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadUserAccessLevelForTarget", arg0, arg1, arg2)
	ret0, _ := ret[0].(permission.Access)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadUserAccessLevelForTarget indicates an expected call of ReadUserAccessLevelForTarget.
func (mr *MockAccessServiceMockRecorder) ReadUserAccessLevelForTarget(arg0, arg1, arg2 any) *MockAccessServiceReadUserAccessLevelForTargetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadUserAccessLevelForTarget", reflect.TypeOf((*MockAccessService)(nil).ReadUserAccessLevelForTarget), arg0, arg1, arg2)
	return &MockAccessServiceReadUserAccessLevelForTargetCall{Call: call}
}

// MockAccessServiceReadUserAccessLevelForTargetCall wrap *gomock.Call
type MockAccessServiceReadUserAccessLevelForTargetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceReadUserAccessLevelForTargetCall) Return(arg0 permission.Access, arg1 error) *MockAccessServiceReadUserAccessLevelForTargetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceReadUserAccessLevelForTargetCall) Do(f func(context.Context, string, permission.ID) (permission.Access, error)) *MockAccessServiceReadUserAccessLevelForTargetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceReadUserAccessLevelForTargetCall) DoAndReturn(f func(context.Context, string, permission.ID) (permission.Access, error)) *MockAccessServiceReadUserAccessLevelForTargetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdatePermission mocks base method.
func (m *MockAccessService) UpdatePermission(arg0 context.Context, arg1 access.UpdatePermissionArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePermission", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePermission indicates an expected call of UpdatePermission.
func (mr *MockAccessServiceMockRecorder) UpdatePermission(arg0, arg1 any) *MockAccessServiceUpdatePermissionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePermission", reflect.TypeOf((*MockAccessService)(nil).UpdatePermission), arg0, arg1)
	return &MockAccessServiceUpdatePermissionCall{Call: call}
}

// MockAccessServiceUpdatePermissionCall wrap *gomock.Call
type MockAccessServiceUpdatePermissionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceUpdatePermissionCall) Return(arg0 error) *MockAccessServiceUpdatePermissionCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceUpdatePermissionCall) Do(f func(context.Context, access.UpdatePermissionArgs) error) *MockAccessServiceUpdatePermissionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceUpdatePermissionCall) DoAndReturn(f func(context.Context, access.UpdatePermissionArgs) error) *MockAccessServiceUpdatePermissionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockSecretBackendService is a mock of SecretBackendService interface.
type MockSecretBackendService struct {
	ctrl     *gomock.Controller
	recorder *MockSecretBackendServiceMockRecorder
}

// MockSecretBackendServiceMockRecorder is the mock recorder for MockSecretBackendService.
type MockSecretBackendServiceMockRecorder struct {
	mock *MockSecretBackendService
}

// NewMockSecretBackendService creates a new mock instance.
func NewMockSecretBackendService(ctrl *gomock.Controller) *MockSecretBackendService {
	mock := &MockSecretBackendService{ctrl: ctrl}
	mock.recorder = &MockSecretBackendServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretBackendService) EXPECT() *MockSecretBackendServiceMockRecorder {
	return m.recorder
}

// BackendSummaryInfoForModel mocks base method.
func (m *MockSecretBackendService) BackendSummaryInfoForModel(arg0 context.Context, arg1 model.UUID) ([]*service.SecretBackendInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BackendSummaryInfoForModel", arg0, arg1)
	ret0, _ := ret[0].([]*service.SecretBackendInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BackendSummaryInfoForModel indicates an expected call of BackendSummaryInfoForModel.
func (mr *MockSecretBackendServiceMockRecorder) BackendSummaryInfoForModel(arg0, arg1 any) *MockSecretBackendServiceBackendSummaryInfoForModelCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BackendSummaryInfoForModel", reflect.TypeOf((*MockSecretBackendService)(nil).BackendSummaryInfoForModel), arg0, arg1)
	return &MockSecretBackendServiceBackendSummaryInfoForModelCall{Call: call}
}

// MockSecretBackendServiceBackendSummaryInfoForModelCall wrap *gomock.Call
type MockSecretBackendServiceBackendSummaryInfoForModelCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSecretBackendServiceBackendSummaryInfoForModelCall) Return(arg0 []*service.SecretBackendInfo, arg1 error) *MockSecretBackendServiceBackendSummaryInfoForModelCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSecretBackendServiceBackendSummaryInfoForModelCall) Do(f func(context.Context, model.UUID) ([]*service.SecretBackendInfo, error)) *MockSecretBackendServiceBackendSummaryInfoForModelCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSecretBackendServiceBackendSummaryInfoForModelCall) DoAndReturn(f func(context.Context, model.UUID) ([]*service.SecretBackendInfo, error)) *MockSecretBackendServiceBackendSummaryInfoForModelCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockModelService is a mock of ModelService interface.
type MockModelService struct {
	ctrl     *gomock.Controller
	recorder *MockModelServiceMockRecorder
}

// MockModelServiceMockRecorder is the mock recorder for MockModelService.
type MockModelServiceMockRecorder struct {
	mock *MockModelService
}

// NewMockModelService creates a new mock instance.
func NewMockModelService(ctrl *gomock.Controller) *MockModelService {
	mock := &MockModelService{ctrl: ctrl}
	mock.recorder = &MockModelServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModelService) EXPECT() *MockModelServiceMockRecorder {
	return m.recorder
}

// CreateModel mocks base method.
func (m *MockModelService) CreateModel(arg0 context.Context, arg1 model0.ModelCreationArgs) (func(context.Context) error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateModel", arg0, arg1)
	ret0, _ := ret[0].(func(context.Context) error)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateModel indicates an expected call of CreateModel.
func (mr *MockModelServiceMockRecorder) CreateModel(arg0, arg1 any) *MockModelServiceCreateModelCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateModel", reflect.TypeOf((*MockModelService)(nil).CreateModel), arg0, arg1)
	return &MockModelServiceCreateModelCall{Call: call}
}

// MockModelServiceCreateModelCall wrap *gomock.Call
type MockModelServiceCreateModelCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelServiceCreateModelCall) Return(arg0 func(context.Context) error, arg1 error) *MockModelServiceCreateModelCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelServiceCreateModelCall) Do(f func(context.Context, model0.ModelCreationArgs) (func(context.Context) error, error)) *MockModelServiceCreateModelCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelServiceCreateModelCall) DoAndReturn(f func(context.Context, model0.ModelCreationArgs) (func(context.Context) error, error)) *MockModelServiceCreateModelCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// DefaultModelCloudNameAndCredential mocks base method.
func (m *MockModelService) DefaultModelCloudNameAndCredential(arg0 context.Context) (string, credential.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DefaultModelCloudNameAndCredential", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(credential.Key)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DefaultModelCloudNameAndCredential indicates an expected call of DefaultModelCloudNameAndCredential.
func (mr *MockModelServiceMockRecorder) DefaultModelCloudNameAndCredential(arg0 any) *MockModelServiceDefaultModelCloudNameAndCredentialCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultModelCloudNameAndCredential", reflect.TypeOf((*MockModelService)(nil).DefaultModelCloudNameAndCredential), arg0)
	return &MockModelServiceDefaultModelCloudNameAndCredentialCall{Call: call}
}

// MockModelServiceDefaultModelCloudNameAndCredentialCall wrap *gomock.Call
type MockModelServiceDefaultModelCloudNameAndCredentialCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelServiceDefaultModelCloudNameAndCredentialCall) Return(arg0 string, arg1 credential.Key, arg2 error) *MockModelServiceDefaultModelCloudNameAndCredentialCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelServiceDefaultModelCloudNameAndCredentialCall) Do(f func(context.Context) (string, credential.Key, error)) *MockModelServiceDefaultModelCloudNameAndCredentialCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelServiceDefaultModelCloudNameAndCredentialCall) DoAndReturn(f func(context.Context) (string, credential.Key, error)) *MockModelServiceDefaultModelCloudNameAndCredentialCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// DeleteModel mocks base method.
func (m *MockModelService) DeleteModel(arg0 context.Context, arg1 model.UUID, arg2 ...model0.DeleteModelOption) error {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteModel", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteModel indicates an expected call of DeleteModel.
func (mr *MockModelServiceMockRecorder) DeleteModel(arg0, arg1 any, arg2 ...any) *MockModelServiceDeleteModelCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteModel", reflect.TypeOf((*MockModelService)(nil).DeleteModel), varargs...)
	return &MockModelServiceDeleteModelCall{Call: call}
}

// MockModelServiceDeleteModelCall wrap *gomock.Call
type MockModelServiceDeleteModelCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelServiceDeleteModelCall) Return(arg0 error) *MockModelServiceDeleteModelCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelServiceDeleteModelCall) Do(f func(context.Context, model.UUID, ...model0.DeleteModelOption) error) *MockModelServiceDeleteModelCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelServiceDeleteModelCall) DoAndReturn(f func(context.Context, model.UUID, ...model0.DeleteModelOption) error) *MockModelServiceDeleteModelCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListAllModels mocks base method.
func (m *MockModelService) ListAllModels(arg0 context.Context) ([]model.Model, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllModels", arg0)
	ret0, _ := ret[0].([]model.Model)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllModels indicates an expected call of ListAllModels.
func (mr *MockModelServiceMockRecorder) ListAllModels(arg0 any) *MockModelServiceListAllModelsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllModels", reflect.TypeOf((*MockModelService)(nil).ListAllModels), arg0)
	return &MockModelServiceListAllModelsCall{Call: call}
}

// MockModelServiceListAllModelsCall wrap *gomock.Call
type MockModelServiceListAllModelsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelServiceListAllModelsCall) Return(arg0 []model.Model, arg1 error) *MockModelServiceListAllModelsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelServiceListAllModelsCall) Do(f func(context.Context) ([]model.Model, error)) *MockModelServiceListAllModelsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelServiceListAllModelsCall) DoAndReturn(f func(context.Context) ([]model.Model, error)) *MockModelServiceListAllModelsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListModelsForUser mocks base method.
func (m *MockModelService) ListModelsForUser(arg0 context.Context, arg1 user.UUID) ([]model.Model, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListModelsForUser", arg0, arg1)
	ret0, _ := ret[0].([]model.Model)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListModelsForUser indicates an expected call of ListModelsForUser.
func (mr *MockModelServiceMockRecorder) ListModelsForUser(arg0, arg1 any) *MockModelServiceListModelsForUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListModelsForUser", reflect.TypeOf((*MockModelService)(nil).ListModelsForUser), arg0, arg1)
	return &MockModelServiceListModelsForUserCall{Call: call}
}

// MockModelServiceListModelsForUserCall wrap *gomock.Call
type MockModelServiceListModelsForUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelServiceListModelsForUserCall) Return(arg0 []model.Model, arg1 error) *MockModelServiceListModelsForUserCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelServiceListModelsForUserCall) Do(f func(context.Context, user.UUID) ([]model.Model, error)) *MockModelServiceListModelsForUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelServiceListModelsForUserCall) DoAndReturn(f func(context.Context, user.UUID) ([]model.Model, error)) *MockModelServiceListModelsForUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
