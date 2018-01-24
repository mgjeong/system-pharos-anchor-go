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
	"bytes"
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	"db/mongo/registry"
	"encoding/json"
	"messenger"
)

type Command interface {
	// AddDockerRegistry add docker registry to database.
	AddDockerRegistry(body string) (int, map[string]interface{}, error)

	DeleteDockerRegistry(registryId string) (int, error)

	GetDockerRegistries() (int, map[string]interface{}, error)

	GetDockerRegistry(registryId string) (int, map[string]interface{}, error)

	GetDockerImages(registryId string) (int, map[string]interface{}, error)

	DockerRegistryEventHandler(body string) (int, error)
}

const (
	ID           = "id"
	IP           = "ip"
	TARGETINFO   = "target"
	REQUESTINFO  = "request"
	REGISTRIES   = "registries"
	REGISTRY     = "registry"
	IMAGES       = "images"
	REPOSITORIES = "repositories"
	HOST         = "host"
	REPOSITORY   = "repository"
	TAG          = "tag"
	SIZE         = "size"
	ACTION       = "action"
	EVENTS       = "events"
	TIMESTAMP    = "timestamp"
	PUSH         = "push"
	DELETE       = "delete"
)

type Executor struct{}

var dbExecutor registry.Command
var httpExecutor messenger.Command

func init() {
	dbExecutor = registry.Executor{}
	httpExecutor = messenger.NewExecutor()
}

func (Executor) AddDockerRegistry(body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	reqBody, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Check the URL is valiadated or not with catalog API.
	urls := makeRequestUrl(reqBody[IP].(string), URL.Catalog())

	codes, respStr := httpExecutor.SendHttpRequest("GET", urls, nil)
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	result := codes[0]
	if isSuccessCode(result) {
		// Add new registry to database with given url.
		registry, err := dbExecutor.AddDockerRegistry(reqBody[IP].(string))
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}

		imagesList := respMap[REPOSITORIES]

		// if registry has repository, add the list of repository.
		if len(imagesList.([]interface{})) > 0 {
			imagesInfo := make([]map[string]interface{}, len(imagesList.([]interface{})))

			for i, imageInfo := range imagesList.([]interface{}) {
				imagesInfo[i][REPOSITORY] = imageInfo.(string)
				imagesInfo[i][TAG] = ""
				imagesInfo[i][SIZE] = ""
				imagesInfo[i][ACTION] = ""
				imagesInfo[i][TIMESTAMP] = ""
			}

			err = dbExecutor.AddDockerImages(registry[ID].(string), imagesInfo)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, nil, err
			}
		}

		res := make(map[string]interface{})
		res[ID] = registry[ID]

		return results.OK, nil, err
	}

	return result, nil, err
}

func (Executor) DeleteDockerRegistry(registryId string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Delete registry specified by registryId parameter.
	err := dbExecutor.DeleteDockerRegistry(registryId)
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
	registries, err := dbExecutor.GetDockerRegistries()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[REGISTRIES] = registries

	return results.OK, res, err
}

func (Executor) GetDockerRegistry(registryId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get registry specified by registryId parameter.
	registry, err := dbExecutor.GetDockerRegistry(registryId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[REGISTRY] = registry

	return results.OK, res, err
}

func (Executor) GetDockerImages(registryId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Get all of images list on registry specified by registryId parameter.
	images, err := dbExecutor.GetDockerImages(registryId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[IMAGES] = images

	return results.OK, res, err
}

func (Executor) DockerRegistryEventHandler(body string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	convertedBody, err := convertJsonToMap(body)
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

		switch parsedEvent[ACTION] {
		case PUSH:
			err := addDockerImage(parsedEvent)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, err
			}
		case DELETE:
			err := deleteDockerImage(parsedEvent)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, err
			}
		}
	}

	return results.OK, nil
}

func addDockerImage(imageInfo map[string]interface{}) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	host, err := dbExecutor.GetDockerRegistry(imageInfo[HOST].(string))
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return err
	}

	delete(imageInfo, HOST)
	images := make([]map[string]interface{}, 0)
	images = append(images, imageInfo)

	err = dbExecutor.AddDockerImages(host[ID].(string), images)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return err
	}

	return nil
}

func deleteDockerImage(imageInfo map[string]interface{}) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	host, err := dbExecutor.GetDockerRegistry(imageInfo[HOST].(string))
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return err
	}

	delete(imageInfo, HOST)

	err = dbExecutor.DeleteDockerImage(host[ID].(string), imageInfo)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return err
	}

	return nil
}

// convertRespToMap converts a response in the form of JSON data into a map.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func convertRespToMap(respStr []string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	resp, err := convertJsonToMap(respStr[0])
	if err != nil {
		logger.Logging(logger.ERROR, "Failed to convert response from string to map")
		return nil, errors.InternalServerError{"Json Converting Failed"}
	}
	return resp, err
}

// convertJsonToMap converts JSON data into a map.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func convertJsonToMap(jsonStr string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, errors.InvalidJSON{"Unmarshalling Failed"}
	}
	return result, err
}

// makeRequestUrl make a list of urls that can be used to send a http request.
func makeRequestUrl(ip string, api_parts ...string) (urls []string) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	var httpTag string = "http://"
	var full_url bytes.Buffer

	full_url.Reset()
	full_url.WriteString(httpTag + ip)

	for _, api_part := range api_parts {
		full_url.WriteString(api_part)
	}
	urls = append(urls, full_url.String())

	return urls
}

// isSuccessCode returns true in case of success and false otherwise.
func isSuccessCode(code int) bool {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	if code >= 200 && code <= 299 {
		return true
	}
	return false
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

	parsedEvent[ACTION] = eventInfo[ACTION]
	parsedEvent[TIMESTAMP] = eventInfo[TIMESTAMP]
	parsedEvent[REPOSITORY] = targetInfoEvent[REPOSITORY]
	parsedEvent[TAG] = targetInfoEvent[TAG]
	parsedEvent[SIZE] = targetInfoEvent[SIZE]
	parsedEvent[HOST] = requestInfoEvent[HOST]

	return parsedEvent, nil
}
