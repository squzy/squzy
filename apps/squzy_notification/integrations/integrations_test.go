package integrations

import (
	"context"
	"errors"
	"github.com/squzy/squzy/apps/squzy_notification/database"
	api "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"testing"
	"time"
)

type mockCfg struct {
}

func (m mockCfg) GetPort() int32 {
	panic("implement me")
}

func (m mockCfg) GetMongoURI() string {
	panic("implement me")
}

func (m mockCfg) GetMongoDB() string {
	panic("implement me")
}

func (m mockCfg) GetNotificationMethodCollection() string {
	panic("implement me")
}

func (m mockCfg) GetNotificationListCollection() string {
	panic("implement me")
}

func (m mockCfg) GetStorageHost() string {
	panic("implement me")
}

func (m mockCfg) GetDashboardHost() string {
	return ""
}

type mock struct {
}

type mockError struct {
}

func (m mockError) SendRequest(req *http.Request) (int, []byte, error) {
	return 0, nil, errors.New("")
}

func (m mockError) SendRequestTimeout(req *http.Request, timeout time.Duration) (int, []byte, error) {
	panic("implement me")
}

func (m mockError) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	panic("implement me")
}

func (m mockError) SendRequestTimeoutStatusCode(req *http.Request, timeout time.Duration, expectedCode int) (int, []byte, error) {
	panic("implement me")
}

func (m mockError) CreateRequest(method string, url string, headers *map[string]string, schedulerID string) *http.Request {
	panic("implement me")
}

func (m mock) SendRequest(req *http.Request) (int, []byte, error) {
	return 0, []byte{}, nil
}

func (m mock) SendRequestTimeout(req *http.Request, timeout time.Duration) (int, []byte, error) {
	panic("implement me")
}

func (m mock) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	panic("implement me")
}

func (m mock) SendRequestTimeoutStatusCode(req *http.Request, timeout time.Duration, expectedCode int) (int, []byte, error) {
	panic("implement me")
}

func (m mock) CreateRequest(method string, url string, headers *map[string]string, schedulerID string) *http.Request {
	panic("implement me")
}

func TestNew(t *testing.T) {
	t.Run("Shuld: not be nil", func(t *testing.T) {
		s := New(nil, nil)
		assert.NotNil(t, s)
	})
}

func TestIntegrations_Webhook(t *testing.T) {
	t.Run("Should: not throw panic", func(t *testing.T) {
		s := New(&mock{}, &mockCfg{})
		assert.NotPanics(t, func() {
			s.Webhook(context.Background(), &api.Incident{
				Histories: []*api.Incident_HistoryItem{
					{
						Timestamp: timestamp.Now(),
						Status:    api.IncidentStatus_INCIDENT_STATUS_CAN_BE_CLOSED,
					},
				},
			}, &database.WebHookConfig{})
		})
	})
	t.Run("Should: not throw panic", func(t *testing.T) {
		s := New(&mock{}, &mockCfg{})
		assert.NotPanics(t, func() {
			s.Webhook(context.Background(), &api.Incident{
				Histories: []*api.Incident_HistoryItem{
					{
						Timestamp: timestamp.Now(),
						Status:    123,
					},
				},
			}, &database.WebHookConfig{})
		})
	})
	t.Run("Should: not throw panic", func(t *testing.T) {
		s := New(&mockError{}, &mockCfg{})
		assert.NotPanics(t, func() {
			s.Webhook(context.Background(), &api.Incident{
				Histories: []*api.Incident_HistoryItem{
					{
						Timestamp: timestamp.Now(),
						Status:    api.IncidentStatus_INCIDENT_STATUS_CAN_BE_CLOSED,
					},
				},
			}, &database.WebHookConfig{})
		})
	})
}

func TestIntegrations_Slack(t *testing.T) {
	t.Run("Should: not throw panic", func(t *testing.T) {
		s := New(&mock{}, &mockCfg{})
		assert.NotPanics(t, func() {
			s.Slack(context.Background(), &api.Incident{
				Histories: []*api.Incident_HistoryItem{
					{
						Timestamp: timestamp.Now(),
						Status:    api.IncidentStatus_INCIDENT_STATUS_CAN_BE_CLOSED,
					},
				},
			}, &database.SlackConfig{})
		})
	})
}
