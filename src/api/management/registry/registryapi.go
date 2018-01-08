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

// Package api/registry provides functionality to handle request related to agent.
package registry

import (
	"api/common"
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	"controller/management/registry"
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
	registerDockerRegistry(w http.ResponseWriter, req *http.Request)
	getDockerRegistries(w http.ResponseWriter, req *http.Request)
	getDockerRegistry(w http.ResponseWriter, req *http.Request, registryID string)
	handleDockerRegistryEvent(w http.ResponseWriter, req *http.Request)
}

type registryHandler struct{}
type registryAPIExecutor struct{
	Command
}

var Handler registryHandler
var registryAPI registryAPIExecutor
var registryExecutor registry.Command

func init() {
	registryAPI = registryAPIExecutor{}
	registryExecutor = registry.Executor{}
	Handler = registryHandler{}
}

// Handle calls a proper function according to the url and method received from remote device.
func (registryHandler) Handle(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	url := strings.Replace(req.URL.Path, URL.Base()+URL.Registry(), "", -1)
	split := strings.Split(url, "/")

	switch len(split) {
	case 1:
		if req.Method == GET {
			registryAPI.getDockerRegistries(w, req)
		} else if req.Method == POST {
			registryAPI.registerDockerRegistry(w, req)
		} else {
			common.WriteError(w, errors.InvalidMethod{req.Method})
		}

	case 2:
		if "/"+split[1] == URL.Events() {
			if req.Method == POST {
				registryAPI.handleDockerRegistryEvent(w, req)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			if req.Method == GET {
				registryID := split[1]
				registryAPI.getDockerRegistry(w, req, registryID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		}
	}
}

func (registryAPIExecutor) registerDockerRegistry(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := registryExecutor.AddDockerRegistry(body)

	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

func (registryAPIExecutor) getDockerRegistries(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	result, resp, err := registryExecutor.GetDockerRegistries()

	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

func (registryAPIExecutor) getDockerRegistry(w http.ResponseWriter, req *http.Request, registryID string) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	result, resp, err := registryExecutor.GetDockerRegistry(registryID)

	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

func (registryAPIExecutor) handleDockerRegistryEvent(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, err := registryExecutor.DockerRegistryEventHandler(body)

	common.MakeResponse(w, result, common.ChangeToJson(nil), err)
}
