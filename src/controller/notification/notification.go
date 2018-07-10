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

package notification

import (
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	"commons/util"
	nodeSearch "controller/search/node"
	"crypto/sha1"
	appEventDB "db/mongo/event/app"
	nodeEventDB "db/mongo/event/node"
	subsDB "db/mongo/event/subscriber"
	nodeDB "db/mongo/node"
	"encoding/hex"
	"encoding/json"
	"messenger"
	"strings"
)

// Command is an interface of notification operations.
type Command interface {
	Register(body string, query map[string][]string) (int, map[string]interface{}, error)
	UnRegister(eventId string) (int, error)
	UpdateSubscriber()
	NotificationHandler(eventType string, body string) (int, error)
}

const (
	ID                = "id"
	GROUP_ID          = "groupid"
	NODE_ID           = "nodeid"
	APP_ID            = "appid"
	IMAGE_NAME        = "imagename"
	APP               = "app"
	NODE              = "node"
	NODES             = "nodes"
	SUBS              = "subscriber"
	EVENT             = "event"
	EVENT_ID          = "eventid"
	RESPONSES         = "response"
	DEFAULT_NODE_PORT = "48098"
	RESPONSE_CODE     = "code"
	ERROR_MESSAGE     = "message"
	TYPE              = "type"
	STATUS            = "status"
)

// Executor implements the Command interface.
type Executor struct{}

var subsDbExecutor subsDB.Command
var appEventDbExecutor appEventDB.Command
var nodeEventDbExecutor nodeEventDB.Command
var nodeSearchExecutor nodeSearch.Command
var httpExecutor messenger.Command
var nodeDbExecutor nodeDB.Command

func init() {
	subsDbExecutor = subsDB.Executor{}
	appEventDbExecutor = appEventDB.Executor{}
	nodeEventDbExecutor = nodeEventDB.Executor{}
	nodeSearchExecutor = nodeSearch.Executor{}
	httpExecutor = messenger.NewExecutor()
	nodeDbExecutor = nodeDB.Executor{}
}

func (Executor) Register(body string,
	query map[string][]string) (int, map[string]interface{}, error) {

	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	bodyMap, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	// Check whether 'URL' is included.
	url, exists := bodyMap["url"].(string)
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"url field is required"}
	}

	// Check whether 'event' is included.
	event, exists := bodyMap["event"]
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"event field is required"}
	}

	switch parseEventType(event.(map[string]interface{})) {
	default:
		return results.ERROR, nil, err
	case APP:
		result, resp, err := registerAppEvent(url, event.(map[string]interface{}), query)
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
		return result, resp, err
	case NODE:
		result, resp, err := registerNodeEvent(url, event.(map[string]interface{}), query)
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
		return result, resp, err
	}
}

func (Executor) UpdateSubscriber() {
	subscribers, err := subsDbExecutor.GetSubscribers()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return
	}

	for _, subscriber := range subscribers {
		status := make([]interface{}, len(subscriber["status"].([]string)))
		for i, v := range subscriber["status"].([]string) {
			status[i] = v
		}
		subscriber["status"] = status

		switch subscriber["type"] {
		case NODE:
			registerNodeEvent(subscriber["url"].(string), subscriber, subscriber["query"].(map[string][]string))
		case APP:
			registerAppEvent(subscriber["url"].(string), subscriber, subscriber["query"].(map[string][]string))
		}
	}
}

func (Executor) UnRegister(eventId string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	subs, err := subsDbExecutor.GetSubscriber(eventId)
	if err != nil {
		return results.ERROR, err
	}

	switch subs[TYPE] {
	case APP:
		for _, appEventId := range subs[EVENT_ID].([]string) {
			err = appEventDbExecutor.UnRegisterEvent(appEventId, subs[ID].(string))
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, err
			}
			appEvent, err := appEventDbExecutor.GetEvent(appEventId)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, err
			}

			if len(appEvent[SUBS].([]string)) == 0 {
				requestUnRegisterAppEvent(appEvent[NODES].([]string), appEventId)
				err = appEventDbExecutor.DeleteEvent(appEventId)
				if err != nil {
					logger.Logging(logger.ERROR, err.Error())
					return results.ERROR, err
				}
			}
		}
	case NODE:
		for _, nodeEventId := range subs[EVENT_ID].([]string) {
			err = nodeEventDbExecutor.UnRegisterEvent(nodeEventId, subs[ID].(string))
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, err
			}
			nodeEvent, err := nodeEventDbExecutor.GetEvent(nodeEventId)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, err
			}

			if len(nodeEvent[SUBS].([]string)) == 0 {
				err = nodeEventDbExecutor.DeleteEvent(nodeEventId)
				if err != nil {
					logger.Logging(logger.ERROR, err.Error())
					return results.ERROR, err
				}
			}
		}
	}
	err = subsDbExecutor.DeleteSubscriber(eventId)
	if err != nil {
		return results.ERROR, err
	}

	return results.OK, nil
}

func (Executor) NotificationHandler(eventType string, body string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	bodyMap, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}
	// Check whether 'EventId' is included.
	eventIds, exists := bodyMap[EVENT_ID]
	if !exists {
		logger.Logging(logger.ERROR, "eventid field is required")
		return results.ERROR, nil
	}

	// Check whether 'Event' is included.
	event, exists := bodyMap[EVENT]
	if !exists {
		logger.Logging(logger.ERROR, "event field is required")
		return results.ERROR, nil
	}

	switch eventType {
	case APP:
		for _, eventId := range eventIds.([]interface{}) {
			appEvent, err := appEventDbExecutor.GetEvent(eventId.(string))
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, err
			}
			for _, subscriberId := range appEvent[SUBS].([]string) {
				subs, err := subsDbExecutor.GetSubscriber(subscriberId)
				if err != nil {
					logger.Logging(logger.ERROR, err.Error())
					return results.ERROR, err
				}

				for _, status := range subs[STATUS].([]string) {
					if strings.Compare(status, event.(map[string]interface{})[STATUS].(string)) == 0 {
						urls := make([]string, 0)
						urls = append(urls, subs["url"].(string))
						reqBody := make(map[string]interface{})
						reqBody[EVENT] = event
						body, err := convertMapToJson(reqBody)
						if err != nil {
							logger.Logging(logger.ERROR, err.Error())
							return results.ERROR, err
						}
						httpExecutor.SendHttpRequest("POST", urls, nil, []byte(body))
						return results.OK, nil
					}
				}
			}
		}
	case NODE:
		for _, eventId := range eventIds.([]interface{}) {
			nodeEvent, err := nodeEventDbExecutor.GetEvent(eventId.(string))
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, err
			}

			for _, subscriberId := range nodeEvent[SUBS].([]string) {
				subs, err := subsDbExecutor.GetSubscriber(subscriberId)
				if err != nil {
					logger.Logging(logger.ERROR, err.Error())
					return results.ERROR, err
				}

				for _, status := range subs[STATUS].([]string) {
					if strings.Compare(status, event.(map[string]interface{})[STATUS].(string)) == 0 {
						urls := make([]string, 0)
						urls = append(urls, subs["url"].(string))
						reqBody := make(map[string]interface{})
						reqBody[EVENT] = event
						body, err := convertMapToJson(reqBody)
						if err != nil {
							return results.ERROR, err
						}
						httpExecutor.SendHttpRequest("POST", urls, nil, []byte(body))
						return results.OK, nil
					}
				}
			}
		}
	}
	return results.ERROR, nil
}

func registerAppEvent(url string, event map[string]interface{},
	query map[string][]string) (int, map[string]interface{}, error) {

	eventId := make([]string, 0)
	eventId = append(eventId, generateEventId(query))

	nodes, err := getTargetNodes(query)
	if err != nil {
		switch err.(type) {
		default:
			return results.ERROR, nil, err
		case errors.NotFound:
			break
		}
	}

	nodeIds := make([]string, 0)
	for _, node := range nodes[NODES].([]map[string]interface{}) {
		nodeIds = append(nodeIds, node["id"].(string))
	}

	addedNodes := make([]map[string]interface{}, 0)
	removedNodes := make([]map[string]interface{}, 0)

	eventInfo, err := appEventDbExecutor.GetEvent(eventId[0])
	if err != nil {
		switch err.(type) {
		default:
			return results.ERROR, nil, err
		case errors.NotFound:
			addedNodes = nodes[NODES].([]map[string]interface{})
		}
	} else {
		for _, node := range nodes[NODES].([]map[string]interface{}) {
			if !util.IsContainedStringInList(eventInfo["nodes"].([]string), node["id"].(string)) {
				addedNodes = append(addedNodes, node)
			}
		}

		for _, id := range eventInfo["nodes"].([]string) {
			if !util.IsContainedStringInList(nodeIds, id) {
				node, err := nodeDbExecutor.GetNode(id)
				if err != nil {
					logger.Logging(logger.ERROR, err.Error())
				} else {
					removedNodes = append(removedNodes, node)
				}
			}
		}
	}

	// Request unregister event target application.
	removedNodeAddress := getNodesAddress(removedNodes)
	urls := util.MakeRequestUrl(removedNodeAddress, URL.Notification(), URL.Apps(), URL.Watch())
	requestUnRegisterAppEvent(urls, eventId[0])

	// Request register event target application.
	addedNodeAddress := getNodesAddress(addedNodes)
	urls = util.MakeRequestUrl(addedNodeAddress, URL.Notification(), URL.Apps(), URL.Watch())
	codes, respStr := requestRegisterAppEvent(urls, query, eventId[0])

	// Convert the received response from string to map.
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		requestUnRegisterAppEvent(urls, eventId[0])
		return results.ERROR, nil, err
	}

	eventStatus := parseEventStatus(event)
	subsId := generateSubsId(eventId[0], url, eventStatus)
	resp := make(map[string]interface{})

	result := decideResultCode(codes)
	if result == results.ERROR {
		return results.ERROR, nil, err
	} else if result == results.MULTI_STATUS {
		// Make separate responses to represent partial failure case.
		resp[RESPONSES] = makeSeparateResponses(nodes[NODES].([]map[string]interface{}), codes, respMap)
	}

	err = subsDbExecutor.AddSubscriber(subsId, APP, url, eventStatus, eventId, query)
	if err != nil {
		return results.ERROR, nil, err
	}

	for i, node := range addedNodes {
		if !util.IsSuccessCode(codes[i]) {
			for _, id := range nodeIds {
				if strings.Compare(id, node["id"].(string)) == 0 {
					nodeIds = append(nodeIds[:i], nodeIds[i+1:]...)
				}
			}
		}
	}

	err = appEventDbExecutor.AddEvent(eventId[0], subsId, nodeIds)
	if err != nil {
		return results.ERROR, nil, err
	}

	resp[ID] = subsId

	return result, resp, err
}

func registerNodeEvent(url string, event map[string]interface{},
	query map[string][]string) (int, map[string]interface{}, error) {

	nodes, err := getTargetNodes(query)
	if err != nil {
		return results.ERROR, nil, err
	}

	nodeIds := make([]string, 0)
	for _, node := range nodes[NODES].([]map[string]interface{}) {
		nodeIds = append(nodeIds, node["id"].(string))
	}

	eventStatus := parseEventStatus(event)
	subsId := generateSubsId(generateEventId(query), url, eventStatus)
	err = subsDbExecutor.AddSubscriber(subsId, NODE, url, eventStatus, nodeIds, query)
	if err != nil {
		return results.ERROR, nil, err
	}

	for _, nodeId := range nodeIds {
		err = nodeEventDbExecutor.AddEvent(nodeId, subsId)
		if err != nil {
			return results.ERROR, nil, err
		}
	}
	resp := make(map[string]interface{})
	resp[ID] = subsId

	return results.OK, resp, err
}

func requestRegisterAppEvent(urls []string, query map[string][]string, eventId string) ([]int, []string) {
	if len(urls) == 0 {
		return nil, nil
	}
	reqBody := makeRequestBody(query, eventId)
	body, err := convertMapToJson(reqBody)
	if err != nil {
		return nil, nil
	}
	// Request register event target application.
	return httpExecutor.SendHttpRequest("POST", urls, nil, []byte(body))
}

func requestUnRegisterAppEvent(urls []string, eventId string) {
	if len(urls) == 0 {
		return
	}
	reqBody := makeRequestBody(nil, eventId)
	body, _ := convertMapToJson(reqBody)

	// Request unregister event target nodes.
	httpExecutor.SendHttpRequest("DELETE", urls, nil, []byte(body))
}

func getTargetNodes(query map[string][]string) (map[string]interface{}, error) {
	_, nodes, err := nodeSearchExecutor.SearchNodes(query)
	if err != nil {
		return nil, err
	}
	if nodes == nil {
		logger.Logging(logger.DEBUG, "No matched nodes with query")
		return nil, errors.NotFound{"No matched nodes with query"}
	}

	return nodes, err
}

func getSucceedNodesId(nodes []map[string]interface{}, codes []int) []string {
	nodeId := make([]string, 0)
	for i, node := range nodes {
		if util.IsSuccessCode(codes[i]) {
			nodeId = append(nodeId, node[ID].(string))
		}
	}
	return nodeId
}

func getNodesAddress(nodes []map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(nodes))
	for i, node := range nodes {
		result[i] = map[string]interface{}{
			"ip":     node["ip"],
			"config": node["config"],
		}
	}
	return result
}

func parseEventStatus(event map[string]interface{}) []string {
	statusList := make([]string, 0)
	for _, status := range event["status"].([]interface{}) {
		statusList = append(statusList, status.(string))
	}
	return statusList
}

func parseEventType(event map[string]interface{}) string {
	return event["type"].(string)
}

func generateSubsId(eventId string, url string, eventStatus []string) string {
	var source string
	source += eventId
	source += url
	for _, status := range eventStatus {
		source += status
	}
	return makeHash(source)
}

func generateEventId(query map[string][]string) string {
	var source string
	if _, ok := query[GROUP_ID]; ok {
		source += query[GROUP_ID][0]
	} else {
		source += "ALL"
	}
	if _, ok := query[NODE_ID]; ok {
		source += query[NODE_ID][0]
	} else {
		source += "ALL"
	}
	if _, ok := query[IMAGE_NAME]; ok {
		source += query[IMAGE_NAME][0]
	} else {
		source += "ALL"
	}
	if _, ok := query[APP_ID]; ok {
		source += query[APP_ID][0]
	} else {
		source += "ALL"
	}
	return makeHash(source)
}

func makeRequestBody(query map[string][]string, eventId string) map[string]interface{} {
	reqBody := make(map[string]interface{})

	if query != nil {
		if _, ok := query[IMAGE_NAME]; ok {
			reqBody[IMAGE_NAME] = query[IMAGE_NAME][0]
		}
		if _, ok := query[APP_ID]; ok {
			reqBody[APP_ID] = query[APP_ID][0]
		}
	}
	reqBody[EVENT_ID] = eventId

	return reqBody
}

// Making hash code by hash value.
// return hash code as slice of byte
func makeHash(source string) string {
	h := sha1.New()
	h.Write([]byte(source))
	return hex.EncodeToString(h.Sum(nil))
}

// convertMapToJson converts map data into a JSON.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func convertMapToJson(reqBody map[string]interface{}) (string, error) {
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return string(""), errors.InvalidJSON{"Marshalling Failed"}
	}
	return string(jsonBody), err
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

// decideResultCode returns a result of group operations.
// OK: Returned when all members of the group send a success response.
// MULTI_STATUS: Partial success for multiple requests. Some requests succeeded
//               but at least one failed.
// ERROR: Returned when all members of the gorup send an error response.
func decideResultCode(codes []int) int {
	successCounts := 0
	for _, code := range codes {
		if util.IsSuccessCode(code) {
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

// makeSeparateResponses used to make a separate response.
func makeSeparateResponses(nodes []map[string]interface{}, codes []int,
	respMap []map[string]interface{}) []map[string]interface{} {

	respValue := make([]map[string]interface{}, len(nodes))

	for i, node := range nodes {
		respValue[i] = make(map[string]interface{})
		respValue[i][ID] = node[ID].(string)
		respValue[i][RESPONSE_CODE] = codes[i]

		if !util.IsSuccessCode(codes[i]) {
			respValue[i][ERROR_MESSAGE] = respMap[i][ERROR_MESSAGE]
		}
	}

	return respValue
}
