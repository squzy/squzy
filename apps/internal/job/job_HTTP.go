package job

import (
	"errors"
	"net/http"
	"time"
)

type jobHTTP struct {
	methodType string
	url        string
	headers    map[string]string
	statusCode int
}

const (
	timeout = 5 * time.Second
)

var (
	wrongStatusError = errors.New("WRONG_STATUS_CODE")
)

func (j *jobHTTP) Do() error {
	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(j.methodType, j.url, nil)
	if err != nil {
		return err
	}

	for name, val := range j.headers {
		req.Header.Set(name, val)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != j.statusCode {
		return wrongStatusError
	}

	return nil
}

func NewJob(method, url string, headers map[string]string, status int) *jobHTTP {
	return &jobHTTP{
		methodType: method,
		url:        url,
		headers:    headers,
		statusCode: status,
	}
}