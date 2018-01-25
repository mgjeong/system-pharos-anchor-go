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
package registry

import (
	"commons/errors"
	"commons/results"
	"commons/url"
	appmocks "controller/management/app/mocks"
	nodemocks "controller/management/node/mocks"
	dbmocks "db/mongo/registry/mocks"
	"github.com/golang/mock/gomock"
	msgmocks "messenger/mocks"
	"reflect"
	"testing"
)

const (
	status            = "connected"
	registryId        = "000000000000000000000001"
	invalidregistryId = "0"
	nodeId            = "000000000000000000000001"
	refcnt            = 1
)

var (
	app = map[string]interface{}{
		"id":       appId,
		"images":   []string{},
		"services": []string{},
		"refcnt":   refcnt,
	}
	node = map[string]interface{}{
		"id":     nodeId,
		"ip":     ip,
		"apps":   []string{},
		"config": configuration,
	}
	configuration = map[string]interface{}{
		"key": "value",
	}
	registryModel = map[string]interface{}{
		"id": registryId,
		"ip": ip,
	}
	respCode           = []int{results.OK}
	respStr            = []string{`{"response":"response"}`}
	notFoundError      = errors.NotFound{}
	ip                 = "127.0.0.1"
	appId              = "000000000000000000000001"
	seperator          = "/"
	dummy_repositry    = "dummy_repository"
	dummy_host         = "dummy_host"
	dummyRegistryEvent = `{"events":[{"target":{"repository": "dummy_repository"},"request": {"host": "dummy_host"}}]}`
)

var manager Command

func init() {
	manager = Executor{}
}

func TestCalledAddDockerRegistry_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := `{"ip":"127.0.0.1"}`
	expectedRes := map[string]interface{}{
		"id": "000000000000000000000001",
	}

	registryDbExecutorMockObj := dbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		registryDbExecutorMockObj.EXPECT().AddDockerRegistry(ip).Return(registryModel, nil),
	)

	// pass mockObj to a real object.
	registryDbExecutor = registryDbExecutorMockObj

	code, res, err := manager.AddDockerRegistry(body)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(res, expectedRes) {
		t.Error()
	}
}

func TestCalledAddDockerRegistryWithInValidJsonFormatBody_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidBody := `{"ip"}`

	code, _, err := manager.AddDockerRegistry(invalidBody)

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

func TestCalledAddDockerRegistryWithInvalidBodyNotIncludingIPField_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidBody := `{"key":"value"}`

	code, _, err := manager.AddDockerRegistry(invalidBody)

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

func TestCalledAddDockerRegistryWhenFailedToInsertNewRegistryToDB_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := `{"ip":"127.0.0.1"}`

	registryDbExecutorMockObj := dbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		registryDbExecutorMockObj.EXPECT().AddDockerRegistry(ip).Return(nil, notFoundError),
	)

	// pass mockObj to a real object.
	registryDbExecutor = registryDbExecutorMockObj

	code, _, err := manager.AddDockerRegistry(body)

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

func TestCalledDeleteDockerRegistry_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registryDbExecutorMockObj := dbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		registryDbExecutorMockObj.EXPECT().DeleteDockerRegistry(registryId).Return(nil),
	)

	// pass mockObj to a real object.
	registryDbExecutor = registryDbExecutorMockObj

	code, err := manager.DeleteDockerRegistry(registryId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledDeleteDockerRegistryWithInvalidId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registryDbExecutorMockObj := dbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		registryDbExecutorMockObj.EXPECT().DeleteDockerRegistry(invalidregistryId).Return(notFoundError),
	)

	// pass mockObj to a real object.
	registryDbExecutor = registryDbExecutorMockObj

	code, err := manager.DeleteDockerRegistry(invalidregistryId)

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

func TestCalledGetDockerRegistries_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registries := []map[string]interface{}{registryModel}

	registryDbExecutorMockObj := dbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		registryDbExecutorMockObj.EXPECT().GetDockerRegistries().Return(registries, nil),
	)

	// pass mockObj to a real object.
	registryDbExecutor = registryDbExecutorMockObj

	code, res, err := manager.GetDockerRegistries()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(res["registries"].([]map[string]interface{}), registries) {
		t.Error()
	}
}

func TestCalledGetDockerRegistriesWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registryDbExecutorMockObj := dbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		registryDbExecutorMockObj.EXPECT().GetDockerRegistries().Return(nil, notFoundError),
	)

	// pass mockObj to a real object.
	registryDbExecutor = registryDbExecutorMockObj

	code, _, err := manager.GetDockerRegistries()

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

func TestCalledDockerRegistryEventHandlerWithValidRegistryEvent_ExpectSendRequestToMatchedNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appmanagementExecutorMockObj := appmocks.NewMockCommand(ctrl)
	nodemanagementExecutorMockObj := nodemocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	apps := []map[string]interface{}{app}
	matchedApplist := make(map[string]interface{})
	matchedApplist["apps"] = apps
	nodes := []map[string]interface{}{node}
	matchedNodelist := make(map[string]interface{})
	matchedNodelist["nodes"] = nodes
	dummy_imagename := dummy_host + seperator + dummy_repositry
	dummy_urls := make([]string, 0)
	dummy_url := "http://" + ip + ":48098/api/v1" + url.Management() + url.Apps() + seperator + appId + url.Events()
	dummy_urls = append(dummy_urls, dummy_url)

	gomock.InOrder(
		appmanagementExecutorMockObj.EXPECT().GetAppsWithImageName(dummy_imagename).Return(results.OK, matchedApplist, nil),
		nodemanagementExecutorMockObj.EXPECT().GetNodesWithAppID(appId).Return(results.OK, matchedNodelist, nil),
		msgMockObj.EXPECT().SendHttpRequest(POST, dummy_urls, nil, []byte(dummyRegistryEvent)),
	)

	// pass mockObj to a real object.
	appmanagementExecutor = appmanagementExecutorMockObj
	nodemanagementExecutor = nodemanagementExecutorMockObj
	httpExecutor = msgMockObj

	code, err := manager.DockerRegistryEventHandler(dummyRegistryEvent)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}
