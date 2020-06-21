package main

import (
	"context"
	"github.com/squzy/mongo_helper"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"squzy/apps/squzy_application_monitoring/application"
	"squzy/apps/squzy_application_monitoring/config"
	"squzy/apps/squzy_application_monitoring/database"
	"squzy/apps/squzy_application_monitoring/server"
	_ "squzy/apps/squzy_application_monitoring/version"
	"squzy/internal/grpctools"
	"squzy/internal/helpers"
)

func main() {
	cfg := config.New()
	tools := grpctools.New()
	conn, err := tools.GetConnection(cfg.GetStorageHost(), 0, grpc.WithInsecure())
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
	defer func() {
		_ = client.Disconnect(context.Background())
	}()
	connector := mongo_helper.New(client.Database(cfg.GetMongoDb()).Collection(cfg.GetMongoCollection()))

	log.Fatal(application.New(server.New(database.New(connector), cfg, storageClient)).Run(cfg.GetPort()))
}
