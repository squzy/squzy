package main

import (
	"context"
	"github.com/squzy/mongo_helper"
	"github.com/squzy/squzy/apps/squzy_agent_server/application"
	"github.com/squzy/squzy/apps/squzy_agent_server/config"
	"github.com/squzy/squzy/apps/squzy_agent_server/database"
	"github.com/squzy/squzy/apps/squzy_agent_server/server"
	_ "github.com/squzy/squzy/apps/squzy_agent_server/version"
	"github.com/squzy/squzy/internal/grpctools"
	"github.com/squzy/squzy/internal/helpers"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func main() {
	cfg := config.New()
	tools := grpctools.New()
	conn, err := tools.GetConnection(cfg.GetStorageAddress(), 0, grpc.WithInsecure())
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
	db := database.New(connector)

	// Before start left setup all agents to unregister
	unregisterAll(db)
	app := application.New(server.New(db, storageClient))
	logger.Fatal(app.Run(cfg.GetPort()).Error())
}

func unregisterAll(db database.Database) {
	ctx, cancel := helpers.TimeoutContext(context.Background(), time.Minute)
	defer cancel()
	list, err := db.GetAll(ctx, bson.M{
		"status": bson.M{
			"$ne": apiPb.AgentStatus_UNREGISTRED,
		},
	})
	if err != nil {
		return
	}
	for _, v := range list {
		timeoutRq := time.Second * 10
		rqCtx, cn := helpers.TimeoutContext(context.Background(), timeoutRq)
		defer cn()
		id, err := primitive.ObjectIDFromHex(v.Id)
		if err != nil {
			continue
		}
		_ = db.UpdateStatus(rqCtx, id, apiPb.AgentStatus_UNREGISTRED, timestamp.Now())
	}
}
