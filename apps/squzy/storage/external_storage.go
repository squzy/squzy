package storage

import (
	"context"
	"errors"
	clientPb "github.com/squzy/squzy_generated/generated/logger"
	"squzy/apps/internal/job"
	"time"
)

type externalStorage struct {
	client clientPb.LoggerClient
}

const (
	loggerConnTimeout = time.Second * 5
)

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
	ctx, cancel := context.WithTimeout(context.Background(), loggerConnTimeout)
	defer cancel()
	res, err := s.client.SendLogMessage(ctx, req)
	if err != nil {
		return connectionExternalStorageError
	}
	if !res.Success {
		return storageNotSaveLog
	}
	return nil
}
