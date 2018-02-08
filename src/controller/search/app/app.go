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
	"commons/logger"
	appMgmt "controller/management/app"
	groupMgmt "controller/management/group"
	nodeMgmt "controller/management/node"
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

var appMgmtExecutor appMgmt.Command
var nodeMgmtExecutor nodeMgmt.Command
var groupMgmtExecutor groupMgmt.Command

func init() {
	appMgmtExecutor = appMgmt.Executor{}
	nodeMgmtExecutor = nodeMgmt.Executor{}
	groupMgmtExecutor = groupMgmt.Executor{}
}

func (Executor) Search(query map[string]interface{}) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	apps := make(map[string]interface{})
	appsInfo := make([]map[string]interface{}, 0)

	if query == nil {
		return appMgmtExecutor.GetApps()
	} else {
		groupId, exists := query[GROUPID]
		if !exists {
			nodeId, exists := query[NODEID]
			if !exists {
				appId, exists := query[APPID]
				if !exists {
					imageName, exists := query[IMAGENAME]
					if !exists {
						//empty
						return appMgmtExecutor.GetApps()
					} else {
						//imageName
						return appMgmtExecutor.GetAppsWithImageName(imageName.([]string)[0])
					}
				} else {
					imageName, exists := query[IMAGENAME]
					if !exists {
						//appId
						return appMgmtExecutor.GetApp(appId.([]string)[0])
					} else {
						//appId, imageName
						result, app, err := appMgmtExecutor.GetApp(appId.([]string)[0])
						if err != nil {
							return result, nil, err
						}
						if searchStringFromSlice(app[IMAGES].([]string), imageName.([]string)[0]) {
							appsInfo = append(appsInfo, app)
						}
						apps[APPS] = removeDuplicateInfo(appsInfo)
						return result, apps, err
					}
				}
			} else {
				appId, exists := query[APPID]
				if !exists {
					imageName, exists := query[IMAGENAME]
					if !exists {
						//nodeId
						result, node, err := nodeMgmtExecutor.GetNode(nodeId.([]string)[0])
						if err != nil {
							return result, nil, err
						}
						for _, id := range node[APPS].([]string) {
							result, app, err := appMgmtExecutor.GetApp(id)
							if err != nil {
								return result, nil, err
							}
							appsInfo = append(appsInfo, app)
						}
						apps[APPS] = removeDuplicateInfo(appsInfo)
						return result, apps, err
					} else {
						//nodeId, imageName
						result, node, err := nodeMgmtExecutor.GetNode(nodeId.([]string)[0])
						if err != nil {
							return result, nil, err
						}
						for _, id := range node[APPS].([]string) {
							result, app, err := appMgmtExecutor.GetApp(id)
							if err != nil {
								return result, nil, err
							}
							if searchStringFromSlice(app[IMAGES].([]string), imageName.([]string)[0]) {
								appsInfo = append(appsInfo, app)
							}
						}
						apps[APPS] = removeDuplicateInfo(appsInfo)
						return result, apps, err
					}
				} else {
					imageName, exists := query[IMAGENAME]
					if !exists {
						//nodeId, appId
						result, node, err := nodeMgmtExecutor.GetNode(nodeId.([]string)[0])
						if err != nil {
							return result, nil, err
						}
						if searchStringFromSlice(node[APPS].([]string), appId.([]string)[0]) {
							result, app, err := appMgmtExecutor.GetApp(appId.([]string)[0])
							if err != nil {
								return result, nil, err
							}
							appsInfo = append(appsInfo, app)
						}
						apps[APPS] = removeDuplicateInfo(appsInfo)
						return result, apps, err
					} else {
						//nodeId, appId, imageName
						result, node, err := nodeMgmtExecutor.GetNode(nodeId.([]string)[0])
						if err != nil {
							return result, nil, err
						}
						if searchStringFromSlice(node[APPS].([]string), appId.([]string)[0]) {
							result, app, err := appMgmtExecutor.GetApp(appId.([]string)[0])
							if err != nil {
								return result, nil, err
							}
							if searchStringFromSlice(app[IMAGES].([]string), imageName.([]string)[0]) {
								appsInfo = append(appsInfo, app)
							}
						}
						apps[APPS] = removeDuplicateInfo(appsInfo)
						return result, apps, err
					}
				}
			}
		} else {
			nodeId, exists := query[NODEID]
			if !exists {
				appId, exists := query[APPID]
				if !exists {
					imageName, exists := query[IMAGENAME]
					if !exists {
						//groupId
						result, group, err := groupMgmtExecutor.GetGroup(groupId.([]string)[0])
						if err != nil {
							return result, nil, err
						}
						for _, id := range group[MEMBERS].([]string) {
							result, node, err := nodeMgmtExecutor.GetNode(id)
							if err != nil {
								return result, nil, err
							}
							for _, id = range node[APPS].([]string) {
								result, app, err := appMgmtExecutor.GetApp(id)
								if err != nil {
									return result, nil, err
								}
								appsInfo = append(appsInfo, app)
							}
						}
						apps[APPS] = removeDuplicateInfo(appsInfo)
						return result, apps, err
					} else {
						//groupId, imageName
						result, group, err := groupMgmtExecutor.GetGroup(groupId.([]string)[0])
						if err != nil {
							return result, nil, err
						}
						for _, id := range group[MEMBERS].([]string) {
							result, node, err := nodeMgmtExecutor.GetNode(id)
							if err != nil {
								return result, nil, err
							}
							for _, id = range node[APPS].([]string) {
								result, app, err := appMgmtExecutor.GetApp(id)
								if err != nil {
									return result, nil, err
								}
								if searchStringFromSlice(app[IMAGES].([]string), imageName.([]string)[0]) {
									appsInfo = append(appsInfo, app)
								}
							}
						}
						apps[APPS] = removeDuplicateInfo(appsInfo)
						return result, apps, err
					}
				} else {
					imageName, exists := query[IMAGENAME]
					if !exists {
						//groupId, appId
						result, group, err := groupMgmtExecutor.GetGroup(groupId.([]string)[0])
						if err != nil {
							return result, nil, err
						}
						for _, id := range group[MEMBERS].([]string) {
							result, node, err := nodeMgmtExecutor.GetNode(id)
							if err != nil {
								return result, nil, err
							}
							if searchStringFromSlice(node[APPS].([]string), appId.([]string)[0]) {
								result, app, err := appMgmtExecutor.GetApp(appId.([]string)[0])
								if err != nil {
									return result, nil, err
								}
								appsInfo = append(appsInfo, app)
							}
						}
						apps[APPS] = removeDuplicateInfo(appsInfo)
						return result, apps, err
					} else {
						//groupId, appId, imageName
						result, group, err := groupMgmtExecutor.GetGroup(groupId.([]string)[0])
						if err != nil {
							return result, nil, err
						}
						for _, id := range group[MEMBERS].([]string) {
							result, node, err := nodeMgmtExecutor.GetNode(id)
							if err != nil {
								return result, nil, err
							}
							if searchStringFromSlice(node[APPS].([]string), appId.([]string)[0]) {
								result, app, err := appMgmtExecutor.GetApp(appId.([]string)[0])
								if err != nil {
									return result, nil, err
								}
								if searchStringFromSlice(app[IMAGES].([]string), imageName.([]string)[0]) {
									appsInfo = append(appsInfo, app)
								}
							}
						}
						apps[APPS] = removeDuplicateInfo(appsInfo)
						return result, apps, err
					}
				}
			} else {
				appId, exists := query[APPID]
				if !exists {
					imageName, exists := query[IMAGENAME]
					if !exists {
						//groupId, nodeId
						result, group, err := groupMgmtExecutor.GetGroup(groupId.([]string)[0])
						if err != nil {
							return result, nil, err
						}
						if searchStringFromSlice(group[MEMBERS].([]string), nodeId.([]string)[0]) {
							result, node, err := nodeMgmtExecutor.GetNode(nodeId.([]string)[0])
							if err != nil {
								return result, nil, err
							}
							for _, id := range node[APPS].([]string) {
								result, app, err := appMgmtExecutor.GetApp(id)
								if err != nil {
									return result, nil, err
								}
								appsInfo = append(appsInfo, app)
							}
						}
						apps[APPS] = removeDuplicateInfo(appsInfo)
						return result, apps, err
					} else {
						//groupId, nodeId, imageName
						result, group, err := groupMgmtExecutor.GetGroup(groupId.([]string)[0])
						if err != nil {
							return result, nil, err
						}
						if searchStringFromSlice(group[MEMBERS].([]string), nodeId.([]string)[0]) {
							result, node, err := nodeMgmtExecutor.GetNode(nodeId.([]string)[0])
							if err != nil {
								return result, nil, err
							}
							for _, id := range node[APPS].([]string) {
								result, app, err := appMgmtExecutor.GetApp(id)
								if err != nil {
									return result, nil, err
								}
								if searchStringFromSlice(app[IMAGES].([]string), imageName.([]string)[0]) {
									appsInfo = append(appsInfo, app)
								}
							}
						}
						apps[APPS] = removeDuplicateInfo(appsInfo)
						return result, apps, err
					}
				} else {
					imageName, exists := query[IMAGENAME]
					if !exists {
						//groupId, nodeId, appId
						result, group, err := groupMgmtExecutor.GetGroup(groupId.([]string)[0])
						if err != nil {
							return result, nil, err
						}
						if searchStringFromSlice(group[MEMBERS].([]string), nodeId.([]string)[0]) {
							result, node, err := nodeMgmtExecutor.GetNode(nodeId.([]string)[0])
							if err != nil {
								return result, nil, err
							}
							if searchStringFromSlice(node[APPS].([]string), appId.([]string)[0]) {
								result, app, err := appMgmtExecutor.GetApp(appId.([]string)[0])
								if err != nil {
									return result, nil, err
								}
								appsInfo = append(appsInfo, app)
							}
						}
						apps[APPS] = removeDuplicateInfo(appsInfo)
						return result, apps, err
					} else {
						//groupId, nodeId, appId, imageName
						result, group, err := groupMgmtExecutor.GetGroup(groupId.([]string)[0])
						if err != nil {
							return result, nil, err
						}
						if searchStringFromSlice(group[MEMBERS].([]string), nodeId.([]string)[0]) {
							result, node, err := nodeMgmtExecutor.GetNode(nodeId.([]string)[0])
							if err != nil {
								return result, nil, err
							}
							if searchStringFromSlice(node[APPS].([]string), appId.([]string)[0]) {
								result, app, err := appMgmtExecutor.GetApp(appId.([]string)[0])
								if err != nil {
									return result, nil, err
								}
								if searchStringFromSlice(app[IMAGES].([]string), imageName.([]string)[0]) {
									appsInfo = append(appsInfo, app)
								}
							}
						}
						apps[APPS] = removeDuplicateInfo(appsInfo)
						return result, apps, err
					}
				}
			}
		}
	}
}

func searchStringFromSlice(slice []string, str string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	return false
}

func removeDuplicateInfo(appsInfo []map[string]interface{}) []map[string]interface{} {
	mapAppsInfo := make(map[string]interface{})
	arrangedAppsInfo := make([]map[string]interface{}, 0)

	for _, appInfo := range appsInfo {
		mapAppsInfo[appInfo[ID].(string)] = appInfo
	}

	for _, value := range arrangedAppsInfo {
		arrangedAppsInfo = append(arrangedAppsInfo, value)
	}

	return arrangedAppsInfo
}
