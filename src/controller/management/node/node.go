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

// Package node provides an interfaces to add, delete, get
// an target edge device.
package node

import (
	"bytes"
	"commons/errors"
	"commons/logger"
	"commons/results"
	"commons/url"
	noti "controller/notification"
	nodeDB "db/mongo/node"
	"encoding/json"
	"messenger"
	"strings"
	"time"
)

// Command is an interface of node operations.
type Command interface {
	RegisterNode(body string) (int, map[string]interface{}, error)
	UnRegisterNode(nodeId string) (int, error)
	GetNode(nodeId string) (int, map[string]interface{}, error)
	GetNodes() (int, map[string]interface{}, error)
	GetNodesWithAppID(appId string) (int, map[string]interface{}, error)
	UpdateNodeStatus(nodeId string, status string) error
	PingNode(nodeId string, body string) (int, error)
	GetNodeConfiguration(nodeId string) (int, map[string]interface{}, error)
	SetNodeConfiguration(nodeId string, body string) (int, error)
}

const (
	NODE                        = "node"         // used to indicate a node.
	NODES                       = "nodes"        // used to indicate a list of nodes.
	ID                          = "id"           // used to indicate an node id.
	APPS                        = "apps"         // used to indicate a list of apps.
	EVENT                       = "event"        // used to indicate an event.
	EVENT_ID                    = "eventid"      // used to indicate an event id.
	HOST                        = "host"         // used to indicate an node address.
	PORT                        = "port"         // used to indicate an node port.
	STATUS                      = "status"       // used to update node status.
	STATUS_CONNECTED            = "connected"    // used to update node status with connected.
	STATUS_DISCONNECTED         = "disconnected" // used to update node status with disconnected.
	INTERVAL                    = "interval"     // a period between two healthcheck message.
	MAXIMUM_NETWORK_LATENCY_SEC = 3             // the term used to indicate any kind of delay that happens in data communication over a network.
	TIME_UNIT                   = time.Minute    // the minute is a unit of time for healthcheck.
	DEFAULT_NODE_PORT           = "48098"        // used to indicate a default pharos node port.
	PROPERTIES                  = "properties"
)

// Executor implements the Command interface.
type Executor struct{}

var nodeDbExecutor nodeDB.Command
var httpExecutor messenger.Command
var notiExecutor noti.Command

func init() {
	nodeDbExecutor = nodeDB.Executor{}
	httpExecutor = messenger.NewExecutor()
	notiExecutor = noti.Executor{}
}

// RegisterNode inserts a new node with ip which is passed in call to function.
// If successful, a unique id that is created automatically will be returned.
// otherwise, an appropriate error will be returned.
func (Executor) RegisterNode(body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// If body is not empty, try to get node id from body.
	// This code will be used to update the information of node without changing id.
	bodyMap, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Check whether 'ip' is included.
	ip, exists := bodyMap["ip"].(string)
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"ip field is required"}
	}

	// Check whether 'config' is included.
	config, exists := bodyMap["config"]
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"config field is required"}
	}

	// check whether device already exists.
	node, err := nodeDbExecutor.GetNodeByIP(ip)
	if err == nil {
		logger.Logging(logger.DEBUG, "This device is already registered")
		res := make(map[string]interface{})
		res[ID] = node[ID]
		return results.OK, res, err
	}

	// Add new node to database with given ip, port, status.
	node, err = nodeDbExecutor.AddNode(ip, STATUS_CONNECTED, config.(map[string]interface{}))
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[ID] = node[ID]
	return results.OK, res, err
}

// UnRegisterNode deletes the node with a primary key matching the nodeId argument.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UnRegisterNode(nodeId string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node specified by nodeId parameter.
	node, err := nodeDbExecutor.GetNode(nodeId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	address, err := getNodeAddress(node)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	urls := makeRequestUrl(address, url.Management(), url.Unregister())
	httpExecutor.SendHttpRequest("POST", urls, nil)

	// Stop timer and close the channel for ping.
	if common.timers[nodeId] != nil {
		common.timers[nodeId] <- true
		close(common.timers[nodeId])
	}
	delete(common.timers, nodeId)

	// Delete node specified by nodeId parameter.
	err = nodeDbExecutor.DeleteNode(nodeId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	return results.OK, err
}

// GetNode returns the node with a primary key matching the nodeId argument.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetNode(nodeId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node specified by nodeId parameter.
	node, err := nodeDbExecutor.GetNode(nodeId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Copy from the original map to the target map, except "config" field.
	res := make(map[string]interface{})
	for key, value := range node {
		if strings.Compare(key, "config") != 0 {
			res[key] = value
		}
	}

	return results.OK, res, err
}

// GetNodes returns all nodes in databases as an array.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetNodes() (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get all nodes stored in the database.
	nodes, err := nodeDbExecutor.GetNodes()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	copiedNodes := make([]map[string]interface{}, 0)
	for _, node := range nodes {
		// Copy from the original map to the target map, except "config" field.
		res := make(map[string]interface{})
		for key, value := range node {
			if strings.Compare(key, "config") != 0 {
				res[key] = value
			}
		}
		copiedNodes = append(copiedNodes, res)
	}

	res := make(map[string]interface{})
	res[NODES] = copiedNodes

	return results.OK, res, err
}

func (Executor) GetNodesWithAppID(appID string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	query := make(map[string]interface{})
	query[APPS] = appID

	// Get matched nodes with query stored in the database.
	nodes, err := nodeDbExecutor.GetNodes(query)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	copiedNodes := make([]map[string]interface{}, 0)
	for _, node := range nodes {
		// Copy from the original map to the target map, except "config" field.
		res := make(map[string]interface{})
		for key, value := range node {
			if strings.Compare(key, "config") != 0 {
				res[key] = value
			}
		}
		copiedNodes = append(copiedNodes, res)
	}

	res := make(map[string]interface{})
	res[NODES] = copiedNodes

	return results.OK, res, err
}

// UpdateNodeStatus returns the node's status.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UpdateNodeStatus(nodeId string, status string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node specified by nodeId parameter.
	err := nodeDbExecutor.UpdateNodeStatus(nodeId, status)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return err
	}

	return err
}

func (Executor) GetNodeConfiguration(nodeId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get matched nodes with query stored in the database.
	node, err := nodeDbExecutor.GetNode(nodeId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := node["config"].(map[string]interface{})
	return results.OK, res, err
}

func (Executor) SetNodeConfiguration(nodeId string, body string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node specified by nodeId parameter.
	node, err := nodeDbExecutor.GetNode(nodeId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	address, err := getNodeAddress(node)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	urls := makeRequestUrl(address, url.Management(), url.Device(), url.Configuration())

	codes, _ := httpExecutor.SendHttpRequest("POST", urls, nil, []byte(body))

	result := codes[0]
	if !isSuccessCode(result) {
		return results.ERROR, err
	}

	// Update configuration information.
	updatedProps, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	originProps := node["config"].(map[string]interface{})["properties"]
	for _, originProp := range originProps.([]interface{}) {
		for originKey, _ := range originProp.(map[string]interface{}) {
			for _, updatedProp := range updatedProps["properties"].([]interface{}) {
				for updatedKey, updatedValue := range updatedProp.(map[string]interface{}) {
					if strings.Compare(originKey, updatedKey) == 0 {
						originProp.(map[string]interface{})[originKey] = updatedValue
					}
				}
			}
		}
	}

	err = nodeDbExecutor.UpdateNodeConfiguration(nodeId, node["config"].(map[string]interface{}))
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	return results.OK, nil
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

// isSuccessCode returns true in case of success and false otherwise.
func isSuccessCode(code int) bool {
	if code >= 200 && code <= 299 {
		return true
	}
	return false
}

// makeRequestUrl make a list of urls that can be used to send a http request.
func makeRequestUrl(address []map[string]interface{}, api_parts ...string) (urls []string) {
	var httpTag string = "http://"
	var full_url bytes.Buffer

	for i := range address {
		full_url.Reset()
		full_url.WriteString(httpTag + address[i]["ip"].(string) +
			":" + DEFAULT_NODE_PORT + url.Base())
		for _, api_part := range api_parts {
			full_url.WriteString(api_part)
		}
		urls = append(urls, full_url.String())
	}
	return urls
}

// getNodeAddress returns an address as an array.
func getNodeAddress(node map[string]interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 1)

	_, exists := node["ip"]
	if !exists {
		return nil, errors.InvalidJSON{"ip field is required"}
	}

	result[0] = map[string]interface{}{
		"ip": node["ip"],
	}
	return result, nil
}
