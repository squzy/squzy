package job

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	structType "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	httpPb "github.com/squzy/squzy_generated/generated/server/proto/v1"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/tidwall/gjson"
	"squzy/apps/internal/helpers"
	"squzy/apps/internal/httpTools"
	"time"
)

type jsonHttpValueJob struct {
	method    string
	url       string
	headers   map[string]string
	httpTool  httpTools.HttpTool
	selectors []*httpPb.HttpJsonValueCheck_Selectors
}

type jsonHttpError struct {
	logId       string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        clientPb.StatusCode
	description string
	location    string
	value       *structType.Value
}

var (
	valueNotExistErrorFn = func(path string) error {
		return errors.New(fmt.Sprintf("Value by path=`%s` not exist", path))
	}
)

func (e *jsonHttpError) GetLogData() *clientPb.Log {
	return &clientPb.Log{
		Code:        e.code,
		Description: e.description,
		Meta: &clientPb.MetaData{
			Id:        e.logId,
			Location:  e.location,
			Port:      helpers.GetPortByUrl(e.location),
			StartTime: e.startTime,
			EndTime:   e.endTime,
			Type:      clientPb.Type_HttpJsonValue,
		},
		Value: e.value,
	}
}

func (j *jsonHttpValueJob) Do() CheckError {
	logId := uuid.New().String()
	startTime := ptypes.TimestampNow()
	req := j.httpTool.CreateRequest(j.method, j.url, &j.headers, logId)

	_, data, err := j.httpTool.SendRequest(req)

	if err != nil {
		return newJsonHttpError(
			logId,
			startTime,
			ptypes.TimestampNow(),
			clientPb.StatusCode_Error,
			err.Error(),
			j.url,
			nil,
		)
	}

	jsonString := string(data)

	results := []*structType.Value{}

	if len(j.selectors) == 0 {
		return newJsonHttpError(
			logId,
			startTime,
			ptypes.TimestampNow(),
			clientPb.StatusCode_OK,
			"",
			j.url,
			nil,
		)
	}

	for _, value := range j.selectors {
		res := gjson.Get(jsonString, value.Path)
		if !res.Exists() {
			return newJsonHttpError(
				logId,
				startTime,
				ptypes.TimestampNow(),
				clientPb.StatusCode_Error,
				valueNotExistErrorFn(value.Path).Error(),
				j.url,
				nil,
			)
		}
		switch value.Type {
		case httpPb.HttpJsonValueCheck_String:
			results = append(results, &structType.Value{
				Kind: &structType.Value_StringValue{
					StringValue: res.String(),
				},
			})
		case httpPb.HttpJsonValueCheck_Bool:
			results = append(results, &structType.Value{
				Kind: &structType.Value_BoolValue{
					BoolValue: res.Bool(),
				},
			})
		case httpPb.HttpJsonValueCheck_Number:
			results = append(results, &structType.Value{
				Kind: &structType.Value_NumberValue{
					NumberValue: res.Float(),
				},
			})
		case httpPb.HttpJsonValueCheck_Time:
			results = append(results, &structType.Value{
				Kind: &structType.Value_StringValue{
					StringValue: res.Time().Format(time.RFC3339),
				},
			})
		case httpPb.HttpJsonValueCheck_Any:
			results = append(results, &structType.Value{
				Kind: &structType.Value_StringValue{
					StringValue: fmt.Sprintf("%v", res.Value()),
				},
			})
		case httpPb.HttpJsonValueCheck_Raw:
			results = append(results, &structType.Value{
				Kind: &structType.Value_StringValue{
					StringValue: res.Raw,
				},
			})
		}
	}

	if len(j.selectors) == 1 {
		return newJsonHttpError(
			logId,
			startTime,
			ptypes.TimestampNow(),
			clientPb.StatusCode_OK,
			"",
			j.url,
			results[0],
		)
	}

	return newJsonHttpError(
		logId,
		startTime,
		ptypes.TimestampNow(),
		clientPb.StatusCode_OK,
		"",
		j.url,
		&structType.Value{
			Kind: &structType.Value_ListValue{
				ListValue: &structType.ListValue{
					Values: results,
				},
			},
		},
	)
}

func newJsonHttpError(logId string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code clientPb.StatusCode, description string, location string, value *structType.Value) CheckError {
	return &jsonHttpError{
		logId:       logId,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		location:    location,
		value:       value,
	}
}

func NewJsonHttpValueJob(method, url string, headers map[string]string, httpTool httpTools.HttpTool, selectors []*httpPb.HttpJsonValueCheck_Selectors) *jsonHttpValueJob {
	return &jsonHttpValueJob{
		method:    method,
		url:       url,
		headers:   headers,
		httpTool:  httpTool,
		selectors: selectors,
	}
}
