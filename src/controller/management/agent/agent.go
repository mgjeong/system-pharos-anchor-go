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
	"commons/errors"
	"commons/logger"
	"commons/results"
	"db/modelinterface"
	agentDB "db/mongo/model/agent"
	"encoding/json"
)

type Command interface {
	AddAgent(body string) (int, map[string]interface{}, error)
	DeleteAgent(agentId string) (int, error)
	GetAgent(agentId string) (int, map[string]interface{}, error)
	GetAgents() (int, map[string]interface{}, error)
	UpdateAgentStatus(agentId string, status string) error
}

const (
	AGENTS           = "agents"    // used to indicate a list of agents.
	ID               = "id"        // used to indicate an agent id.
	HOST             = "host"      // used to indicate an agent address.
	PORT             = "port"      // used to indicate an agent port.
	STATUS_CONNECTED = "connected" // used to update agent status with connected.
)

type Executor struct{}

var dbManager modelinterface.AgentInterface

func init() {
	dbManager = agentDB.DBManager{}
}

// AddAgent inserts a new agent with ip which is passed in call to function.
// If successful, a unique id that is created automatically will be returned.
// otherwise, an appropriate error will be returned.
func (Executor) AddAgent(body string) (int, map[string]interface{}, error) {
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
	agent, err := dbManager.AddAgent(ip, STATUS_CONNECTED, config.(map[string]interface{}))
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
func (Executor) DeleteAgent(agentId string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Delete agent specified by agentId parameter.
	err := dbManager.DeleteAgent(agentId)
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
	agent, err := dbManager.GetAgent(agentId)
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
	agents, err := dbManager.GetAllAgents()
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
	err := dbManager.UpdateAgentStatus(agentId, status)
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
