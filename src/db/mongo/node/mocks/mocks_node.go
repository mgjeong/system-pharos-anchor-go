// Automatically generated by MockGen. DO NOT EDIT!
// Source: node/node.go

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

// AddNode mocks base method
func (m *MockCommand) AddNode(ip, status string, config map[string]interface{}, apps []string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "AddNode", ip, status, config, apps)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNode indicates an expected call of AddNode
func (mr *MockCommandMockRecorder) AddNode(ip, status, config, apps interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNode", reflect.TypeOf((*MockCommand)(nil).AddNode), ip, status, config, apps)
}

// UpdateNodeAddress mocks base method
func (m *MockCommand) UpdateNodeAddress(nodeId, host, port string) error {
	ret := m.ctrl.Call(m, "UpdateNodeAddress", nodeId, host, port)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNodeAddress indicates an expected call of UpdateNodeAddress
func (mr *MockCommandMockRecorder) UpdateNodeAddress(nodeId, host, port interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNodeAddress", reflect.TypeOf((*MockCommand)(nil).UpdateNodeAddress), nodeId, host, port)
}

// UpdateNodeStatus mocks base method
func (m *MockCommand) UpdateNodeStatus(nodeId, status string) error {
	ret := m.ctrl.Call(m, "UpdateNodeStatus", nodeId, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNodeStatus indicates an expected call of UpdateNodeStatus
func (mr *MockCommandMockRecorder) UpdateNodeStatus(nodeId, status interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNodeStatus", reflect.TypeOf((*MockCommand)(nil).UpdateNodeStatus), nodeId, status)
}

// UpdateNodeConfiguration mocks base method
func (m *MockCommand) UpdateNodeConfiguration(nodeId string, config map[string]interface{}) error {
	ret := m.ctrl.Call(m, "UpdateNodeConfiguration", nodeId, config)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNodeConfiguration indicates an expected call of UpdateNodeConfiguration
func (mr *MockCommandMockRecorder) UpdateNodeConfiguration(nodeId, config interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNodeConfiguration", reflect.TypeOf((*MockCommand)(nil).UpdateNodeConfiguration), nodeId, config)
}

// GetNode mocks base method
func (m *MockCommand) GetNode(nodeId string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetNode", nodeId)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNode indicates an expected call of GetNode
func (mr *MockCommandMockRecorder) GetNode(nodeId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNode", reflect.TypeOf((*MockCommand)(nil).GetNode), nodeId)
}

// GetNodes mocks base method
func (m *MockCommand) GetNodes(queryOptional ...map[string]interface{}) ([]map[string]interface{}, error) {
	varargs := []interface{}{}
	for _, a := range queryOptional {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetNodes", varargs...)
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodes indicates an expected call of GetNodes
func (mr *MockCommandMockRecorder) GetNodes(queryOptional ...interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodes", reflect.TypeOf((*MockCommand)(nil).GetNodes), queryOptional...)
}

// GetNodeByAppID mocks base method
func (m *MockCommand) GetNodeByAppID(nodeId, appId string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetNodeByAppID", nodeId, appId)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeByAppID indicates an expected call of GetNodeByAppID
func (mr *MockCommandMockRecorder) GetNodeByAppID(nodeId, appId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeByAppID", reflect.TypeOf((*MockCommand)(nil).GetNodeByAppID), nodeId, appId)
}

// GetNodeByIP mocks base method
func (m *MockCommand) GetNodeByIP(ip string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetNodeByIP", ip)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeByIP indicates an expected call of GetNodeByIP
func (mr *MockCommandMockRecorder) GetNodeByIP(ip interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeByIP", reflect.TypeOf((*MockCommand)(nil).GetNodeByIP), ip)
}

// AddAppToNode mocks base method
func (m *MockCommand) AddAppToNode(nodeId, appId string) error {
	ret := m.ctrl.Call(m, "AddAppToNode", nodeId, appId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAppToNode indicates an expected call of AddAppToNode
func (mr *MockCommandMockRecorder) AddAppToNode(nodeId, appId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAppToNode", reflect.TypeOf((*MockCommand)(nil).AddAppToNode), nodeId, appId)
}

// DeleteAppFromNode mocks base method
func (m *MockCommand) DeleteAppFromNode(nodeId, appId string) error {
	ret := m.ctrl.Call(m, "DeleteAppFromNode", nodeId, appId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAppFromNode indicates an expected call of DeleteAppFromNode
func (mr *MockCommandMockRecorder) DeleteAppFromNode(nodeId, appId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAppFromNode", reflect.TypeOf((*MockCommand)(nil).DeleteAppFromNode), nodeId, appId)
}

// DeleteNode mocks base method
func (m *MockCommand) DeleteNode(nodeId string) error {
	ret := m.ctrl.Call(m, "DeleteNode", nodeId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNode indicates an expected call of DeleteNode
func (mr *MockCommandMockRecorder) DeleteNode(nodeId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNode", reflect.TypeOf((*MockCommand)(nil).DeleteNode), nodeId)
}
