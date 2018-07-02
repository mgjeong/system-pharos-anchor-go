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
	"commons/errors"
	"commons/logger"
	"commons/results"
	"commons/url"
	"commons/util"
	noti "controller/notification"
	groupSearch "controller/search/group"
	groupDB "db/mongo/group"
	nodeDB "db/mongo/node"
	"github.com/satori/go.uuid"
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
	Reboot(nodeId string) (int, error)
	Restore(nodeId string) (int, error)
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
	MAXIMUM_NETWORK_LATENCY_SEC = 3              // the term used to indicate any kind of delay that happens in data communication over a network.
	TIME_UNIT                   = time.Minute    // the minute is a unit of time for healthcheck.
	PROPERTIES                  = "properties"
)

// Executor implements the Command interface.
type Executor struct{}

var nodeDbExecutor nodeDB.Command
var groupDbExecutor groupDB.Command
var httpExecutor messenger.Command
var notiExecutor noti.Command
var groupSearchExecutor groupSearch.Command

func init() {
	nodeDbExecutor = nodeDB.Executor{}
	groupDbExecutor = groupDB.Executor{}
	httpExecutor = messenger.NewExecutor()
	notiExecutor = noti.Executor{}
	groupSearchExecutor = groupSearch.Executor{}
}

// RegisterNode inserts a new node with ip which is passed in call to function.
// If successful, a unique id that is created automatically will be returned.
// otherwise, an appropriate error will be returned.
func (Executor) RegisterNode(body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// If body is not empty, try to get node id from body.
	// This code will be used to update the information of node without changing id.
	bodyMap, err := util.ConvertJsonToMap(body)
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
	config, exists := bodyMap["config"].(map[string]interface{})
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"config field is required"}
	}

	// Check whether 'apps' is included.
	appIds := make([]string, 0)
	if apps, exists := bodyMap["apps"]; exists {
		for _, app := range apps.([]interface{}) {
			appIds = append(appIds, app.(string))
		}
	}

	// check whether deviceId already exists.
	var deviceId string
	for _, prop := range config["properties"].([]interface{}) {
		if value, exists := prop.(map[string]interface{})["deviceid"]; exists {
			deviceId = value.(string)
		}
	}

	// Generate a unique deviceId.
	for len(deviceId) == 0 {
		uuid, err := generateUUIDv4()
		if err != nil {
			return results.ERROR, nil, err
		}

		_, err = nodeDbExecutor.GetNode(uuid)
		if err != nil {
			switch err.(type) {
			default:
				return results.ERROR, nil, err
			case errors.NotFound:
				deviceId = uuid
			}
		}
	}

	// Add new node to database with given ip, port, status.
	node, err := nodeDbExecutor.AddNode(deviceId, ip, STATUS_CONNECTED, config, appIds)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Send notification to subscribers.
	go func() {
		notiExecutor.UpdateSubscriber()
		sendNotification(node[ID].(string), STATUS_CONNECTED)
	}()

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

	urls := util.MakeRequestUrl(address, url.Management(), url.Unregister())
	httpExecutor.SendHttpRequest("POST", urls, nil)

	// Stop timer and close the channel for ping.
	common.Lock()
	if common.timers[nodeId] != nil {
		common.timers[nodeId] <- true
		close(common.timers[nodeId])
	}
	delete(common.timers, nodeId)
	common.Unlock()

	// Delete node specified by nodeId parameter.
	err = nodeDbExecutor.DeleteNode(nodeId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	// Remove the node from a list of group members.
	query := make(map[string]interface{})
	query["nodeId"] = []string{nodeId}

	_, groups, _ := groupSearchExecutor.SearchGroups(query)
	for _, group := range groups["groups"].([]map[string]interface{}) {
		groupDbExecutor.LeaveGroup(group["id"].(string), nodeId)
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

// Reboot reboots the device with nodeId.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) Reboot(nodeId string) (int, error) {
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

	urls := util.MakeRequestUrl(address, url.Management(), url.Device(), url.Reboot())
	httpExecutor.SendHttpRequest("POST", urls, nil)

	return results.OK, err
}

// Restore restore the device with nodeId to initial state.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) Restore(nodeId string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get matched nodes with query stored in the database.
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

	urls := util.MakeRequestUrl(address, url.Management(), url.Device(), url.Restore())
	httpExecutor.SendHttpRequest("POST", urls, nil)

	return results.OK, err
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

	urls := util.MakeRequestUrl(address, url.Management(), url.Device(), url.Configuration())

	codes, _ := httpExecutor.SendHttpRequest("POST", urls, nil, []byte(body))

	result := codes[0]
	if !util.IsSuccessCode(result) {
		return results.ERROR, err
	}

	// Update configuration information.
	updatedProps, err := util.ConvertJsonToMap(body)
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

// generateUUIDv4 generates a random UUID.
func generateUUIDv4() (string, error) {
	ret, err := uuid.NewV4()
	if err != nil {
		return "", errors.IOError{"generating uuid is fail"}
	}
	return ret.String(), nil
}
