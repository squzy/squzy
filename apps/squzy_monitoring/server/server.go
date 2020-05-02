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
	invalidTypeError = errors.New("Invalid type of config")
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
	for i, _ := range list {
		index := i
		errGroup.Go(func() error {
			res, err := s.GetSchedulerById(ctx, &apiPb.GetSchedulerByIdRequest{
				Id: list[index].Id.Hex(),
			})
			if err != nil {
				return err
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
		List: arr,
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
	case apiPb.SchedulerType_Tcp:
		return &apiPb.Scheduler{
			Id:       id,
			Type:     config.Type,
			Status:   config.Status,
			Interval: config.Interval,
			Timeout:  config.Timeout,
			Config: &apiPb.Scheduler_Tcp{
				Tcp: &apiPb.TcpConfig{
					Host: config.TcpConfig.Host,
					Port: config.TcpConfig.Port,
				},
			},
		}, nil
	case apiPb.SchedulerType_Grpc:
		return &apiPb.Scheduler{
			Id:       id,
			Type:     config.Type,
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
	case apiPb.SchedulerType_Http:
		return &apiPb.Scheduler{
			Id:       id,
			Type:     config.Type,
			Status:   config.Status,
			Interval: config.Interval,
			Timeout:  config.Timeout,
			Config: &apiPb.Scheduler_Http{
				Http: &apiPb.HttpConfig{
					Method:     config.HttpConfig.Method,
					Url:        config.HttpConfig.Url,
					Headers:    config.HttpConfig.Headers,
					StatusCode: config.HttpConfig.StatusCode,
				},
			},
		}, nil
	case apiPb.SchedulerType_SiteMap:
		return &apiPb.Scheduler{
			Id:       id,
			Type:     config.Type,
			Status:   config.Status,
			Interval: config.Interval,
			Timeout:  config.Timeout,
			Config: &apiPb.Scheduler_Sitemap{
				Sitemap: &apiPb.SiteMapConfig{
					Url:         config.SiteMapConfig.Url,
					Concurrency: config.SiteMapConfig.Concurrency,
				},
			},
		}, nil
	case apiPb.SchedulerType_HttpJsonValue:
		return &apiPb.Scheduler{
			Id:       id,
			Type:     config.Type,
			Status:   config.Status,
			Interval: config.Interval,
			Timeout:  config.Timeout,
			Config: &apiPb.Scheduler_HttpValue{
				HttpValue: &apiPb.HttpJsonValueConfig{
					Method:    config.HttpValueConfig.Method,
					Url:       config.HttpValueConfig.Url,
					Headers:   config.HttpValueConfig.Headers,
					Selectors: helpers.SelectorsToProto(config.HttpValueConfig.Selectors),
				},
			},
		}, nil
	default:
		return nil, invalidTypeError
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
			Id:       schld.GetIdBson(),
			Type:     apiPb.SchedulerType_Tcp,
			Status:   apiPb.SchedulerStatus_STOPPED,
			Interval: rq.Interval,
			Timeout:  rq.Timeout,
			TcpConfig: &scheduler_config_storage.TcpConfig{
				Host: config.Tcp.Host,
				Port: config.Tcp.Port,
			},
		}
	case *apiPb.AddRequest_Sitemap:
		schedulerConfig = &scheduler_config_storage.SchedulerConfig{
			Id:       schld.GetIdBson(),
			Type:     apiPb.SchedulerType_SiteMap,
			Status:   apiPb.SchedulerStatus_STOPPED,
			Interval: rq.Interval,
			Timeout:  rq.Timeout,
			SiteMapConfig: &scheduler_config_storage.SiteMapConfig{
				Url:         config.Sitemap.Url,
				Concurrency: config.Sitemap.Concurrency,
			},
		}
	case *apiPb.AddRequest_Grpc:
		schedulerConfig = &scheduler_config_storage.SchedulerConfig{
			Id:       schld.GetIdBson(),
			Type:     apiPb.SchedulerType_Grpc,
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
			Id:       schld.GetIdBson(),
			Type:     apiPb.SchedulerType_Http,
			Status:   apiPb.SchedulerStatus_STOPPED,
			Interval: rq.Interval,
			Timeout:  rq.Timeout,
			HttpConfig: &scheduler_config_storage.HttpConfig{
				Method:     config.Http.Method,
				Url:        config.Http.Url,
				Headers:    config.Http.Headers,
				StatusCode: config.Http.StatusCode,
			},
		}
	case *apiPb.AddRequest_HttpValue:
		schedulerConfig = &scheduler_config_storage.SchedulerConfig{
			Id:       schld.GetIdBson(),
			Type:     apiPb.SchedulerType_HttpJsonValue,
			Status:   apiPb.SchedulerStatus_STOPPED,
			Interval: rq.Interval,
			Timeout:  rq.Timeout,
			HttpValueConfig: &scheduler_config_storage.HttpValueConfig{
				Method:    config.HttpValue.Method,
				Url:       config.HttpValue.Url,
				Headers:   config.HttpValue.Headers,
				Selectors: helpers.SelectorsToDb(config.HttpValue.Selectors),
			},
		}
	default:
		return nil, invalidTypeError
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
		Id: schld.GetId(),
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
