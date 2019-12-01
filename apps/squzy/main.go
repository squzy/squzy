package main

import (
	"context"
	"fmt"
	serverPb "github.com/squzy/squzy_generated/generated/server/proto/v1"
	storagePb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"squzy/apps/internal/httpTools"
	"squzy/apps/internal/parsers"
	scheduler_storage "squzy/apps/internal/scheduler-storage"
	sitemap_storage "squzy/apps/internal/sitemap-storage"
	"squzy/apps/internal/storage"
	"squzy/apps/squzy/server"
	"strconv"
	"time"
)

func runServer(
	port int32,
	schedulerStorage scheduler_storage.SchedulerStorage,
	externalStorage storage.Storage,
	siteMapStorage sitemap_storage.SiteMapStorage,
	tool httpTools.HttpTool,
) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	serverPb.RegisterServerServer(
		grpcServer,
		server.New(
			schedulerStorage,
			externalStorage,
			siteMapStorage,
			tool,
		),
	)
	return grpcServer.Serve(lis)
}

type app struct {
	schedulerStorage scheduler_storage.SchedulerStorage
	externalStorage  storage.Storage
	siteMapStorage   sitemap_storage.SiteMapStorage
	tool             httpTools.HttpTool
}

func New(
	schedulerStorage scheduler_storage.SchedulerStorage,
	externalStorage storage.Storage,
	siteMapStorage sitemap_storage.SiteMapStorage,
	tool httpTools.HttpTool,
) *app {
	return &app{
		schedulerStorage,
		externalStorage,
		siteMapStorage,
		tool,
	}
}

func (s *app) Run(port int32) error {
	return runServer(
		port,
		s.schedulerStorage,
		s.externalStorage,
		s.siteMapStorage,
		s.tool,
	)
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
	var port int32
	portValue := os.Getenv("PORT")
	port = 8080
	if portValue != "" {
		i, err := strconv.ParseInt(portValue, 10, 32)
		if err != nil {
			panic(err)
		}
		port = int32(i)
	}
	log.Fatal(application.Run(port))
}

func getExternalStorage() storage.Storage {
	clientAddress := os.Getenv("STORAGE_HOST")
	timeoutValue := os.Getenv("STORAGE_TIMEOUT")
	var timeout int32
	timeout = 5
	if timeoutValue != "" {
		i, err := strconv.ParseInt(timeoutValue, 10, 32)
		if err != nil {
			panic(err)
		}
		timeout = int32(i)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
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
