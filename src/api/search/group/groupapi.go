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

// Package api/search/group provides functionality to handle request related to node.
package group

import (
	"api/common"
	"commons/errors"
	"commons/logger"
	URL "commons/url"
	groupSearcher "controller/search/group"
	"net/http"
	"strings"
)

const (
	GET string = "GET"
)

type Command interface {
	Handle(w http.ResponseWriter, req *http.Request)
}

type groupSearchAPI interface {
	searchGroups(w http.ResponseWriter, req *http.Request)
}

type RequestHandler struct{}
type groupAPIExecutor struct {
	groupSearchAPI
}

var searchExecutor groupSearcher.Command
var groupAPI groupAPIExecutor

func init() {
	searchExecutor = groupSearcher.Executor{}
	groupAPI = groupAPIExecutor{}
}

// Handle calls a proper function according to the url and method received from remote device.
func (RequestHandler) Handle(w http.ResponseWriter, req *http.Request) {
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Search()+URL.Groups(), "", -1)
	split := strings.Split(url, "/")

	switch len(split) {
	case 1:
		if req.Method == GET {
			groupAPI.searchGroups(w, req)
		} else {
			common.WriteError(w, errors.InvalidMethod{req.Method})
		}
	}
}

func (groupAPIExecutor) searchGroups(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[GROUP] Search Group")

	result, resp, err := searchExecutor.SearchGroups(parseQuery(req))
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

func parseQuery(req *http.Request) map[string]interface{} {
	query := make(map[string]interface{})

	keys := req.URL.Query()
	if len(keys) == 0 {
		return nil
	}

	for key, value := range req.URL.Query() {
		query[key] = value
	}

	return query
}
