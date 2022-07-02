package main

import (
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jinzhu/gorm"
	"github.com/squzy/squzy/apps/squzy_storage/application"
	"github.com/squzy/squzy/apps/squzy_storage/config"
	"github.com/squzy/squzy/apps/squzy_storage/server"
	_ "github.com/squzy/squzy/apps/squzy_storage/version"
	"github.com/squzy/squzy/internal/database"
	"github.com/squzy/squzy/internal/grpctools"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"google.golang.org/grpc"
	"os"
)

func main() {
	tools := grpctools.New()
	cfg := config.New()

	db, err := getDatabase(cfg)

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

func getDatabase(cfg config.Config) (database.Database, error) {
	var db database.Database

	if dt, ok := os.LookupEnv("DB_TYPE"); ok && dt == "postgres" {
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

		db, err = database.New(postgresDb)
		if err != nil {
			logger.Fatal(err.Error())
		}
	} else {
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
		db, err = database.New(connect)
		if err != nil {
			logger.Fatal(err.Error())
		}
	}

	err := db.Migrate()
	if err != nil {
		logger.Fatal(err.Error())
	}
	return db, err
}
