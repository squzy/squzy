package job

import (
	"github.com/PxyUp/fitter/lib"
	fconfig "github.com/PxyUp/fitter/pkg/config"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"google.golang.org/protobuf/types/known/structpb"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
)

type jsonHTMLError struct {
	schedulerID string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerCode
	description string
	value       *structpb.Value
}

//var (
//	valueNotExistErrorFn = func(path string) error {
//		return fmt.Errorf("value by path=`%s` not exist", path)
//	}
//)

func (e *jsonHTMLError) GetLogData() *apiPb.SchedulerResponse {
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

func ExecHTML(schedulerID string, timeout int32, config *scheduler_config_storage.HTMLValueConfig) CheckError {
	startTime := timestamp.Now()

	var fields map[string]*fconfig.Field
	for i, f := range config.Selectors {
		fields[strconv.Itoa(i)] = &fconfig.Field{
			BaseField: &fconfig.BaseField{
				Type: scheduler_config_storage.FieldToHtml[f.Type],
				Path: f.Path,
			},
		}
	}

	res, err := lib.Parse(&fconfig.Item{
		ConnectorConfig: &fconfig.ConnectorConfig{
			ResponseType: fconfig.HTML,
			Url:          config.URL,
			ServerConfig: &fconfig.ServerConnectorConfig{
				Method:  config.Method,
				Headers: config.Headers,
				Timeout: uint32(timeout),
			},
		},
		Model: &fconfig.Model{
			ObjectConfig: &fconfig.ObjectConfig{
				Fields: fields,
			},
		},
	}, nil, nil)

	if err != nil {
		return newHTMLError(
			schedulerID,
			startTime,
			timestamp.Now(),
			apiPb.SchedulerCode_ERROR,
			err.Error(),
			nil,
		)
	}

	return newHTMLError(
		schedulerID,
		startTime,
		timestamp.Now(),
		apiPb.SchedulerCode_OK,
		"",
		&structpb.Value{
			Kind: &structpb.Value_ListValue{
				ListValue: &structpb.ListValue{
					Values: []*structpb.Value{
						{
							Kind: &structpb.Value_StringValue{
								StringValue: res.ToJson(),
							},
						},
					},
				},
			},
		},
	)
}

func newHTMLError(schedulerID string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerCode, description string, value *structpb.Value) CheckError {
	return &jsonHTMLError{
		schedulerID: schedulerID,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		value:       value,
	}
}
