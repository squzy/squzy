package job

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"squzy/internal/helpers"
	"squzy/internal/httpTools"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
)

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

func ExecHttp(schedulerId string, timeout int32, config *scheduler_config_storage.HttpConfig, httpTool httpTools.HttpTool) CheckError {
	startTime := ptypes.TimestampNow()
	req := httpTool.CreateRequest(config.Method, config.Url, &config.Headers, schedulerId)

	_, _, err := httpTool.SendRequestTimeoutStatusCode(req, helpers.DurationFromSecond(timeout), int(config.StatusCode))

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
