package router

import (
	"bytes"
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockOk struct {
}

func (m mockOk) GetSchedulerHistoryByID(ctx context.Context, rq *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error) {
	return &apiPb.GetSchedulerInformationResponse{}, nil
}

func (m mockOk) GetAgentHistoryByID(ctx context.Context, rq *apiPb.GetAgentInformationRequest) (*apiPb.GetAgentInformationResponse, error) {
	return &apiPb.GetAgentInformationResponse{}, nil
}

func (m mockOk) GetAgentList(ctx context.Context) ([]*apiPb.AgentItem, error) {
	return []*apiPb.AgentItem{}, nil
}

func (m mockOk) GetAgentByID(ctx context.Context, id string) (*apiPb.AgentItem, error) {
	return &apiPb.AgentItem{}, nil
}

func (m mockOk) GetSchedulerList(ctx context.Context) ([]*apiPb.Scheduler, error) {
	return []*apiPb.Scheduler{}, nil
}

func (m mockOk) GetSchedulerByID(ctx context.Context, id string) (*apiPb.Scheduler, error) {
	return &apiPb.Scheduler{}, nil
}

func (m mockOk) RunScheduler(ctx context.Context, id string) error {
	return nil
}

func (m mockOk) StopScheduler(ctx context.Context, id string) error {
	return nil
}

func (m mockOk) RemoveScheduler(ctx context.Context, id string) error {
	return nil
}

func (m mockOk) AddScheduler(ctx context.Context, scheduler *apiPb.AddRequest) (*apiPb.AddResponse, error) {
	return &apiPb.AddResponse{}, nil
}

type mockError struct {
}

func (m mockError) GetSchedulerHistoryByID(ctx context.Context, rq *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error) {
	return nil, errors.New("")
}

func (m mockError) GetAgentHistoryByID(ctx context.Context, rq *apiPb.GetAgentInformationRequest) (*apiPb.GetAgentInformationResponse, error) {
	return nil, errors.New("")
}

func (m mockError) GetAgentList(ctx context.Context) ([]*apiPb.AgentItem, error) {
	return nil, errors.New("")
}

func (m mockError) GetAgentByID(ctx context.Context, id string) (*apiPb.AgentItem, error) {
	return nil, errors.New("")
}

func (m mockError) GetSchedulerList(ctx context.Context) ([]*apiPb.Scheduler, error) {
	return nil, errors.New("")
}

func (m mockError) GetSchedulerByID(ctx context.Context, id string) (*apiPb.Scheduler, error) {
	return nil, errors.New("")
}

func (m mockError) RunScheduler(ctx context.Context, id string) error {
	return errors.New("")
}

func (m mockError) StopScheduler(ctx context.Context, id string) error {
	return errors.New("")
}

func (m mockError) RemoveScheduler(ctx context.Context, id string) error {
	return errors.New("")
}

func (m mockError) AddScheduler(ctx context.Context, scheduler *apiPb.AddRequest) (*apiPb.AddResponse, error) {
	return nil, errors.New("")
}

func TestNew(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		r := New(nil)
		assert.NotEqual(t, nil, r)
	})
}

func TestRouter_GetEngine(t *testing.T) {
	t.Run("Should: create router without error", func(t *testing.T) {
		r := New(nil)
		engine := r.GetEngine()
		assert.NotEqual(t, nil, engine)
	})
	t.Run("Should: return error with mockError", func(t *testing.T) {
		r := New(&mockError{}).GetEngine()
		type TestCase struct {
			Path         string
			Method       string
			ExpectedCode int
			Body         io.Reader
		}
		tt := []*TestCase{
			{
				Path:         "/v1/agents",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path:         "/v1/agents/agent",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusNotFound,
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path:         "/v1/schedulers/scheduler",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusNotFound,
			},
			{
				Path:         "/v1/schedulers/scheduler/run",
				Method:       http.MethodPut,
				ExpectedCode: http.StatusNotFound,
			},
			{
				Path:         "/v1/schedulers/scheduler",
				Method:       http.MethodDelete,
				ExpectedCode: http.StatusNotFound,
			},
			{
				Path:         "/v1/schedulers/scheduler/stop",
				Method:       http.MethodPut,
				ExpectedCode: http.StatusNotFound,
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 0
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 1
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 2
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 3,
							"siteMapConfig": {}
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 4,
							"httpValueConfig": {}
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 1000
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 1,
							"tcpConfig": {
								"host": "GET",
								"port": 32
							}
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 2,
							"grpcConfig": {
								"host": "GET",
								"port": 3
							}
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 3,
							"httpConfig": {
								"method": "GET",
								"url": "https://google.ru"
							}
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 4
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 5
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers/schdeduler/history?dateFrom=2020-05-7T19:17:05.899Z&dateTo=2020-05-17T19:17:05.899Z&page=2&limit=4",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path:         "/v1/agents/schdeduler/history?dateFrom=2020-05-7T19:17:05.899Z&dateTo=2020-05-17T19:17:05.899Z&page=2&limit=4",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path:         "/v1/schedulers/schdeduler/history?dateFrom=0000-01-01T00:00:00.899Z&dateTo=0000-01-01T00:00:00.899Z&page=2&limit=4",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusUnprocessableEntity,
			},
			{
				Path:         "/v1/agents/schdeduler/history?dateFrom=0000-01-01T00:00:00.899Z&dateTo=0000-01-01T00:00:00.899Z&page=2&limit=4",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusUnprocessableEntity,
			},
			{
				Path:         "/v1/schedulers/schdeduler/history?dateFrom=2020-05-07T19:17:05.899Z&dateTo=2020-05-17T19:17:05.899Z&page=2&limit=4",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path:         "/v1/agents/schdeduler/history?dateFrom=2020-05-07T19:17:05.899Z&dateTo=2020-05-17T19:17:05.899Z&page=2&limit=4",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
		}

		for _, test := range tt {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.Method, test.Path, test.Body)
			r.ServeHTTP(w, req)
			assert.Equal(t, test.ExpectedCode, w.Code)
		}
	})
	t.Run("Should: return success with mockOk", func(t *testing.T) {
		r := New(&mockOk{}).GetEngine()
		type TestCase struct {
			Path         string
			Method       string
			ExpectedCode int
			Body         io.Reader
		}
		tt := []*TestCase{
			{
				Path:         "/v1/agents",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path:         "/v1/agents/agent",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path:         "/v1/schedulers/scheduler",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path:         "/v1/schedulers/scheduler/run",
				Method:       http.MethodPut,
				ExpectedCode: http.StatusAccepted,
			},
			{
				Path:         "/v1/schedulers/scheduler",
				Method:       http.MethodDelete,
				ExpectedCode: http.StatusAccepted,
			},
			{
				Path:         "/v1/schedulers/scheduler/stop",
				Method:       http.MethodPut,
				ExpectedCode: http.StatusAccepted,
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusCreated,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 1,
							"tcpConfig": {
								"host": "GET",
								"port": 32
							}
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusCreated,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 2,
							"grpcConfig": {
								"host": "GET",
								"port": 3
							}
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusCreated,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 3,
							"httpConfig": {
								"method": "GET",
								"url": "https://google.ru"
							}
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusCreated,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 4,
							"siteMapConfig": {}
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusCreated,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 5,
							"httpValueConfig": {}
						}
					`,
				)),
			},
			{
				Path:         "/v1/schedulers/schdeduler/history?dateFrom=2020-05-17T19:17:05.899Z&dateTo=2020-05-17T19:17:05.899Z&page=2&limit=4",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path:         "/v1/agents/schdeduler/history?dateFrom=2020-05-17T19:17:05.899Z&dateTo=2020-05-17T19:17:05.899Z&page=2&limit=4",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
		}

		for _, test := range tt {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.Method, test.Path, test.Body)
			r.ServeHTTP(w, req)
			assert.Equal(t, test.ExpectedCode, w.Code)
		}
	})
}

func TestGetFilters(t *testing.T) {
	t.Run("Should: return error because dateFrom", func(t *testing.T) {
		r := time.Unix(-62135596801, 0)
		_, _, err := GetFilters(&HistoryFilterRequest{
			DateFrom: &r,
			DateTo:   nil,
			Page:     0,
			Limit:    0,
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because dateTo", func(t *testing.T) {
		r := time.Unix(-62135596801, 0)
		_, _, err := GetFilters(&HistoryFilterRequest{
			DateFrom: nil,
			DateTo:   &r,
			Page:     0,
			Limit:    0,
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: parse time correct", func(t *testing.T) {
		tim := time.Now()
		pag, tF, err := GetFilters(&HistoryFilterRequest{
			DateFrom: &tim,
			DateTo:   &tim,
			Page:     2,
			Limit:    24,
		})
		assert.Nil(t, err)
		res, _ := ptypes.TimestampProto(tim)
		assert.Equal(t, int32(24), pag.Limit)
		assert.Equal(t, int32(2), pag.Page)
		assert.Equal(t, res, tF.To)
		assert.Equal(t, res, tF.From)
	})
	t.Run("Should: parse time correct", func(t *testing.T) {
		r, tf, err := GetFilters(&HistoryFilterRequest{
			DateFrom: nil,
			DateTo:   nil,
			Page:     2,
			Limit:    24,
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(24), r.Limit)
		assert.Equal(t, int32(2), r.Page)
		assert.Nil(t, tf)
	})
}
