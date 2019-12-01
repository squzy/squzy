package main

import (
	"context"
	storagePb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"google.golang.org/grpc"
	"log"
	"squzy/apps/internal/httpTools"
	"squzy/apps/internal/parsers"
	scheduler_storage "squzy/apps/internal/scheduler-storage"
	sitemap_storage "squzy/apps/internal/sitemap-storage"
	"squzy/apps/internal/storage"
	"time"
)

func getExternalStorage() storage.Storage {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutStorage)
	defer cancel()
	conn, err := grpc.DialContext(ctx, clientAddress, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Println("Will be use in memory storage")
		return storage.GetInMemoryStorage()
	}
	log.Println("Will be use client storage " + clientAddress)
	client := storagePb.NewLoggerClient(conn)
	return storage.NewExternalStorage(client)
}

func main() {
	httpPackage := httpTools.New()
	application := New(
		scheduler_storage.New(),
		getExternalStorage(),
		sitemap_storage.New(
			time.Hour*24,
			httpPackage,
			parsers.NewSiteMapParser(),
		),
		httpPackage,
	)
	log.Fatal(application.Run(port))
}
