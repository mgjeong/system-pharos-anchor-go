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
	"bytes"
	"encoding/json"
	registrymanagermocks "controller/management/registry/mocks"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testBodyString  = `{"test":"body"}`
)

var testBody = map[string]interface{}{
	"test":   "body",
}

var Handler Command

func init() {
	Handler = RequestHandler{}
}

func TestCalledHandleWithInvalidURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registrymanageMockObj := registrymanagermocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/invalid", nil)

	// pass mockObj to a real object.
	registryExecutor = registrymanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithExcludedBaseURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registrymanageMockObj := registrymanagermocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/management/registries", nil)

	// pass mockObj to a real object.
	registryExecutor = registrymanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetRegistriesRequest_ExpectCalledGetRegistries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registrymanageMockObj := registrymanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		registrymanageMockObj.EXPECT().GetDockerRegistries(),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/registries", nil)

	// pass mockObj to a real object.
	registryExecutor = registrymanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetRegistryRequest_ExpectCalledGetRegistry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registrymanageMockObj := registrymanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		registrymanageMockObj.EXPECT().GetDockerRegistry("registryID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/registries/registryID", nil)

	// pass mockObj to a real object.
	registryExecutor = registrymanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithAddRegistryRequest_ExpectCalledAddDockerRegistry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registrymanageMockObj := registrymanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		registrymanageMockObj.EXPECT().AddDockerRegistry(testBodyString),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/management/registries", bytes.NewReader(body))

	// pass mockObj to a real object.
	registryExecutor = registrymanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithDeleteRegistryRequest_ExpectCalledDeleteDockerRegistry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registrymanageMockObj := registrymanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		registrymanageMockObj.EXPECT().DeleteDockerRegistry("registryID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/management/registries/registryID", nil)

	// pass mockObj to a real object.
	registryExecutor = registrymanageMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetImagesRequest_ExpectCalledGetDockerImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registrymanageMockObj := registrymanagermocks.NewMockCommand(ctrl)

	gomock.InOrder(
		registrymanageMockObj.EXPECT().DockerRegistryEventHandler(testBodyString),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/management/registries/events", bytes.NewReader(body))

	// pass mockObj to a real object.
	registryExecutor = registrymanageMockObj

	Handler.Handle(w, req)
}

