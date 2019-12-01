package storage

import (
	"context"
	"errors"
	squzy_logger_v1_service "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
)

type client_mock_error struct {

}

type client_mock_error_not_save struct {

}

func (c client_mock_error_not_save) SendLogMessage(ctx context.Context, in *squzy_logger_v1_service.SendLogMessageRequest, opts ...grpc.CallOption) (*squzy_logger_v1_service.SendLogMessageResponse, error) {
	return &squzy_logger_v1_service.SendLogMessageResponse{
		Success: false,
	}, nil
}

type client_mock struct {

}

func (c client_mock) SendLogMessage(ctx context.Context, in *squzy_logger_v1_service.SendLogMessageRequest, opts ...grpc.CallOption) (*squzy_logger_v1_service.SendLogMessageResponse, error) {
	return &squzy_logger_v1_service.SendLogMessageResponse{
		Success: true,
	}, nil
}

type mock struct {

}

func (m mock) GetLogData() *squzy_logger_v1_service.Log {
	return &squzy_logger_v1_service.Log{}
}

func (c client_mock_error) SendLogMessage(ctx context.Context, in *squzy_logger_v1_service.SendLogMessageRequest, opts ...grpc.CallOption) (*squzy_logger_v1_service.SendLogMessageResponse, error) {
	return nil, errors.New("Error")
}

func TestNewExternalStorage(t *testing.T) {
	t.Run("Test: Create new storage", func(t *testing.T) {
		s := NewExternalStorage(nil)
		assert.Implements(t, (*Storage)(nil), s)
	})
}

func TestExternalStorage_Write(t *testing.T) {
	t.Run("Should: return connectionExternalStorageError", func(t *testing.T) {
		s := NewExternalStorage(&client_mock_error{})
		assert.Equal(t, connectionExternalStorageError, s.Write("", &mock{}))
	})

	t.Run("Should: return nil", func(t *testing.T) {
		s := NewExternalStorage(&client_mock{})
		assert.Equal(t, nil, s.Write("", &mock{}))
	})

	t.Run("Should: return storageNotSaveLog", func(t *testing.T) {
		s := NewExternalStorage(&client_mock_error_not_save{})
		assert.Equal(t, storageNotSaveLog, s.Write("", &mock{}))
	})
}

func TestMemory_Write(t *testing.T) {
	t.Run("Memory storage", func(t *testing.T) {
		s := GetInMemoryStorage()
		assert.Implements(t, (*Storage)(nil), s)
	})
}