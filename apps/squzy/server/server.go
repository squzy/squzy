package server

import (
	clientPb "github.com/squzy/squzy_generated/generated/squzy"
	"context"
	scheduler_storage "squzy/apps/internal/scheduler-storage"
)

type server struct {
	schedulerStorage scheduler_storage.SchedulerStorage
}

func (s server) GetList(context.Context, *clientPb.GetListRequest) (*clientPb.GetListResponse, error) {
	panic("implement me")
}

func New(schedulerStorage scheduler_storage.SchedulerStorage) clientPb.ServerServer {
	return &server{
		schedulerStorage: schedulerStorage,
	}
}