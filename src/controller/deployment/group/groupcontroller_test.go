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
	"commons/errors"
	"commons/results"
	appdbmocks "db/mongo/app/mocks"
	groupdbmocks "db/mongo/group/mocks"
	nodedbmocks "db/mongo/node/mocks"
	"github.com/golang/mock/gomock"
	msgmocks "messenger/mocks"
	"reflect"
	"testing"
)

const (
	appId   = "000000000000000000000000"
	nodeId  = "000000000000000000000001"
	groupId = "000000000000000000000002"
	ip      = "192.168.0.1"
	port    = "48098"
)

var (
	node = map[string]interface{}{
		"id":   nodeId,
		"ip":   ip,
		"apps": []string{appId},
	}
	members = []map[string]interface{}{node, node}
	address = map[string]interface{}{
		"ip": ip,
	}
	membersAddress = []map[string]interface{}{address, address}
	group          = map[string]interface{}{
		"id":      groupId,
		"members": []string{nodeId, nodeId},
	}

	body                   = `{"description":"description"}`
	deployUrl              = "http://" + ip + ":" + port + "/api/v1/management/apps/deploy"
	baseUrl                = "http://" + ip + ":" + port + "/api/v1/management/apps/" + appId
	respCode               = []int{results.OK, results.OK}
	partialSuccessRespCode = []int{results.OK, results.ERROR}
	errorRespCode          = []int{results.ERROR, results.ERROR}
	invalidRespStr         = []string{`{"invalidJson"}`}
	notFoundError          = errors.NotFound{}
	connectionError        = errors.DBConnectionError{}
)

var executor Command

func init() {
	executor = Executor{}
}

func TestCalledDeployApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{
		`{"id":"000000000000000000000000", "description":"description"}`,
		`{"id":"000000000000000000000000", "description":"description"}`,
	}
	expectedUrl := []string{deployUrl, deployUrl}
	expectedRes := map[string]interface{}{
		"id": "000000000000000000000000",
	}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	nodeDbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	appDbExecutorMockObj := appdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembers(groupId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, []byte(body)).Return(respCode, respStr),
		appDbExecutorMockObj.EXPECT().AddApp(appId, gomock.Any()).Return(nil),
		nodeDbExecutorMockObj.EXPECT().AddAppToNode(nodeId, appId).Return(nil),
		appDbExecutorMockObj.EXPECT().AddApp(appId, gomock.Any()).Return(nil),
		nodeDbExecutorMockObj.EXPECT().AddAppToNode(nodeId, appId).Return(nil),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	appDbExecutor = appDbExecutorMockObj
	nodeDbExecutor = nodeDbExecutorMockObj
	httpExecutor = msgMockObj

	code, res, err := executor.DeployApp(groupId, body)

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

func TestCalledDeployAppWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembers(groupId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, _, err := executor.DeployApp(groupId, body)

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

	expectedUrl := []string{deployUrl, deployUrl}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembers(groupId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, []byte(body)).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.DeployApp(groupId, body)

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

	respStr := []string{`{"id":"000000000000000000000000", "description":"description"}`}
	expectedUrl := []string{deployUrl, deployUrl}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	nodeDbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	appDbExecutorMockObj := appdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembers(groupId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, []byte(body)).Return(respCode, respStr),
		appDbExecutorMockObj.EXPECT().AddApp(appId, []byte("description")).Return(nil).AnyTimes(),
		nodeDbExecutorMockObj.EXPECT().AddAppToNode(nodeId, appId).Return(notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	appDbExecutor = appDbExecutorMockObj
	nodeDbExecutor = nodeDbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.DeployApp(groupId, body)

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

func TestCalledDeployAppWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"id":"000000000000000000000000", "description":"description"}`, `{"message":"errorMsg"}`}
	expectedUrl := []string{deployUrl, deployUrl}
	expectedRes := map[string]interface{}{
		"id": "000000000000000000000000",
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":   nodeId,
				"code": results.OK,
			},
			map[string]interface{}{
				"id":      nodeId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	nodeDbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	appDbExecutorMockObj := appdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembers(groupId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, []byte(body)).Return(partialSuccessRespCode, partialSuccessRespStr),
		appDbExecutorMockObj.EXPECT().AddApp(appId, []byte("description")).Return(nil).AnyTimes(),
		nodeDbExecutorMockObj.EXPECT().AddAppToNode(nodeId, appId).Return(nil),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	appDbExecutor = appDbExecutorMockObj
	nodeDbExecutor = nodeDbExecutorMockObj
	httpExecutor = msgMockObj

	code, res, err := executor.DeployApp(groupId, body)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetApps_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedRes := map[string]interface{}{
		"apps": []map[string]interface{}{{
			"id":      appId,
			"members": []string{nodeId, nodeId},
		}},
	}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembers(groupId).Return(members, nil),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, res, err := executor.GetApps(groupId)

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

func TestCalledGetAppsWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembers(groupId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, _, err := executor.GetApps(groupId)

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

func TestCalledGetApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{`{"description":"description"}`, `{"description":"description"}`}
	expectedUrl := []string{baseUrl, baseUrl}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{{
			"description": "description",
			"id":          members[0]["id"],
		},
			{
				"description": "description",
				"id":          members[0]["id"],
			}},
	}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl, nil).Return(respCode, respStr),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, res, err := executor.GetApp(groupId, appId)

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

func TestCalledGetAppWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, _, err := executor.GetApp(groupId, appId)

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

func TestCalledGetAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidRespStr := []string{`{"invalidJson"}`, `{"invalidJson"}`}
	expectedUrl := []string{baseUrl, baseUrl}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl, nil).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.GetApp(groupId, appId)

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

func TestCalledGetAppWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"description": "description"}`, `{"message":"errorMsg"}`}
	expectedUrl := []string{baseUrl, baseUrl}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":          nodeId,
				"code":        results.OK,
				"description": "description",
			},
			map[string]interface{}{
				"id":      nodeId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("GET", expectedUrl, nil).Return(partialSuccessRespCode, partialSuccessRespStr),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, res, err := executor.GetApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledUpdateAppInfo_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{baseUrl, baseUrl}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, []byte(body)).Return(respCode, nil),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.UpdateAppInfo(groupId, appId, body)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledUpdateAppInfoWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, _, err := executor.UpdateAppInfo(groupId, appId, body)

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

	invalidRespStr := []string{`{"invalidJson"}`, `{"invalidJson"}`}
	expectedUrl := []string{baseUrl, baseUrl}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, []byte(body)).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.UpdateAppInfo(groupId, appId, body)

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

func TestCalledUpdateAppInfoWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"message": "successMsg"}`, `{"message":"errorMsg"}`}
	expectedUrl := []string{baseUrl, baseUrl}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":   nodeId,
				"code": results.OK,
			},
			map[string]interface{}{
				"id":      nodeId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, []byte(body)).Return(partialSuccessRespCode, partialSuccessRespStr),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, res, err := executor.UpdateAppInfo(groupId, appId, body)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledUpdateApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{baseUrl + "/update", baseUrl + "/update"}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, nil),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.UpdateApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledUpdateAppWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, _, err := executor.UpdateApp(groupId, appId)

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

	invalidRespStr := []string{`{"invalidJson"}`, `{"invalidJson"}`}
	expectedUrl := []string{baseUrl + "/update", baseUrl + "/update"}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.UpdateApp(groupId, appId)

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

func TestCalledUpdateAppWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"message": "successMsg"}`, `{"message":"errorMsg"}`}
	expectedUrl := []string{baseUrl + "/update", baseUrl + "/update"}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":   nodeId,
				"code": results.OK,
			},
			map[string]interface{}{
				"id":      nodeId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(partialSuccessRespCode, partialSuccessRespStr),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, res, err := executor.UpdateApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledStartApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{baseUrl + "/start", baseUrl + "/start"}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, nil),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.StartApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledStartAppWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, _, err := executor.StartApp(groupId, appId)

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

func TestCalledStartAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidRespStr := []string{`{"invalidJson"}`, `{"invalidJson"}`}
	expectedUrl := []string{baseUrl + "/start", baseUrl + "/start"}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.StartApp(groupId, appId)

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

func TestCalledStartAppWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"message": "successMsg"}`, `{"message":"errorMsg"}`}
	expectedUrl := []string{baseUrl + "/start", baseUrl + "/start"}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":   nodeId,
				"code": results.OK,
			},
			map[string]interface{}{
				"id":      nodeId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(partialSuccessRespCode, partialSuccessRespStr),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, res, err := executor.StartApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledStopApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{baseUrl + "/stop", baseUrl + "/stop"}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, nil),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.StopApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledStopAppWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, _, err := executor.StopApp(groupId, appId)

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

func TestCalledStopAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidRespStr := []string{`{"invalidJson"}`, `{"invalidJson"}`}
	expectedUrl := []string{baseUrl + "/stop", baseUrl + "/stop"}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.StopApp(groupId, appId)

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

func TestCalledStopAppWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"message": "successMsg"}`, `{"message":"errorMsg"}`}
	expectedUrl := []string{baseUrl + "/stop", baseUrl + "/stop"}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":   nodeId,
				"code": results.OK,
			},
			map[string]interface{}{
				"id":      nodeId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil).Return(partialSuccessRespCode, partialSuccessRespStr),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, res, err := executor.StopApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledDeleteApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{baseUrl, baseUrl}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	nodeDbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	appDbExecutorMockObj := appdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("DELETE", expectedUrl, nil).Return(respCode, nil),
		nodeDbExecutorMockObj.EXPECT().DeleteAppFromNode(nodeId, appId).Return(nil),
		appDbExecutorMockObj.EXPECT().DeleteApp(appId).Return(nil),
		nodeDbExecutorMockObj.EXPECT().DeleteAppFromNode(nodeId, appId).Return(nil),
		appDbExecutorMockObj.EXPECT().DeleteApp(appId).Return(nil).AnyTimes(),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	appDbExecutor = appDbExecutorMockObj
	nodeDbExecutor = nodeDbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.DeleteApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledDeleteAppWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, _, err := executor.DeleteApp(groupId, appId)

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

func TestCalledDeleteAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidRespStr := []string{`{"invalidJson"}`, `{"invalidJson"}`}
	expectedUrl := []string{baseUrl, baseUrl}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("DELETE", expectedUrl, nil).Return(respCode, invalidRespStr),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	httpExecutor = msgMockObj

	code, _, err := executor.DeleteApp(groupId, appId)

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

func TestCalledDeleteAppWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"message": "successMsg"}`, `{"message":"errorMsg"}`}
	expectedUrl := []string{baseUrl, baseUrl}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":   nodeId,
				"code": results.OK,
			},
			map[string]interface{}{
				"id":      nodeId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	nodeDbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)
	appDbExecutorMockObj := appdbmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().SendHttpRequest("DELETE", expectedUrl, nil).Return(partialSuccessRespCode, partialSuccessRespStr),
		nodeDbExecutorMockObj.EXPECT().DeleteAppFromNode(nodeId, appId).Return(nil),
		appDbExecutorMockObj.EXPECT().DeleteApp(appId).Return(nil).AnyTimes(),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	appDbExecutor = appDbExecutorMockObj
	nodeDbExecutor = nodeDbExecutorMockObj
	httpExecutor = msgMockObj

	code, res, err := executor.DeleteApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}
