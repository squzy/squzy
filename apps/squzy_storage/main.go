package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"log"
	"squzy/apps/squzy_storage/application"
	"squzy/apps/squzy_storage/config"
	"squzy/apps/squzy_storage/server"
	_ "squzy/apps/squzy_storage/version"
	"squzy/internal/database"
	"squzy/internal/grpctools"
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
		log.Fatal(err)
	}

	db := database.New(postgresDb.LogMode(true))

	err = db.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	incidentConn, err := tools.GetConnection(cfg.GetIncidentServerAddress(), 0, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = incidentConn.Close()
	}()
	incidentClient := apiPb.NewIncidentServerClient(incidentConn)

	apiService := server.NewServer(db, incidentClient, cfg)
	storageServ := application.NewApplication(cfg, apiService)
	log.Fatal(storageServ.Run())
}
