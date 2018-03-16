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
	Handle(w http.ResponseWriter, req *http.Request)
}

type resourceMonitoringAPI interface {
	getNodeResourceInfo(w http.ResponseWriter, req *http.Request, nodeId string)
	getAppResourceInfo(w http.ResponseWriter, req *http.Request, nodeId string)
}

type RequestHandler struct{}
type resourceAPIExecutor struct {
	resourceMonitoringAPI
}

var resourceAPI resourceAPIExecutor
var resourceExecutor resource.Command

func init() {
	resourceExecutor = resource.Executor{}
}

func (RequestHandler) Handle(w http.ResponseWriter, req *http.Request) {
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Monitoring()+URL.Nodes(), "", -1)
	split := strings.Split(url, "/")
	switch len(split) {
	case 3: // [,{nodeId},resource]
		nodeId := split[1]
		if "/"+split[2] == URL.Resource() {
			if req.Method == GET {
				resourceAPI.getNodeResourceInfo(w, req, nodeId)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 5: // [,{nodeId},apps,{appId},resource]
		if "/"+split[2] == URL.Apps() {
			nodeId := split[1]
			appId := split[3]
			if req.Method == GET {
				resourceAPI.getAppResourceInfo(w, req, nodeId, appId)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}
	}
}

// getNodeResourceInfo handles requests related to get node's resource informaion
// identified by the given nodeId.
//
//    paths: '/api/v1/monitoring/nodes/{nodeId}/resource'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (resourceAPIExecutor) getNodeResourceInfo(w http.ResponseWriter, req *http.Request, nodeId string) {
	logger.Logging(logger.DEBUG, "[NODE] Get Resource Info")
	result, resp, err := resourceExecutor.GetNodeResourceInfo(nodeId)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// getAppResourceInfo handles requests related to get app's resource informaion deployed on the specific node.
// identified by the given nodeId, appId.
//
//    paths: '/api/v1/monitoring/nodes/{nodeId}/apps/{appId}/resource'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (resourceAPIExecutor) getAppResourceInfo(w http.ResponseWriter, req *http.Request, nodeId string, appId string) {
	logger.Logging(logger.DEBUG, "[NODE] Get Performance Info")
	result, resp, err := resourceExecutor.GetAppResourceInfo(nodeId, appId)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}
