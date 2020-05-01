package main

import (
	"google.golang.org/grpc"
	"log"
	"squzy/internal/grpcTools"
	"squzy/internal/httpTools"
	"squzy/internal/parsers"
	scheduler_storage "squzy/internal/scheduler-storage"
	"squzy/internal/semaphore"
	sitemap_storage "squzy/internal/sitemap-storage"
	"squzy/internal/storage"
	"squzy/apps/squzy_monitoring/application"
	"squzy/apps/squzy_monitoring/config"
	"squzy/apps/squzy_monitoring/version"
	"time"
)

func main() {
	httpPackage := httpTools.New(version.GetVersion())
	grpcTool := grpcTools.New()
	cfg := config.New()
	app := application.New(
		scheduler_storage.New(),
		storage.NewExternalStorage(grpcTool, cfg.GetClientAddress(), cfg.GetStorageTimeout(), storage.GetInMemoryStorage(), grpc.WithInsecure()),
		sitemap_storage.New(
			time.Hour*24,
			httpPackage,
			parsers.NewSiteMapParser(),
		),
		httpPackage,
		semaphore.NewSemaphore,
	)
	log.Fatal(app.Run(cfg.GetPort()))
}
