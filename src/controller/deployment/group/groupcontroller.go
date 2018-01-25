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

// Package group provides an interfaces to deploy, update, start, stop, delete
// an application to target edge group.
// and also provides operations to manage edge device group (e.g., create, join, leave, delete...).
package group

import (
	"bytes"
	"commons/errors"
	"commons/logger"
	"commons/results"
	"commons/url"
	appDB "db/mongo/app"
	groupDB "db/mongo/group"
	nodeDB "db/mongo/node"
	"encoding/json"
	"messenger"
)

const (
	NODES              = "nodes"       // used to indicate a list of nodes.
	GROUPS             = "groups"      // used to indicate a list of groups.
	MEMBERS            = "members"     // used to indicate a list of members.
	APPS               = "apps"        // used to indicate a list of apps.
	ID                 = "id"          // used to indicate an id.
	RESPONSE_CODE      = "code"        // used to indicate a code.
	ERROR_MESSAGE      = "message"     // used to indicate a message.
	RESPONSES          = "responses"   // used to indicate a list of responses.
	DESCRIPTION        = "description" // used to indicate a description.
	DEFAULT_AGENT_PORT = "48098"       // used to indicate a default system-management-node port.
)

type Executor struct{}

var appDbExecutor appDB.Command
var nodeDbExecutor nodeDB.Command
var groupDbExecutor groupDB.Command
var httpExecutor messenger.Command

func init() {
	appDbExecutor = appDB.Executor{}
	nodeDbExecutor = nodeDB.Executor{}
	groupDbExecutor = groupDB.Executor{}
	httpExecutor = messenger.NewExecutor()
}

// Command is an interface of group deployment operations.
type Command interface {
	// DeployApp request an deployment of edge services to a group specified by groupId parameter.
	DeployApp(groupId string, body string) (int, map[string]interface{}, error)

	// GetApps request a list of applications that is deployed to a group specified by groupId parameter.
	GetApps(groupId string) (int, map[string]interface{}, error)

	// GetApp gets the application's information of the group specified by groupId parameter.
	GetApp(groupId string, appId string) (int, map[string]interface{}, error)

	// UpdateApp request to update an application specified by appId parameter to all members of the group.
	UpdateAppInfo(groupId string, appId string, body string) (int, map[string]interface{}, error)

	// DeleteApp request to delete an application specified by appId parameter to all members of the group.
	DeleteApp(groupId string, appId string) (int, map[string]interface{}, error)

	// UpdateAppInfo request to update all of images which is included an application specified by
	// appId parameter to all members of the group.
	UpdateApp(groupId string, appId string) (int, map[string]interface{}, error)

	// StartApp request to start an application specified by appId parameter to all members of the group.
	StartApp(groupId string, appId string) (int, map[string]interface{}, error)

	// StopApp request to stop an application specified by appId parameter to all members of the group.
	StopApp(groupId string, appId string) (int, map[string]interface{}, error)
}

// DeployApp request an deployment of edge services to a group specified by groupId parameter.
// If response code represents success, add an app id to a list of installed app and returns it.
// Otherwise, an appropriate error will be returned.
func (Executor) DeployApp(groupId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get group members from the database.
	members, err := groupDbExecutor.GetGroupMembers(groupId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getMemberAddress(members)
	urls := makeRequestUrl(address, url.Deploy())

	// Request an deployment of edge services to a specific group.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls, nil, []byte(body))

	// Convert the received response from string to map.
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// if response code represents success, insert the installed appId into groupDbExecutor.
	installedAppId := ""
	for i, node := range members {
		if isSuccessCode(codes[i]) {
			err = appDbExecutor.AddApp(respMap[i]["id"].(string), []byte(respMap[i]["description"].(string)))
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, nil, err
			}

			err = nodeDbExecutor.AddAppToNode(node[ID].(string), respMap[i][ID].(string))
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, nil, err
			}
			installedAppId = respMap[i][ID].(string)
		}
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)
		if installedAppId != "" {
			resp[ID] = installedAppId
		}
		return result, resp, err
	}

	resp := make(map[string]interface{})
	resp[ID] = installedAppId

	return result, resp, err
}

// GetApps request a list of applications that is deployed to a group
// specified by groupId parameter.
// If response code represents success, returns a list of applications.
// Otherwise, an appropriate error will be returned.
func (Executor) GetApps(groupId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get group members from the database.
	members, err := groupDbExecutor.GetGroupMembers(groupId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	contains := func(list []map[string]interface{}, appId string) map[string]interface{} {
		for _, item := range list {
			if item[ID] == appId {
				return item
			}
		}
		return nil
	}

	respValue := make([]map[string]interface{}, 0)
	for _, node := range members {
		for _, appId := range node[APPS].([]string) {
			item := contains(respValue, appId)
			if item != nil {
				item[MEMBERS] = append(item[MEMBERS].([]string), node[ID].(string))
			} else {
				item = map[string]interface{}{
					ID:      appId,
					MEMBERS: []string{node[ID].(string)},
				}
				respValue = append(respValue, item)
			}
		}
	}

	res := make(map[string]interface{})
	res[APPS] = respValue

	return results.OK, res, err
}

// GetApp gets the application's information of the group specified by groupId parameter.
// If response code represents success, returns information of application.
// Otherwise, an appropriate error will be returned.
func (Executor) GetApp(groupId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get group members including app specified by appId parameter.
	members, err := groupDbExecutor.GetGroupMembersByAppID(groupId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getMemberAddress(members)
	urls := makeRequestUrl(address, url.Apps(), "/", appId)

	// Request get target application's information.
	codes, respStr := httpExecutor.SendHttpRequest("GET", urls, nil)

	// Convert the received response from string to map.
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)

		for i, _ := range members {
			if isSuccessCode(codes[i]) {
				respValue := resp[RESPONSES].([]map[string]interface{})
				for key, value := range respMap[i] {
					respValue[i][key] = value
				}
			}
		}
		return result, resp, err
	}

	resp := make(map[string]interface{})
	respValue := make([]map[string]interface{}, len(members))
	resp[RESPONSES] = respValue

	for i, node := range members {
		respValue[i] = make(map[string]interface{})
		respValue[i][ID] = node[ID].(string)
		for key, value := range respMap[i] {
			respValue[i][key] = value
		}
	}

	return result, resp, err
}

// UpdateApp request to update an application specified by appId parameter
// to all members of the group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UpdateAppInfo(groupId string, appId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get group members including app specified by appId parameter.
	members, err := groupDbExecutor.GetGroupMembersByAppID(groupId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getMemberAddress(members)
	urls := makeRequestUrl(address, url.Apps(), "/", appId)

	// Request update target application's information.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls, nil, []byte(body))

	// Convert the received response from string to map.
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)
		return result, resp, err
	}

	return result, nil, err
}

// DeleteApp request to delete an application specified by appId parameter
// to all members of the group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) DeleteApp(groupId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get group members including app specified by appId parameter.
	members, err := groupDbExecutor.GetGroupMembersByAppID(groupId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getMemberAddress(members)
	urls := makeRequestUrl(address, url.Apps(), "/", appId)

	// Request delete target application.
	codes, respStr := httpExecutor.SendHttpRequest("DELETE", urls, nil)

	// Convert the received response from string to map.
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// if response code represents success, delete the appId from groupDbExecutor.
	for i, node := range members {
		if isSuccessCode(codes[i]) {
			err = nodeDbExecutor.DeleteAppFromNode(node[ID].(string), appId)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, nil, err
			}

			err = appDbExecutor.DeleteApp(appId)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, nil, err
			}
		}
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)
		return result, resp, err
	}

	return result, nil, err
}

// UpdateAppInfo request to update all of images which is included an application
// specified by appId parameter to all members of the group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UpdateApp(groupId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get group members including app specified by appId parameter.
	members, err := groupDbExecutor.GetGroupMembersByAppID(groupId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getMemberAddress(members)
	urls := makeRequestUrl(address, url.Apps(), "/", appId, url.Update())

	// Request checking and updating all of images which is included target.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls, nil)

	// Convert the received response from string to map.
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)
		return result, resp, err
	}

	return result, nil, err
}

// StartApp request to start an application specified by appId parameter
// to all members of the group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) StartApp(groupId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get group members including app specified by appId parameter.
	members, err := groupDbExecutor.GetGroupMembersByAppID(groupId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getMemberAddress(members)
	urls := makeRequestUrl(address, url.Apps(), "/", appId, url.Start())

	// Request start target application.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls, nil)

	// Convert the received response from string to map.
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)
		return result, resp, err
	}

	return result, nil, err
}

// StopApp request to stop an application specified by appId parameter
// to all members of the group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) StopApp(groupId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get group members including app specified by appId parameter.
	members, err := groupDbExecutor.GetGroupMembersByAppID(groupId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	address := getMemberAddress(members)
	urls := makeRequestUrl(address, url.Apps(), "/", appId, url.Stop())

	// Request stop target application.
	codes, respStr := httpExecutor.SendHttpRequest("POST", urls, nil)

	// Convert the received response from string to map.
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)
		return result, resp, err
	}

	return result, nil, err
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

// getNodeAddress returns an member's address as an array.
func getMemberAddress(members []map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(members))
	for i, node := range members {
		result[i] = map[string]interface{}{
			"ip": node["ip"],
		}
	}
	return result
}

// convertRespToMap converts a response in the form of JSON data into a map.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func convertRespToMap(respStr []string) ([]map[string]interface{}, error) {
	respMap := make([]map[string]interface{}, len(respStr))
	for i, v := range respStr {
		resp, err := convertJsonToMap(v)
		if err != nil {
			logger.Logging(logger.ERROR, "Failed to convert response from string to map")
			return nil, errors.InternalServerError{"Json Converting Failed"}
		}
		respMap[i] = resp
	}

	return respMap, nil
}

// isSuccessCode returns true in case of success and false otherwise.
func isSuccessCode(code int) bool {
	if code >= 200 && code <= 299 {
		return true
	}
	return false
}

// decideResultCode returns a result of group operations.
// OK: Returned when all members of the group send a success response.
// MULTI_STATUS: Partial success for multiple requests. Some requests succeeded
//               but at least one failed.
// ERROR: Returned when all members of the gorup send an error response.
func decideResultCode(codes []int) int {
	successCounts := 0
	for _, code := range codes {
		if isSuccessCode(code) {
			successCounts++
		}
	}

	result := results.OK
	switch successCounts {
	case len(codes):
		result = results.OK
	case 0:
		result = results.ERROR
	default:
		result = results.MULTI_STATUS
	}
	return result
}

// makeSeparateResponses used to make a separate response
// when the group operations is a partial success.
func makeSeparateResponses(members []map[string]interface{}, codes []int,
	respMap []map[string]interface{}) []map[string]interface{} {

	respValue := make([]map[string]interface{}, len(members))

	for i, node := range members {
		respValue[i] = make(map[string]interface{})
		respValue[i][ID] = node[ID].(string)
		respValue[i][RESPONSE_CODE] = codes[i]

		if !isSuccessCode(codes[i]) {
			respValue[i][ERROR_MESSAGE] = respMap[i][ERROR_MESSAGE]
		}
	}

	return respValue
}

// makeRequestUrl make a list of urls that can be used to send a http request.
func makeRequestUrl(address []map[string]interface{}, api_parts ...string) (urls []string) {
	var httpTag string = "http://"
	var full_url bytes.Buffer

	for i := range address {
		full_url.Reset()
		full_url.WriteString(httpTag + address[i]["ip"].(string) +
			":" + DEFAULT_AGENT_PORT + url.Base())
		for _, api_part := range api_parts {
			full_url.WriteString(api_part)
		}
		urls = append(urls, full_url.String())
	}
	return urls
}
