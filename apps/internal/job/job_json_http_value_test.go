package job

import (
	"errors"
	structType "github.com/golang/protobuf/ptypes/struct"
	httpPb "github.com/squzy/squzy_generated/generated/server/proto/v1"
	storagePb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type mockError struct {

}

type mockSuccess struct {

}

func (m mockSuccess) SendRequest(req *http.Request) (int, []byte, error) {
	return 0, []byte(`{"name":"John", "age":31, "city":"New York", "success": true, "time": "2012-04-23T18:25:43.511Z"}`), nil
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

func TestNewJsonHttpValueJob(t *testing.T) {
	t.Run("Should: implement interface", func(t *testing.T) {
		s := NewJsonHttpValueJob(http.MethodGet, "", map[string]string{}, nil, nil)
		assert.Implements(t, (*Job)(nil), s)
	})
}

func TestJsonHttpValueJob_Do(t *testing.T) {
	t.Run("Should: return error on http request", func(t *testing.T) {
		s := NewJsonHttpValueJob(http.MethodGet, "", map[string]string{}, &mockError{}, nil)
		assert.Equal(t, storagePb.StatusCode_Error, s.Do().GetLogData().Code)
	})
	t.Run("Should: not return error because selectors is missing", func(t *testing.T) {
		s := NewJsonHttpValueJob(http.MethodGet, "", map[string]string{}, &mockSuccess{}, nil)
		assert.Equal(t, storagePb.StatusCode_OK, s.Do().GetLogData().Code)
	})
	t.Run("Should: parse single bool value", func(t *testing.T) {
		s := NewJsonHttpValueJob(http.MethodGet, "", map[string]string{}, &mockSuccess{}, []*httpPb.HttpJsonValueCheck_Selectors{
			{
				Type: httpPb.HttpJsonValueCheck_Bool,
				Path: "success",
			},
		})
		res := s.Do()
		assert.Equal(t, storagePb.StatusCode_OK, res.GetLogData().Code)
		assert.Equal(t, true, res.GetLogData().Value.GetBoolValue())
	})
	t.Run("Should: parse single string value", func(t *testing.T) {
		s := NewJsonHttpValueJob(http.MethodGet, "", map[string]string{}, &mockSuccess{}, []*httpPb.HttpJsonValueCheck_Selectors{
			{
				Type: httpPb.HttpJsonValueCheck_String,
				Path: "name",
			},
		})
		res := s.Do()
		assert.Equal(t, storagePb.StatusCode_OK, res.GetLogData().Code)
		assert.Equal(t, "John", res.GetLogData().Value.GetStringValue())
	})
	t.Run("Should: parse single number value", func(t *testing.T) {
		s := NewJsonHttpValueJob(http.MethodGet, "", map[string]string{}, &mockSuccess{}, []*httpPb.HttpJsonValueCheck_Selectors{
			{
				Type: httpPb.HttpJsonValueCheck_Number,
				Path: "age",
			},
		})
		res := s.Do()
		assert.Equal(t, storagePb.StatusCode_OK, res.GetLogData().Code)
		assert.Equal(t, float64(31), res.GetLogData().Value.GetNumberValue())
	})
	t.Run("Should: parse single any value", func(t *testing.T) {
		s := NewJsonHttpValueJob(http.MethodGet, "", map[string]string{}, &mockSuccess{}, []*httpPb.HttpJsonValueCheck_Selectors{
			{
				Type: httpPb.HttpJsonValueCheck_Any,
				Path: "age",
			},
		})
		res := s.Do()
		assert.Equal(t, storagePb.StatusCode_OK, res.GetLogData().Code)
		assert.Equal(t, "31", res.GetLogData().Value.GetStringValue())
	})
	t.Run("Should: parse single time value", func(t *testing.T) {
		s := NewJsonHttpValueJob(http.MethodGet, "", map[string]string{}, &mockSuccess{}, []*httpPb.HttpJsonValueCheck_Selectors{
			{
				Type: httpPb.HttpJsonValueCheck_Time,
				Path: "time",
			},
		})
		res := s.Do()
		assert.Equal(t, storagePb.StatusCode_OK, res.GetLogData().Code)
		assert.Equal(t, "2012-04-23T18:25:43Z", res.GetLogData().Value.GetStringValue())
	})
	t.Run("Should: parse multipile value", func(t *testing.T) {
		s := NewJsonHttpValueJob(http.MethodGet, "", map[string]string{}, &mockSuccess{}, []*httpPb.HttpJsonValueCheck_Selectors{
			{
				Type: httpPb.HttpJsonValueCheck_Time,
				Path: "time",
			},
			{
				Type: httpPb.HttpJsonValueCheck_Number,
				Path: "age",
			},
		})
		res := s.Do()
		assert.Equal(t, storagePb.StatusCode_OK, res.GetLogData().Code)
		assert.EqualValues(t, &structType.ListValue{
			Values:               []*structType.Value{
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

		}, res.GetLogData().Value.GetListValue())
	})
}