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
	"db/mongo/model/agent"
	"encoding/json"
)

const (
	AGENTS                      = "agents"       // used to indicate a list of agents.
	ID                          = "id"           // used to indicate an agent id.
	HOST                        = "host"         // used to indicate an agent address.
	PORT                        = "port"         // used to indicate an agent port.
	DEFAULT_SDA_PORT            = "48098"        // default service deployment agent port.
)

type AgentController struct{}

var dbManager agent.DBManager

func init() {
	dbManager = agent.DBManager{}
}

// AddAgent inserts a new agent with ip which is passed in call to function.
// If successful, a unique id that is created automatically will be returned.
// otherwise, an appropriate error will be returned.
func (AgentController) AddAgent(ip string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// If body is not empty, try to get agent id from body.
	// This code will be used to update the information of agent without changing id.
	if body != "" {
		bodyMap, err := convertJsonToMap(body)
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}

		// Check whether 'id' is included.
		_, exists := bodyMap["id"]
		if !exists {
			return results.ERROR, nil, errors.InvalidJSON{"id field is required"}
		}

		// Update the information of agent with given ip.
		agentId := bodyMap["id"].(string)
		err = dbManager.UpdateAgentAddress(agentId, ip, DEFAULT_SDA_PORT)
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}

		// Status is updated with 'connected'.
		err = dbManager.UpdateAgentStatus(agentId, STATUS_CONNECTED)
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}

		res := make(map[string]interface{})
		res[ID] = agentId
		return results.OK, res, err
	}

	// Add new agent to database with given ip, port, status.
	agent, err := dbManager.AddAgent(ip, DEFAULT_SDA_PORT, STATUS_CONNECTED)
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
func (AgentController) DeleteAgent(agentId string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Delete agent specified by agentId parameter.
	err = dbManager.DeleteAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	if timers[agentId] != nil {
		timers[agentId] <- true
	}
	delete(timers, agentId)

	return results.OK, err
}

// GetAgent returns the agent with a primary key matching the agentId argument.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (AgentController) GetAgent(agentId string) (int, map[string]interface{}, error) {
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
func (AgentController) GetAgents() (int, map[string]interface{}, error) {
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
