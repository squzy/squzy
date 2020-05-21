package main

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/squzy/mongo_helper"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	"time"
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
	defer func() {
		_ = client.Disconnect(context.Background())
	}()
	connector := mongo_helper.New(client.Database(cfg.GetMongoDb()).Collection(cfg.GetMongoCollection()))
	db := database.New(connector)

	// Before start left setup all agents to unregister
	unregisterAll(db)
	app := application.New(server.New(db, storageClient))
	log.Fatal(app.Run(cfg.GetPort()))
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
		_ = db.UpdateStatus(rqCtx, id, apiPb.AgentStatus_UNREGISTRED, ptypes.TimestampNow())
	}
}
