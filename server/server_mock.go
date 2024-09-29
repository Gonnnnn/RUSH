// Code generated by MockGen. DO NOT EDIT.
// Source: server.go
//
// Generated by this command:
//
//	mockgen -source=server.go -destination=server_mock.go -package=server
//

// Package server is a generated GoMock package.
package server

import (
	reflect "reflect"
	attendance "rush/attendance"
	auth "rush/auth"
	permission "rush/permission"
	session "rush/session"
	user "rush/user"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockoauthClient is a mock of oauthClient interface.
type MockoauthClient struct {
	ctrl     *gomock.Controller
	recorder *MockoauthClientMockRecorder
}

// MockoauthClientMockRecorder is the mock recorder for MockoauthClient.
type MockoauthClientMockRecorder struct {
	mock *MockoauthClient
}

// NewMockoauthClient creates a new mock instance.
func NewMockoauthClient(ctrl *gomock.Controller) *MockoauthClient {
	mock := &MockoauthClient{ctrl: ctrl}
	mock.recorder = &MockoauthClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockoauthClient) EXPECT() *MockoauthClientMockRecorder {
	return m.recorder
}

// GetEmail mocks base method.
func (m *MockoauthClient) GetEmail(token string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmail", token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmail indicates an expected call of GetEmail.
func (mr *MockoauthClientMockRecorder) GetEmail(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmail", reflect.TypeOf((*MockoauthClient)(nil).GetEmail), token)
}

// MockauthHandler is a mock of authHandler interface.
type MockauthHandler struct {
	ctrl     *gomock.Controller
	recorder *MockauthHandlerMockRecorder
}

// MockauthHandlerMockRecorder is the mock recorder for MockauthHandler.
type MockauthHandlerMockRecorder struct {
	mock *MockauthHandler
}

// NewMockauthHandler creates a new mock instance.
func NewMockauthHandler(ctrl *gomock.Controller) *MockauthHandler {
	mock := &MockauthHandler{ctrl: ctrl}
	mock.recorder = &MockauthHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockauthHandler) EXPECT() *MockauthHandlerMockRecorder {
	return m.recorder
}

// GetSession mocks base method.
func (m *MockauthHandler) GetSession(token string) (auth.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", token)
	ret0, _ := ret[0].(auth.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockauthHandlerMockRecorder) GetSession(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockauthHandler)(nil).GetSession), token)
}

// SignIn mocks base method.
func (m *MockauthHandler) SignIn(userId string, role permission.Role) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", userId, role)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockauthHandlerMockRecorder) SignIn(userId, role any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockauthHandler)(nil).SignIn), userId, role)
}

// MockuserRepo is a mock of userRepo interface.
type MockuserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockuserRepoMockRecorder
}

// MockuserRepoMockRecorder is the mock recorder for MockuserRepo.
type MockuserRepoMockRecorder struct {
	mock *MockuserRepo
}

// NewMockuserRepo creates a new mock instance.
func NewMockuserRepo(ctrl *gomock.Controller) *MockuserRepo {
	mock := &MockuserRepo{ctrl: ctrl}
	mock.recorder = &MockuserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuserRepo) EXPECT() *MockuserRepoMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockuserRepo) Get(id string) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockuserRepoMockRecorder) Get(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockuserRepo)(nil).Get), id)
}

// GetAll mocks base method.
func (m *MockuserRepo) GetAll() ([]user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockuserRepoMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockuserRepo)(nil).GetAll))
}

// GetAllByExternalNames mocks base method.
func (m *MockuserRepo) GetAllByExternalNames(externalNames []string) ([]user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByExternalNames", externalNames)
	ret0, _ := ret[0].([]user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByExternalNames indicates an expected call of GetAllByExternalNames.
func (mr *MockuserRepoMockRecorder) GetAllByExternalNames(externalNames any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByExternalNames", reflect.TypeOf((*MockuserRepo)(nil).GetAllByExternalNames), externalNames)
}

// GetByEmail mocks base method.
func (m *MockuserRepo) GetByEmail(email string) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", email)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockuserRepoMockRecorder) GetByEmail(email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockuserRepo)(nil).GetByEmail), email)
}

// List mocks base method.
func (m *MockuserRepo) List(offset, pageSize int) (*user.ListResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", offset, pageSize)
	ret0, _ := ret[0].(*user.ListResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockuserRepoMockRecorder) List(offset, pageSize any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockuserRepo)(nil).List), offset, pageSize)
}

// MockuserAdder is a mock of userAdder interface.
type MockuserAdder struct {
	ctrl     *gomock.Controller
	recorder *MockuserAdderMockRecorder
}

// MockuserAdderMockRecorder is the mock recorder for MockuserAdder.
type MockuserAdderMockRecorder struct {
	mock *MockuserAdder
}

// NewMockuserAdder creates a new mock instance.
func NewMockuserAdder(ctrl *gomock.Controller) *MockuserAdder {
	mock := &MockuserAdder{ctrl: ctrl}
	mock.recorder = &MockuserAdderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuserAdder) EXPECT() *MockuserAdderMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockuserAdder) Add(name string, generation float64, isActive bool, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", name, generation, isActive, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockuserAdderMockRecorder) Add(name, generation, isActive, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockuserAdder)(nil).Add), name, generation, isActive, email)
}

// MocksessionRepo is a mock of sessionRepo interface.
type MocksessionRepo struct {
	ctrl     *gomock.Controller
	recorder *MocksessionRepoMockRecorder
}

// MocksessionRepoMockRecorder is the mock recorder for MocksessionRepo.
type MocksessionRepoMockRecorder struct {
	mock *MocksessionRepo
}

// NewMocksessionRepo creates a new mock instance.
func NewMocksessionRepo(ctrl *gomock.Controller) *MocksessionRepo {
	mock := &MocksessionRepo{ctrl: ctrl}
	mock.recorder = &MocksessionRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocksessionRepo) EXPECT() *MocksessionRepoMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MocksessionRepo) Add(name, description, createdBy string, startsAt time.Time, score int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", name, description, createdBy, startsAt, score)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MocksessionRepoMockRecorder) Add(name, description, createdBy, startsAt, score any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MocksessionRepo)(nil).Add), name, description, createdBy, startsAt, score)
}

// Delete mocks base method.
func (m *MocksessionRepo) Delete(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MocksessionRepoMockRecorder) Delete(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MocksessionRepo)(nil).Delete), id)
}

// Get mocks base method.
func (m *MocksessionRepo) Get(id string) (session.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MocksessionRepoMockRecorder) Get(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MocksessionRepo)(nil).Get), id)
}

// GetAll mocks base method.
func (m *MocksessionRepo) GetAll() ([]session.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MocksessionRepoMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MocksessionRepo)(nil).GetAll))
}

// List mocks base method.
func (m *MocksessionRepo) List(offset, pageSize int) (*session.ListResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", offset, pageSize)
	ret0, _ := ret[0].(*session.ListResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MocksessionRepoMockRecorder) List(offset, pageSize any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MocksessionRepo)(nil).List), offset, pageSize)
}

// Update mocks base method.
func (m *MocksessionRepo) Update(id string, updateForm session.UpdateForm) (session.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", id, updateForm)
	ret0, _ := ret[0].(session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MocksessionRepoMockRecorder) Update(id, updateForm any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MocksessionRepo)(nil).Update), id, updateForm)
}

// MockattendanceFormHandler is a mock of attendanceFormHandler interface.
type MockattendanceFormHandler struct {
	ctrl     *gomock.Controller
	recorder *MockattendanceFormHandlerMockRecorder
}

// MockattendanceFormHandlerMockRecorder is the mock recorder for MockattendanceFormHandler.
type MockattendanceFormHandlerMockRecorder struct {
	mock *MockattendanceFormHandler
}

// NewMockattendanceFormHandler creates a new mock instance.
func NewMockattendanceFormHandler(ctrl *gomock.Controller) *MockattendanceFormHandler {
	mock := &MockattendanceFormHandler{ctrl: ctrl}
	mock.recorder = &MockattendanceFormHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockattendanceFormHandler) EXPECT() *MockattendanceFormHandlerMockRecorder {
	return m.recorder
}

// GenerateForm mocks base method.
func (m *MockattendanceFormHandler) GenerateForm(title, description string, users []user.User) (attendance.Form, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateForm", title, description, users)
	ret0, _ := ret[0].(attendance.Form)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateForm indicates an expected call of GenerateForm.
func (mr *MockattendanceFormHandlerMockRecorder) GenerateForm(title, description, users any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateForm", reflect.TypeOf((*MockattendanceFormHandler)(nil).GenerateForm), title, description, users)
}

// GetSubmissions mocks base method.
func (m *MockattendanceFormHandler) GetSubmissions(formId string) ([]attendance.FormSubmission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubmissions", formId)
	ret0, _ := ret[0].([]attendance.FormSubmission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubmissions indicates an expected call of GetSubmissions.
func (mr *MockattendanceFormHandlerMockRecorder) GetSubmissions(formId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubmissions", reflect.TypeOf((*MockattendanceFormHandler)(nil).GetSubmissions), formId)
}

// MockattendanceRepo is a mock of attendanceRepo interface.
type MockattendanceRepo struct {
	ctrl     *gomock.Controller
	recorder *MockattendanceRepoMockRecorder
}

// MockattendanceRepoMockRecorder is the mock recorder for MockattendanceRepo.
type MockattendanceRepoMockRecorder struct {
	mock *MockattendanceRepo
}

// NewMockattendanceRepo creates a new mock instance.
func NewMockattendanceRepo(ctrl *gomock.Controller) *MockattendanceRepo {
	mock := &MockattendanceRepo{ctrl: ctrl}
	mock.recorder = &MockattendanceRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockattendanceRepo) EXPECT() *MockattendanceRepoMockRecorder {
	return m.recorder
}

// BulkInsert mocks base method.
func (m *MockattendanceRepo) BulkInsert(requests []attendance.AddAttendanceReq) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BulkInsert", requests)
	ret0, _ := ret[0].(error)
	return ret0
}

// BulkInsert indicates an expected call of BulkInsert.
func (mr *MockattendanceRepoMockRecorder) BulkInsert(requests any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkInsert", reflect.TypeOf((*MockattendanceRepo)(nil).BulkInsert), requests)
}

// FindByUserId mocks base method.
func (m *MockattendanceRepo) FindByUserId(userId string) ([]attendance.Attendance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserId", userId)
	ret0, _ := ret[0].([]attendance.Attendance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserId indicates an expected call of FindByUserId.
func (mr *MockattendanceRepoMockRecorder) FindByUserId(userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserId", reflect.TypeOf((*MockattendanceRepo)(nil).FindByUserId), userId)
}

// GetAll mocks base method.
func (m *MockattendanceRepo) GetAll() ([]attendance.Attendance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]attendance.Attendance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockattendanceRepoMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockattendanceRepo)(nil).GetAll))
}
