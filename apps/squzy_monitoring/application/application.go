package application

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"net"
	"squzy/internal/httpTools"
	scheduler_storage "squzy/internal/scheduler-storage"
	"squzy/internal/semaphore"
	sitemap_storage "squzy/internal/sitemap-storage"
	"squzy/internal/storage"
	"squzy/apps/squzy_monitoring/server"
)

type app struct {
	schedulerStorage scheduler_storage.SchedulerStorage
	externalStorage  storage.Storage
	siteMapStorage   sitemap_storage.SiteMapStorage
	tool             httpTools.HttpTool
	semaphoreFactory semaphore.SemaphoreFactory
}

func New(
	schedulerStorage scheduler_storage.SchedulerStorage,
	externalStorage storage.Storage,
	siteMapStorage sitemap_storage.SiteMapStorage,
	tool httpTools.HttpTool,
	semaphoreFactory semaphore.SemaphoreFactory,
) *app {
	return &app{
		schedulerStorage,
		externalStorage,
		siteMapStorage,
		tool,
		semaphoreFactory,
	}
}

func (s *app) Run(port int32) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			),
		),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			),
		),
	)
	apiPb.RegisterSchedulersExecutorServer(
		grpcServer,
		server.New(
			s.schedulerStorage,
			s.externalStorage,
			s.siteMapStorage,
			s.tool,
			s.semaphoreFactory,
		),
	)
	return grpcServer.Serve(lis)
}
