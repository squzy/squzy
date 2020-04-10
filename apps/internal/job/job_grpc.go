package job

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"google.golang.org/grpc"
	health_check "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"squzy/apps/internal/helpers"
)

const (
	logMetaData = "Squzy_log_id"
)

type grpcJob struct {
	service     string
	host        string
	port        int32
	timeout     int32
	connOptions []grpc.DialOption
	callOptions []grpc.CallOption
}

func NewGrpcJob(service string, host string, port int32, timeout int32, connOptions []grpc.DialOption, callOptions []grpc.CallOption) Job {
	return &grpcJob{
		service:     service,
		host:        host,
		port:        port,
		timeout:     timeout,
		callOptions: callOptions,
		connOptions: connOptions,
	}
}

type grpcError struct {
	logId       string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        clientPb.StatusCode
	description string
	location    string
	port        int32
}

func (s *grpcError) GetLogData() *clientPb.Log {
	return &clientPb.Log{
		Code:        s.code,
		Description: s.description,
		Meta: &clientPb.MetaData{
			Id:        s.logId,
			Location:  s.location,
			Port:      s.port,
			StartTime: s.startTime,
			EndTime:   s.endTime,
			Type:      clientPb.Type_Grpc,
		},
	}
}

func newGrpcError(logId string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code clientPb.StatusCode, description string, location string, port int32) CheckError {
	return &grpcError{
		logId:       logId,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		location:    location,
		port:        port,
	}
}

func (j *grpcJob) Do() CheckError {
	logId := uuid.New().String()
	startTime := ptypes.TimestampNow()

	ctx, cancel := helpers.TimeoutContext(context.Background(), helpers.DurationFromSecond(j.timeout))

	defer cancel()

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", j.host, j.port), j.connOptions...)

	if err != nil {
		return newGrpcError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_Error, wrongConnectConfigError.Error(), j.host, j.port)
	}

	defer func() {
		_ = conn.Close()
	}()

	client := health_check.NewHealthClient(conn)

	md := metadata.New(map[string]string{
		logMetaData: logId,
	})

	res, err := client.Check(metadata.NewOutgoingContext(ctx, md), &health_check.HealthCheckRequest{Service: j.service}, j.callOptions...)

	if err != nil {
		return newGrpcError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_Error, connTimeoutError.Error(), j.host, j.port)
	}

	if res.Status != health_check.HealthCheckResponse_SERVING {
		return newGrpcError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_Error, grpcNotServing.Error(), j.host, j.port)
	}
	return newGrpcError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_OK, "", j.host, j.port)
}
