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
	agentServerClient := apiPb.NewAgentServerClient(agentServerConn)
	monitoringConn, err := tools.GetConnection(cfg.GetMonitoringServerAddress(), 0, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	monitoringClient := apiPb.NewSchedulersExecutorClient(monitoringConn)
	log.Fatal(
		router.New(
			handlers.New(agentServerClient, monitoringClient),
		).GetEngine().Run(fmt.Sprintf(":%d", cfg.GetPort())),
	)
}
