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

package registration

import (
	"commons/errors"
	"commons/results"
	agentmocks "controller/management/agent/mocks"
	"github.com/golang/mock/gomock"
	msgmocks "messenger/mocks"
	"testing"
)

const (
	status  = "connected"
	appId   = "000000000000000000000000"
	agentId = "000000000000000000000001"
	ip      = "127.0.0.1"
	port    = "48098"
)

var (
	agent = map[string]interface{}{
		"id":     agentId,
		"ip":     ip,
		"apps":   []string{},
		"config": configuration,
	}
	address = []map[string]interface{}{
		map[string]interface{}{
			"ip": ip,
		}}
	configuration = map[string]interface{}{
		"key": "value",
	}
	invalidKeyBody   = `{"key":"value"}`
	invalidValueBody = `{"interval":"value"}`
	respCode         = []int{results.OK}
	errorRespCode    = []int{results.ERROR}
	respStr          = []string{`{"response":"response"}`}
	invalidRespStr   = []string{`{"invalidJson"}`}
	notFoundError    = errors.NotFound{}
	dbOperationError = errors.DBOperationError{}
	invalidJsonError = errors.InvalidJSON{}
)

var registrator RegistrationInterface

func init() {
	registrator = AgentRegistrator{}
}

func TestCalledRegisterAgentWithValidBody_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := `{"ip":"127.0.0.1"}`
	agentManagerMockObj := agentmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		agentManagerMockObj.EXPECT().AddAgent(body).Return(results.OK, nil, nil),
	)
	// pass mockObj to a real object.
	agentManager = agentManagerMockObj

	code, _, err := registrator.RegisterAgent(body)

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledRegisterAgentWhenDBHasDuplicatedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := `{"ip":"127.0.0.1"}`
	agentManagerMockObj := agentmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		agentManagerMockObj.EXPECT().AddAgent(body).Return(results.ERROR, nil, dbOperationError),
	)
	// pass mockObj to a real object.
	agentManager = agentManagerMockObj

	code, _, err := registrator.RegisterAgent(body)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBOperationError", err.Error())
	case errors.DBOperationError:
	}
}

func TestCalledUnRegisterAgentWithValidBody_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUrl := []string{"http://" + ip + ":" + port + "/api/v1/unregister"}

	agentManagerMockObj := agentmocks.NewMockAgentInterface(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		agentManagerMockObj.EXPECT().GetAgent(agentId).Return(results.OK, agent, nil),
		msgMockObj.EXPECT().SendHttpRequest("POST", expectedUrl).Return(respCode, respStr),
		agentManagerMockObj.EXPECT().DeleteAgent(agentId).Return(results.OK, nil),
	)
	// pass mockObj to a real object.
	agentManager = agentManagerMockObj
	httpRequester = msgMockObj

	code, err := registrator.UnRegisterAgent(agentId)

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledUnRegisterAgentWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	agentManagerMockObj := agentmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		agentManagerMockObj.EXPECT().GetAgent(agentId).Return(results.ERROR, nil, notFoundError),
	)
	// pass mockObj to a real object.
	agentManager = agentManagerMockObj

	code, err := registrator.UnRegisterAgent(agentId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledPingAgentWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	agentManagerMockObj := agentmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		agentManagerMockObj.EXPECT().GetAgent(agentId).Return(results.ERROR, nil, notFoundError),
	)
	// pass mockObj to a real object.
	agentManager = agentManagerMockObj

	code, err := registrator.PingAgent(agentId, "")

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledPingAgentWithInvalidBody_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	agentManagerMockObj := agentmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		agentManagerMockObj.EXPECT().GetAgent(agentId).Return(results.OK, nil, nil),
	)
	// pass mockObj to a real object.
	agentManager = agentManagerMockObj

	code, err := registrator.PingAgent(agentId, invalidKeyBody)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", err.Error())
	case errors.InvalidJSON:
	}
}

func TestCalledPingAgentWithInvalidValueBody_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	agentManagerMockObj := agentmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		agentManagerMockObj.EXPECT().GetAgent(agentId).Return(results.OK, nil, nil),
	)
	// pass mockObj to a real object.
	agentManager = agentManagerMockObj

	code, err := registrator.PingAgent(agentId, invalidValueBody)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", err.Error())
	case errors.InvalidJSON:
	}
}
