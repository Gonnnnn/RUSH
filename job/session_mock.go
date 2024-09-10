// Code generated by MockGen. DO NOT EDIT.
// Source: session.go
//
// Generated by this command:
//
//	mockgen -source=session.go -destination=session_mock.go -package=job
//

// Package job is a generated GoMock package.
package job

import (
	reflect "reflect"
	session "rush/session"

	gomock "go.uber.org/mock/gomock"
)

// MocksessionCloser is a mock of sessionCloser interface.
type MocksessionCloser struct {
	ctrl     *gomock.Controller
	recorder *MocksessionCloserMockRecorder
}

// MocksessionCloserMockRecorder is the mock recorder for MocksessionCloser.
type MocksessionCloserMockRecorder struct {
	mock *MocksessionCloser
}

// NewMocksessionCloser creates a new mock instance.
func NewMocksessionCloser(ctrl *gomock.Controller) *MocksessionCloser {
	mock := &MocksessionCloser{ctrl: ctrl}
	mock.recorder = &MocksessionCloserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocksessionCloser) EXPECT() *MocksessionCloserMockRecorder {
	return m.recorder
}

// CloseSession mocks base method.
func (m *MocksessionCloser) CloseSession(sessionId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSession", sessionId)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSession indicates an expected call of CloseSession.
func (mr *MocksessionCloserMockRecorder) CloseSession(sessionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSession", reflect.TypeOf((*MocksessionCloser)(nil).CloseSession), sessionId)
}

// MocksessionGetter is a mock of sessionGetter interface.
type MocksessionGetter struct {
	ctrl     *gomock.Controller
	recorder *MocksessionGetterMockRecorder
}

// MocksessionGetterMockRecorder is the mock recorder for MocksessionGetter.
type MocksessionGetterMockRecorder struct {
	mock *MocksessionGetter
}

// NewMocksessionGetter creates a new mock instance.
func NewMocksessionGetter(ctrl *gomock.Controller) *MocksessionGetter {
	mock := &MocksessionGetter{ctrl: ctrl}
	mock.recorder = &MocksessionGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocksessionGetter) EXPECT() *MocksessionGetterMockRecorder {
	return m.recorder
}

// GetOpenSessions mocks base method.
func (m *MocksessionGetter) GetOpenSessions() ([]*session.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenSessions")
	ret0, _ := ret[0].([]*session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenSessions indicates an expected call of GetOpenSessions.
func (mr *MocksessionGetterMockRecorder) GetOpenSessions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenSessions", reflect.TypeOf((*MocksessionGetter)(nil).GetOpenSessions))
}
