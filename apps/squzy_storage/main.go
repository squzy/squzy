package main

import (
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"io/ioutil"
	"squzy/apps/squzy_storage/application"
	"squzy/apps/squzy_storage/config"
	"squzy/apps/squzy_storage/server"
	_ "squzy/apps/squzy_storage/version"
	"squzy/internal/database"
	"squzy/internal/grpctools"
	"squzy/internal/logger"
)

func main() {
	var cfg config.Config

	filename := flag.String("config", "", "path to configFile")
	if filename == nil || *filename == "" {
		logger.Info(fmt.Sprintf("Empty config file param. Reading os env."))
		cfg = config.New()
	} else {
		// Reading config in case when flag provided
		cfgFromFile, err := readConfigFile(filename)
		if err != nil {
			logger.Fatal(fmt.Sprintf("error reading env file: %s", err.Error()))
		}
		cfg = cfgFromFile
	}

	tools := grpctools.New()
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

	db := database.New(postgresDb.LogMode(cfg.GetWithDbLogs()))

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

func readConfigFile(filename *string) (config.Config, error) {
	bytes, err := ioutil.ReadFile(*filename)
	if err != nil {
		return nil, fmt.Errorf("error reading cfg file: %w", err)
	}
	cfg, err := config.NewConfigFromYaml(bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing cfg file: %w", err)
	}
	return cfg, nil
}
