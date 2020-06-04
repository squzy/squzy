package server

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"squzy/apps/squzy_application_monitoring/database"
	"testing"
	"time"
)

type mockStorage struct {

}

func (m mockStorage) SaveResponseFromScheduler(ctx context.Context, in *apiPb.SchedulerResponse, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockStorage) SaveResponseFromAgent(ctx context.Context, in *apiPb.Metric, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockStorage) SaveTransaction(ctx context.Context, in *apiPb.TransactionInfo, opts ...grpc.CallOption) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (m mockStorage) GetSchedulerInformation(ctx context.Context, in *apiPb.GetSchedulerInformationRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerInformationResponse, error) {
	panic("implement me")
}

func (m mockStorage) GetSchedulerUptime(ctx context.Context, in *apiPb.GetSchedulerUptimeRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerUptimeResponse, error) {
	panic("implement me")
}

func (m mockStorage) GetAgentInformation(ctx context.Context, in *apiPb.GetAgentInformationRequest, opts ...grpc.CallOption) (*apiPb.GetAgentInformationResponse, error) {
	panic("implement me")
}

func (m mockStorage) GetTransactionsGroup(ctx context.Context, in *apiPb.GetTransactionGroupRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionGroupResponse, error) {
	panic("implement me")
}

func (m mockStorage) GetTransactions(ctx context.Context, in *apiPb.GetTransactionsRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionsResponse, error) {
	panic("implement me")
}

func (m mockStorage) GetTransactionById(ctx context.Context, in *apiPb.GetTransactionByIdRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionByIdResponse, error) {
	panic("implement me")
}

type mockCfg struct {

}

func (m mockCfg) GetTracingHeader() string {
	return ""
}

func (m mockCfg) GetPort() int32 {
	panic("implement me")
}

func (m mockCfg) GetMongoURI() string {
	panic("implement me")
}

func (m mockCfg) GetMongoDb() string {
	panic("implement me")
}

func (m mockCfg) GetMongoCollection() string {
	panic("implement me")
}

func (m mockCfg) GetStorageTimeout() time.Duration {
	return time.Second
}

func (m mockCfg) GetStorageHost() string {
	panic("implement me")
}

type dbMockError struct {

}

func (d dbMockError) FindOrCreate(ctx context.Context, name string, host string) (*database.Application, error) {
	return nil, errors.New("as")
}

func (d dbMockError) FindApplicationByName(ctx context.Context, name string) (*database.Application, error) {
	return nil, errors.New("as")
}

func (d dbMockError) FindApplicationById(ctx context.Context, id primitive.ObjectID) (*database.Application, error) {
	return nil, errors.New("as")
}

func (d dbMockError) FindAllApplication(ctx context.Context) ([]*database.Application, error) {
	return nil, errors.New("as")
}

type dbMockOk struct {

}

func (d dbMockOk) FindOrCreate(ctx context.Context, name string, host string) (*database.Application, error) {
	return &database.Application{

	}, nil
}

func (d dbMockOk) FindApplicationByName(ctx context.Context, name string) (*database.Application, error) {
	return &database.Application{

	}, nil
}

func (d dbMockOk) FindApplicationById(ctx context.Context, id primitive.ObjectID) (*database.Application, error) {
	return &database.Application{

	}, nil
}

func (d dbMockOk) FindAllApplication(ctx context.Context) ([]*database.Application, error) {
	return []*database.Application{
		{

		},
		{

		},
	}, nil
}

func TestNew(t *testing.T) {
	t.Run("Should: not be nil", func(t *testing.T) {
		s := New(nil, nil, nil)
		assert.NotNil(t, s)
	})
}

func TestServer_GetApplicationById(t *testing.T) {
	t.Run("Should: return application without error", func(t *testing.T) {
		s := New(&dbMockOk{}, nil, nil)
		_, err := s.GetApplicationById(context.Background(), &apiPb.GetApplicationByIdRequest{
			ApplicationId: primitive.NewObjectID().Hex(),
		})
		assert.Nil(t, err)
	})
	t.Run("Should: return error because objectId", func(t *testing.T) {
		s := New(&dbMockOk{}, nil, nil)
		_, err := s.GetApplicationById(context.Background(), &apiPb.GetApplicationByIdRequest{
			ApplicationId: "primitive.NewObjectID().Hex()",
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: return error database", func(t *testing.T) {
		s := New(&dbMockError{}, nil, nil)
		_, err := s.GetApplicationById(context.Background(), &apiPb.GetApplicationByIdRequest{
			ApplicationId:  primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
}

func TestServer_GetApplicationList(t *testing.T) {
	t.Run("Should: return application without error", func(t *testing.T) {
		s := New(&dbMockOk{}, nil, nil)
		_, err := s.GetApplicationList(context.Background(), &empty.Empty{})
		assert.Nil(t, err)
	})
	t.Run("Should: return error database", func(t *testing.T) {
		s := New(&dbMockError{}, nil, nil)
		_, err := s.GetApplicationList(context.Background(), &empty.Empty{})
		assert.NotNil(t, err)
	})
}

func TestServer_InitializeApplication(t *testing.T) {
	t.Run("Should: return id without error", func(t *testing.T) {
		s := New(&dbMockOk{}, &mockCfg{}, nil)
		_, err := s.InitializeApplication(context.Background(), &apiPb.ApplicationInfo{
			Name:                 "asfsf",
			HostName:             "",
		})
		assert.Nil(t, err)
	})
	t.Run("Should: return error because name", func(t *testing.T) {
		s := New(&dbMockOk{}, &mockCfg{}, nil)
		_, err := s.InitializeApplication(context.Background(), &apiPb.ApplicationInfo{
			Name:                 "",
			HostName:             "",
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: return error database", func(t *testing.T) {
		s := New(&dbMockError{}, nil, nil)
		_, err := s.InitializeApplication(context.Background(), &apiPb.ApplicationInfo{
			Name:                 "asfsf",
			HostName:             "",
		})
		assert.NotNil(t, err)
	})
}

func TestServer_SaveTransaction(t *testing.T) {
	t.Run("Should: return error because appId", func(t *testing.T) {
		s := New(&dbMockOk{}, nil, nil)
		_, err := s.SaveTransaction(context.Background(), &apiPb.TransactionInfo{
			Id:                   "",
			ApplicationId:        "",
			ParentId:             "",
			Meta:                 nil,
			Name:                 "",
			StartTime:            nil,
			EndTime:              nil,
			Status:               0,
			Type:                 0,
			Error:                nil,
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: return error because app not exist", func(t *testing.T) {
		s := New(&dbMockError{}, nil, nil)
		_, err := s.SaveTransaction(context.Background(), &apiPb.TransactionInfo{
			Id:                   primitive.NewObjectID().Hex(),
			ApplicationId:        primitive.NewObjectID().Hex(),
			ParentId:             "",
			Meta:                 nil,
			Name:                 "",
			StartTime:            nil,
			EndTime:              nil,
			Status:               0,
			Type:                 0,
			Error:                nil,
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: save without error", func(t *testing.T) {
		s := New(&dbMockOk{}, &mockCfg{}, &mockStorage{})
		_, err := s.SaveTransaction(context.Background(), &apiPb.TransactionInfo{
			Id:                   primitive.NewObjectID().Hex(),
			ApplicationId:        primitive.NewObjectID().Hex(),
			ParentId:             "",
			Meta:                 nil,
			Name:                 "",
			StartTime:            nil,
			EndTime:              nil,
			Status:               0,
			Type:                 0,
			Error:                nil,
		})
		assert.Nil(t, err)
	})
}