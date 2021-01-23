package storage

import (
	"context"
	"errors"
	"fmt"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"squzy/internal/grpctools"
	"squzy/internal/job"
	"squzy/internal/logger"
	"time"
)

type externalStorage struct {
	client   apiPb.StorageClient
	fallback Storage
	address  string
}

const (
	loggerConnTimeout = time.Second * 5
)

var (
	errConnectionExternalStorageError = errors.New("CANT_CONNECT_TO_EXTERNAL_STORAGE")
	errStorageNotSaveLog              = errors.New("EXTERNAL_STORAGE_NOT_SAVE_LOG")
)

func NewExternalStorage(grpcTools grpctools.GrpcTool, address string, timeout time.Duration, fallBack Storage, options ...grpc.DialOption) Storage {
	conn, err := grpcTools.GetConnection(address, timeout, options...)
	if err != nil {
		logger.Info("Will wrote to in memory storage")
		return fallBack
	}
	logger.Info(fmt.Sprintf("Will send log to client %s", address))
	return &externalStorage{
		client:   apiPb.NewStorageClient(conn),
		fallback: fallBack,
		address:  address,
	}
}

func (s *externalStorage) Write(checkerLog job.CheckError) error {
	req := checkerLog.GetLogData()
	ctx, cancel := context.WithTimeout(context.Background(), loggerConnTimeout)
	defer cancel()
	_, err := s.client.SaveResponseFromScheduler(ctx, req)
	if err != nil {
		if s.fallback != nil {
			_ = s.fallback.Write(checkerLog)
		}
		return errConnectionExternalStorageError
	}
	return nil
}
