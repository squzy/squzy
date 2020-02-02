package storage

import (
	"context"
	"errors"
	"fmt"
	storagePb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"google.golang.org/grpc"
	"log"
	"squzy/apps/internal/grpcTools"
	"squzy/apps/internal/job"
	"time"
)

type externalStorage struct {
	client   storagePb.LoggerClient
	fallback Storage
	address  string
}

const (
	loggerConnTimeout = time.Second * 5
)

var (
	connectionExternalStorageError = errors.New("CANT_CONNECT_TO_EXTERNAL_STORAGE")
	storageNotSaveLog              = errors.New("EXTERNAL_STORAGE_NOT_SAVE_LOG")
)

func NewExternalStorage(grpcTools grpcTools.GrpcTool, address string, timeout time.Duration, fallBack Storage, options ...grpc.DialOption) Storage {
	conn, err := grpcTools.GetConnection(address, timeout, options...)
	if err != nil {
		log.Println("Will wrote to in memory storage")
		return fallBack
	}
	log.Println(fmt.Sprintf("Will send log to client %s", address))
	return &externalStorage{
		client:   storagePb.NewLoggerClient(conn),
		fallback: fallBack,
		address:  address,
	}
}

func (s *externalStorage) Write(id string, checkerLog job.CheckError) error {
	req := &storagePb.SendLogMessageRequest{
		SchedulerId: id,
		Log:         checkerLog.GetLogData(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), loggerConnTimeout)
	defer cancel()
	res, err := s.client.SendLogMessage(ctx, req)
	if err != nil {
		if s.fallback != nil {
			log.Println(fmt.Sprintf("Cant connect to %s will use fallback to std", s.address))
			_ = s.fallback.Write(id, checkerLog)
		}
		return connectionExternalStorageError
	}
	if !res.Success {
		return storageNotSaveLog
	}
	return nil
}
