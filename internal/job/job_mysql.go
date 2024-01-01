package job

import (
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type mysqlJob struct {
	schedulerID string
	dbConfig    *scheduler_config_storage.DbConfig
	db          DBConnector
}

type mysqlError struct {
	schedulerID string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerCode
	description string
	location    string
	port        int32
}

func newSqlError(schedulerID string, startTime, endTime *timestamp.Timestamp, code apiPb.SchedulerCode, description, location string, port int32) CheckError {
	return &mysqlError{
		schedulerID: schedulerID,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		location:    location,
		port:        port,
	}
}

func (s *mysqlError) GetLogData() *apiPb.SchedulerResponse {
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
			Type:  apiPb.SchedulerType_MYSQL,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: s.startTime,
				EndTime:   s.endTime,
			},
		},
	}
}

func ExecMysql(schedulerId string, config *scheduler_config_storage.DbConfig, dbC DBConnector) CheckError {
	startTime := timestamppb.Now()

	sqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.User, config.Password, config.Host, config.Port, config.DbName)
	err := dbC.Connect("mysql", sqlInfo)
	if err != nil {
		return newSqlError(schedulerId, startTime, timestamppb.Now(), apiPb.SchedulerCode_ERROR, mysqlConnectionError.Error(), config.Host, config.Port)
	}
	defer dbC.Close()

	err = dbC.Ping()
	if err != nil {
		return newSqlError(schedulerId, startTime, timestamppb.Now(), apiPb.SchedulerCode_ERROR, mysqlPingError.Error(), config.Host, config.Port)
	}

	return newSqlError(schedulerId, startTime, timestamppb.Now(), apiPb.SchedulerCode_OK, "", config.Host, config.Port)
}

type DBConnector interface {
	Connect(string, string) error
	Ping() error
	Close() error
}

type DbClient interface {
	Ping() error
	Close() error
}

type DBConnection struct {
	Client DbClient
	Open   func(driverName, dataSourceName string) (*sql.DB, error)
}

func NewDBConnection() DBConnection {
	return DBConnection{
		Open: sql.Open,
	}
}

func (m DBConnection) Connect(driver string, dataSource string) error {
	if client, err := m.Open(driver, dataSource); err == nil {
		m.Client = client
		return err
	} else {
		return err
	}
}

func (m DBConnection) Ping() error {
	return m.Client.Ping()
}

func (m DBConnection) Close() error {
	return m.Client.Close()
}
