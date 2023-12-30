package job

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type postgresJob struct {
	dbConfig *scheduler_config_storage.DbConfig
	db       DBConnector
}

type postgresError struct {
	schedulerID string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerCode
	description string
	location    string
	port        int32
}

func newPostgresError(schedulerID string, startTime, endTime *timestamp.Timestamp, code apiPb.SchedulerCode, description, location string, port int32) CheckError {
	return &postgresError{
		schedulerID: schedulerID,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		location:    location,
		port:        port,
	}
}

func (s *postgresError) GetLogData() *apiPb.SchedulerResponse {
	var err *apiPb.SchedulerSnapshot_Error
	if s.code == apiPb.SchedulerCode_ERROR {
		err = &apiPb.SchedulerSnapshot_Error{
			Message: s.description,
		}
	}
	return &apiPb.SchedulerResponse{
		SchedulerId: s.schedulerID,
		Snapshot: &apiPb.SchedulerSnapshot{
			Code:  s.code,
			Error: err,
			Type:  apiPb.SchedulerType_POSTGRES,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: s.startTime,
				EndTime:   s.endTime,
			},
		},
	}
}

func ExecPostgres(config *scheduler_config_storage.DbConfig, dbC DBConnector) CheckError {
	logId := uuid.New().String()
	startTime := timestamppb.Now()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DbName)
	err := dbC.Connect("postgres", psqlInfo)
	if err != nil {
		return newPostgresError(logId, startTime, timestamppb.Now(), apiPb.SchedulerCode_ERROR, postgresConnectionError.Error(), config.Host, config.Port)
	}
	defer dbC.Close()

	err = dbC.Ping()
	if err != nil {
		return newPostgresError(logId, startTime, timestamppb.Now(), apiPb.SchedulerCode_ERROR, postgresPingError.Error(), config.Host, config.Port)
	}

	return newPostgresError(logId, startTime, timestamppb.Now(), apiPb.SchedulerCode_OK, "", config.Host, config.Port)
}
