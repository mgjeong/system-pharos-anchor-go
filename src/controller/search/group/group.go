/*******************************************************************************
 * Copyright 2018 Samsung Electronics All Rights Reserved.
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

import (
	"commons/errors"
	"commons/logger"
	"commons/results"
	appmanager "controller/management/app"
	groupmanager "controller/management/group"
	nodemanager "controller/management/node"
)

type Command interface {
	SearchGroups(query map[string]interface{}) (int, map[string]interface{}, error)
}

const (
	GROUPID   = "groupId"
	NODEID    = "nodeId"
	APPID     = "appId"
	IMAGENAME = "imageName"
)

type Executor struct{}

var appmanagementExecutor appmanager.Command
var nodemanagementExecutor nodemanager.Command
var groupmanagementExecutor groupmanager.Command

func init() {
	appmanagementExecutor = appmanager.Executor{}
	nodemanagementExecutor = nodemanager.Executor{}
	groupmanagementExecutor = groupmanager.Executor{}
}

func (Executor) SearchGroups(query map[string]interface{}) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Validate query parameters.
	err := checkQueryParam(query)
	if err != nil {
		logger.Logging(logger.DEBUG, err.Error())
		return results.ERROR, nil, err
	}

	// Checks for 'groupId' query parameter existence.
	groupList := make([]map[string]interface{}, 0)
	if groupId, exists := query[GROUPID]; exists {
		groupList = filterByGroupId(groupId.([]string)[0])
	} else {
		groupList = filterByGroupId()
	}

	// Checks for 'imageName' query parameter existence.
	if imageName, exists := query[IMAGENAME]; exists {
		groupList = filterByImageName(groupList, imageName.([]string)[0])
	}

	// Checks for 'appId' query parameter existence.
	if appId, exists := query[APPID]; exists {
		groupList = filterByAppId(groupList, appId.([]string)[0])
	}

	// Checks for 'nodeId' query parameter existence.
	if nodeId, exists := query[NODEID]; exists {
		groupList = filterByNodeId(groupList, nodeId.([]string)[0])
	}

	res := make(map[string]interface{})
	res["groups"] = groupList
	return results.OK, res, nil
}

func filterByImageName(groups []map[string]interface{}, imageName string) []map[string]interface{} {
	filteredGroups := make([]map[string]interface{}, 0)

	for _, group := range groups {
		members := group["members"].([]string)
		for _, member := range members {
			_, node, err := nodemanagementExecutor.GetNode(member)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return nil
			}

			for _, appId := range node["apps"].([]string) {
				_, app, err := appmanagementExecutor.GetApp(appId)
				if err != nil {
					logger.Logging(logger.ERROR, err.Error())
					return nil
				}

				if searchStringFromSlice(app["images"].([]string), imageName) {
					filteredGroups = append(filteredGroups, group)
					break
				}
			}
		}
	}
	return filteredGroups
}

func filterByAppId(groups []map[string]interface{}, appId string) []map[string]interface{} {
	filteredGroups := make([]map[string]interface{}, 0)

	for _, group := range groups {
		members := group["members"].([]string)
		for _, member := range members {
			_, node, err := nodemanagementExecutor.GetNode(member)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return nil
			}

			if searchStringFromSlice(node["apps"].([]string), appId) {
				filteredGroups = append(filteredGroups, group)
				break
			}
		}
	}
	return filteredGroups
}

func filterByNodeId(groups []map[string]interface{}, nodeId string) []map[string]interface{} {
	filteredGroups := make([]map[string]interface{}, 0)

	for _, group := range groups {
		if searchStringFromSlice(group["members"].([]string), nodeId) {
			filteredGroups = append(filteredGroups, group)
		}
	}
	return filteredGroups
}

func filterByGroupId(groupId ...string) []map[string]interface{} {
	filteredGroups := make([]map[string]interface{}, 0)

	switch len(groupId) {
	case 1:
		_, group, err := groupmanagementExecutor.GetGroup(groupId[0])
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return nil
		}
		filteredGroups = append(filteredGroups, group)
	case 0:
		_, groups, err := groupmanagementExecutor.GetGroups()
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return nil
		}
		filteredGroups = groups["groups"].([]map[string]interface{})
	}
	return filteredGroups
}

func checkQueryParam(query map[string]interface{}) error {
	supportedQueries := []string{GROUPID, NODEID, APPID, IMAGENAME}

	for key, _ := range query {
		if !searchStringFromSlice(supportedQueries, key) {
			return errors.NotFoundURL{Message: "not supported query parameter"}
		}
	}
	return nil
}

func searchStringFromSlice(slice []string, str string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	return false
}
