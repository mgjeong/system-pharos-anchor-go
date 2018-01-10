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

// Package api/node/apps provides functionality to handle request related to apps.
package apps

import (
	"api/common"
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	deployment "controller/deployment/node"
	"net/http"
	"strings"
)

const (
	GET    string = "GET"
	POST   string = "POST"
	DELETE string = "DELETE"
)

type Command interface {
	nodeDeployApp(w http.ResponseWriter, req *http.Request, nodeID string)
	nodeInfoApps(w http.ResponseWriter, req *http.Request, nodeID string)
	nodeInfoApp(w http.ResponseWriter, req *http.Request, nodeID string, appID string)
	nodeUpdateAppInfo(w http.ResponseWriter, req *http.Request, nodeID string, appID string)
	nodeDeleteApp(w http.ResponseWriter, req *http.Request, nodeID string, appID string)
	nodeStartApp(w http.ResponseWriter, req *http.Request, nodeID string, appID string)
	nodeStopApp(w http.ResponseWriter, req *http.Request, nodeID string, appID string)
	nodeUpdateApp(w http.ResponseWriter, req *http.Request, nodeID string, appID string)
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
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Management()+URL.Nodes(), "", -1)
	split := strings.Split(url, "/")
	switch len(split) {
	case 3:
		nodeID := split[1]
		if "/"+split[2] == URL.Apps() {
			if req.Method == GET {
				appsAPI.nodeInfoApps(w, req, nodeID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 4:
		if "/"+split[2] == URL.Apps() {
			if "/"+split[3] == URL.Deploy() {
				if req.Method == POST {
					nodeID := split[1]
					appsAPI.nodeDeployApp(w, req, nodeID)
				} else {
					common.WriteError(w, errors.InvalidMethod{req.Method})
				}
			} else {
				nodeID, appID := split[1], split[3]
				switch req.Method {
				case GET:
					appsAPI.nodeInfoApp(w, req, nodeID, appID)

				case POST:
					appsAPI.nodeUpdateAppInfo(w, req, nodeID, appID)

				case DELETE:
					appsAPI.nodeDeleteApp(w, req, nodeID, appID)

				default:
					common.WriteError(w, errors.InvalidMethod{req.Method})
				}
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 5:
		if "/"+split[2] == URL.Apps() {
			nodeID, appID := split[1], split[3]
			switch {
			case "/"+split[4] == URL.Start() && req.Method == POST:
				appsAPI.nodeStartApp(w, req, nodeID, appID)

			case "/"+split[4] == URL.Stop() && req.Method == POST:
				appsAPI.nodeStopApp(w, req, nodeID, appID)

			case "/"+split[4] == URL.Update() && req.Method == POST:
				appsAPI.nodeUpdateApp(w, req, nodeID, appID)

			default:
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}
	}
}

// nodeDeployApp handles requests which is used to deploy new application to node
// identified by the given nodeID.
//
//    paths: '/api/v1/management/nodes/{nodeID}/apps/deploy'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) nodeDeployApp(w http.ResponseWriter, req *http.Request, nodeID string) {
	logger.Logging(logger.DEBUG, "[NODE] Deploy App")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := deploymentExecutor.DeployApp(nodeID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// nodeInfoApps handles requests which is used to get information of all applications
// installed on node identified by the given nodeID.
//
//    paths: '/api/v1/management/nodes/{nodeID}/apps'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) nodeInfoApps(w http.ResponseWriter, req *http.Request, nodeID string) {
	logger.Logging(logger.DEBUG, "[NODE] Get Info Apps")
	result, resp, err := deploymentExecutor.GetApps(nodeID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// nodeInfoApp handles requests which is used to get information of application
// identified by the given appID.
//
//    paths: '/api/v1/management/nodes/{nodeID}/apps/{appID}'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) nodeInfoApp(w http.ResponseWriter, req *http.Request, nodeID string, appID string) {
	logger.Logging(logger.DEBUG, "[NODE] Get Info App")
	result, resp, err := deploymentExecutor.GetApp(nodeID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// nodeUpdateAppInfo handles requests related to updating the application with given yaml in body.
//
//    paths: '/api/v1/management/nodes/{nodeID}/apps/{appID}'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) nodeUpdateAppInfo(w http.ResponseWriter, req *http.Request, nodeID string, appID string) {
	logger.Logging(logger.DEBUG, "[NODE] Update App Info")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := deploymentExecutor.UpdateAppInfo(nodeID, appID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// nodeDeleteApp handles requests related to delete application installed on node
// identified by the given nodeID.
//
//    paths: '/api/v1/management/nodes/{nodeID}/apps/{appID}'
//    method: DELETE
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) nodeDeleteApp(w http.ResponseWriter, req *http.Request, nodeID string, appID string) {
	logger.Logging(logger.DEBUG, "[NODE] Delete App")
	result, resp, err := deploymentExecutor.DeleteApp(nodeID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// nodeStartApp handles requests related to start application installed on node
// identified by the given nodeID.
//
//    paths: '/api/v1/management/nodes/{nodeID}/apps/{appID}/start'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) nodeStartApp(w http.ResponseWriter, req *http.Request, nodeID string, appID string) {
	logger.Logging(logger.DEBUG, "[NODE] Start App")
	result, resp, err := deploymentExecutor.StartApp(nodeID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// nodeStopApp handles requests related to stop application installed on node
// identified by the given nodeID.
//
//    paths: '/api/v1/management/nodes/{nodeID}/apps/{appID}/stop'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) nodeStopApp(w http.ResponseWriter, req *http.Request, nodeID string, appID string) {
	logger.Logging(logger.DEBUG, "[NODE] Stop App")
	result, resp, err := deploymentExecutor.StopApp(nodeID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// nodeUpdateApp handles requests related to updating application installed on node
// identified by the given nodeID.
//
//    paths: '/api/v1/management/nodes/{nodeID}/apps/{appID}/update'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) nodeUpdateApp(w http.ResponseWriter, req *http.Request, nodeID string, appID string) {
	logger.Logging(logger.DEBUG, "[NODE] Update App")
	result, resp, err := deploymentExecutor.UpdateApp(nodeID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}
