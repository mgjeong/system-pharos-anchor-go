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
	resource "controller/resource/node"
	"net/http"
	"strings"
)

const (
	GET    string = "GET"
)

type Command interface {
	nodeGetResourceInfo(w http.ResponseWriter, req *http.Request, nodeId string)
	nodeGetPerformanceInfo(w http.ResponseWriter, req *http.Request, nodeId string)
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
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Monitoring()+URL.Nodes(), "", -1)
	split := strings.Split(url, "/")
	switch len(split) {
	case 3:
		nodeID := split[1]
		if "/"+split[2] == URL.Resource() {
			if req.Method == GET {
				resourceAPI.nodeGetResourceInfo(w, req, nodeID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 4:
		if "/"+split[3] == URL.Performance() {
			nodeID := split[1]
			if req.Method == GET {
				resourceAPI.nodeGetPerformanceInfo(w, req, nodeID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}
	}
}

// nodeGetResourceInfo handles requests related to get node's resource informaion
// identified by the given nodeID.
//
//    paths: '/api/v1/monitoring/nodes/{nodeID}/resource'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (resourceAPIExecutor) nodeGetResourceInfo(w http.ResponseWriter, req *http.Request, nodeId string) {
	logger.Logging(logger.DEBUG, "[NODE] Get Resource Info")
	result, resp, err := resourceExecutor.GetResourceInfo(nodeId)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// nodeGetPerformanceInfo handles requests related to get node's resource performance informaion
// identified by the given nodeID.
//
//    paths: '/api/v1/monitoring/nodes/{nodeID}/resource/performance'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (resourceAPIExecutor) nodeGetPerformanceInfo(w http.ResponseWriter, req *http.Request, nodeId string) {
	logger.Logging(logger.DEBUG, "[NODE] Get Performance Info")
	result, resp, err := resourceExecutor.GetPerformanceInfo(nodeId)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}
