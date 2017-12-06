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
package messenger

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"io/ioutil"
	msgmocks "messenger/mocks"
	"net/http"
	"testing"
)

func TestCalledSendHttpRequestWithoutData_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpMockObj := msgmocks.NewMockhttpInterface(ctrl)

	gomock.InOrder(
		httpMockObj.EXPECT().DoWrapper(gomock.Any()).Return(&http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(""))}, nil),
	)

	messengerObj := NewMessenger()
	messengerObj.client = httpMockObj

	_, _, err := messengerObj.SendHttpRequest("POST", "/test/url")

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledSendHttpRequestWithData_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpMockObj := msgmocks.NewMockhttpInterface(ctrl)

	gomock.InOrder(
		httpMockObj.EXPECT().DoWrapper(gomock.Any()).Return(&http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(""))}, nil),
	)

	messengerObj := NewMessenger()
	messengerObj.client = httpMockObj

	_, _, err := messengerObj.SendHttpRequest("POST", "/test/url", []byte("data"))

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledSendHttpRequestWhenFailedToSendHttpRequest_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpMockObj := msgmocks.NewMockhttpInterface(ctrl)

	gomock.InOrder(
		httpMockObj.EXPECT().DoWrapper(gomock.Any()).Return(&http.Response{}, errors.New("Error")),
	)

	messengerObj := NewMessenger()
	messengerObj.client = httpMockObj

	_, _, err := messengerObj.SendHttpRequest("POST", "/test/url", []byte("data"))

	if err == nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}
