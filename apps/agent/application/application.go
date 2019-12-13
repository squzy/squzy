package application

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/host"
	agentPb "github.com/squzy/squzy_generated/generated/agent/proto/v1"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"squzy/apps/agent/config"
	agent_executor "squzy/apps/internal/agent-executor"
	"squzy/apps/internal/grpcTools"
	"syscall"
)

type application struct {
	id          string
	server      *grpc.Server
	executor    agent_executor.AgentExecutor
	grpcTools   grpcTools.GrpcTool
	config      config.Config
	hostStatFn  func() (*host.InfoStat, error)
	getStreamFn func(agent agentPb.AgentServerClient) (agentPb.AgentServer_SendStatClient, error)
}

func (a *application) Run() error {
	hostStat, err := a.hostStatFn()
	if err != nil {
		return err
	}
	conn, err := a.grpcTools.GetConnection(a.config.GetSquzyServer(), a.config.GetSquzyServerTimeout(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("Can't connect to squzy server %s", a.config.GetSquzyServer())
	}
	ctx, cancel := context.WithTimeout(context.Background(), a.config.GetSquzyServerTimeout())
	defer cancel()
	client := agentPb.NewAgentServerClient(conn)
	res, err := client.Register(ctx, &agentPb.RegisterRequest{
		HostInfo: &agentPb.HostInfo{
			HostName: hostStat.Hostname,
			Os:       hostStat.OS,
			PlatformInfo: &agentPb.PlatformInfo{
				Name:    hostStat.Platform,
				Family:  hostStat.PlatformVersion,
				Version: hostStat.PlatformVersion,
			},
		},
	})
	if err != nil {
		return err
	}

	a.id = res.Id

	log.Printf("Registred with ID=%s", a.id)

	stream, err := a.getStreamFn(client)

	if err != nil {
		return err
	}

	go func() {
		for stat := range a.executor.Execute() {
			// what we should do if squzy server cant get msg
			_ = stream.Send(stat)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	<-interrupt

	_ = stream.CloseSend()

	ctxClose, cancelClose := context.WithTimeout(context.Background(), a.config.GetSquzyServerTimeout())
	defer cancelClose()
	_, _ = client.UnRegister(ctxClose, &agentPb.UnRegisterRequest{
		Id: a.id,
	})

	a.server.GracefulStop()

	return nil
}

type Application interface {
	Run() error
}

func New(executor agent_executor.AgentExecutor, grpcTools grpcTools.GrpcTool, config config.Config, hostStatFn func() (*host.InfoStat, error), getStreamFn func(agent agentPb.AgentServerClient) (agentPb.AgentServer_SendStatClient, error)) Application {
	return &application{
		config:      config,
		executor:    executor,
		grpcTools:   grpcTools,
		server:      grpc.NewServer(),
		hostStatFn:  hostStatFn,
		getStreamFn: getStreamFn,
	}
}

func NewStream(agent agentPb.AgentServerClient) (agentPb.AgentServer_SendStatClient, error) {
	return agent.SendStat(context.Background())
}
