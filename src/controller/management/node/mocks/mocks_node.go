/*******************************************************************************
 * Copyright 2017 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/

// Code generated by MockGen. DO NOT EDIT.
// Source: node.go

// Package mock_node is a generated GoMock package.
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

// RegisterNode mocks base method
func (m *MockCommand) RegisterNode(body string) (int, map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "RegisterNode", body)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(map[string]interface{})
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RegisterNode indicates an expected call of RegisterNode
func (mr *MockCommandMockRecorder) RegisterNode(body interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterNode", reflect.TypeOf((*MockCommand)(nil).RegisterNode), body)
}

// UnRegisterNode mocks base method
func (m *MockCommand) UnRegisterNode(nodeId string) (int, error) {
	ret := m.ctrl.Call(m, "UnRegisterNode", nodeId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnRegisterNode indicates an expected call of UnRegisterNode
func (mr *MockCommandMockRecorder) UnRegisterNode(nodeId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnRegisterNode", reflect.TypeOf((*MockCommand)(nil).UnRegisterNode), nodeId)
}

// GetNode mocks base method
func (m *MockCommand) GetNode(nodeId string) (int, map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetNode", nodeId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(map[string]interface{})
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNode indicates an expected call of GetNode
func (mr *MockCommandMockRecorder) GetNode(nodeId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNode", reflect.TypeOf((*MockCommand)(nil).GetNode), nodeId)
}

// GetNodes mocks base method
func (m *MockCommand) GetNodes() (int, map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetNodes")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(map[string]interface{})
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNodes indicates an expected call of GetNodes
func (mr *MockCommandMockRecorder) GetNodes() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodes", reflect.TypeOf((*MockCommand)(nil).GetNodes))
}

// GetNodesWithAppID mocks base method
func (m *MockCommand) GetNodesWithAppID(appId string) (int, map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetNodesWithAppID", appId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(map[string]interface{})
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNodesWithAppID indicates an expected call of GetNodesWithAppID
func (mr *MockCommandMockRecorder) GetNodesWithAppID(appId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodesWithAppID", reflect.TypeOf((*MockCommand)(nil).GetNodesWithAppID), appId)
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

// PingNode mocks base method
func (m *MockCommand) PingNode(nodeId, body string) (int, error) {
	ret := m.ctrl.Call(m, "PingNode", nodeId, body)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PingNode indicates an expected call of PingNode
func (mr *MockCommandMockRecorder) PingNode(nodeId, body interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingNode", reflect.TypeOf((*MockCommand)(nil).PingNode), nodeId, body)
}
