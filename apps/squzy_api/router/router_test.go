package router

import (
	"bytes"
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
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

func (m mockOk) ArchivedApplicationById(ctx context.Context, id string) (*apiPb.Application, error) {
	return &apiPb.Application{}, nil
}

func (m mockOk) EnabledApplicationById(ctx context.Context, id string) (*apiPb.Application, error) {
	return &apiPb.Application{}, nil
}

func (m mockOk) DisabledApplicationById(ctx context.Context, id string) (*apiPb.Application, error) {
	return &apiPb.Application{}, nil
}

func (m mockOk) GetSchedulerUptime(ctx context.Context, rq *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	return &apiPb.GetSchedulerUptimeResponse{}, nil
}

func (m mockOk) GetTransactionGroups(ctx context.Context, req *apiPb.GetTransactionGroupRequest) (*apiPb.GetTransactionGroupResponse, error) {
	return &apiPb.GetTransactionGroupResponse{}, nil
}

func (m mockOk) GetTransactionsList(ctx context.Context, req *apiPb.GetTransactionsRequest) (*apiPb.GetTransactionsResponse, error) {
	return &apiPb.GetTransactionsResponse{}, nil
}

func (m mockOk) GetApplicationById(ctx context.Context, id string) (*apiPb.Application, error) {
	return &apiPb.Application{}, nil
}

func (m mockOk) GetApplicationList(ctx context.Context) ([]*apiPb.Application, error) {
	return []*apiPb.Application{}, nil
}

func (m mockOk) GetTransactionById(ctx context.Context, id string) (*apiPb.GetTransactionByIdResponse, error) {
	return &apiPb.GetTransactionByIdResponse{}, nil
}

func (m mockOk) RegisterApplication(ctx context.Context, rq *apiPb.ApplicationInfo) (*apiPb.InitializeApplicationResponse, error) {
	return &apiPb.InitializeApplicationResponse{}, nil
}

func (m mockOk) SaveTransaction(ctx context.Context, rq *apiPb.TransactionInfo) (*empty.Empty, error) {
	return &empty.Empty{}, nil
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

func (m mockError) ArchivedApplicationById(ctx context.Context, id string) (*apiPb.Application, error) {
	return nil, errors.New("")
}

func (m mockError) EnabledApplicationById(ctx context.Context, id string) (*apiPb.Application, error) {
	return nil, errors.New("")
}

func (m mockError) DisabledApplicationById(ctx context.Context, id string) (*apiPb.Application, error) {
	return nil, errors.New("")
}

func (m mockError) GetSchedulerUptime(ctx context.Context, rq *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	return nil, errors.New("")
}

func (m mockError) GetTransactionGroups(ctx context.Context, req *apiPb.GetTransactionGroupRequest) (*apiPb.GetTransactionGroupResponse, error) {
	return nil, errors.New("")
}

func (m mockError) GetTransactionsList(ctx context.Context, req *apiPb.GetTransactionsRequest) (*apiPb.GetTransactionsResponse, error) {
	return nil, errors.New("")
}

func (m mockError) GetApplicationById(ctx context.Context, id string) (*apiPb.Application, error) {
	return nil, errors.New("")
}

func (m mockError) GetApplicationList(ctx context.Context) ([]*apiPb.Application, error) {
	return nil, errors.New("")
}

func (m mockError) GetTransactionById(ctx context.Context, id string) (*apiPb.GetTransactionByIdResponse, error) {
	return nil, errors.New("")
}

func (m mockError) RegisterApplication(ctx context.Context, rq *apiPb.ApplicationInfo) (*apiPb.InitializeApplicationResponse, error) {
	return nil, errors.New("")
}

func (m mockError) SaveTransaction(ctx context.Context, rq *apiPb.TransactionInfo) (*empty.Empty, error) {
	return nil, errors.New("")
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
			{
				Path:         "/v1/schedulers/scheduler/uptime",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path:         "/v1/schedulers/scheduler/uptime?dateFrom=0000-01-01T00:00:00.899Z&dateTo=0000-01-01T00:00:00.899Z",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusUnprocessableEntity,
			},
			{
				Path:         "/v1/schedulers/scheduler/uptime?dateFrom=12321323&dateTo=12321323",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusUnprocessableEntity,
			},
			{
				Path:         "/v1/applications",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path:         "/v1/applications",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusInternalServerError,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"name": "sf"
						}
					`,
				)),
			},
			{
				Path:         "/v1/applications",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusUnprocessableEntity,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"name": ""
						}
					`,
				)),
			},
			{
				Path:         "/v1/applications/app",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path:         "/v1/applications/app/transactions/list?dateFrom=12321323&dateTo=12321323",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusBadRequest,
			},
			{
				Path:         "/v1/applications/app/transactions/list?dateFrom=0000-01-01T00:00:00.899Z&dateTo=0000-01-01T00:00:00.899Z",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusUnprocessableEntity,
			},
			{
				Path:         "/v1/applications/app/transactions/list?dateFrom=2020-05-07T19:17:05.899Z&dateTo=2020-05-17T19:17:05.899Z",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path:         "/v1/applications/app/transactions/group?dateFrom=12321323&dateTo=12321323",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusBadRequest,
			},
			{
				Path:         "/v1/applications/app/transactions/group?dateFrom=0000-01-01T00:00:00.899Z&dateTo=0000-01-01T00:00:00.899Z",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusUnprocessableEntity,
			},
			{
				Path:         "/v1/applications/app/transactions/group?dateFrom=2020-05-07T19:17:05.899Z&dateTo=2020-05-17T19:17:05.899Z",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path:         "/v1/applications/app/transactions/single/trra",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path:         "/v1/applications/app/transactions",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusAccepted,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"name": ""
						}
					`,
				)),
			},
			{
				Path:         "/v1/applications/app/transactions",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusAccepted,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"name": "safasf",
							"id": "asfasfasf",
							"dateFrom": "0000-01-01T00:00:00.899Z",
							"dateTo": "0000-01-01T00:00:00.899Z"
						}
					`,
				)),
			},
			{
				Path:         "/v1/applications/app/transactions",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusAccepted,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"name": "safasf",
							"id": "asfasfasf",
							"dateFrom": "2020-05-07T19:17:05.899Z",
							"dateTo": "0000-01-01T00:00:00.899Z"
						}
					`,
				)),
			},
			{
				Path:         "/v1/applications/app/transactions",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusAccepted,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"name": "safasf",
							"id": "asfasfasf",
							"dateFrom": "2020-05-07T19:17:05.899Z",
							"dateTo": "2020-05-07T19:17:05.899Z"
						}
					`,
				)),
			},
			{
				Path:         "/v1/applications/app/enabled",
				Method:       http.MethodPut,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path:         "/v1/applications/app/disabled",
				Method:       http.MethodPut,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Path:         "/v1/applications/app/archived",
				Method:       http.MethodPut,
				ExpectedCode: http.StatusInternalServerError,
			},
		}

		for _, test := range tt {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.Method, test.Path, test.Body)
			r.ServeHTTP(w, req)
			assert.Equal(t, test.ExpectedCode, w.Code, test.Path)
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
			{
				Path:         "/v1/applications/app/transactions/list?dateFrom=2020-05-07T19:17:05.899Z&dateTo=2020-05-17T19:17:05.899Z",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path:         "/v1/applications/app/transactions/group?dateFrom=2020-05-07T19:17:05.899Z&dateTo=2020-05-17T19:17:05.899Z",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path:         "/v1/schedulers/scheduler/uptime",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path:         "/v1/applications/app/transactions",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusAccepted,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"name": "safasf",
							"id": "asfasfasf",
							"dateFrom": "2020-05-07T19:17:05.899Z",
							"dateTo": "2020-05-07T19:17:05.899Z"
						}
					`,
				)),
			},
			{
				Path:         "/v1/applications/app",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path:         "/v1/applications",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path:         "/v1/applications",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusOK,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"name": "sf"
						}
					`,
				)),
			},
			{
				Path:         "/v1/applications/app/transactions",
				Method:       http.MethodPost,
				ExpectedCode: http.StatusAccepted,
				Body: bytes.NewBuffer([]byte(
					`
						{
							"name": "safasf",
							"id": "asfasfasf",
							"dateFrom": "2020-05-07T19:17:05.899Z",
							"dateTo": "2020-05-07T19:17:05.899Z",
							"error": {
								"message": "asffsaf"
							},
							"meta": {
								"host": "asfsf"
							}
						}
					`,
				)),
			},
			{
				Path:         "/v1/applications/app/transactions/single/trra",
				Method:       http.MethodGet,
				ExpectedCode: http.StatusOK,
			},
			{
				Path:         "/v1/applications/app/enabled",
				Method:       http.MethodPut,
				ExpectedCode: http.StatusAccepted,
			},
			{
				Path:         "/v1/applications/app/disabled",
				Method:       http.MethodPut,
				ExpectedCode: http.StatusAccepted,
			},
			{
				Path:         "/v1/applications/app/archived",
				Method:       http.MethodPut,
				ExpectedCode: http.StatusAccepted,
			},
		}

		for _, test := range tt {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.Method, test.Path, test.Body)
			r.ServeHTTP(w, req)
			assert.Equal(t, test.ExpectedCode, w.Code, test.Path)
		}
	})
}

func TestGetStringValueFromString(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		assert.Nil(t, GetStringValueFromString(""))
	})
	t.Run("Should: not return nil", func(t *testing.T) {
		assert.NotNil(t, GetStringValueFromString("1"))
	})
}

func TestGetSchedulerListSorting(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		assert.Nil(t, GetSchedulerListSorting(0, 0))
	})
	t.Run("Should: not return nil", func(t *testing.T) {
		assert.NotNil(t, GetSchedulerListSorting(0, 1))
	})
}

func TestGetTransactionListSorting(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		assert.Nil(t, GetTransactionListSorting(0, 0))
	})
	t.Run("Should: not return nil", func(t *testing.T) {
		assert.NotNil(t, GetTransactionListSorting(0, 1))
	})
}

func TestGetFilters(t *testing.T) {
	t.Run("Should: return all nils", func(t *testing.T) {
		p, f, err := GetFilters(nil, nil)
		assert.Nil(t, p)
		assert.Nil(t, f)
		assert.Nil(t, err)
	})
	t.Run("Should: return error because dateFrom", func(t *testing.T) {
		r := time.Unix(-62135596801, 0)
		_, _, err := GetFilters(&PaginationRequest{
			Page:  0,
			Limit: 0,
		}, &TimeFilterRequest{
			DateFrom: &r,
			DateTo:   nil,
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because dateTo", func(t *testing.T) {
		r := time.Unix(-62135596801, 0)
		_, _, err := GetFilters(&PaginationRequest{
			Page:  0,
			Limit: 0,
		}, &TimeFilterRequest{
			DateFrom: nil,
			DateTo:   &r,
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: parse time correct", func(t *testing.T) {
		tim := time.Now()
		pag, tF, err := GetFilters(&PaginationRequest{
			Page:  2,
			Limit: 24,
		}, &TimeFilterRequest{
			DateFrom: &tim,
			DateTo:   &tim,
		})
		assert.Nil(t, err)
		res, _ := ptypes.TimestampProto(tim)
		assert.Equal(t, int32(24), pag.Limit)
		assert.Equal(t, int32(2), pag.Page)
		assert.Equal(t, res, tF.To)
		assert.Equal(t, res, tF.From)
	})
	t.Run("Should: parse time correct", func(t *testing.T) {
		r, tf, err := GetFilters(&PaginationRequest{
			Page:  2,
			Limit: 24,
		}, nil)
		assert.Nil(t, err)
		assert.Equal(t, int32(24), r.Limit)
		assert.Equal(t, int32(2), r.Page)
		assert.Nil(t, tf)
	})
}
