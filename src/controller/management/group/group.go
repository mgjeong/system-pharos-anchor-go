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
	"commons/util"
	groupDB "db/mongo/group"
	nodeDB "db/mongo/node"
)

type Command interface {
	// CreateGroup inserts a new group to databases.
	CreateGroup(body string) (int, map[string]interface{}, error)

	// GetGroup returns the information of the group specified by groupId parameter.
	GetGroup(groupId string) (int, map[string]interface{}, error)

	// GetGroups returns a list of groups that is created on databases.
	GetGroups() (int, map[string]interface{}, error)

	// JoinGroup adds the node to a list of members.
	JoinGroup(groupId string, body string) (int, map[string]interface{}, error)

	// LeaveGroup removes the node from a list of members.
	LeaveGroup(groupId string, body string) (int, map[string]interface{}, error)

	// DeleteGroup deletes the group with a primary key matching the groupId argument.
	DeleteGroup(groupId string) (int, map[string]interface{}, error)
}

const (
	AGENTS     = "nodes"  // used to indicate a list of nodes.
	GROUPS     = "groups" // used to indicate a list of groups.
	GROUP_NAME = "name"   // used to indicate a group name.
)

type Executor struct{}

var groupDbExecutor groupDB.Command
var nodeDbExecutor nodeDB.Command

func init() {
	groupDbExecutor = groupDB.Executor{}
	nodeDbExecutor = nodeDB.Executor{}
}

// CreateGroup inserts a new group to databases.
// This function returns a unique id in case of success and an error otherwise.
func (Executor) CreateGroup(body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	bodyMap, err := util.ConvertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Check whether 'name' is included.
	_, exists := bodyMap[GROUP_NAME]
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"name field is required"}
	}

	name := bodyMap[GROUP_NAME].(string)
	group, err := groupDbExecutor.CreateGroup(name)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return results.OK, group, err
}

// GetGroup returns the information of the group specified by groupId parameter.
// If response code represents success, returns information about the group.
// Otherwise, an appropriate error will be returned.
func (Executor) GetGroup(groupId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	group, err := groupDbExecutor.GetGroup(groupId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return results.OK, group, err
}

// GetGroups returns a list of groups that is created on databases.
// If response code represents success, returns a list of groups.
// Otherwise, an appropriate error will be returned.
func (Executor) GetGroups() (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	groups, err := groupDbExecutor.GetGroups()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[GROUPS] = groups

	return results.OK, res, err
}

// JoinGroup adds the node to a list of members.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) JoinGroup(groupId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	bodyMap, err := util.ConvertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Check whether 'nodes' is included.
	_, exists := bodyMap[AGENTS]
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"nodes field is required"}
	}

	// Validate nodeIds in request body.
	for _, nodeId := range bodyMap[AGENTS].([]interface{}) {
		_, err := nodeDbExecutor.GetNode(nodeId.(string))
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
	}

	for _, nodeId := range bodyMap[AGENTS].([]interface{}) {
		err = groupDbExecutor.JoinGroup(groupId, nodeId.(string))
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
	}

	return results.OK, nil, err
}

// LeaveGroup removes the node from a list of members.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) LeaveGroup(groupId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	bodyMap, err := util.ConvertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Check whether 'nodes' is included.
	_, exists := bodyMap[AGENTS]
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"nodes field is required"}
	}

	for _, nodeId := range bodyMap[AGENTS].([]interface{}) {
		err = groupDbExecutor.LeaveGroup(groupId, nodeId.(string))
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
func (Executor) DeleteGroup(groupId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	err := groupDbExecutor.DeleteGroup(groupId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return results.OK, nil, err
}
