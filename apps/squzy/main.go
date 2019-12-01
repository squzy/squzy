package main

import (
	"log"
	"squzy/apps/internal/grpcTools"
	"squzy/apps/internal/httpTools"
	"squzy/apps/internal/parsers"
	scheduler_storage "squzy/apps/internal/scheduler-storage"
	sitemap_storage "squzy/apps/internal/sitemap-storage"
	"squzy/apps/internal/storage"
	"time"
)

func init() {
	ReadConfig()
}

func main() {
	httpPackage := httpTools.New()
	grpcTool := grpcTools.New()
	application := New(
		scheduler_storage.New(),
		storage.NewExternalStorage(grpcTool, clientAddress, timeoutStorage, storage.GetInMemoryStorage()),
		sitemap_storage.New(
			time.Hour*24,
			httpPackage,
			parsers.NewSiteMapParser(),
		),
		httpPackage,
	)
	log.Fatal(application.Run(port))
}
