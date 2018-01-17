// Automatically generated by MockGen. DO NOT EDIT!
// Source: group/group.go

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

// CreateGroup mocks base method
func (m *MockCommand) CreateGroup(name string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "CreateGroup", name)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroup indicates an expected call of CreateGroup
func (mr *MockCommandMockRecorder) CreateGroup(name interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockCommand)(nil).CreateGroup), name)
}

// GetGroup mocks base method
func (m *MockCommand) GetGroup(groupId string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetGroup", groupId)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroup indicates an expected call of GetGroup
func (mr *MockCommandMockRecorder) GetGroup(groupId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockCommand)(nil).GetGroup), groupId)
}

// GetGroups mocks base method
func (m *MockCommand) GetGroups() ([]map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetGroups")
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroups indicates an expected call of GetGroups
func (mr *MockCommandMockRecorder) GetGroups() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroups", reflect.TypeOf((*MockCommand)(nil).GetGroups))
}

// GetGroupMembers mocks base method
func (m *MockCommand) GetGroupMembers(groupId string) ([]map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetGroupMembers", groupId)
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupMembers indicates an expected call of GetGroupMembers
func (mr *MockCommandMockRecorder) GetGroupMembers(groupId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupMembers", reflect.TypeOf((*MockCommand)(nil).GetGroupMembers), groupId)
}

// GetGroupMembersByAppID mocks base method
func (m *MockCommand) GetGroupMembersByAppID(groupId, appId string) ([]map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetGroupMembersByAppID", groupId, appId)
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupMembersByAppID indicates an expected call of GetGroupMembersByAppID
func (mr *MockCommandMockRecorder) GetGroupMembersByAppID(groupId, appId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupMembersByAppID", reflect.TypeOf((*MockCommand)(nil).GetGroupMembersByAppID), groupId, appId)
}

// JoinGroup mocks base method
func (m *MockCommand) JoinGroup(groupId, nodeId string) error {
	ret := m.ctrl.Call(m, "JoinGroup", groupId, nodeId)
	ret0, _ := ret[0].(error)
	return ret0
}

// JoinGroup indicates an expected call of JoinGroup
func (mr *MockCommandMockRecorder) JoinGroup(groupId, nodeId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JoinGroup", reflect.TypeOf((*MockCommand)(nil).JoinGroup), groupId, nodeId)
}

// LeaveGroup mocks base method
func (m *MockCommand) LeaveGroup(groupId, nodeId string) error {
	ret := m.ctrl.Call(m, "LeaveGroup", groupId, nodeId)
	ret0, _ := ret[0].(error)
	return ret0
}

// LeaveGroup indicates an expected call of LeaveGroup
func (mr *MockCommandMockRecorder) LeaveGroup(groupId, nodeId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaveGroup", reflect.TypeOf((*MockCommand)(nil).LeaveGroup), groupId, nodeId)
}

// DeleteGroup mocks base method
func (m *MockCommand) DeleteGroup(groupId string) error {
	ret := m.ctrl.Call(m, "DeleteGroup", groupId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGroup indicates an expected call of DeleteGroup
func (mr *MockCommandMockRecorder) DeleteGroup(groupId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroup", reflect.TypeOf((*MockCommand)(nil).DeleteGroup), groupId)
}
