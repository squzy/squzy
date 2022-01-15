package server

import (
	"context"
	"errors"
	"github.com/squzy/squzy/apps/squzy_notification/database"
	"github.com/squzy/squzy/apps/squzy_notification/integrations"
	"github.com/squzy/squzy/internal/helpers"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type server struct {
	nlDb         database.NotificationListDb
	nmDb         database.NotificationMethodDb
	client       apiPb.StorageClient
	integrations integrations.Integrations
}

var (
	errTypeNotExist = errors.New("notification type not exist")
)

func dbMethodToProto(method *database.NotificationMethod) (*apiPb.NotificationMethod, error) {
	switch method.Type {
	case apiPb.NotificationMethodType_NOTIFICATION_METHOD_WEBHOOK:
		return &apiPb.NotificationMethod{
			Id:     method.Id.Hex(),
			Status: method.Status,
			Name:   method.Name,
			Type:   apiPb.NotificationMethodType_NOTIFICATION_METHOD_WEBHOOK,
			Method: &apiPb.NotificationMethod_Webhook{
				Webhook: &apiPb.WebHookMethod{
					Url: method.WebHook.Url,
				},
			},
		}, nil
	case apiPb.NotificationMethodType_NOTIFICATION_METHOD_SLACK:
		return &apiPb.NotificationMethod{
			Id:     method.Id.Hex(),
			Status: method.Status,
			Name:   method.Name,
			Type:   apiPb.NotificationMethodType_NOTIFICATION_METHOD_SLACK,
			Method: &apiPb.NotificationMethod_Slack{
				Slack: &apiPb.SlackMethod{
					Url: method.Slack.Url,
				},
			},
		}, nil
	default:
		return nil, errTypeNotExist
	}
}

func (s *server) GetNotificationMethods(ctx context.Context, e *empty.Empty) (*apiPb.GetListResponse, error) {
	list, err := s.nmDb.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	arr := []*apiPb.NotificationMethod{}
	for _, value := range list {
		item, err := dbMethodToProto(value)
		if err != nil {
			return nil, err
		}
		arr = append(arr, item)
	}
	return &apiPb.GetListResponse{
		Methods: arr,
	}, nil
}

func (s *server) Notify(ctx context.Context, request *apiPb.NotifyRequest) (*empty.Empty, error) {
	ownerId, err := primitive.ObjectIDFromHex(request.OwnerId)
	if err != nil {
		return nil, err
	}
	incident, err := s.client.GetIncidentById(ctx, &apiPb.IncidentIdRequest{
		IncidentId: request.IncidentId,
	})
	if err != nil {
		return nil, err
	}
	methods, err := s.nlDb.GetList(ctx, ownerId, request.OwnerType)
	if err != nil {
		return nil, err
	}

	var eg errgroup.Group
	for _, method := range methods {
		m := method
		eg.Go(func() error {
			c, cancel := helpers.TimeoutContext(ctx, time.Second*5)
			defer cancel()
			config, errGet := s.nmDb.Get(c, m.NotificationMethodId)
			if errGet != nil {
				logger.Error(errGet.Error())
				return nil
			}
			if config.Status != apiPb.NotificationMethodStatus_NOTIFICATION_STATUS_ACTIVE {
				return nil
			}
			switch config.Type {
			case apiPb.NotificationMethodType_NOTIFICATION_METHOD_SLACK:
				s.integrations.Slack(ctx, incident, config.Slack)
				return nil
			case apiPb.NotificationMethodType_NOTIFICATION_METHOD_WEBHOOK:
				s.integrations.Webhook(ctx, incident, config.WebHook)
				return nil
			}
			return nil
		})
	}
	_ = eg.Wait()
	return &empty.Empty{}, nil
}

func (s *server) CreateNotificationMethod(ctx context.Context, request *apiPb.CreateNotificationMethodRequest) (*apiPb.NotificationMethod, error) {
	var req *database.NotificationMethod
	switch request.Type {
	case apiPb.NotificationMethodType_NOTIFICATION_METHOD_SLACK:
		req = &database.NotificationMethod{
			Id:     primitive.NewObjectID(),
			Status: apiPb.NotificationMethodStatus_NOTIFICATION_STATUS_ACTIVE,
			Type:   request.Type,
			Name:   request.Name,
			Slack: &database.SlackConfig{
				Url: request.GetSlack().Url,
			},
		}
	case apiPb.NotificationMethodType_NOTIFICATION_METHOD_WEBHOOK:
		req = &database.NotificationMethod{
			Id:     primitive.NewObjectID(),
			Status: apiPb.NotificationMethodStatus_NOTIFICATION_STATUS_ACTIVE,
			Type:   request.Type,
			Name:   request.Name,
			WebHook: &database.WebHookConfig{
				Url: request.GetWebhook().Url,
			},
		}
	default:
		return nil, errTypeNotExist
	}
	err := s.nmDb.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return dbMethodToProto(req)
}

func (s *server) GetById(ctx context.Context, request *apiPb.NotificationMethodIdRequest) (*apiPb.NotificationMethod, error) {
	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	method, err := s.nmDb.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return dbMethodToProto(method)
}

func (s *server) DeleteById(ctx context.Context, request *apiPb.NotificationMethodIdRequest) (*apiPb.NotificationMethod, error) {
	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	err = s.nmDb.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	method, err := s.nmDb.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return dbMethodToProto(method)
}

func (s *server) Activate(ctx context.Context, request *apiPb.NotificationMethodIdRequest) (*apiPb.NotificationMethod, error) {
	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	err = s.nmDb.Activate(ctx, id)
	if err != nil {
		return nil, err
	}
	method, err := s.nmDb.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return dbMethodToProto(method)
}

func (s *server) Deactivate(ctx context.Context, request *apiPb.NotificationMethodIdRequest) (*apiPb.NotificationMethod, error) {
	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	err = s.nmDb.Deactivate(ctx, id)
	if err != nil {
		return nil, err
	}
	method, err := s.nmDb.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return dbMethodToProto(method)
}

func (s *server) Add(ctx context.Context, request *apiPb.NotificationMethodRequest) (*apiPb.NotificationMethod, error) {
	ownerId, err := primitive.ObjectIDFromHex(request.OwnerId)
	if err != nil {
		return nil, err
	}
	methodId, err := primitive.ObjectIDFromHex(request.NotificationMethodId)
	if err != nil {
		return nil, err
	}
	err = s.nlDb.Add(ctx, &database.Notification{
		Id:                   primitive.NewObjectID(),
		OwnerId:              ownerId,
		Type:                 request.OwnerType,
		NotificationMethodId: methodId,
	})
	if err != nil {
		return nil, err
	}
	method, err := s.nmDb.Get(ctx, methodId)
	if err != nil {
		return nil, err
	}
	return dbMethodToProto(method)
}

func (s *server) Remove(ctx context.Context, request *apiPb.NotificationMethodRequest) (*apiPb.NotificationMethod, error) {
	methodId, err := primitive.ObjectIDFromHex(request.NotificationMethodId)
	if err != nil {
		return nil, err
	}
	err = s.nlDb.Delete(ctx, methodId)
	if err != nil {
		return nil, err
	}
	method, err := s.nmDb.Get(ctx, methodId)
	if err != nil {
		return nil, err
	}
	return dbMethodToProto(method)
}

func (s *server) GetList(ctx context.Context, request *apiPb.GetListRequest) (*apiPb.GetListResponse, error) {
	ownerId, err := primitive.ObjectIDFromHex(request.OwnerId)
	if err != nil {
		return nil, err
	}
	list, err := s.nlDb.GetList(ctx, ownerId, request.OwnerType)
	if err != nil {
		return nil, err
	}
	arr := []*apiPb.NotificationMethod{}
	for _, item := range list {
		method, err := s.nmDb.Get(ctx, item.NotificationMethodId)
		if err != nil {
			return nil, err
		}
		m, err := dbMethodToProto(method)
		if err != nil {
			return nil, err
		}
		arr = append(arr, m)
	}
	return &apiPb.GetListResponse{
		Methods: arr,
	}, nil
}

func New(
	nlDb database.NotificationListDb,
	nmDb database.NotificationMethodDb,
	client apiPb.StorageClient,
	integrations integrations.Integrations) apiPb.NotificationManagerServer {
	return &server{
		nlDb:         nlDb,
		nmDb:         nmDb,
		client:       client,
		integrations: integrations,
	}
}
