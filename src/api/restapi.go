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

// Package api provides web server for Service Deployment Agent Manager.
// and also provides functionality of request processing and response making.
package api

import (
	"api/agent"
	"api/common"
	"api/group"
	"commons/logger"
	"commons/errors"
	URL "commons/url"
	"net/http"
	"strconv"
	"strings"
)

// RunSDAMWebServer starts web server service with given address and port number.
func RunSDAMWebServer(addr string, port int) {
	http.ListenAndServe(addr+":"+strconv.Itoa(port), &_SDAMApis)
}

var _SDAMApis _SDAMApisHandler

type _SDAMApisHandler struct{}

// ServeHTTP implements a http serve interface.
// Check if the url contains a given string and call a proper function.
//
//    agents: agent.SdamAgentHandle.Handle will be called.
//	  groups: group.SdamGroupHandle.Handle will be called.
//    others: NotFoundURL error will be used to send an error message.
func (_SDAMApis *_SDAMApisHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "receive msg", req.Method, req.URL.Path)
	defer logger.Logging(logger.DEBUG, "OUT")

	switch url := req.URL.Path; {
	default:
		logger.Logging(logger.DEBUG, "Unknown URL")
		common.WriteError(w, errors.NotFoundURL{})

	case !strings.Contains(url, URL.Base()):
		logger.Logging(logger.DEBUG, "Unknown URL")
		common.WriteError(w, errors.NotFoundURL{})

	case strings.Contains(url, URL.Agents()):
		logger.Logging(logger.DEBUG, "Request Agents APIs")
		agent.SdamAgentHandle.Handle(w, req)

	case strings.Contains(url, URL.Groups()):
		logger.Logging(logger.DEBUG, "Request Groups APIs")
		group.SdamGroupHandle.Handle(w, req)
	}
}
