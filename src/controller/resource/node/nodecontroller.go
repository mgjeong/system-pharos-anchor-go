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
package node

import (
	"bytes"
	"commons/errors"
	"commons/logger"
	"commons/results"
	"commons/url"
	nodeDB "db/mongo/node"
	"encoding/json"
	"messenger"
)

const (
	DEFAULT_AGENT_PORT = "48098" // used to indicate a default system-management-node port.
	IP                 = "ip"
	GET                = "GET"
)

type Command interface {
	GetNodeResourceInfo(nodeId string) (int, map[string]interface{}, error)
	GetAppResourceInfo(nodeId string, appId string) (int, map[string]interface{}, error)
}

type Executor struct{}

var nodeDbExecutor nodeDB.Command

var httpExecutor messenger.Command

func init() {
	nodeDbExecutor = nodeDB.Executor{}
	httpExecutor = messenger.NewExecutor()
}

// GetNodeResourceInfo request an node resource (cpu, mem, disk, network usage) information.
// If response code represents success, returns resource information.
// Otherwise, an appropriate error will be returned.
func (Executor) GetNodeResourceInfo(nodeId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node specified by nodeId parameter.
	node, err := nodeDbExecutor.GetNode(nodeId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getNodeAddress(node)
	urls := makeRequestUrl(address, url.Monitoring(), url.Resource())

	// Request to return node's resource information.
	codes, respStr := httpExecutor.SendHttpRequest(GET, urls, nil)

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// GetAppResourceInfo request an node resource (cpu, mem, net i/o, block i/o) information.
// If response code represents success, returns resource information.
// Otherwise, an appropriate error will be returned.
func (Executor) GetAppResourceInfo(nodeId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node specified by nodeId parameter.
	node, err := nodeDbExecutor.GetNode(nodeId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getNodeAddress(node)
	urls := makeRequestUrl(address, url.Monitoring(), url.Apps(), "/", appId, url.Resource())

	// Request to return node's resource information.
	codes, respStr := httpExecutor.SendHttpRequest(GET, urls, nil)

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
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

// getNodeAddress returns an address as an array.
func getNodeAddress(node map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 1)
	result[0] = map[string]interface{}{
		IP: node[IP],
	}
	return result
}

func makeRequestUrl(address []map[string]interface{}, api_parts ...string) (urls []string) {
	var httpTag string = "http://"
	var full_url bytes.Buffer

	for i := range address {
		full_url.Reset()
		full_url.WriteString(httpTag + address[i][IP].(string) +
			":" + DEFAULT_AGENT_PORT + url.Base())
		for _, api_part := range api_parts {
			full_url.WriteString(api_part)
		}
		urls = append(urls, full_url.String())
	}
	return urls
}

// convertRespToMap converts a response in the form of JSON data into a map.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func convertRespToMap(respStr []string) (map[string]interface{}, error) {
	resp, err := convertJsonToMap(respStr[0])
	if err != nil {
		logger.Logging(logger.ERROR, "Failed to convert response from string to map")
		return nil, errors.InternalServerError{"Json Converting Failed"}
	}
	return resp, err
}