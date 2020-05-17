package main

import (
	"context"
	"github.com/squzy/mongo_helper"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"squzy/apps/squzy_agent_server/application"
	"squzy/apps/squzy_agent_server/config"
	"squzy/apps/squzy_agent_server/database"
	"squzy/apps/squzy_agent_server/server"
	_ "squzy/apps/squzy_agent_server/version"
	"squzy/internal/grpctools"
	"squzy/internal/helpers"
)

func main() {
	cfg := config.New()
	tools := grpctools.New()
	conn, err := tools.GetConnection(cfg.GetStorageAddress(), 0, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = conn.Close()
	}()
	storageClient := apiPb.NewStorageClient(conn)
	ctx, cancel := helpers.TimeoutContext(context.Background(), 0)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.GetMongoURI()))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	connector := mongo_helper.New(client.Database(cfg.GetMongoDb()).Collection(cfg.GetMongoCollection()))
	app := application.New(server.New(database.New(connector), storageClient))
	log.Fatal(app.Run(cfg.GetPort()))
}
