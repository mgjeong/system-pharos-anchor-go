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

package api

import (
	managementmocks "api/management/mocks"
	monitoringmocks "api/monitoring/mocks"
	searchmocks "api/search/mocks"
	healthmocks "api/health/mocks"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalledServeHTTPWithInvalidURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	managementHandlerMockObj := managementmocks.NewMockCommand(ctrl)
	monitoringHandlerMockObj := monitoringmocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/invalid", nil)

	// pass mockObj to a real object.
	managementHandler = managementHandlerMockObj
	monitoringHandler = monitoringHandlerMockObj

	Handler.ServeHTTP(w, req)
}

func TestCalledServeHTTPWithExcludedBaseURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	managementHandlerMockObj := managementmocks.NewMockCommand(ctrl)
	monitoringHandlerMockObj := monitoringmocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/monitoring/resource", nil)

	// pass mockObj to a real object.
	managementHandler = managementHandlerMockObj
	monitoringHandler = monitoringHandlerMockObj

	Handler.ServeHTTP(w, req)
}

func TestCalledServeHTTPWithManagementRequest_ExpectCalledManagementHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	managementHandlerMockObj := managementmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		managementHandlerMockObj.EXPECT().Handle(gomock.Any(), gomock.Any()),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/management/nodes", nil)

	// pass mockObj to a real object.
	managementHandler = managementHandlerMockObj

	Handler.ServeHTTP(w, req)
}

func TestCalledServeHTTPWithMonitoringRequest_ExpectCalledMonitoringHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	monitoringHandlerMockObj := monitoringmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		monitoringHandlerMockObj.EXPECT().Handle(gomock.Any(), gomock.Any()),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/monitoring/resource", nil)

	// pass mockObj to a real object.
	monitoringHandler = monitoringHandlerMockObj

	Handler.ServeHTTP(w, req)
}

func TestCalledServeHTTPWithSearchRequest_ExpectCalledSearchHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	searchHandlerMockObj := searchmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		searchHandlerMockObj.EXPECT().Handle(gomock.Any(), gomock.Any()),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search", nil)

	// pass mockObj to a real object.
	searchHandler = searchHandlerMockObj

	Handler.ServeHTTP(w, req)
}

func TestCalledServeHTTPWithPingRequest_ExpectCalledPingHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	healthHandlerMockObj := healthmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		healthHandlerMockObj.EXPECT().Handle(gomock.Any(), gomock.Any()),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ping", nil)

	// pass mockObj to a real object.
	healthHandler = healthHandlerMockObj

	Handler.ServeHTTP(w, req)
}
