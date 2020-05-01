package server

import (
	"context"
	serverPb "github.com/squzy/squzy_generated/generated/server/proto/v1"
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

func (s server) RemoveScheduler(ctx context.Context, rq *serverPb.RemoveSchedulerRequest) (*serverPb.RemoveSchedulerResponse, error) {
	err := s.schedulerStorage.Remove(rq.Id)
	if err != nil {
		return nil, err
	}
	return &serverPb.RemoveSchedulerResponse{
		Id: rq.Id,
	}, nil
}

func (s server) RunScheduler(ctx context.Context, rq *serverPb.RunSchedulerRequest) (*serverPb.RunSchedulerResponse, error) {
	schld, err := s.schedulerStorage.Get(rq.Id)
	if err != nil {
		return nil, err
	}
	err = schld.Run()
	if err != nil {
		return nil, err
	}
	return &serverPb.RunSchedulerResponse{
		Id: schld.GetId(),
	}, nil
}

func (s server) StopScheduler(ctx context.Context, rq *serverPb.StopSchedulerRequest) (*serverPb.StopSchedulerResponse, error) {
	schld, err := s.schedulerStorage.Get(rq.Id)
	if err != nil {
		return nil, err
	}
	err = schld.Stop()
	if err != nil {
		return nil, err
	}
	return &serverPb.StopSchedulerResponse{
		Id: schld.GetId(),
	}, nil
}

func (s server) AddScheduler(ctx context.Context, rq *serverPb.AddSchedulerRequest) (*serverPb.AddSchedulerResponse, error) {
	interval := helpers.DurationFromSecond(rq.Interval)
	switch check := rq.Check.(type) {
	case *serverPb.AddSchedulerRequest_TcpCheck:
		tcpCheck := check.TcpCheck
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
		return &serverPb.AddSchedulerResponse{
			Id: schld.GetId(),
		}, nil
	case *serverPb.AddSchedulerRequest_SitemapCheck:
		siteMapCheck := check.SitemapCheck
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
		return &serverPb.AddSchedulerResponse{
			Id: schld.GetId(),
		}, nil
	case *serverPb.AddSchedulerRequest_GrpcCheck:
		grcpCheck := check.GrpcCheck
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
		return &serverPb.AddSchedulerResponse{
			Id: schld.GetId(),
		}, nil
	case *serverPb.AddSchedulerRequest_HttpCheck:
		httpCheck := check.HttpCheck
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
		return &serverPb.AddSchedulerResponse{
			Id: schld.GetId(),
		}, nil

	case *serverPb.AddSchedulerRequest_HttpJsonValue:
		httpJsonValueCheck := check.HttpJsonValue
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
		return &serverPb.AddSchedulerResponse{
			Id: schld.GetId(),
		}, nil
	default:
		return &serverPb.AddSchedulerResponse{
			Id: "",
		}, nil
	}
}

func (s server) GetList(context.Context, *serverPb.GetListRequest) (*serverPb.GetListResponse, error) {
	sMap := s.schedulerStorage.GetList()
	list := []*serverPb.SchedulerListItem{}
	for k, v := range sMap {
		status := serverPb.Status_STOPPED
		if v {
			status = serverPb.Status_RUNNED
		}
		list = append(list, &serverPb.SchedulerListItem{
			Id:     k,
			Status: status,
		})
	}
	return &serverPb.GetListResponse{
		List: list,
	}, nil
}

func New(
	schedulerStorage scheduler_storage.SchedulerStorage,
	externalStorage storage.Storage,
	siteMapStorage sitemap_storage.SiteMapStorage,
	httpTools httpTools.HttpTool,
	semaphoreFactory semaphore.SemaphoreFactory,
) serverPb.ServerServer {
	return &server{
		schedulerStorage: schedulerStorage,
		externalStorage:  externalStorage,
		siteMapStorage:   siteMapStorage,
		httpTools:        httpTools,
		semaphoreFactory: semaphoreFactory,
	}
}
