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

package search

import (
	appsearchmocks "api/search/app/mocks"
	groupsearchmocks "api/search/group/mocks"
	nodessearchmocks "api/search/node/mocks"
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

	appsSearchHandlerMockObj := appsearchmocks.NewMockCommand(ctrl)
	nodesSearchHandlerMockObj := nodessearchmocks.NewMockCommand(ctrl)
	groupsSearchHandlerMockObj := groupsearchmocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/invalid", nil)

	// pass mockObj to a real object.
	appSearchHandler = appsSearchHandlerMockObj
	nodeSearchHandler = nodesSearchHandlerMockObj
	groupSearchHandler = groupsSearchHandlerMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithExcludedBaseURL_UnExpectCalledAnyHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appsSearchHandlerMockObj := appsearchmocks.NewMockCommand(ctrl)
	nodesSearchHandlerMockObj := nodessearchmocks.NewMockCommand(ctrl)
	groupsSearchHandlerMockObj := groupsearchmocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search", nil)

	// pass mockObj to a real object.
	appSearchHandler = appsSearchHandlerMockObj
	nodeSearchHandler = nodesSearchHandlerMockObj
	groupSearchHandler = groupsSearchHandlerMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithAppsSearchRequest_ExpectCalledAppsSearchHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appsSearchHandlerMockObj := appsearchmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		appsSearchHandlerMockObj.EXPECT().Handle(gomock.Any(), gomock.Any()),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search/apps", nil)

	// pass mockObj to a real object.
	appSearchHandler = appsSearchHandlerMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithNodesSearchRequest_ExpectCalledNodesSearchHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodesSearchHandlerMockObj := nodessearchmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodesSearchHandlerMockObj.EXPECT().Handle(gomock.Any(), gomock.Any()),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search/nodes", nil)

	// pass mockObj to a real object.
	nodeSearchHandler = nodesSearchHandlerMockObj

	Handler.Handle(w, req)
}

func TestCalledHandleWithGroupsSearchRequest_ExpectCalledGroupsSearchHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupsSearchHandlerMockObj := groupsearchmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupsSearchHandlerMockObj.EXPECT().Handle(gomock.Any(), gomock.Any()),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search/groups", nil)

	// pass mockObj to a real object.
	groupSearchHandler = groupsSearchHandlerMockObj

	Handler.Handle(w, req)
}
