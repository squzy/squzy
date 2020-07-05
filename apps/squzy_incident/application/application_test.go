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

type configErrorMock struct {
}

func (*configErrorMock) GetPort() int32 {
	return 1000000
}

func (*configErrorMock) GetStorageHost() string {
	panic("implement me!")
}

func (*configErrorMock) GetMongoURI() string {
	panic("implement me!")
}

func (*configErrorMock) GetMongoDb() string {
	panic("implement me!")
}

func (*configErrorMock) GetMongoCollection() string {
	panic("implement me!")
}

type configMock struct {
}

func (*configMock) GetPort() int32 {
	return 23233
}

func (*configMock) GetStorageHost() string {
	panic("implement me!")
}

func (*configMock) GetMongoURI() string {
	panic("implement me!")
}

func (*configMock) GetMongoDb() string {
	panic("implement me!")
}

func (*configMock) GetMongoCollection() string {
	panic("implement me!")
}

func TestNewServer(t *testing.T) {
	t.Run("Should: work", func(t *testing.T) {
		s := NewApplication(nil, nil)
		assert.NotNil(t, s)
	})
}

type mockApiIncident struct {
}

func (m mockApiIncident) CreateRule(context.Context, *apiPb.CreateRuleRequest) (*apiPb.Rule, error) {
	panic("implement me")
}

func (m mockApiIncident) GetRuleById(context.Context, *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	panic("implement me")
}

func (m mockApiIncident) GetRulesByOwnerId(context.Context, *apiPb.GetRulesByOwnerIdRequest) (*apiPb.Rules, error) {
	panic("implement me")
}

func (m mockApiIncident) RemoveRule(context.Context, *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	panic("implement me")
}

func (m mockApiIncident) ValidateRule(context.Context, *apiPb.ValidateRuleRequest) (*apiPb.ValidateRuleResponse, error) {
	panic("implement me")
}

func (m mockApiIncident) ProcessRecordFromStorage(context.Context, *apiPb.StorageRecord) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockApiIncident) CloseIncident(context.Context, *apiPb.IncidentIdRequest) (*apiPb.Incident, error) {
	panic("implement me")
}

func (m mockApiIncident) ActivateRule(context.Context, *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	panic("implement me")
}

func (m mockApiIncident) DeactivateRule(context.Context, *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	panic("implement me")
}

func (m mockApiIncident) StudyIncident(context.Context, *apiPb.IncidentIdRequest) (*apiPb.Incident, error) {
	panic("implement me")
}

func TestServer_Run(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := &application{
			config:  &configErrorMock{},
			apiServ: nil,
		}
		assert.Error(t, s.Run())
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := &application{
			config:  &configMock{},
			apiServ: &mockApiIncident{},
		}
		go func() {
			_ = s.Run()
		}()
		time.Sleep(time.Second)
		_, err := net.Dial("tcp", "localhost:23233")
		assert.Equal(t, nil, err)
	})
}
