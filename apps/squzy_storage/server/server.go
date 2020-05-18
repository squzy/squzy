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
	config    config.Config
	newServer func(cnfg config.Config) (apiPb.StorageServer, error)
}

func NewServer(cnfg config.Config, newServer func(cnfg config.Config) (apiPb.StorageServer, error)) Server {
	return &server{
		config:    cnfg,
		newServer: newServer,
	}
}

func (s *server) Run() error {
	serv, err := s.newServer(s.config) //call in main.go НАХУЙ НЕ НУЖНА
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.config.GetPort()))
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
	apiPb.RegisterStorageServer(grpcServer, serv)
	return grpcServer.Serve(lis)
}
