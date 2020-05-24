package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"log"
	"squzy/apps/squzy_api/config"
	"squzy/apps/squzy_api/handlers"
	"squzy/apps/squzy_api/router"
	_ "squzy/apps/squzy_api/version"
	"squzy/internal/grpctools"
)

func main() {
	cfg := config.New()
	gin.SetMode(gin.ReleaseMode)
	tools := grpctools.New()
	agentServerConn, err := tools.GetConnection(cfg.GetAgentServerAddress(), 0, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = agentServerConn.Close()
	}()
	agentServerClient := apiPb.NewAgentServerClient(agentServerConn)
	monitoringConn, err := tools.GetConnection(cfg.GetMonitoringServerAddress(), 0, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = monitoringConn.Close()
	}()
	monitoringClient := apiPb.NewSchedulersExecutorClient(monitoringConn)
	storageConn, err := tools.GetConnection(cfg.GetStorageServerAddress(), 0, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = storageConn.Close()
	}()
	storageClient := apiPb.NewStorageClient(storageConn)

	appMonConn, err := tools.GetConnection(cfg.GetApplicationMonitoringAddress(), 0, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = appMonConn.Close()
	}()
	appMonClient := apiPb.NewApplicationMonitoringClient(appMonConn)

	log.Fatal(
		router.New(
			handlers.New(agentServerClient, monitoringClient, storageClient, appMonClient),
		).GetEngine().Run(fmt.Sprintf(":%d", cfg.GetPort())),
	)
}
