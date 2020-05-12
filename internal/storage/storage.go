package storage

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	logger "log"
	"os"
	"squzy/internal/job"
	"time"
)

type Storage interface {
	Write(log job.CheckError) error
}

type memory struct {
	infoLogger *logger.Logger
	errLogger  *logger.Logger
}

func (m *memory) Write(log job.CheckError) error {
	logData := log.GetLogData()
	logID := uuid.New().String()
	startTime, err := ptypes.Timestamp(logData.Meta.StartTime)

	if err != nil {
		return err
	}

	endTime, err := ptypes.Timestamp(logData.Meta.EndTime)

	if err != nil {
		return err
	}

	if logData.Code == apiPb.SchedulerResponseCode_OK {
		m.infoLogger.Println(fmt.Sprintf(
			"SchedulerId: %s, Value: %s, LogId: %s, Status: Ok, Type: %s, startTime: %s, endTime: %s, duration: %s",
			logData.SchedulerId,
			logData.Meta.Value,
			logID,
			logData.Type.String(),
			startTime.Format(time.RFC3339),
			endTime.Format(time.RFC3339),
			fmt.Sprintf("%f", endTime.Sub(startTime).Seconds()),
		))
		return nil
	}
	m.errLogger.Println(fmt.Sprintf(
		"SchedulerId: %s, LogId: %s, Error msg: %s, Type: %s, startTime: %s, endTime: %s, duration: %s",
		logData.SchedulerId,
		logID,
		logData.Error.Message,
		logData.Type.String(),
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		fmt.Sprintf("%f", endTime.Sub(startTime).Seconds()),
	))
	return nil
}

func GetInMemoryStorage() Storage {
	return &memory{
		infoLogger: logger.New(os.Stdout, "INFO: ", logger.Ldate|logger.Ltime),
		errLogger:  logger.New(os.Stderr, "ERROR: ", logger.Ldate|logger.Ltime),
	}
}
