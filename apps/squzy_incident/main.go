package main

import (
	"context"
	"github.com/squzy/mongo_helper"
	"github.com/squzy/squzy/apps/squzy_incident/application"
	"github.com/squzy/squzy/apps/squzy_incident/config"
	"github.com/squzy/squzy/apps/squzy_incident/database"
	"github.com/squzy/squzy/apps/squzy_incident/server"
	_ "github.com/squzy/squzy/apps/squzy_incident/version"
	"github.com/squzy/squzy/internal/grpctools"
	"github.com/squzy/squzy/internal/helpers"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
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
	notificationConn, err := tools.GetConnection(cfg.GetNoticationServerHost(), 0, grpc.WithInsecure())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		_ = notificationConn.Close()
	}()
	notificationClient := apiPb.NewNotificationManagerClient(notificationConn)
	apiService := server.NewIncidentServer(notificationClient, storageClient, database.New(connector))
	storageServ := application.NewApplication(apiService)
	logger.Fatal(storageServ.Run(cfg.GetPort()).Error())
}
