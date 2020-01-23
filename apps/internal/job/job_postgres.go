package job

import (
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
)

type postgresJob struct {
	host         string
	port         int32
	user         string
	password     string
	dbname       string
	postgresOpen func(string, string) (*sql.DB, error)
	postgresPing func(*sql.DB) error
}

func NewPosgresDbJob(
	host string,
	port int32,
	user, password, dbname string,
	mySqlPing func(*sql.DB) error) Job {
	return &postgresJob{
		host:         host,
		port:         port,
		user:         user,
		password:     password,
		dbname:       dbname,
		postgresOpen: sql.Open,
		postgresPing: mySqlPing,
	}
}

type postgresError struct {
	time        *timestamp.Timestamp
	code        clientPb.StatusCode
	description string
	location    string
	port        int32
}

func newPostgresError(time *timestamp.Timestamp, code clientPb.StatusCode, description, location string, port int32) CheckError {
	return &postgresError{
		time:        time,
		code:        code,
		description: description,
		location:    location,
		port:        port,
	}
}

func (m *postgresError) GetLogData() *clientPb.Log {
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

func (j *postgresJob) Do() CheckError {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		j.host, j.port, j.user, j.password, j.dbname)
	db, err := j.postgresOpen("postgres", psqlInfo)
	if err != nil {
		return newPostgresError(ptypes.TimestampNow(), clientPb.StatusCode_Error, postgresConnectionError.Error(), j.host, j.port)
	}
	defer db.Close()

	err = j.postgresPing(db)
	if err != nil {
		return newPostgresError(ptypes.TimestampNow(), clientPb.StatusCode_Error, postgresPingError.Error(), j.host, j.port)
	}

	return newPostgresError(ptypes.TimestampNow(), clientPb.StatusCode_OK, "", j.host, j.port)
}
