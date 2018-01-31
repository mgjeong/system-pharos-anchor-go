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

package search

import (
	"api/common"
	"api/search/app"
	"api/search/group"
	"api/search/node"
	"commons/errors"
	"commons/logger"
	URL "commons/url"
	"net/http"
	"strings"
)

type Command interface {
	Handle(w http.ResponseWriter, req *http.Request)
}

type RequestHandler struct{}

var nodeSearchHandler node.Command
var groupSearchHandler group.Command
var appSearchHandler app.Command

func init() {
	nodeSearchHandler = node.RequestHandler{}
	groupSearchHandler = group.RequestHandler{}
	appSearchHandler = app.RequestHandler{}
}

func (RequestHandler) Handle(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "receive msg", req.Method, req.URL.Path)
	defer logger.Logging(logger.DEBUG, "OUT")

	switch url := req.URL.Path; {
	default:
		logger.Logging(logger.DEBUG, "Unknown URL")
		common.WriteError(w, errors.NotFoundURL{})

	case !strings.Contains(url, URL.Base()):
		logger.Logging(logger.DEBUG, "Unknown URL")
		common.WriteError(w, errors.NotFoundURL{})

	case strings.Contains(url, URL.Nodes()):
		logger.Logging(logger.DEBUG, "Request Nodes APIs")
		nodeSearchHandler.Handle(w, req)

	case strings.Contains(url, URL.Groups()):
		logger.Logging(logger.DEBUG, "Request Groups APIs")
		groupSearchHandler.Handle(w, req)

	case strings.Contains(url, URL.Apps()):
		logger.Logging(logger.DEBUG, "Request Apps APIs")
		appSearchHandler.Handle(w, req)
	}
}
