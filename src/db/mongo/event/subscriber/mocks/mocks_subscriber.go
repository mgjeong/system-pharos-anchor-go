// Code generated by MockGen. DO NOT EDIT.
// Source: src/db/mongo/event/subscriber/subscriber.go

// Package mock_subscriber is a generated GoMock package.
package mock_subscriber

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCommand is a mock of Command interface
type MockCommand struct {
	ctrl     *gomock.Controller
	recorder *MockCommandMockRecorder
}

// MockCommandMockRecorder is the mock recorder for MockCommand
type MockCommandMockRecorder struct {
	mock *MockCommand
}

// NewMockCommand creates a new mock instance
func NewMockCommand(ctrl *gomock.Controller) *MockCommand {
	mock := &MockCommand{ctrl: ctrl}
	mock.recorder = &MockCommandMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommand) EXPECT() *MockCommandMockRecorder {
	return m.recorder
}

// AddSubscriber mocks base method
func (m *MockCommand) AddSubscriber(id, eventType, URL string, Status, eventId []string) error {
	ret := m.ctrl.Call(m, "AddSubscriber", id, eventType, URL, Status, eventId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSubscriber indicates an expected call of AddSubscriber
func (mr *MockCommandMockRecorder) AddSubscriber(id, eventType, URL, Status, eventId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSubscriber", reflect.TypeOf((*MockCommand)(nil).AddSubscriber), id, eventType, URL, Status, eventId)
}

// GetSubscriber mocks base method
func (m *MockCommand) GetSubscriber(id string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetSubscriber", id)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriber indicates an expected call of GetSubscriber
func (mr *MockCommandMockRecorder) GetSubscriber(id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriber", reflect.TypeOf((*MockCommand)(nil).GetSubscriber), id)
}

// DeleteSubscriber mocks base method
func (m *MockCommand) DeleteSubscriber(id string) error {
	ret := m.ctrl.Call(m, "DeleteSubscriber", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSubscriber indicates an expected call of DeleteSubscriber
func (mr *MockCommandMockRecorder) DeleteSubscriber(id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSubscriber", reflect.TypeOf((*MockCommand)(nil).DeleteSubscriber), id)
}
