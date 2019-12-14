package job

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"time"

	"github.com/gocql/gocql"
)

type cassandraJob struct {
	cluster  string
	user     string
	password string
}

func NewCassandraJob(cluster, user, password string) Job {
	return &cassandraJob{
		cluster:  cluster,
		user:     user,
		password: password,
	}
}

type cassandraError struct {
	time        *timestamp.Timestamp
	code        clientPb.StatusCode
	description string
	cluster     string
}

func newCassandraError(time *timestamp.Timestamp, code clientPb.StatusCode, description, cluster string) CheckError {
	return &cassandraError{
		time:        time,
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
			Time:     m.time,
		},
	}
}

func (j *cassandraJob) Do() CheckError {
	cluster := gocql.NewCluster(j.cluster, j.cluster, j.cluster)
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: j.user, Password: j.password}
	session, err := cluster.CreateSession()
	if err != nil {
		return newCassandraError(ptypes.TimestampNow(), clientPb.StatusCode_Error, postgresConnectionError.Error(), j.cluster)
	}
	defer session.Close()

	err = session.ExecuteBatch(session.NewBatch(gocql.UnloggedBatch)) //TODO: check correctness
	if err != nil {
		return newCassandraError(ptypes.TimestampNow(), clientPb.StatusCode_Error, postgresPingError.Error(), j.cluster)
	}

	return newCassandraError(ptypes.TimestampNow(), clientPb.StatusCode_OK, "", j.cluster)
}
