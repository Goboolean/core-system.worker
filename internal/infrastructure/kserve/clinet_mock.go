// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Goboolean/core-system.worker/internal/infrastructure/kserve (interfaces: Client)
//
// Generated by this command:
//
//	mockgen -destination=clinet_mock.go -package=kserve --build_flags=--mod=mod . Client
//

// Package kserve is a generated GoMock package.
package kserve

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// RequestInference mocks base method.
func (m *MockClient) RequestInference(arg0 context.Context, arg1 []int, arg2 []float32) ([]float32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestInference", arg0, arg1, arg2)
	ret0, _ := ret[0].([]float32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RequestInference indicates an expected call of RequestInference.
func (mr *MockClientMockRecorder) RequestInference(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestInference", reflect.TypeOf((*MockClient)(nil).RequestInference), arg0, arg1, arg2)
}

// SetModelName mocks base method.
func (m *MockClient) SetModelName(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetModelName", arg0)
}

// SetModelName indicates an expected call of SetModelName.
func (mr *MockClientMockRecorder) SetModelName(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetModelName", reflect.TypeOf((*MockClient)(nil).SetModelName), arg0)
}