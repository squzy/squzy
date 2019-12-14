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
	host     string
	port     int32
	user     string
	password string
	dbname   string
}

func NewPosgresDbJob(host string, port int32, user, password, dbname string) Job {
	return &postgresJob{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
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
			Time:     m.time,
			Port:     m.port,
		},
	}
}

func (j *postgresJob) Do() CheckError {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		j.host, j.port, j.user, j.password, j.dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return newPostgresError(ptypes.TimestampNow(), clientPb.StatusCode_Error, postgresConnectionError.Error(), j.host, j.port)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return newPostgresError(ptypes.TimestampNow(), clientPb.StatusCode_Error, postgresPingError.Error(), j.host, j.port)
	}

	return newPostgresError(ptypes.TimestampNow(), clientPb.StatusCode_OK, "", j.host, j.port)
}
