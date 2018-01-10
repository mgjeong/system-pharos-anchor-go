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

// Package api/management/group provides functionality to handle request related to group.
package apps

import (
	"api/common"
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	deployment "controller/deployment/group"
	"net/http"
	"strings"
)

const (
	GET    string = "GET"
	PUT    string = "PUT"
	POST   string = "POST"
	DELETE string = "DELETE"
)

type Command interface {
	groupDeployApp(w http.ResponseWriter, req *http.Request, groupID string)
	groupInfoApps(w http.ResponseWriter, req *http.Request, groupID string)
	groupInfoApp(w http.ResponseWriter, req *http.Request, groupID string, appID string)
	groupUpdateAppInfo(w http.ResponseWriter, req *http.Request, groupID string, appID string)
	groupDeleteApp(w http.ResponseWriter, req *http.Request, groupID string, appID string)
	groupStartApp(w http.ResponseWriter, req *http.Request, groupID string, appID string)
	groupStopApp(w http.ResponseWriter, req *http.Request, groupID string, appID string)
	groupUpdateApp(w http.ResponseWriter, req *http.Request, groupID string, appID string)
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
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Management()+URL.Groups(), "", -1)
	split := strings.Split(url, "/")
	switch len(split) {
	case 3:
		groupID := split[1]
		switch {
		case "/"+split[2] == URL.Deploy():
			if req.Method == POST {
				appsAPI.groupDeployApp(w, req, groupID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}

		case "/"+split[2] == URL.Apps():
			if req.Method == GET {
				appsAPI.groupInfoApps(w, req, groupID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}

		default:
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 4:
		if "/"+split[2] == URL.Apps() {
			if "/"+split[3] == URL.Deploy() {
				if req.Method == POST {
					groupID := split[1]
					appsAPI.groupDeployApp(w, req, groupID)
				} else {
					common.WriteError(w, errors.InvalidMethod{req.Method})
				}
			} else {
				groupID, appID := split[1], split[3]
				switch req.Method {
				case GET:
					appsAPI.groupInfoApp(w, req, groupID, appID)

				case POST:
					appsAPI.groupUpdateAppInfo(w, req, groupID, appID)

				case DELETE:
					appsAPI.groupDeleteApp(w, req, groupID, appID)

				default:
					common.WriteError(w, errors.InvalidMethod{req.Method})
				}
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 5:
		if "/"+split[2] == URL.Apps() {
			groupID, appID := split[1], split[3]
			switch {
			case "/"+split[4] == URL.Start() && req.Method == POST:
				appsAPI.groupStartApp(w, req, groupID, appID)

			case "/"+split[4] == URL.Stop() && req.Method == POST:
				appsAPI.groupStopApp(w, req, groupID, appID)

			case "/"+split[4] == URL.Update() && req.Method == POST:
				appsAPI.groupUpdateApp(w, req, groupID, appID)

			default:
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}
	}
}

// groupDeployApp handles requests which is used to deploy new application to group
// identified by the given groupID.
//
//    paths: '/api/v1/management/groups/{groupID}/apps/deploy'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) groupDeployApp(w http.ResponseWriter, req *http.Request, groupID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Deploy App")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := deploymentExecutor.DeployApp(groupID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupInfoApps handles requests which is used to get information of all applications
// installed on group identified by the given groupID.
//
//    paths: '/api/v1/management/groups/{groupID}/apps'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) groupInfoApps(w http.ResponseWriter, req *http.Request, groupID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Get Info Apps")
	result, resp, err := deploymentExecutor.GetApps(groupID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupInfoApp handles requests which is used to get information of application
// identified by the given appID.
//
//    paths: '/api/v1/management/groups/{groupID}/apps/{appID}'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) groupInfoApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Get Info App")
	result, resp, err := deploymentExecutor.GetApp(groupID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupUpdateAppInfo handles requests related to updating application installed on group
// with given yaml in body.
//
//    paths: '/api/v1/management/groups/{groupID}/apps/{appID}'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) groupUpdateAppInfo(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Update App Info")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := deploymentExecutor.UpdateAppInfo(groupID, appID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupDeleteApp handles requests related to delete application installed on group
// identified by the given groupID.
//
//    paths: '/api/v1/management/groups/{groupID}/apps/{appID}'
//    method: DELETE
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) groupDeleteApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Delete App")
	result, resp, err := deploymentExecutor.DeleteApp(groupID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupStartApp handles requests related to start application installed on group
// identified by the given groupID.
//
//    paths: '/api/v1/management/groups/{groupID}/apps/{appID}/start'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) groupStartApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Start App")
	result, resp, err := deploymentExecutor.StartApp(groupID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupStopApp handles requests related to stop application installed on group
// identified by the given groupID.
//
//    paths: '/api/v1/management/groups/{groupID}/apps/{appID}/stop'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) groupStopApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Stop App")
	result, resp, err := deploymentExecutor.StopApp(groupID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupUpdateApp handles requests related to updating application installed on group
// identified by the given groupID.
//
//    paths: '/api/v1/management/groups/{groupID}/apps/{appID}/update'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (appsAPIExecutor) groupUpdateApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Update App")
	result, resp, err := deploymentExecutor.UpdateApp(groupID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}
