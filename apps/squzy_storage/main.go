package main

import (
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"squzy/apps/squzy_storage/application"
	"squzy/apps/squzy_storage/config"
	"squzy/apps/squzy_storage/server"
	_ "squzy/apps/squzy_storage/version"
	"squzy/internal/database"
	"squzy/internal/grpctools"
	"squzy/internal/logger"
)

func main() {
	tools := grpctools.New()
	cfg := config.New()
	connect, err := sql.Open("clickhouse", fmt.Sprintf("tcp://%s:%s?username=%s&password=%s&database=%s&read_timeout=10&write_timeout=20",
		cfg.GetDbHost(),
		cfg.GetDbPort(),
		cfg.GetDbUser(),
		cfg.GetDbPassword(),
		cfg.GetDbName(),
	))
	if err != nil {
		logger.Fatal(err.Error())
	}
	db := database.New(connect)

	err = db.Migrate()
	if err != nil {
		logger.Fatal(err.Error())
	}

	incidentConn, err := tools.GetConnection(cfg.GetIncidentServerAddress(), 0, grpc.WithInsecure())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		_ = incidentConn.Close()
	}()
	incidentClient := apiPb.NewIncidentServerClient(incidentConn)

	apiService := server.NewServer(db, incidentClient, cfg)
	storageServ := application.NewApplication(cfg, apiService)
	logger.Fatal(storageServ.Run().Error())
}
