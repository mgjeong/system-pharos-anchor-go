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
	"commons/logger"
	"net/http"
	"sort"
	"sync"
)

type httpInterface interface {
	DoWrapper(req *http.Request) (*http.Response, error)
}

type httpClient struct{}

// DoWrapper is a wrapper around DefaultClient.Do.
func (httpClient) DoWrapper(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

type MessengerInterface interface {
	SendHttpRequest(method string, urls []string, dataOptional ...[]byte) ([]int, []string)
}

type HttpRequester struct {
	client httpInterface
}

func NewMessenger() *HttpRequester {
	return &HttpRequester{
		client: httpClient{},
	}
}

// A httpResponse represents an HTTP response received from remote device.
type httpResponse struct {
	index int
	resp  *http.Response
	err   string
}

type sortRespSlice []httpResponse

// Len returns length of httpResponse.
func (arr sortRespSlice) Len() int {
	return len(arr)
}

// Less returns whether the its first argument compares less than the second.
func (arr sortRespSlice) Less(i, j int) bool {
	return arr[i].index < arr[j].index
}

// Swap exchange its first argument with the second.
func (arr sortRespSlice) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

// sendHttpRequest creates a new request and sends it to target device.
func (requester HttpRequester) SendHttpRequest(method string, urls []string, dataOptional ...[]byte) ([]int, []string) {
	var wg sync.WaitGroup
	wg.Add(len(urls))

	respChannel := make(chan httpResponse, len(urls))
	for i := range urls {
		go func(idx int) {
			logger.Logging(logger.DEBUG, "sending http request:", urls[idx])

			var err error
			var req *http.Request
			var resp httpResponse

			resp.index = idx

			switch len(dataOptional) {
			case 0:
				req, err = http.NewRequest(method, urls[idx], bytes.NewBuffer(nil))
			case 1:
				req, err = http.NewRequest(method, urls[idx], bytes.NewBuffer([]byte(dataOptional[0])))
			}

			if err != nil {
				resp.resp = nil
				resp.err = err.Error()
				respChannel <- resp
			} else {
				resp.resp, err = requester.client.DoWrapper(req)
				if err != nil {
					resp.err = err.Error()
				} else {
					resp.err = ""
				}
				respChannel <- resp
			}
			defer wg.Done()
		}(i)
	}
	wg.Wait()

	var respList []httpResponse
	for range urls {
		respList = append(respList, <-respChannel)
	}
	sort.Sort(sortRespSlice(respList))
	return changeToReturnValue(respList)
}

// changeToReturnValue parses a response code and body from httpResponse structure.
func changeToReturnValue(respList []httpResponse) (respCode []int, respBody []string) {
	var buf bytes.Buffer

	for i := 0; i < len(respList); i++ {
		buf.Reset()
		if respList[i].resp == nil {
			message := `{"message":"` + respList[i].err + `"}`
			respBody = append(respBody, message)
			respCode = append(respCode, 500)
		} else {
			buf.ReadFrom(respList[i].resp.Body)
			respBody = append(respBody, buf.String())
			respCode = append(respCode, respList[i].resp.StatusCode)
		}
	}
	return respCode, respBody
}
