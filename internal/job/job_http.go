package job

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"squzy/internal/helpers"
	"squzy/internal/httptools"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
)

type httpError struct {
	schedulerID string
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
		SchedulerId: e.schedulerID,
		Code:        e.code,
		Error:       err,
		Type:        apiPb.SchedulerType_Http,
		Meta: &apiPb.SchedulerResponse_MetaData{
			StartTime: e.startTime,
			EndTime:   e.endTime,
		},
	}
}

func newHTTPError(schedulerID string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerResponseCode, description string) CheckError {
	return &httpError{
		schedulerID: schedulerID,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
	}
}

func ExecHTTP(schedulerID string, timeout int32, config *scheduler_config_storage.HTTPConfig, httpTool httptools.HTTPTool) CheckError {
	startTime := ptypes.TimestampNow()
	req := httpTool.CreateRequest(config.Method, config.URL, &config.Headers, schedulerID)

	_, _, err := httpTool.SendRequestTimeoutStatusCode(req, helpers.DurationFromSecond(timeout), int(config.StatusCode))

	if err != nil {
		return newHTTPError(
			schedulerID,
			startTime,
			ptypes.TimestampNow(),
			apiPb.SchedulerResponseCode_Error,
			err.Error(),
		)
	}

	return newHTTPError(
		schedulerID,
		startTime,
		ptypes.TimestampNow(),
		apiPb.SchedulerResponseCode_OK,
		"",
	)
}
