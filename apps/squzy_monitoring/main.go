package main

import (
	"context"
	"github.com/squzy/mongo_helper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"squzy/apps/squzy_monitoring/application"
	"squzy/apps/squzy_monitoring/config"
	"squzy/apps/squzy_monitoring/version"
	"squzy/internal/grpcTools"
	"squzy/internal/helpers"
	"squzy/internal/httpTools"
	"squzy/internal/job"
	job_executor "squzy/internal/job-executor"
	"squzy/internal/parsers"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
	scheduler_storage "squzy/internal/scheduler-storage"
	"squzy/internal/semaphore"
	sitemap_storage "squzy/internal/sitemap-storage"
	"squzy/internal/storage"
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
	httpPackage := httpTools.New(version.GetVersion())
	grpcTool := grpcTools.New()
	externalStorage := storage.NewExternalStorage(grpcTool, cfg.GetClientAddress(), cfg.GetStorageTimeout(), storage.GetInMemoryStorage(), grpc.WithInsecure())
	siteMapStorage := sitemap_storage.New(
		time.Hour*24,
		httpPackage,
		parsers.NewSiteMapParser(),
	)
	configStorage := scheduler_config_storage.New(connector)
	jobExecutor := job_executor.NewExecutor(
		externalStorage,
		siteMapStorage,
		httpPackage,
		semaphore.NewSemaphore,
		configStorage,
		job.ExecTcp,
		job.ExecGrpc,
		job.ExecHttp,
		job.ExecSiteMap,
		job.ExecHttpValue,
	)
	app := application.New(
		scheduler_storage.New(),
		jobExecutor,
		configStorage,
	)
	log.Fatal(app.Run(cfg.GetPort()))
}
