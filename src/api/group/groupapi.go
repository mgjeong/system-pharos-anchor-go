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

// Package api/group provides functionality to handle request related to group.
package group

import (
	"api/common"
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	groupmanager "controller/management/group"
	"net/http"
	"strings"
)

const (
	GET    string = "GET"
	PUT    string = "PUT"
	POST   string = "POST"
	DELETE string = "DELETE"
)

type _SDAMGroupApisHandler struct{}
type _SDAMGroupApis struct{}

var sdamH _SDAMGroupApisHandler
var sdam _SDAMGroupApis
var sdamGroupManager groupmanager.GroupInterface

func init() {
	SdamGroupHandle = sdamH
	SdamGroup = sdam
	sdamGroupManager = groupmanager.GroupController{}
}

// Handle calls a proper function according to the url and method received from remote device.
func (sdamH _SDAMGroupApisHandler) Handle(w http.ResponseWriter, req *http.Request) {
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Groups(), "", -1)
	split := strings.Split(url, "/")
	switch len(split) {
	case 1:
		if req.Method == GET {
			SdamGroup.groups(w, req)
		} else {
			common.WriteError(w, errors.InvalidMethod{req.Method})
		}

	case 2:
		if "/"+split[1] == URL.Create() {
			if req.Method == POST {
				SdamGroup.createGroup(w, req)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			if req.Method == GET || req.Method == DELETE {
				groupID := split[1]
				SdamGroup.group(w, req, groupID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		}

	case 3:
		groupID := split[1]
		switch {
		case "/"+split[2] == URL.Deploy():
			if req.Method == POST {
				SdamGroup.groupDeployApp(w, req, groupID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}

		case "/"+split[2] == URL.Join():
			if req.Method == POST {
				SdamGroup.groupJoin(w, req, groupID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}

		case "/"+split[2] == URL.Leave():
			if req.Method == POST {
				SdamGroup.groupLeave(w, req, groupID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}

		case "/"+split[2] == URL.Apps():
			if req.Method == GET {
				SdamGroup.groupInfoApps(w, req, groupID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}

		default:
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 4:
		if "/"+split[2] == URL.Apps() {
			groupID, appID := split[1], split[3]
			switch req.Method {
			case GET:
				SdamGroup.groupInfoApp(w, req, groupID, appID)

			case POST:
				SdamGroup.groupUpdateAppInfo(w, req, groupID, appID)

			case DELETE:
				SdamGroup.groupDeleteApp(w, req, groupID, appID)

			default:
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 5:
		if "/"+split[2] == URL.Apps() {
			groupID, appID := split[1], split[3]
			switch {
			case "/"+split[4] == URL.Start() && req.Method == POST:
				SdamGroup.groupStartApp(w, req, groupID, appID)

			case "/"+split[4] == URL.Stop() && req.Method == POST:
				SdamGroup.groupStopApp(w, req, groupID, appID)

			case "/"+split[4] == URL.Update() && req.Method == POST:
				SdamGroup.groupUpdateApp(w, req, groupID, appID)

			default:
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}
	}
}

// createGroup handles requests which is used to create new group.
//
//    paths: '/api/v1/groups/create'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (Groupasdam _SDAMGroupApis) createGroup(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[GROUP] Create SDA Group")
	result, resp, err := sdamGroupManager.CreateGroup()
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// group handles requests which is used to get information of group identified by the given groupID.
//
//    paths: '/api/v1/groups/{groupID}'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (Groupasdam _SDAMGroupApis) group(w http.ResponseWriter, req *http.Request, groupID string) {
	var result int
	var resp map[string]interface{}
	var err error
	switch req.Method {
	case GET:
		logger.Logging(logger.DEBUG, "[GROUP] Get SDA Group")
		result, resp, err = sdamGroupManager.GetGroup(groupID)
	case DELETE:
		logger.Logging(logger.DEBUG, "[GROUP] Delete SDA Group")
		result, resp, err = sdamGroupManager.DeleteGroup(groupID)
	}

	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groups handles requests which is used to get information of all groups created.
//
//    paths: '/api/v1/groups'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (Groupasdam _SDAMGroupApis) groups(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[GROUP] Get All SDA Groups")
	result, resp, err := sdamGroupManager.GetGroups()
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupJoin handles requests which is used to add an agent to a list of group members
// identified by the given groupID.
//
//    paths: '/api/v1/groups/{groupID}/join'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (Groupasdam _SDAMGroupApis) groupJoin(w http.ResponseWriter, req *http.Request, groupID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Join SDA Group")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := sdamGroupManager.JoinGroup(groupID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupLeave handles requests which is used to delete an agent from a list of group members
// identified by the given groupID.
//
//    paths: '/api/v1/groups/{groupID}/leave'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (Groupasdam _SDAMGroupApis) groupLeave(w http.ResponseWriter, req *http.Request, groupID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Leave SDA Group")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := sdamGroupManager.LeaveGroup(groupID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupDeployApp handles requests which is used to deploy new application to group
// identified by the given groupID.
//
//    paths: '/api/v1/groups/{groupID}/apps/deploy'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (Groupasdam _SDAMGroupApis) groupDeployApp(w http.ResponseWriter, req *http.Request, groupID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Deploy App")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := sdamGroupController.DeployApp(groupID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupInfoApps handles requests which is used to get information of all applications
// installed on group identified by the given groupID.
//
//    paths: '/api/v1/groups/{groupID}/apps'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (Groupasdam _SDAMGroupApis) groupInfoApps(w http.ResponseWriter, req *http.Request, groupID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Get Info Apps")
	result, resp, err := sdamGroupController.GetApps(groupID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupInfoApp handles requests which is used to get information of application
// identified by the given appID.
//
//    paths: '/api/v1/groups/{groupID}/apps/{appID}'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (Groupasdam _SDAMGroupApis) groupInfoApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Get Info App")
	result, resp, err := sdamGroupController.GetApp(groupID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupUpdateAppInfo handles requests related to updating application installed on group
// with given yaml in body.
//
//    paths: '/api/v1/groups/{groupID}/apps/{appID}'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (Groupasdam _SDAMGroupApis) groupUpdateAppInfo(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Update App Info")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := sdamGroupController.UpdateAppInfo(groupID, appID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupDeleteApp handles requests related to delete application installed on group
// identified by the given groupID.
//
//    paths: '/api/v1/groups/{groupID}/apps/{appID}'
//    method: DELETE
//    responses: if successful, 200 status code will be returned.
func (Groupasdam _SDAMGroupApis) groupDeleteApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Delete App")
	result, resp, err := sdamGroupController.DeleteApp(groupID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupStartApp handles requests related to start application installed on group
// identified by the given groupID.
//
//    paths: '/api/v1/groups/{groupID}/apps/{appID}/start'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (Groupasdam _SDAMGroupApis) groupStartApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Start App")
	result, resp, err := sdamGroupController.StartApp(groupID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupStopApp handles requests related to stop application installed on group
// identified by the given groupID.
//
//    paths: '/api/v1/groups/{groupID}/apps/{appID}/stop'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (Groupasdam _SDAMGroupApis) groupStopApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Stop App")
	result, resp, err := sdamGroupController.StopApp(groupID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupUpdateApp handles requests related to updating application installed on group
// identified by the given groupID.
//
//    paths: '/api/v1/groups/{groupID}/apps/{appID}/update'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (Groupasdam _SDAMGroupApis) groupUpdateApp(w http.ResponseWriter, req *http.Request, groupID string, appID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Update App")
	result, resp, err := sdamGroupController.UpdateApp(groupID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}
