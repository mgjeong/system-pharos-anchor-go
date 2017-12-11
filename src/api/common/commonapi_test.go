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
package common

import (
	"bytes"
	Errors "commons/errors"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	WriteSuccess(w, http.StatusOK, nil)
	if w.Code != http.StatusOK {
		t.Error("WriteSuccess is invalid")
	}
}

func TestWriteError(t *testing.T) {
	dummyError := errors.New("")
	w := httptest.NewRecorder()
	WriteError(w, dummyError)
	if w.Code != http.StatusInternalServerError {
		t.Error("WriteError is invalid")
	}
}

func TestMakeResponse(t *testing.T) {
	w := httptest.NewRecorder()
	MakeResponse(w, http.StatusOK, nil, nil)
	if w.Code != http.StatusOK {
		t.Error("MakeResponse is invalid")
	}
}

func TestChangeToJson(t *testing.T) {
	dummyError := errors.New("")
	w := httptest.NewRecorder()
	WriteError(w, dummyError)
	if w.Code != http.StatusInternalServerError {
		t.Error("WriteError is invalid")
	}
}

func TestGetBodyFromReq(t *testing.T) {
	body := []byte("body")
	req, _ := http.NewRequest("POST", "/api/v1/test/url", bytes.NewReader(body))
	_, err := GetBodyFromReq(req)
	if err != nil {
		t.Error("GetBodyFromReq is invalid")
	}
}

func TestGetBodyFromReqWithEmptyBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v1/test/url", nil)
	_, err := GetBodyFromReq(req)
	if err == nil {
		t.Error("GetBodyFromReq is invalid")
	}
}

func TestConvertToHttpStatusCodeWithInvalidParam(t *testing.T) {
	err := Errors.InvalidParam{}
	code := convertToHttpStatusCode(err)
	if code != http.StatusBadRequest {
		t.Error("convertToHttpStatusCode is invalid")
	}
}

func TestConvertToHttpStatusCodeWithInvalidJSON(t *testing.T) {
	err := Errors.InvalidJSON{}
	code := convertToHttpStatusCode(err)
	if code != http.StatusBadRequest {
		t.Error("convertToHttpStatusCode is invalid")
	}
}

func TestConvertToHttpStatusCodeWithInvalidMethod(t *testing.T) {
	err := Errors.InvalidMethod{}
	code := convertToHttpStatusCode(err)
	if code != http.StatusBadRequest {
		t.Error("convertToHttpStatusCode is invalid")
	}
}

func TestConvertToHttpStatusCodeWithInvalidObjectId(t *testing.T) {
	err := Errors.InvalidObjectId{}
	code := convertToHttpStatusCode(err)
	if code != http.StatusBadRequest {
		t.Error("convertToHttpStatusCode is invalid")
	}
}

func TestConvertToHttpStatusCodeWithNotFoundURL(t *testing.T) {
	err := Errors.NotFoundURL{}
	code := convertToHttpStatusCode(err)
	if code != http.StatusNotFound {
		t.Error("convertToHttpStatusCode is invalid")
	}
}

func TestConvertToHttpStatusCodeWithNotFound(t *testing.T) {
	err := Errors.NotFound{}
	code := convertToHttpStatusCode(err)
	if code != http.StatusNotFound {
		t.Error("convertToHttpStatusCode is invalid")
	}
}

func TestConvertToHttpStatusCodeWithDBConnectionError(t *testing.T) {
	err := Errors.DBConnectionError{}
	code := convertToHttpStatusCode(err)
	if code != http.StatusServiceUnavailable {
		t.Error("convertToHttpStatusCode is invalid")
	}
}

func TestConvertToHttpStatusCodeWithDBOperationError(t *testing.T) {
	err := Errors.DBOperationError{}
	code := convertToHttpStatusCode(err)
	if code != http.StatusServiceUnavailable {
		t.Error("convertToHttpStatusCode is invalid")
	}
}

func TestConvertToHttpStatusCodeWithIOError(t *testing.T) {
	err := Errors.IOError{}
	code := convertToHttpStatusCode(err)
	if code != http.StatusInternalServerError {
		t.Error("convertToHttpStatusCode is invalid")
	}
}

func TestConvertToHttpStatusCodeWithInternalServerError(t *testing.T) {
	err := Errors.InternalServerError{}
	code := convertToHttpStatusCode(err)
	if code != http.StatusInternalServerError {
		t.Error("convertToHttpStatusCode is invalid")
	}
}
