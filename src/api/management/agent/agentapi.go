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
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	agentmanager "controller/management/agent"
	"net/http"
	"strings"
)

const (
	GET string = "GET"
)

type Command interface {
	agent(w http.ResponseWriter, req *http.Request, agentID string)
	agents(w http.ResponseWriter, req *http.Request)
	register(w http.ResponseWriter, req *http.Request)
	ping(w http.ResponseWriter, req *http.Request, agentID string)
	unregister(w http.ResponseWriter, req *http.Request, agentID string)
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
			agentAPI.register(w, req)
		} else {
			if req.Method == GET {
				agentID := split[1]
				agentAPI.agent(w, req, agentID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		}

	case 3:
		if strings.Contains(url, URL.Apps()) {
			apps.Handler.Handle(w, req)
		} else if "/"+split[2] == URL.Unregister() {
			agentID := split[1]
			agentAPI.unregister(w, req, agentID)
		} else if "/"+split[2] == URL.Ping() {
			agentID := split[1]
			agentAPI.ping(w, req, agentID)
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 4:
	case 5:
		if strings.Contains(url, URL.Apps()) {
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
//    paths: '/api/v1/management/agents'
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

// register handles requests which is used to register agent to a list of agents.
//
//    paths: '/api/v1/management/agents/register'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (agentAPIExecutor) register(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[AGENT] Register New Service Deployment Agent")

	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := managementExecutor.RegisterAgent(body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// unregister handles requests which is used to unregister agent from a list of agents.
//
//    paths: '/api/v1/management/agents/{agentID}/unregister'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (agentAPIExecutor) unregister(w http.ResponseWriter, req *http.Request, agentID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Unregister New Service Deployment Agent")

	result, err := managementExecutor.UnRegisterAgent(agentID)
	common.MakeResponse(w, result, nil, err)
}

// ping handles requests which is used to check whether an agent is up.
//
//    paths: '/api/v1/management/agents/{agentID}/ping'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (agentAPIExecutor) ping(w http.ResponseWriter, req *http.Request, agentID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Ping From Service Deployment Agent")

	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, err := managementExecutor.PingAgent(agentID, body)
	common.MakeResponse(w, result, nil, err)
}
