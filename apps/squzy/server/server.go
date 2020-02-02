package server

import (
	"context"
	"database/sql"
	serverPb "github.com/squzy/squzy_generated/generated/server/proto/v1"
	"google.golang.org/grpc"
	"squzy/apps/internal/httpTools"
	"squzy/apps/internal/job"
	"squzy/apps/internal/scheduler"
	scheduler_storage "squzy/apps/internal/scheduler-storage"
	"squzy/apps/internal/semaphore"
	sitemap_storage "squzy/apps/internal/sitemap-storage"
	"squzy/apps/internal/storage"
	"time"
)

type server struct {
	schedulerStorage scheduler_storage.SchedulerStorage
	externalStorage  storage.Storage
	siteMapStorage   sitemap_storage.SiteMapStorage
	httpTools        httpTools.HttpTool
	mySqlPing        func(db *sql.DB) error
	semaphoreFactory semaphore.SemaphoreFactory
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
	interval := rq.Interval
	switch check := rq.Check.(type) {
	case *serverPb.AddSchedulerRequest_TcpCheck:
		tcpCheck := check.TcpCheck
		schld, err := scheduler.New(
			time.Second*time.Duration(interval),
			job.NewTcpJob(tcpCheck.Host, tcpCheck.Port),
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
			time.Second*time.Duration(interval),
			job.NewSiteMapJob(siteMapCheck.Url, s.siteMapStorage, s.httpTools, s.semaphoreFactory, siteMapCheck.Concurrency),
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
			time.Second*time.Duration(interval),
			job.NewGrpcJob(grcpCheck.Service, grcpCheck.Host, grcpCheck.Port, []grpc.DialOption{grpc.WithInsecure()}, []grpc.CallOption{}),
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
			time.Second*time.Duration(interval),
			job.NewHttpJob(httpCheck.Method, httpCheck.Url, httpCheck.Headers, httpCheck.StatusCode, s.httpTools),
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
	case *serverPb.AddSchedulerRequest_MongoCheck:
		mongoCheck := check.MongoCheck
		schld, err := scheduler.New(
			time.Second*time.Duration(interval),
			job.NewMongoJob(mongoCheck.Url),
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
	case *serverPb.AddSchedulerRequest_PostgresCheck:
		postgresCheck := check.PostgresCheck
		schld, err := scheduler.New(
			time.Second*time.Duration(interval),
			job.NewPosgresDbJob(
				postgresCheck.Host,
				postgresCheck.Port,
				postgresCheck.User,
				postgresCheck.Password,
				postgresCheck.DbName,
				s.mySqlPing),
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
	case *serverPb.AddSchedulerRequest_CassandraCheck:
		cassandraCheck := check.CassandraCheck
		schld, err := scheduler.New(
			time.Second*time.Duration(interval),
			job.NewCassandraJob(
				cassandraCheck.Cluster,
				cassandraCheck.User,
				cassandraCheck.Password),
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
	case *serverPb.AddSchedulerRequest_MysqlCheck:
		mysqlCheck := check.MysqlCheck
		schld, err := scheduler.New(
			time.Second*time.Duration(interval),
			job.NewMysqlJob(
				mysqlCheck.Host,
				mysqlCheck.Port,
				mysqlCheck.User,
				mysqlCheck.Password,
				mysqlCheck.DbName,
				s.mySqlPing),
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
	mySqlPing func(db *sql.DB) error,
	semaphoreFactory semaphore.SemaphoreFactory,
) serverPb.ServerServer {
	return &server{
		schedulerStorage: schedulerStorage,
		externalStorage:  externalStorage,
		siteMapStorage:   siteMapStorage,
		httpTools:        httpTools,
		mySqlPing:        mySqlPing,
		semaphoreFactory: semaphoreFactory,
	}
}
