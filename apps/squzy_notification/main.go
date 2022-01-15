package main

import (
	"context"
	"github.com/squzy/mongo_helper"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"github.com/squzy/squzy/apps/squzy_notification/application"
	"github.com/squzy/squzy/apps/squzy_notification/config"
	"github.com/squzy/squzy/apps/squzy_notification/database"
	"github.com/squzy/squzy/apps/squzy_notification/integrations"
	"github.com/squzy/squzy/apps/squzy_notification/server"
	"github.com/squzy/squzy/apps/squzy_notification/version"
	"github.com/squzy/squzy/internal/grpctools"
	"github.com/squzy/squzy/internal/helpers"
	"github.com/squzy/squzy/internal/httptools"
	"github.com/squzy/squzy/internal/logger"
)

func main() {
	tools := httptools.New(version.GetVersion())
	gtools := grpctools.New()
	cfg := config.New()
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
	conn, err := gtools.GetConnection(cfg.GetStorageHost(), 0, grpc.WithInsecure())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		_ = conn.Close()
	}()
	storageClient := apiPb.NewStorageClient(conn)
	logger.Fatal(
		application.New(
			server.New(
				database.NewList(mongo_helper.New(client.Database(cfg.GetMongoDB()).Collection(cfg.GetNotificationListCollection()))),
				database.NewMethods(mongo_helper.New(client.Database(cfg.GetMongoDB()).Collection(cfg.GetNotificationMethodCollection()))),
				storageClient,
				integrations.New(tools, cfg),
			),
		).Run(cfg.GetPort()).Error())
}
