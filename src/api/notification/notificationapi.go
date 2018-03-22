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

package notification

import (
	"api/common"
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	noti "controller/notification"
	"net/http"
	"strings"
)

const (
	POST   string = "POST"
	DELETE string = "DELETE"
)

type Command interface {
	Handle(w http.ResponseWriter, req *http.Request)
}

type notificationEventAPI interface {
	registerNotificationEvent(w http.ResponseWriter, req *http.Request)
	unRegisterNotificationEvent(w http.ResponseWriter, req *http.Request, eventId string)
}

type RequestHandler struct{}
type notificationAPIExecutor struct {
	notificationEventAPI
}

var notiExecutor noti.Command
var notificationAPI notificationAPIExecutor

func init() {
	notiExecutor = noti.Executor{}
}

func (RequestHandler) Handle(w http.ResponseWriter, req *http.Request) {
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Notification(), "", -1)
	split := strings.Split(url, "/")

	switch len(split) {
	default:
		logger.Logging(logger.DEBUG, "Unknown URL")
		common.WriteError(w, errors.NotFoundURL{})
	case 1:
		if req.Method == POST {
			notificationAPI.registerNotificationEvent(w, req)
		} else {
			common.WriteError(w, errors.InvalidMethod{req.Method})
		}
	case 2:
		eventId := split[1]
		if req.Method == DELETE {
			notificationAPI.unRegisterNotificationEvent(w, req, eventId)
		} else {
			common.WriteError(w, errors.InvalidMethod{req.Method})
		}
	}
}

func (notificationAPIExecutor) registerNotificationEvent(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[Notification] registration")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := notiExecutor.Register(body, parseQuery(req))
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

func (notificationAPIExecutor) unRegisterNotificationEvent(w http.ResponseWriter, req *http.Request, eventId string) {
	logger.Logging(logger.DEBUG, "[Notification] registration")

	result, err := notiExecutor.UnRegister(eventId)
	common.MakeResponse(w, result, nil, err)
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
