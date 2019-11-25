package storage

import (
	"context"
	"errors"
	clientPb "github.com/squzy/squzy_generated/generated/logger"
	"squzy/apps/internal/job"
)

type externalStorage struct {
	client clientPb.LoggerClient
}

var (
	connectionExternalStorageError = errors.New("CANT_CONNECT_TO_EXTERNAL_STORAGE")
	storageNotSaveLog              = errors.New("EXTERNAL_STORAGE_NOT_SAVE_LOG")
)

func NewExternalStorage(client clientPb.LoggerClient) Storage {
	return &externalStorage{client: client}
}

func (s *externalStorage) Write(id string, log job.CheckError) error {
	req := &clientPb.SendLogMessageRequest{
		Log: log.GetLogData(),
	}

	res, err := s.client.SendLogMessage(context.Background(), req)
	if err != nil {
		return connectionExternalStorageError
	}
	if !res.Success {
		return storageNotSaveLog
	}
	return nil
}
