package job

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	structType "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/tidwall/gjson"
	"squzy/internal/helpers"
	"squzy/internal/httpTools"
	"time"
)

type jsonHttpValueJob struct {
	method    string
	url       string
	timeout   int32
	headers   map[string]string
	httpTool  httpTools.HttpTool
	selectors []*apiPb.HttpJsonValueConfig_Selectors
}

type jsonHttpError struct {
	schedulerId string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerResponseCode
	description string
	value       *structType.Value
}

var (
	valueNotExistErrorFn = func(path string) error {
		return errors.New(fmt.Sprintf("Value by path=`%s` not exist", path))
	}
)

func (e *jsonHttpError) GetLogData() *apiPb.SchedulerResponse {
	var err *apiPb.SchedulerResponse_Error
	if e.code == apiPb.SchedulerResponseCode_Error {
		err = &apiPb.SchedulerResponse_Error{
			Message: e.description,
		}
	}
	return &apiPb.SchedulerResponse{
		SchedulerId: e.schedulerId,
		Code:  e.code,
		Error: err,
		Type:  apiPb.SchedulerType_HttpJsonValue,
		Meta: &apiPb.SchedulerResponse_MetaData{
			StartTime: e.startTime,
			EndTime:   e.endTime,
			Value:     e.value,
		},
	}
}

func (j *jsonHttpValueJob) Do(schedulerId string) CheckError {
	startTime := ptypes.TimestampNow()
	req := j.httpTool.CreateRequest(j.method, j.url, &j.headers, schedulerId)

	_, data, err := j.httpTool.SendRequestTimeout(req, helpers.DurationFromSecond(j.timeout))

	if err != nil {
		return newJsonHttpError(
			schedulerId,
			startTime,
			ptypes.TimestampNow(),
			apiPb.SchedulerResponseCode_Error,
			err.Error(),
			nil,
		)
	}

	jsonString := string(data)

	results := []*structType.Value{}

	if len(j.selectors) == 0 {
		return newJsonHttpError(
			schedulerId,
			startTime,
			ptypes.TimestampNow(),
			apiPb.SchedulerResponseCode_OK,
			"",
			nil,
		)
	}

	for _, value := range j.selectors {
		res := gjson.Get(jsonString, value.Path)
		if !res.Exists() {
			return newJsonHttpError(
				schedulerId,
				startTime,
				ptypes.TimestampNow(),
				apiPb.SchedulerResponseCode_Error,
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

	if len(j.selectors) == 1 {
		return newJsonHttpError(
			schedulerId,
			startTime,
			ptypes.TimestampNow(),
			apiPb.SchedulerResponseCode_OK,
			"",
			results[0],
		)
	}

	return newJsonHttpError(
		schedulerId,
		startTime,
		ptypes.TimestampNow(),
		apiPb.SchedulerResponseCode_OK,
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

func newJsonHttpError(schedulerId string,startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerResponseCode, description string, value *structType.Value) CheckError {
	return &jsonHttpError{
		schedulerId: schedulerId,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		value:       value,
	}
}

func NewJsonHttpValueJob(method, url string, headers map[string]string, timeout int32, httpTool httpTools.HttpTool, selectors []*apiPb.HttpJsonValueConfig_Selectors) *jsonHttpValueJob {
	return &jsonHttpValueJob{
		method:    method,
		url:       url,
		headers:   headers,
		timeout:   timeout,
		httpTool:  httpTool,
		selectors: selectors,
	}
}
