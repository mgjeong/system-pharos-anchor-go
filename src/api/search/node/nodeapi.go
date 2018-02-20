/*******************************************************************************
 * Copyright 2018 Samsung Electronics All Rights Reserved.
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
// Package api/search/node provides functionality to handle request related to node.
package node

import (
	"api/common"
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	nodeSearch "controller/search/node"
	"net/http"
	"strings"
)

const (
	GET string = "GET"
)

type Command interface {
	Handle(w http.ResponseWriter, req *http.Request)
}

type nodeSearchAPI interface {
	searchNodes(w http.ResponseWriter, req *http.Request)
}

type RequestHandler struct{}
type nodeAPIExecutor struct {
	nodeSearchAPI
}

var searchExecutor nodeSearch.Command
var nodeAPI nodeAPIExecutor

func init() {
	searchExecutor = nodeSearch.Executor{}
}

// Handle calls a proper function according to the url and method received from remote device.
func (RequestHandler) Handle(w http.ResponseWriter, req *http.Request) {
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Search()+URL.Nodes(), "", -1)
	split := strings.Split(url, "/")

	switch len(split) {
	default:
		logger.Logging(logger.DEBUG, "Unknown URL")
		common.WriteError(w, errors.NotFoundURL{})
	case 1:
		if req.Method == GET {
			nodeAPI.searchNodes(w, req)
		} else {
			common.WriteError(w, errors.InvalidMethod{req.Method})
		}
	}
}

func (nodeAPIExecutor) searchNodes(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[NODE] Get Nodes maching the condition")

	result, resp, err := searchExecutor.SearchNodes(req.URL.Query())
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}
