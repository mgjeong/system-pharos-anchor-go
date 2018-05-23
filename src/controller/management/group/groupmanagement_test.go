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
	groupdbmocks "db/mongo/group/mocks"
	nodedbmocks "db/mongo/node/mocks"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

const (
	appId     = "000000000000000000000000"
	nodeId    = "000000000000000000000001"
	groupId   = "000000000000000000000002"
	host      = "192.168.0.1"
	port      = "8888"
	groupName = "testGroup"
)

var (
	node = map[string]interface{}{
		"id":   nodeId,
		"host": host,
		"port": port,
		"apps": []string{appId},
	}
	members = []map[string]interface{}{node, node}
	address = map[string]interface{}{
		"host": host,
		"port": port,
	}
	membersAddress = []map[string]interface{}{address, address}
	group          = map[string]interface{}{
		"id":      groupId,
		"name":    groupName,
		"members": []string{},
	}

	body                   = `{"description":"description"}`
	respCode               = []int{results.OK, results.OK}
	partialSuccessRespCode = []int{results.OK, results.ERROR}
	errorRespCode          = []int{results.ERROR, results.ERROR}
	invalidRespStr         = []string{`{"invalidJson"}`}
	notFoundError          = errors.NotFound{}
	connectionError        = errors.DBConnectionError{}
)

var manager Command

func init() {
	manager = Executor{}
}

func TestCalledCreateGroup_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().CreateGroup(groupName).Return(group, nil),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	body := `{"name":"testGroup"}`
	code, res, err := manager.CreateGroup(body)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(group, res) {
		t.Errorf("Expected res: %s, actual res: %s", group, res)
	}
}

func TestCalledCreateGroupWhenFailedToInsertGroupToDB_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().CreateGroup(groupName).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	body := `{"name":"testGroup"}`
	code, _, err := manager.CreateGroup(body)

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

func TestCalledGetGroup_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroup(groupId).Return(group, nil),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, res, err := manager.GetGroup(groupId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(group, res) {
		t.Errorf("Expected res: %s, actual res: %s", group, res)
	}
}

func TestCalledGetGroupWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroup(groupId).Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, _, err := manager.GetGroup(groupId)

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

func TestCalledGetGroups_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groups := []map[string]interface{}{group}

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroups().Return(groups, nil),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, res, err := manager.GetGroups()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(groups, res["groups"].([]map[string]interface{})) {
		t.Errorf("Expected res: %s, actual res: %s", groups, res["groups"].([]map[string]interface{}))
	}
}

func TestCalledGetGroupsWhenFailedToGetGroupsFromDB_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroups().Return(nil, notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, _, err := manager.GetGroups()

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

func TestCalledJoinGroup_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	nodeDbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodeDbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		groupDbExecutorMockObj.EXPECT().JoinGroup(groupId, nodeId).Return(nil),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	nodeDbExecutor = nodeDbExecutorMockObj

	nodes := `{"nodes":["000000000000000000000001"]}`
	code, _, err := manager.JoinGroup(groupId, nodes)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledJoinGroupWithInvalidRequestBody_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidJsonStr := `{"invalidJson"}`
	code, _, err := manager.JoinGroup(groupId, invalidJsonStr)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InvalidParamError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InvalidParamError", err.Error())
	case errors.InvalidJSON:
	}
}

func TestCalledJoinGroupWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)
	nodeDbExecutorMockObj := nodedbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodeDbExecutorMockObj.EXPECT().GetNode(nodeId).Return(node, nil),
		groupDbExecutorMockObj.EXPECT().JoinGroup(groupId, nodeId).Return(notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj
	nodeDbExecutor = nodeDbExecutorMockObj

	nodes := `{"nodes":["000000000000000000000001"]}`
	code, _, err := manager.JoinGroup(groupId, nodes)

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

func TestCalledLeaveGroup_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().LeaveGroup(groupId, nodeId).Return(nil),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	nodes := `{"nodes":["000000000000000000000001"]}`
	code, _, err := manager.LeaveGroup(groupId, nodes)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledLeaveGroupWithInvalidRequestBody_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidJsonStr := `{"invalidJson"}`
	code, _, err := manager.LeaveGroup(groupId, invalidJsonStr)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InvalidParamError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InvalidParamError", err.Error())
	case errors.InvalidJSON:
	}
}

func TestCalledLeaveGroupWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().LeaveGroup(groupId, nodeId).Return(notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	nodes := `{"nodes":["000000000000000000000001"]}`
	code, _, err := manager.LeaveGroup(groupId, nodes)

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

func TestCalledDeleteGroup_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().DeleteGroup(groupId).Return(nil),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, _, err := manager.DeleteGroup(groupId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledDeleteGroupWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupdbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().DeleteGroup(groupId).Return(notFoundError),
	)
	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	code, _, err := manager.DeleteGroup(groupId)

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
