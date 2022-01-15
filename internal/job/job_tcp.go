package job

import (
	"fmt"
	"github.com/squzy/squzy/internal/helpers"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"net"
)

type tcpError struct {
	schedulerID string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerCode
	description string
}

func (s *tcpError) GetLogData() *apiPb.SchedulerResponse {
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
			Type:  apiPb.SchedulerType_TCP,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: s.startTime,
				EndTime:   s.endTime,
			},
		},
	}
}

func newTCPError(schedulerID string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerCode, description string) CheckError {
	return &tcpError{
		schedulerID: schedulerID,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
	}
}

func ExecTCP(schedulerID string, timeout int32, config *scheduler_config_storage.TCPConfig) CheckError {
	startTime := timestamp.Now()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(config.Host, fmt.Sprintf("%d", config.Port)), helpers.DurationFromSecond(timeout))
	if err != nil {
		return newTCPError(schedulerID, startTime, timestamp.Now(), apiPb.SchedulerCode_ERROR, errWrongConnectConfigError.Error())
	}
	if conn != nil {
		defer func() {
			_ = conn.Close()
		}()
	}
	return newTCPError(schedulerID, startTime, timestamp.Now(), apiPb.SchedulerCode_OK, "")
}
