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
	"api/management/group/apps"
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

type Command interface {
	createGroup(w http.ResponseWriter, req *http.Request)
	group(w http.ResponseWriter, req *http.Request, groupID string)
	groups(w http.ResponseWriter, req *http.Request)
	groupJoin(w http.ResponseWriter, req *http.Request, groupID string)
	groupLeave(w http.ResponseWriter, req *http.Request, groupID string)
}

type groupHandler struct{}
type groupAPIExecutor struct {
	Command
}

var managementExecutor groupmanager.Command
var groupAPI groupAPIExecutor
var Handler groupHandler

func init() {
	managementExecutor = groupmanager.Executor{}
	groupAPI = groupAPIExecutor{}
	Handler = groupHandler{}
}

// Handle calls a proper function according to the url and method received from remote device.
func (groupHandler) Handle(w http.ResponseWriter, req *http.Request) {
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Management()+URL.Groups(), "", -1)
	split := strings.Split(url, "/")
	switch len(split) {
	case 1:
		if req.Method == GET {
			groupAPI.groups(w, req)
		} else {
			common.WriteError(w, errors.InvalidMethod{req.Method})
		}

	case 2:
		if "/"+split[1] == URL.Create() {
			if req.Method == POST {
				groupAPI.createGroup(w, req)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			if req.Method == GET || req.Method == DELETE {
				groupID := split[1]
				groupAPI.group(w, req, groupID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		}

	case 3:
		groupID := split[1]
		switch {
		case "/"+split[2] == URL.Apps():
			apps.Handler.Handle(w, req)

		case "/"+split[2] == URL.Join():
			if req.Method == POST {
				groupAPI.groupJoin(w, req, groupID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}

		case "/"+split[2] == URL.Leave():
			if req.Method == POST {
				groupAPI.groupLeave(w, req, groupID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}

		default:
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

// createGroup handles requests which is used to create new group.
//
//    paths: '/api/v1/management/groups/create'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (groupAPIExecutor) createGroup(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[GROUP] Create Group")
	result, resp, err := managementExecutor.CreateGroup()
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// group handles requests which is used to get information of group identified by the given groupID.
//
//    paths: '/api/v1/management/groups/{groupID}'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (groupAPIExecutor) group(w http.ResponseWriter, req *http.Request, groupID string) {
	var result int
	var resp map[string]interface{}
	var err error
	switch req.Method {
	case GET:
		logger.Logging(logger.DEBUG, "[GROUP] Get Group")
		result, resp, err = managementExecutor.GetGroup(groupID)
	case DELETE:
		logger.Logging(logger.DEBUG, "[GROUP] Delete Group")
		result, resp, err = managementExecutor.DeleteGroup(groupID)
	}

	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groups handles requests which is used to get information of all groups created.
//
//    paths: '/api/v1/management/groups'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (groupAPIExecutor) groups(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[GROUP] Get All Groups")
	result, resp, err := managementExecutor.GetGroups()
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupJoin handles requests which is used to add an agent to a list of group members
// identified by the given groupID.
//
//    paths: '/api/v1/management/groups/{groupID}/join'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (groupAPIExecutor) groupJoin(w http.ResponseWriter, req *http.Request, groupID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Join Group")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := managementExecutor.JoinGroup(groupID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// groupLeave handles requests which is used to delete an agent from a list of group members
// identified by the given groupID.
//
//    paths: '/api/v1/management/groups/{groupID}/leave'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (groupAPIExecutor) groupLeave(w http.ResponseWriter, req *http.Request, groupID string) {
	logger.Logging(logger.DEBUG, "[GROUP] Leave Group")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := managementExecutor.LeaveGroup(groupID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}
