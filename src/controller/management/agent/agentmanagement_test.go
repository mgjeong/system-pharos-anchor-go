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
package agent

import (
	"commons/errors"
	"commons/results"
	dbmocks "db/modelinterface/mocks"
	"github.com/golang/mock/gomock"
	"reflect"
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
	body             = `{"description":"description"}`
	respCode         = []int{results.OK}
	errorRespCode    = []int{results.ERROR}
	respStr          = []string{`{"response":"response"}`}
	invalidRespStr   = []string{`{"invalidJson"}`}
	notFoundError    = errors.NotFound{}
	connectionError  = errors.DBConnectionError{}
	invalidJsonError = errors.InvalidJSON{}
)

var manager AgentInterface

func init() {
	manager = AgentManager{}
}

func TestCalledAddAgentWithValidBody_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := `{"ip":"127.0.0.1", "config":{"key":"value"}}`
	expectedRes := map[string]interface{}{
		"id": "000000000000000000000001",
	}

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().AddAgent(ip, status, configuration).Return(agent, nil),
	)
	// pass mockObj to a real object.
	dbManager = dbManagerMockObj

	code, res, err := manager.AddAgent(body)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(res, expectedRes) {
		t.Error()
	}
}

func TestCalledAddAgentWithInValidJsonFormatBody_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidBody := `{"ip"}`

	code, _, err := manager.AddAgent(invalidBody)

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

func TestCalledAddAgentWithInvalidBodyNotIncludingIPField_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidBody := `{"key":"value"}`

	code, _, err := manager.AddAgent(invalidBody)

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

func TestCalledAddAgentWhenFailedToInsertNewAgentToDB_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().AddAgent(ip, status, configuration).Return(nil, notFoundError),
	)

	// pass mockObj to a real object.
	dbManager = dbManagerMockObj

	body := `{"ip":"127.0.0.1", "config":{"key":"value"}}`
	code, _, err := manager.AddAgent(body)

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

func TestCalledDeleteAgent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().DeleteAgent(agentId).Return(nil),
	)

	// pass mockObj to a real object.
	dbManager = dbManagerMockObj

	code, err := manager.DeleteAgent(agentId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledDeleteAgentWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().DeleteAgent(agentId).Return(notFoundError),
	)

	// pass mockObj to a real object.
	dbManager = dbManagerMockObj

	code, err := manager.DeleteAgent(agentId)

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

func TestCalledGetAgent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(agent, nil),
	)

	// pass mockObj to a real object.
	dbManager = dbManagerMockObj

	code, res, err := manager.GetAgent(agentId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(res, agent) {
		t.Error()
	}
}

func TestCalledGetAgentWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(nil, notFoundError),
	)

	// pass mockObj to a real object.
	dbManager = dbManagerMockObj

	code, _, err := manager.GetAgent(agentId)

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

func TestCalledGetAgents_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	agents := []map[string]interface{}{agent}

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAllAgents().Return(agents, nil),
	)

	// pass mockObj to a real object.
	dbManager = dbManagerMockObj

	code, res, err := manager.GetAgents()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(res["agents"].([]map[string]interface{}), agents) {
		t.Error()
	}
}

func TestCalledGetAgentsWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().GetAllAgents().Return(nil, notFoundError),
	)

	// pass mockObj to a real object.
	dbManager = dbManagerMockObj

	code, _, err := manager.GetAgents()

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

func TestCalledUpdateAgentStatus_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().UpdateAgentStatus(agentId, status).Return(nil),
	)

	// pass mockObj to a real object.
	dbManager = dbManagerMockObj

	err := manager.UpdateAgentStatus(agentId, status)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledUpdateAgentStatusWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManagerMockObj := dbmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		dbManagerMockObj.EXPECT().UpdateAgentStatus(agentId, status).Return(notFoundError),
	)

	// pass mockObj to a real object.
	dbManager = dbManagerMockObj

	err := manager.UpdateAgentStatus(agentId, status)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}
