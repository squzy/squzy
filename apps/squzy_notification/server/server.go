package server

import (
	"context"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

type server struct {

}

func (s *server) CreateNotificationMethod(ctx context.Context, request *apiPb.CreateNotificationMethodRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (s *server) GetById(ctx context.Context, request *apiPb.NotificationMethodIdRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (s *server) DeleteById(ctx context.Context, request *apiPb.NotificationMethodIdRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (s *server) Activate(ctx context.Context, request *apiPb.NotificationMethodIdRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (s *server) Deactivate(ctx context.Context, request *apiPb.NotificationMethodIdRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (s *server) Add(ctx context.Context, request *apiPb.NotificationMethodRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (s *server) Remove(ctx context.Context, request *apiPb.NotificationMethodRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (s *server) GetList(ctx context.Context, request *apiPb.GetListRequest) (*apiPb.GetListResponse, error) {
	panic("implement me")
}

func New() apiPb.NotificationManagerServer {
	return &server{}
}