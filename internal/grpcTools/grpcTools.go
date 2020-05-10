package grpcTools

import (
	"context"
	"google.golang.org/grpc"
	"squzy/internal/helpers"
	"time"
)

type grpcTool struct {
}

func (g grpcTool) GetConnection(address string, timeout time.Duration, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := helpers.TimeoutContext(context.Background(), timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, options...)

	if err != nil {
		return nil, err
	}
	return conn, nil
}

type GrpcTool interface {
	GetConnection(address string, timeout time.Duration, option ...grpc.DialOption) (*grpc.ClientConn, error)
}

func New() GrpcTool {
	return &grpcTool{}
}
