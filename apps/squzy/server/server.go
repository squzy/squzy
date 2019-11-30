package server

import (
	clientPb "github.com/squzy/squzy_generated/generated/server/proto/v1"
	"context"
	scheduler_storage "squzy/apps/internal/scheduler-storage"
)

type server struct {
	schedulerStorage scheduler_storage.SchedulerStorage
}

func (s server) GetList(context.Context, *clientPb.GetListRequest) (*clientPb.GetListResponse, error) {
	panic("implement me")
}
