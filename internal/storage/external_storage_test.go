package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"net"
	"squzy/internal/grpcTools"
	"squzy/internal/job"
	"testing"
	"time"
)

type server struct {
}

func (s server) SendResponseFromAgent(context.Context, *apiPb.SendMetricsRequest) (*empty.Empty, error) {
	panic("implement me")
}

type serverErrorThrow struct {
}

func (s serverErrorThrow) SendResponseFromScheduler(context.Context, *apiPb.SchedulerResponse) (*empty.Empty, error) {
	return nil, errors.New("saf")
}

func (s serverErrorThrow) SendResponseFromAgent(context.Context, *apiPb.SendMetricsRequest) (*empty.Empty, error) {
	panic("implement me")
}

func (s server) SendResponseFromScheduler(context.Context, *apiPb.SchedulerResponse) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

type mockStorage struct {
}

type mockStorageError struct {
}

func (m mockStorageError) Write(log job.CheckError) error {
	return storageNotSaveLog
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

	t.Run("Should: return storageNotSaveLog", func(t *testing.T) {
		s := NewExternalStorage(&grpcMockError{}, "", time.Second, &mockStorageError{}, grpc.WithInsecure(), grpc.WithBlock())
		assert.Equal(t, storageNotSaveLog, s.Write(&mock{}))
	})
	t.Run("Should: not return error on write real storage", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 12122))
		grpcServer := grpc.NewServer()
		apiPb.RegisterStorageServer(grpcServer, &server{})
		go func() {
			_ = grpcServer.Serve(lis)
		}()
		time.Sleep(time.Second * 2)
		s := NewExternalStorage(grpcTools.New(), "localhost:12122", time.Second*2, &mockStorage{}, grpc.WithInsecure(), grpc.WithBlock())
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
		s := NewExternalStorage(grpcTools.New(), "localhost:12124", time.Second*2, &mockStorage{}, grpc.WithInsecure(), grpc.WithBlock())
		assert.Equal(t, connectionExternalStorageError, s.Write(&mock{}))
	})
}
