package server

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"squzy/apps/squzy_notification/database"
	"testing"
)

type mockMethodSuccessActiveWebhook struct {

}

func (m mockMethodSuccessActiveWebhook) Create(ctx context.Context, nm *database.NotificationMethod) error {
	panic("implement me")
}

func (m mockMethodSuccessActiveWebhook) Delete(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}

func (m mockMethodSuccessActiveWebhook) Activate(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}

func (m mockMethodSuccessActiveWebhook) Deactivate(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}

func (m mockMethodSuccessActiveWebhook) Get(ctx context.Context, id primitive.ObjectID) (*database.NotificationMethod, error) {
	return &database.NotificationMethod{
		Status: api.NotificationMethodStatus_NOTIFICATION_STATUS_ACTIVE,
		Type:    api.NotificationMethodType_NOTIFICATION_METHOD_WEBHOOK,
		WebHook:   &database.WebHookConfig{Url: ""},
	}, nil
}

type mockMethodSuccessActiveSlack struct {

}

func (m mockMethodSuccessActiveSlack) Create(ctx context.Context, nm *database.NotificationMethod) error {
	panic("implement me")
}

func (m mockMethodSuccessActiveSlack) Delete(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}

func (m mockMethodSuccessActiveSlack) Activate(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}

func (m mockMethodSuccessActiveSlack) Deactivate(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}

func (m mockMethodSuccessActiveSlack) Get(ctx context.Context, id primitive.ObjectID) (*database.NotificationMethod, error) {
	return &database.NotificationMethod{
		Status: api.NotificationMethodStatus_NOTIFICATION_STATUS_ACTIVE,
		Type:    api.NotificationMethodType_NOTIFICATION_METHOD_SLACK,
		Slack:   &database.SlackConfig{Url: ""},
	}, nil
}

type mockIntegration struct {

}

func (m mockIntegration) Slack(ctx context.Context, incident *api.Incident, config *database.SlackConfig) {
	return
}

func (m mockIntegration) Webhook(ctx context.Context, incident *api.Incident, config *database.WebHookConfig) {
	return
}

type mockStorageError struct {

}

func (m mockStorageError) SaveResponseFromScheduler(ctx context.Context, in *api.SchedulerResponse, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockStorageError) SaveResponseFromAgent(ctx context.Context, in *api.Metric, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockStorageError) SaveTransaction(ctx context.Context, in *api.TransactionInfo, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockStorageError) GetSchedulerInformation(ctx context.Context, in *api.GetSchedulerInformationRequest, opts ...grpc.CallOption) (*api.GetSchedulerInformationResponse, error) {
	panic("implement me")
}

func (m mockStorageError) GetSchedulerUptime(ctx context.Context, in *api.GetSchedulerUptimeRequest, opts ...grpc.CallOption) (*api.GetSchedulerUptimeResponse, error) {
	panic("implement me")
}

func (m mockStorageError) GetAgentInformation(ctx context.Context, in *api.GetAgentInformationRequest, opts ...grpc.CallOption) (*api.GetAgentInformationResponse, error) {
	panic("implement me")
}

func (m mockStorageError) GetTransactionsGroup(ctx context.Context, in *api.GetTransactionGroupRequest, opts ...grpc.CallOption) (*api.GetTransactionGroupResponse, error) {
	panic("implement me")
}

func (m mockStorageError) GetTransactions(ctx context.Context, in *api.GetTransactionsRequest, opts ...grpc.CallOption) (*api.GetTransactionsResponse, error) {
	panic("implement me")
}

func (m mockStorageError) GetTransactionById(ctx context.Context, in *api.GetTransactionByIdRequest, opts ...grpc.CallOption) (*api.GetTransactionByIdResponse, error) {
	panic("implement me")
}

func (m mockStorageError) SaveIncident(ctx context.Context, in *api.Incident, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockStorageError) UpdateIncidentStatus(ctx context.Context, in *api.UpdateIncidentStatusRequest, opts ...grpc.CallOption) (*api.Incident, error) {
	panic("implement me")
}

func (m mockStorageError) GetIncidentById(ctx context.Context, in *api.IncidentIdRequest, opts ...grpc.CallOption) (*api.Incident, error) {
	return nil, errors.New("")
}

func (m mockStorageError) GetIncidentByRuleId(ctx context.Context, in *api.RuleIdRequest, opts ...grpc.CallOption) (*api.Incident, error) {
	panic("implement me")
}

func (m mockStorageError) GetIncidentsList(ctx context.Context, in *api.GetIncidentsListRequest, opts ...grpc.CallOption) (*api.GetIncidentsListResponse, error) {
	panic("implement me")
}

type mockStorageOk struct {

}

func (m mockStorageOk) SaveResponseFromScheduler(ctx context.Context, in *api.SchedulerResponse, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockStorageOk) SaveResponseFromAgent(ctx context.Context, in *api.Metric, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockStorageOk) SaveTransaction(ctx context.Context, in *api.TransactionInfo, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockStorageOk) GetSchedulerInformation(ctx context.Context, in *api.GetSchedulerInformationRequest, opts ...grpc.CallOption) (*api.GetSchedulerInformationResponse, error) {
	panic("implement me")
}

func (m mockStorageOk) GetSchedulerUptime(ctx context.Context, in *api.GetSchedulerUptimeRequest, opts ...grpc.CallOption) (*api.GetSchedulerUptimeResponse, error) {
	panic("implement me")
}

func (m mockStorageOk) GetAgentInformation(ctx context.Context, in *api.GetAgentInformationRequest, opts ...grpc.CallOption) (*api.GetAgentInformationResponse, error) {
	panic("implement me")
}

func (m mockStorageOk) GetTransactionsGroup(ctx context.Context, in *api.GetTransactionGroupRequest, opts ...grpc.CallOption) (*api.GetTransactionGroupResponse, error) {
	panic("implement me")
}

func (m mockStorageOk) GetTransactions(ctx context.Context, in *api.GetTransactionsRequest, opts ...grpc.CallOption) (*api.GetTransactionsResponse, error) {
	panic("implement me")
}

func (m mockStorageOk) GetTransactionById(ctx context.Context, in *api.GetTransactionByIdRequest, opts ...grpc.CallOption) (*api.GetTransactionByIdResponse, error) {
	panic("implement me")
}

func (m mockStorageOk) SaveIncident(ctx context.Context, in *api.Incident, opts ...grpc.CallOption) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockStorageOk) UpdateIncidentStatus(ctx context.Context, in *api.UpdateIncidentStatusRequest, opts ...grpc.CallOption) (*api.Incident, error) {
	panic("implement me")
}

func (m mockStorageOk) GetIncidentById(ctx context.Context, in *api.IncidentIdRequest, opts ...grpc.CallOption) (*api.Incident, error) {
	return &api.Incident{
		Id:                   "",
		Status:               0,
		RuleId:               "",
		Histories:            nil,
	}, nil
}

func (m mockStorageOk) GetIncidentByRuleId(ctx context.Context, in *api.RuleIdRequest, opts ...grpc.CallOption) (*api.Incident, error) {
	panic("implement me")
}

func (m mockStorageOk) GetIncidentsList(ctx context.Context, in *api.GetIncidentsListRequest, opts ...grpc.CallOption) (*api.GetIncidentsListResponse, error) {
	panic("implement me")
}

type mockMethodSuccessTypeWrong struct {

}

func (m mockMethodSuccessTypeWrong) Create(ctx context.Context, nm *database.NotificationMethod) error {
	panic("implement me")
}

func (m mockMethodSuccessTypeWrong) Delete(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}

func (m mockMethodSuccessTypeWrong) Activate(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}

func (m mockMethodSuccessTypeWrong) Deactivate(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}

func (m mockMethodSuccessTypeWrong) Get(ctx context.Context, id primitive.ObjectID) (*database.NotificationMethod, error) {
	return  &database.NotificationMethod{
		Id:      primitive.ObjectID{},
		Status:  0,
		Type:    0,
	}, nil
}

type mockMethodSuccessSecond struct {

}

func (m mockMethodSuccessSecond) Create(ctx context.Context, nm *database.NotificationMethod) error {
	panic("implement me")
}

func (m mockMethodSuccessSecond) Delete(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}

func (m mockMethodSuccessSecond) Activate(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}

func (m mockMethodSuccessSecond) Deactivate(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}

func (m mockMethodSuccessSecond) Get(ctx context.Context, id primitive.ObjectID) (*database.NotificationMethod, error) {
	return &database.NotificationMethod{
		Id:      primitive.ObjectID{},
		Status:  0,
		Type:    api.NotificationMethodType_NOTIFICATION_METHOD_SLACK,
		Slack:   &database.SlackConfig{Url: ""},
	}, nil
}

type mockListSuccess struct {

}

func (m mockListSuccess) Add(ctx context.Context, notification *database.Notification) error {
	return nil
}

func (m mockListSuccess) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m mockListSuccess) GetList(ctx context.Context, OwnerId primitive.ObjectID, Type api.ComponentOwnerType) ([]*database.Notification, error) {
	return []*database.Notification{
		{
			Id: primitive.NewObjectID(),
			NotificationMethodId: primitive.NewObjectID(),
		},
	}, nil
}

type mockListInternalError struct {

}

func (m mockListInternalError) Add(ctx context.Context, notification *database.Notification) error {
	return errors.New("")
}

func (m mockListInternalError) Delete(ctx context.Context, id primitive.ObjectID) error {
	return errors.New("")
}

func (m mockListInternalError) GetList(ctx context.Context, OwnerId primitive.ObjectID, Type api.ComponentOwnerType) ([]*database.Notification, error) {
	return nil, errors.New("")
}

type mockMethodSuccess struct {
	
}

func (m mockMethodSuccess) Create(ctx context.Context, nm *database.NotificationMethod) error {
	return nil
}

func (m mockMethodSuccess) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m mockMethodSuccess) Activate(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m mockMethodSuccess) Deactivate(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m mockMethodSuccess) Get(ctx context.Context, id primitive.ObjectID) (*database.NotificationMethod, error) {
	return &database.NotificationMethod{
		Type: api.NotificationMethodType_NOTIFICATION_METHOD_WEBHOOK,
		WebHook: &database.WebHookConfig{Url: ""},
	}, nil
}

type mockMethodInternalError struct {
	
}

func (m mockMethodInternalError) Create(ctx context.Context, nm *database.NotificationMethod) error {
	return errors.New("")
}

func (m mockMethodInternalError) Delete(ctx context.Context, id primitive.ObjectID) error {
	return errors.New("")
}

func (m mockMethodInternalError) Activate(ctx context.Context, id primitive.ObjectID) error {
	return errors.New("")
}

func (m mockMethodInternalError) Deactivate(ctx context.Context, id primitive.ObjectID) error {
	return errors.New("")
}

func (m mockMethodInternalError) Get(ctx context.Context, id primitive.ObjectID) (*database.NotificationMethod, error) {
	return nil, errors.New("")
}

type mockMethodNotFoundError struct {
	
}

func (m mockMethodNotFoundError) Create(ctx context.Context, nm *database.NotificationMethod) error {
	return nil
}

func (m mockMethodNotFoundError) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m mockMethodNotFoundError) Activate(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m mockMethodNotFoundError) Deactivate(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m mockMethodNotFoundError) Get(ctx context.Context, id primitive.ObjectID) (*database.NotificationMethod, error) {
	return nil, errors.New("")
}

func TestNew(t *testing.T) {
	t.Run("Should: not be nil", func(t *testing.T) {
		s := New(nil, nil, nil, nil)
		assert.NotNil(t, s)
	})
}

func TestServer_Activate(t *testing.T) {
	t.Run("Should: throw error because not bson", func(t *testing.T) {
		s := New(nil, nil, nil, nil)
		_, err := s.Activate(context.Background(), &api.NotificationMethodIdRequest{Id: ""})
		assert.NotNil(t, err)
	})
	t.Run("Should: throw error because internal", func(t *testing.T) {
		s := New(nil, &mockMethodInternalError{}, nil, nil)
		_, err := s.Activate(context.Background(), &api.NotificationMethodIdRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: return error not found", func(t *testing.T) {
		s := New(nil, &mockMethodNotFoundError{}, nil, nil)
		_, err := s.Activate(context.Background(), &api.NotificationMethodIdRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, &mockMethodSuccess{}, nil, nil)
		_, err := s.Activate(context.Background(), &api.NotificationMethodIdRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.Nil(t, err)
	})
}

func TestServer_Deactivate(t *testing.T) {
	t.Run("Should: throw error because not bson", func(t *testing.T) {
		s := New(nil, nil, nil, nil)
		_, err := s.Deactivate(context.Background(), &api.NotificationMethodIdRequest{Id: ""})
		assert.NotNil(t, err)
	})
	t.Run("Should: throw error because internal", func(t *testing.T) {
		s := New(nil, &mockMethodInternalError{}, nil, nil)
		_, err := s.Deactivate(context.Background(), &api.NotificationMethodIdRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: return error not found", func(t *testing.T) {
		s := New(nil, &mockMethodNotFoundError{}, nil, nil)
		_, err := s.Deactivate(context.Background(), &api.NotificationMethodIdRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, &mockMethodSuccess{}, nil, nil)
		_, err := s.Deactivate(context.Background(), &api.NotificationMethodIdRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.Nil(t, err)
	})
}

func TestServer_DeleteById(t *testing.T) {
	t.Run("Should: throw error because not bson", func(t *testing.T) {
		s := New(nil, nil, nil, nil)
		_, err := s.DeleteById(context.Background(), &api.NotificationMethodIdRequest{Id: ""})
		assert.NotNil(t, err)
	})
	t.Run("Should: throw error because internal", func(t *testing.T) {
		s := New(nil, &mockMethodInternalError{}, nil, nil)
		_, err := s.DeleteById(context.Background(), &api.NotificationMethodIdRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: return error not found", func(t *testing.T) {
		s := New(nil, &mockMethodNotFoundError{}, nil, nil)
		_, err := s.DeleteById(context.Background(), &api.NotificationMethodIdRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, &mockMethodSuccess{}, nil, nil)
		_, err := s.DeleteById(context.Background(), &api.NotificationMethodIdRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.Nil(t, err)
	})
}

func TestServer_GetById(t *testing.T) {
	t.Run("Should: throw error because not bson", func(t *testing.T) {
		s := New(nil, nil, nil, nil)
		_, err := s.GetById(context.Background(), &api.NotificationMethodIdRequest{Id: ""})
		assert.NotNil(t, err)
	})
	t.Run("Should: throw error because not exist type", func(t *testing.T) {
		s := New(nil, &mockMethodSuccessTypeWrong{}, nil, nil)
		_, err := s.GetById(context.Background(), &api.NotificationMethodIdRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.Equal(t, errTypeNotExist, err)
	})
	t.Run("Should: throw internal error", func(t *testing.T) {
		s := New(nil, &mockMethodInternalError{}, nil, nil)
		_, err := s.GetById(context.Background(), &api.NotificationMethodIdRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, &mockMethodSuccess{}, nil, nil)
		_, err := s.GetById(context.Background(), &api.NotificationMethodIdRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.Nil(t, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, &mockMethodSuccessSecond{}, nil, nil)
		_, err := s.GetById(context.Background(), &api.NotificationMethodIdRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.Nil(t, err)
	})
}

func TestServer_Remove(t *testing.T) {
	t.Run("Should: throw error because not bson", func(t *testing.T) {
		s := New(nil, nil, nil, nil)
		_, err := s.Remove(context.Background(), &api.NotificationMethodRequest{
			NotificationMethodId: "",
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: throw error because internal error", func(t *testing.T) {
		s := New(&mockListInternalError{}, nil, nil, nil)
		_, err := s.Remove(context.Background(), &api.NotificationMethodRequest{
			NotificationMethodId: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: throw error because internal error", func(t *testing.T) {
		s := New(&mockListSuccess{}, &mockMethodNotFoundError{}, nil, nil)
		_, err := s.Remove(context.Background(), &api.NotificationMethodRequest{
			NotificationMethodId: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: not throw error", func(t *testing.T) {
		s := New(&mockListSuccess{}, &mockMethodSuccess{}, nil, nil)
		_, err := s.Remove(context.Background(), &api.NotificationMethodRequest{
			NotificationMethodId: primitive.NewObjectID().Hex(),
		})
		assert.Nil(t, err)
	})
}

func TestServer_Add(t *testing.T) {
	t.Run("Should: return error because not bson", func(t *testing.T) {
		s := New(nil, nil, nil, nil)
		_, err := s.Add(context.Background(), &api.NotificationMethodRequest{
			NotificationMethodId: "",
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: return error because not bson", func(t *testing.T) {
		s := New(nil, nil, nil, nil)
		_, err := s.Add(context.Background(), &api.NotificationMethodRequest{
			OwnerId: primitive.NewObjectID().Hex(),
			NotificationMethodId: "",
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: throw error because internal error", func(t *testing.T) {
		s := New(&mockListInternalError{}, nil, nil, nil)
		_, err := s.Add(context.Background(), &api.NotificationMethodRequest{
			OwnerId: primitive.NewObjectID().Hex(),
			NotificationMethodId: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: throw error because internal error", func(t *testing.T) {
		s := New(&mockListSuccess{}, &mockMethodNotFoundError{}, nil, nil)
		_, err := s.Add(context.Background(), &api.NotificationMethodRequest{
			OwnerId: primitive.NewObjectID().Hex(),
			NotificationMethodId: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: not throw error", func(t *testing.T) {
		s := New(&mockListSuccess{}, &mockMethodSuccess{}, nil, nil)
		_, err := s.Add(context.Background(), &api.NotificationMethodRequest{
			OwnerId: primitive.NewObjectID().Hex(),
			NotificationMethodId: primitive.NewObjectID().Hex(),
		})
		assert.Nil(t, err)
	})
}

func TestServer_GetList(t *testing.T) {
	t.Run("Should: return error because not bson", func(t *testing.T) {
		s := New(nil, nil, nil, nil)
		_, err := s.GetList(context.Background(), &api.GetListRequest{
			OwnerId: "",
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: return internal error", func(t *testing.T) {
		s := New(&mockListInternalError{}, nil, nil, nil)
		_, err := s.GetList(context.Background(), &api.GetListRequest{
			OwnerId: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: return internal error", func(t *testing.T) {
		s := New(&mockListSuccess{}, &mockMethodNotFoundError{}, nil, nil)
		_, err := s.GetList(context.Background(), &api.GetListRequest{
			OwnerId: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: return internal error", func(t *testing.T) {
		s := New(&mockListSuccess{}, &mockMethodSuccessTypeWrong{}, nil, nil)
		_, err := s.GetList(context.Background(), &api.GetListRequest{
			OwnerId: primitive.NewObjectID().Hex(),
		})
		assert.Equal(t, errTypeNotExist, err)
	})
	t.Run("Should: return response", func(t *testing.T) {
		s := New(&mockListSuccess{}, &mockMethodSuccess{}, nil, nil)
		_, err := s.GetList(context.Background(), &api.GetListRequest{
			OwnerId: primitive.NewObjectID().Hex(),
		})
		assert.Nil(t, err)
	})
}

func TestServer_CreateNotificationMethod(t *testing.T) {
	t.Run("Should: return error because wrong type", func(t *testing.T) {
		s := New(nil, nil, nil, nil)
		_, err := s.CreateNotificationMethod(context.Background(), &api.CreateNotificationMethodRequest{})
		assert.Equal(t, errTypeNotExist, err)
	})
	t.Run("Should: return error because internal", func(t *testing.T) {
		s := New(nil, &mockMethodInternalError{}, nil, nil)
		_, err := s.CreateNotificationMethod(context.Background(), &api.CreateNotificationMethodRequest{
			Type: api.NotificationMethodType_NOTIFICATION_METHOD_SLACK,
			Method: &api.CreateNotificationMethodRequest_Slack{
					Slack: &api.SlackMethod{
						Url:                  "",
					},
			},
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, &mockMethodSuccess{}, nil, nil)
		_, err := s.CreateNotificationMethod(context.Background(), &api.CreateNotificationMethodRequest{
			Type: api.NotificationMethodType_NOTIFICATION_METHOD_SLACK,
			Method: &api.CreateNotificationMethodRequest_Slack{
				Slack: &api.SlackMethod{
					Url:                  "",
				},
			},
		})
		assert.Nil(t, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(nil, &mockMethodSuccess{}, nil, nil)
		_, err := s.CreateNotificationMethod(context.Background(), &api.CreateNotificationMethodRequest{
			Type: api.NotificationMethodType_NOTIFICATION_METHOD_WEBHOOK,
			Method: &api.CreateNotificationMethodRequest_Webhook{
				Webhook: &api.WebHookMethod{
					Url:                  "",
				},
			},
		})
		assert.Nil(t, err)
	})
}

func TestServer_Notify(t *testing.T) {
	t.Run("Should: return error because not bson", func(t *testing.T) {
		s := New(nil, nil, nil, nil)
		_, err := s.Notify(context.Background(), &api.NotifyRequest{
			OwnerId: "",
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: return error because storage", func(t *testing.T) {
		s := New(nil, nil, &mockStorageError{}, nil)
		_, err := s.Notify(context.Background(), &api.NotifyRequest{
			OwnerId: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: return error because internal", func(t *testing.T) {
		s := New(&mockListInternalError{}, nil, &mockStorageOk{}, nil)
		_, err := s.Notify(context.Background(), &api.NotifyRequest{
			OwnerId: primitive.NewObjectID().Hex(),
		})
		assert.NotNil(t, err)
	})
	t.Run("Should: return not return ", func(t *testing.T) {
		s := New(&mockListSuccess{}, &mockMethodNotFoundError{}, &mockStorageOk{}, &mockIntegration{})
		_, err := s.Notify(context.Background(), &api.NotifyRequest{
			OwnerId: primitive.NewObjectID().Hex(),
		})
		assert.Nil(t, err)
	})
	t.Run("Should: return not return ", func(t *testing.T) {
		s := New(&mockListSuccess{}, &mockMethodSuccess{}, &mockStorageOk{}, &mockIntegration{})
		_, err := s.Notify(context.Background(), &api.NotifyRequest{
			OwnerId: primitive.NewObjectID().Hex(),
		})
		assert.Nil(t, err)
	})
	t.Run("Should: return not return ", func(t *testing.T) {
		s := New(&mockListSuccess{}, &mockMethodSuccessActiveSlack{}, &mockStorageOk{}, &mockIntegration{})
		_, err := s.Notify(context.Background(), &api.NotifyRequest{
			OwnerId: primitive.NewObjectID().Hex(),
		})
		assert.Nil(t, err)
	})
	t.Run("Should: return not return ", func(t *testing.T) {
		s := New(&mockListSuccess{}, &mockMethodSuccessActiveWebhook{}, &mockStorageOk{}, &mockIntegration{})
		_, err := s.Notify(context.Background(), &api.NotifyRequest{
			OwnerId: primitive.NewObjectID().Hex(),
		})
		assert.Nil(t, err)
	})
}