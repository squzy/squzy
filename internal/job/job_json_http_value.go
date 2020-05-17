package job

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	structType "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/tidwall/gjson"
	"squzy/internal/helpers"
	"squzy/internal/httptools"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
	"time"
)

type jsonHTTPError struct {
	schedulerID string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerCode
	description string
	value       *structType.Value
}

var (
	valueNotExistErrorFn = func(path string) error {
		return fmt.Errorf("value by path=`%s` not exist", path)
	}
)

func (e *jsonHTTPError) GetLogData() *apiPb.SchedulerResponse {
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
			Type:  apiPb.SchedulerType_HttpJsonValue,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: e.startTime,
				EndTime:   e.endTime,
				Value:     e.value,
			},
		},
	}
}

func ExecHTTPValue(schedulerID string, timeout int32, config *scheduler_config_storage.HTTPValueConfig, httpTool httptools.HTTPTool) CheckError {
	startTime := ptypes.TimestampNow()
	req := httpTool.CreateRequest(config.Method, config.URL, &config.Headers, schedulerID)

	_, data, err := httpTool.SendRequestTimeout(req, helpers.DurationFromSecond(timeout))

	if err != nil {
		return newJSONHTTPError(
			schedulerID,
			startTime,
			ptypes.TimestampNow(),
			apiPb.SchedulerCode_Error,
			err.Error(),
			nil,
		)
	}

	jsonString := string(data)

	results := []*structType.Value{}

	if len(config.Selectors) == 0 {
		return newJSONHTTPError(
			schedulerID,
			startTime,
			ptypes.TimestampNow(),
			apiPb.SchedulerCode_OK,
			"",
			nil,
		)
	}

	for _, value := range config.Selectors {
		res := gjson.Get(jsonString, value.Path)
		if !res.Exists() {
			return newJSONHTTPError(
				schedulerID,
				startTime,
				ptypes.TimestampNow(),
				apiPb.SchedulerCode_Error,
				valueNotExistErrorFn(value.Path).Error(),
				nil,
			)
		}
		switch value.Type {
		case apiPb.HttpJsonValueConfig_String:
			results = append(results, &structType.Value{
				Kind: &structType.Value_StringValue{
					StringValue: res.String(),
				},
			})
		case apiPb.HttpJsonValueConfig_Bool:
			results = append(results, &structType.Value{
				Kind: &structType.Value_BoolValue{
					BoolValue: res.Bool(),
				},
			})
		case apiPb.HttpJsonValueConfig_Number:
			results = append(results, &structType.Value{
				Kind: &structType.Value_NumberValue{
					NumberValue: res.Float(),
				},
			})
		case apiPb.HttpJsonValueConfig_Time:
			results = append(results, &structType.Value{
				Kind: &structType.Value_StringValue{
					StringValue: res.Time().Format(time.RFC3339),
				},
			})
		case apiPb.HttpJsonValueConfig_Any:
			results = append(results, &structType.Value{
				Kind: &structType.Value_StringValue{
					StringValue: fmt.Sprintf("%v", res.Value()),
				},
			})
		case apiPb.HttpJsonValueConfig_Raw:
			results = append(results, &structType.Value{
				Kind: &structType.Value_StringValue{
					StringValue: res.Raw,
				},
			})
		}
	}

	if len(config.Selectors) == 1 {
		return newJSONHTTPError(
			schedulerID,
			startTime,
			ptypes.TimestampNow(),
			apiPb.SchedulerCode_OK,
			"",
			results[0],
		)
	}

	return newJSONHTTPError(
		schedulerID,
		startTime,
		ptypes.TimestampNow(),
		apiPb.SchedulerCode_OK,
		"",
		&structType.Value{
			Kind: &structType.Value_ListValue{
				ListValue: &structType.ListValue{
					Values: results,
				},
			},
		},
	)
}

func newJSONHTTPError(schedulerID string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerCode, description string, value *structType.Value) CheckError {
	return &jsonHTTPError{
		schedulerID: schedulerID,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		value:       value,
	}
}
