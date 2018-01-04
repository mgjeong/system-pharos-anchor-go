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
	"bytes"
	"commons/errors"
	"commons/logger"
	"commons/results"
	"commons/url"
	"encoding/json"
	"messenger"
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

type (
	// Registerer is an interface of agent registration operation.
	Registerer interface {
		RegisterAgent(body string) (int, map[string]interface{}, error)
	}

	// Unregisterer is an interface of agent un-registration operation.
	Unregisterer interface {
		UnRegisterAgent(agentId string) (int, error)
	}

	// Command is an interface of health operations.
	Command interface {
		Registerer
		Unregisterer
		Checker
	}

	// Executor implements the Command interface.
	Executor struct{}
)

var httpExecutor messenger.Command

func init() {
	httpExecutor = messenger.NewExecutor()
}

// RegisterAgent inserts a new agent with ip which is passed in call to function.
// If successful, a unique id that is created automatically will be returned.
// otherwise, an appropriate error will be returned.
func (Executor) RegisterAgent(body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	result, res, err := common.agentManager.AddAgent(body)
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
func (Executor) UnRegisterAgent(agentId string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent specified by agentId parameter.
	_, agent, err := common.agentManager.GetAgent(agentId)
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

	codes, _ := httpExecutor.SendHttpRequest("POST", urls)

	result := codes[0]
	if !isSuccessCode(result) {
		return results.ERROR, err
	}

	// Stop timer and close the channel for ping.
	if common.timers[agentId] != nil {
		common.timers[agentId] <- true
		close(common.timers[agentId])
	}
	delete(common.timers, agentId)

	// Delete agent
	result, err = common.agentManager.DeleteAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

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
