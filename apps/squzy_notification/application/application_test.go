package application

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	t.Run("Should: work", func(t *testing.T) {
		s := New(nil)
		assert.NotNil(t, s)
	})
}

type mockNotification struct {
	
}

func (m mockNotification) CreateNotificationMethod(ctx context.Context, request *apiPb.CreateNotificationMethodRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (m mockNotification) GetById(ctx context.Context, request *apiPb.NotificationMethodIdRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (m mockNotification) DeleteById(ctx context.Context, request *apiPb.NotificationMethodIdRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (m mockNotification) Activate(ctx context.Context, request *apiPb.NotificationMethodIdRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (m mockNotification) Deactivate(ctx context.Context, request *apiPb.NotificationMethodIdRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (m mockNotification) Add(ctx context.Context, request *apiPb.NotificationMethodRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (m mockNotification) Remove(ctx context.Context, request *apiPb.NotificationMethodRequest) (*apiPb.NotificationMethod, error) {
	panic("implement me")
}

func (m mockNotification) GetList(ctx context.Context, request *apiPb.GetListRequest) (*apiPb.GetListResponse, error) {
	panic("implement me")
}

func (m mockNotification) Notify(ctx context.Context, request *apiPb.NotifyRequest) (*empty.Empty, error) {
	panic("implement me")
}

func TestServer_Run(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := New(nil)
		assert.Error(t, s.Run(124124))
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockNotification{})
		go func() {
			_ = s.Run(23235)
		}()
		time.Sleep(time.Second * 2)
		_, err := net.Dial("tcp", "localhost:23235")
		assert.Equal(t, nil, err)
	})
}
