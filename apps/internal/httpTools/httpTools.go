package httpTools

import (
	"errors"
	"io/ioutil"
	"net/http"
)

type httpTool struct {
	client http.Client
}

var (
	notExpectedStatusCode = errors.New("NOT_EXPECTED_STATUS_CODE")
)

func (h *httpTool) SendRequest(req *http.Request) ([]byte, error) {
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (h *httpTool) SendRequestWithStatusCode(req *http.Request, expectedCode int) ([]byte, error) {
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != expectedCode {
		return nil, notExpectedStatusCode
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type HttpTool interface {
	SendRequest(req *http.Request) ([]byte, error)
	SendRequestWithStatusCode(req *http.Request, expectedCode int) ([]byte, error)
}

func New(client http.Client) HttpTool {
	return &httpTool{
		client: client,
	}
}
