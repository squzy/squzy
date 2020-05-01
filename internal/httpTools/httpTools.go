package httpTools

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"squzy/internal/helpers"
	"time"
)

type httpTool struct {
	userAgent string
	client    *http.Client
}

const (
	MaxIdleConnections        int   = 30
	MaxIdleConnectionsPerHost int   = 30
	RequestTimeout            int32 = 10
	userAgentPrefix                 = "Squzy_monitoring"
	logHeader                       = "Squzy_log_id"
	userAgentHeaderKey              = "User-Agent"
)

var (
	notExpectedStatusCode   = errors.New("NOT_EXPECTED_STATUS_CODE")
	notExpectedStatusCodeFn = func(url string, statusCode int, expectedStatusCode int) error {
		return errors.New(
			fmt.Sprintf(
				"ErrCode: %s, Location: %s, StatusCode: %d, ExpectedStatusCode: %d, Port: %d",
				notExpectedStatusCode,
				url,
				statusCode,
				expectedStatusCode,
				helpers.GetPortByUrl(url),
			),
		)
	}
	defaultTimeout = helpers.DurationFromSecond(10)
)

type HttpTool interface {
	SendRequest(req *http.Request) (int, []byte, error)
	SendRequestTimeout(req *http.Request, timeout time.Duration) (int, []byte, error)
	SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error)
	SendRequestTimeoutStatusCode(req *http.Request, timeout time.Duration, expectedCode int, ) (int, []byte, error)
	CreateRequest(method string, url string, headers *map[string]string, logId string) *http.Request
}

func (h *httpTool) CreateRequest(method string, url string, headers *map[string]string, logId string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)

	// Set user agent
	req.Header.Set(userAgentHeaderKey, h.userAgent)

	if logId != "" {
		req.Header.Set(logHeader, logId)
	}

	if headers == nil {
		return req
	}
	for k, v := range *headers {
		req.Header.Set(k, v)
	}

	return req
}

func (h *httpTool) SendRequest(req *http.Request) (int, []byte, error) {
	return sendReq(h.client, req, false, 0)
}

func (h *httpTool) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	return sendReq(h.client, req, true, expectedCode)
}

func (h *httpTool) SendRequestTimeout(req *http.Request, timeout time.Duration) (int, []byte, error) {
	return h.sendRequestTimeout(req, timeout, false, 0)
}

func (h *httpTool) SendRequestTimeoutStatusCode(req *http.Request, timeout time.Duration, expectedCode int) (int, []byte, error) {
	return h.sendRequestTimeout(req, timeout, true, expectedCode)
}

func (h *httpTool) sendRequestTimeout(req *http.Request, timeout time.Duration, checkCode bool, code int) (int, []byte, error) {
	// If timeout not present will be use method with custom http client
	if timeout.Seconds() <= 0 {
		return sendReq(h.client, req, checkCode, code)
	}
	ctx, cancel := helpers.TimeoutContext(context.Background(), timeout)
	defer cancel()
	reqTimeout := req.WithContext(ctx)
	return sendReq(http.DefaultClient, reqTimeout, checkCode, code)
}

func sendReq(client *http.Client, req *http.Request, checkCode bool, statusCode int) (int, []byte, error) {
	resp, err := client.Do(req)

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

func getUserAgent(version string) string {
	return fmt.Sprintf("%s_%s", userAgentPrefix, version)
}

func New(userAgentVersion string) HttpTool {
	return &httpTool{
		userAgent: getUserAgent(userAgentVersion),
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: MaxIdleConnectionsPerHost,
				MaxIdleConns:        MaxIdleConnections,
			},
			Timeout: defaultTimeout,
		},
	}
}
