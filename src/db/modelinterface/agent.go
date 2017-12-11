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
package modelinterface

type AgentInterface interface {
	// AddAgent insert new Agent.
	AddAgent(host string, port string, status string) (map[string]interface{}, error)

	// UpdateAgentAddress updates ip,port of agent from db related to agent.
	UpdateAgentAddress(agent_id string, host string, port string) error

	// UpdateAgentStatus updates status of agent from db related to agent.
	UpdateAgentStatus(agent_id string, status string) error

	// GetAgent returns single document from db related to agent.
	GetAgent(agent_id string) (map[string]interface{}, error)

	// GetAllAgents returns all documents from db related to agent.
	GetAllAgents() ([]map[string]interface{}, error)

	// GetAgentByAppID returns single document including specific app.
	GetAgentByAppID(agent_id string, app_id string) (map[string]interface{}, error)

	// AddAppToAgent add specific app to the target agent.
	AddAppToAgent(agent_id string, app_id string) error

	// DeleteAppFromAgent delete specific app from the target agent.
	DeleteAppFromAgent(agent_id string, app_id string) error

	// DeleteAgent delete single document from db related to agent.
	DeleteAgent(agent_id string) error
}
