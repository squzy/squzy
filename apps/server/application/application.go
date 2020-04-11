package application

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	agentPb "github.com/squzy/squzy_generated/generated/agent/proto/v1"
	storagePb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"google.golang.org/grpc"
	"net"
	"squzy/apps/internal/database"
	"squzy/apps/server/services/agents"
	"squzy/apps/server/services/checks"
)

type application struct {
	db database.Database
}

func New(db database.Database) *application {
	return &application{
		db: db,
	}
}

func (a *application) Run(port int32) error {
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
	storagePb.RegisterLoggerServer(grpcServer, checks.New(a.db))
	agentPb.RegisterAgentServerServer(grpcServer, agents.New(a.db))
	return grpcServer.Serve(lis)
}