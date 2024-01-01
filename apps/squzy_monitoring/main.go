package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/squzy/mongo_helper"
	"github.com/squzy/squzy/apps/squzy_monitoring/application"
	"github.com/squzy/squzy/apps/squzy_monitoring/config"
	"github.com/squzy/squzy/apps/squzy_monitoring/version"
	"github.com/squzy/squzy/internal/cache"
	"github.com/squzy/squzy/internal/grpctools"
	"github.com/squzy/squzy/internal/helpers"
	"github.com/squzy/squzy/internal/httptools"
	"github.com/squzy/squzy/internal/job"
	job_executor "github.com/squzy/squzy/internal/job-executor"
	"github.com/squzy/squzy/internal/logger"
	"github.com/squzy/squzy/internal/parsers"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	scheduler_storage "github.com/squzy/squzy/internal/scheduler-storage"
	"github.com/squzy/squzy/internal/semaphore"
	sitemap_storage "github.com/squzy/squzy/internal/sitemap-storage"
	"github.com/squzy/squzy/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"time"
)

const (
	day = time.Hour * 24
)

func main() {
	cfg := config.New()
	ctx, cancel := helpers.TimeoutContext(context.Background(), 0)
	defer cancel()
	cache := getCache(cfg)

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
		job.ExecSSL,
		job.ExecCassandra,
		job.ExecMongo,
		job.ExecMysql,
		job.ExecPostgres,
	)
	app := application.New(
		scheduler_storage.New(),
		jobExecutor,
		configStorage,
		cache,
	)
	logger.Fatal(app.Run(cfg.GetPort()).Error())
}

func getCache(cfg config.Config) cache.Cache {
	var rdb *redis.Client

	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.GetCacheAddr(),
		Password: cfg.GetCachePassword(),
		DB:       int(cfg.GetCacheDB()),
	})

	c, err := cache.New(rdb)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return c
}
