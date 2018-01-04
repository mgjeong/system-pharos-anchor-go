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
	"bytes"
	"commons/errors"
	"commons/logger"
	"commons/results"
	"commons/url"
	"controller/management/agent"
	"encoding/json"
	"messenger"
	"strconv"
	"time"
)

const (
	STATUS_CONNECTED            = "connected"    // used to update agent status with connected.
	STATUS_DISCONNECTED         = "disconnected" // used to update agent status with disconnected.
	INTERVAL                    = "interval"     // a period between two healthcheck message.
	MAXIMUM_NETWORK_LATENCY_SEC = 3              // the term used to indicate any kind of delay that happens in data communication over a network.
	TIME_UNIT                   = time.Minute    // the minute is a unit of time for healthcheck.
	DEFAULT_AGENT_PORT          = "48098"        // used to indicate a default system-management-agent port.
)

type AgentRegistrator struct{}

var agentManager agent.Command
var httpRequester messenger.Command
var timers map[string]chan bool

func init() {
	agentManager = agent.Executor{}
	timers = make(map[string]chan bool)
	httpRequester = messenger.NewMessenger()
}

// RegisterAgent inserts a new agent with ip which is passed in call to function.
// If successful, a unique id that is created automatically will be returned.
// otherwise, an appropriate error will be returned.
func (AgentRegistrator) RegisterAgent(body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	result, res, err := agentManager.AddAgent(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, res, err
}

// UnregisterAgent send unregister request to target-agent.
// And then stop the ping process and delete agent in db.
// IF successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (AgentRegistrator) UnRegisterAgent(agentId string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent specified by agentId parameter.
	_, agent, err := agentManager.GetAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	address, err := getAgentAddress(agent)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	urls := makeRequestUrl(address, url.Unregister())

	codes, _ := httpRequester.SendHttpRequest("POST", urls)

	result := codes[0]
	if !isSuccessCode(result) {
		return results.ERROR, err
	}

	// Stop timer and close the channel for ping.
	if timers[agentId] != nil {
		timers[agentId] <- true
		close(timers[agentId])
	}
	delete(timers, agentId)

	// Delete agent
	result, err = agentManager.DeleteAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	return results.OK, err
}

// PingAgent starts timer with received interval.
// If agent does not send next healthcheck message in interval time,
// change the status of device from connected to disconnected.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (AgentRegistrator) PingAgent(agentId string, body string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent specified by agentId parameter.
	_, _, err := agentManager.GetAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	bodyMap, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	// Check whether 'interval' is included.
	_, exists := bodyMap[INTERVAL]
	if !exists {
		return results.ERROR, errors.InvalidJSON{"interval field is required"}
	}

	interval, err := strconv.Atoi(bodyMap[INTERVAL].(string))
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, errors.InvalidJSON{"invalid value type(interval must be integer)"}
	}

	_, exists = timers[agentId]
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

// isSuccessCode returns true in case of success and false otherwise.
func isSuccessCode(code int) bool {
	if code >= 200 && code <= 299 {
		return true
	}
	return false
}

// makeRequestUrl make a list of urls that can be used to send a http request.
func makeRequestUrl(address []map[string]interface{}, api_parts ...string) (urls []string) {
	var httpTag string = "http://"
	var full_url bytes.Buffer

	for i := range address {
		full_url.Reset()
		full_url.WriteString(httpTag + address[i]["ip"].(string) +
			":" + DEFAULT_AGENT_PORT + url.Base())
		for _, api_part := range api_parts {
			full_url.WriteString(api_part)
		}
		urls = append(urls, full_url.String())
	}
	return urls
}

// getAgentAddress returns an address as an array.
func getAgentAddress(agent map[string]interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 1)

	_, exists := agent["ip"]
	if !exists {
		return nil, errors.InvalidJSON{"ip field is required"}
	}

	result[0] = map[string]interface{}{
		"ip": agent["ip"],
	}
	return result, nil
}
