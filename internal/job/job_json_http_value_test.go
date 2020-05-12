package job

import (
	"errors"
	structType "github.com/golang/protobuf/ptypes/struct"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
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
		assert.Equal(t, apiPb.SchedulerResponseCode_Error, s.GetLogData().Code)
	})
	t.Run("Should: return error because value not exist", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_String,
				Path: "asfasf",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerResponseCode_Error, s.GetLogData().Code)
		assert.Equal(t, "", s.GetLogData().Meta.Value.GetStringValue())
	})
	t.Run("Should: not return error because selectors is missing", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerResponseCode_OK, s.GetLogData().Code)
	})
	t.Run("Should: parse single bool value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_Bool,
				Path: "success",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerResponseCode_OK, s.GetLogData().Code)
		assert.Equal(t, true, s.GetLogData().Meta.Value.GetBoolValue())
	})
	t.Run("Should: parse single string value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_String,
				Path: "name",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerResponseCode_OK, s.GetLogData().Code)
		assert.Equal(t, "John", s.GetLogData().Meta.Value.GetStringValue())
	})
	t.Run("Should: parse single number value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_Number,
				Path: "age",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerResponseCode_OK, s.GetLogData().Code)
		assert.Equal(t, float64(31), s.GetLogData().Meta.Value.GetNumberValue())
	})
	t.Run("Should: parse single any value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_Any,
				Path: "age",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerResponseCode_OK, s.GetLogData().Code)
		assert.Equal(t, "31", s.GetLogData().Meta.Value.GetStringValue())
	})
	t.Run("Should: parse single raw value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_Raw,
				Path: "raw",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerResponseCode_OK, s.GetLogData().Code)
		assert.Equal(t, `{"name":"ahha"}`, s.GetLogData().Meta.Value.GetStringValue())
	})
	t.Run("Should: parse single time value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_Time,
				Path: "time",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerResponseCode_OK, s.GetLogData().Code)
		assert.Equal(t, "2012-04-23T18:25:43Z", s.GetLogData().Meta.Value.GetStringValue())
	})
	t.Run("Should: parse multipile value", func(t *testing.T) {
		s := ExecHTTPValue("", 0, &scheduler_config_storage.HTTPValueConfig{Method: http.MethodGet, Headers: map[string]string{}, Selectors: []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_Time,
				Path: "time",
			},
			{
				Type: apiPb.HttpJsonValueConfig_Number,
				Path: "age",
			},
		}}, &mockSuccess{})
		assert.Equal(t, apiPb.SchedulerResponseCode_OK, s.GetLogData().Code)
		assert.EqualValues(t, &structType.ListValue{
			Values: []*structType.Value{
				{
					Kind: &structType.Value_StringValue{
						StringValue: "2012-04-23T18:25:43Z",
					},
				},
				{
					Kind: &structType.Value_NumberValue{
						NumberValue: float64(31),
					},
				},
			},
		}, s.GetLogData().Meta.Value.GetListValue())
	})
}
