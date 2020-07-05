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
		s := NewApplication(nil)
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
		s := NewApplication(nil)
		assert.Error(t, s.Run(124124))
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := NewApplication(mockApiIncident{})
		go func() {
			_ = s.Run(23233)
		}()
		time.Sleep(time.Second)
		_, err := net.Dial("tcp", "localhost:23233")
		assert.Equal(t, nil, err)
	})
}
