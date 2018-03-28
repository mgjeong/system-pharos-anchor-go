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

// Package api/node provides functionality to handle request related to node.
package node

import (
	"api/common"
	"api/management/node/apps"
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	nodemanager "controller/management/node"
	"net/http"
	"strings"
)

const (
	GET  string = "GET"
	POST string = "POST"
)

type Command interface {
	Handle(w http.ResponseWriter, req *http.Request)
}

type nodeManagementAPI interface {
	node(w http.ResponseWriter, req *http.Request, nodeID string)
	nodes(w http.ResponseWriter, req *http.Request)
	register(w http.ResponseWriter, req *http.Request)
	ping(w http.ResponseWriter, req *http.Request, nodeID string)
	unregister(w http.ResponseWriter, req *http.Request, nodeID string)
	configuration(w http.ResponseWriter, req *http.Request, nodeID string)
	reboot(w http.ResponseWriter, req *http.Request)
	restore(w http.ResponseWriter, req *http.Request)
}

type RequestHandler struct{}
type nodeAPIExecutor struct {
	nodeManagementAPI
}

var deploymentHandler apps.Command
var managementExecutor nodemanager.Command
var nodeAPI nodeAPIExecutor

func init() {
	deploymentHandler = apps.RequestHandler{}
	managementExecutor = nodemanager.Executor{}
	nodeAPI = nodeAPIExecutor{}
}

// Handle calls a proper function according to the url and method received from remote device.
func (RequestHandler) Handle(w http.ResponseWriter, req *http.Request) {
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Management()+URL.Nodes(), "", -1)
	split := strings.Split(url, "/")

	if strings.Contains(url, URL.Apps()) {
		deploymentHandler.Handle(w, req)
	} else if strings.Contains(url, URL.Reboot()) {
		nodeId := split[1]
		nodeAPI.reboot(w, req, nodeId)
	} else if strings.Contains(url, URL.Restore()) {
		nodeId := split[1]
		nodeAPI.restore(w, req, nodeId)
	} else {
		switch len(split) {
		case 1:
			if req.Method == GET {
				nodeAPI.nodes(w, req)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}

		case 2:
			if "/"+split[1] == URL.Register() {
				nodeAPI.register(w, req)
			} else {
				if req.Method == GET {
					nodeID := split[1]
					nodeAPI.node(w, req, nodeID)
				} else {
					common.WriteError(w, errors.InvalidMethod{req.Method})
				}
			}

		case 3:
			if "/"+split[2] == URL.Unregister() {
				agentID := split[1]
				nodeAPI.unregister(w, req, agentID)
			} else if "/"+split[2] == URL.Ping() {
				nodeID := split[1]
				nodeAPI.ping(w, req, nodeID)
			} else if "/"+split[2] == URL.Configuration() {
				nodeID := split[1]
				nodeAPI.configuration(w, req, nodeID)
			} else {
				common.WriteError(w, errors.NotFoundURL{})
			}
		}
	}
}

// nodes handles requests which is used to reboot a device with node.
//
//    paths: '/api/v1/management/nodes/{nodeID}/reboot'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (nodeAPIExecutor) reboot(w http.ResponseWriter, req *http.Request, nodeId string) {
	logger.Logging(logger.DEBUG, "[NODE] Reboot Pharos Nodes")
	result, err := managementExecutor.Reboot(nodeId)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	common.MakeResponse(w, result, nil, err)
}

// nodes handles requests which is used to restore a device to initial state.
//
//    paths: '/api/v1/management/nodes/{nodeID}/restore'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (nodeAPIExecutor) restore(w http.ResponseWriter, req *http.Request, nodeId string) {
	logger.Logging(logger.DEBUG, "[NODE] Restore Pharos Nodes")
	result, err := managementExecutor.Restore(nodeId)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	common.MakeResponse(w, result, nil, err)
}

// nodes handles requests which is used to get information of node identified by the given nodeID.
//
//    paths: '/api/v1/management/nodes/{nodeID}'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (nodeAPIExecutor) node(w http.ResponseWriter, req *http.Request, nodeID string) {
	logger.Logging(logger.DEBUG, "[NODE] Get Pharos Nodes")
	result, resp, err := managementExecutor.GetNode(nodeID)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// nodes handles requests which is used to get information of all nodes registered.
//
//    paths: '/api/v1/management/nodes'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (nodeAPIExecutor) nodes(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[NODE] Get All Pharos Nodes")
	result, resp, err := managementExecutor.GetNodes()
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// register handles requests which is used to register node to a list of nodes.
//
//    paths: '/api/v1/management/nodes/register'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (nodeAPIExecutor) register(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[NODE] Register New Pharos Node")

	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := managementExecutor.RegisterNode(body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// unregister handles requests which is used to unregister node from a list of nodes.
//
//    paths: '/api/v1/management/nodes/{nodeID}/unregister'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (nodeAPIExecutor) unregister(w http.ResponseWriter, req *http.Request, nodeID string) {
	logger.Logging(logger.DEBUG, "[NODE] Unregister New Pharos Node")

	result, err := managementExecutor.UnRegisterNode(nodeID)
	common.MakeResponse(w, result, nil, err)
}

// ping handles requests which is used to check whether a node is up.
//
//    paths: '/api/v1/management/nodes/{nodeID}/ping'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (nodeAPIExecutor) ping(w http.ResponseWriter, req *http.Request, nodeID string) {
	logger.Logging(logger.DEBUG, "[NODE] Ping From Pharos Node")

	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, err := managementExecutor.PingNode(nodeID, body)
	common.MakeResponse(w, result, nil, err)
}

//  configuration handles requests which is used to get/set a node configuration.
//
//    paths: '/api/v1/management/nodes/{nodeID}/configuration'
//    method: GET, POST
//    responses: if successful, 200 status code will be returned.
func (nodeAPIExecutor) configuration(w http.ResponseWriter, req *http.Request, nodeID string) {
	logger.Logging(logger.DEBUG, "[NODE] Configure Pharos Node")

	response := make(map[string]interface{})
	var result int
	var err error
	switch req.Method {
	case GET:
		result, response, err = managementExecutor.GetNodeConfiguration(nodeID)
	case POST:
		var bodyStr string
		bodyStr, err = common.GetBodyFromReq(req)
		if err != nil {
			common.MakeResponse(w, results.ERROR, nil, err)
			return
		}
		result, err = managementExecutor.SetNodeConfiguration(nodeID, bodyStr)
	}
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	common.MakeResponse(w, result, common.ChangeToJson(response), err)
}