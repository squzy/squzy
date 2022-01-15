package job

import (
	"github.com/squzy/squzy/internal/helpers"
	"github.com/squzy/squzy/internal/httptools"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
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
	if e.code == apiPb.SchedulerCode_ERROR {
		err = &apiPb.SchedulerSnapshot_Error{
			Message: e.description,
		}
	}
	return &apiPb.SchedulerResponse{
		SchedulerId: e.schedulerID,
		Snapshot: &apiPb.SchedulerSnapshot{
			Code:  e.code,
			Error: err,
			Type:  apiPb.SchedulerType_HTTP,
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
	startTime := timestamp.Now()
	req := httpTool.CreateRequest(config.Method, config.URL, &config.Headers, schedulerID)

	_, _, err := httpTool.SendRequestTimeoutStatusCode(req, helpers.DurationFromSecond(timeout), int(config.StatusCode))

	if err != nil {
		return newHTTPError(
			schedulerID,
			startTime,
			timestamp.Now(),
			apiPb.SchedulerCode_ERROR,
			err.Error(),
		)
	}

	return newHTTPError(
		schedulerID,
		startTime,
		timestamp.Now(),
		apiPb.SchedulerCode_OK,
		"",
	)
}
