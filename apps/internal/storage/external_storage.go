package storage

import (
	"context"
	"errors"
	storagePb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"google.golang.org/grpc"
	"squzy/apps/internal/grpcTools"
	"squzy/apps/internal/job"
	"time"
)

type externalStorage struct {
	client storagePb.LoggerClient
}

const (
	loggerConnTimeout = time.Second * 5
)

var (
	connectionExternalStorageError = errors.New("CANT_CONNECT_TO_EXTERNAL_STORAGE")
	storageNotSaveLog              = errors.New("EXTERNAL_STORAGE_NOT_SAVE_LOG")
)

func NewExternalStorage(grpcTools grpcTools.GrpcTool, address string, timeout time.Duration, fallBack Storage) Storage {
	conn, err := grpcTools.GetConnection(address, timeout, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fallBack
	}
	return &externalStorage{client: storagePb.NewLoggerClient(conn)}
}

func (s *externalStorage) Write(id string, log job.CheckError) error {

	req := &storagePb.SendLogMessageRequest{
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
