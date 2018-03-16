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

package monitoring

import (
	resourcemocks "api/monitoring/resource/mocks"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

var Handler Command

func init() {
	Handler = RequestHandler{}
}

func TestCalledHandleWithInvalidURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceHandlerMockObj := resourcemocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/invalid", nil)

	// pass mockObj to a real object.
	resourceMonitoringHandler = resourceHandlerMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithExcludedBaseURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceHandlerMockObj := resourcemocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nodes/nodeId/resource", nil)

	// pass mockObj to a real object.
	resourceMonitoringHandler = resourceHandlerMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetNodeResourceNodeRequest_ExpectCalledNodeHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceHandlerMockObj := resourcemocks.NewMockCommand(ctrl)

	gomock.InOrder(
		resourceHandlerMockObj.EXPECT().Handle(gomock.Any(), gomock.Any()),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/monitoring/nodes/nodeId/resource", nil)

	// pass mockObj to a real object.
	resourceMonitoringHandler = resourceHandlerMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetAppResourceNodeRequest_ExpectCalledNodeHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceHandlerMockObj := resourcemocks.NewMockCommand(ctrl)

	gomock.InOrder(
		resourceHandlerMockObj.EXPECT().Handle(gomock.Any(), gomock.Any()),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/monitoring/nodes/nodeId/apps/appId/resource", nil)

	// pass mockObj to a real object.
	resourceMonitoringHandler = resourceHandlerMockObj

	Handler.Handle(w, req)
}