package application

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"google.golang.org/grpc"
	"net"
)

type app struct {
	server apiPb.ApplicationMonitoringServer
}

func New(server apiPb.ApplicationMonitoringServer) *app {
	return &app{
		server: server,
	}
}

func (a *app) Run(port int32) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
		),
		),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		),
		),
	)
	apiPb.RegisterApplicationMonitoringServer(grpcServer, a.server)
	return grpcServer.Serve(lis)
}
