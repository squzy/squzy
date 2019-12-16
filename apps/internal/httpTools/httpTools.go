package httpTools

import (
	"errors"
	"github.com/valyala/fasthttp"
	"time"
)

type httpTool struct {
	client *fasthttp.Client
}

const (
	MaxIdleConnections int = 30
	RequestTimeout     int = 10
)

var (
	notExpectedStatusCode = errors.New("NOT_EXPECTED_STATUS_CODE")
)

func (h *httpTool) CreateRequest(method string, url string, headers *map[string]string) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	if headers == nil {
		return req
	}
	for k, v := range *headers {
		req.Header.Set(k, v)
	}
	return req
}

func (h *httpTool) SendRequest(req *fasthttp.Request) (int, []byte, error) {
	return h.sendReq(req, false, 0)
}

func (h *httpTool) SendRequestWithStatusCode(req *fasthttp.Request, expectedCode int) (int, []byte, error) {
	return h.sendReq(req, true, expectedCode)
}

type HttpTool interface {
	SendRequest(req *fasthttp.Request) (int, []byte, error)
	SendRequestWithStatusCode(req *fasthttp.Request, expectedCode int) (int, []byte, error)
	CreateRequest(method string, url string, headers *map[string]string) *fasthttp.Request
}

func (h *httpTool) sendReq(req *fasthttp.Request, checkCode bool, statusCode int) (int, []byte, error) {
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	err := h.client.Do(req, resp)

	if err != nil {
		return 0, nil, err
	}

	if checkCode {
		if statusCode != resp.StatusCode() {
			return resp.StatusCode(), nil, notExpectedStatusCode
		}
		return resp.StatusCode(), resp.Body(), nil
	}

	bodyBytes := resp.Body()
	return resp.StatusCode(), bodyBytes, nil
}

func New() HttpTool {
	return &httpTool{
		client: &fasthttp.Client{
			ReadTimeout:                   time.Duration(RequestTimeout) * time.Second,
			DisableHeaderNamesNormalizing: true,
			MaxIdleConnDuration:           time.Duration(MaxIdleConnections) * time.Second,
			MaxConnsPerHost:               10000,
			ReadBufferSize:                4096,
			NoDefaultUserAgentHeader:      true,
		},
	}
}
