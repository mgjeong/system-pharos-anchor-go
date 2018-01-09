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
	"api/common"
	"commons/errors"
	"commons/logger"
	URL "commons/url"
	resource "controller/resource/agent"
	"net/http"
	"strings"
)

const (
	GET    string = "GET"
)

type Command interface {
	agentGetResourceInfo(w http.ResponseWriter, req *http.Request, agentId string)
	agentGetPerformanceInfo(w http.ResponseWriter, req *http.Request, agentId string)
}

type resourceHandler struct{}
type resourceAPIExecutor struct {
	Command
}

var Handler resourceHandler
var resourceAPI resourceAPIExecutor
var resourceExecutor resource.Command

func init() {
	resourceExecutor = resource.Executor{}
}

func (resourceHandler) Handle(w http.ResponseWriter, req *http.Request) {
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Monitoring()+URL.Agents(), "", -1)
	split := strings.Split(url, "/")
	switch len(split) {
	case 3:
		agentID := split[1]
		if "/"+split[2] == URL.Resource() {
			if req.Method == GET {
				resourceAPI.agentGetResourceInfo(w, req, agentID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 4:
		if "/"+split[3] == URL.Performance() {
			agentID := split[1]
			if req.Method == GET {
				resourceAPI.agentGetPerformanceInfo(w, req, agentID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}
	}
}

// agentGetResourceInfo handles requests related to get agent's resource informaion
// identified by the given agentID.
//
//    paths: '/api/v1/agents/{agentID}/resource'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (resourceAPIExecutor) agentGetResourceInfo(w http.ResponseWriter, req *http.Request, agentId string) {
	logger.Logging(logger.DEBUG, "[AGENT] Get Resource Info")
	result, resp, err := resourceExecutor.GetResourceInfo(agentId)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentGetPerformanceInfo handles requests related to get agent's resource performance informaion
// identified by the given agentID.
//
//    paths: '/api/v1/agents/{agentID}/resource/performance'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (resourceAPIExecutor) agentGetPerformanceInfo(w http.ResponseWriter, req *http.Request, agentId string) {
	logger.Logging(logger.DEBUG, "[AGENT] Get Performance Info")
	result, resp, err := resourceExecutor.GetPerformanceInfo(agentId)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}
