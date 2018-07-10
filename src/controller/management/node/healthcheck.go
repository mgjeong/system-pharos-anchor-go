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
	"commons/errors"
	"commons/logger"
	"commons/results"
	"commons/util"
	"encoding/json"
	"strconv"
	"time"
)

// PingNode starts timer with received interval.
// If node does not send next healthcheck message in interval time,
// change the status of device from connected to disconnected.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (executor Executor) PingNode(nodeId string, body string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get node specified by nodeId parameter.
	_, _, err := executor.GetNode(nodeId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	bodyMap, err := util.ConvertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	// Check whether 'interval' is included.
	_, exists := bodyMap[INTERVAL]
	if !exists {
		return results.ERROR, errors.InvalidJSON{"interval field is required"}
	}

	interval, err := strconv.Atoi(bodyMap[INTERVAL].(string))
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, errors.InvalidJSON{"invalid value type(interval must be integer)"}
	}

	common.Lock()
	_, exists = common.timers[nodeId]
	if !exists {
		logger.Logging(logger.DEBUG, "first ping request is received from node")
	} else {
		// Status is updated with 'connected'.
		executor.UpdateNodeStatus(nodeId, STATUS_CONNECTED)

		if common.timers[nodeId] != nil {
			// If ping request is received in interval time, send signal to stop timer.
			common.timers[nodeId] <- true
			logger.Logging(logger.DEBUG, "ping request is received in interval time")
		} else {
			logger.Logging(logger.DEBUG, "ping request is received after interval time-out")
			sendNotification(nodeId, STATUS_CONNECTED)
		}
	}
	common.Unlock()

	// Start timer with received interval time.
	timeDurationMin := time.Duration(interval+MAXIMUM_NETWORK_LATENCY_SEC) * TIME_UNIT
	timer := time.NewTimer(timeDurationMin)
	go func() {
		quit := make(chan bool)
		common.Lock()
		common.timers[nodeId] = quit
		common.Unlock()

		select {
		// Block until timer finishes.
		case <-timer.C:
			logger.Logging(logger.ERROR, "ping request is not received in interval time")

			// Status is updated with 'disconnected'.
			err = executor.UpdateNodeStatus(nodeId, STATUS_DISCONNECTED)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
			}
			sendNotification(nodeId, STATUS_DISCONNECTED)

		case <-quit:
			timer.Stop()
			return
		}

		common.Lock()
		common.timers[nodeId] = nil
		close(quit)
		common.Unlock()
	}()

	return results.OK, err
}

func sendNotification(nodeId string, status string) {
	eventIds := make([]string, 0)
	eventIds = append(eventIds, nodeId)
	event := make(map[string]interface{})
	event[ID] = nodeId
	event[STATUS] = status
	event[TIMESTAMP] = time.Now().String()

	notification := make(map[string]interface{})
	notification[EVENT_ID] = eventIds
	notification[EVENT] = event

	notiStr, err := convertMapToJson(notification)
	if err != nil {
		return
	}
	notiExecutor.NotificationHandler(NODE, notiStr)
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
