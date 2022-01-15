package storage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/squzy/squzy/internal/job"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"time"
)

type Storage interface {
	Write(log job.CheckError) error
}

type memory struct {
}

func (m *memory) Write(log job.CheckError) error {
	logData := log.GetLogData()
	logID := uuid.New().String()
	startTime := logData.Snapshot.Meta.StartTime.AsTime()
	err := logData.Snapshot.Meta.StartTime.CheckValid()
	if err != nil {
		return err
	}

	endTime := logData.Snapshot.Meta.EndTime.AsTime()
	err = logData.Snapshot.Meta.EndTime.CheckValid()
	if err != nil {
		return err
	}

	if logData.Snapshot.Code == apiPb.SchedulerCode_OK {
		logger.Info(fmt.Sprintf(
			"SchedulerId: %s, Value: %s, LogId: %s, Status: Ok, Type: %s, startTime: %s, endTime: %s, duration: %s",
			logData.SchedulerId,
			logData.Snapshot.Meta.Value,
			logID,
			logData.Snapshot.Type.String(),
			startTime.Format(time.RFC3339),
			endTime.Format(time.RFC3339),
			fmt.Sprintf("%f", endTime.Sub(startTime).Seconds()),
		))
		return nil
	}
	logger.Error(fmt.Sprintf(
		"SchedulerId: %s, LogId: %s, Error msg: %s, Type: %s, startTime: %s, endTime: %s, duration: %s",
		logData.SchedulerId,
		logID,
		logData.Snapshot.Error.Message,
		logData.Snapshot.Type.String(),
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		fmt.Sprintf("%f", endTime.Sub(startTime).Seconds()),
	))
	return nil
}

func GetInMemoryStorage() Storage {
	return &memory{}
}
