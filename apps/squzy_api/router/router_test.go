package router

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-playground/assert/v2"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockOk struct {

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

func (m mockOk) AddScheduler(ctx context.Context, scheduler *apiPb.AddRequest) error {
	return nil
}

type mockError struct {

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

func (m mockError) AddScheduler(ctx context.Context, scheduler *apiPb.AddRequest) error {
	return errors.New("")
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
			Path string
			Method string
			ExpectedCode int
			Body io.Reader
		}
		tt := []*TestCase{
			{
				Path: "/v1/agents",
				Method: http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path: "/v1/agents/agent",
				Method: http.MethodGet,
				ExpectedCode: http.StatusNotFound,
			},
			{
				Path: "/v1/schedulers",
				Method: http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path: "/v1/schedulers/scheduler",
				Method: http.MethodGet,
				ExpectedCode: http.StatusNotFound,
			},
			{
				Path: "/v1/schedulers/scheduler/run",
				Method: http.MethodPut,
				ExpectedCode: http.StatusNotFound,
			},
			{
				Path: "/v1/schedulers/scheduler",
				Method: http.MethodDelete,
				ExpectedCode: http.StatusNotFound,
			},
			{
				Path: "/v1/schedulers/scheduler/stop",
				Method: http.MethodPut,
				ExpectedCode: http.StatusNotFound,
			},
			{
				Path: "/v1/schedulers",
				Method: http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
			},
			{
				Path: "/v1/schedulers",
				Method: http.MethodPost,
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
				Path: "/v1/schedulers",
				Method: http.MethodPost,
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
				Path: "/v1/schedulers",
				Method: http.MethodPost,
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
				Path: "/v1/schedulers",
				Method: http.MethodPost,
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
				Path: "/v1/schedulers",
				Method: http.MethodPost,
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
				Path: "/v1/schedulers",
				Method: http.MethodPost,
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
				Path: "/v1/schedulers",
				Method: http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 0,
							"tcpConfig": {
								"host": "GET",
								"port": 32
							}
						}
					`,
				)),
			},
			{
				Path: "/v1/schedulers",
				Method: http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 1,
							"grpcConfig": {
								"host": "GET",
								"port": 3
							}
						}
					`,
				)),
			},
			{
				Path: "/v1/schedulers",
				Method: http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 2,
							"httpConfig": {
								"method": "GET",
								"url": "https://google.ru"
							}
						}
					`,
				)),
			},
			{
				Path: "/v1/schedulers",
				Method: http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 3
						}
					`,
				)),
			},
			{
				Path: "/v1/schedulers",
				Method: http.MethodPost,
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
			Path string
			Method string
			ExpectedCode int
			Body io.Reader
		}
		tt := []*TestCase{
			{
				Path: "/v1/agents",
				Method: http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path: "/v1/agents/agent",
				Method: http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path: "/v1/schedulers",
				Method: http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path: "/v1/schedulers/scheduler",
				Method: http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path: "/v1/schedulers/scheduler/run",
				Method: http.MethodPut,
				ExpectedCode: http.StatusAccepted,
			},
			{
				Path: "/v1/schedulers/scheduler",
				Method: http.MethodDelete,
				ExpectedCode: http.StatusAccepted,
			},
			{
				Path: "/v1/schedulers/scheduler/stop",
				Method: http.MethodPut,
				ExpectedCode: http.StatusAccepted,
			},
			{
				Path: "/v1/schedulers",
				Method: http.MethodPost,
				ExpectedCode: http.StatusCreated,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 0,
							"tcpConfig": {
								"host": "GET",
								"port": 32
							}
						}
					`,
				)),
			},
			{
				Path: "/v1/schedulers",
				Method: http.MethodPost,
				ExpectedCode: http.StatusCreated,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 1,
							"grpcConfig": {
								"host": "GET",
								"port": 3
							}
						}
					`,
				)),
			},
			{
				Path: "/v1/schedulers",
				Method: http.MethodPost,
				ExpectedCode: http.StatusCreated,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"interval": 10,
							"timeout": 10,
							"type": 2,
							"httpConfig": {
								"method": "GET",
								"url": "https://google.ru"
							}
						}
					`,
				)),
			},
			{
				Path: "/v1/schedulers",
				Method: http.MethodPost,
				ExpectedCode: http.StatusCreated,
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
				Path: "/v1/schedulers",
				Method: http.MethodPost,
				ExpectedCode: http.StatusCreated,
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
		}

		for _, test := range tt {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.Method, test.Path, test.Body)
			r.ServeHTTP(w, req)
			assert.Equal(t, test.ExpectedCode, w.Code)
		}
	})
}