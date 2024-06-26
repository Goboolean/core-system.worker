// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Goboolean/core-system.worker/internal/job/fetcher (interfaces: TradeRepository,FetchingSession,TradeStream)
//
// Generated by this command:
//
//	mockgen -destination=fetching_infra_mock.go -package=fetcher --build_flags=--mod=mod . TradeRepository,FetchingSession,TradeStream
//

// Package fetcher is a generated GoMock package.
package fetcher

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockTradeRepository is a mock of TradeRepository interface.
type MockTradeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTradeRepositoryMockRecorder
}

// MockTradeRepositoryMockRecorder is the mock recorder for MockTradeRepository.
type MockTradeRepositoryMockRecorder struct {
	mock *MockTradeRepository
}

// NewMockTradeRepository creates a new mock instance.
func NewMockTradeRepository(ctrl *gomock.Controller) *MockTradeRepository {
	mock := &MockTradeRepository{ctrl: ctrl}
	mock.recorder = &MockTradeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTradeRepository) EXPECT() *MockTradeRepositoryMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockTradeRepository) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockTradeRepositoryMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockTradeRepository)(nil).Close))
}

// SelectProduct mocks base method.
func (m *MockTradeRepository) SelectProduct(arg0, arg1, arg2 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SelectProduct", arg0, arg1, arg2)
}

// SelectProduct indicates an expected call of SelectProduct.
func (mr *MockTradeRepositoryMockRecorder) SelectProduct(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectProduct", reflect.TypeOf((*MockTradeRepository)(nil).SelectProduct), arg0, arg1, arg2)
}

// Session mocks base method.
func (m *MockTradeRepository) Session() (TradeCursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Session")
	ret0, _ := ret[0].(TradeCursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Session indicates an expected call of Session.
func (mr *MockTradeRepositoryMockRecorder) Session() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Session", reflect.TypeOf((*MockTradeRepository)(nil).Session))
}

// SetRangeAll mocks base method.
func (m *MockTradeRepository) SetRangeAll() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetRangeAll")
}

// SetRangeAll indicates an expected call of SetRangeAll.
func (mr *MockTradeRepositoryMockRecorder) SetRangeAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRangeAll", reflect.TypeOf((*MockTradeRepository)(nil).SetRangeAll))
}

// SetRangeByNumberAndEndTime mocks base method.
func (m *MockTradeRepository) SetRangeByNumberAndEndTime(arg0 int, arg1 time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetRangeByNumberAndEndTime", arg0, arg1)
}

// SetRangeByNumberAndEndTime indicates an expected call of SetRangeByNumberAndEndTime.
func (mr *MockTradeRepositoryMockRecorder) SetRangeByNumberAndEndTime(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRangeByNumberAndEndTime", reflect.TypeOf((*MockTradeRepository)(nil).SetRangeByNumberAndEndTime), arg0, arg1)
}

// SetRangeByTime mocks base method.
func (m *MockTradeRepository) SetRangeByTime(arg0, arg1 time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetRangeByTime", arg0, arg1)
}

// SetRangeByTime indicates an expected call of SetRangeByTime.
func (mr *MockTradeRepositoryMockRecorder) SetRangeByTime(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRangeByTime", reflect.TypeOf((*MockTradeRepository)(nil).SetRangeByTime), arg0, arg1)
}

// MockFetchingSession is a mock of FetchingSession interface.
type MockFetchingSession struct {
	ctrl     *gomock.Controller
	recorder *MockFetchingSessionMockRecorder
}

// MockFetchingSessionMockRecorder is the mock recorder for MockFetchingSession.
type MockFetchingSessionMockRecorder struct {
	mock *MockFetchingSession
}

// NewMockFetchingSession creates a new mock instance.
func NewMockFetchingSession(ctrl *gomock.Controller) *MockFetchingSession {
	mock := &MockFetchingSession{ctrl: ctrl}
	mock.recorder = &MockFetchingSessionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFetchingSession) EXPECT() *MockFetchingSessionMockRecorder {
	return m.recorder
}

// Next mocks base method.
func (m *MockFetchingSession) Next() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Next indicates an expected call of Next.
func (mr *MockFetchingSessionMockRecorder) Next() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockFetchingSession)(nil).Next))
}

// Value mocks base method.
func (m *MockFetchingSession) Value(arg0 context.Context) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Value", arg0)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Value indicates an expected call of Value.
func (mr *MockFetchingSessionMockRecorder) Value(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Value", reflect.TypeOf((*MockFetchingSession)(nil).Value), arg0)
}

// MockTradeStream is a mock of TradeStream interface.
type MockTradeStream struct {
	ctrl     *gomock.Controller
	recorder *MockTradeStreamMockRecorder
}

// MockTradeStreamMockRecorder is the mock recorder for MockTradeStream.
type MockTradeStreamMockRecorder struct {
	mock *MockTradeStream
}

// NewMockTradeStream creates a new mock instance.
func NewMockTradeStream(ctrl *gomock.Controller) *MockTradeStream {
	mock := &MockTradeStream{ctrl: ctrl}
	mock.recorder = &MockTradeStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTradeStream) EXPECT() *MockTradeStreamMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockTradeStream) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockTradeStreamMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockTradeStream)(nil).Close))
}

// SelectProduct mocks base method.
func (m *MockTradeStream) SelectProduct(arg0, arg1, arg2 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SelectProduct", arg0, arg1, arg2)
}

// SelectProduct indicates an expected call of SelectProduct.
func (mr *MockTradeStreamMockRecorder) SelectProduct(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectProduct", reflect.TypeOf((*MockTradeStream)(nil).SelectProduct), arg0, arg1, arg2)
}

// Session mocks base method.
func (m *MockTradeStream) Session() (TradeCursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Session")
	ret0, _ := ret[0].(TradeCursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Session indicates an expected call of Session.
func (mr *MockTradeStreamMockRecorder) Session() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Session", reflect.TypeOf((*MockTradeStream)(nil).Session))
}
