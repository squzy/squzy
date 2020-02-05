package job

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"squzy/apps/internal/cassandraTools"
)

type cassandraJob struct {
	cluster        string
	user           string
	password       string
	cassandraTools cassandraTools.CassandraTools
}

func NewCassandraJob(cluster, user, password string) Job {
	return &cassandraJob{
		cluster:        cluster,
		user:           user,
		password:       password,
		cassandraTools: cassandraTools.NewCassandraTools(cluster, user, password),
	}
}

type cassandraError struct {
	logId       string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        clientPb.StatusCode
	description string
	cluster     string
}

func newCassandraError(logId string, startTime, endTime *timestamp.Timestamp, code clientPb.StatusCode, description, cluster string) CheckError {
	return &cassandraError{
		logId:       logId,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		cluster:     cluster,
	}
}

func (m *cassandraError) GetLogData() *clientPb.Log {
	return &clientPb.Log{
		Code:        m.code,
		Description: m.description,
		Meta: &clientPb.MetaData{
			Id:       uuid.New().String(),
			Location: m.cluster,
		},
	}
}

func (j *cassandraJob) Do() CheckError {
	logId := uuid.New().String()
	startTime := ptypes.TimestampNow()

	session, err := j.cassandraTools.CreateSession()
	if err != nil {
		return newCassandraError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_Error, postgresConnectionError.Error(), j.cluster)
	}
	defer j.cassandraTools.Close(session)

	err = j.cassandraTools.ExecuteBatch(session, j.cassandraTools.NewBatch(session)) //TODO: check correctness
	if err != nil {
		return newCassandraError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_Error, postgresPingError.Error(), j.cluster)
	}

	return newCassandraError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_OK, "", j.cluster)
}
