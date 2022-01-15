package job_executor

import (
	"context"
	"crypto/tls"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"github.com/squzy/squzy/internal/httptools"
	"github.com/squzy/squzy/internal/job"
	"github.com/squzy/squzy/internal/logger"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	"github.com/squzy/squzy/internal/semaphore"
	sitemap_storage "github.com/squzy/squzy/internal/sitemap-storage"
	"github.com/squzy/squzy/internal/storage"
)

type HTTPExecutor func(schedulerId string,
	timeout int32,
	config *scheduler_config_storage.HTTPConfig,
	httpTool httptools.HTTPTool) job.CheckError

type SSLExpirationExecutor func(
	schedulerId string,
	timeout int32,
	config *scheduler_config_storage.SslExpirationConfig,
	cfg *tls.Config,
) job.CheckError

type GrpcExecutor func(schedulerId string,
	timeout int32,
	config *scheduler_config_storage.GrpcConfig,
	opts ...grpc.DialOption) job.CheckError
type TCPExecutor func(schedulerId string, timeout int32, config *scheduler_config_storage.TCPConfig) job.CheckError

type SiteMapExecutor func(
	schedulerId string,
	timeout int32,
	config *scheduler_config_storage.SiteMapConfig,
	siteMapStorage sitemap_storage.SiteMapStorage,
	httpTools httptools.HTTPTool,
	semaphoreFactoryFn func(n int) semaphore.Semaphore) job.CheckError

type HTTPValueExecutor func(
	schedulerId string,
	timeout int32,
	config *scheduler_config_storage.HTTPValueConfig,
	httpTool httptools.HTTPTool) job.CheckError

type executor struct {
	externalStorage    storage.Storage
	siteMapStorage     sitemap_storage.SiteMapStorage
	httpTool           httptools.HTTPTool
	semaphoreFactoryFn func(n int) semaphore.Semaphore
	configStorage      scheduler_config_storage.Storage
	execTCP            TCPExecutor
	execGrpc           GrpcExecutor
	execHTTP           HTTPExecutor
	execSiteMap        SiteMapExecutor
	execHTTPValue      HTTPValueExecutor
	execSSLExpiration  SSLExpirationExecutor
}

func (e *executor) Execute(schedulerID primitive.ObjectID) {
	config, err := e.configStorage.Get(context.Background(), schedulerID)
	if err != nil || config == nil {
		msg := schedulerID.Hex()
		if err != nil {
			msg += err.Error()
		}
		logger.Errorf("Could not get config for schedulerID: %s", msg)
		return
	}
	id := schedulerID.Hex()
	switch config.Type {
	case apiPb.SchedulerType_TCP:
		_ = e.externalStorage.Write(e.execTCP(id, config.Timeout, config.TCPConfig))
		logger.Infof("TCP job executed is used for scheduler id %s", schedulerID)
	case apiPb.SchedulerType_GRPC:
		_ = e.externalStorage.Write(e.execGrpc(id, config.Timeout, config.GrpcConfig, grpc.WithInsecure()))
		logger.Infof("gRPC job executed is used for scheduler id %s", schedulerID)
	case apiPb.SchedulerType_HTTP:
		_ = e.externalStorage.Write(e.execHTTP(id, config.Timeout, config.HTTPConfig, e.httpTool))
		logger.Infof("HTTP job executed is used for scheduler id %s", schedulerID)
	case apiPb.SchedulerType_SITE_MAP:
		_ = e.externalStorage.Write(e.execSiteMap(id, config.Timeout, config.SiteMapConfig, e.siteMapStorage, e.httpTool, e.semaphoreFactoryFn))
		logger.Infof("Site map job executed is used for scheduler id %s", schedulerID)
	case apiPb.SchedulerType_HTTP_JSON_VALUE:
		_ = e.externalStorage.Write(e.execHTTPValue(id, config.Timeout, config.HTTPValueConfig, e.httpTool))
		logger.Infof("HTTP JSON job executed is used for scheduler id %s", schedulerID)
	case apiPb.SchedulerType_SSL_EXPIRATION:
		_ = e.externalStorage.Write(e.execSSLExpiration(id, config.Timeout, config.SslExpirationConfig, nil))
		logger.Infof("SSL Expiration job executed is used for scheduler id %s", schedulerID)
	default:
		logger.Errorf("Incorrect config type passed to job executor: %s", config.Type)
	}
}

type JobExecutor interface {
	Execute(schedulerID primitive.ObjectID)
}

func NewExecutor(
	externalStorage storage.Storage,
	siteMapStorage sitemap_storage.SiteMapStorage,
	httpTool httptools.HTTPTool,
	semaphoreFactoryFn func(n int) semaphore.Semaphore,
	configStorage scheduler_config_storage.Storage,
	execTCP TCPExecutor,
	execGrpc GrpcExecutor,
	execHTTP HTTPExecutor,
	execSiteMap SiteMapExecutor,
	execHTTPValue HTTPValueExecutor,
	execSSLExpiration SSLExpirationExecutor,
) JobExecutor {
	return &executor{
		externalStorage:    externalStorage,
		siteMapStorage:     siteMapStorage,
		httpTool:           httpTool,
		semaphoreFactoryFn: semaphoreFactoryFn,
		configStorage:      configStorage,
		execTCP:            execTCP,
		execGrpc:           execGrpc,
		execHTTP:           execHTTP,
		execSiteMap:        execSiteMap,
		execHTTPValue:      execHTTPValue,
		execSSLExpiration:  execSSLExpiration,
	}
}
