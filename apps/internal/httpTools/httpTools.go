package httpTools

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"squzy/apps/internal/helpers"
	"time"
)

type httpTool struct {
	client *http.Client
}

const (
	MaxIdleConnections        int = 30
	MaxIdleConnectionsPerHost int = 30
	RequestTimeout            int = 10
)

var (
	notExpectedStatusCode   = errors.New("NOT_EXPECTED_STATUS_CODE")
	notExpectedStatusCodeFn = func(url string, statusCode int, expectedStatusCode int) error {
		return errors.New(
			fmt.Sprintf(
				"ErrCode: %s, Location: %s, StatusCoe: %d, ExpectedStatusCode: %d, Port: %d",
				notExpectedStatusCode,
				url,
				statusCode,
				expectedStatusCode,
				helpers.GetPortByUrl(url),
			),
		)
	}
)

func (h *httpTool) CreateRequest(method string, url string, headers *map[string]string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	if headers == nil {
		return req
	}
	for k, v := range *headers {
		req.Header.Set(k, v)
	}
	return req
}

func (h *httpTool) SendRequest(req *http.Request) (int, []byte, error) {
	return h.sendReq(req, false, 0)
}

func (h *httpTool) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	return h.sendReq(req, true, expectedCode)
}

type HttpTool interface {
	SendRequest(req *http.Request) (int, []byte, error)
	SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error)
	CreateRequest(method string, url string, headers *map[string]string) *http.Request
}

func (h *httpTool) sendReq(req *http.Request, checkCode bool, statusCode int) (int, []byte, error) {
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

	if checkCode {
		if statusCode != resp.StatusCode {
			return resp.StatusCode, nil, notExpectedStatusCodeFn(req.URL.String(), resp.StatusCode, statusCode)
		}
		return resp.StatusCode, data, nil
	}

	return resp.StatusCode, data, nil
}

func New() HttpTool {
	return &httpTool{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: MaxIdleConnectionsPerHost,
				MaxIdleConns:        MaxIdleConnections,
			},
			Timeout: time.Duration(RequestTimeout) * time.Second,
		},
	}
}
