package job

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/squzy/squzy/internal/cassandra-tools"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type cassandraJob struct {
	cluster        string
	user           string
	password       string
	cassandraTools cassandra_tools.CassandraTools
}

type cassandraError struct {
	schedulerID string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerCode
	description string
	cluster     string
}

func newCassandraError(schedulerID string, startTime, endTime *timestamp.Timestamp, code apiPb.SchedulerCode, description, cluster string) CheckError {
	return &cassandraError{
		schedulerID: schedulerID,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		cluster:     cluster,
	}
}

func (s *cassandraError) GetLogData() *apiPb.SchedulerResponse {
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
			Type:  apiPb.SchedulerType_CASSANDRA,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: s.startTime,
				EndTime:   s.endTime,
			},
		},
	}
}

func ExecCassandra(schedulerID string, config *scheduler_config_storage.DbConfig, cTools cassandra_tools.CassandraTools) CheckError {
	startTime := timestamppb.Now()

	session, err := cTools.CreateSession()
	if err != nil {
		return newCassandraError(schedulerID, startTime, timestamppb.Now(), apiPb.SchedulerCode_ERROR, cassandraConnectionError.Error(), config.Cluster)
	}
	defer cTools.Close(session)

	err = cTools.ExecuteBatch(session, cTools.NewBatch(session)) //TODO: check correctness
	if err != nil {
		return newCassandraError(schedulerID, startTime, timestamppb.Now(), apiPb.SchedulerCode_ERROR, cassandraConnectionError.Error(), config.Cluster)
	}

	return newCassandraError(schedulerID, startTime, timestamppb.Now(), apiPb.SchedulerCode_OK, "", config.Cluster)
}
