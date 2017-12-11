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

// Package common implements some utility functions for http transaction.
package common

import (
	"commons/errors"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// WriteSuccess writes the data to the connection as part of an HTTP reply.
func WriteSuccess(w http.ResponseWriter, code int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

// WriteError writes the data to the connection as part of an HTTP reply.
// The http status code depend on an error type.
// An error message will be included as a body.
func WriteError(w http.ResponseWriter, err error) {
	code := convertToHttpStatusCode(err)
	data := make(map[string]interface{})
	data["message"] = err.Error()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(ChangeToJson(data))
}

// MakeResponse calls WriteSuccess or MakeResponse function to respond to the request.
// If err is nil, WriteSuccess will be called.
// otherwise, WriteError will be called.
func MakeResponse(w http.ResponseWriter, code int, data []byte, err error) {
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteSuccess(w, code, data)
}

// ChangeToJson converts map to []byte.
func ChangeToJson(src map[string]interface{}) []byte {
	dst, err := json.Marshal(src)
	if err != nil {
		return nil
	}
	return dst
}

// GetBodyFromReq reads a body from http request object.
// A successful call returns the body type of string.
// If request does not include body, InvalidParam will be returned.
// In other cases, an appropriate error will be returned.
func GetBodyFromReq(req *http.Request) (string, error) {
	if req.Body == nil {
		return "", errors.InvalidParam{"body is empty"}
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", errors.IOError{err.Error()}
	}
	return string(body), nil
}

// convertToHttpStatusCode converts an error object to http status code.
// The following codes are used.
//
//    400 (Bad Request)
//    404 (Not Found)
// 	  500 (Internal Server Error)
//    503 (Service Unavailable)
func convertToHttpStatusCode(err error) int {
	code := http.StatusInternalServerError
	switch err.(type) {
	case errors.InvalidParam,
		errors.InvalidJSON,
		errors.InvalidMethod,
		errors.InvalidObjectId:
		code = http.StatusBadRequest
	case errors.NotFoundURL,
		errors.NotFound:
		code = http.StatusNotFound
	case errors.DBConnectionError,
		errors.DBOperationError:
		code = http.StatusServiceUnavailable
	case errors.IOError,
		errors.InternalServerError:
		code = http.StatusInternalServerError
	}

	return code
}
