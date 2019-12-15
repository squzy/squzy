package httpTools

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type httpTool struct {
	client *http.Client
}

const (
	MaxIdleConnections int = 30
	RequestTimeout     int = 10
)

var (
	notExpectedStatusCode = errors.New("NOT_EXPECTED_STATUS_CODE")
)

func (h *httpTool) SendRequest(req *http.Request) (int, []byte, error) {
	resp, err := h.client.Do(req)

	if err != nil {
		return 0, nil, err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, data, nil
}

func (h *httpTool) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	resp, err := h.client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != expectedCode {
		return resp.StatusCode, nil, notExpectedStatusCode
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, data, nil
}

type HttpTool interface {
	SendRequest(req *http.Request) (int, []byte, error)
	SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error)
}

func New() HttpTool {
	return &httpTool{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: MaxIdleConnections,
			},
			Timeout: time.Duration(RequestTimeout) * time.Second,
		},
	}
}
