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
package agent

import (
	"commons/errors"
	"commons/results"
	"github.com/golang/mock/gomock"
	agentdbmocks "db/mongo/agent/mocks"
	msgmocks "messenger/mocks"
	"reflect"
	"testing"
)

const (
	AGENTID = "000000000000000000000001"
	IP      = "127.0.0.1"
	PORT    = "48098"
)

var (
	agent = map[string]interface{}{
		"id":   AGENTID,
		"ip":   IP,
		"apps": []string{},
	}

	respCode         = []int{results.OK}
	respStr          = []string{`{"os":"os","processor":"processor","cpu":"00%","mem":"00%","disk":"00%"}`}
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

func TestGetResourceInfo_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + IP + ":" + PORT + "/api/v1/resource"}
	expectedRes := map[string]interface{}{
		"os":        "os",
		"processor": "processor",
		"cpu":       "00%",
		"mem":       "00%",
		"disk":      "00%",
	}
	dbExecutorMockObj := agentdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetAgent(AGENTID).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	agentDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj
	code, res, err := resourceMonitor.GetResourceInfo(AGENTID)

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

func TestGetResourceInfoWithGetAgentError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbExecutorMockObj := agentdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetAgent(AGENTID).Return(nil, errors.NotFound{}),
	)
	// pass mockObj to a real object.
	agentDbExecutor = dbExecutorMockObj
	code, _, err := resourceMonitor.GetResourceInfo(AGENTID)

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

func TestGetResourceInfoWhenSendhttpRequestReturnErrorCode_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + IP + ":" + PORT + "/api/v1/resource"}

	dbExecutorMockObj := agentdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetAgent(AGENTID).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl).Return(errorRespCode, respStr),
	)
	// pass mockObj to a real object.
	agentDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj
	code, _, err := resourceMonitor.GetResourceInfo(AGENTID)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}
}

func TestGetResourceInfoWhenSendhttpRequestReturnErrorCodeAndInvalidResponse_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + IP + ":" + PORT + "/api/v1/resource"}

	dbExecutorMockObj := agentdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetAgent(AGENTID).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl).Return(errorRespCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	agentDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj
	code, _, err := resourceMonitor.GetResourceInfo(AGENTID)

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

func TestGetPerformanceInfo_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{`{"cpu":"00%","mem":"00%","disk":"00%"}`}
	expectedUrl := []string{"http://" + IP + ":" + PORT + "/api/v1/resource/performance"}
	expectedRes := map[string]interface{}{
		"cpu":  "00%",
		"mem":  "00%",
		"disk": "00%",
	}
	dbExecutorMockObj := agentdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetAgent(AGENTID).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	agentDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj
	code, res, err := resourceMonitor.GetPerformanceInfo(AGENTID)

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

func TestGetPerformanceInfoWithGetAgentError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbExecutorMockObj := agentdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetAgent(AGENTID).Return(nil, errors.NotFound{}),
	)
	// pass mockObj to a real object.
	agentDbExecutor = dbExecutorMockObj
	code, _, err := resourceMonitor.GetPerformanceInfo(AGENTID)

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

func TestGetPerformanceInfoWhenSendHttpRequestReturnErrorCode_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + IP + ":" + PORT + "/api/v1/resource/performance"}

	dbExecutorMockObj := agentdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetAgent(AGENTID).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl).Return(errorRespCode, respStr),
	)
	// pass mockObj to a real object.
	agentDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj
	code, _, err := resourceMonitor.GetPerformanceInfo(AGENTID)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}
}

func TestGetPerformanceInfoWhenSendhttpRequestReturnErrorCodeAndInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + IP + ":" + PORT + "/api/v1/resource/performance"}

	dbExecutorMockObj := agentdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetAgent(AGENTID).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl).Return(errorRespCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	agentDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj
	code, _, err := resourceMonitor.GetPerformanceInfo(AGENTID)

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
