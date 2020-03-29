package storage

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	storagePb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	logger "log"
	"os"
	"squzy/apps/internal/job"
	"time"
)

type Storage interface {
	Write(id string, log job.CheckError) error
}

type memory struct {
	infoLogger *logger.Logger
	errLogger  *logger.Logger
}

func (m *memory) Write(id string, log job.CheckError) error {
	logData := log.GetLogData()

	startTime, err := ptypes.Timestamp(logData.Meta.StartTime)

	if err != nil {
		return err
	}

	endTime, err := ptypes.Timestamp(logData.Meta.EndTime)

	if err != nil {
		return err
	}

	if logData.Code == storagePb.StatusCode_OK {
		m.infoLogger.Println(fmt.Sprintf(
			"SchedulerId: %s, Value: %s, LogId: %s, Status: Ok, Type: %s, startTime: %s, endTime: %s, duration: %s",
			id,
			logData.Value.GetStringValue(),
			logData.Meta.Id,
			logData.Meta.Type.String(),
			startTime.Format(time.RFC3339),
			endTime.Format(time.RFC3339),
			fmt.Sprintf("%f", endTime.Sub(startTime).Seconds()),
		))
		return nil
	}
	m.errLogger.Println(fmt.Sprintf(
		"SchedulerId: %s, LogId: %s, Error msg: %s, Location: %s, port: %d, Type: %s, startTime: %s, endTime: %s, duration: %s",
		id,
		logData.Meta.Id,
		logData.Description,
		logData.Meta.Location,
		logData.Meta.Port,
		logData.Meta.Type.String(),
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
