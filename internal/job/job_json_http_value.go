package job

import (
	"fmt"
	"github.com/squzy/squzy/internal/helpers"
	"github.com/squzy/squzy/internal/httptools"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/tidwall/gjson"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type jsonHTTPError struct {
	schedulerID string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerCode
	description string
	value       *structpb.Value
}

var (
	valueNotExistErrorFn = func(path string) error {
		return fmt.Errorf("value by path=`%s` not exist", path)
	}
)

func (e *jsonHTTPError) GetLogData() *apiPb.SchedulerResponse {
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
			Type:  apiPb.SchedulerType_HTTP_JSON_VALUE,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: e.startTime,
				EndTime:   e.endTime,
				Value:     e.value,
			},
		},
	}
}

func ExecHTTPValue(schedulerID string, timeout int32, config *scheduler_config_storage.HTTPValueConfig, httpTool httptools.HTTPTool) CheckError {
	startTime := timestamp.Now()
	req := httpTool.CreateRequest(config.Method, config.URL, &config.Headers, schedulerID)

	_, data, err := httpTool.SendRequestTimeout(req, helpers.DurationFromSecond(timeout))

	if err != nil {
		return newJSONHTTPError(
			schedulerID,
			startTime,
			timestamp.Now(),
			apiPb.SchedulerCode_ERROR,
			err.Error(),
			nil,
		)
	}

	jsonString := string(data)

	results := []*structpb.Value{}

	if len(config.Selectors) == 0 {
		return newJSONHTTPError(
			schedulerID,
			startTime,
			timestamp.Now(),
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
				timestamp.Now(),
				apiPb.SchedulerCode_ERROR,
				valueNotExistErrorFn(value.Path).Error(),
				nil,
			)
		}
		switch value.Type {
		case apiPb.HttpJsonValueConfig_STRING:
			results = append(results, &structpb.Value{
				Kind: &structpb.Value_StringValue{
					StringValue: res.String(),
				},
			})
		case apiPb.HttpJsonValueConfig_BOOL:
			results = append(results, &structpb.Value{
				Kind: &structpb.Value_BoolValue{
					BoolValue: res.Bool(),
				},
			})
		case apiPb.HttpJsonValueConfig_NUMBER:
			results = append(results, &structpb.Value{
				Kind: &structpb.Value_NumberValue{
					NumberValue: res.Float(),
				},
			})
		case apiPb.HttpJsonValueConfig_TIME:
			results = append(results, &structpb.Value{
				Kind: &structpb.Value_StringValue{
					StringValue: res.Time().Format(time.RFC3339),
				},
			})
		case apiPb.HttpJsonValueConfig_ANY:
			results = append(results, &structpb.Value{
				Kind: &structpb.Value_StringValue{
					StringValue: fmt.Sprintf("%v", res.Value()),
				},
			})
		case apiPb.HttpJsonValueConfig_RAW:
			results = append(results, &structpb.Value{
				Kind: &structpb.Value_StringValue{
					StringValue: res.Raw,
				},
			})
		}
	}

	if len(config.Selectors) == 1 {
		return newJSONHTTPError(
			schedulerID,
			startTime,
			timestamp.Now(),
			apiPb.SchedulerCode_OK,
			"",
			results[0],
		)
	}

	return newJSONHTTPError(
		schedulerID,
		startTime,
		timestamp.Now(),
		apiPb.SchedulerCode_OK,
		"",
		&structpb.Value{
			Kind: &structpb.Value_ListValue{
				ListValue: &structpb.ListValue{
					Values: results,
				},
			},
		},
	)
}

func newJSONHTTPError(schedulerID string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerCode, description string, value *structpb.Value) CheckError {
	return &jsonHTTPError{
		schedulerID: schedulerID,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		value:       value,
	}
}
