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

package health

import (
	"commons/errors"
	"commons/results"
	agentmocks "controller/management/agent/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestCalledPingAgentWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	agentManagerMockObj := agentmocks.NewMockAgentInterface(ctrl)

	gomock.InOrder(
		agentManagerMockObj.EXPECT().GetAgent(agentId).Return(results.ERROR, nil, notFoundError),
	)
	// pass mockObj to a real object.
	common.agentManager = agentManagerMockObj

	healthExecutor := Executor{}
	code, err := healthExecutor.PingAgent(agentId, "")

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
	common.agentManager = agentManagerMockObj

	healthExecutor := Executor{}
	code, err := healthExecutor.PingAgent(agentId, invalidKeyBody)

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
	common.agentManager = agentManagerMockObj

	healthExecutor := Executor{}
	code, err := healthExecutor.PingAgent(agentId, invalidValueBody)

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
