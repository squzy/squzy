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
	var db database.Database
	tools := grpctools.New()
	cfg := config.New()

	db, err := getDatabase(cfg, db)

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

func getDatabase(cfg config.Config, db database.Database) (database.Database, error) {
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

		db, err = database.New(postgresDb)
		if err != nil {
			logger.Fatal(err.Error())
		}
	}
	connect, err := sql.Open("clickhouse", fmt.Sprintf("tcp://%s:%s?username=%s&password=%s&database=%s&read_timeout=10&write_timeout=20",
		cfg.GetDbHost(),
		cfg.GetDbPort(),
		cfg.GetDbUser(),
		cfg.GetDbPassword(),
		cfg.GetDbName(),
	))
	db, err = database.New(connect)

	err = db.Migrate()
	if err != nil {
		logger.Fatal(err.Error())
	}
	return db, err
}
