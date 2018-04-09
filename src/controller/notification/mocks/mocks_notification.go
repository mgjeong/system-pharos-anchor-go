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
// Source: notificationapi.go

// Package mock_notification is a generated GoMock package.
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

// MocknotificationEventAPI is a mock of notificationEventAPI interface
type MocknotificationEventAPI struct {
	ctrl     *gomock.Controller
	recorder *MocknotificationEventAPIMockRecorder
}

// MocknotificationEventAPIMockRecorder is the mock recorder for MocknotificationEventAPI
type MocknotificationEventAPIMockRecorder struct {
	mock *MocknotificationEventAPI
}

// NewMocknotificationEventAPI creates a new mock instance
func NewMocknotificationEventAPI(ctrl *gomock.Controller) *MocknotificationEventAPI {
	mock := &MocknotificationEventAPI{ctrl: ctrl}
	mock.recorder = &MocknotificationEventAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MocknotificationEventAPI) EXPECT() *MocknotificationEventAPIMockRecorder {
	return m.recorder
}

// registerNotificationEvent mocks base method
func (m *MocknotificationEventAPI) registerNotificationEvent(w http.ResponseWriter, req *http.Request) {
	m.ctrl.Call(m, "registerNotificationEvent", w, req)
}

// registerNotificationEvent indicates an expected call of registerNotificationEvent
func (mr *MocknotificationEventAPIMockRecorder) registerNotificationEvent(w, req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "registerNotificationEvent", reflect.TypeOf((*MocknotificationEventAPI)(nil).registerNotificationEvent), w, req)
}

// unRegisterNotificationEvent mocks base method
func (m *MocknotificationEventAPI) unRegisterNotificationEvent(w http.ResponseWriter, req *http.Request, eventId string) {
	m.ctrl.Call(m, "unRegisterNotificationEvent", w, req, eventId)
}

// unRegisterNotificationEvent indicates an expected call of unRegisterNotificationEvent
func (mr *MocknotificationEventAPIMockRecorder) unRegisterNotificationEvent(w, req, eventId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "unRegisterNotificationEvent", reflect.TypeOf((*MocknotificationEventAPI)(nil).unRegisterNotificationEvent), w, req, eventId)
}

// receiveNotificationEvnet mocks base method
func (m *MocknotificationEventAPI) receiveNotificationEvnet(w http.ResponseWriter, req *http.Request) {
	m.ctrl.Call(m, "receiveNotificationEvnet", w, req)
}

// receiveNotificationEvnet indicates an expected call of receiveNotificationEvnet
func (mr *MocknotificationEventAPIMockRecorder) receiveNotificationEvnet(w, req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "receiveNotificationEvnet", reflect.TypeOf((*MocknotificationEventAPI)(nil).receiveNotificationEvnet), w, req)
}
