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
	appmocks "controller/management/app/mocks"
	groupmocks "controller/management/group/mocks"
	nodemocks "controller/management/node/mocks"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

const (
	imageName1 = "testImage1"
	imageName2 = "testImage2"
	appId1     = "000000000000000000000001"
	appId2     = "000000000000000000000002"
	nodeId1    = "000000000000000000000011"
	nodeId2    = "000000000000000000000022"
	groupId1   = "000000000000000000000111"
	groupId2   = "000000000000000000000222"
	groupName  = "testGroup"
	host       = "192.168.0.1"
	port       = "8888"
)

var (
	app1 = map[string]interface{}{
		"id":       appId1,
		"images":   []string{imageName1},
		"services": []string{},
	}
	app2 = map[string]interface{}{
		"id":       appId2,
		"images":   []string{imageName2},
		"services": []string{},
	}
	node1 = map[string]interface{}{
		"id":   nodeId1,
		"host": host,
		"port": port,
		"apps": []string{appId1},
	}
	node2 = map[string]interface{}{
		"id":   nodeId2,
		"host": host,
		"port": port,
		"apps": []string{appId2},
	}
	group1 = map[string]interface{}{
		"id":      groupId1,
		"name":    groupName,
		"members": []string{nodeId1},
	}
	group2 = map[string]interface{}{
		"id":      groupId2,
		"name":    groupName,
		"members": []string{nodeId2},
	}
	groups = map[string]interface{}{
		"groups": []map[string]interface{}{group1, group2},
	}
	allQuery = map[string]interface{}{
		GROUPID:   []string{groupId1},
		NODEID:    []string{nodeId1},
		APPID:     []string{appId1},
		IMAGENAME: []string{imageName1},
	}
	queryWithoutGroupId = map[string]interface{}{
		NODEID:    []string{nodeId1},
		APPID:     []string{appId1},
		IMAGENAME: []string{imageName1},
	}
	invalidQuery = map[string]interface{}{
		"invalid": []string{"value"},
	}
	notFoundError = errors.NotFound{}
)

var executor Command

func init() {
	executor = Executor{}
}

func TestCalledSearchGroupsWithInvalidQuery_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	code, _, err := executor.SearchGroups(invalidQuery)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundURL", "nil")
	}

	if code != results.ERROR {
		t.Errorf("Expected return code : %d, actual err: %d", 500, code)
	}
}

func TestCalledSearchGroupsWithAllQuery_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appExecutorMockObj := appmocks.NewMockCommand(ctrl)
	nodeExecutorMockObj := nodemocks.NewMockCommand(ctrl)
	groupExecutorMockObj := groupmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupExecutorMockObj.EXPECT().GetGroup(groupId1).Return(results.OK, group1, nil),
		nodeExecutorMockObj.EXPECT().GetNode(nodeId1).Return(results.OK, node1, nil),
		appExecutorMockObj.EXPECT().GetApp(appId1).Return(results.OK, app1, nil),
		nodeExecutorMockObj.EXPECT().GetNode(nodeId1).Return(results.OK, node1, nil),
	)
	// pass mockObj to a real object
	appmanagementExecutor = appExecutorMockObj
	nodemanagementExecutor = nodeExecutorMockObj
	groupmanagementExecutor = groupExecutorMockObj

	code, res, err := executor.SearchGroups(allQuery)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	expectedResult := make(map[string]interface{})
	expectedResult["groups"] = make([]map[string]interface{}, 1)
	expectedResult["groups"].([]map[string]interface{})[0] = group1

	if !reflect.DeepEqual(expectedResult, res) {
		t.Errorf("Expected res: %s\n actual res: %s", expectedResult, res)
	}
}

func TestCalledSearchGroupsWithoutGroupId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appExecutorMockObj := appmocks.NewMockCommand(ctrl)
	nodeExecutorMockObj := nodemocks.NewMockCommand(ctrl)
	groupExecutorMockObj := groupmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupExecutorMockObj.EXPECT().GetGroups().Return(results.OK, groups, nil),
		nodeExecutorMockObj.EXPECT().GetNode(nodeId1).Return(results.OK, node1, nil),
		appExecutorMockObj.EXPECT().GetApp(appId1).Return(results.OK, app1, nil),
		nodeExecutorMockObj.EXPECT().GetNode(nodeId2).Return(results.OK, node2, nil),
		appExecutorMockObj.EXPECT().GetApp(appId2).Return(results.OK, app2, nil),
		nodeExecutorMockObj.EXPECT().GetNode(nodeId1).Return(results.OK, node1, nil),
	)
	// pass mockObj to a real object
	appmanagementExecutor = appExecutorMockObj
	nodemanagementExecutor = nodeExecutorMockObj
	groupmanagementExecutor = groupExecutorMockObj

	code, res, err := executor.SearchGroups(queryWithoutGroupId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	expectedResult := make(map[string]interface{})
	expectedResult["groups"] = make([]map[string]interface{}, 1)
	expectedResult["groups"].([]map[string]interface{})[0] = group1

	if !reflect.DeepEqual(expectedResult, res) {
		t.Errorf("Expected res: %s\n actual res: %s", expectedResult, res)
	}
}
