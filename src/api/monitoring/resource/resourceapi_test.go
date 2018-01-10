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

package resource

import (
	resourcemocks "controller/resource/node/mocks"
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

	resourceMockObj := resourcemocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/monitoring/invalid", nil)

	// pass mockObj to a real object.
	resourceExecutor = resourceMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithExcludedBaseURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceMockObj := resourcemocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/monitoring/nodes", nil)

	// pass mockObj to a real object.
	resourceExecutor = resourceMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetResourceRequest_ExpectCalledGetResourceInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceMockObj := resourcemocks.NewMockCommand(ctrl)

	gomock.InOrder(
		resourceMockObj.EXPECT().GetResourceInfo("nodeID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/monitoring/nodes/nodeID/resource", nil)

	// pass mockObj to a real object.
	resourceExecutor = resourceMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGetResourcePerformanceRequest_ExpectCalledGetPerformanceInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceMockObj := resourcemocks.NewMockCommand(ctrl)

	gomock.InOrder(
		resourceMockObj.EXPECT().GetPerformanceInfo("nodeID"),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/monitoring/nodes/nodeID/resource/performance", nil)

	// pass mockObj to a real object.
	resourceExecutor = resourceMockObj

	Handler.Handle(w, req)
}