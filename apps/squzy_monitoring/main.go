package main

import (
	"context"
	"github.com/squzy/mongo_helper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"squzy/internal/logger"
	"squzy/apps/squzy_monitoring/application"
	"squzy/apps/squzy_monitoring/config"
	"squzy/apps/squzy_monitoring/version"
	"squzy/internal/grpctools"
	"squzy/internal/helpers"
	"squzy/internal/httptools"
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

const (
	day = time.Hour * 24
)

func main() {
	cfg := config.New()
	ctx, cancel := helpers.TimeoutContext(context.Background(), 0)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.GetMongoURI()))
	if err != nil {
		logger.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()
	connector := mongo_helper.New(client.Database(cfg.GetMongoDb()).Collection(cfg.GetMongoCollection()))
	httpPackage := httptools.New(version.GetVersion())
	grpcTool := grpctools.New()
	externalStorage := storage.NewExternalStorage(
		grpcTool,
		cfg.GetClientAddress(),
		cfg.GetStorageTimeout(),
		storage.GetInMemoryStorage(),
		grpc.WithInsecure(),
	)
	siteMapStorage := sitemap_storage.New(
		day,
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
		job.ExecTCP,
		job.ExecGrpc,
		job.ExecHTTP,
		job.ExecSiteMap,
		job.ExecHTTPValue,
	)
	app := application.New(
		scheduler_storage.New(),
		jobExecutor,
		configStorage,
	)
	logger.Fatal(app.Run(cfg.GetPort()))
}
