package application

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"net"
)

type Application interface {
	Run(port int32) error
}

type application struct {
	incidentServ apiPb.IncidentServerServer
}

func NewApplication(incidentServ apiPb.IncidentServerServer) Application {
	return &application{
		incidentServ: incidentServ,
	}
}

func (s *application) Run(port int32) error {
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
	apiPb.RegisterIncidentServerServer(grpcServer, s.incidentServ)
	return grpcServer.Serve(lis)
}
