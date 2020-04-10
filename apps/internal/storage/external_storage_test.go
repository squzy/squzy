package storage

import (
	"context"
	"errors"
	"fmt"
	squzy_logger_v1_service "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"net"
	"squzy/apps/internal/grpcTools"
	"squzy/apps/internal/job"
	"testing"
	"time"
)

type server struct {
}

type serverError struct {
}

type serverErrorThrow struct {
}

func (s serverErrorThrow) SendLogMessage(context.Context, *squzy_logger_v1_service.SendLogMessageRequest) (*squzy_logger_v1_service.SendLogMessageResponse, error) {
	return nil, errors.New("asfasf")
}

func (s serverError) SendLogMessage(context.Context, *squzy_logger_v1_service.SendLogMessageRequest) (*squzy_logger_v1_service.SendLogMessageResponse, error) {
	return &squzy_logger_v1_service.SendLogMessageResponse{
		Success: false,
	}, nil
}

func (s server) SendLogMessage(context.Context, *squzy_logger_v1_service.SendLogMessageRequest) (*squzy_logger_v1_service.SendLogMessageResponse, error) {
	return &squzy_logger_v1_service.SendLogMessageResponse{
		Success: true,
	}, nil
}

type mockStorage struct {
}

type mockStorageError struct {
}

func (m mockStorageError) Write(id string, log job.CheckError) error {
	return storageNotSaveLog
}

func (m mockStorage) Write(id string, log job.CheckError) error {
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

func (m mock) GetLogData() *squzy_logger_v1_service.Log {
	return &squzy_logger_v1_service.Log{}
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
		assert.Equal(t, nil, s.Write("", &mock{}))
	})

	t.Run("Should: return storageNotSaveLog", func(t *testing.T) {
		s := NewExternalStorage(&grpcMockError{}, "", time.Second, &mockStorageError{}, grpc.WithInsecure(), grpc.WithBlock())
		assert.Equal(t, storageNotSaveLog, s.Write("", &mock{}))
	})
	t.Run("Should: not return error on write real storage", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 12122))
		grpcServer := grpc.NewServer()
		squzy_logger_v1_service.RegisterLoggerServer(grpcServer, &server{})
		go func() {
			_ = grpcServer.Serve(lis)
		}()
		s := NewExternalStorage(grpcTools.New(), "localhost:12122", time.Second*2, &mockStorage{}, grpc.WithInsecure(), grpc.WithBlock())
		assert.Equal(t, nil, s.Write("", &mock{}))
	})
	t.Run("Should: return error on write real storage", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 12123))
		grpcServer := grpc.NewServer()
		squzy_logger_v1_service.RegisterLoggerServer(grpcServer, &serverError{})
		go func() {
			_ = grpcServer.Serve(lis)
		}()
		s := NewExternalStorage(grpcTools.New(), "localhost:12123", time.Second*2, &mockStorage{}, grpc.WithInsecure(), grpc.WithBlock())
		assert.Equal(t, storageNotSaveLog, s.Write("", &mock{}))
	})
	t.Run("Should: return error connection error on write real storage", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 12124))
		grpcServer := grpc.NewServer()
		squzy_logger_v1_service.RegisterLoggerServer(grpcServer, &serverErrorThrow{})
		go func() {
			_ = grpcServer.Serve(lis)
		}()
		s := NewExternalStorage(grpcTools.New(), "localhost:12124", time.Second*2, &mockStorage{}, grpc.WithInsecure(), grpc.WithBlock())
		assert.Equal(t, connectionExternalStorageError, s.Write("", &mock{}))
	})
}
