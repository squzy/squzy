package job

import (
	"context"
	"google.golang.org/grpc"
	health_check "google.golang.org/grpc/health/grpc_health_v1"
)

type grpcJob struct {
	service     string
	address     string
	connOptions []grpc.DialOption
	callOptions []grpc.CallOption
}

func NewGrpcJob(service string, address string, connOptions []grpc.DialOption, callOptions []grpc.CallOption) Job {
	return &grpcJob{
		service:     service,
		address:     address,
		callOptions: callOptions,
		connOptions: connOptions,
	}
}

func (j *grpcJob) Do() error {
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)

	defer cancel()

	conn, err := grpc.DialContext(ctx, j.address, j.connOptions...)

	if err != nil {
		return wrongConnectConfigError
	}

	defer func() {
		_ = conn.Close()
	}()

	client := health_check.NewHealthClient(conn)

	reqCtx, cancelCtx := context.WithTimeout(context.Background(), connTimeout)

	defer cancelCtx()

	res, err := client.Check(reqCtx, &health_check.HealthCheckRequest{Service: j.service}, j.callOptions...)

	if err != nil {
		return connTimeoutError
	}

	if res.Status != health_check.HealthCheckResponse_SERVING {
		return grpcNotServing
	}
	return nil
}
