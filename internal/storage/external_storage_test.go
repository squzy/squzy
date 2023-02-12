package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/squzy/squzy/internal/grpctools"
	"github.com/squzy/squzy/internal/job"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"net"
	"testing"
	"time"
)

type server struct {
	apiPb.UnimplementedStorageServer
}

func (s server) SaveIncident(context.Context, *apiPb.Incident) (*empty.Empty, error) {
	panic("implement me")
}

func (s server) UpdateIncidentStatus(context.Context, *apiPb.UpdateIncidentStatusRequest) (*apiPb.Incident, error) {
	panic("implement me")
}

func (s server) GetIncidentById(context.Context, *apiPb.IncidentIdRequest) (*apiPb.Incident, error) {
	panic("implement me")
}

func (s server) GetIncidentByRuleId(context.Context, *apiPb.RuleIdRequest) (*apiPb.Incident, error) {
	panic("implement me")
}

func (s server) GetIncidentsList(context.Context, *apiPb.GetIncidentsListRequest) (*apiPb.GetIncidentsListResponse, error) {
	panic("implement me")
}

func (s server) GetSchedulerUptime(ctx context.Context, request *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	panic("implement me")
}

func (s server) SaveTransaction(ctx context.Context, info *apiPb.TransactionInfo) (*empty.Empty, error) {
	panic("implement me")
}

func (s server) GetTransactionsGroup(ctx context.Context, request *apiPb.GetTransactionGroupRequest) (*apiPb.GetTransactionGroupResponse, error) {
	panic("implement me")
}

func (s server) GetTransactions(ctx context.Context, request *apiPb.GetTransactionsRequest) (*apiPb.GetTransactionsResponse, error) {
	panic("implement me")
}

func (s server) GetTransactionById(ctx context.Context, request *apiPb.GetTransactionByIdRequest) (*apiPb.GetTransactionByIdResponse, error) {
	panic("implement me")
}

func (s server) GetSchedulerInformation(ctx context.Context, request *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error) {
	panic("implement me")
}

func (s server) GetAgentInformation(ctx context.Context, request *apiPb.GetAgentInformationRequest) (*apiPb.GetAgentInformationResponse, error) {
	panic("implement me")
}

func (s server) SaveResponseFromAgent(context.Context, *apiPb.Metric) (*empty.Empty, error) {
	panic("implement me")
}

type serverErrorThrow struct {
	apiPb.UnimplementedStorageServer
}

func (s serverErrorThrow) SaveIncident(context.Context, *apiPb.Incident) (*empty.Empty, error) {
	panic("implement me")
}

func (s serverErrorThrow) UpdateIncidentStatus(context.Context, *apiPb.UpdateIncidentStatusRequest) (*apiPb.Incident, error) {
	panic("implement me")
}

func (s serverErrorThrow) GetIncidentById(context.Context, *apiPb.IncidentIdRequest) (*apiPb.Incident, error) {
	panic("implement me")
}

func (s serverErrorThrow) GetIncidentByRuleId(context.Context, *apiPb.RuleIdRequest) (*apiPb.Incident, error) {
	panic("implement me")
}

func (s serverErrorThrow) GetIncidentsList(context.Context, *apiPb.GetIncidentsListRequest) (*apiPb.GetIncidentsListResponse, error) {
	panic("implement me")
}

func (s serverErrorThrow) GetSchedulerUptime(ctx context.Context, request *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	panic("implement me")
}

func (s serverErrorThrow) SaveTransaction(ctx context.Context, info *apiPb.TransactionInfo) (*empty.Empty, error) {
	panic("implement me")
}

func (s serverErrorThrow) GetTransactionsGroup(ctx context.Context, request *apiPb.GetTransactionGroupRequest) (*apiPb.GetTransactionGroupResponse, error) {
	panic("implement me")
}

func (s serverErrorThrow) GetTransactions(ctx context.Context, request *apiPb.GetTransactionsRequest) (*apiPb.GetTransactionsResponse, error) {
	panic("implement me")
}

func (s serverErrorThrow) GetTransactionById(ctx context.Context, request *apiPb.GetTransactionByIdRequest) (*apiPb.GetTransactionByIdResponse, error) {
	panic("implement me")
}

func (s serverErrorThrow) GetSchedulerInformation(ctx context.Context, request *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error) {
	panic("implement me")
}

func (s serverErrorThrow) GetAgentInformation(ctx context.Context, request *apiPb.GetAgentInformationRequest) (*apiPb.GetAgentInformationResponse, error) {
	panic("implement me")
}

func (s serverErrorThrow) SaveResponseFromScheduler(context.Context, *apiPb.SchedulerResponse) (*empty.Empty, error) {
	return nil, errors.New("saf")
}

func (s serverErrorThrow) SaveResponseFromAgent(context.Context, *apiPb.Metric) (*empty.Empty, error) {
	panic("implement me")
}

func (s server) SaveResponseFromScheduler(context.Context, *apiPb.SchedulerResponse) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

type mockStorage struct {
}

type mockStorageError struct {
}

func (m mockStorageError) Write(log job.CheckError) error {
	return errStorageNotSaveLog
}

func (m mockStorage) Write(log job.CheckError) error {
	return nil
}

type grpcMock struct {
}

func (g grpcMock) GetConnection(address string, timeout time.Duration, option ...grpc.DialOption) (*grpc.ClientConn, error) {
	return &grpc.ClientConn{}, nil
}

type grpcMockError struct {
}

func (g grpcMockError) GetConnection(address string, timeout time.Duration, option ...grpc.DialOption) (*grpc.ClientConn, error) {
	return nil, errors.New("error")
}

type mock struct {
}

func (m mock) GetLogData() *apiPb.SchedulerResponse {
	return &apiPb.SchedulerResponse{}
}

func TestNewExternalStorage(t *testing.T) {
	t.Run("Test: Create new storage", func(t *testing.T) {
		s := NewExternalStorage(&grpcMock{}, "", time.Second, &mockStorage{}, grpc.WithInsecure(), grpc.WithBlock())
		assert.Implements(t, (*Storage)(nil), s)
	})
}

func TestExternalStorage_Write(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		s := NewExternalStorage(&grpcMockError{}, "", time.Second, &mockStorage{}, grpc.WithInsecure(), grpc.WithBlock())
		assert.Equal(t, nil, s.Write(&mock{}))
	})

	t.Run("Should: return errStorageNotSaveLog", func(t *testing.T) {
		s := NewExternalStorage(&grpcMockError{}, "", time.Second, &mockStorageError{}, grpc.WithInsecure(), grpc.WithBlock())
		assert.Equal(t, errStorageNotSaveLog, s.Write(&mock{}))
	})
	t.Run("Should: not return error on write real storage", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 12122))
		grpcServer := grpc.NewServer()
		apiPb.RegisterStorageServer(grpcServer, &server{})
		go func() {
			_ = grpcServer.Serve(lis)
		}()
		time.Sleep(time.Second * 2)
		s := NewExternalStorage(grpctools.New(), "localhost:12122", time.Second*2, &mockStorage{}, grpc.WithInsecure(), grpc.WithBlock())
		assert.Equal(t, nil, s.Write(&mock{}))
	})
	t.Run("Should: return error connection error on write real storage", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 12124))
		grpcServer := grpc.NewServer()
		apiPb.RegisterStorageServer(grpcServer, &serverErrorThrow{})
		go func() {
			_ = grpcServer.Serve(lis)
		}()
		time.Sleep(time.Second * 2)
		s := NewExternalStorage(grpctools.New(), "localhost:12124", time.Second*2, &mockStorage{}, grpc.WithInsecure(), grpc.WithBlock())
		assert.Equal(t, errConnectionExternalStorageError, s.Write(&mock{}))
	})
}
