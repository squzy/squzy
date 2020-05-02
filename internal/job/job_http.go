package job

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"squzy/internal/helpers"
	"squzy/internal/httpTools"
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
	schedulerId string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerResponseCode
	description string
}

func (e *httpError) GetLogData() *apiPb.SchedulerResponse {
	var err *apiPb.SchedulerResponse_Error
	if e.code == apiPb.SchedulerResponseCode_Error {
		err = &apiPb.SchedulerResponse_Error{
			Message: e.description,
		}
	}
	return &apiPb.SchedulerResponse{
		SchedulerId: e.schedulerId,
		Code:        e.code,
		Error:       err,
		Type:        apiPb.SchedulerType_Http,
		Meta: &apiPb.SchedulerResponse_MetaData{
			StartTime: e.startTime,
			EndTime:   e.endTime,
		},
	}
}

func newHttpError(schedulerId string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerResponseCode, description string) CheckError {
	return &httpError{
		schedulerId: schedulerId,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
	}
}

func (j *jobHTTP) Do(schedulerId string) CheckError {
	startTime := ptypes.TimestampNow()
	req := j.httpTool.CreateRequest(j.methodType, j.url, &j.headers, schedulerId)

	_, _, err := j.httpTool.SendRequestTimeoutStatusCode(req, helpers.DurationFromSecond(j.timeout), int(j.expectedStatus))

	if err != nil {
		return newHttpError(
			schedulerId,
			startTime,
			ptypes.TimestampNow(),
			apiPb.SchedulerResponseCode_Error,
			err.Error(),
		)
	}

	return newHttpError(
		schedulerId,
		startTime,
		ptypes.TimestampNow(),
		apiPb.SchedulerResponseCode_OK,
		"",
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
