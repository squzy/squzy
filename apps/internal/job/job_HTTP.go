package job

import (
	"errors"
	"fmt"
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
	timeout          time.Duration = 5
	errorWrongStatus string        = "WRONG_STATUS_CODE"
	errorNoResponse string         = "No_RESPONSE_FROM_SERVER"
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
	if resp == nil {
		return errors.New(errorNoResponse)
	}

	fmt.Println(resp.StatusCode)

	defer resp.Body.Close()

	if resp.StatusCode == j.statusCode {
		return nil
	}
	return errors.New(errorWrongStatus)
}

func NewJob(method, url string, headers map[string]string, status int) *jobHTTP {
	return &jobHTTP{
		methodType: method,
		url:        url,
		headers:    headers,
		statusCode: status,
	}
}