package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"github.com/squzy/squzy/apps/squzy_storage/application"
	"github.com/squzy/squzy/apps/squzy_storage/config"
	"github.com/squzy/squzy/apps/squzy_storage/server"
	_ "github.com/squzy/squzy/apps/squzy_storage/version"
	"github.com/squzy/squzy/internal/database"
	"github.com/squzy/squzy/internal/grpctools"
	"github.com/squzy/squzy/internal/logger"
)

func main() {
	tools := grpctools.New()
	cfg := config.New()
	postgresDb, err := gorm.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s dbname=%s user=%s  password=%s connect_timeout=10 sslmode=disable",
			cfg.GetDbHost(),
			cfg.GetDbPort(),
			cfg.GetDbName(),
			cfg.GetDbUser(),
			cfg.GetDbPassword(),
		))

	if err != nil {
		logger.Fatal(err.Error())
	}

	db := database.New(postgresDb.LogMode(cfg.WithDbLogs()))

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
