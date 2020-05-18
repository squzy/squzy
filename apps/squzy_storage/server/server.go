package server

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"net"
	"squzy/apps/squzy_storage/config"
)

type Server interface {
	Run() error
}

type server struct {
	config  config.Config
	apiServ apiPb.StorageServer
}

func NewServer(cnfg config.Config, apiServ apiPb.StorageServer) Server {
	return &server{
		config:  cnfg,
		apiServ: apiServ,
	}
}

func (s *server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GetPort()))
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
	apiPb.RegisterStorageServer(grpcServer, s.apiServ)
	return grpcServer.Serve(lis)
}
