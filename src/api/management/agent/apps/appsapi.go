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

// Package api/agent/apps provides functionality to handle request related to apps.
package apps

import (
	"api/common"
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	deployment "controller/deployment/agent"
	"net/http"
	"strings"
)

const (
	GET    string = "GET"
	POST   string = "POST"
	DELETE string = "DELETE"
)

type Command interface {
	agentDeployApp(w http.ResponseWriter, req *http.Request, agentID string)
	agentInfoApps(w http.ResponseWriter, req *http.Request, agentID string)
	agentInfoApp(w http.ResponseWriter, req *http.Request, agentID string, appID string)
	agentUpdateAppInfo(w http.ResponseWriter, req *http.Request, agentID string, appID string)
	agentDeleteApp(w http.ResponseWriter, req *http.Request, agentID string, appID string)
	agentStartApp(w http.ResponseWriter, req *http.Request, agentID string, appID string)
	agentStopApp(w http.ResponseWriter, req *http.Request, agentID string, appID string)
	agentUpdateApp(w http.ResponseWriter, req *http.Request, agentID string, appID string)
}

type appsHandler struct{}
type appsAPIExecutor struct {
	Command
}

var deploymentExecutor deployment.Command
var appsAPI appsAPIExecutor
var Handler appsHandler

func init() {
	deploymentExecutor = deployment.Executor{}
	appsAPI = appsAPIExecutor{}
	Handler = appsHandler{}
}

// Handle calls a proper function according to the url and method received from remote device.
func (appsHandler) Handle(w http.ResponseWriter, req *http.Request) {
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Management()+URL.Agents(), "", -1)
	split := strings.Split(url, "/")
	switch len(split) {
	case 3:
		agentID := split[1]
		if "/"+split[2] == URL.Deploy() {
			if req.Method == POST {
				appsAPI.agentDeployApp(w, req, agentID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else if "/"+split[2] == URL.Apps() {
			if req.Method == GET {
				appsAPI.agentInfoApps(w, req, agentID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 4:
		if "/"+split[2] == URL.Apps() {
			agentID, appID := split[1], split[3]
			switch req.Method {
			case GET:
				appsAPI.agentInfoApp(w, req, agentID, appID)

			case POST:
				appsAPI.agentUpdateAppInfo(w, req, agentID, appID)

			case DELETE:
				appsAPI.agentDeleteApp(w, req, agentID, appID)

			default:
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 5:
		if "/"+split[2] == URL.Apps() {
			agentID, appID := split[1], split[3]
			switch {
			case "/"+split[4] == URL.Start() && req.Method == POST:
				appsAPI.agentStartApp(w, req, agentID, appID)

			case "/"+split[4] == URL.Stop() && req.Method == POST:
				appsAPI.agentStopApp(w, req, agentID, appID)

			case "/"+split[4] == URL.Update() && req.Method == POST:
				appsAPI.agentUpdateApp(w, req, agentID, appID)

			default:
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}
	}
}

// agentDeployApp handles requests which is used to deploy new application to agent
// identified by the given agentID.
//
//    paths: '/api/v1/management/agents/{agentID}/apps/deploy'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) agentDeployApp(w http.ResponseWriter, req *http.Request, agentID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Deploy App")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := deploymentExecutor.DeployApp(agentID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentInfoApps handles requests which is used to get information of all applications
// installed on agent identified by the given agentID.
//
//    paths: '/api/v1/management/agents/{agentID}/apps'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) agentInfoApps(w http.ResponseWriter, req *http.Request, agentID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Get Info Apps")
	result, resp, err := deploymentExecutor.GetApps(agentID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentInfoApp handles requests which is used to get information of application
// identified by the given appID.
//
//    paths: '/api/v1/management/agents/{agentID}/apps/{appID}'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) agentInfoApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Get Info App")
	result, resp, err := deploymentExecutor.GetApp(agentID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentUpdateAppInfo handles requests related to updating the application with given yaml in body.
//
//    paths: '/api/v1/management/agents/{agentID}/apps/{appID}'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) agentUpdateAppInfo(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Update App Info")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := deploymentExecutor.UpdateAppInfo(agentID, appID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentDeleteApp handles requests related to delete application installed on agent
// identified by the given agentID.
//
//    paths: '/api/v1/management/agents/{agentID}/apps/{appID}'
//    method: DELETE
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) agentDeleteApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Delete App")
	result, resp, err := deploymentExecutor.DeleteApp(agentID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentStartApp handles requests related to start application installed on agent
// identified by the given agentID.
//
//    paths: '/api/v1/management/agents/{agentID}/apps/{appID}/start'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) agentStartApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Start App")
	result, resp, err := deploymentExecutor.StartApp(agentID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentStopApp handles requests related to stop application installed on agent
// identified by the given agentID.
//
//    paths: '/api/v1/management/agents/{agentID}/apps/{appID}/stop'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) agentStopApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Stop App")
	result, resp, err := deploymentExecutor.StopApp(agentID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentUpdateApp handles requests related to updating application installed on agent
// identified by the given agentID.
//
//    paths: '/api/v1/management/agents/{agentID}/apps/{appID}/update'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) agentUpdateApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Update App")
	result, resp, err := deploymentExecutor.UpdateApp(agentID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}
