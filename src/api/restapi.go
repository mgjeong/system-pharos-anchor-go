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

// Package api provides web server for Pharos Anchor(Edge Manager).
// and also provides functionality of request processing and response making.
package api

import (
	"api/management"
	"api/monitoring"
	"api/common"
	"commons/errors"
	"commons/logger"
	URL "commons/url"
	"net/http"
	"strconv"
	"strings"
)

var managementHandler management.Command
var monitoringHandler monitoring.Command

func init() {
	managementHandler = management.RequestHandler{}
	monitoringHandler = monitoring.RequestHandler{}
}

// RunWebServer starts web server service with given address and port number.
func RunWebServer(addr string, port int) {
	http.ListenAndServe(addr+":"+strconv.Itoa(port), &Handler)
}

var Handler RequestHandler

type RequestHandler struct{}

func (RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "receive msg", req.Method, req.URL.Path)
	defer logger.Logging(logger.DEBUG, "OUT")

	switch url := req.URL.Path; {
	default:
		logger.Logging(logger.DEBUG, "Unknown URL")
		common.WriteError(w, errors.NotFoundURL{})

	case !strings.Contains(url, URL.Base()):
		logger.Logging(logger.DEBUG, "Unknown URL")
		common.WriteError(w, errors.NotFoundURL{})

	case strings.Contains(url, URL.Management()):
		logger.Logging(logger.DEBUG, "Request Management APIs")
		managementHandler.Handle(w, req)

	case strings.Contains(url, URL.Monitoring()):
		logger.Logging(logger.DEBUG, "Request Monitoring APIs")
		monitoringHandler.Handle(w, req)
	}
}
