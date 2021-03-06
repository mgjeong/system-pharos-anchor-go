// Automatically generated by MockGen. DO NOT EDIT!
// Source: registry/registry.go

package mocks

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

// AddDockerRegistry mocks base method
func (m *MockCommand) AddDockerRegistry(url string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "AddDockerRegistry", url)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDockerRegistry indicates an expected call of AddDockerRegistry
func (mr *MockCommandMockRecorder) AddDockerRegistry(url interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDockerRegistry", reflect.TypeOf((*MockCommand)(nil).AddDockerRegistry), url)
}

// GetDockerRegistries mocks base method
func (m *MockCommand) GetDockerRegistries() ([]map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetDockerRegistries")
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDockerRegistries indicates an expected call of GetDockerRegistries
func (mr *MockCommandMockRecorder) GetDockerRegistries() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDockerRegistries", reflect.TypeOf((*MockCommand)(nil).GetDockerRegistries))
}

// DeleteDockerRegistry mocks base method
func (m *MockCommand) DeleteDockerRegistry(registryId string) error {
	ret := m.ctrl.Call(m, "DeleteDockerRegistry", registryId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDockerRegistry indicates an expected call of DeleteDockerRegistry
func (mr *MockCommandMockRecorder) DeleteDockerRegistry(registryId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDockerRegistry", reflect.TypeOf((*MockCommand)(nil).DeleteDockerRegistry), registryId)
}
