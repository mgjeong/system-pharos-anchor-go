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

package group

type GroupInterface interface {
	// CreateGroup inserts a new group to databases.
	CreateGroup() (int, map[string]interface{}, error)

	// GetGroup returns the information of the group specified by groupId parameter.
	GetGroup(groupId string) (int, map[string]interface{}, error)

	// GetGroups returns a list of groups that is created on databases.
	GetGroups() (int, map[string]interface{}, error)

	// JoinGroup adds the agent to a list of members.
	JoinGroup(groupId string, body string) (int, map[string]interface{}, error)

	// LeaveGroup removes the agent from a list of members.
	LeaveGroup(groupId string, body string) (int, map[string]interface{}, error)

	// DeleteGroup deletes the group with a primary key matching the groupId argument.
	DeleteGroup(groupId string) (int, map[string]interface{}, error)
}

