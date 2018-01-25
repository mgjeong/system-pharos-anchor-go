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
package errors

import (
	"strings"
	"testing"
)

func TestTError(t *testing.T) {
	msg := "Test"

	type commonsError interface {
		Error() string
	}

	type testObj struct {
		testName   string
		testPrefix string
		testError  commonsError
	}

	testList := []testObj{
		{testName: "Unknown", testPrefix: "unknown error",
			testError: &Unknown{msg}},
		{testName: "NotFoundURL", testPrefix: "unsupported url",
			testError: &NotFoundURL{msg}},
		{testName: "InvalidMethod", testPrefix: "invalid method",
			testError: &InvalidMethod{msg}},
		{testName: "InvalidParam", testPrefix: "invalid parameter",
			testError: &InvalidParam{msg}},
		{testName: "InvalidJSON", testPrefix: "invalid json format",
			testError: &InvalidJSON{msg}},
		{testName: "InvalidYaml", testPrefix: "invalid yaml file",
			testError: &InvalidYaml{msg}},
		{testName: "InvalidObjectId", testPrefix: "invalid objectId",
			testError: &InvalidObjectId{msg}},
		{testName: "NotFound", testPrefix: "not found target",
			testError: &NotFound{msg}},
		{testName: "DBConnectionError", testPrefix: "db connection failed",
			testError: &DBConnectionError{msg}},
		{testName: "DBOperationError", testPrefix: "db operation failed",
			testError: &DBOperationError{msg}},
		{testName: "IOError", testPrefix: "io error",
			testError: &IOError{msg}},
		{testName: "InternalServerError", testPrefix: "internal server error",
			testError: &InternalServerError{msg}},
	}

	testFunc := func(err commonsError, prefix string) {
		ret := err.Error()
		if !strings.HasPrefix(ret, prefix) {
			t.Error()
		} else if !strings.HasSuffix(ret, msg) {
			t.Error()
		}
	}

	for _, test := range testList {
		t.Run(test.testName, func(t *testing.T) {
			testFunc(test.testError, test.testPrefix)
		})
	}
}
