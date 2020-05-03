package job

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	health_check "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"squzy/internal/helpers"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
)

const (
	logMetaData = "squzy_scheduler_id"
)

type grpcError struct {
	schedulerId string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerResponseCode
	description string
}

func (s *grpcError) GetLogData() *apiPb.SchedulerResponse {
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
		Type:        apiPb.SchedulerType_Grpc,
		Meta: &apiPb.SchedulerResponse_MetaData{
			StartTime: s.startTime,
			EndTime:   s.endTime,
		},
	}
}

func newGrpcError(schedulerId string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerResponseCode, description string) CheckError {
	return &grpcError{
		schedulerId: schedulerId,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
	}
}

func ExecGrpc(schedulerId string, timeout int32, config *scheduler_config_storage.GrpcConfig, opts ...grpc.DialOption) CheckError {
	startTime := ptypes.TimestampNow()

	ctx, cancel := helpers.TimeoutContext(context.Background(), helpers.DurationFromSecond(timeout))

	defer cancel()

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", config.Host, config.Port), opts...)

	if err != nil {
		return newGrpcError(schedulerId, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_Error, wrongConnectConfigError.Error())
	}

	defer func() {
		_ = conn.Close()
	}()

	client := health_check.NewHealthClient(conn)

	md := metadata.New(map[string]string{
		logMetaData: schedulerId,
	})

	res, err := client.Check(metadata.NewOutgoingContext(ctx, md), &health_check.HealthCheckRequest{Service: config.Service})

	if err != nil {
		return newGrpcError(schedulerId, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_Error, connTimeoutError.Error())
	}

	if res.Status != health_check.HealthCheckResponse_SERVING {
		return newGrpcError(schedulerId, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_Error, grpcNotServing.Error())
	}
	return newGrpcError(schedulerId, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_OK, "")
}
