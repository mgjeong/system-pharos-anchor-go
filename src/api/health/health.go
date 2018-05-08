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
package health

import (
	"api/common"
	"commons/logger"
	"commons/url"
	"controller/health"
	"net/http"
	"strings"
)

type Command interface {
	Handle(w http.ResponseWriter, req *http.Request)
}

type apiInnerCommand interface {
	ping(w http.ResponseWriter, req *http.Request)
}

type RequestHandler struct{}
type innerExecutorImpl struct{}

var apiInnerExecutor apiInnerCommand
var healthExecutor health.Command

func init() {
	apiInnerExecutor = innerExecutorImpl{}
	healthExecutor = health.Executor{}
}

func (RequestHandler) Handle(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG)
	defer logger.Logging(logger.DEBUG, "OUT")

	switch reqUrl := req.URL.Path; {
	case strings.Contains(reqUrl, url.Ping()):
		apiInnerExecutor.ping(w, req)
	}
}

// ping handles requests which is used to check whether a pharos-anchor is up.
func (innerExecutorImpl) ping(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG)
	defer logger.Logging(logger.DEBUG, "OUT")

	result, err := healthExecutor.Ping()

	common.MakeResponse(w, result, common.ChangeToJson(nil), err)
}
