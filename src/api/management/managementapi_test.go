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

package management

import (
	nodemocks "api/management/node/mocks"
	groupmocks "api/management/group/mocks"
	registrymocks "api/management/registry/mocks"
	"net/http"
	"net/http/httptest"
	"github.com/golang/mock/gomock"
	"testing"
)

var Handler Command

func init() {
	Handler = RequestHandler{}
}

func TestCalledHandleWithInvalidURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodeHandlerMockObj := nodemocks.NewMockCommand(ctrl)
	groupHandlerMockObj := groupmocks.NewMockCommand(ctrl)
	registryHandlerMockObj := registrymocks.NewMockCommand(ctrl)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/invalid", nil)
	
	// pass mockObj to a real object.
	nodeManagementHandler = nodeHandlerMockObj
	groupManagementHandler = groupHandlerMockObj
	registryManagementHandler = registryHandlerMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithExcludedBaseURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodeHandlerMockObj := nodemocks.NewMockCommand(ctrl)
	groupHandlerMockObj := groupmocks.NewMockCommand(ctrl)
	registryHandlerMockObj := registrymocks.NewMockCommand(ctrl)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nodes/resource", nil)
	
	// pass mockObj to a real object.
	nodeManagementHandler = nodeHandlerMockObj
	groupManagementHandler = groupHandlerMockObj
	registryManagementHandler = registryHandlerMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithNodeRequest_ExpectCalledNodeHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodeHandlerMockObj := nodemocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodeHandlerMockObj.EXPECT().Handle(gomock.Any(), gomock.Any()),
	)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/nodes", nil)
	
	// pass mockObj to a real object.
	nodeManagementHandler = nodeHandlerMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGroupRequest_ExpectCalledGroupHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupHandlerMockObj := groupmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupHandlerMockObj.EXPECT().Handle(gomock.Any(), gomock.Any()),
	)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/groups", nil)
	
	// pass mockObj to a real object.
	groupManagementHandler = groupHandlerMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithRegistryRequest_ExpectCalledRegistryHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registryHandlerMockObj := registrymocks.NewMockCommand(ctrl)

	gomock.InOrder(
		registryHandlerMockObj.EXPECT().Handle(gomock.Any(), gomock.Any()),
	)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/registries", nil)
	
	// pass mockObj to a real object.
	registryManagementHandler = registryHandlerMockObj

	Handler.Handle(w, req)
}
