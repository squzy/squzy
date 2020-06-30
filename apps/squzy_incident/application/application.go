package application

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"net"
	"squzy/apps/squzy_incident/config"
)

type Application interface {
	Run() error
}

type application struct {
	config  config.Config
	apiServ apiPb.IncidentServerServer
}

func NewApplication(cnfg config.Config, apiServ apiPb.IncidentServerServer) Application {
	return &application{
		config:  cnfg,
		apiServ: apiServ,
	}
}

func (s *application) Run() error {
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
	apiPb.RegisterIncidentServerServer(grpcServer, s.apiServ)
	return grpcServer.Serve(lis)
}
