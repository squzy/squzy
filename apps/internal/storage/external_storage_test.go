package storage

import (
	"errors"
	squzy_logger_v1_service "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"squzy/apps/internal/job"
	"testing"
	"time"
)

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
		s := NewExternalStorage(&grpcMock{}, "", time.Second, &mockStorage{})
		assert.Implements(t, (*Storage)(nil), s)
	})
}

func TestExternalStorage_Write(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		s := NewExternalStorage(&grpcMockError{}, "", time.Second, &mockStorage{})
		assert.Equal(t, nil, s.Write("", &mock{}))
	})

	t.Run("Should: return storageNotSaveLog", func(t *testing.T) {
		s := NewExternalStorage(&grpcMockError{}, "", time.Second,&mockStorageError{})
		assert.Equal(t, storageNotSaveLog, s.Write("", &mock{}))
	})
}

func TestMemory_Write(t *testing.T) {
	t.Run("Memory storage", func(t *testing.T) {
		s := GetInMemoryStorage()
		assert.Implements(t, (*Storage)(nil), s)
	})
}
