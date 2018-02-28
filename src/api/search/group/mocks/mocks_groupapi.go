/*******************************************************************************
 * Copyright 2018 Samsung Electronics All Rights Reserved.
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
// Source: groupapi.go

// Package mock_group is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	http "net/http"
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

// Handle mocks base method
func (m *MockCommand) Handle(w http.ResponseWriter, req *http.Request) {
	m.ctrl.Call(m, "Handle", w, req)
}

// Handle indicates an expected call of Handle
func (mr *MockCommandMockRecorder) Handle(w, req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockCommand)(nil).Handle), w, req)
}

// MockgroupSearchAPI is a mock of groupSearchAPI interface
type MockgroupSearchAPI struct {
	ctrl     *gomock.Controller
	recorder *MockgroupSearchAPIMockRecorder
}

// MockgroupSearchAPIMockRecorder is the mock recorder for MockgroupSearchAPI
type MockgroupSearchAPIMockRecorder struct {
	mock *MockgroupSearchAPI
}

// NewMockgroupSearchAPI creates a new mock instance
func NewMockgroupSearchAPI(ctrl *gomock.Controller) *MockgroupSearchAPI {
	mock := &MockgroupSearchAPI{ctrl: ctrl}
	mock.recorder = &MockgroupSearchAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockgroupSearchAPI) EXPECT() *MockgroupSearchAPIMockRecorder {
	return m.recorder
}

// searchGroups mocks base method
func (m *MockgroupSearchAPI) searchGroups(w http.ResponseWriter, req *http.Request) {
	m.ctrl.Call(m, "searchGroups", w, req)
}

// searchGroups indicates an expected call of searchGroups
func (mr *MockgroupSearchAPIMockRecorder) searchGroups(w, req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "searchGroups", reflect.TypeOf((*MockgroupSearchAPI)(nil).searchGroups), w, req)
}
