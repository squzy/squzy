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
	schedulerID string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerCode
	description string
}

func (s *grpcError) GetLogData() *apiPb.SchedulerResponse {
	var err *apiPb.SchedulerSnapshot_Error
	if s.code == apiPb.SchedulerCode_Error {
		err = &apiPb.SchedulerSnapshot_Error{
			Message: s.description,
		}
	}
	return &apiPb.SchedulerResponse{
		SchedulerId: s.schedulerID,
		Snapshot: &apiPb.SchedulerSnapshot{
			Code:  s.code,
			Error: err,
			Type:  apiPb.SchedulerType_Grpc,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: s.startTime,
				EndTime:   s.endTime,
			},
		},
	}
}

func newGrpcError(schedulerID string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerCode, description string) CheckError {
	return &grpcError{
		schedulerID: schedulerID,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
	}
}

func ExecGrpc(schedulerID string, timeout int32, config *scheduler_config_storage.GrpcConfig, opts ...grpc.DialOption) CheckError {
	startTime := ptypes.TimestampNow()

	ctx, cancel := helpers.TimeoutContext(context.Background(), helpers.DurationFromSecond(timeout))

	defer cancel()

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", config.Host, config.Port), opts...)

	if err != nil {
		return newGrpcError(schedulerID, startTime, ptypes.TimestampNow(), apiPb.SchedulerCode_Error, errWrongConnectConfigError.Error())
	}

	defer func() {
		_ = conn.Close()
	}()

	client := health_check.NewHealthClient(conn)

	md := metadata.New(map[string]string{
		logMetaData: schedulerID,
	})

	res, err := client.Check(metadata.NewOutgoingContext(ctx, md), &health_check.HealthCheckRequest{Service: config.Service})

	if err != nil {
		return newGrpcError(schedulerID, startTime, ptypes.TimestampNow(), apiPb.SchedulerCode_Error, errConnTimeoutError.Error())
	}

	if res.Status != health_check.HealthCheckResponse_SERVING {
		return newGrpcError(schedulerID, startTime, ptypes.TimestampNow(), apiPb.SchedulerCode_Error, errGrpcNotServing.Error())
	}
	return newGrpcError(schedulerID, startTime, ptypes.TimestampNow(), apiPb.SchedulerCode_OK, "")
}
