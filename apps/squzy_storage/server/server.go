package server

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	panichandler "github.com/kazegusuri/grpc-panic-handler"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"net"
	"squzy/apps/squzy_storage/config"
)

type Server interface {
	Run(opts ...grpc.DialOption) error
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

func (s *server) Run(opts ...grpc.DialOption) error {
	serv, err := s.newServer(s.config)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.config.GetPort()))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			panichandler.UnaryPanicHandler,
		),
		grpc_middleware.WithStreamServerChain(
			panichandler.StreamPanicHandler,
		),
	)
	apiPb.RegisterStorageServer(grpcServer, serv)
	return grpcServer.Serve(lis)
}
