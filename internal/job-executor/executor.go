package job_executor

import (
	"context"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"squzy/internal/httpTools"
	"squzy/internal/job"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
	"squzy/internal/semaphore"
	sitemap_storage "squzy/internal/sitemap-storage"
	"squzy/internal/storage"
)

type HttpExecutor func(schedulerId string, timeout int32, config *scheduler_config_storage.HttpConfig, httpTool httpTools.HttpTool) job.CheckError
type GrpcExecutor func(schedulerId string, timeout int32, config *scheduler_config_storage.GrpcConfig, opts... grpc.DialOption) job.CheckError
type TcpExecutor func(schedulerId string, timeout int32, config *scheduler_config_storage.TcpConfig) job.CheckError
type SiteMapExecutor func(schedulerId string, timeout int32, config *scheduler_config_storage.SiteMapConfig, siteMapStorage sitemap_storage.SiteMapStorage, httpTools httpTools.HttpTool, semaphoreFactoryFn func(n int) semaphore.Semaphore) job.CheckError
type HttpValueExecutor func(schedulerId string, timeout int32, config *scheduler_config_storage.HttpValueConfig, httpTool httpTools.HttpTool) job.CheckError

type executor struct {
	externalStorage    storage.Storage
	siteMapStorage     sitemap_storage.SiteMapStorage
	httpTool           httpTools.HttpTool
	semaphoreFactoryFn func(n int) semaphore.Semaphore
	configStorage      scheduler_config_storage.Storage
	execTcp            TcpExecutor
	execGrpc           GrpcExecutor
	execHttp           HttpExecutor
	execSiteMap        SiteMapExecutor
	execHttpValue      HttpValueExecutor
}

func (e *executor) Execute(schedulerId primitive.ObjectID) {
	config, err := e.configStorage.Get(context.Background(), schedulerId)
	if err != nil || config == nil {
		// @TODO log error
		return
	}
	id := schedulerId.Hex()
	switch config.Type {
	case apiPb.SchedulerType_Tcp:
		_ = e.externalStorage.Write(e.execTcp(id, config.Timeout, config.TcpConfig))
		// @TODO logger
	case apiPb.SchedulerType_Grpc:
		_ = e.externalStorage.Write(e.execGrpc(id, config.Timeout, config.GrpcConfig, grpc.WithInsecure()))
		// @TODO logger
	case apiPb.SchedulerType_Http:
		_ = e.externalStorage.Write(e.execHttp(id, config.Timeout, config.HttpConfig, e.httpTool))
		// @TODO logger
	case apiPb.SchedulerType_SiteMap:
		_ = e.externalStorage.Write(e.execSiteMap(id, config.Timeout, config.SiteMapConfig, e.siteMapStorage, e.httpTool, e.semaphoreFactoryFn))
		// @TODO logger
	case apiPb.SchedulerType_HttpJsonValue:
		_ = e.externalStorage.Write(e.execHttpValue(id, config.Timeout, config.HttpValueConfig, e.httpTool))
		// @TODO logger
	default:
		// @TODO log incorrect type
	}
}

type JobExecutor interface {
	Execute(schedulerId primitive.ObjectID)
}

func NewExecutor(
	externalStorage storage.Storage,
	siteMapStorage sitemap_storage.SiteMapStorage,
	httpTool httpTools.HttpTool,
	semaphoreFactoryFn func(n int) semaphore.Semaphore,
	configStorage scheduler_config_storage.Storage,
	execTcp TcpExecutor,
	execGrpc GrpcExecutor,
	execHttp HttpExecutor,
	execSiteMap SiteMapExecutor,
	execHttpValue HttpValueExecutor,
) JobExecutor {
	return &executor{
		externalStorage:    externalStorage,
		siteMapStorage:     siteMapStorage,
		httpTool:           httpTool,
		semaphoreFactoryFn: semaphoreFactoryFn,
		configStorage:      configStorage,
		execTcp:            execTcp,
		execGrpc:           execGrpc,
		execHttp:           execHttp,
		execSiteMap:        execSiteMap,
		execHttpValue:      execHttpValue,
	}
}
