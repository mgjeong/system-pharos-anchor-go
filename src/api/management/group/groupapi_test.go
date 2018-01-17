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

package group

import (
	"bytes"
	groupmanagermocks "controller/management/group/mocks"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testBodyString = `{"test":"body"}`
)

var testBody = map[string]interface{}{
	"test": "body",
}

var Handler Command

func init() {
	Handler = RequestHandler{}
}

func TestCalledHandleWithInvalidURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupmanageMockObj := groupmanagermocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/invalid", nil)

	// pass mockObj to a real object.
	managementExecutor = groupmanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithExcludedBaseURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupmanageMockObj := groupmanagermocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/management/groups", nil)

	// pass mockObj to a real object.
	managementExecutor = groupmanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGroupsRequest_ExpectCalledGetGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupmanageMockObj := groupmanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupmanageMockObj.EXPECT().GetGroups(),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/groups", nil)

	// pass mockObj to a real object.
	managementExecutor = groupmanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetGroupRequest_ExpectCalledGetGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupmanageMockObj := groupmanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupmanageMockObj.EXPECT().GetGroup("groupID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/groups/groupID", nil)

	// pass mockObj to a real object.
	managementExecutor = groupmanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithDeleteGroupRequest_ExpectCalledDeleteGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupmanageMockObj := groupmanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupmanageMockObj.EXPECT().DeleteGroup("groupID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/management/groups/groupID", nil)

	// pass mockObj to a real object.
	managementExecutor = groupmanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithCreateGroupRequest_ExpectCalledCreateGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupmanageMockObj := groupmanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupmanageMockObj.EXPECT().CreateGroup(testBodyString),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/management/groups/create", bytes.NewReader(body))

	// pass mockObj to a real object.
	managementExecutor = groupmanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithJoinGroupRequest_ExpectCalledJoinGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupmanageMockObj := groupmanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupmanageMockObj.EXPECT().JoinGroup("groupID", testBodyString),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/management/groups/groupID/join", bytes.NewReader(body))

	// pass mockObj to a real object.
	managementExecutor = groupmanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithLeaveGroupRequest_ExpectCalledLeaveGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupmanageMockObj := groupmanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupmanageMockObj.EXPECT().LeaveGroup("groupID", testBodyString),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/management/groups/groupID/leave", bytes.NewReader(body))

	// pass mockObj to a real object.
	managementExecutor = groupmanageMockObj

	Handler.Handle(w, req)
}
