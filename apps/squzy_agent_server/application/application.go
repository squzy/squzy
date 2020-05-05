package application

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"net"
	"squzy/apps/squzy_agent_server/database"
	"squzy/apps/squzy_agent_server/server"
)

type app struct {
	db database.Database
}

func New(db database.Database) *app {
	return &app{
		db: db,
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
	apiPb.RegisterAgentServerServer(grpcServer, server.New(a.db, nil))
	return grpcServer.Serve(lis)
}
