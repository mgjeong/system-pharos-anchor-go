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

package node

import (
	"bytes"
	nodemanagermocks "controller/management/node/mocks"
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

	nodemanageMockObj := nodemanagermocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/invalid", nil)

	// pass mockObj to a real object.
	managementExecutor = nodemanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithExcludedBaseURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodemanageMockObj := nodemanagermocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/management/nodes", nil)

	// pass mockObj to a real object.
	managementExecutor = nodemanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetNodesRequest_ExpectCalledGetNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodemanageMockObj := nodemanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodemanageMockObj.EXPECT().GetNodes(),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/nodes", nil)

	// pass mockObj to a real object.
	managementExecutor = nodemanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetNodeRequest_ExpectCalledGetNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodemanageMockObj := nodemanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodemanageMockObj.EXPECT().GetNode("nodeID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/nodes/nodeID", nil)

	// pass mockObj to a real object.
	managementExecutor = nodemanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithRegisterNodeRequest_ExpectCalledRegisterNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodemanageMockObj := nodemanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodemanageMockObj.EXPECT().RegisterNode(testBodyString),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/management/nodes/register", bytes.NewReader(body))

	// pass mockObj to a real object.
	managementExecutor = nodemanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithUnregisterNodeRequest_ExpectCalledUnregisterNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodemanageMockObj := nodemanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodemanageMockObj.EXPECT().UnRegisterNode("nodeID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/management/nodes/nodeID/unregister", nil)

	// pass mockObj to a real object.
	managementExecutor = nodemanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithPingNodeRequest_ExpectCalledPingNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodemanageMockObj := nodemanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodemanageMockObj.EXPECT().PingNode("nodeID", testBodyString),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/management/nodes/nodeID/ping", bytes.NewReader(body))

	// pass mockObj to a real object.
	managementExecutor = nodemanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithConfigurationNodeRequest_ExpectCalledGetNodeConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodemanageMockObj := nodemanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodemanageMockObj.EXPECT().GetNodeConfiguration("nodeID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/nodes/nodeID/configuration", nil)

	// pass mockObj to a real object.
	managementExecutor = nodemanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithConfigurationNodeRequest_ExpectCalledSetNodeConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodemanageMockObj := nodemanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodemanageMockObj.EXPECT().SetNodeConfiguration("nodeID", testBodyString),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/management/nodes/nodeID/configuration", bytes.NewReader(body))

	// pass mockObj to a real object.
	managementExecutor = nodemanageMockObj

	Handler.Handle(w, req)
}

func TestRebootRequest_ExpectRebootCalled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodemanageMockObj := nodemanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodemanageMockObj.EXPECT().Reboot("nodeID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/management/nodes/nodeID/reboot", nil)

	// pass mockObj to a real object.
	managementExecutor = nodemanageMockObj

	Handler.Handle(w, req)
}

func TestRestoreRequest_ExpectRestoreCalled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodemanageMockObj := nodemanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodemanageMockObj.EXPECT().Restore("nodeID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/management/nodes/nodeID/restore", nil)

	// pass mockObj to a real object.
	managementExecutor = nodemanageMockObj

	Handler.Handle(w, req)
}
