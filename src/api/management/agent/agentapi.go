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

// Package api/agent provides functionality to handle request related to agent.
package agent

import (
	"api/common"
	"api/management/agent/apps"
	"api/management/agent/health"
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	agentmanager "controller/management/agent"
	"net/http"
	"strings"
)

const (
	GET    string = "GET"
)

type Command interface {
	agent(w http.ResponseWriter, req *http.Request, agentID string)
	agents(w http.ResponseWriter, req *http.Request)
}

type agentHandler struct{}
type agentAPIExecutor struct {
	Command
}

var managementExecutor agentmanager.Command
var agentAPI agentAPIExecutor
var Handler agentHandler

func init() {
	managementExecutor = agentmanager.Executor{}
	agentAPI = agentAPIExecutor{}
	Handler = agentHandler{}
}

// Handle calls a proper function according to the url and method received from remote device.
func (agentHandler) Handle(w http.ResponseWriter, req *http.Request) {
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Management()+URL.Agents(), "", -1)
	split := strings.Split(url, "/")
	switch len(split) {
	case 1:
		if req.Method == GET {
			agentAPI.agents(w, req)
		} else {
			common.WriteError(w, errors.InvalidMethod{req.Method})
		}

	case 2:
		if "/"+split[1] == URL.Register() {
			health.Handler.Handle(w, req)
		} else {
			if req.Method == GET {
				agentID := split[1]
				agentAPI.agent(w, req, agentID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		}

	case 3:
		if "/"+split[2] == URL.Deploy() ||
			"/"+split[2] == URL.Apps() {
			apps.Handler.Handle(w, req)
		} else if "/"+split[2] == URL.Unregister() ||
			"/"+split[2] == URL.Ping() {
			health.Handler.Handle(w, req)
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 4:
	case 5:
		if "/"+split[2] == URL.Apps() {
			apps.Handler.Handle(w, req)
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}
	}
}

// agents handles requests which is used to get information of agent identified by the given agentID.
//
//    paths: '/api/v1/management/agents/{agentID}'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (agentAPIExecutor) agent(w http.ResponseWriter, req *http.Request, agentID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Get Service Deployment Agent")
	result, resp, err := managementExecutor.GetAgent(agentID)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agents handles requests which is used to get information of all agents registered.
//
//    paths: '/api/v1/agents'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (agentAPIExecutor) agents(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[AGENT] Get All Service Deployment Agents")
	result, resp, err := managementExecutor.GetAgents()
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}
