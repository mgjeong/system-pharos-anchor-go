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

// Package api/agent/health provides functionality to handle request related to health.
package health

import (
	"api/common"
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	"controller/health"
	"net/http"
	"strings"
)

const (
	POST   string = "POST"
)

type Command interface {
	register(w http.ResponseWriter, req *http.Request)
	ping(w http.ResponseWriter, req *http.Request, agentID string)
	unregister(w http.ResponseWriter, req *http.Request, agentID string)
}

type healthHandler struct{}
type healthAPIExecutor struct {
	Command
}

var healthExecutor health.Command
var healthAPI healthAPIExecutor
var Handler healthHandler

func init() {
	healthExecutor = health.Executor{}
	healthAPI = healthAPIExecutor{}
	Handler = healthHandler{}
}

// Handle calls a proper function according to the url and method received from remote device.
func (healthHandler) Handle(w http.ResponseWriter, req *http.Request) {
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Management()+URL.Agents(), "", -1)
	split := strings.Split(url, "/")
	switch len(split) {
	case 2:
		if "/"+split[1] == URL.Register() {
			if req.Method == POST {
				healthAPI.register(w, req)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		}
	case 3:
		agentID := split[1]
		if "/"+split[2] == URL.Unregister() {
			if req.Method == POST {
				healthAPI.unregister(w, req, agentID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else if "/"+split[2] == URL.Ping() {
			if req.Method == POST {
				healthAPI.ping(w, req, agentID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}
	}
}

// register handles requests which is used to register agent to a list of agents.
//
//    paths: '/api/v1/management/agents/register'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (healthAPIExecutor) register(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[AGENT] Register New Service Deployment Agent")

	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := healthExecutor.RegisterAgent(body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// unregister handles requests which is used to unregister agent from a list of agents.
//
//    paths: '/api/v1/management/agents/{agentID}/unregister'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (healthAPIExecutor) unregister(w http.ResponseWriter, req *http.Request, agentID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Unregister New Service Deployment Agent")

	result, err := healthExecutor.UnRegisterAgent(agentID)
	common.MakeResponse(w, result, nil, err)
}

// ping handles requests which is used to check whether an agent is up.
//
//    paths: '/api/v1/management/agents/{agentID}/ping'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (healthAPIExecutor) ping(w http.ResponseWriter, req *http.Request, agentID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Ping From Service Deployment Agent")

	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, err := healthExecutor.PingAgent(agentID, body)
	common.MakeResponse(w, result, nil, err)
}
