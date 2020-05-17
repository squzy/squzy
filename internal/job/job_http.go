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
	code        apiPb.SchedulerCode
	description string
}

func (e *httpError) GetLogData() *apiPb.SchedulerResponse {
	var err *apiPb.SchedulerSnapshot_Error
	if e.code == apiPb.SchedulerCode_Error {
		err = &apiPb.SchedulerSnapshot_Error{
			Message: e.description,
		}
	}
	return &apiPb.SchedulerResponse{
		SchedulerId: e.schedulerID,
		Snapshot: &apiPb.SchedulerSnapshot{
			Code:  e.code,
			Error: err,
			Type:  apiPb.SchedulerType_Http,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: e.startTime,
				EndTime:   e.endTime,
			},
		},
	}
}

func newHTTPError(schedulerID string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerCode, description string) CheckError {
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
			apiPb.SchedulerCode_Error,
			err.Error(),
		)
	}

	return newHTTPError(
		schedulerID,
		startTime,
		ptypes.TimestampNow(),
		apiPb.SchedulerCode_OK,
		"",
	)
}
