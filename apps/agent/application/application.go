package application

import (
	"context"
	"fmt"
	agentPb "github.com/squzy/squzy_generated/generated/agent/proto/v1"
	"google.golang.org/grpc"
	"log"
	"squzy/apps/agent/config"
	agent_executor "squzy/apps/internal/agent-executor"
	"squzy/apps/internal/grpcTools"
)

type application struct {
	id        string
	server    *grpc.Server
	executor  agent_executor.AgentExecutor
	grpcTools grpcTools.GrpcTool
	config    config.Config
	client    agentPb.AgentServerClient
}

func (a *application) Run() error {
	conn, err := a.grpcTools.GetConnection(a.config.GetSquzyServer(), a.config.GetSquzyServerTimeout(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("Can't connect to squzy server %s", a.config.GetSquzyServer())
	}
	a.client = agentPb.NewAgentServerClient(conn)
	stream, _ := a.client.Register(context.Background())

	response, err := stream.Recv()
	if err != nil {
		return err
	}

	a.id = response.Id
	log.Printf("Registred with ID=%s", a.id)

	go func() {
		for stat := range a.executor.Execute() {
			// what we should do if squzy server cant get msg
			_ = stream.Send(stat)
		}
	}()

	return nil
}

type Application interface {
	Run() error
}

func New(executor agent_executor.AgentExecutor, grpcTools grpcTools.GrpcTool, config config.Config) Application {
	return &application{
		config:    config,
		executor:  executor,
		grpcTools: grpcTools,
		server:    grpc.NewServer(),
		client:    nil,
	}
}
