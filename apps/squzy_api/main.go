package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"github.com/squzy/squzy/apps/squzy_api/config"
	"github.com/squzy/squzy/apps/squzy_api/handlers"
	"github.com/squzy/squzy/apps/squzy_api/router"
	_ "github.com/squzy/squzy/apps/squzy_api/version"
	"github.com/squzy/squzy/internal/grpctools"
	"github.com/squzy/squzy/internal/logger"
)

func main() {
	cfg := config.New()
	gin.SetMode(gin.ReleaseMode)
	tools := grpctools.New()
	agentServerConn, err := tools.GetConnection(cfg.GetAgentServerAddress(), 0, grpc.WithInsecure())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		_ = agentServerConn.Close()
	}()
	agentServerClient := apiPb.NewAgentServerClient(agentServerConn)
	monitoringConn, err := tools.GetConnection(cfg.GetMonitoringServerAddress(), 0, grpc.WithInsecure())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		_ = monitoringConn.Close()
	}()
	monitoringClient := apiPb.NewSchedulersExecutorClient(monitoringConn)
	storageConn, err := tools.GetConnection(cfg.GetStorageServerAddress(), 0, grpc.WithInsecure())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		_ = storageConn.Close()
	}()
	storageClient := apiPb.NewStorageClient(storageConn)

	appMonConn, err := tools.GetConnection(cfg.GetApplicationMonitoringAddress(), 0, grpc.WithInsecure())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		_ = appMonConn.Close()
	}()
	appMonClient := apiPb.NewApplicationMonitoringClient(appMonConn)

	incidentConn, err := tools.GetConnection(cfg.GetIncidentServerAddress(), 0, grpc.WithInsecure())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		_ = incidentConn.Close()
	}()
	incidentClient := apiPb.NewIncidentServerClient(incidentConn)

	notificationConn, err := tools.GetConnection(cfg.GetNotificationServerAddress(), 0, grpc.WithInsecure())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		_ = notificationConn.Close()
	}()
	notificicationClient := apiPb.NewNotificationManagerClient(notificationConn)

	logger.Fatal(
		router.New(
			handlers.New(agentServerClient, monitoringClient, storageClient, appMonClient, incidentClient, notificicationClient),
		).GetEngine().Run(fmt.Sprintf(":%d", cfg.GetPort())).Error(),
	)
}
