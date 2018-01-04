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

// Package agent provides an interfaces to deploy, update, start, stop, delete
// an application to target edge device.
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

type Executor struct{}

var agentDbExecutor agentDB.Command
var httpExecutor messenger.Command

func init() {
	agentDbExecutor = agentDB.Executor{}
  httpExecutor = messenger.NewExecutor()
}

// Command is an interface of agent deployment operations.
type Command interface {
	// DeployApp request an deployment of edge services to an agent specified by
	// agentId parameter.
	DeployApp(agentId string, body string) (int, map[string]interface{}, error)

	// GetApps request a list of applications that is deployed to an agent specified
	// by agentId parameter.
	GetApps(agentId string) (int, map[string]interface{}, error)

	// GetApp gets the application's information of the agent specified by agentId parameter.
	GetApp(agentId string, appId string) (int, map[string]interface{}, error)

	// UpdateApp request to update an application specified by appId parameter.
	UpdateAppInfo(agentId string, appId string, body string) (int, map[string]interface{}, error)

	// DeleteApp request to delete an application specified by appId parameter.
	DeleteApp(agentId string, appId string) (int, map[string]interface{}, error)

	// UpdateAppInfo request to update all of images which is included an application
	// specified by appId parameter.
	UpdateApp(agentId string, appId string) (int, map[string]interface{}, error)

	// StartApp request to start an application specified by appId parameter.
	StartApp(agentId string, appId string) (int, map[string]interface{}, error)

	// StopApp request to stop an application specified by appId parameter.
	StopApp(agentId string, appId string) (int, map[string]interface{}, error)
}

// DeployApp request an deployment of edge services to an agent specified by agentId parameter.
// If response code represents success, add an app id to a list of installed app and returns it.
// Otherwise, an appropriate error will be returned.
func (Executor) DeployApp(agentId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent specified by agentId parameter.
	agent, err := agentDbExecutor.GetAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getAgentAddress(agent)
	urls := makeRequestUrl(address, url.Deploy())

	// Request an deployment of edge services to a specific agent.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls, []byte(body))

	// Convert the received response from string to map.
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// if response code represents success, insert the installed appId into agentDbExecutor.
	result := codes[0]
	if isSuccessCode(result) {
		err = agentDbExecutor.AddAppToAgent(agentId, respMap["id"].(string))
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
	}

	return result, respMap, err
}

// GetApps request a list of applications that is deployed to an agent
// specified by agentId parameter.
// If response code represents success, returns a list of applications.
// Otherwise, an appropriate error will be returned.
func (Executor) GetApps(agentId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent specified by agentId parameter.
	agent, err := agentDbExecutor.GetAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getAgentAddress(agent)
	urls := makeRequestUrl(address, url.Apps())

	// Request list of applications that is deployed to agent.
	codes, respStr := httpExecutor.SendHttpRequest("GET", urls)

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// GetApp gets the application's information of the agent specified by agentId parameter.
// If response code represents success, returns information of application.
// Otherwise, an appropriate error will be returned.
func (Executor) GetApp(agentId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent including app specified by appId parameter.
	agent, err := agentDbExecutor.GetAgentByAppID(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getAgentAddress(agent)
	urls := makeRequestUrl(address, url.Apps(), "/", appId)

	// Request get target application's information
	codes, respStr := httpExecutor.SendHttpRequest("GET", urls)

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// UpdateApp request to update an application specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UpdateAppInfo(agentId string, appId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent including app specified by appId parameter.
	agent, err := agentDbExecutor.GetAgentByAppID(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getAgentAddress(agent)
	urls := makeRequestUrl(address, url.Apps(), "/", appId)

	// Request update target application's information.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls, []byte(body))

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// DeleteApp request to delete an application specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) DeleteApp(agentId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent including app specified by appId parameter.
	agent, err := agentDbExecutor.GetAgentByAppID(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getAgentAddress(agent)
	urls := makeRequestUrl(address, url.Apps(), "/", appId)

	// Request delete target application
	codes, respStr := httpExecutor.SendHttpRequest("DELETE", urls)

	// Convert the received response from string to map.
	result := codes[0]
	if !isSuccessCode(result) {
		respMap, err := convertRespToMap(respStr)
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
		return result, respMap, err
	}

	// if response code represents success, delete the appId from agentDbExecutor.
	err = agentDbExecutor.DeleteAppFromAgent(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, nil, err
}

// UpdateAppInfo request to update all of images which is included an application
// specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UpdateApp(agentId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent including app specified by appId parameter.
	agent, err := agentDbExecutor.GetAgentByAppID(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getAgentAddress(agent)
	urls := makeRequestUrl(address, url.Apps(), "/", appId, url.Update())

	// Request checking and updating all of images which is included target.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls)

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// StartApp request to start an application specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) StartApp(agentId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent including app specified by appId parameter.
	agent, err := agentDbExecutor.GetAgentByAppID(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getAgentAddress(agent)
	urls := makeRequestUrl(address, url.Apps(), "/", appId, url.Start())

	// Request start target application.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls)

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// StopApp request to stop an application specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) StopApp(agentId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent including app specified by appId parameter.
	agent, err := agentDbExecutor.GetAgentByAppID(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getAgentAddress(agent)
	urls := makeRequestUrl(address, url.Apps(), "/", appId, url.Stop())

	// Request stop target application.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls)

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