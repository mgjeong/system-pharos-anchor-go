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

package apps

import (
	"bytes"
	deploymentmocks "controller/deployment/group/mocks"
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

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/groups/groupID/invalid", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithExcludedBaseURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/groups/groupID/invalid", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithDeployRequest_ExpectCalledDeployApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().DeployApp("groupID", testBodyString),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/management/groups/groupID/deploy", bytes.NewReader(body))

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithUpdateAppInfoRequest_ExpectCalledUpdateAppInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().UpdateAppInfo("groupID", "appID", testBodyString),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/management/groups/groupID/apps/appID", bytes.NewReader(body))

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetAppsRequest_ExpectCalledGetApps(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().GetApps("groupID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/groups/groupID/apps", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetAppRequest_ExpectCalledGetApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().GetApp("groupID", "appID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/groups/groupID/apps/appID", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithDeleteAppRequest_ExpectCalledDeleteApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().DeleteApp("groupID", "appID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/management/groups/groupID/apps/appID", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithUpdateAppRequest_ExpectCalledUpdateApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().UpdateApp("groupID", "appID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/management/groups/groupID/apps/appID/update", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithStartAppRequest_ExpectCalledStartApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().StartApp("groupID", "appID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/management/groups/groupID/apps/appID/start", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithStopAppRequest_ExpectCalledStopApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().StopApp("groupID", "appID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/management/groups/groupID/apps/appID/stop", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}