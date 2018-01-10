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
package node

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Test functions for Agent API Handler.

type handleFunc struct {
	functionCall string
}

func TestHandle(t *testing.T) {
	w := httptest.NewRecorder()
	mockApis := handleFunc{}
	defaultApis := SdamAgent
	SdamAgent = &mockApis
	Input := [][]string{
		{GET, "/api/v1/agents", "agents"},
		{GET, "/api/v1/agents/agentID", "agent"},
		{POST, "/api/v1/agents/agentID/deploy", "agentDeployApp"},
		{GET, "/api/v1/agents/agentID/apps", "agentInfoApps"},
		{GET, "/api/v1/agents/agentID/apps/appID", "agentInfoApp"},
		{POST, "/api/v1/agents/agentID/apps/appID", "agentUpdateAppInfo"},
		{DELETE, "/api/v1/agents/agentID/apps/appID", "agentDeleteApp"},
		{POST, "/api/v1/agents/agentID/apps/appID/start", "agentStartApp"},
		{POST, "/api/v1/agents/agentID/apps/appID/stop", "agentStopApp"},
		{POST, "/api/v1/agents/agentID/apps/appID/update", "agentUpdateApp"},
		{POST, "/api/v1/agents/register", "agentRegister"},
		{POST, "/api/v1/agents/agentID/unregister", "agentUnregister"},
		{POST, "/api/v1/agents/agentID/ping", "agentPing"},
	}
	for _, val := range Input {
		method, url, funcname := val[0], val[1], val[2]
		req, _ := http.NewRequest(method, url, nil)
		SdamAgentHandle.Handle(w, req)
		if mockApis.functionCall != funcname {
			t.Error("[SDAM][Agent]Handle is invalid about " + funcname)
		}
	}
	SdamAgent = defaultApis
}

func TestHandle_Invalid_Method(t *testing.T) {
	w := httptest.NewRecorder()
	Input := map[string][]string{
		"/api/v1/agents":                           {POST, DELETE, PUT},
		"/api/v1/agents/agentID":                   {POST, DELETE, PUT},
		"/api/v1/agents/agentID/deploy":            {GET, DELETE, PUT},
		"/api/v1/agents/agentID/apps":              {PUT},
		"/api/v1/agents/agentID/apps/appID":        {PUT},
		"/api/v1/agents/agentID/apps/appID/start":  {GET, DELETE, PUT},
		"/api/v1/agents/agentID/apps/appID/stop":   {GET, DELETE, PUT},
		"/api/v1/agents/agentID/apps/appID/update": {GET, DELETE, PUT},
		"/api/v1/agents/register":                  {GET, DELETE, PUT},
		"/api/v1/agents/agentID/unregister":        {GET, DELETE, PUT},
		"/api/v1/agents/agentID/ping":              {GET, DELETE, PUT},
	}
	for key, vals := range Input {
		for _, val := range vals {
			req, _ := http.NewRequest(val, key, nil)
			SdamAgentHandle.Handle(w, req)
			if w.Code != http.StatusBadRequest {
				t.Error("[SDAM][Agent]Handle is invalid")
			}
		}
	}
}

//Mock functions for Agent APIs.

func (mockApis *handleFunc) agentRegister(w http.ResponseWriter, req *http.Request) {
	mockApis.functionCall = "agentRegister"
}

func (mockApis *handleFunc) agentPing(w http.ResponseWriter, req *http.Request, agentID string) {
	mockApis.functionCall = "agentPing"
}

func (mockApis *handleFunc) agentUnregister(w http.ResponseWriter, req *http.Request, agentID string) {
	mockApis.functionCall = "agentUnregister"
}

func (mockApis *handleFunc) agent(w http.ResponseWriter, req *http.Request, agentID string) {
	mockApis.functionCall = "agent"
}

func (mockApis *handleFunc) agents(w http.ResponseWriter, req *http.Request) {
	mockApis.functionCall = "agents"
}

func (mockApis *handleFunc) agentDeployApp(w http.ResponseWriter, req *http.Request, agentID string) {
	mockApis.functionCall = "agentDeployApp"
}

func (mockApis *handleFunc) agentInfoApps(w http.ResponseWriter, req *http.Request, agentID string) {
	mockApis.functionCall = "agentInfoApps"
}

func (mockApis *handleFunc) agentInfoApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	mockApis.functionCall = "agentInfoApp"
}

func (mockApis *handleFunc) agentUpdateAppInfo(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	mockApis.functionCall = "agentUpdateAppInfo"
}

func (mockApis *handleFunc) agentDeleteApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	mockApis.functionCall = "agentDeleteApp"
}

func (mockApis *handleFunc) agentStartApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	mockApis.functionCall = "agentStartApp"
}

func (mockApis *handleFunc) agentStopApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	mockApis.functionCall = "agentStopApp"
}

func (mockApis *handleFunc) agentUpdateApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	mockApis.functionCall = "agentUpdateApp"
}

//Test functions for Agent APIs.

type controllerFunc struct {
	functionCall  string
	occurredError bool
}

func newCtrlFunc() *controllerFunc {
	cf := controllerFunc{}
	cf.functionCall = ""
	cf.occurredError = false
	return &cf
}

var testBody = map[string]interface{}{
	"id":   "testID",
	"host": "0.0.0.0",
	"port": "8080",
	"apps": nil,
}

func TestAgent(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/agents/testAgentID", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agent(w, req, "testAgentID")
	if mockCtrl.functionCall != "GetAgent" || w.Code != http.StatusOK {
		t.Error("[SDAM][Agent]agent is invalid")
	}
}

func TestAgent_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/agents/testAgentID", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agent(w, req, "testAgentID")
	if mockCtrl.functionCall != "GetAgent" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Agent]agent is invalid about controller occurred error")
	}
}

func TestAgents(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/agents", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agents(w, req)
	if mockCtrl.functionCall != "GetAgents" || w.Code != http.StatusOK {
		t.Error("[SDAM][Agent]agents is invalid")
	}
}

func TestAgents_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/agents", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agents(w, req)
	if mockCtrl.functionCall != "GetAgents" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Agent]agents is invalid about controller occurred error")
	}
}

func TestAgentDeployApp(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	bod, _ := json.Marshal(testBody)
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/deploy", bytes.NewReader(bod))
	sdamAgentController = mockCtrl
	SdamAgent.agentDeployApp(w, req, "testAgentID")
	if mockCtrl.functionCall != "DeployApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Agent]agentDeployApp is invalid")
	}
}

func TestAgentDeployApp_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	bod := []byte("error")
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/deploy", bytes.NewReader(bod))
	sdamAgentController = mockCtrl
	SdamAgent.agentDeployApp(w, req, "testAgentID")
	if mockCtrl.functionCall != "DeployApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Agent]agentDeployApp is invalid about controller occurred error")
	}
}

func TestAgentDeployApp_empty_body(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/deploy", nil)
	SdamAgent.agentDeployApp(w, req, "testAgentID")
	if w.Code != http.StatusBadRequest {
		t.Error("[SDAM][Agent]agentDeployApp is invalid about emtpy body")
	}
}

func TestAgentInfoApps(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/agents/testAgentID/apps", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agentInfoApps(w, req, "testAgentID")
	if mockCtrl.functionCall != "GetApps" || w.Code != http.StatusOK {
		t.Error("[SDAM][Agent]agentInfoApps is invalid")
	}
}

func TestAgentInfoApps_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/agents/testAgentID/apps", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agentInfoApps(w, req, "testAgentID")
	if mockCtrl.functionCall != "GetApps" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Agent]agentInfoApps is invalid about controller occurred error")
	}
}

func TestAgentInfoApp(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/agents/testAgentID/apps/testAppID", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agentInfoApp(w, req, "testAgentID", "testAppID")
	if mockCtrl.functionCall != "GetApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Agent]agentInfoApp is invalid")
	}
}

func TestAgentInfoApp_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(GET, "/api/v1/agents/testAgentID/apps/testAppID", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agentInfoApp(w, req, "testAgentID", "testAppID")
	if mockCtrl.functionCall != "GetApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Agent]agentInfoApp is invalid about controller occurred error")
	}
}

func TestAgentUpdateAppInfo(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	bod, _ := json.Marshal(testBody)
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/apps/testAppID", bytes.NewReader(bod))
	sdamAgentController = mockCtrl
	SdamAgent.agentUpdateAppInfo(w, req, "testAgentID", "testAppID")
	if mockCtrl.functionCall != "UpdateAppInfo" || w.Code != http.StatusOK {
		t.Error("[SDAM][Agent]agentUpdateAppInfo is invalid")
	}
}

func TestAgentUpdateAppInfo_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	bod := []byte("error")
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/apps/testAppID", bytes.NewReader(bod))
	sdamAgentController = mockCtrl
	SdamAgent.agentUpdateAppInfo(w, req, "testAgentID", "testAppID")
	if mockCtrl.functionCall != "UpdateAppInfo" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Agent]agentUpdateAppInfo is invalid about controller occurred error")
	}
}

func TestAgentUpdateAppInfo_empty_body(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/apps/testAppID", nil)
	SdamAgent.agentUpdateAppInfo(w, req, "testAgentID", "testAppID")
	if w.Code != http.StatusBadRequest {
		t.Error("[SDAM][Agent]agentUpdateAppInfo is invalid about empty body")
	}
}

func TestDeleteApp(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(DELETE, "/api/v1/agents/testAgentID/apps/testAppID", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agentDeleteApp(w, req, "testAgentID", "testAppID")
	if mockCtrl.functionCall != "DeleteApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Agent]agentDeleteApps is invalid")
	}
}

func TestDeleteApp_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(DELETE, "/api/v1/agents/testAgentID/apps/testAppID", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agentDeleteApp(w, req, "testAgentID", "testAppID")
	if mockCtrl.functionCall != "DeleteApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Agent]agentDeleteApps is invalid about controller occurred error")
	}
}

func TestAgentStartApp(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/apps/testAppID/start", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agentStartApp(w, req, "testAgentID", "testAppID")
	if mockCtrl.functionCall != "StartApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Agent]agentStartApps is invalid")
	}
}

func TestAgentStartApp_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/apps/testAppID/start", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agentStartApp(w, req, "testAgentID", "testAppID")
	if mockCtrl.functionCall != "StartApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Agent]agentStartApps is invalid about controller occurred error")
	}
}

func TestAgentStopApp(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/apps/testAppID/stop", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agentStopApp(w, req, "testAgentID", "testAppID")
	if mockCtrl.functionCall != "StopApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Agent]agentStopApps is invalid")
	}
}

func TestAgentStopApp_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/apps/testAppID/stop", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agentStopApp(w, req, "testAgentID", "testAppID")
	if mockCtrl.functionCall != "StopApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Agent]agentStopApps is invalid about controller occurred error")
	}
}

func TestAgentUpdateApp(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/apps/testAppID/update", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agentUpdateApp(w, req, "testAgentID", "testAppID")
	if mockCtrl.functionCall != "UpdateApp" || w.Code != http.StatusOK {
		t.Error("[SDAM][Agent]agentUpdateApps is invalid")
	}
}

func TestAgentUpdateApp_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/apps/testAppID/update", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agentUpdateApp(w, req, "testAgentID", "testAppID")
	if mockCtrl.functionCall != "UpdateApp" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Agent]agentUpdateApps is invalid about controller occurred error")
	}
}

func TestAgentRegister(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	bod, _ := json.Marshal(testBody)
	req, _ := http.NewRequest(POST, "/api/v1/agents/register", bytes.NewReader(bod))
	sdamAgentController = mockCtrl
	SdamAgent.agentRegister(w, req)
	if mockCtrl.functionCall != "AddAgent" || w.Code != http.StatusOK {
		t.Error("[SDAM][Agent]agentRegister is invalid")
	}
}

func TestAgentRegister_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	bod := []byte("error")
	req, _ := http.NewRequest(POST, "/api/v1/agents/register", bytes.NewReader(bod))
	sdamAgentController = mockCtrl
	SdamAgent.agentRegister(w, req)
	if mockCtrl.functionCall != "AddAgent" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Agent]agentRegister is invalid about controller occurred error")
	}
}

func TestAgentRegister_empty_body(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/agents/register", nil)
	SdamAgent.agentRegister(w, req)
	if w.Code != http.StatusBadRequest {
		t.Error("[SDAM][Agent]agentRegister is invalid about empty body")
	}
}

func TestAgentUnregister(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/unregister", nil)
	sdamAgentController = mockCtrl
	SdamAgent.agentUnregister(w, req, "testAgentID")
	if mockCtrl.functionCall != "DeleteAgent" || w.Code != http.StatusOK {
		t.Error("[SDAM][Agent]agentUnregister is invalid")
	}
}

func TestAgentUnregister_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	bod := []byte("error")
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/unregister", bytes.NewReader(bod))
	sdamAgentController = mockCtrl
	SdamAgent.agentUnregister(w, req, "testAgentID")
	if mockCtrl.functionCall != "DeleteAgent" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Agent]agentUnregister is invalid about controller occurred error")
	}
}

func TestAgentPing(t *testing.T) {
	mockCtrl := newCtrlFunc()
	w := httptest.NewRecorder()
	bod, _ := json.Marshal(testBody)
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/ping", bytes.NewReader(bod))
	sdamAgentController = mockCtrl
	SdamAgent.agentPing(w, req, "testAgentID")
	if mockCtrl.functionCall != "PingAgent" || w.Code != http.StatusOK {
		t.Error("[SDAM][Agent]agentPing is invalid")
	}
}

func TestAgentPing_controller_occurred_error(t *testing.T) {
	mockCtrl := newCtrlFunc()
	mockCtrl.occurredError = true
	w := httptest.NewRecorder()
	bod := []byte("error")
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/ping", bytes.NewReader(bod))
	sdamAgentController = mockCtrl
	SdamAgent.agentPing(w, req, "testAgentID")
	if mockCtrl.functionCall != "PingAgent" || w.Code != http.StatusNotFound {
		t.Error("[SDAM][Agent]agentPing is invalid about controller occurred error")
	}
}

func TestAgentPing_empty_body(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(POST, "/api/v1/agents/testAgentID/ping", nil)
	SdamAgent.agentPing(w, req, "testAgentID")
	if w.Code != http.StatusBadRequest {
		t.Error("[SDAM][Agent]agentPing is invalid about empty body")
	}
}

//Mock functions for Agent Controller Functions.

func (mockCtrl *controllerFunc) AddAgent(ip string, body string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "AddAgent"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) PingAgent(agentId string, ip string, body string) (int, error) {
	mockCtrl.functionCall = "PingAgent"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil
	}
	return http.StatusNotFound, nil
}

func (mockCtrl *controllerFunc) DeleteAgent(agentId string) (int, error) {
	mockCtrl.functionCall = "DeleteAgent"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil
	}
	return http.StatusNotFound, nil
}

func (mockCtrl *controllerFunc) GetAgent(agentID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "GetAgent"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) GetAgents() (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "GetAgents"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) DeployApp(agentID string, body string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "DeployApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) GetApps(agentID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "GetApps"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) GetApp(agentID string, appID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "GetApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) UpdateAppInfo(agentID string, appID string, body string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "UpdateAppInfo"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) DeleteApp(agentID string, appID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "DeleteApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) UpdateApp(agentID string, appID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "UpdateApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) StartApp(agentID string, appID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "StartApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}

func (mockCtrl *controllerFunc) StopApp(agentID string, appID string) (int, map[string]interface{}, error) {
	mockCtrl.functionCall = "StopApp"
	if !mockCtrl.occurredError {
		return http.StatusOK, nil, nil
	}
	return http.StatusNotFound, nil, nil
}
