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

// Package group provides operations to manage edge device group (e.g., create, join, leave, delete...).
package group

import (
	"commons/errors"
	"commons/logger"
	"commons/results"
	"db/mongo/model/group"
	"encoding/json"
)

const (
	AGENTS        = "agents"      // used to indicate a list of agents.
	GROUPS        = "groups"      // used to indicate a list of groups.
)

type GroupController struct{}

var dbManager group.DBManager

func init() {
	dbManager = group.DBManager{}
}

// CreateGroup inserts a new group to databases.
// This function returns a unique id in case of success and an error otherwise.
func (GroupController) CreateGroup() (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	group, err := dbManager.CreateGroup()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return results.OK, group, err
}

// GetGroup returns the information of the group specified by groupId parameter.
// If response code represents success, returns information about the group.
// Otherwise, an appropriate error will be returned.
func (GroupController) GetGroup(groupId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	group, err := dbManager.GetGroup(groupId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return results.OK, group, err
}

// GetGroups returns a list of groups that is created on databases.
// If response code represents success, returns a list of groups.
// Otherwise, an appropriate error will be returned.
func (GroupController) GetGroups() (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	groups, err := dbManager.GetAllGroups()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[GROUPS] = groups

	return results.OK, res, err
}

// JoinGroup adds the agent to a list of members.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (GroupController) JoinGroup(groupId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	bodyMap, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Check whether 'agents' is included.
	_, exists := bodyMap[AGENTS]
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"agents field is required"}
	}

	for _, agentId := range bodyMap[AGENTS].([]interface{}) {
		err = dbManager.JoinGroup(groupId, agentId.(string))
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
	}

	return results.OK, nil, err
}

// LeaveGroup removes the agent from a list of members.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (GroupController) LeaveGroup(groupId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	bodyMap, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Check whether 'agents' is included.
	_, exists := bodyMap[AGENTS]
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"agents field is required"}
	}

	for _, agentId := range bodyMap[AGENTS].([]interface{}) {
		err = dbManager.LeaveGroup(groupId, agentId.(string))
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
	}

	return results.OK, nil, err
}

// DeleteGroup deletes the group with a primary key matching the groupId argument.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (GroupController) DeleteGroup(groupId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	err = dbManager.DeleteGroup(groupId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return results.OK, nil, err
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
