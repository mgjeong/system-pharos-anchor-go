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
	"commons/logger"
	"commons/results"
	"controller/management/agent"
	"encoding/json"
	"strconv"
	"time"
)

const (
	STATUS_CONNECTED            = "connected"    // used to update agent status with connected.
	STATUS_DISCONNECTED         = "disconnected" // used to update agent status with disconnected.
	INTERVAL                    = "interval"     // a period between two healthcheck message.
	MAXIMUM_NETWORK_LATENCY_SEC = 3              // the term used to indicate any kind of delay that happens in data communication over a network.
	TIME_UNIT                   = time.Minute    // the minute is a unit of time for healthcheck.
)

type AgentRegistrator struct{}

var agentManager agent.AgentManager
var timers map[string]chan bool

func init() {
	agentManager = agent.AgentManager{}
	timers = make(map[string]chan bool)
}

// RegisterAgent inserts a new agent with ip which is passed in call to function.
// If successful, a unique id that is created automatically will be returned.
// otherwise, an appropriate error will be returned.
func (AgentRegistrator) RegisterAgent(ip string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	result, res, err := agentManager.AddAgent(ip, body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, res, err
}

// PingAgent starts timer with received interval.
// If agent does not send next healthcheck message in interval time,
// change the status of device from connected to disconnected.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (AgentRegistrator) PingAgent(agentId string, ip string, body string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent specified by agentId parameter.
	_, agent, err := agentManager.GetAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	storedIP := agent["host"].(string)
	if ip != storedIP {
		logger.Logging(logger.ERROR, "address does not match")
		return results.ERROR, errors.NotFound{"address does not match, try registration again"}
	}

	bodyMap, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	_, exists := timers[agentId]
	if !exists {
		logger.Logging(logger.DEBUG, "first ping request is received from agent")
	} else {
		if timers[agentId] != nil {
			// If ping request is received in interval time, send signal to stop timer.
			timers[agentId] <- true
			logger.Logging(logger.DEBUG, "ping request is received in interval time")
		} else {
			logger.Logging(logger.DEBUG, "ping request is received after interval time-out")
			err = agentManager.UpdateAgentStatus(agentId, STATUS_CONNECTED)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
			}
		}
	}

	// Start timer with received interval time.
	interval, err := strconv.Atoi(bodyMap[INTERVAL].(string))
	timeDurationMin := time.Duration(interval+MAXIMUM_NETWORK_LATENCY_SEC) * TIME_UNIT
	timer := time.NewTimer(timeDurationMin)
	go func() {
		quit := make(chan bool)
		timers[agentId] = quit

		select {
		// Block until timer finishes.
		case <-timer.C:
			logger.Logging(logger.ERROR, "ping request is not received in interval time")

			// Status is updated with 'disconnected'.
			err = agentManager.UpdateAgentStatus(agentId, STATUS_DISCONNECTED)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
			}

		case <-quit:
			timer.Stop()
			return
		}

		timers[agentId] = nil
		close(quit)
	}()

	return results.OK, err
}

// convertJsonToMap converts JSON data into a map.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func convertJsonToMap(jsonStr string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, errors.InvalidJSON{"Unmarshalling Failed"}
	}
	return result, err
}
