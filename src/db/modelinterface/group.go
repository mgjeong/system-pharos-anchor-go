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

type GroupInterface interface {
	// CreateGroup insert new Group.
	CreateGroup() (map[string]interface{}, error)

	// GetGroup returns single document from db related to group.
	GetGroup(group_id string) (map[string]interface{}, error)

	// GetAllGroups returns all documents from db related to group.
	GetAllGroups() ([]map[string]interface{}, error)

	// GetGroupMembers returns all agents who belong to the target group.
	GetGroupMembers(group_id string) ([]map[string]interface{}, error)

	// GetGroupMembersByAppID returns all agents including specific app on the target group.
	GetGroupMembersByAppID(group_id string, app_id string) ([]map[string]interface{}, error)

	// JoinGroup add specific agent to the target group.
	JoinGroup(group_id string, agent_id string) error

	// LeaveGroup delete specific agent from the target group.
	LeaveGroup(group_id string, agent_id string) error

	// DeleteGroup delete single document from db related to group.
	DeleteGroup(group_id string) error
}
