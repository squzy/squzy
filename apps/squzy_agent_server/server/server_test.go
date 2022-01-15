package server

import (
	"context"
	"errors"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	empty "google.golang.org/protobuf/types/known/emptypb"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"testing"
)

type storageMock struct {
}

func (s storageMock) SaveIncident(ctx context.Context, in *apiPb.Incident, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (s storageMock) UpdateIncidentStatus(ctx context.Context, in *apiPb.UpdateIncidentStatusRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	panic("implement me")
}

func (s storageMock) GetIncidentById(ctx context.Context, in *apiPb.IncidentIdRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	panic("implement me")
}

func (s storageMock) GetIncidentByRuleId(ctx context.Context, in *apiPb.RuleIdRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	panic("implement me")
}

func (s storageMock) GetIncidentsList(ctx context.Context, in *apiPb.GetIncidentsListRequest, opts ...grpc.CallOption) (*apiPb.GetIncidentsListResponse, error) {
	panic("implement me")
}

func (s storageMock) GetSchedulerUptime(ctx context.Context, in *apiPb.GetSchedulerUptimeRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerUptimeResponse, error) {
	panic("implement me")
}

func (s storageMock) SaveTransaction(ctx context.Context, in *apiPb.TransactionInfo, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (s storageMock) GetTransactionsGroup(ctx context.Context, in *apiPb.GetTransactionGroupRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionGroupResponse, error) {
	panic("implement me")
}

func (s storageMock) GetTransactions(ctx context.Context, in *apiPb.GetTransactionsRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionsResponse, error) {
	panic("implement me")
}

func (s storageMock) GetTransactionById(ctx context.Context, in *apiPb.GetTransactionByIdRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionByIdResponse, error) {
	panic("implement me")
}

func (s storageMock) SaveResponseFromScheduler(ctx context.Context, in *apiPb.SchedulerResponse, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, nil
}

func (s storageMock) SaveResponseFromAgent(ctx context.Context, in *apiPb.Metric, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, nil
}

func (s storageMock) GetSchedulerInformation(ctx context.Context, in *apiPb.GetSchedulerInformationRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerInformationResponse, error) {
	panic("implement me")
}

func (s storageMock) GetAgentInformation(ctx context.Context, in *apiPb.GetAgentInformationRequest, opts ...grpc.CallOption) (*apiPb.GetAgentInformationResponse, error) {
	panic("implement me")
}

type mockStreamClose struct {
}

func (m mockStreamClose) SendAndClose(e *empty.Empty) error {
	return nil
}

type mockStreamContinueWork struct {
	count int
}

func (m mockStreamContinueWork) SendAndClose(e *empty.Empty) error {
	return nil
}

func (m *mockStreamContinueWork) Recv() (*apiPb.SendMetricsRequest, error) {
	defer func() {
		m.count += 1
	}()
	if m.count == 3 {
		return nil, io.EOF
	}
	if m.count == 2 {
		return &apiPb.SendMetricsRequest{
			Msg: &apiPb.SendMetricsRequest_Disconnect_{
				Disconnect: &apiPb.SendMetricsRequest_Disconnect{
					AgentId: primitive.NewObjectID().Hex(),
					Time:    timestamp.Now(),
				},
			},
		}, nil
	}
	return &apiPb.SendMetricsRequest{
		Msg: &apiPb.SendMetricsRequest_Metric{
			Metric: &apiPb.Metric{
				AgentId: primitive.NewObjectID().Hex(),
			},
		},
	}, nil
}

func (m mockStreamContinueWork) SetHeader(md metadata.MD) error {
	panic("implement me")
}

func (m mockStreamContinueWork) SendHeader(md metadata.MD) error {
	panic("implement me")
}

func (m mockStreamContinueWork) SetTrailer(md metadata.MD) {
	panic("implement me")
}

func (m mockStreamContinueWork) Context() context.Context {
	panic("implement me")
}

func (mockStreamContinueWork) SendMsg(m interface{}) error {
	panic("implement me")
}

func (mockStreamContinueWork) RecvMsg(m interface{}) error {
	panic("implement me")
}

func (m mockStreamClose) Recv() (*apiPb.SendMetricsRequest, error) {
	return &apiPb.SendMetricsRequest{
		Msg: &apiPb.SendMetricsRequest_Disconnect_{
			Disconnect: &apiPb.SendMetricsRequest_Disconnect{
				AgentId: primitive.NewObjectID().Hex(),
				Time:    timestamp.Now(),
			},
		},
	}, nil
}

func (m mockStreamClose) SetHeader(md metadata.MD) error {
	panic("implement me")
}

func (m mockStreamClose) SendHeader(md metadata.MD) error {
	panic("implement me")
}

func (m mockStreamClose) SetTrailer(md metadata.MD) {
	panic("implement me")
}

func (m mockStreamClose) Context() context.Context {
	panic("implement me")
}

func (mockStreamClose) SendMsg(m interface{}) error {
	panic("implement me")
}

func (mockStreamClose) RecvMsg(m interface{}) error {
	panic("implement me")
}

type mockStreamOk struct {
	exec bool
}

func (m mockStreamOk) SendAndClose(*empty.Empty) error {
	return nil
}

func (m *mockStreamOk) Recv() (*apiPb.SendMetricsRequest, error) {
	if m.exec {
		return nil, errors.New("")
	}
	m.exec = true
	return &apiPb.SendMetricsRequest{
		Msg: &apiPb.SendMetricsRequest_Metric{
			Metric: &apiPb.Metric{
				AgentId: primitive.NewObjectID().Hex(),
			},
		},
	}, nil
}

func (m mockStreamOk) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (m mockStreamOk) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (m mockStreamOk) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (m mockStreamOk) Context() context.Context {
	panic("implement me")
}

func (m mockStreamOk) SendMsg(_ interface{}) error {
	panic("implement me")
}

func (m mockStreamOk) RecvMsg(_ interface{}) error {
	panic("implement me")
}

type mockInternalStreamError struct {
}

func (m mockInternalStreamError) SendAndClose(*empty.Empty) error {
	return nil
}

func (m mockInternalStreamError) Recv() (*apiPb.SendMetricsRequest, error) {
	return nil, errors.New("")
}

func (m mockInternalStreamError) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (m mockInternalStreamError) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (m mockInternalStreamError) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (m mockInternalStreamError) Context() context.Context {
	panic("implement me")
}

func (m mockInternalStreamError) SendMsg(_ interface{}) error {
	panic("implement me")
}

func (m mockInternalStreamError) RecvMsg(_ interface{}) error {
	panic("implement me")
}

type mockStreamError struct {
}

func (m mockStreamError) SendAndClose(*empty.Empty) error {
	panic("implement me")
}

func (m mockStreamError) Recv() (*apiPb.SendMetricsRequest, error) {
	return &apiPb.SendMetricsRequest{
		Msg: &apiPb.SendMetricsRequest_Metric{
			Metric: &apiPb.Metric{
				AgentId: "asf",
			},
		},
	}, nil
}

func (m mockStreamError) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (m mockStreamError) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (m mockStreamError) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (m mockStreamError) Context() context.Context {
	panic("implement me")
}

func (m mockStreamError) SendMsg(_ interface{}) error {
	panic("implement me")
}

func (m mockStreamError) RecvMsg(_ interface{}) error {
	panic("implement me")
}

type dbMockOk struct {
}

func (d dbMockOk) GetByID(ctx context.Context, id primitive.ObjectID) (*apiPb.AgentItem, error) {
	return &apiPb.AgentItem{}, nil
}

func (d dbMockOk) Add(ctx context.Context, agent *apiPb.RegisterRequest) (string, error) {
	return "", nil
}

func (d dbMockOk) UpdateStatus(ctx context.Context, agentId primitive.ObjectID, status apiPb.AgentStatus, time *timestamp.Timestamp) error {
	return nil
}

func (d dbMockOk) GetAll(ctx context.Context, filter bson.M) ([]*apiPb.AgentItem, error) {
	return []*apiPb.AgentItem{}, nil
}

type dbMockError struct {
}

func (d dbMockError) GetByID(ctx context.Context, id primitive.ObjectID) (*apiPb.AgentItem, error) {
	return nil, errors.New("")
}

func (d dbMockError) Add(ctx context.Context, agent *apiPb.RegisterRequest) (string, error) {
	return "", errors.New("")
}

func (d dbMockError) UpdateStatus(ctx context.Context, agentId primitive.ObjectID, status apiPb.AgentStatus, time *timestamp.Timestamp) error {
	return errors.New("")
}

func (d dbMockError) GetAll(ctx context.Context, filter bson.M) ([]*apiPb.AgentItem, error) {
	return nil, errors.New("")
}

func TestNew(t *testing.T) {
	t.Run("Should: not nil", func(t *testing.T) {
		s := New(nil, nil)
		assert.NotEqual(t, nil, s)
	})
}

func TestServer_Register(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&dbMockError{}, &storageMock{})
		_, err := s.Register(context.Background(), &apiPb.RegisterRequest{
			AgentName: "",
			HostInfo:  nil,
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&dbMockOk{}, &storageMock{})
		_, err := s.Register(context.Background(), &apiPb.RegisterRequest{
			AgentName: "",
			HostInfo:  nil,
		})
		assert.Equal(t, nil, err)
	})
}

func TestServer_UnRegister(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&dbMockError{}, &storageMock{})
		_, err := s.UnRegister(context.Background(), &apiPb.UnRegisterRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because id", func(t *testing.T) {
		s := New(&dbMockOk{}, &storageMock{})
		_, err := s.UnRegister(context.Background(), &apiPb.UnRegisterRequest{})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: not return  error", func(t *testing.T) {
		s := New(&dbMockOk{}, &storageMock{})
		_, err := s.UnRegister(context.Background(), &apiPb.UnRegisterRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.Equal(t, nil, err)
	})
}

func TestServer_GetAgentList(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&dbMockError{}, &storageMock{})
		_, err := s.GetAgentList(context.Background(), &empty.Empty{})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&dbMockOk{}, &storageMock{})
		_, err := s.GetAgentList(context.Background(), &empty.Empty{})
		assert.Equal(t, nil, err)
	})
}

func TestServer_GetByAgentName(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&dbMockError{}, &storageMock{})
		_, err := s.GetByAgentName(context.Background(), &apiPb.GetByAgentNameRequest{})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&dbMockOk{}, &storageMock{})
		_, err := s.GetByAgentName(context.Background(), &apiPb.GetByAgentNameRequest{})
		assert.Equal(t, nil, err)
	})
}

func TestServer_SendMetrics(t *testing.T) {
	t.Run("Should: return error if id is wrong", func(t *testing.T) {
		s := New(&dbMockOk{}, &storageMock{})
		assert.NotEqual(t, nil, s.SendMetrics(&mockStreamError{}))
	})
	t.Run("Should: close stream without error if error", func(t *testing.T) {
		s := New(&dbMockOk{}, &storageMock{})
		assert.NotEqual(t, nil, s.SendMetrics(&mockInternalStreamError{}))
	})
	t.Run("Should: works as excpected", func(t *testing.T) {
		s := New(&dbMockOk{}, &storageMock{})
		assert.Equal(t, nil, s.SendMetrics(&mockStreamOk{}))
	})
	t.Run("Should: works as excpected", func(t *testing.T) {
		s := New(&dbMockOk{}, &storageMock{})
		assert.Equal(t, nil, s.SendMetrics(&mockStreamClose{}))
	})
	t.Run("Should: works as excpected", func(t *testing.T) {
		s := New(&dbMockOk{}, &storageMock{})
		assert.Equal(t, nil, s.SendMetrics(&mockStreamContinueWork{}))
	})
}

func TestServer_GetAgentById(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&dbMockError{}, &storageMock{})
		_, err := s.GetAgentById(context.Background(), &apiPb.GetAgentByIdRequest{
			AgentId: primitive.NewObjectID().Hex(),
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because bson", func(t *testing.T) {
		s := New(&dbMockOk{}, &storageMock{})
		_, err := s.GetAgentById(context.Background(), &apiPb.GetAgentByIdRequest{})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&dbMockOk{}, &storageMock{})
		_, err := s.GetAgentById(context.Background(), &apiPb.GetAgentByIdRequest{
			AgentId: primitive.NewObjectID().Hex(),
		})
		assert.Equal(t, nil, err)
	})
}
