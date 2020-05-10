package server

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/metadata"
	"testing"
)

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
		AgentId: primitive.NewObjectID().Hex(),
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
		AgentId: "asf",
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

func (d dbMockOk) GetById(ctx context.Context, id primitive.ObjectID) (*apiPb.AgentItem, error) {
	return &apiPb.AgentItem{}, nil
}

func (d dbMockOk) Add(ctx context.Context, agent *apiPb.RegisterRequest) (string, error) {
	return "", nil
}

func (d dbMockOk) UpdateStatus(ctx context.Context, agentId primitive.ObjectID, status apiPb.AgentStatus) error {
	return nil
}

func (d dbMockOk) GetAll(ctx context.Context, filter bson.M) ([]*apiPb.AgentItem, error) {
	return []*apiPb.AgentItem{}, nil
}

type dbMockError struct {
}

func (d dbMockError) GetById(ctx context.Context, id primitive.ObjectID) (*apiPb.AgentItem, error) {
	return nil, errors.New("")
}

func (d dbMockError) Add(ctx context.Context, agent *apiPb.RegisterRequest) (string, error) {
	return "", errors.New("")
}

func (d dbMockError) UpdateStatus(ctx context.Context, agentId primitive.ObjectID, status apiPb.AgentStatus) error {
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
		s := New(&dbMockError{}, nil)
		_, err := s.Register(context.Background(), &apiPb.RegisterRequest{
			AgentName: "",
			HostInfo:  nil,
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&dbMockOk{}, nil)
		_, err := s.Register(context.Background(), &apiPb.RegisterRequest{
			AgentName: "",
			HostInfo:  nil,
		})
		assert.Equal(t, nil, err)
	})
}

func TestServer_UnRegister(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&dbMockError{}, nil)
		_, err := s.UnRegister(context.Background(), &apiPb.UnRegisterRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because id", func(t *testing.T) {
		s := New(&dbMockOk{}, nil)
		_, err := s.UnRegister(context.Background(), &apiPb.UnRegisterRequest{})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: not return  error", func(t *testing.T) {
		s := New(&dbMockOk{}, nil)
		_, err := s.UnRegister(context.Background(), &apiPb.UnRegisterRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.Equal(t, nil, err)
	})
}

func TestServer_GetAgentList(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&dbMockError{}, nil)
		_, err := s.GetAgentList(context.Background(), &empty.Empty{})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&dbMockOk{}, nil)
		_, err := s.GetAgentList(context.Background(), &empty.Empty{})
		assert.Equal(t, nil, err)
	})
}

func TestServer_GetByAgentName(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&dbMockError{}, nil)
		_, err := s.GetByAgentName(context.Background(), &apiPb.GetByAgentNameRequest{})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&dbMockOk{}, nil)
		_, err := s.GetByAgentName(context.Background(), &apiPb.GetByAgentNameRequest{})
		assert.Equal(t, nil, err)
	})
}

func TestServer_SendMetrics(t *testing.T) {
	t.Run("Should: return error if id is wrong", func(t *testing.T) {
		s := New(&dbMockOk{}, nil)
		assert.NotEqual(t, nil, s.SendMetrics(&mockStreamError{}))
	})
	t.Run("Should: close stream without error if error", func(t *testing.T) {
		s := New(&dbMockOk{}, nil)
		assert.NotEqual(t, nil, s.SendMetrics(&mockInternalStreamError{}))
	})
	t.Run("Should: works as excpected", func(t *testing.T) {
		s := New(&dbMockOk{}, nil)
		assert.Equal(t, nil, s.SendMetrics(&mockStreamOk{}))
	})
}

func TestServer_GetAgentById(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&dbMockError{}, nil)
		_, err := s.GetAgentById(context.Background(),  &apiPb.GetAgentByIdRequest{
			AgentId: primitive.NewObjectID().Hex(),
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because bson", func(t *testing.T) {
		s := New(&dbMockOk{}, nil)
		_, err := s.GetAgentById(context.Background(), &apiPb.GetAgentByIdRequest{})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&dbMockOk{}, nil)
		_, err := s.GetAgentById(context.Background(), &apiPb.GetAgentByIdRequest{
			AgentId: primitive.NewObjectID().Hex(),
		})
		assert.Equal(t, nil, err)
	})
}
