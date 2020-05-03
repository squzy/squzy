package application

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"net"
	"os"
	"squzy/apps/squzy_monitoring/server"
	"squzy/internal/helpers"
	job_executor "squzy/internal/job-executor"
	"squzy/internal/scheduler"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
	scheduler_storage "squzy/internal/scheduler-storage"
)

type app struct {
	schedulerStorage scheduler_storage.SchedulerStorage
	jobExecutor      job_executor.JobExecutor
	configStorage    scheduler_config_storage.Storage
}

func New(
	schedulerStorage scheduler_storage.SchedulerStorage,
	jobExecutor job_executor.JobExecutor,
	configStorage scheduler_config_storage.Storage,
) *app {
	return &app{
		schedulerStorage: schedulerStorage,
		jobExecutor:      jobExecutor,
		configStorage:    configStorage,
	}
}

func (s *app) SyncOne(config *scheduler_config_storage.SchedulerConfig) error {
	sched, err := scheduler.New(config.Id, helpers.DurationFromSecond(config.Interval), s.jobExecutor)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("SchedulerId: %s cant synced, error in config", config.Id.Hex()))
		// @TODO logger here
		return err
	}
	err = s.schedulerStorage.Set(sched)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("SchedulerId: %s cant synced, error in memory storage", config.Id.Hex()))
		// @TODO logger here
		return err
	}
	if config.Status == apiPb.SchedulerStatus_STOPPED {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("SchedulerId: %s synced and STOP", config.Id.Hex()))
		return nil
	}
	if config.Status == apiPb.SchedulerStatus_RUNNED {
		sched.Run()
		fmt.Fprintln(os.Stdout, fmt.Sprintf("SchedulerId: %s synced and RUN", config.Id.Hex()))
	}
	return nil
}

func (s *app) sync() error {
	configs, err := s.configStorage.GetAllForSync(context.Background())
	if err != nil {
		return err
	}
	for _, config := range configs {
		_ = s.SyncOne(config)
	}
	fmt.Fprintln(os.Stdout, "Sync done")
	return nil
}

func (s *app) Run(port int32) error {
	err := s.sync()
	if err != nil {
		return err
	}
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
			s.jobExecutor,
			s.configStorage,
		),
	)
	return grpcServer.Serve(lis)
}
