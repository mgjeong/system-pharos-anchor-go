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

// Package commons/util defines utility functions used by Pharos Node.
package util

import (
	"bytes"
	"commons/errors"
	"commons/logger"
	"commons/url"
	"encoding/json"
	"net"
	"strings"
)

const (
	DEFAULT_NODE_PORT                      = "48098"
	UNSECURED_NODE_PORT_WITH_REVERSE_PROXY = "80"
	SECURED_NODE_PORT_WITH_REVERSE_PROXY   = "443"
	NODE_URL_PREFIX                        = "/pharos-node"
)

// convertJsonToMap converts JSON data into a map.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func ConvertJsonToMap(jsonStr string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, errors.InvalidJSON{"Unmarshalling Failed"}
	}
	return result, err
}

// ConvertMapToJson converts Map data to json data.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func ConvertMapToJson(data map[string]interface{}) (string, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	result, err := json.Marshal(data)
	if err != nil {
		return "", errors.Unknown{"json marshalling failed"}
	}
	return string(result), nil
}

// isSuccessCode returns true in case of success and false otherwise.
func IsSuccessCode(code int) bool {
	if code >= 200 && code <= 299 {
		return true
	}
	return false
}

func IsContainedStringInList(list []string, str string) bool {
	for _, value := range list {
		if strings.Compare(value, str) == 0 {
			return true
		}
	}
	return false
}

// MakeRequestUrl make a list of urls that can be used to send a http request.
func MakeRequestUrl(address []map[string]interface{}, api_parts ...string) (urls []string) {
	var httpTag string = "http://"
	var full_url bytes.Buffer

	for i := range address {
		full_url.Reset()

		ipTest := net.ParseIP(address[i]["ip"].(string))
		if ipTest == nil {
			logger.Logging(logger.ERROR, address[i]["ip"].(string), "Validation check on anchor's ip  failed")
			urls = append(urls, "")
			continue
		}

		config := address[i]["config"].(map[string]interface{})
		properties := config["properties"].([]interface{})

		reverseproxy := make(map[string]interface{}, 0)
		for i := range properties {
			property := make(map[string]interface{}, 0)

			marshaled, _ := json.Marshal(properties[i])
			err := json.Unmarshal(marshaled, &property)
			if err != nil {
				logger.Logging(logger.ERROR, "Unmarshal configuration property failed")
				urls = append(urls, "")
				continue
			}

			if _, exists := property["reverseproxy"]; exists {
				reverseproxy = property["reverseproxy"].(map[string]interface{})
			}
		}

		v, exists := reverseproxy["enabled"]
		if !exists {
			logger.Logging(logger.ERROR, "No node's reverse proxy environment")
			urls = append(urls, "")
			continue
		}

		switch v.(type) {
		case bool:
		default:
			logger.Logging(logger.ERROR, "Invalid type for reverseproxy value")
			urls = append(urls, "")
			continue
		}

		if reverseproxy["enabled"].(bool) == true {
			full_url.WriteString(httpTag + address[i]["ip"].(string) + ":" + UNSECURED_NODE_PORT_WITH_REVERSE_PROXY + url.PharosNode() + url.Base())
		} else {
			full_url.WriteString(httpTag + address[i]["ip"].(string) + ":" + DEFAULT_NODE_PORT + url.Base())
		}

		for _, api_part := range api_parts {
			full_url.WriteString(api_part)
		}
		urls = append(urls, full_url.String())
	}
	return urls
}
