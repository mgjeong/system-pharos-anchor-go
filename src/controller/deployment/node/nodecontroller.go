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

// Package node provides an interfaces to deploy, update, start, stop, delete
// an application to target edge device.
package node

import (
	"commons/errors"
	"commons/logger"
	"commons/results"
	"commons/url"
	"commons/util"
	noti "controller/notification"
	appDB "db/mongo/app"
	appEventDB "db/mongo/event/app"
	subsDB "db/mongo/event/subscriber"
	nodeDB "db/mongo/node"
	"math/rand"
	"messenger"
	"time"
)

const (
	LETTERBYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	EVENTID     = "eventId"
	EVENT       = "event"
	PULLED      = "pulled"
	CREATED     = "created"
	STARTED     = "started"
	APP         = "app"
)

type Executor struct{}

var appDbExecutor appDB.Command
var nodeDbExecutor nodeDB.Command
var appEventDbExecutor appEventDB.Command
var subsDbExecutor subsDB.Command
var httpExecutor messenger.Command
var notiExecutor noti.Command

func init() {
	rand.Seed(time.Now().UnixNano())

	appDbExecutor = appDB.Executor{}
	nodeDbExecutor = nodeDB.Executor{}
	httpExecutor = messenger.NewExecutor()
	appEventDbExecutor = appEventDB.Executor{}
	subsDbExecutor = subsDB.Executor{}
	notiExecutor = noti.Executor{}
}

// Command is an interface of node deployment operations.
type Command interface {
	// DeployApp request an deployment of edge services to an node specified by
	// nodeId parameter.
	DeployApp(nodeId string, body string, query map[string]interface{}) (int, map[string]interface{}, error)

	// GetApps request a list of applications that is deployed to an node specified
	// by nodeId parameter.
	GetApps(nodeId string) (int, map[string]interface{}, error)

	// GetApp gets the application's information of the node specified by nodeId parameter.
	GetApp(nodeId string, appId string) (int, map[string]interface{}, error)

	// UpdateApp request to update an application specified by appId parameter.
	UpdateAppInfo(nodeId string, appId string, body string) (int, map[string]interface{}, error)

	// DeleteApp request to delete an application specified by appId parameter.
	DeleteApp(nodeId string, appId string) (int, map[string]interface{}, error)

	// UpdateAppInfo request to update all of images which is included an application
	// specified by appId parameter.
	UpdateApp(nodeId string, appId string, query map[string]interface{}) (int, map[string]interface{}, error)

	// StartApp request to start an application specified by appId parameter.
	StartApp(nodeId string, appId string) (int, map[string]interface{}, error)

	// StopApp request to stop an application specified by appId parameter.
	StopApp(nodeId string, appId string) (int, map[string]interface{}, error)
}

// DeployApp request an deployment of edge services to an node specified by nodeId parameter.
// If response code represents success, add an app id to a list of installed app and returns it.
// Otherwise, an appropriate error will be returned.
func (Executor) DeployApp(nodeId string, body string, query map[string]interface{}) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node specified by nodeId parameter.
	node, err := nodeDbExecutor.GetNode(nodeId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getNodeAddress(node)
	urls := util.MakeRequestUrl(address, url.Management(), url.Apps(), url.Deploy())

	codes := make([]int, 0)
	respStr := make([]string, 0)
	if eventUrl, exists := query[EVENT]; exists {
		eventId := generateRandStringBytes(39)
		subsId := generateRandStringBytes(39)

		err = subsDbExecutor.AddSubscriber(subsId, APP, eventUrl.([]string)[0],
			[]string{PULLED, CREATED, STARTED}, []string{eventId}, make(map[string][]string))
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
		err = appEventDbExecutor.AddEvent(eventId, subsId, []string{nodeId})
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			subsDbExecutor.DeleteSubscriber(subsId)
			return results.ERROR, nil, err
		}

		eventIDQuery := make(map[string]interface{})
		eventIDQuery[EVENTID] = []string{eventId}

		// Request an deployment of edge services to a specific node.
		codes, respStr = httpExecutor.SendHttpRequest("POST", urls, eventIDQuery, []byte(body))

		err = subsDbExecutor.DeleteSubscriber(subsId)
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
		err = appEventDbExecutor.DeleteEvent(eventId)
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
	} else {
		// Request an deployment of edge services to a specific node.
		codes, respStr = httpExecutor.SendHttpRequest("POST", urls, nil, []byte(body))
	}

	// Convert the received response from string to map.
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// if response code represents success, insert the installed appId into nodeDbExecutor.
	result := codes[0]
	if util.IsSuccessCode(result) {
		err = appDbExecutor.AddApp(respMap["id"].(string), []byte(respMap["description"].(string)))
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
		err = nodeDbExecutor.AddAppToNode(nodeId, respMap["id"].(string))
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
	}

	notiExecutor.UpdateSubscriber()

	return result, respMap, err
}

// GetApps request a list of applications that is deployed to an node
// specified by nodeId parameter.
// If response code represents success, returns a list of applications.
// Otherwise, an appropriate error will be returned.
func (Executor) GetApps(nodeId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node specified by nodeId parameter.
	node, err := nodeDbExecutor.GetNode(nodeId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getNodeAddress(node)
	urls := util.MakeRequestUrl(address, url.Management(), url.Apps())

	// Request list of applications that is deployed to node.
	codes, respStr := httpExecutor.SendHttpRequest("GET", urls, nil)

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// GetApp gets the application's information of the node specified by nodeId parameter.
// If response code represents success, returns information of application.
// Otherwise, an appropriate error will be returned.
func (Executor) GetApp(nodeId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node including app specified by appId parameter.
	node, err := nodeDbExecutor.GetNodeByAppID(nodeId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getNodeAddress(node)
	urls := util.MakeRequestUrl(address, url.Management(), url.Apps(), "/", appId)

	// Request get target application's information
	codes, respStr := httpExecutor.SendHttpRequest("GET", urls, nil)

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// UpdateApp request to update an application specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UpdateAppInfo(nodeId string, appId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node including app specified by appId parameter.
	node, err := nodeDbExecutor.GetNodeByAppID(nodeId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getNodeAddress(node)
	urls := util.MakeRequestUrl(address, url.Management(), url.Apps(), "/", appId)

	// Request update target application's information.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls, nil, []byte(body))

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// DeleteApp request to delete an application specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) DeleteApp(nodeId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node including app specified by appId parameter.
	node, err := nodeDbExecutor.GetNodeByAppID(nodeId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getNodeAddress(node)
	urls := util.MakeRequestUrl(address, url.Management(), url.Apps(), "/", appId)

	// Request delete target application
	codes, respStr := httpExecutor.SendHttpRequest("DELETE", urls, nil)

	// Convert the received response from string to map.
	result := codes[0]
	if !util.IsSuccessCode(result) {
		respMap, err := convertRespToMap(respStr)
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
		return result, respMap, err
	}

	// if response code represents success, delete the appId from nodeDbExecutor.
	err = nodeDbExecutor.DeleteAppFromNode(nodeId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	err = appDbExecutor.DeleteApp(appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	notiExecutor.UpdateSubscriber()
	
	return result, nil, err
}

// UpdateAppInfo request to update all of images which is included an application
// specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UpdateApp(nodeId string, appId string, query map[string]interface{}) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node including app specified by appId parameter.
	node, err := nodeDbExecutor.GetNodeByAppID(nodeId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getNodeAddress(node)
	urls := util.MakeRequestUrl(address, url.Management(), url.Apps(), "/", appId, url.Update())

	// Request checking and updating all of images which is included target.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls, query)

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// StartApp request to start an application specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) StartApp(nodeId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node including app specified by appId parameter.
	node, err := nodeDbExecutor.GetNodeByAppID(nodeId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getNodeAddress(node)
	urls := util.MakeRequestUrl(address, url.Management(), url.Apps(), "/", appId, url.Start())

	// Request start target application.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls, nil)

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// StopApp request to stop an application specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) StopApp(nodeId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node including app specified by appId parameter.
	node, err := nodeDbExecutor.GetNodeByAppID(nodeId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getNodeAddress(node)
	urls := util.MakeRequestUrl(address, url.Management(), url.Apps(), "/", appId, url.Stop())

	// Request stop target application.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls, nil)

	// Convert the received response from string to map.
	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// getNodeAddress returns an address as an array.
func getNodeAddress(node map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 1)
	result[0] = map[string]interface{}{
		"ip":     node["ip"],
		"config": node["config"],
	}
	return result
}

// convertRespToMap converts a response in the form of JSON data into a map.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func convertRespToMap(respStr []string) (map[string]interface{}, error) {
	resp, err := util.ConvertJsonToMap(respStr[0])
	if err != nil {
		logger.Logging(logger.ERROR, "Failed to convert response from string to map")
		return nil, errors.InternalServerError{"Json Converting Failed"}
	}
	return resp, err
}

func generateRandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = LETTERBYTES[rand.Intn(len(LETTERBYTES))]
	}
	return string(b)
}
