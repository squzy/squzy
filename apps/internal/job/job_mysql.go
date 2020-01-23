package job

import (
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
)

type mysqlJob struct {
	host      string
	port      int32
	user      string
	password  string
	dbname    string
	mySqlOpen func(string, string) (*sql.DB, error)
	mySqlPing func(*sql.DB) error
}

func NewMysqlJob(
	host string,
	port int32,
	user, password, dbname string,
	mySqlPing func(*sql.DB) error) Job {
	return &mysqlJob{
		host:      host,
		port:      port,
		user:      user,
		password:  password,
		dbname:    dbname,
		mySqlOpen: sql.Open,
		mySqlPing: mySqlPing,
	}
}

type mysqlError struct {
	time        *timestamp.Timestamp
	code        clientPb.StatusCode
	description string
	location    string
	port        int32
}

func newSqlError(time *timestamp.Timestamp, code clientPb.StatusCode, description, location string, port int32) CheckError {
	return &mysqlError{
		time:        time,
		code:        code,
		description: description,
		location:    location,
		port:        port,
	}
}

func (m *mysqlError) GetLogData() *clientPb.Log {
	return &clientPb.Log{
		Code:        m.code,
		Description: m.description,
		Meta: &clientPb.MetaData{
			Id:       uuid.New().String(),
			Location: m.location,
			Port:     m.port,
		},
	}
}

func (j *mysqlJob) Do() CheckError {
	sqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		j.user, j.password, j.host, j.port, j.dbname)
	db, err := j.mySqlOpen("mysql", sqlInfo)
	if err != nil {
		return newSqlError(ptypes.TimestampNow(), clientPb.StatusCode_Error, mysqlConnectionError.Error(), j.host, j.port)
	}
	defer db.Close()

	err = j.mySqlPing(db)
	if err != nil {
		return newSqlError(ptypes.TimestampNow(), clientPb.StatusCode_Error, mysqlPingError.Error(), j.host, j.port)
	}

	return newSqlError(ptypes.TimestampNow(), clientPb.StatusCode_OK, "", j.host, j.port)
}
