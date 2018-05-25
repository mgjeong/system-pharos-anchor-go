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
	appdbmocks "db/mongo/app/mocks"
	appeventdbmocks "db/mongo/event/app/mocks"
	subsdbmocks "db/mongo/event/subscriber/mocks"
	dbmocks "db/mongo/node/mocks"
	"github.com/golang/mock/gomock"
	msgmocks "messenger/mocks"
	"reflect"
	"testing"
)

const (
	status = "connected"
	appId  = "000000000000000000000000"
	nodeId = "000000000000000000000001"
	ip     = "127.0.0.1"
	port   = "48098"
)

var (
	node = map[string]interface{}{
		"id":   nodeId,
		"ip":   ip,
		"apps": []string{},
	}
	address = []map[string]interface{}{
		map[string]interface{}{
			"ip": ip,
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

var executor Command

func init() {
	executor = Executor{}
}

func TestCalledDeployAppWithEventQuery_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testEventUrl := []string{"http://0.0.0.0:0000"}
	testQuery := map[string]interface{}{
		EVENT: testEventUrl,
	}
	respStr := []string{`{"id":"000000000000000000000000", "description":"description"}`}
	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/deploy"}
	expectedRes := map[string]interface{}{
		"id":          "000000000000000000000000",
		"description": "description",
	}

	subsDbMockObj := subsdbmocks.NewMockCommand(ctrl)
	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)
	appDbMockObj := appdbmocks.NewMockCommand(ctrl)
	appEventDbMockObj := appeventdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		subsDbMockObj.EXPECT().AddSubscriber(gomock.Any(), APP, testEventUrl[0],
			[]string{PULLED, CREATED, STARTED}, gomock.Any(), make(map[string][]string)).Return(nil),
		appEventDbMockObj.EXPECT().AddEvent(gomock.Any(), gomock.Any(), []string{nodeId}).Return(nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, gomock.Any(), []byte(body)).Return(respCode, respStr),
		subsDbMockObj.EXPECT().DeleteSubscriber(gomock.Any()),
		appEventDbMockObj.EXPECT().DeleteEvent(gomock.Any()),
		appDbMockObj.EXPECT().AddApp(appId, []byte("description")).Return(nil),
		dbExecutorMockObj.EXPECT().AddAppToNode(nodeId, appId).Return(nil),
	)
	// pass mockObj to a real object.
	subsDbExecutor = subsDbMockObj
	appEventDbExecutor = appEventDbMockObj
	appDbExecutor = appDbMockObj
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, res, err := executor.DeployApp(nodeId, body, testQuery)

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

func TestCalledDeployAppWithEventQueryWhenAddSubscriberFailed_ExpectReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testEventUrl := []string{"http://0.0.0.0:0000"}
	testQuery := map[string]interface{}{
		EVENT: testEventUrl,
	}

	subsDbMockObj := subsdbmocks.NewMockCommand(ctrl)
	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		subsDbMockObj.EXPECT().AddSubscriber(gomock.Any(), APP, testEventUrl[0],
			[]string{PULLED, CREATED, STARTED}, gomock.Any(), make(map[string][]string)).Return(errors.Unknown{}),
	)
	// pass mockObj to a real object.
	subsDbExecutor = subsDbMockObj
	nodeDbExecutor = dbExecutorMockObj

	_, _, err := executor.DeployApp(nodeId, body, testQuery)

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "Unknwon", err.Error())
	case errors.Unknown:
	}
}

func TestCalledDeployAppWithEventQueryWhenAddEventFailed_ExpectReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testEventUrl := []string{"http://0.0.0.0:0000"}
	testQuery := map[string]interface{}{
		EVENT: testEventUrl,
	}

	subsDbMockObj := subsdbmocks.NewMockCommand(ctrl)
	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	appEventDbMockObj := appeventdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		subsDbMockObj.EXPECT().AddSubscriber(gomock.Any(), APP, testEventUrl[0],
			[]string{PULLED, CREATED, STARTED}, gomock.Any(), make(map[string][]string)).Return(nil),
		appEventDbMockObj.EXPECT().AddEvent(gomock.Any(), gomock.Any(), []string{nodeId}).Return(errors.Unknown{}),
		subsDbMockObj.EXPECT().DeleteSubscriber(gomock.Any()).Return(nil),
	)
	// pass mockObj to a real object.
	subsDbExecutor = subsDbMockObj
	nodeDbExecutor = dbExecutorMockObj
	appEventDbExecutor = appEventDbMockObj

	_, _, err := executor.DeployApp(nodeId, body, testQuery)

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "Unknwon", err.Error())
	case errors.Unknown:
	}
}

func TestCalledDeployApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{`{"id":"000000000000000000000000", "description":"description"}`}
	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/deploy"}
	expectedRes := map[string]interface{}{
		"id":          "000000000000000000000000",
		"description": "description",
	}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)
	appDbMockObj := appdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, []byte(body)).Return(respCode, respStr),
		appDbMockObj.EXPECT().AddApp(appId, []byte("description")).Return(nil),
		dbExecutorMockObj.EXPECT().AddAppToNode(nodeId, appId).Return(nil),
	)
	// pass mockObj to a real object.
	appDbExecutor = appDbMockObj
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, res, err := executor.DeployApp(nodeId, body, nil)

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

func TestCalledDeployAppWhenDBHasNotMatchedNode_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj

	code, _, err := executor.DeployApp(nodeId, body, nil)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/deploy"}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, []byte(body)).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.DeployApp(nodeId, body, nil)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/deploy"}
	respStr := []string{`{"id":"000000000000000000000000", "description":"description"}`}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)
	appDbMockObj := appdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, []byte(body)).Return(respCode, respStr),
		appDbMockObj.EXPECT().AddApp(appId, []byte("description")).Return(nil),
		dbExecutorMockObj.EXPECT().AddAppToNode(nodeId, appId).Return(notFoundError),
	)
	// pass mockObj to a real object.
	appDbExecutor = appDbMockObj
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.DeployApp(nodeId, body, nil)

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
	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps"}
	expectedRes := map[string]interface{}{
		"description": "description",
	}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl, nil).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, res, err := executor.GetApps(nodeId)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps"}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl, nil).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.GetApps(nodeId)

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

func TestCalledGetAppsWhenDBHasNotMatchedNode_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNode(nodeId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj

	code, _, err := executor.GetApps(nodeId)

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
	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId}
	expectedRes := map[string]interface{}{
		"description": "description",
	}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl, nil).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, res, err := executor.GetApp(nodeId, appId)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl, nil).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.GetApp(nodeId, appId)

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

func TestCalledGetAppWhenDBHasNotMatchedNode_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj

	code, _, err := executor.GetApp(nodeId, appId)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, []byte(body)).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.UpdateAppInfo(nodeId, appId, body)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledUpdateAppInfoWhenDBHasNotMatchedNode_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.UpdateAppInfo(nodeId, appId, body)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, []byte(body)).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.UpdateAppInfo(nodeId, appId, body)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId + "/update"}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.UpdateApp(nodeId, appId, nil)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledUpdateAppWhenDBHasNotMatchedNode_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.UpdateApp(nodeId, appId, nil)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId + "/update"}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.UpdateApp(nodeId, appId, nil)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId + "/start"}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.StartApp(nodeId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledStartAppWhenDBHasNotMatchedNode_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj

	code, _, err := executor.StartApp(nodeId, appId)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId + "/start"}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.StartApp(nodeId, appId)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId + "/stop"}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.StopApp(nodeId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledStopAppWhenDBHasNotMatchedNode_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.StopApp(nodeId, appId)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId + "/stop"}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.StopApp(nodeId, appId)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)
	appDbMockObj := appdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("DELETE", expectedUrl, nil).Return(respCode, respStr),
		dbExecutorMockObj.EXPECT().DeleteAppFromNode(nodeId, appId).Return(nil),
		appDbMockObj.EXPECT().DeleteApp(appId).Return(nil),
	)
	// pass mockObj to a real object.
	appDbExecutor = appDbMockObj
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.DeleteApp(nodeId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledDeleteAppWhenDBHasNotMatchedNode_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.DeleteApp(nodeId, appId)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("DELETE", expectedUrl, nil).Return(errorRespCode, respStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.DeleteApp(nodeId, appId)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("DELETE", expectedUrl, nil).Return(errorRespCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.DeleteApp(nodeId, appId)

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

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/management/apps/" + appId}

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetNodeByAppID(nodeId, appId).Return(node, nil),
		msgMockObj.EXPECT().SendHttpRequest("DELETE", expectedUrl, nil).Return(respCode, nil),
		dbExecutorMockObj.EXPECT().DeleteAppFromNode(nodeId, appId).Return(notFoundError),
	)
	// pass mockObj to a real object.
	nodeDbExecutor = dbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.DeleteApp(nodeId, appId)

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

func TestGenerateRandStringBytes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testStrLen := 50
	ret := generateRandStringBytes(testStrLen)

	if len(ret) != testStrLen {
		t.Errorf("Expected length of string : %d, actual length of string : %d", testStrLen, len(ret))
	}
}
