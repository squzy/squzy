package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"squzy/apps/internal/database"
	"squzy/apps/server/application"
	"squzy/apps/server/config"
	_ "squzy/apps/server/version"
)

func main() {
	cfg := config.New()
	getDB := func() (*gorm.DB, error) {
		return gorm.Open(
			"postgres",
			fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s connect_timeout=10 sslmode=disable",
				cfg.GetDbHost(),
				cfg.GetDbPort(),
				cfg.GetDbUser(),
				cfg.GetDbName(),
				cfg.GetDbPassword(),
			),
		)
	}
	db, err := database.New(getDB)
	if err != nil {
		log.Fatalf(err.Error())
	}
	app := application.New(db)
	log.Fatal(app.Run(cfg.GetPort()))
}
