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
	deploymentmocks "controller/deployment/node/mocks"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testBodyString = `{"test":"body"}`
	testQueryKey   = "testKey"
	testQueryValue = "testValue"
)

var testBody = map[string]interface{}{
	"test": "body",
}

var testQuery map[string]interface{}

var Handler Command

func init() {
	Handler = RequestHandler{}
	testQuery = make(map[string]interface{})
	testQueryValueList := make([]string, 1)
	testQueryValueList[0] = testQueryValue
	testQuery[testQueryKey] = testQueryValueList
}

func TestCalledHandleWithInvalidURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/nodes/nodeID/invalid", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithExcludedBaseURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nodes/nodeID/apps/invalid", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithDeployRequest_ExpectCalledDeployApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().DeployApp("nodeID", testBodyString),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/management/nodes/nodeID/apps/deploy", bytes.NewReader(body))

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithUpdateAppInfoRequest_ExpectCalledUpdateAppInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().UpdateAppInfo("nodeID", "appID", testBodyString),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/management/nodes/nodeID/apps/appID", bytes.NewReader(body))

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetAppsRequest_ExpectCalledGetApps(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().GetApps("nodeID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/nodes/nodeID/apps", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetAppRequest_ExpectCalledGetApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().GetApp("nodeID", "appID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/nodes/nodeID/apps/appID", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithDeleteAppRequest_ExpectCalledDeleteApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().DeleteApp("nodeID", "appID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/management/nodes/nodeID/apps/appID", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithUpdateAppRequestWithoutQuery_ExpectCalledUpdateApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().UpdateApp("nodeID", "appID", nil),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/management/nodes/nodeID/apps/appID/update", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithUpdateAppRequestWithQuery_ExpectCalledUpdateApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().UpdateApp("nodeID", "appID", testQuery),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/management/nodes/nodeID/apps/appID/update", nil)

	query := req.URL.Query()
	query.Add(testQueryKey, testQueryValue)
	req.URL.RawQuery = query.Encode()
	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithStartAppRequest_ExpectCalledStartApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().StartApp("nodeID", "appID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/management/nodes/nodeID/apps/appID/start", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithStopAppRequest_ExpectCalledStopApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deploymentMockObj := deploymentmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		deploymentMockObj.EXPECT().StopApp("nodeID", "appID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/management/nodes/nodeID/apps/appID/stop", nil)

	// pass mockObj to a real object.
	deploymentExecutor = deploymentMockObj

	Handler.Handle(w, req)
}
