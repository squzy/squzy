package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"squzy/apps/squzy_storage/application"
	"squzy/apps/squzy_storage/config"
	"squzy/apps/squzy_storage/server"
	_ "squzy/apps/squzy_storage/version"
	"squzy/internal/database"
)

func main() {
	cnfg := config.New()
	postgresDb, err := gorm.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s dbname=%s user=%s  password=%s connect_timeout=10 sslmode=disable",
			cnfg.GetDbHost(),
			cnfg.GetDbPort(),
			cnfg.GetDbName(),
			cnfg.GetDbUser(),
			cnfg.GetDbPassword(),
		))

	if err != nil {
		log.Fatal(err)
	}

	db := database.New(postgresDb.LogMode(false))

	err = db.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	apiService := server.NewServer(db)
	storageServ := application.NewApplication(cnfg, apiService)
	log.Fatal(storageServ.Run())
}
