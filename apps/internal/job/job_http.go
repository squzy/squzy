package job

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"squzy/apps/internal/httpTools"
)

type jobHTTP struct {
	methodType     string
	url            string
	headers        map[string]string
	expectedStatus int32
	httpTool       httpTools.HttpTool
}

type httpError struct {
	time        *timestamp.Timestamp
	code        clientPb.StatusCode
	description string
	location    string
}

func (e *httpError) GetLogData() *clientPb.Log {
	return &clientPb.Log{
		Code:        e.code,
		Description: e.description,
		Meta: &clientPb.MetaData{
			Id:       uuid.New().String(),
			Location: e.location,
			Port:     GetPortByUrl(e.location),
			Time:     e.time,
			Type:     clientPb.Type_Http,
		},
	}
}

func newHttpError(time *timestamp.Timestamp, code clientPb.StatusCode, description string, location string) CheckError {
	return &httpError{
		time:        time,
		code:        code,
		description: description,
		location:    location,
	}
}

func (j *jobHTTP) Do() CheckError {
	req := j.httpTool.CreateRequest(j.methodType, j.url, &j.headers)

	_, _, err := j.httpTool.SendRequestWithStatusCode(req, int(j.expectedStatus))

	if err != nil {
		return newHttpError(
			ptypes.TimestampNow(),
			clientPb.StatusCode_Error,
			err.Error(),
			j.url,
		)
	}

	return newHttpError(
		ptypes.TimestampNow(),
		clientPb.StatusCode_OK,
		"",
		j.url,
	)
}

func NewHttpJob(method, url string, headers map[string]string, expectedStatus int32, httpTool httpTools.HttpTool) *jobHTTP {
	return &jobHTTP{
		methodType:     method,
		url:            url,
		headers:        headers,
		expectedStatus: expectedStatus,
		httpTool:       httpTool,
	}
}
