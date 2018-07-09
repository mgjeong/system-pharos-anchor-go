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
	"commons/errors"
	"commons/results"
	"commons/util"
	searchmocks "controller/search/group/mocks"
	nodedbmocks "db/mongo/node/mocks"
	"encoding/json"
	"github.com/golang/mock/gomock"
	msgmocks "messenger/mocks"
	"reflect"
	"testing"
)

const (
	status       = "connected"
	appId        = "000000000000000000000000"
	invalidAppId = "000000000000000000000001"
	nodeId       = "54919CA5-4101-4AE4-595B-353C51AA983C"
	ip           = "127.0.0.1"
	port         = "48098"
)

var (
	registrationBody = map[string]interface{}{
		"ip":     ip,
		"config": config,
		"apps":   []string{},
	}
	property = map[string]interface{}{
		"key": "value",
	}
	reverseproxy = map[string]interface{}{
		"reverseproxy": map[string]interface{}{
			"enabled": "false",
		},
	}
	properties = []interface{}{property, reverseproxy}
	config     = map[string]interface{}{
		"properties": properties,
	}
	node = map[string]interface{}{
		"id":     nodeId,
		"ip":     ip,
		"apps":   []string{},
		"config": config,
	}
	nodeWithoutConfig = map[string]interface{}{
		"id":   nodeId,
		"ip":   ip,
		"apps": []string{},
	}
	groups = map[string]interface{}{
		"groups": []map[string]interface{}{},
	}
	body          = `{"description":"description"}`
	respCode      = []int{results.OK}
	respStr       = []string{`{"response":"response"}`}
	notFoundError = errors.NotFound{}
)

var manager Command

func init() {
	manager = Executor{}
}

func TestCalledRegisterNodeWithValidBody_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(gomock.Any()).Return(nil, notFoundError),
		nodedDBExecutorMockObj.EXPECT().AddNode(gomock.Any(), ip, status, gomock.Any(), []string{}).Return(node, nil),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	jsonString, _ := json.Marshal(registrationBody)
	code, _, err := manager.RegisterNode(string(jsonString))

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledRegisterNodeWithInValidJsonFormatBody_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidBody := `{"ip"}`

	code, _, err := manager.RegisterNode(invalidBody)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", err.Error())
	case errors.InvalidJSON:
	}
}

func TestCalledRegisterNodeWithInvalidBodyNotIncludingIPField_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidBody := `{"key":"value"}`

	code, _, err := manager.RegisterNode(invalidBody)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", err.Error())
	case errors.InvalidJSON:
	}
}

func TestCalledRegisterNodeWhenFailedToInsertNewNodeToDB_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(gomock.Any()).Return(nil, notFoundError),
		nodedDBExecutorMockObj.EXPECT().AddNode(gomock.Any(), ip, status, gomock.Any(), []string{}).Return(nil, notFoundError),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	jsonString, _ := json.Marshal(registrationBody)
	code, _, err := manager.RegisterNode(string(jsonString))

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledUnRegisterNodeWithValidBody_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := make(map[string]interface{})
	query["nodeId"] = []string{nodeId}

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/unregister"}

	msgMockObj := msgmocks.NewMockCommand(ctrl)
	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	searchExecutorMockObj := searchmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, respStr),
		nodedDBExecutorMockObj.EXPECT().DeleteNode(nodeId).Return(nil),
		searchExecutorMockObj.EXPECT().SearchGroups(query).Return(results.OK, groups, nil),
	)
	// pass mockObj to a real object.
	httpExecutor = msgMockObj
	nodeDbExecutor = nodedDBExecutorMockObj
	groupSearchExecutor = searchExecutorMockObj

	code, err := manager.UnRegisterNode(nodeId)

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledUnRegisterNodeWhenDBHasNotMatchedNode_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	code, err := manager.UnRegisterNode(nodeId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetNode_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	code, res, err := manager.GetNode(nodeId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(res, nodeWithoutConfig) {
		t.Error()
	}
}

func TestCalledGetNodeWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(nil, notFoundError),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	code, _, err := manager.GetNode(nodeId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetNodes_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodes := []map[string]interface{}{node}
	nodesWithoutConfig := []map[string]interface{}{nodeWithoutConfig}

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNodes().Return(nodes, nil),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	code, res, err := manager.GetNodes()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(res["nodes"].([]map[string]interface{}), nodesWithoutConfig) {
		t.Error()
	}
}

func TestCalledGetNodesWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNodes().Return(nil, notFoundError),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	code, _, err := manager.GetNodes()

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetNodesWithAppId_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodes := []map[string]interface{}{node}

	//make the query
	query := make(map[string]interface{})
	query[APPS] = appId

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNodes(query).Return(nodes, nil),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	code, res, err := manager.GetNodesWithAppID(appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(res["nodes"].([]map[string]interface{}), nodes) {
		t.Error()
	}
}

func TestCalledGetNodesWithInvalidAppId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	//make the query
	query := make(map[string]interface{})
	query[APPS] = invalidAppId

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNodes(query).Return(nil, notFoundError),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	code, _, err := manager.GetNodesWithAppID(invalidAppId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledUpdateNodeStatus_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().UpdateNodeStatus(nodeId, status).Return(nil),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	err := manager.UpdateNodeStatus(nodeId, status)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledUpdateNodeStatusWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().UpdateNodeStatus(nodeId, status).Return(notFoundError),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	err := manager.UpdateNodeStatus(nodeId, status)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledPingNodeWhenDBHasNotMatchedNode_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	code, err := manager.PingNode(nodeId, "")

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledPingNodeWithInvalidBody_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	invalidKeyBody := `{"key":"value"}`
	code, err := manager.PingNode(nodeId, invalidKeyBody)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", err.Error())
	case errors.InvalidJSON:
	}
}

func TestCalledPingNodeWithInvalidValueBody_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	invalidValueBody := `{"interval":"value"}`
	code, err := manager.PingNode(nodeId, invalidValueBody)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", err.Error())
	case errors.InvalidJSON:
	}
}

func TestCalledGetNodeConfiguration_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	code, res, err := manager.GetNodeConfiguration(nodeId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(res, config) {
		t.Error()
	}
}

func TestCalledGetNodeConfigurationWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(nil, notFoundError),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodedDBExecutorMockObj

	code, _, err := manager.GetNodeConfiguration(nodeId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestRestore_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/device/restore"}

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, respStr),
	)

	httpExecutor = msgMockObj
	nodeDbExecutor = nodedDBExecutorMockObj

	code, err := manager.Restore(nodeId)

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestRestoreWhenGetNodeFailed_ExpectReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(nil, notFoundError),
	)

	nodeDbExecutor = nodedDBExecutorMockObj

	code, err := manager.Restore(nodeId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestReboot_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/device/reboot"}

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, respStr),
	)

	httpExecutor = msgMockObj
	nodeDbExecutor = nodedDBExecutorMockObj

	code, err := manager.Reboot(nodeId)

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestRebootWhenGetNodeFailed_ExpectReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(nil, notFoundError),
	)

	nodeDbExecutor = nodedDBExecutorMockObj

	code, err := manager.Reboot(nodeId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledSetNodeConfiguration_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgMockObj := msgmocks.NewMockCommand(ctrl)
	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	jsonBody, _ := json.Marshal(config)
	jsonNodeData, _ := json.Marshal(node)
	nodeDataMap, _ := util.ConvertJsonToMap(string(jsonNodeData))
	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/device/configuration"}

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(nodeDataMap, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, jsonBody).Return(respCode, respStr),
		nodedDBExecutorMockObj.EXPECT().UpdateNodeConfiguration(nodeId, gomock.Any()).Return(nil),
	)
	// pass mockObj to a real object.
	httpExecutor = msgMockObj
	nodeDbExecutor = nodedDBExecutorMockObj

	code, err := manager.SetNodeConfiguration(nodeId, string(jsonBody))

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledSetNodeConfigurationWhenDBHasNotMatchedNode_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgMockObj := msgmocks.NewMockCommand(ctrl)
	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	httpExecutor = msgMockObj
	nodeDbExecutor = nodedDBExecutorMockObj

	body, _ := json.Marshal(config)
	code, err := manager.SetNodeConfiguration(nodeId, string(body))

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledSetNodeConfigurationWhenFailedToUpdateConfiguration_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgMockObj := msgmocks.NewMockCommand(ctrl)
	nodedDBExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	jsonBody, _ := json.Marshal(config)
	jsonNodeData, _ := json.Marshal(node)
	nodeDataMap, _ := util.ConvertJsonToMap(string(jsonNodeData))
	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/device/configuration"}

	gomock.InOrder(
		nodedDBExecutorMockObj.EXPECT().GetNode(nodeId).Return(nodeDataMap, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, jsonBody).Return(respCode, respStr),
		nodedDBExecutorMockObj.EXPECT().UpdateNodeConfiguration(nodeId, gomock.Any()).Return(notFoundError),
	)
	// pass mockObj to a real object.
	httpExecutor = msgMockObj
	nodeDbExecutor = nodedDBExecutorMockObj

	code, err := manager.SetNodeConfiguration(nodeId, string(jsonBody))

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}
