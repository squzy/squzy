package job

import (
	"crypto/tls"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	structType "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"net"
	"squzy/internal/helpers"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
)

type sslError struct {
	schedulerID string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerCode
	description string
	value       *structType.Value
}

func (s *sslError) GetLogData() *apiPb.SchedulerResponse {
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
			Type:  apiPb.SchedulerType_SSL_EXPIRATION,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: s.startTime,
				EndTime:   s.endTime,
				Value:     s.value,
			},
		},
	}
}

func newSSLError(schedulerID string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerCode, description string, value *structType.Value) CheckError {
	return &sslError{
		schedulerID: schedulerID,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		value:       value,
	}
}

func ExecSSL(schedulerID string, timeout int32, config *scheduler_config_storage.SslExpirationConfig, cfg *tls.Config) CheckError {
	startTime := ptypes.TimestampNow()

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: helpers.DurationNotNegative(timeout)}, "tcp", net.JoinHostPort(config.Host, fmt.Sprintf("%d", config.Port)), cfg)

	if err != nil {
		return newSSLError(schedulerID, startTime, ptypes.TimestampNow(), apiPb.SchedulerCode_ERROR, err.Error(), nil)
	}

	defer func() {
		_ = conn.Close()
	}()

	crt := conn.ConnectionState().PeerCertificates[0]

	return newSSLError(schedulerID, startTime, ptypes.TimestampNow(), apiPb.SchedulerCode_OK, "", &structType.Value{
		Kind: &structType.Value_NumberValue{
			NumberValue: float64(crt.NotAfter.UnixNano()),
		},
	})
}
