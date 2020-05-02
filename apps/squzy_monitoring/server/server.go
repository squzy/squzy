package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc"
	"squzy/internal/helpers"
	"squzy/internal/httpTools"
	"squzy/internal/job"
	"squzy/internal/scheduler"
	scheduler_storage "squzy/internal/scheduler-storage"
	"squzy/internal/semaphore"
	sitemap_storage "squzy/internal/sitemap-storage"
	"squzy/internal/storage"
)

type server struct {
	semaphoreFactory semaphore.SemaphoreFactory
	schedulerStorage scheduler_storage.SchedulerStorage
	externalStorage  storage.Storage
	siteMapStorage   sitemap_storage.SiteMapStorage
	httpTools        httpTools.HttpTool
}

func (s *server) GetSchedulerList(context.Context, *empty.Empty) (*apiPb.GetSchedulerListResponse, error) {
	return &apiPb.GetSchedulerListResponse{}, nil
}

func (s *server) GetSchedulerById(context.Context, *apiPb.GetSchedulerByIdRequest) (*apiPb.Scheduler, error) {
	return &apiPb.Scheduler{}, nil
}

func (s *server) Remove(ctx context.Context, rq *apiPb.RemoveRequest) (*apiPb.RemoveResponse, error) {
	err := s.schedulerStorage.Remove(rq.Id)
	if err != nil {
		return nil, err
	}
	return &apiPb.RemoveResponse{
		Id: rq.Id,
	}, nil
}

func (s *server) Run(ctx context.Context, rq *apiPb.RunRequest) (*apiPb.RunResponse, error) {
	schld, err := s.schedulerStorage.Get(rq.Id)
	if err != nil {
		return nil, err
	}
	err = schld.Run()
	if err != nil {
		return nil, err
	}
	return &apiPb.RunResponse{
		Id: schld.GetId(),
	}, nil
}

func (s *server) Stop(ctx context.Context, rq *apiPb.StopRequest) (*apiPb.StopResponse, error) {
	schld, err := s.schedulerStorage.Get(rq.Id)
	if err != nil {
		return nil, err
	}
	err = schld.Stop()
	if err != nil {
		return nil, err
	}
	return &apiPb.StopResponse{
		Id: schld.GetId(),
	}, nil
}

func (s *server) Add(ctx context.Context, rq *apiPb.AddRequest) (*apiPb.AddResponse, error) {
	interval := helpers.DurationFromSecond(rq.Interval)
	switch check := rq.Config.(type) {
	case *apiPb.AddRequest_Tcp:
		tcpCheck := check.Tcp
		schld, err := scheduler.New(
			interval,
			job.NewTcpJob(tcpCheck.Host, tcpCheck.Port, rq.Timeout),
			s.externalStorage,
		)
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
	case *apiPb.AddRequest_Sitemap:
		siteMapCheck := check.Sitemap
		schld, err := scheduler.New(
			interval,
			job.NewSiteMapJob(siteMapCheck.Url, rq.Timeout, s.siteMapStorage, s.httpTools, s.semaphoreFactory, siteMapCheck.Concurrency),
			s.externalStorage,
		)
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
	case *apiPb.AddRequest_Grpc:
		grcpCheck := check.Grpc
		schld, err := scheduler.New(
			interval,
			job.NewGrpcJob(grcpCheck.Service, grcpCheck.Host, grcpCheck.Port, rq.Timeout, []grpc.DialOption{grpc.WithInsecure()}, []grpc.CallOption{}),
			s.externalStorage,
		)
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
	case *apiPb.AddRequest_Http:
		httpCheck := check.Http
		schld, err := scheduler.New(
			interval,
			job.NewHttpJob(httpCheck.Method, httpCheck.Url, httpCheck.Headers, rq.Timeout, httpCheck.StatusCode, s.httpTools),
			s.externalStorage,
		)
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

	case *apiPb.AddRequest_HttpValue:
		httpJsonValueCheck := check.HttpValue
		schld, err := scheduler.New(
			interval,
			job.NewJsonHttpValueJob(httpJsonValueCheck.Method, httpJsonValueCheck.Url, httpJsonValueCheck.Headers, rq.Timeout, s.httpTools, httpJsonValueCheck.Selectors),
			s.externalStorage,
		)
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
	default:
		return &apiPb.AddResponse{
			Id: "",
		}, nil
	}
}

func New(
	schedulerStorage scheduler_storage.SchedulerStorage,
	externalStorage storage.Storage,
	siteMapStorage sitemap_storage.SiteMapStorage,
	httpTools httpTools.HttpTool,
	semaphoreFactory semaphore.SemaphoreFactory,
) apiPb.SchedulersExecutorServer {
	return &server{
		schedulerStorage: schedulerStorage,
		externalStorage:  externalStorage,
		siteMapStorage:   siteMapStorage,
		httpTools:        httpTools,
		semaphoreFactory: semaphoreFactory,
	}
}
