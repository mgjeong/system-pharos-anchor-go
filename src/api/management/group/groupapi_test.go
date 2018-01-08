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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Test functions for Group API Handler.

type handleFunc struct {
	functionCall string
}

func TestHandle(t *testing.T) {
	w := httptest.NewRecorder()
	mockHandle := handleFunc{}
	defaultApis := SdamGroup
	SdamGroup = &mockHandle
	Input := [][]string{
		{GET, "/api/v1/groups", "groups"},
		{POST, "/api/v1/groups/create", "createGroup"},
		{GET, "/api/v1/groups/groupID", "group"},
		{DELETE, "/api/v1/groups/groupID", "group"},
		{POST, "/api/v1/groups/groupID/deploy", "groupDeployApp"},
		{POST, "/api/v1/groups/groupID/join", "groupJoin"},
		{POST, "/api/v1/groups/groupID/leave", "groupLeave"},
		{GET, "/api/v1/groups/groupID/apps", "groupInfoApps"},
		{GET, "/api/v1/groups/groupID/apps/appID", "groupInfoApp"},
		{POST, "/api/v1/groups/groupID/apps/appID", "groupUpdateAppInfo"},
		{DELETE, "/api/v1/groups/groupID/apps/appID", "groupDeleteApp"},
		{POST, "/api/v1/groups/groupID/apps/appID/start", "groupStartApp"},
		{POST, "/api/v1/groups/groupID/apps/appID/stop", "groupStopApp"},
		{POST, "/api/v1/groups/groupID/apps/appID/update", "groupUpdateApp"},
	}
	for _, val := range Input {
		method, url, funcname := val[0], val[1], val[2]
		req, _ := http.NewRequest(method, url, nil)
		SdamGroupHandle.Handle(w, req)
		if mockHandle.functionCall != funcname {
			t.Error("[SDAM][Group]Handle is invalid about " + funcname)
		}
	}
	SdamGroup = defaultApis
}

func TestHandle_Invalid_Method(t *testing.T) {
	w := httptest.NewRecorder()
	Input := map[string][]string{
		"/api/v1/groups":                           {POST, DELETE, PUT},
		"/api/v1/groups/create":                    {GET, DELETE, PUT},
		"/api/v1/groups/groupID":                   {POST, PUT},
		"/api/v1/groups/groupID/deploy":            {GET, DELETE, PUT},
		"/api/v1/groups/groupID/join":              {GET, DELETE, PUT},
		"/api/v1/groups/groupID/leave":             {GET, DELETE, PUT},
		"/api/v1/groups/groupID/apps":              {POST, DELETE, PUT},
		"/api/v1/groups/groupID/apps/appID":        {PUT},
		"/api/v1/groups/groupID/apps/appID/start":  {GET, DELETE, PUT},
		"/api/v1/groups/groupID/apps/appID/stop":   {GET, DELETE, PUT},
		"/api/v1/groups/groupID/apps/appID/update": {GET, DELETE, PUT},
	}
	for key, vals := range Input {
		for _, val := range vals {
			req, _ := http.NewRequest(val, key, nil)
			SdamGroupHandle.Handle(w, req)
			if w.Code != http.StatusBadRequest {
				t.Error("[SDAM][Group]Handle is invalid")
			}
		}
	}
}

//Mock functions for Group APIs.

func (mockHandle *handleFunc) createGroup(w http.ResponseWriter, req *http.Request) {
	mockHandle.functionCall = "createGroup"
}

func (mockHandle *handleFunc) group(w http.ResponseWriter, req *http.Request, groupID string) {
	mockHandle.functionCall = "group"
}

func (mockHandle *handleFunc) groups(w http.ResponseWriter, req *http.Request) {
	mockHandle.functionCall = "groups"
}

func (mockHandle *handleFunc) groupJoin(w http.ResponseWriter, req *http.Request, groupID string) {
	mockHandle.functionCall = "groupJoin"
}

func (mockHandle *handleFunc) groupLeave(w http.ResponseWriter, req *http.Request, groupID string) {
	mockHandle.functionCall = "groupLeave"
}

func (mockHandle *handleFunc) groupDeployApp(w http.ResponseWriter, req *http.Request, groupID string) {
	mockHandle.functionCall = "groupDeployApp"
}

func (mockHandle *handleFunc) groupInfoApps(w http.ResponseWriter, req *http.Request, groupID string) {
	mockHandle.functionCall = "groupInfoApps"
}

func (mockHandle *handleFunc) groupInfoApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	mockHandle.functionCall = "groupInfoApp"
}

func (mockHandle *handleFunc) groupUpdateAppInfo(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	mockHandle.functionCall = "groupUpdateAppInfo"
}

func (mockHandle *handleFunc) groupDeleteApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	mockHandle.functionCall = "groupDeleteApp"
}

func (mockHandle *handleFunc) groupStartApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	mockHandle.functionCall = "groupStartApp"
}

func (mockHandle *handleFunc) groupStopApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	mockHandle.functionCall = "groupStopApp"
}

func (mockHandle *handleFunc) groupUpdateApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	mockHandle.functionCall = "groupUpdateApp"
}

//Test functions for Group APIs.

type controllerFunc struct {
	functionCall  string
	occurredError bool
}

func newCtrlFunc() *controllerFunc {
	cf := controllerFunc{}
	cf.functionCall = ""
	cf.occurredError = false
	return &cf
}

var testBody = map[string]interface{}{
	"id":   "testID",
	"host": "0.0.0.0",
	"port": "8080",
	"apps": nil,
}

func TestCreateGroup(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/groups/create", nil)
	sdamGroupController = mockCtrl
	SdamGroup.createGroup(w, req)
	if mockCtrl.functionCall != "CreateGroup" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]createGroup is invalid")
	}
}

func TestCreateGroup_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/groups/create", nil)
	sdamGroupController = mockCtrl
	SdamGroup.createGroup(w, req)
	if mockCtrl.functionCall != "CreateGroup" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]createGroup is invalid about controller occurred error")
	}
}

func TestGroupGET(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/groups/testGroupID", nil)
	sdamGroupController = mockCtrl
	SdamGroup.group(w, req, "testGroupID")
	if mockCtrl.functionCall != "GetGroup" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]group is invalid")
	}
}

func TestGroupGET_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/groups/testGroupID", nil)
	sdamGroupController = mockCtrl
	SdamGroup.group(w, req, "testGroupID")
	if mockCtrl.functionCall != "GetGroup" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]group is invalid about controller occurred error")
	}
}

func TestGroupDELETE(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(DELETE, "/api/v1/groups/testGroupID", nil)
	sdamGroupController = mockCtrl
	SdamGroup.group(w, req, "testGroupID")
	if mockCtrl.functionCall != "DeleteGroup" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]group is invalid")
	}
}

func TestGroupDELETE_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/groups/testGroupID", nil)
	sdamGroupController = mockCtrl
	SdamGroup.group(w, req, "testGroupID")
	if mockCtrl.functionCall != "GetGroup" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]group is invalid about controller occurred error")
	}
}

func TestGroups(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/groups", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groups(w, req)
	if mockCtrl.functionCall != "GetGroups" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]groups is invalid")
	}
}

func TestGroups_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/groups", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groups(w, req)
	if mockCtrl.functionCall != "GetGroups" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]groups is invalid about controller occurred error")
	}
}

func TestGroupJoin(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	bod, _ := json.Marshal(testBody)
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/join", bytes.NewReader(bod))
	sdamGroupController = mockCtrl
	SdamGroup.groupJoin(w, req, "testGroupID")
	if mockCtrl.functionCall != "JoinApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]groupJoin is invalid")
	}
}

func TestGroupJoin_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	bod := []byte("error")
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/join", bytes.NewReader(bod))
	sdamGroupController = mockCtrl
	SdamGroup.groupJoin(w, req, "testGroupID")
	if mockCtrl.functionCall != "JoinApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]groupJoin is invalid about controller occurred error")
	}
}

func TestGroupJoin_empty_body(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/join", nil)
	SdamGroup.groupJoin(w, req, "testGroupID")
	if w.Code != http.StatusBadRequest {
		t.Error("[SDAM][Group]groupJoin is invalid about emtpy body")
	}
}

func TestGroupLeave(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	bod, _ := json.Marshal(testBody)
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/leave", bytes.NewReader(bod))
	sdamGroupController = mockCtrl
	SdamGroup.groupLeave(w, req, "testGroupID")
	if mockCtrl.functionCall != "LeaveApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]groupLeave is invalid")
	}
}

func TestGroupLeave_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	bod := []byte("error")
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/leave", bytes.NewReader(bod))
	sdamGroupController = mockCtrl
	SdamGroup.groupLeave(w, req, "testGroupID")
	if mockCtrl.functionCall != "LeaveApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]groupLeave is invalid about controller occurred error")
	}
}

func TestGroupLeave_empty_body(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/leave", nil)
	SdamGroup.groupLeave(w, req, "testGroupID")
	if w.Code != http.StatusBadRequest {
		t.Error("[SDAM][Group]groupLeave is invalid about emtpy body")
	}
}

func TestGroupDeployApp(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	bod, _ := json.Marshal(testBody)
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/deploy", bytes.NewReader(bod))
	sdamGroupController = mockCtrl
	SdamGroup.groupDeployApp(w, req, "testGroupID")
	if mockCtrl.functionCall != "DeployApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]groupDeployApp is invalid")
	}
}

func TestGroupDeployApp_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	bod := []byte("error")
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/deploy", bytes.NewReader(bod))
	sdamGroupController = mockCtrl
	SdamGroup.groupDeployApp(w, req, "testGroupID")
	if mockCtrl.functionCall != "DeployApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]groupDeployApp is invalid about controller occurred error")
	}
}

func TestGroupDeployApp_empty_body(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/deploy", nil)
	SdamGroup.groupDeployApp(w, req, "testGroupID")
	if w.Code != http.StatusBadRequest {
		t.Error("[SDAM][Group]groupDeployApp is invalid about emtpy body")
	}
}

func TestGroupInfoApps(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/groups/testGroupID/apps", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groupInfoApps(w, req, "testGroupID")
	if mockCtrl.functionCall != "GetApps" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]groupInfoApps is invalid")
	}
}

func TestGroupInfoApps_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/groups/testGroupID/apps", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groupInfoApps(w, req, "testGroupID")
	if mockCtrl.functionCall != "GetApps" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]groupInfoApps is invalid about controller occurred error")
	}
}

func TestGroupInfoApp(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/groups/testGroupID/apps/testAppID", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groupInfoApp(w, req, "testGroupID", "testAppID")
	if mockCtrl.functionCall != "GetApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]groupInfoApps is invalid")
	}
}

func TestGroupInfoApp_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/groups/testGroupID/apps/testAppID", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groupInfoApp(w, req, "testGroupID", "testAppID")
	if mockCtrl.functionCall != "GetApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]groupInfoApps is invalid about controller occurred error")
	}
}

func TestGroupUpdateAppInfo(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	bod, _ := json.Marshal(testBody)
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/apps/testAppID", bytes.NewReader(bod))
	sdamGroupController = mockCtrl
	SdamGroup.groupUpdateAppInfo(w, req, "testGroupID", "testAppID")
	if mockCtrl.functionCall != "UpdateAppInfo" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]groupUpdateAppInfo is invalid")
	}
}

func TestGroupUpdateAppInfo_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	bod := []byte("error")
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/apps/testAppID", bytes.NewReader(bod))
	sdamGroupController = mockCtrl
	SdamGroup.groupUpdateAppInfo(w, req, "testGroupID", "testAppID")
	if mockCtrl.functionCall != "UpdateAppInfo" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]groupUpdateAppInfo is invalid about controller occurred error")
	}
}

func TestGroupUpdateAppInfo_empty_body(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/apps/testAppID", nil)
	SdamGroup.groupUpdateAppInfo(w, req, "testGroupID", "testAppID")
	if w.Code != http.StatusBadRequest {
		t.Error("[SDAM][Group]groupUpdateAppInfo is invalid about empty body")
	}
}

func TestDeleteApp(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(DELETE, "/api/v1/groups/testGroupID/apps/testAppID", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groupDeleteApp(w, req, "testGroupID", "testAppID")
	if mockCtrl.functionCall != "DeleteApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]groupDeleteApps is invalid")
	}
}

func TestDeleteApp_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(DELETE, "/api/v1/groups/testGroupID/apps/testAppID", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groupDeleteApp(w, req, "testGroupID", "testAppID")
	if mockCtrl.functionCall != "DeleteApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]groupDeleteApps is invalid about controller occurred error")
	}
}

func TestGroupStartApp(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/apps/testAppID/start", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groupStartApp(w, req, "testGroupID", "testAppID")
	if mockCtrl.functionCall != "StartApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]groupStartApps is invalid")
	}
}

func TestGroupStartApp_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/apps/testAppID/start", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groupStartApp(w, req, "testGroupID", "testAppID")
	if mockCtrl.functionCall != "StartApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]groupStartApps is invalid about controller occurred error")
	}
}

func TestGroupStopApp(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/apps/testAppID/stop", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groupStopApp(w, req, "testGroupID", "testAppID")
	if mockCtrl.functionCall != "StopApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]groupStopApps is invalid")
	}
}

func TestGroupStopApp_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/apps/testAppID/stop", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groupStopApp(w, req, "testGroupID", "testAppID")
	if mockCtrl.functionCall != "StopApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]groupStopApps is invalid about controller occurred error")
	}
}

func TestGroupUpdateApp(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/apps/testAppID/update", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groupUpdateApp(w, req, "testGroupID", "testAppID")
	if mockCtrl.functionCall != "UpdateApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Group]groupUpdateApps is invalid")
	}
}

func TestGroupUpdateApp_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/groups/testGroupID/apps/testAppID/update", nil)
	sdamGroupController = mockCtrl
	SdamGroup.groupUpdateApp(w, req, "testGroupID", "testAppID")
	if mockCtrl.functionCall != "UpdateApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Group]groupUpdateApps is invalid about controller occurred error")
	}
}

//Mock functions for Group Controller Functions.

func (mockCtrl *controllerFunc) CreateGroup() (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "CreateGroup"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) GetGroup(groupID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "GetGroup"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) DeleteGroup(groupID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "DeleteGroup"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) GetGroups() (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "GetGroups"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) JoinGroup(groupID string, body string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "JoinApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) LeaveGroup(groupID string, body string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "LeaveApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) DeployApp(groupID string, body string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "DeployApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) GetApps(groupID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "GetApps"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) GetApp(groupID string, appID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "GetApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) UpdateAppInfo(groupID string, appID string, body string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "UpdateAppInfo"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) DeleteApp(groupID string, appID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "DeleteApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) UpdateApp(groupID string, appID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "UpdateApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) StartApp(groupID string, appID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "StartApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) StopApp(groupID string, appID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "StopApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}
