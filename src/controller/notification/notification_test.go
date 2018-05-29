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
	//	"commons/errors"
	"commons/results"
	nodeSearchmocks "controller/search/node/mocks"
	appEventDBmocks "db/mongo/event/app/mocks"
	nodeEventDBmocks "db/mongo/event/node/mocks"
	subsDBmocks "db/mongo/event/subscriber/mocks"
	"github.com/golang/mock/gomock"
	msgmocks "messenger/mocks"
	"testing"
)

const (
	URL_KEY  = "url"
	TEST_URL = "test-url"
	IP       = "192.168.0.1"
	PORT     = "48098"
)

var (
	watchUrl               = "http://" + IP + ":" + PORT + "/api/v1/notification/apps/watch"
	respCode               = []int{results.OK, results.OK}
	partialSuccessRespCode = []int{results.OK, results.ERROR}
	errorRespCode          = []int{results.ERROR, results.ERROR}
	eventId                = "92a1407cb237d05b4a985b34070ddad135bf8a0c"
	appsubsId              = "8f834351058adfffb19fc1e2f9ea3facd316ddff"
	nodesubsId             = "494d526547ed41d07f718b5bea633b4fd9181285"
	nodeIds                = []string{"nodeid", "nodeid"}
	appState               = []string{"stop"}
	nodeState              = []string{"disconnected"}

	validBody = map[string]interface{}{
		APP_ID:     "appid",
		EVENT_ID:   eventId,
		IMAGE_NAME: "imagename",
	}

	appEventBody = map[string]interface{}{
		URL_KEY: string(TEST_URL),
		EVENT: map[string]interface{}{
			TYPE:   APP,
			STATUS: appState,
		},
	}
	nodeEventBody = map[string]interface{}{
		URL_KEY: string(TEST_URL),
		EVENT: map[string]interface{}{
			TYPE:   NODE,
			STATUS: nodeState,
		},
	}
	eventBodyWithoutURL = map[string]interface{}{
		EVENT: map[string]interface{}{
			TYPE:   APP,
			STATUS: []string{},
		},
	}
	eventBodyWithoutEvent = map[string]interface{}{
		URL_KEY: string(TEST_URL),
	}
	invalidBody = string("invalid_body")

	allQuery = map[string][]string{
		GROUP_ID:   []string{"groupid"},
		NODE_ID:    []string{"nodeid"},
		APP_ID:     []string{"appid"},
		IMAGE_NAME: []string{"imagename"},
	}
	node = map[string]interface{}{
		"apps": []string{"appid"},
		ID:     "nodeid",
		"ip":   IP,
		STATUS: "status",
	}
	appSubs = map[string]interface{}{
		ID:       appsubsId,
		TYPE:     "app",
		URL_KEY:  TEST_URL,
		STATUS:   appState,
		EVENT_ID: []string{eventId},
	}
	nodeSubs = map[string]interface{}{
		ID:       nodesubsId,
		TYPE:     "node",
		URL_KEY:  TEST_URL,
		STATUS:   nodeState,
		EVENT_ID: []string{eventId},
	}
	appEvent = map[string]interface{}{
		ID:    eventId,
		SUBS:  []string{appsubsId},
		NODES: []string{"nodeid"},
	}
	lastAppEvent = map[string]interface{}{
		ID:    eventId,
		SUBS:  []string{},
		NODES: []string{"nodeid"},
	}
	nodeEvent = map[string]interface{}{
		ID:   eventId,
		SUBS: []string{nodesubsId},
	}
	lastNodeEvent = map[string]interface{}{
		ID:   eventId,
		SUBS: []string{},
	}
)

var executor Executor

func init() {
	executor = Executor{}
}

func TestCalledRegisterWithInvaildBody_ExpectReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	code, _, err := executor.Register(invalidBody, allQuery)

	if err == nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledRegisterWithBodyWithoutURL_ExpectReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	strBody, _ := convertMapToJson(eventBodyWithoutURL)

	code, _, err := executor.Register(strBody, allQuery)

	if err == nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledRegisterWithBodyWithoutEvent_ExpectReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	strBody, _ := convertMapToJson(eventBodyWithoutEvent)

	code, _, err := executor.Register(strBody, allQuery)

	if err == nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledRegisterWithAppEventBody_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodeSearchExecutorMockObj := nodeSearchmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)
	subsDbMockObj := subsDBmocks.NewMockCommand(ctrl)
	appEventDbMockObj := appEventDBmocks.NewMockCommand(ctrl)

	respStr := []string{`{"message":"valid_message"}`, `{"message":"valid_message"}`}
	expectedUrl := []string{watchUrl, watchUrl}

	nodes := make(map[string]interface{})
	nodes["nodes"] = make([]map[string]interface{}, 2)
	nodes["nodes"].([]map[string]interface{})[0] = node
	nodes["nodes"].([]map[string]interface{})[1] = node

	body, _ := convertMapToJson(validBody)

	gomock.InOrder(
		nodeSearchExecutorMockObj.EXPECT().SearchNodes(allQuery).Return(results.OK, nodes, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl, nil, []byte(body)).Return(respCode, respStr),
		subsDbMockObj.EXPECT().AddSubscriber(appsubsId, APP, TEST_URL, appState, []string{eventId}, allQuery).Return(nil),
		appEventDbMockObj.EXPECT().AddEvent(eventId, appsubsId, nodeIds).Return(nil),
	)

	// pass mockObj to a real object.
	nodeSearchExecutor = nodeSearchExecutorMockObj
	httpExecutor = msgMockObj
	subsDbExecutor = subsDbMockObj
	appEventDbExecutor = appEventDbMockObj

	strBody, _ := convertMapToJson(appEventBody)
	code, _, err := executor.Register(strBody, allQuery)
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledRegisterWithNodeEventBody_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodeSearchExecutorMockObj := nodeSearchmocks.NewMockCommand(ctrl)
	subsDbMockObj := subsDBmocks.NewMockCommand(ctrl)
	nodeEventDbMockObj := nodeEventDBmocks.NewMockCommand(ctrl)

	nodes := make(map[string]interface{})
	nodes["nodes"] = make([]map[string]interface{}, 2)
	nodes["nodes"].([]map[string]interface{})[0] = node
	nodes["nodes"].([]map[string]interface{})[1] = node

	gomock.InOrder(
		nodeSearchExecutorMockObj.EXPECT().SearchNodes(allQuery).Return(results.OK, nodes, nil),
		subsDbMockObj.EXPECT().AddSubscriber(nodesubsId, NODE, TEST_URL, nodeState, nodeIds, allQuery).Return(nil),
		nodeEventDbMockObj.EXPECT().AddEvent(NODE_ID, nodesubsId).Return(nil).AnyTimes(),
	)

	// pass mockObj to a real object.
	nodeSearchExecutor = nodeSearchExecutorMockObj
	subsDbExecutor = subsDbMockObj
	nodeEventDbExecutor = nodeEventDbMockObj

	strBody, _ := convertMapToJson(nodeEventBody)
	code, _, err := executor.Register(strBody, allQuery)
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledUnRegisterAppEvent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subsDbMockObj := subsDBmocks.NewMockCommand(ctrl)
	appEventDbMockObj := appEventDBmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		subsDbMockObj.EXPECT().GetSubscriber(eventId).Return(appSubs, nil),
		appEventDbMockObj.EXPECT().UnRegisterEvent(eventId, appsubsId).Return(nil),
		appEventDbMockObj.EXPECT().GetEvent(eventId).Return(appEvent, nil),
		subsDbMockObj.EXPECT().DeleteSubscriber(eventId).Return(nil),
	)

	// pass mockObj to a real object.
	subsDbExecutor = subsDbMockObj
	appEventDbExecutor = appEventDbMockObj

	code, err := executor.UnRegister(eventId)
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledUnRegisterLastAppEvent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subsDbMockObj := subsDBmocks.NewMockCommand(ctrl)
	appEventDbMockObj := appEventDBmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	reqBody := makeRequestBody(nil, eventId)
	body, _ := convertMapToJson(reqBody)

	gomock.InOrder(
		subsDbMockObj.EXPECT().GetSubscriber(eventId).Return(appSubs, nil),
		appEventDbMockObj.EXPECT().UnRegisterEvent(eventId, appsubsId).Return(nil),
		appEventDbMockObj.EXPECT().GetEvent(eventId).Return(lastAppEvent, nil),
		msgMockObj.EXPECT().SendHttpRequest("DELETE", lastAppEvent[NODES].([]string), nil, []byte(body)),
		appEventDbMockObj.EXPECT().DeleteEvent(eventId).Return(nil),
		subsDbMockObj.EXPECT().DeleteSubscriber(eventId).Return(nil),
	)

	// pass mockObj to a real object.
	subsDbExecutor = subsDbMockObj
	appEventDbExecutor = appEventDbMockObj
	httpExecutor = msgMockObj

	code, err := executor.UnRegister(eventId)
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledUnRegisterNodeEvent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subsDbMockObj := subsDBmocks.NewMockCommand(ctrl)
	nodeEventDbMockObj := nodeEventDBmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		subsDbMockObj.EXPECT().GetSubscriber(eventId).Return(nodeSubs, nil),
		nodeEventDbMockObj.EXPECT().UnRegisterEvent(eventId, nodesubsId).Return(nil),
		nodeEventDbMockObj.EXPECT().GetEvent(eventId).Return(nodeEvent, nil),
		subsDbMockObj.EXPECT().DeleteSubscriber(eventId).Return(nil),
	)

	// pass mockObj to a real object.
	subsDbExecutor = subsDbMockObj
	nodeEventDbExecutor = nodeEventDbMockObj

	code, err := executor.UnRegister(eventId)
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledUnRegisterLastNodeEvent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subsDbMockObj := subsDBmocks.NewMockCommand(ctrl)
	nodeEventDbMockObj := nodeEventDBmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		subsDbMockObj.EXPECT().GetSubscriber(eventId).Return(nodeSubs, nil),
		nodeEventDbMockObj.EXPECT().UnRegisterEvent(eventId, nodesubsId).Return(nil),
		nodeEventDbMockObj.EXPECT().GetEvent(eventId).Return(lastNodeEvent, nil),
		nodeEventDbMockObj.EXPECT().DeleteEvent(eventId).Return(nil),
		subsDbMockObj.EXPECT().DeleteSubscriber(eventId).Return(nil),
	)

	// pass mockObj to a real object.
	subsDbExecutor = subsDbMockObj
	nodeEventDbExecutor = nodeEventDbMockObj

	code, err := executor.UnRegister(eventId)
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledNotificationHandlerWithNodeEvent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventIds := make([]string, 0)
	eventIds = append(eventIds, NODE_ID)
	event := make(map[string]interface{})
	event[ID] = NODE_ID
	event[STATUS] = nodeState[0]

	notification := make(map[string]interface{})
	notification[EVENT_ID] = eventIds
	notification[EVENT] = event

	notiStr, _ := convertMapToJson(notification)

	reqBody := make(map[string]interface{})
	reqBody[EVENT] = event
	body, _ := convertMapToJson(reqBody)

	urls := make([]string, 0)
	urls = append(urls, nodeSubs[URL_KEY].(string))

	subsDbMockObj := subsDBmocks.NewMockCommand(ctrl)
	nodeEventDbMockObj := nodeEventDBmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodeEventDbMockObj.EXPECT().GetEvent(NODE_ID).Return(nodeEvent, nil),
		subsDbMockObj.EXPECT().GetSubscriber(nodesubsId).Return(nodeSubs, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", urls, nil, []byte(body)),
	)

	// pass mockObj to a real object.
	subsDbExecutor = subsDbMockObj
	nodeEventDbExecutor = nodeEventDbMockObj
	httpExecutor = msgMockObj

	executor.NotificationHandler(NODE, notiStr)
}

func TestCalledNotificationHandlerWithAppEvent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventIds := make([]string, 0)
	eventIds = append(eventIds, eventId)
	event := make(map[string]interface{})
	event[ID] = eventId
	event[STATUS] = appState[0]

	notification := make(map[string]interface{})
	notification[EVENT_ID] = eventIds
	notification[EVENT] = event

	notiStr, _ := convertMapToJson(notification)

	reqBody := make(map[string]interface{})
	reqBody[EVENT] = event
	body, _ := convertMapToJson(reqBody)

	urls := make([]string, 0)
	urls = append(urls, appSubs[URL_KEY].(string))

	subsDbMockObj := subsDBmocks.NewMockCommand(ctrl)
	appEventDbMockObj := appEventDBmocks.NewMockCommand(ctrl)
	msgMockObj := msgmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		appEventDbMockObj.EXPECT().GetEvent(eventId).Return(appEvent, nil),
		subsDbMockObj.EXPECT().GetSubscriber(appsubsId).Return(appSubs, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", urls, nil, []byte(body)),
	)

	// pass mockObj to a real object.
	subsDbExecutor = subsDbMockObj
	appEventDbExecutor = appEventDbMockObj
	httpExecutor = msgMockObj

	executor.NotificationHandler(APP, notiStr)
}

