package job

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/logger"
	"google.golang.org/grpc"
	health_check "google.golang.org/grpc/health/grpc_health_v1"
)

type grpcJob struct {
	service     string
	host        string
	port        int32
	connOptions []grpc.DialOption
	callOptions []grpc.CallOption
}

func NewGrpcJob(service string, host string, port int32, connOptions []grpc.DialOption, callOptions []grpc.CallOption) Job {
	return &grpcJob{
		service:     service,
		host:        host,
		port:        port,
		callOptions: callOptions,
		connOptions: connOptions,
	}
}

type grpcError struct {
	time        *timestamp.Timestamp
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
			Id:       uuid.New().String(),
			Location: s.location,
			Port:     s.port,
			Time:     s.time,
			Type:     clientPb.Type_Grpc,
		},
	}
}

func newGrpcError(time *timestamp.Timestamp, code clientPb.StatusCode, description string, location string, port int32) CheckError {
	return &grpcError{
		time:        time,
		code:        code,
		description: description,
		location:    location,
		port:        port,
	}
}

func (j *grpcJob) Do() CheckError {
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)

	defer cancel()

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", j.host, j.port), j.connOptions...)

	if err != nil {
		return newGrpcError(ptypes.TimestampNow(), clientPb.StatusCode_Error, wrongConnectConfigError.Error(), j.host, j.port)
	}

	defer func() {
		_ = conn.Close()
	}()

	client := health_check.NewHealthClient(conn)

	reqCtx, cancelCtx := context.WithTimeout(context.Background(), connTimeout)

	defer cancelCtx()

	res, err := client.Check(reqCtx, &health_check.HealthCheckRequest{Service: j.service}, j.callOptions...)

	if err != nil {
		return newGrpcError(ptypes.TimestampNow(), clientPb.StatusCode_Error, connTimeoutError.Error(), j.host, j.port)
	}

	if res.Status != health_check.HealthCheckResponse_SERVING {
		return newGrpcError(ptypes.TimestampNow(), clientPb.StatusCode_Error, grpcNotServing.Error(), j.host, j.port)
	}
	return newGrpcError(ptypes.TimestampNow(), clientPb.StatusCode_OK, "", j.host, j.port)
}
