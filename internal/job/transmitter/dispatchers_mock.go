// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Goboolean/core-system.worker/internal/job/transmitter (interfaces: OrderEventDispatcher,AnnotationDispatcher)
//
// Generated by this command:
//
//	mockgen -destination=dispatchers_mock.go -package=transmitter --build_flags=--mod=mod . OrderEventDispatcher,AnnotationDispatcher
//

// Package transmitter is a generated GoMock package.
package transmitter

import (
	context "context"
	reflect "reflect"

	model "github.com/Goboolean/core-system.worker/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockOrderEventDispatcher is a mock of OrderEventDispatcher interface.
type MockOrderEventDispatcher struct {
	ctrl     *gomock.Controller
	recorder *MockOrderEventDispatcherMockRecorder
}

// MockOrderEventDispatcherMockRecorder is the mock recorder for MockOrderEventDispatcher.
type MockOrderEventDispatcherMockRecorder struct {
	mock *MockOrderEventDispatcher
}

// NewMockOrderEventDispatcher creates a new mock instance.
func NewMockOrderEventDispatcher(ctrl *gomock.Controller) *MockOrderEventDispatcher {
	mock := &MockOrderEventDispatcher{ctrl: ctrl}
	mock.recorder = &MockOrderEventDispatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderEventDispatcher) EXPECT() *MockOrderEventDispatcherMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockOrderEventDispatcher) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockOrderEventDispatcherMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockOrderEventDispatcher)(nil).Close))
}

// Dispatch mocks base method.
func (m *MockOrderEventDispatcher) Dispatch(arg0 *model.OrderEvent) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dispatch", arg0)
}

// Dispatch indicates an expected call of Dispatch.
func (mr *MockOrderEventDispatcherMockRecorder) Dispatch(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispatch", reflect.TypeOf((*MockOrderEventDispatcher)(nil).Dispatch), arg0)
}

// Flush mocks base method.
func (m *MockOrderEventDispatcher) Flush(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Flush", arg0)
}

// Flush indicates an expected call of Flush.
func (mr *MockOrderEventDispatcherMockRecorder) Flush(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockOrderEventDispatcher)(nil).Flush), arg0)
}

// MockAnnotationDispatcher is a mock of AnnotationDispatcher interface.
type MockAnnotationDispatcher struct {
	ctrl     *gomock.Controller
	recorder *MockAnnotationDispatcherMockRecorder
}

// MockAnnotationDispatcherMockRecorder is the mock recorder for MockAnnotationDispatcher.
type MockAnnotationDispatcherMockRecorder struct {
	mock *MockAnnotationDispatcher
}

// NewMockAnnotationDispatcher creates a new mock instance.
func NewMockAnnotationDispatcher(ctrl *gomock.Controller) *MockAnnotationDispatcher {
	mock := &MockAnnotationDispatcher{ctrl: ctrl}
	mock.recorder = &MockAnnotationDispatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAnnotationDispatcher) EXPECT() *MockAnnotationDispatcherMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockAnnotationDispatcher) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockAnnotationDispatcherMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAnnotationDispatcher)(nil).Close))
}

// Dispatch mocks base method.
func (m *MockAnnotationDispatcher) Dispatch(arg0 any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dispatch", arg0)
}

// Dispatch indicates an expected call of Dispatch.
func (mr *MockAnnotationDispatcherMockRecorder) Dispatch(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispatch", reflect.TypeOf((*MockAnnotationDispatcher)(nil).Dispatch), arg0)
}

// Flush mocks base method.
func (m *MockAnnotationDispatcher) Flush(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Flush", arg0)
}

// Flush indicates an expected call of Flush.
func (mr *MockAnnotationDispatcherMockRecorder) Flush(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockAnnotationDispatcher)(nil).Flush), arg0)
}