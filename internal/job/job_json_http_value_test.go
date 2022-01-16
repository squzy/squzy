package job

import (
	"errors"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	structpb "google.golang.org/protobuf/types/known/structpb"
	"net/http"
	"testing"
	"time"
)

type mockError struct {
}

func (m mockError) SendRequestTimeoutStatusCode(req *http.Request, timeout time.Duration, expectedCode int) (int, []byte, error) {
	panic("implement me")
}

func (m mockError) SendRequestTimeout(req *http.Request, timeout time.Duration) (int, []byte, error) {
	return 0, nil, errors.New("afsaf")
}

type mockSuccess struct {
}

func (m mockSuccess) SendRequestTimeoutStatusCode(req *http.Request, timeout time.Duration, expectedCode int) (int, []byte, error) {
	panic("implement me")
}

func (m mockSuccess) SendRequestTimeout(req *http.Request, timeout time.Duration) (int, []byte, error) {
	return 0, []byte(`{"name":"John", "age":31, "city":"New York", "success": true, "time": "2012-04-23T18:25:43.511Z", "raw": {"name":"ahha"}}`), nil
}

func (m mockSuccess) SendRequest(req *http.Request) (int, []byte, error) {
	return 0, []byte(`{"name":"John", "age":31, "city":"New York", "success": true, "time": "2012-04-23T18:25:43.511Z", "raw": {"name":"ahha"}}`), nil
}

func (m mockSuccess) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	panic("implement me")
}

func (m mockSuccess) CreateRequest(method string, url string, headers *map[string]string, logId string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	return req
}

func (m mockError) SendRequest(req *http.Request) (int, []byte, error) {
	return 0, nil, errors.New("afsaf")
}

func (m mockError) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	panic("implement me")
}

func (m mockError) CreateRequest(method string, url string, headers *map[string]string, logId string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	return req
}

func TestExecHttpValue(t *testing.T) {
	t.Run("Should: return error on http request", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}}, &mockError{})
		assert.Equal(t, apiPb.SchedulerCode_ERROR, s.GetLogData().Snapshot.Code)
	})
	t.Run("Should: return error because value not exist", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_STRING,
				Path: "asfasf",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerCode_ERROR, s.GetLogData().Snapshot.Code)
		assert.Equal(t, "", s.GetLogData().Snapshot.Meta.Value.GetStringValue())
	})
	t.Run("Should: not return error because selectors is missing", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerCode_OK, s.GetLogData().Snapshot.Code)
	})
	t.Run("Should: parse single bool value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_BOOL,
				Path: "success",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerCode_OK, s.GetLogData().Snapshot.Code)
		assert.Equal(t, true, s.GetLogData().Snapshot.Meta.Value.GetBoolValue())
	})
	t.Run("Should: parse single string value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_STRING,
				Path: "name",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerCode_OK, s.GetLogData().Snapshot.Code)
		assert.Equal(t, "John", s.GetLogData().Snapshot.Meta.Value.GetStringValue())
	})
	t.Run("Should: parse single number value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_NUMBER,
				Path: "age",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerCode_OK, s.GetLogData().Snapshot.Code)
		assert.Equal(t, float64(31), s.GetLogData().Snapshot.Meta.Value.GetNumberValue())
	})
	t.Run("Should: parse single any value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_ANY,
				Path: "age",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerCode_OK, s.GetLogData().Snapshot.Code)
		assert.Equal(t, "31", s.GetLogData().Snapshot.Meta.Value.GetStringValue())
	})
	t.Run("Should: parse single raw value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_RAW,
				Path: "raw",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerCode_OK, s.GetLogData().Snapshot.Code)
		assert.Equal(t, `{"name":"ahha"}`, s.GetLogData().Snapshot.Meta.Value.GetStringValue())
	})
	t.Run("Should: parse single time value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_TIME,
				Path: "time",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerCode_OK, s.GetLogData().Snapshot.Code)
		assert.Equal(t, "2012-04-23T18:25:43Z", s.GetLogData().Snapshot.Meta.Value.GetStringValue())
	})
	t.Run("Should: parse multipile value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_TIME,
				Path: "time",
			},
			{
				Type: apiPb.HttpJsonValueConfig_NUMBER,
				Path: "age",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerCode_OK, s.GetLogData().Snapshot.Code)
		assert.EqualValues(t, &structpb.ListValue{
			Values: []*structpb.Value{
				{
					Kind: &structpb.Value_StringValue{
						StringValue: "2012-04-23T18:25:43Z",
					},
				},
				{
					Kind: &structpb.Value_NumberValue{
						NumberValue: float64(31),
					},
				},
			},
		}, s.GetLogData().Snapshot.Meta.Value.GetListValue())
	})
}
