package executable

import (
	"fmt"
	"github.com/jinzhu/gorm"
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

// CliExecute is called from cli with map[string]interface{} parameter
// The last one is read from config file by Viper
func CliExecute(cliConfig map[string]interface{}) {
	fmt.Println(cliConfig)
	cfg, err := config.ReadFromMap(cliConfig)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error reading cfg: %s", err.Error()))
	}
	Execute(cfg)
}

// Execute is the function, which initialize all the storage stuff
func Execute(cfg config.Config) {
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

	db := database.New(postgresDb.LogMode(cfg.WithDbLogs()))

	err = db.Migrate()
	if err != nil {
		logger.Fatal(err.Error())
	}

	var apiService apiPb.StorageServer
	if cfg.WithIncident() {
		// In case we have incident, we need initialize incidentClient
		incidentConn, err := tools.GetConnection(cfg.GetIncidentServerAddress(), 0, grpc.WithInsecure())
		if err != nil {
			logger.Fatal(err.Error())
		}
		defer func() {
			_ = incidentConn.Close()
		}()
		incidentClient := apiPb.NewIncidentServerClient(incidentConn)
		apiService = server.NewServer(db, incidentClient, cfg)
	} else {
		// IncidentClient nil in other case
		apiService = server.NewServer(db, nil, cfg)
	}

	storageServ := application.NewApplication(cfg, apiService)
	logger.Fatal(storageServ.Run().Error())
}

