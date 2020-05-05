package main

import (
	"context"
	"github.com/squzy/mongo_helper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"squzy/apps/squzy_agent_server/application"
	"squzy/apps/squzy_agent_server/config"
	"squzy/apps/squzy_agent_server/database"
	_ "squzy/apps/squzy_agent_server/version"
	"squzy/internal/helpers"
	"time"
)

func main() {
	cfg := config.New()
	ctx, cancel := helpers.TimeoutContext(context.Background(), time.Second*10)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.GetMongoUri()))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	connector := mongo_helper.New(client.Database(cfg.GetMongoDb()).Collection(cfg.GetMongoCollection()))
	app := application.New(database.New(connector))
	log.Fatal(app.Run(cfg.GetPort()))
}
