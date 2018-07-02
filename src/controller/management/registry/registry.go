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

package registry

import (
	"commons/errors"
	"commons/logger"
	"commons/results"
	"commons/url"
	"commons/util"
	appmanager "controller/management/app"
	nodemanager "controller/management/node"
	"db/mongo/registry"
	"messenger"
)

type Command interface {
	// AddDockerRegistry add docker registry to database.
	AddDockerRegistry(body string) (int, map[string]interface{}, error)
	DeleteDockerRegistry(registryId string) (int, error)
	GetDockerRegistries() (int, map[string]interface{}, error)
	DockerRegistryEventHandler(body string) (int, error)
}

const (
	ID          = "id"
	IP          = "ip"
	POST        = "POST"
	APPS        = "apps"
	NODES       = "nodes" // used to indicate a list of nodes.
	HOST        = "host"
	REPOSITORY  = "repository"
	TARGETINFO  = "target"
	REQUESTINFO = "request"
	REGISTRIES  = "registries"
	REGISTRY    = "registry"
	EVENTS      = "events"
)

type Executor struct{}

var appmanagementExecutor appmanager.Command
var nodemanagementExecutor nodemanager.Command
var registryDbExecutor registry.Command
var httpExecutor messenger.Command

func init() {
	appmanagementExecutor = appmanager.Executor{}
	nodemanagementExecutor = nodemanager.Executor{}
	registryDbExecutor = registry.Executor{}
	httpExecutor = messenger.NewExecutor()
}

func (Executor) AddDockerRegistry(body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	reqBody, err := util.ConvertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Check whether 'ip' is included.
	ip, exists := reqBody[IP].(string)
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"ip field is required"}
	}

	registry, err := registryDbExecutor.AddDockerRegistry(ip)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[ID] = registry[ID]

	return results.OK, res, err
}

func (Executor) DeleteDockerRegistry(registryId string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Delete registry specified by registryId parameter.
	err := registryDbExecutor.DeleteDockerRegistry(registryId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	return results.OK, err
}

func (Executor) GetDockerRegistries() (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get all of registries list.
	registries, err := registryDbExecutor.GetDockerRegistries()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[REGISTRIES] = registries

	return results.OK, res, err
}

func (Executor) DockerRegistryEventHandler(body string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	convertedBody, err := util.ConvertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}
	events := convertedBody[EVENTS]

	for _, eventInfo := range events.([]interface{}) {
		parsedEvent := make(map[string]interface{})
		parsedEvent, err = parseEventInfo(eventInfo.(map[string]interface{}))
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, err
		}

		imageName := parsedEvent[HOST].(string) + "/" + parsedEvent[REPOSITORY].(string)
		_, apps, err := appmanagementExecutor.GetAppsWithImageName(imageName)
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, err
		}

		for _, app := range apps[APPS].([]map[string]interface{}) {
			appId := app[ID]
			_, nodes, err := nodemanagementExecutor.GetNodesWithAppID(appId.(string))
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, err
			}
			address := getMemberAddress(nodes[NODES].([]map[string]interface{}))
			urls := util.MakeRequestUrl(address, url.Management(), url.Apps(), "/", appId.(string), url.Events())
			_, _ = httpExecutor.SendHttpRequest(POST, urls, nil, []byte(body))
		}
	}
	return results.OK, nil
}

// getNodeAddress returns an member's address as an array.
func getMemberAddress(members []map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(members))
	for i, node := range members {
		result[i] = map[string]interface{}{
			"ip":    node["ip"],
			"config": node["config"],
		}
	}
	return result
}

// parseEventInfo parse data which is matched image-info on DB from event-notification.
func parseEventInfo(eventInfo map[string]interface{}) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	targetInfoEvent := make(map[string]interface{})
	requestInfoEvent := make(map[string]interface{})
	parsedEvent := make(map[string]interface{})

	targetInfoEvent = eventInfo[TARGETINFO].(map[string]interface{})
	requestInfoEvent = eventInfo[REQUESTINFO].(map[string]interface{})
	parsedEvent[HOST] = requestInfoEvent[HOST]
	parsedEvent[REPOSITORY] = targetInfoEvent[REPOSITORY]

	return parsedEvent, nil
}
