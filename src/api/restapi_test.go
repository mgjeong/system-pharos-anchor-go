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
	"api/agent"
	"api/group"
	"net/http"
	"net/http/httptest"
	"testing"
)

type agentMock struct {
	handlerCall bool
}
type groupMock struct {
	handlerCall bool
}

var am agentMock
var gm groupMock

func setUp() func() {
	am.handlerCall = false
	gm.handlerCall = false
	defaultSdamAgentHandle := agent.SdamAgentHandle
	defaultSdamGroupHandle := group.SdamGroupHandle
	agent.SdamAgentHandle = &am
	group.SdamGroupHandle = &gm
	return func() {
		agent.SdamAgentHandle = defaultSdamAgentHandle
		group.SdamGroupHandle = defaultSdamGroupHandle
	}
}

func TestServeHTTPsendAgent(t *testing.T) {
	tearDown := setUp()
	defer tearDown()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/agents", nil)
	_SDAMApis.ServeHTTP(w, req)

	if !am.handlerCall || gm.handlerCall {
		t.Error("ServeHTTPsendAgent is invalid")
	}
}

func TestServeHTTPsendGroup(t *testing.T) {
	tearDown := setUp()
	defer tearDown()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/groups", nil)
	_SDAMApis.ServeHTTP(w, req)

	if !gm.handlerCall || am.handlerCall {
		t.Error("ServeHTTPsendGroup is invalid")
	}
}

func TestServeHTTPURLisEmpty(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "", nil)
	_SDAMApis.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Error("ServeHTTPURLisEmpty is invalid")
	}
}

func TestServeHTTPinvalidURL(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/", nil)
	_SDAMApis.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Error("ServeHTTPinvalidURL is invalid")
	}
}

func (am *agentMock) Handle(w http.ResponseWriter, req *http.Request) {
	am.handlerCall = true
}
func (gm *groupMock) Handle(w http.ResponseWriter, req *http.Request) {
	gm.handlerCall = true
}
