package job

import (
	"context"
	"fmt"
	"github.com/squzy/squzy/internal/helpers"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"google.golang.org/grpc"
	health_check "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
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
			Type:  apiPb.SchedulerType_GRPC,
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
	startTime := timestamp.Now()

	ctx, cancel := helpers.TimeoutContext(context.Background(), helpers.DurationFromSecond(timeout))

	defer cancel()

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", config.Host, config.Port), opts...)

	if err != nil {
		return newGrpcError(schedulerID, startTime, timestamp.Now(), apiPb.SchedulerCode_ERROR, errWrongConnectConfigError.Error())
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
		return newGrpcError(schedulerID, startTime, timestamp.Now(), apiPb.SchedulerCode_ERROR, errConnTimeoutError.Error())
	}

	if res.Status != health_check.HealthCheckResponse_SERVING {
		return newGrpcError(schedulerID, startTime, timestamp.Now(), apiPb.SchedulerCode_ERROR, errGrpcNotServing.Error())
	}
	return newGrpcError(schedulerID, startTime, timestamp.Now(), apiPb.SchedulerCode_OK, "")
}
