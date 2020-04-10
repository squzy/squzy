package job

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"squzy/apps/internal/helpers"
	"squzy/apps/internal/httpTools"
)

type jobHTTP struct {
	methodType     string
	url            string
	headers        map[string]string
	expectedStatus int32
	timeout        int32
	httpTool       httpTools.HttpTool
}

type httpError struct {
	logId       string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        clientPb.StatusCode
	description string
	location    string
}

func (e *httpError) GetLogData() *clientPb.Log {
	return &clientPb.Log{
		Code:        e.code,
		Description: e.description,
		Meta: &clientPb.MetaData{
			Id:        e.logId,
			Location:  e.location,
			Port:      helpers.GetPortByUrl(e.location),
			StartTime: e.startTime,
			EndTime:   e.endTime,
			Type:      clientPb.Type_Http,
		},
	}
}

func newHttpError(logId string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code clientPb.StatusCode, description string, location string) CheckError {
	return &httpError{
		logId:       logId,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		location:    location,
	}
}

func (j *jobHTTP) Do() CheckError {
	logId := uuid.New().String()
	startTime := ptypes.TimestampNow()
	req := j.httpTool.CreateRequest(j.methodType, j.url, &j.headers, logId)

	_, _, err := j.httpTool.SendRequestTimeoutStatusCode(req, helpers.DurationFromSecond(j.timeout), int(j.expectedStatus))

	if err != nil {
		return newHttpError(
			logId,
			startTime,
			ptypes.TimestampNow(),
			clientPb.StatusCode_Error,
			err.Error(),
			j.url,
		)
	}

	return newHttpError(
		logId,
		startTime,
		ptypes.TimestampNow(),
		clientPb.StatusCode_OK,
		"",
		j.url,
	)
}

func NewHttpJob(method, url string, headers map[string]string, timeout int32, expectedStatus int32, httpTool httpTools.HttpTool) *jobHTTP {
	return &jobHTTP{
		methodType:     method,
		url:            url,
		headers:        headers,
		expectedStatus: expectedStatus,
		timeout:        timeout,
		httpTool:       httpTool,
	}
}
