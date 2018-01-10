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
	GET string = "GET"
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
			} else {
				common.WriteError(w, errors.NotFoundURL{})
			}
		}

	}

}

// nodes handles requests which is used to get information of node identified by the given nodeID.
//
//    paths: '/api/v1/management/nodes/{nodeID}'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (nodeAPIExecutor) node(w http.ResponseWriter, req *http.Request, nodeID string) {
	logger.Logging(logger.DEBUG, "[NODE] Get Service Deployment Nodes")
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
	logger.Logging(logger.DEBUG, "[NODE] Get All Service Deployment Nodes")
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
	logger.Logging(logger.DEBUG, "[NODE] Register New Service Deployment Node")

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
	logger.Logging(logger.DEBUG, "[NODE] Unregister New Service Deployment Node")

	result, err := managementExecutor.UnRegisterNode(nodeID)
	common.MakeResponse(w, result, nil, err)
}

// ping handles requests which is used to check whether a node is up.
//
//    paths: '/api/v1/management/nodes/{nodeID}/ping'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (nodeAPIExecutor) ping(w http.ResponseWriter, req *http.Request, nodeID string) {
	logger.Logging(logger.DEBUG, "[NODE] Ping From Service Deployment Node")

	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, err := managementExecutor.PingNode(nodeID, body)
	common.MakeResponse(w, result, nil, err)
}
