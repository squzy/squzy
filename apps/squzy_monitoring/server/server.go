package server

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"
	"squzy/internal/helpers"
	job_executor "squzy/internal/job-executor"
	"squzy/internal/scheduler"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
	scheduler_storage "squzy/internal/scheduler-storage"
)

var (
	errInvalidTypeError = errors.New("invalid type of config")
)

type server struct {
	schedulerStorage scheduler_storage.SchedulerStorage
	jobExecutor      job_executor.JobExecutor
	configStorage    scheduler_config_storage.Storage
}

func (s *server) GetSchedulerList(ctx context.Context, rq *empty.Empty) (*apiPb.GetSchedulerListResponse, error) {
	list, err := s.configStorage.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	arr := make([]*apiPb.Scheduler, len(list))
	var errGroup errgroup.Group
	for i := range list {
		index := i
		errGroup.Go(func() error {
			res, errG := s.GetSchedulerById(ctx, &apiPb.GetSchedulerByIdRequest{
				Id: list[index].ID.Hex(),
			})
			if errG != nil {
				return errG
			}
			arr[index] = res
			return nil
		})
	}
	err = errGroup.Wait()

	if err != nil {
		return nil, err
	}

	return &apiPb.GetSchedulerListResponse{
		Lists: arr,
	}, nil
}

func (s *server) GetSchedulerById(ctx context.Context, rq *apiPb.GetSchedulerByIdRequest) (*apiPb.Scheduler, error) {
	id := rq.Id
	idBson, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	config, err := s.configStorage.Get(ctx, idBson)
	if err != nil {
		return nil, err
	}
	switch config.Type {
	case apiPb.SchedulerType_TCP:
		return &apiPb.Scheduler{
			Id:       id,
			Name:     config.Name,
			Type:     apiPb.SchedulerType_TCP,
			Status:   config.Status,
			Interval: config.Interval,
			Timeout:  config.Timeout,
			Config: &apiPb.Scheduler_Tcp{
				Tcp: &apiPb.TcpConfig{
					Host: config.TCPConfig.Host,
					Port: config.TCPConfig.Port,
				},
			},
		}, nil
	case apiPb.SchedulerType_GRPC:
		return &apiPb.Scheduler{
			Id:       id,
			Name:     config.Name,
			Type:     apiPb.SchedulerType_GRPC,
			Status:   config.Status,
			Interval: config.Interval,
			Timeout:  config.Timeout,
			Config: &apiPb.Scheduler_Grpc{
				Grpc: &apiPb.GrpcConfig{
					Service: config.GrpcConfig.Service,
					Host:    config.GrpcConfig.Host,
					Port:    config.GrpcConfig.Port,
				},
			},
		}, nil
	case apiPb.SchedulerType_HTTP:
		return &apiPb.Scheduler{
			Id:       id,
			Name:     config.Name,
			Type:     apiPb.SchedulerType_HTTP,
			Status:   config.Status,
			Interval: config.Interval,
			Timeout:  config.Timeout,
			Config: &apiPb.Scheduler_Http{
				Http: &apiPb.HttpConfig{
					Method:     config.HTTPConfig.Method,
					Url:        config.HTTPConfig.URL,
					Headers:    config.HTTPConfig.Headers,
					StatusCode: config.HTTPConfig.StatusCode,
				},
			},
		}, nil
	case apiPb.SchedulerType_SITE_MAP:
		return &apiPb.Scheduler{
			Id:       id,
			Name:     config.Name,
			Type:     apiPb.SchedulerType_SITE_MAP,
			Status:   config.Status,
			Interval: config.Interval,
			Timeout:  config.Timeout,
			Config: &apiPb.Scheduler_Sitemap{
				Sitemap: &apiPb.SiteMapConfig{
					Url:         config.SiteMapConfig.URL,
					Concurrency: config.SiteMapConfig.Concurrency,
				},
			},
		}, nil
	case apiPb.SchedulerType_SSL_EXPIRATION:
		return &apiPb.Scheduler{
			Id:       id,
			Name:     config.Name,
			Type:     apiPb.SchedulerType_SSL_EXPIRATION,
			Status:   config.Status,
			Interval: config.Interval,
			Timeout:  config.Timeout,
			Config: &apiPb.Scheduler_SslExpiration{
				SslExpiration: &apiPb.SslExpirationConfig{
					Host: config.SslExpirationConfig.Host,
					Port: config.SslExpirationConfig.Port,
				},
			},
		}, nil
	case apiPb.SchedulerType_HTTP_JSON_VALUE:
		return &apiPb.Scheduler{
			Id:       id,
			Name:     config.Name,
			Type:     apiPb.SchedulerType_HTTP_JSON_VALUE,
			Status:   config.Status,
			Interval: config.Interval,
			Timeout:  config.Timeout,
			Config: &apiPb.Scheduler_HttpValue{
				HttpValue: &apiPb.HttpJsonValueConfig{
					Method:    config.HTTPValueConfig.Method,
					Url:       config.HTTPValueConfig.URL,
					Headers:   config.HTTPValueConfig.Headers,
					Selectors: helpers.SelectorsToProto(config.HTTPValueConfig.Selectors),
				},
			},
		}, nil
	default:
		return nil, errInvalidTypeError
	}
}

func (s *server) Remove(ctx context.Context, rq *apiPb.RemoveRequest) (*apiPb.RemoveResponse, error) {
	id := rq.Id
	idBson, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = s.configStorage.Remove(ctx, idBson)
	if err != nil {
		return nil, err
	}
	err = s.schedulerStorage.Remove(id)
	if err != nil {
		return nil, err
	}
	return &apiPb.RemoveResponse{
		Id: id,
	}, nil
}

func (s *server) Run(ctx context.Context, rq *apiPb.RunRequest) (*apiPb.RunResponse, error) {
	id := rq.Id
	idBson, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = s.configStorage.Run(ctx, idBson)
	if err != nil {
		return nil, err
	}
	schld, err := s.schedulerStorage.Get(id)
	if err != nil {
		return nil, err
	}
	schld.Run()
	return &apiPb.RunResponse{
		Id: id,
	}, nil
}

func (s *server) Stop(ctx context.Context, rq *apiPb.StopRequest) (*apiPb.StopResponse, error) {
	id := rq.Id
	idBson, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = s.configStorage.Stop(ctx, idBson)
	if err != nil {
		return nil, err
	}
	schld, err := s.schedulerStorage.Get(id)
	if err != nil {
		return nil, err
	}
	schld.Stop()
	return &apiPb.StopResponse{
		Id: id,
	}, nil
}

func (s *server) Add(ctx context.Context, rq *apiPb.AddRequest) (*apiPb.AddResponse, error) {
	interval := helpers.DurationFromSecond(rq.Interval)
	schld, err := scheduler.New(
		primitive.NewObjectID(),
		interval,
		s.jobExecutor,
	)
	if err != nil {
		return nil, err
	}
	var schedulerConfig *scheduler_config_storage.SchedulerConfig
	switch config := rq.Config.(type) {
	case *apiPb.AddRequest_Tcp:
		schedulerConfig = &scheduler_config_storage.SchedulerConfig{
			ID:       schld.GetIDBson(),
			Name:     rq.Name,
			Type:     apiPb.SchedulerType_TCP,
			Status:   apiPb.SchedulerStatus_STOPPED,
			Interval: rq.Interval,
			Timeout:  rq.Timeout,
			TCPConfig: &scheduler_config_storage.TCPConfig{
				Host: config.Tcp.Host,
				Port: config.Tcp.Port,
			},
		}
	case *apiPb.AddRequest_Sitemap:
		schedulerConfig = &scheduler_config_storage.SchedulerConfig{
			ID:       schld.GetIDBson(),
			Name:     rq.Name,
			Type:     apiPb.SchedulerType_SITE_MAP,
			Status:   apiPb.SchedulerStatus_STOPPED,
			Interval: rq.Interval,
			Timeout:  rq.Timeout,
			SiteMapConfig: &scheduler_config_storage.SiteMapConfig{
				URL:         config.Sitemap.Url,
				Concurrency: config.Sitemap.Concurrency,
			},
		}
	case *apiPb.AddRequest_Grpc:
		schedulerConfig = &scheduler_config_storage.SchedulerConfig{
			ID:       schld.GetIDBson(),
			Name:     rq.Name,
			Type:     apiPb.SchedulerType_GRPC,
			Status:   apiPb.SchedulerStatus_STOPPED,
			Interval: rq.Interval,
			Timeout:  rq.Timeout,
			GrpcConfig: &scheduler_config_storage.GrpcConfig{
				Service: config.Grpc.Service,
				Host:    config.Grpc.Host,
				Port:    config.Grpc.Port,
			},
		}
	case *apiPb.AddRequest_Http:
		schedulerConfig = &scheduler_config_storage.SchedulerConfig{
			ID:       schld.GetIDBson(),
			Name:     rq.Name,
			Type:     apiPb.SchedulerType_HTTP,
			Status:   apiPb.SchedulerStatus_STOPPED,
			Interval: rq.Interval,
			Timeout:  rq.Timeout,
			HTTPConfig: &scheduler_config_storage.HTTPConfig{
				Method:     config.Http.Method,
				URL:        config.Http.Url,
				Headers:    config.Http.Headers,
				StatusCode: config.Http.StatusCode,
			},
		}
	case *apiPb.AddRequest_HttpValue:
		schedulerConfig = &scheduler_config_storage.SchedulerConfig{
			ID:       schld.GetIDBson(),
			Name:     rq.Name,
			Type:     apiPb.SchedulerType_HTTP_JSON_VALUE,
			Status:   apiPb.SchedulerStatus_STOPPED,
			Interval: rq.Interval,
			Timeout:  rq.Timeout,
			HTTPValueConfig: &scheduler_config_storage.HTTPValueConfig{
				Method:    config.HttpValue.Method,
				URL:       config.HttpValue.Url,
				Headers:   config.HttpValue.Headers,
				Selectors: helpers.SelectorsToDb(config.HttpValue.Selectors),
			},
		}
	case *apiPb.AddRequest_SslExpiration:
		schedulerConfig = &scheduler_config_storage.SchedulerConfig{
			ID:       schld.GetIDBson(),
			Name:     rq.Name,
			Type:     apiPb.SchedulerType_SSL_EXPIRATION,
			Status:   apiPb.SchedulerStatus_STOPPED,
			Interval: rq.Interval,
			Timeout:  rq.Timeout,
			SslExpirationConfig: &scheduler_config_storage.SslExpirationConfig{
				Host: config.SslExpiration.Host,
				Port: config.SslExpiration.Port,
			},
		}

	default:
		return nil, errInvalidTypeError
	}
	err = s.configStorage.Add(ctx, schedulerConfig)
	if err != nil {
		return nil, err
	}
	err = s.schedulerStorage.Set(schld)
	if err != nil {
		return nil, err
	}
	return &apiPb.AddResponse{
		Id: schld.GetID(),
	}, nil
}

func New(
	schedulerStorage scheduler_storage.SchedulerStorage,
	jobExecutor job_executor.JobExecutor,
	configStorage scheduler_config_storage.Storage,
) apiPb.SchedulersExecutorServer {
	return &server{
		schedulerStorage: schedulerStorage,
		jobExecutor:      jobExecutor,
		configStorage:    configStorage,
	}
}
