package messenger

import (
	"bytes"
	"commons/errors"
	"commons/logger"
	"io/ioutil"
	"net/http"
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
	SendHttpRequest(method string, url string, dataOptional ...[]byte) (int, string, error)
}

type HttpRequester struct {
	client httpInterface
}

func NewMessenger() *HttpRequester {
	return &HttpRequester{
		client: httpClient{},
	}
}

// sendHttpRequest creates a new request and sends it to target device.
func (requester HttpRequester) SendHttpRequest(method string, url string, dataOptional ...[]byte) (int, string, error) {
	var err error
	var req *http.Request

	// Make the request with the given method, url, body.
	switch len(dataOptional) {
	case 0:
		req, err = http.NewRequest(method, url, bytes.NewBuffer(nil))
	case 1:
		req, err = http.NewRequest(method, url, bytes.NewBuffer(dataOptional[0]))
	}

	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return http.StatusInternalServerError, "", errors.InternalServerError{err.Error()}
	}

	resp, err := requester.client.DoWrapper(req)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return http.StatusInternalServerError, "", errors.InternalServerError{err.Error()}
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return http.StatusInternalServerError, "", errors.InternalServerError{err.Error()}
	}
	return code, string(responseBody), nil
}
