package job

import (
	"crypto/tls"
	"fmt"
	"github.com/squzy/squzy/internal/helpers"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"net"
)

type sslError struct {
	schedulerID string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerCode
	description string
	value       *structpb.Value
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

func newSSLError(schedulerID string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerCode, description string, value *structpb.Value) CheckError {
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
	startTime := timestamp.Now()

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: helpers.DurationNotNegative(timeout)}, "tcp", net.JoinHostPort(config.Host, fmt.Sprintf("%d", config.Port)), cfg)

	if err != nil {
		return newSSLError(schedulerID, startTime, timestamp.Now(), apiPb.SchedulerCode_ERROR, err.Error(), nil)
	}

	defer func() {
		_ = conn.Close()
	}()

	crt := conn.ConnectionState().PeerCertificates[0]

	return newSSLError(schedulerID, startTime, timestamp.Now(), apiPb.SchedulerCode_OK, "", &structpb.Value{
		Kind: &structpb.Value_NumberValue{
			NumberValue: float64(crt.NotAfter.UnixNano()),
		},
	})
}
