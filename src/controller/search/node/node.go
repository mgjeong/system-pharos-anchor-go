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
package node

import (
	"commons/errors"
	"commons/logger"
	"commons/results"
	app "controller/management/app"
	group "controller/management/group"
	node "controller/management/node"
)

type Command interface {
	SearchNodes(map[string][]string) (int, map[string]interface{}, error)
}

const (
	GROUP_ID   string = "groupId"
	NODE_ID    string = "nodeId"
	APP_ID     string = "appId"
	IMAGE_NAME string = "imageName"
	NODES      string = "nodes"
	GROUPS     string = "groups"
	apps       string = "apps"
)

type Executor struct{}

var nodeExecutor node.Command
var groupExecutor group.Command
var appExecutor app.Command

func init() {
	nodeExecutor = node.Executor{}
	groupExecutor = group.Executor{}
	appExecutor = app.Executor{}
}

func (Executor) SearchNodes(query map[string][]string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	if doesContainInvalidQuery(query) {
		logger.Logging(logger.DEBUG, "Url contains invalid query")
		return results.ERROR, nil, errors.InvalidParam{"Url contains invalid query"}
	}

	_, nodeList, err := nodeExecutor.GetNodes()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	nodes := nodeList[NODES].([]map[string]interface{})

	if nodeIds, ok := query[NODE_ID]; ok {
		for _, node := range nodes {
			nodes = nodes[:0]
			if node["id"] == nodeIds[0] {
				nodes = append(nodes, node)
				break
			}
		}
	}

	if appIds, ok := query[APP_ID]; ok {
		nodes = filterByAppId(nodes, appIds[0])
	}

	if groupIds, ok := query[GROUP_ID]; ok {
		nodes, err = filterByGroupId(nodes, groupIds[0])
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, errors.Unknown{"failed to filter by group ID"}
		}
	}

	if imageNames, ok := query[IMAGE_NAME]; ok {
		nodes, err = filterByImageName(nodes, imageNames[0])
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, errors.Unknown{"failed to filter image name"}
		}
	}

	result := make(map[string]interface{}, 0)
	result["nodes"] = nodes
	return results.OK, result, nil
}

func filterByGroupId(nodes []map[string]interface{}, groupId string) ([]map[string]interface{}, error) {
	filteredNodes := make([]map[string]interface{}, 0)

	_, groupList, err := groupExecutor.GetGroups()
	if err != nil {
		logger.Logging(logger.DEBUG, err.Error())
		return nil, err
	}

	groups := groupList[GROUPS].([]map[string]interface{})
	for _, group := range groups {
		if group["id"] == groupId {
			members := group["members"]
			for _, node := range nodes {
				if doesContain(members.([]string), node["id"].(string)) {
					filteredNodes = append(filteredNodes, node)
				}
			}
		}
	}

	return filteredNodes, nil
}

func filterByAppId(nodes []map[string]interface{}, appId string) []map[string]interface{} {
	filteredNodes := make([]map[string]interface{}, 0)

	for _, node := range nodes {
		if doesContain(node["apps"].([]string), appId) {
			filteredNodes = append(filteredNodes, node)
		}
	}

	return filteredNodes
}

func filterByImageName(nodes []map[string]interface{}, imageName string) ([]map[string]interface{}, error) {
	filteredNodes := make([]map[string]interface{}, 0)

	for _, node := range nodes {
		nodeApps := node["apps"].([]string)
		for _, nodeApp := range nodeApps {
			_, app, err := appExecutor.GetApp(nodeApp)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return nil, err
			}
			if doesContain(app["images"].([]string), imageName) {
				filteredNodes = append(filteredNodes, node)
				break
			}
		}
	}
	return filteredNodes, nil
}

func doesContainInvalidQuery(query map[string][]string) bool {
	queryCnt := 0
	if _, ok := query[NODE_ID]; ok {
		queryCnt++
	}
	if _, ok := query[GROUP_ID]; ok {
		queryCnt++
	}
	if _, ok := query[IMAGE_NAME]; ok {
		queryCnt++
	}
	if _, ok := query[APP_ID]; ok {
		queryCnt++
	}
	if len(query) == queryCnt {
		return false
	} else {
		return true
	}
}

func doesContain(arr []string, item string) bool {
	for _, val := range arr {
		if val == item {
			return true
		}
	}
	return false
}