package handlers

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
)

type storageMockOk struct {
}

func (s storageMockOk) SaveTransaction(ctx context.Context, in *apiPb.TransactionInfo, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (s storageMockOk) GetTransactionsGroup(ctx context.Context, in *apiPb.GetTransactionGroupRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionGroupResponse, error) {
	panic("implement me")
}

func (s storageMockOk) GetTransactions(ctx context.Context, in *apiPb.GetTransactionsRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionsResponse, error) {
	panic("implement me")
}

func (s storageMockOk) GetTransactionById(ctx context.Context, in *apiPb.GetTransactionByIdRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionByIdResponse, error) {
	panic("implement me")
}

func (s storageMockOk) SaveResponseFromScheduler(ctx context.Context, in *apiPb.SchedulerResponse, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (s storageMockOk) SaveResponseFromAgent(ctx context.Context, in *apiPb.Metric, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (s storageMockOk) GetSchedulerInformation(ctx context.Context, in *apiPb.GetSchedulerInformationRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerInformationResponse, error) {
	return &apiPb.GetSchedulerInformationResponse{}, nil
}

func (s storageMockOk) GetAgentInformation(ctx context.Context, in *apiPb.GetAgentInformationRequest, opts ...grpc.CallOption) (*apiPb.GetAgentInformationResponse, error) {
	return &apiPb.GetAgentInformationResponse{}, nil
}

type storageMockError struct {
}

func (s storageMockError) SaveTransaction(ctx context.Context, in *apiPb.TransactionInfo, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (s storageMockError) GetTransactionsGroup(ctx context.Context, in *apiPb.GetTransactionGroupRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionGroupResponse, error) {
	panic("implement me")
}

func (s storageMockError) GetTransactions(ctx context.Context, in *apiPb.GetTransactionsRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionsResponse, error) {
	panic("implement me")
}

func (s storageMockError) GetTransactionById(ctx context.Context, in *apiPb.GetTransactionByIdRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionByIdResponse, error) {
	panic("implement me")
}

func (s storageMockError) SaveResponseFromScheduler(ctx context.Context, in *apiPb.SchedulerResponse, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (s storageMockError) SaveResponseFromAgent(ctx context.Context, in *apiPb.Metric, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (s storageMockError) GetSchedulerInformation(ctx context.Context, in *apiPb.GetSchedulerInformationRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerInformationResponse, error) {
	return nil, errors.New("")
}

func (s storageMockError) GetAgentInformation(ctx context.Context, in *apiPb.GetAgentInformationRequest, opts ...grpc.CallOption) (*apiPb.GetAgentInformationResponse, error) {
	return nil, errors.New("")
}

type agentMockOk struct {
}

func (a agentMockOk) Register(ctx context.Context, in *apiPb.RegisterRequest, opts ...grpc.CallOption) (*apiPb.RegisterResponse, error) {
	panic("implement me")
}

func (a agentMockOk) GetByAgentName(ctx context.Context, in *apiPb.GetByAgentNameRequest, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (a agentMockOk) GetAgentById(ctx context.Context, in *apiPb.GetAgentByIdRequest, opts ...grpc.CallOption) (*apiPb.AgentItem, error) {
	return &apiPb.AgentItem{}, nil
}

func (a agentMockOk) UnRegister(ctx context.Context, in *apiPb.UnRegisterRequest, opts ...grpc.CallOption) (*apiPb.UnRegisterResponse, error) {
	panic("implement me")
}

func (a agentMockOk) GetAgentList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
	return &apiPb.GetAgentListResponse{}, nil
}

func (a agentMockOk) SendMetrics(ctx context.Context, opts ...grpc.CallOption) (apiPb.AgentServer_SendMetricsClient, error) {
	panic("implement me")
}

type agentMockError struct {
}

func (a agentMockError) Register(ctx context.Context, in *apiPb.RegisterRequest, opts ...grpc.CallOption) (*apiPb.RegisterResponse, error) {
	panic("implement me")
}

func (a agentMockError) GetByAgentName(ctx context.Context, in *apiPb.GetByAgentNameRequest, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (a agentMockError) GetAgentById(ctx context.Context, in *apiPb.GetAgentByIdRequest, opts ...grpc.CallOption) (*apiPb.AgentItem, error) {
	return nil, errors.New("")
}

func (a agentMockError) UnRegister(ctx context.Context, in *apiPb.UnRegisterRequest, opts ...grpc.CallOption) (*apiPb.UnRegisterResponse, error) {
	panic("implement me")
}

func (a agentMockError) GetAgentList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
	return nil, errors.New("")
}

func (a agentMockError) SendMetrics(ctx context.Context, opts ...grpc.CallOption) (apiPb.AgentServer_SendMetricsClient, error) {
	panic("implement me")
}

type mockMonitoringError struct {
}

func (m mockMonitoringError) GetSchedulerList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*apiPb.GetSchedulerListResponse, error) {
	return nil, errors.New("")
}

func (m mockMonitoringError) GetSchedulerById(ctx context.Context, in *apiPb.GetSchedulerByIdRequest, opts ...grpc.CallOption) (*apiPb.Scheduler, error) {
	return nil, errors.New("")
}

func (m mockMonitoringError) Add(ctx context.Context, in *apiPb.AddRequest, opts ...grpc.CallOption) (*apiPb.AddResponse, error) {
	return nil, errors.New("")
}

func (m mockMonitoringError) Remove(ctx context.Context, in *apiPb.RemoveRequest, opts ...grpc.CallOption) (*apiPb.RemoveResponse, error) {
	return nil, errors.New("")
}

func (m mockMonitoringError) Run(ctx context.Context, in *apiPb.RunRequest, opts ...grpc.CallOption) (*apiPb.RunResponse, error) {
	return nil, errors.New("")
}

func (m mockMonitoringError) Stop(ctx context.Context, in *apiPb.StopRequest, opts ...grpc.CallOption) (*apiPb.StopResponse, error) {
	return nil, errors.New("")
}

type mockMonitoringOk struct {
}

func (m mockMonitoringOk) GetSchedulerList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*apiPb.GetSchedulerListResponse, error) {
	return &apiPb.GetSchedulerListResponse{}, nil
}

func (m mockMonitoringOk) GetSchedulerById(ctx context.Context, in *apiPb.GetSchedulerByIdRequest, opts ...grpc.CallOption) (*apiPb.Scheduler, error) {
	return &apiPb.Scheduler{}, nil
}

func (m mockMonitoringOk) Add(ctx context.Context, in *apiPb.AddRequest, opts ...grpc.CallOption) (*apiPb.AddResponse, error) {
	return &apiPb.AddResponse{}, nil
}

func (m mockMonitoringOk) Remove(ctx context.Context, in *apiPb.RemoveRequest, opts ...grpc.CallOption) (*apiPb.RemoveResponse, error) {
	return &apiPb.RemoveResponse{}, nil
}

func (m mockMonitoringOk) Run(ctx context.Context, in *apiPb.RunRequest, opts ...grpc.CallOption) (*apiPb.RunResponse, error) {
	return &apiPb.RunResponse{}, nil
}

func (m mockMonitoringOk) Stop(ctx context.Context, in *apiPb.StopRequest, opts ...grpc.CallOption) (*apiPb.StopResponse, error) {
	return &apiPb.StopResponse{}, nil
}

func TestNew(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, nil, nil, nil)
		assert.NotNil(t, s)
	})
}

func TestHandlers_AddScheduler(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, &mockMonitoringOk{}, nil, nil)
		_, err := s.AddScheduler(context.Background(), &apiPb.AddRequest{})
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(nil, &mockMonitoringError{}, nil, nil)
		_, err := s.AddScheduler(context.Background(), &apiPb.AddRequest{})
		assert.NotNil(t, err)
	})
}

func TestHandlers_GetAgentByID(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&agentMockOk{}, nil, nil, nil)
		_, err := s.GetAgentByID(context.Background(), "")
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&agentMockError{}, nil, nil, nil)
		_, err := s.GetAgentByID(context.Background(), "")
		assert.NotNil(t, err)
	})
}

func TestHandlers_GetAgentList(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&agentMockOk{}, nil, nil, nil)
		_, err := s.GetAgentList(context.Background())
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&agentMockError{}, nil, nil, nil)
		_, err := s.GetAgentList(context.Background())
		assert.NotNil(t, err)
	})
}

func TestHandlers_GetAgentHistoryByID(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, nil, &storageMockOk{}, nil)
		_, err := s.GetAgentHistoryByID(context.Background(), nil)
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(nil, nil, &storageMockError{}, nil)
		_, err := s.GetAgentHistoryByID(context.Background(), nil)
		assert.NotNil(t, err)
	})
}

func TestHandlers_GetSchedulerHistoryByID(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, nil, &storageMockOk{}, nil)
		_, err := s.GetSchedulerHistoryByID(context.Background(), nil)
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(nil, nil, &storageMockError{}, nil)
		_, err := s.GetSchedulerHistoryByID(context.Background(), nil)
		assert.NotNil(t, err)
	})
}

func TestHandlers_GetSchedulerByID(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, &mockMonitoringOk{}, nil, nil)
		_, err := s.GetSchedulerByID(context.Background(), "")
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(nil, &mockMonitoringError{}, nil, nil)
		_, err := s.GetSchedulerByID(context.Background(), "")
		assert.NotNil(t, err)
	})
}

func TestHandlers_GetSchedulerList(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, &mockMonitoringOk{}, nil, nil)
		_, err := s.GetSchedulerList(context.Background())
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(nil, &mockMonitoringError{}, nil, nil)
		_, err := s.GetSchedulerList(context.Background())
		assert.NotNil(t, err)
	})
}

func TestHandlers_RemoveScheduler(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, &mockMonitoringOk{}, nil, nil)
		err := s.RemoveScheduler(context.Background(), "nil")
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(nil, &mockMonitoringError{}, nil, nil)
		err := s.RemoveScheduler(context.Background(), "nil")
		assert.NotNil(t, err)
	})
}

func TestHandlers_RunScheduler(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, &mockMonitoringOk{}, nil, nil)
		err := s.RunScheduler(context.Background(), "nil")
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(nil, &mockMonitoringError{}, nil, nil)
		err := s.RunScheduler(context.Background(), "nil")
		assert.NotNil(t, err)
	})
}

func TestHandlers_StopScheduler(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, &mockMonitoringOk{}, nil, nil)
		err := s.StopScheduler(context.Background(), "nil")
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(nil, &mockMonitoringError{}, nil, nil)
		err := s.StopScheduler(context.Background(), "nil")
		assert.NotNil(t, err)
	})
}
