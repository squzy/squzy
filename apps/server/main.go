package main

import (
	"log"
	"squzy/apps/internal/database"
	"squzy/apps/server/application"
	"squzy/apps/server/config"
	_ "squzy/apps/server/version"
)

func main() {
	cfg := config.New()
	db, err := database.New()
	if err != nil {
		log.Fatalf(err.Error())
	}
	app := application.New(db)
	log.Fatal(app.Run(cfg.GetPort()))
}
