package main

import (
	"google.golang.org/grpc"
	"log"
	"squzy/apps/internal/grpcTools"
	"squzy/apps/internal/httpTools"
	"squzy/apps/internal/parsers"
	scheduler_storage "squzy/apps/internal/scheduler-storage"
	"squzy/apps/internal/semaphore"
	sitemap_storage "squzy/apps/internal/sitemap-storage"
	"squzy/apps/internal/storage"
	"squzy/apps/squzy/application"
	"squzy/apps/squzy/config"
	"squzy/apps/squzy/version"
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
