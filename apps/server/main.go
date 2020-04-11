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
	app := application.New(database.New())
	log.Fatal(app.Run(cfg.GetPort()))
}
