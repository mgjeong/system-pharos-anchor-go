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
	dbmocks "db/modelinterface/mocks"
	"github.com/golang/mock/gomock"
	msgmocks "messenger/mocks"
	"reflect"
	"testing"
)

const (
	status  = "connected"
	appId   = "000000000000000000000000"
	agentId = "000000000000000000000001"
	host    = "127.0.0.1"
	port    = "48098"
)

var (
	agent = map[string]interface{}{
		"id":   agentId,
		"host": host,
		"port": port,
		"apps": []string{},
	}
	address = []map[string]interface{}{
		map[string]interface{}{
			"host": host,
			"port": port,
		}}
	body             = `{"description":"description"}`
	respCode         = []int{results.OK}
	errorRespCode    = []int{results.ERROR}
	respStr          = []string{`{"response":"response"}`}
	invalidRespStr   = []string{`{"invalidJson"}`}
	notFoundError    = errors.NotFound{}
	connectionError  = errors.DBConnectionError{}
	invalidJsonError = errors.InvalidJSON{}
)

var controller DeploymentInterface

func init() {
	controller = AgentController{}
}

func TestCalledDeployApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{`{"id":"000000000000000000000000"}`}
	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/deploy"}
	expectedRes := map[string]interface{}{
		"id": "000000000000000000000000",
	}

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, []byte(body)).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().AddAppToAgent(agentId, appId).Return(nil),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, res, err := controller.DeployApp(agentId, body)

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

func TestCalledDeployAppWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj

	code, _, err := controller.DeployApp(agentId, body)

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

func TestCalledDeployAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/deploy"}
	
	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, []byte(body)).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.DeployApp(agentId, body)

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

func TestCalledDeployAppWhenFailedToAddAppIdToDB_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/deploy"}
	respStr := []string{`{"id":"000000000000000000000000"}`}

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, []byte(body)).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().AddAppToAgent(agentId, appId).Return(notFoundError),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.DeployApp(agentId, body)

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

func TestCalledGetApps_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{`{"description":"description"}`}
	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps"}
	expectedRes := map[string]interface{}{
		"description": "description",
	}

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, res, err := controller.GetApps(agentId)

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

func TestCalledGetAppsWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps"}

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.GetApps(agentId)

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

func TestCalledGetAppsWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj

	code, _, err := controller.GetApps(agentId)

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

func TestCalledGetApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{`{"description":"description"}`}
	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId}
	expectedRes := map[string]interface{}{
		"description": "description",
	}

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, res, err := controller.GetApp(agentId, appId)

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

func TestCalledGetAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId}
	
	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.GetApp(agentId, appId)

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

func TestCalledGetAppWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj

	code, _, err := controller.GetApp(agentId, appId)

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

func TestCalledUpdateAppInfo_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId}

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, []byte(body)).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.UpdateAppInfo(agentId, appId, body)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledUpdateAppInfoWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.UpdateAppInfo(agentId, appId, body)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", err.Error())
	case errors.NotFound:
	}
}

func TestCalledUpdateAppInfoWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId}
	
	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, []byte(body)).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.UpdateAppInfo(agentId, appId, body)

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

func TestCalledUpdateApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId + "/update"}
	
	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.UpdateApp(agentId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledUpdateAppWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.UpdateApp(agentId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", err.Error())
	case errors.NotFound:
	}
}

func TestCalledUpdateAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId + "/update"}
	
	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.UpdateApp(agentId, appId)

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

func TestCalledStartApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId + "/start"}
	
	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.StartApp(agentId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledStartAppWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj

	code, _, err := controller.StartApp(agentId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", err.Error())
	case errors.NotFound:
	}
}

func TestCalledStartAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId + "/start"}
	
	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.StartApp(agentId, appId)

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

func TestCalledStopApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId + "/stop"}
	
	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.StopApp(agentId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledStopAppWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.StopApp(agentId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", err.Error())
	case errors.NotFound:
	}
}

func TestCalledStopAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId + "/stop"}
	
	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.StopApp(agentId, appId)

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

func TestCalledDeleteApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId}
	
	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("DELETE", expectedUrl).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().DeleteAppFromAgent(agentId, appId).Return(nil),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.DeleteApp(agentId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledDeleteAppWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.DeleteApp(agentId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", err.Error())
	case errors.NotFound:
	}
}

func TestCalledDeleteAppWhenMessengerReturnsErrorCode_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId}
	
	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("DELETE", expectedUrl).Return(errorRespCode, respStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.DeleteApp(agentId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}
}

func TestCalledDeleteAppWhenMessengerReturnsErrorCodeWithInvalidResponse_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId}

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("DELETE", expectedUrl).Return(errorRespCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.DeleteApp(agentId, appId)

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

func TestCalledDeleteAppWhenFailedToDeleteAppIdFromDB_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + host + ":" + port + "/api/v1/apps/" + appId}

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("DELETE", expectedUrl).Return(respCode, nil),
		dbManagerMockObj.EXPECT().DeleteAppFromAgent(agentId, appId).Return(notFoundError),
	)
	// pass mockObj to a real object.
	agentDbManager = dbManagerMockObj
	httpRequester = msgMockObj

	code, _, err := controller.DeleteApp(agentId, appId)

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
