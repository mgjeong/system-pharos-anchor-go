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

// Package agent provides an interfaces to add, delete, get
// an target edge device.
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
	"time"
)

// Command is an interface of agent operations.
type Command interface {
	RegisterAgent(body string) (int, map[string]interface{}, error)
	UnRegisterAgent(agentId string) (int, error)
	GetAgent(agentId string) (int, map[string]interface{}, error)
	GetAgents() (int, map[string]interface{}, error)
	UpdateAgentStatus(agentId string, status string) error
	Checker
}

const (
	AGENTS                      = "agents"       // used to indicate a list of agents.
	ID                          = "id"           // used to indicate an agent id.
	HOST                        = "host"         // used to indicate an agent address.
	PORT                        = "port"         // used to indicate an agent port.
	STATUS_CONNECTED            = "connected"    // used to update agent status with connected.
	STATUS_DISCONNECTED         = "disconnected" // used to update agent status with disconnected.
	INTERVAL                    = "interval"     // a period between two healthcheck message.
	MAXIMUM_NETWORK_LATENCY_SEC = 3              // the term used to indicate any kind of delay that happens in data communication over a network.
	TIME_UNIT                   = time.Minute    // the minute is a unit of time for healthcheck.
	DEFAULT_AGENT_PORT          = "48098"        // used to indicate a default system-management-agent port.
)

// Executor implements the Command interface.
type Executor struct{}

var dbExecutor agentDB.Command
var httpExecutor messenger.Command

func init() {
	dbExecutor = agentDB.Executor{}
	httpExecutor = messenger.NewExecutor()
}

// AddAgent inserts a new agent with ip which is passed in call to function.
// If successful, a unique id that is created automatically will be returned.
// otherwise, an appropriate error will be returned.
func (Executor) RegisterAgent(body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// If body is not empty, try to get agent id from body.
	// This code will be used to update the information of agent without changing id.
	bodyMap, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Check whether 'ip' is included.
	ip, exists := bodyMap["ip"].(string)
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"ip field is required"}
	}

	// Check whether 'config' is included.
	config, exists := bodyMap["config"]
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"config field is required"}
	}

	// Add new agent to database with given ip, port, status.
	agent, err := dbExecutor.AddAgent(ip, STATUS_CONNECTED, config.(map[string]interface{}))
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[ID] = agent[ID]
	return results.OK, res, err
}

// DeleteAgent deletes the agent with a primary key matching the agentId argument.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UnRegisterAgent(agentId string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent specified by agentId parameter.
	agent, err := dbExecutor.GetAgent(agentId)
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

	// Delete agent specified by agentId parameter.
	err = dbExecutor.DeleteAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	return results.OK, err
}

// GetAgent returns the agent with a primary key matching the agentId argument.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetAgent(agentId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent specified by agentId parameter.
	agent, err := dbExecutor.GetAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return results.OK, agent, err
}

// GetAgents returns all agents in databases as an array.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetAgents() (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get all agents stored in the database.
	agents, err := dbExecutor.GetAllAgents()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[AGENTS] = agents

	return results.OK, res, err
}

// UpdateAgentStatus returns the agent's status.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UpdateAgentStatus(agentId string, status string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get agent specified by agentId parameter.
	err := dbExecutor.UpdateAgentStatus(agentId, status)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return err
	}

	return err
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
