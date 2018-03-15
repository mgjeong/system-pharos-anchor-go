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
	nodedbmocks "db/mongo/node/mocks"
	"github.com/golang/mock/gomock"
	msgmocks "messenger/mocks"
	"reflect"
	"testing"
)

const (
	nodeId = "000000000000000000000001"
	appId = "000000000000000000001"
	testIp = "127.0.0.1"
	port   = "48098"
)

var (
	node = map[string]interface{}{
		"id":   nodeId,
		"ip":   testIp,
		"apps": []string{},
	}

	respCode         = []int{results.OK}
	respStr          = []string{`{"cpu":"00%","mem":"00%","disk":"00%","network":"00%"}`}
	errorRespCode    = []int{results.ERROR}
	invalidRespStr   = []string{`{"invalidJson"}`}
	notFoundError    = errors.NotFound{}
	connectionError  = errors.DBConnectionError{}
	invalidJsonError = errors.InvalidJSON{}
)

var resourceMonitor Command

func init() {
	resourceMonitor = Executor{}
}

func TestGetNodeResourceInfo_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{`{"cpu":"00%","mem":"00%","disk":"00%","network":"00%"}`}
	expectedUrl := []string{"http://" + testIp + ":" + port + "/api/v1/monitoring/resource"}
	expectedRes := map[string]interface{}{
		"cpu":     "00%",
		"mem":     "00%",
		"disk":    "00%",
		"network": "00%",
	}
	dbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl, nil).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj
	code, res, err := resourceMonitor.GetNodeResourceInfo(nodeId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestGetNodeResourceInfoWithGetNodeError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(nil, errors.NotFound{}),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	code, _, err := resourceMonitor.GetNodeResourceInfo(nodeId)

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

func TestGetNodeResourceInfoWhenSendHttpRequestReturnErrorCode_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + testIp + ":" + port + "/api/v1/monitoring/resource"}

	dbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl, nil).Return(errorRespCode, respStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj
	code, _, err := resourceMonitor.GetNodeResourceInfo(nodeId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}
}

func TestGetNodeResourceInfoWhenSendhttpRequestReturnErrorCodeAndInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + testIp + ":" + port + "/api/v1/monitoring/resource"}

	dbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl, nil).Return(errorRespCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj
	code, _, err := resourceMonitor.GetNodeResourceInfo(nodeId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", err.Error())
	case errors.InternalServerError:
	}
}

func TestGetAppResourceInfo_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{`{"services":"test"}`}
	expectedUrl := []string{"http://" + testIp + ":" + port + "/api/v1/monitoring/apps/"+appId+"/resource"}
	expectedRes := map[string]interface{}{
		"services":   "test",
	}
	dbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl, nil).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj
	code, res, err := resourceMonitor.GetAppResourceInfo(nodeId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestGetAppResourceInfoWithGetNodeError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(nil, errors.NotFound{}),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	code, _, err := resourceMonitor.GetAppResourceInfo(nodeId, appId)

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

func TestGetAppResourceInfoWhenSendHttpRequestReturnErrorCode_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + testIp + ":" + port + "/api/v1/monitoring/apps/"+appId+"/resource"}

	dbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl, nil).Return(errorRespCode, respStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj
	code, _, err := resourceMonitor.GetAppResourceInfo(nodeId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}
}

func TestGetAppResourceInfoWhenSendhttpRequestReturnErrorCodeAndInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + testIp + ":" + port + "/api/v1/monitoring/apps/"+appId+"/resource"}

	dbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl, nil).Return(errorRespCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj
	code, _, err := resourceMonitor.GetAppResourceInfo(nodeId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", err.Error())
	case errors.InternalServerError:
	}
}