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

package app

import (
	"commons/errors"
	"commons/logger"
	"commons/results"
	appDB "db/mongo/app"
	nodeDB "db/mongo/node"
	groupDB "db/mongo/group"
)

const (
	ID        string = "id"
	GROUPID   string = "groupId"
	NODEID    string = "nodeId"
	APPID     string = "appId"
	APPS      string = "apps"
	IMAGENAME string = "imageName"
	IMAGES    string = "images"
	MEMBERS   string = "members"
)

// Command is an interface of apps operations.
type Command interface {
	Search(query map[string]interface{}) (int, map[string]interface{}, error)
}

// Executor implements the Command interface.
type Executor struct{}

var appDbExecutor appDB.Command
var nodeDbExecutor nodeDB.Command
var groupDbExecutor groupDB.Command

func init() {
	appDbExecutor = appDB.Executor{}
	nodeDbExecutor = nodeDB.Executor{}
	groupDbExecutor = groupDB.Executor{}
}

func (Executor) Search(query map[string]interface{}) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Validate query parameters.
	err := checkQueryParam(query)
	if err != nil {
		logger.Logging(logger.DEBUG, err.Error())
		return results.ERROR, nil, err
	}

	// Checks for 'appId' query parameter existence.
	appList := make([]map[string]interface{}, 0)
	if appId, exists := query[APPID]; exists {
		appList = filterByAppId(appId.([]string)[0])
	} else {
		appList = filterByAppId()
	}

	// Checks for 'imageName' query parameter existence.
	if imageName, exists := query[IMAGENAME]; exists {
		appList = filterByImageName(appList, imageName.([]string)[0])
	}

	// Checks for 'nodeId' query parameter existence.
	if nodeId, exists := query[NODEID]; exists {
		appList = filterByNodeId(appList, nodeId.([]string)[0])
	}

	// Checks for 'gourId' query parameter existence.
	if groupId, exists := query[GROUPID]; exists {
		appList = filterByGroupId(appList, groupId.([]string)[0])
	}

	res := make(map[string]interface{})
	res[APPS] = appList
	return results.OK, res, nil
}

func filterByGroupId(apps []map[string]interface{}, groupId string) []map[string]interface{} {
	filteredApps := make([]map[string]interface{}, 0)

	group, err := groupDbExecutor.GetGroup(groupId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return nil
	}

	members := group[MEMBERS].([]string)
	for _, member := range members {
		node, err := nodeDbExecutor.GetNode(member)
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return nil
		}

		for _, app := range apps {
			if searchStringFromSlice(node[APPS].([]string), app[ID].(string)) {
				filteredApps = append(filteredApps, app)
				break
			}
		}
	}

	return filteredApps
}

func filterByNodeId(apps []map[string]interface{}, nodeId string) []map[string]interface{} {
	filteredApps := make([]map[string]interface{}, 0)

	node, err := nodeDbExecutor.GetNode(nodeId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return nil
	}

	for _, app := range apps {
		if searchStringFromSlice(node[APPS].([]string), app[ID].(string)) {
			filteredApps = append(filteredApps, app)
			break
		}
	}
	return filteredApps
}

func filterByImageName(apps []map[string]interface{}, imageName string) []map[string]interface{} {
	filteredApps := make([]map[string]interface{}, 0)

	for _, app := range apps {
		if searchStringFromSlice(app[IMAGES].([]string), imageName) {
			filteredApps = append(filteredApps, app)
			break
		}
	}
	return filteredApps
}

func filterByAppId(appId ...string) []map[string]interface{} {
	filteredApps := make([]map[string]interface{}, 0)

	switch len(appId) {
	case 1:
		app, err := appDbExecutor.GetApp(appId[0])
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return nil
		}
		filteredApps = append(filteredApps, app)
	case 0:
		apps, err := appDbExecutor.GetApps()
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return nil
		}
		filteredApps = apps
	}
	return filteredApps
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
