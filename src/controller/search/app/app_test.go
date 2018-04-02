/*******************************************************************************
 * Copyright 2018 Samsung Electronics All Rights Reserved.
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
package app

import (
	"commons/errors"
	"commons/results"
	appDbmocks "db/mongo/app/mocks"
	nodeDbmocks "db/mongo/node/mocks"
	groupDbmocks "db/mongo/group/mocks"
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
	apps = map[string]interface{}{
		APPS: []map[string]interface{}{app1, app2},
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
	queryWithoutAppId = map[string]interface{}{
		GROUPID:   []string{groupId1},
		NODEID:    []string{nodeId1},
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

func TestCalledSearchAppsWithInvalidQuery_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	code, _, err := executor.Search(invalidQuery)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundURL", "nil")
	}

	if code != results.ERROR {
		t.Errorf("Expected return code : %d, actual err: %d", 500, code)
	}
}

func TestCalledSearchAppsWithAllQuery_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appExecutorMockObj := appDbmocks.NewMockCommand(ctrl)
	nodeExecutorMockObj := nodeDbmocks.NewMockCommand(ctrl)
	groupExecutorMockObj := groupDbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		appExecutorMockObj.EXPECT().GetApp(appId1).Return(app1, nil),
		nodeExecutorMockObj.EXPECT().GetNode(nodeId1).Return(node1, nil),
		groupExecutorMockObj.EXPECT().GetGroup(groupId1).Return(group1, nil),
		nodeExecutorMockObj.EXPECT().GetNode(nodeId1).Return(node1, nil),
	)
	// pass mockObj to a real object
	appDbExecutor = appExecutorMockObj
	nodeDbExecutor = nodeExecutorMockObj
	groupDbExecutor = groupExecutorMockObj

	code, res, err := executor.Search(allQuery)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	expectedResult := make(map[string]interface{})
	expectedResult[APPS] = make([]map[string]interface{}, 1)
	expectedResult[APPS].([]map[string]interface{})[0] = app1

	if !reflect.DeepEqual(expectedResult, res) {
		t.Errorf("Expected res: %s\n actual res: %s", expectedResult, res)
	}
}

func TestCalledSearchAppsWithoutAppId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appExecutorMockObj := appDbmocks.NewMockCommand(ctrl)
	nodeExecutorMockObj := nodeDbmocks.NewMockCommand(ctrl)
	groupExecutorMockObj := groupDbmocks.NewMockCommand(ctrl)
	
	apps := make([]map[string]interface{}, 2)
	apps[0] = app1
	apps[1] = app2

	gomock.InOrder(
		appExecutorMockObj.EXPECT().GetApps().Return(apps, nil),
		nodeExecutorMockObj.EXPECT().GetNode(nodeId1).Return(node1, nil),
		groupExecutorMockObj.EXPECT().GetGroup(groupId1).Return(group1, nil),
		nodeExecutorMockObj.EXPECT().GetNode(nodeId1).Return(node1, nil),
	)
	// pass mockObj to a real object
	appDbExecutor = appExecutorMockObj
	nodeDbExecutor = nodeExecutorMockObj
	groupDbExecutor = groupExecutorMockObj

	code, res, err := executor.Search(queryWithoutAppId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	expectedResult := make(map[string]interface{})
	expectedResult[APPS] = make([]map[string]interface{}, 1)
	expectedResult[APPS].([]map[string]interface{})[0] = app1

	if !reflect.DeepEqual(expectedResult, res) {
		t.Errorf("Expected res: %s\n actual res: %s", expectedResult, res)
	}
}

func TestCalledSearchAppsWithoutGroupId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appExecutorMockObj := appDbmocks.NewMockCommand(ctrl)
	nodeExecutorMockObj := nodeDbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		appExecutorMockObj.EXPECT().GetApp(appId1).Return(app1, nil),
		nodeExecutorMockObj.EXPECT().GetNode(nodeId1).Return(node1, nil),
	)
	// pass mockObj to a real object
	appDbExecutor = appExecutorMockObj
	nodeDbExecutor = nodeExecutorMockObj

	code, res, err := executor.Search(queryWithoutGroupId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	expectedResult := make(map[string]interface{})
	expectedResult[APPS] = make([]map[string]interface{}, 1)
	expectedResult[APPS].([]map[string]interface{})[0] = app1

	if !reflect.DeepEqual(expectedResult, res) {
		t.Errorf("Expected res: %s\n actual res: %s", expectedResult, res)
	}
}
