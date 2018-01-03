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
	"bytes"
	"commons/errors"
	"commons/logger"
	"commons/results"
	"commons/url"
	agentDB "db/mongo/agent"
	"encoding/json"
	"messenger"
)

const (
	DEFAULT_AGENT_PORT = "48098" // used to indicate a default system-management-agent port.
)

type agentController struct{}

var agentDbManager agentDB.Command
var AgentController agentController

var httpRequester messenger.Command

func init() {
	agentDbManager = agentDB.Executor{}
	httpRequester = messenger.NewMessenger()
}

// GetResourceInfo request an agent resource (os, processor, performance) information.
// If response code represents success, returns resource information.
// Otherwise, an appropriate error will be returned.
func (agentController) GetResourceInfo(agentId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent specified by agentId parameter.
	agent, err := agentDbManager.GetAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getAgentAddress(agent)
	urls := makeRequestUrl(address, url.Resource())

	// Request to return agent's resource information.
	codes, respStr := httpRequester.SendHttpRequest("GET", urls)

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	
	return result, respMap, err
}

// GetPerformanceInfo request an agent performance(cpu, disk, mem usage) information.
// If response code represents success, returns performance information.
// Otherwise, an appropriate error will be returned.
func (agentController) GetPerformanceInfo(agentId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent specified by agentId parameter.
	agent, err := agentDbManager.GetAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getAgentAddress(agent)
	urls := makeRequestUrl(address, url.Resource(), url.Performance())

	// Request to return agent's performance information.
	codes, respStr := httpRequester.SendHttpRequest("GET", urls)

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	
	return result, respMap, err
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

// getAgentAddress returns an address as an array.
func getAgentAddress(agent map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 1)
	result[0] = map[string]interface{}{
		"ip": agent["ip"],
	}
	return result
}

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

// convertRespToMap converts a response in the form of JSON data into a map.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func convertRespToMap(respStr []string) (map[string]interface{}, error) {
	resp, err := convertJsonToMap(respStr[0])
	if err != nil {
		logger.Logging(logger.ERROR, "Failed to convert response from string to map")
		return nil, errors.InternalServerError{"Json Converting Failed"}
	}
	return resp, err
}
