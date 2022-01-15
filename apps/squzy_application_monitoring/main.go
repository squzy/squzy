package main

import (
	"context"
	"github.com/squzy/mongo_helper"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"github.com/squzy/squzy/apps/squzy_application_monitoring/application"
	"github.com/squzy/squzy/apps/squzy_application_monitoring/config"
	"github.com/squzy/squzy/apps/squzy_application_monitoring/database"
	"github.com/squzy/squzy/apps/squzy_application_monitoring/server"
	_ "github.com/squzy/squzy/apps/squzy_application_monitoring/version"
	"github.com/squzy/squzy/internal/grpctools"
	"github.com/squzy/squzy/internal/helpers"
	logger "github.com/squzy/squzy/internal/logger"
)

func main() {
	cfg := config.New()
	tools := grpctools.New()
	conn, err := tools.GetConnection(cfg.GetStorageHost(), 0, grpc.WithInsecure())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		_ = conn.Close()
	}()
	storageClient := apiPb.NewStorageClient(conn)

	ctx, cancel := helpers.TimeoutContext(context.Background(), 0)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.GetMongoURI()))
	if err != nil {
		logger.Fatal(err.Error())
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()
	connector := mongo_helper.New(client.Database(cfg.GetMongoDb()).Collection(cfg.GetMongoCollection()))

	logger.Fatal(application.New(server.New(database.New(connector), cfg, storageClient)).Run(cfg.GetPort()).Error())
}
