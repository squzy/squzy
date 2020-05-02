package job

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"net"
	"squzy/internal/helpers"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
)

type tcpError struct {
	schedulerId string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerResponseCode
	description string
}

func (s *tcpError) GetLogData() *apiPb.SchedulerResponse {
	var err *apiPb.SchedulerResponse_Error
	if s.code == apiPb.SchedulerResponseCode_Error {
		err = &apiPb.SchedulerResponse_Error{
			Message: s.description,
		}
	}
	return &apiPb.SchedulerResponse{
		SchedulerId: s.schedulerId,
		Code:        s.code,
		Error:       err,
		Type:        apiPb.SchedulerType_Tcp,
		Meta: &apiPb.SchedulerResponse_MetaData{
			StartTime: s.startTime,
			EndTime:   s.endTime,
		},
	}
}

func newTcpError(schedulerId string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerResponseCode, description string) CheckError {
	return &tcpError{
		schedulerId: schedulerId,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
	}
}

func ExecTcp(schedulerId string, timeout int32, config *scheduler_config_storage.TcpConfig) CheckError {
	startTime := ptypes.TimestampNow()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(config.Host, fmt.Sprintf("%d", config.Port)), helpers.DurationFromSecond(timeout))
	if err != nil {
		return newTcpError(schedulerId, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_Error, wrongConnectConfigError.Error())
	}
	if conn != nil {
		defer func() {
			_ = conn.Close()
		}()
	}
	return newTcpError(schedulerId, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_OK, "")
}
