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
	"bytes"
	notificationmocks "controller/notification/mocks"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	EVENT_ID = "eventId"
	BODY     = `{"test":"body"}`
)

var testBody = map[string]interface{}{
	"test": "body",
}

var Handler Command

func init() {
	Handler = RequestHandler{}
}

func TestNotificationHandlerWithInvalidMethod_ExpectReturnInvalidMethodMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	notiMockObj := notificationmocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notification", nil)
	// pass mockObj to a real object.
	notiExecutor = notiMockObj
	Handler.Handle(w, req)

	msg := make(map[string]interface{})
	err := json.Unmarshal(w.Body.Bytes(), &msg)
	if err != nil {
		t.Error("Expected results : invalid method msg, Actual err : json unmarshal failed.")
	}

	if !strings.Contains(msg["message"].(string), "invalid method") {
		t.Errorf("Expected results : invalid method msg, Actual err : %s.", msg["message"])
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/notification", nil)
	// pass mockObj to a real object.
	notiExecutor = notiMockObj
	Handler.Handle(w, req)

	msg = make(map[string]interface{})
	err = json.Unmarshal(w.Body.Bytes(), &msg)
	if err != nil {
		t.Error("Expected results : invalid method msg, Actual err : json unmarshal failed.")
	}

	if !strings.Contains(msg["message"].(string), "invalid method") {
		t.Errorf("Expected results : invalid method msg, Actual err : %s.", msg["message"])
	}
}

func TestNotificationHandlerWithInvalidUrl_ExpectReturnNotFoundURLMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	notiMockObj := notificationmocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/invalid", nil)

	// pass mockObj to a real object.
	notiExecutor = notiMockObj
	Handler.Handle(w, req)

	msg := make(map[string]interface{})
	err := json.Unmarshal(w.Body.Bytes(), &msg)
	if err != nil {
		t.Error("Expected results : unsupported url msg, Actual err : json unmarshal failed.")
	}

	if !strings.Contains(msg["message"].(string), "unsupported url") {
		t.Errorf("Expected results : unsupported url msg, Actual err : %s.", msg["message"])
	}
}

func TestNotificationHandlerWithValidRegisterRequest_ExpectCalledRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	notiMockObj := notificationmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		notiMockObj.EXPECT().Register(BODY, gomock.Any()),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/notification", bytes.NewReader(body))

	// pass mockObj to a real object.
	notiExecutor = notiMockObj

	Handler.Handle(w, req)
}

func TestNotificationHandlerWithValidUnRegisterRequest_ExpectCalledUnRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	notiMockObj := notificationmocks.NewMockCommand(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notification/"+EVENT_ID, nil)

	gomock.InOrder(
		notiMockObj.EXPECT().UnRegister(EVENT_ID),
	)

	// pass mockObj to a real object.
	notiExecutor = notiMockObj

	Handler.Handle(w, req)
}

func TestNotificationHandlerWithValidEventRequest_ExpectCalledHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	notiMockObj := notificationmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		notiMockObj.EXPECT().NotificationHandler(gomock.Any(), BODY),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testBody)
	req, _ := http.NewRequest("POST", "/api/v1/notification/events", bytes.NewReader(body))

	// pass mockObj to a real object.
	notiExecutor = notiMockObj

	Handler.Handle(w, req)
}
