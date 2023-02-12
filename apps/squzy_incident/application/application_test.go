package application

import (
	"context"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	empty "google.golang.org/protobuf/types/known/emptypb"
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
	apiPb.UnimplementedIncidentServerServer
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
			_ = s.Run(23234)
		}()
		time.Sleep(time.Second * 2)
		_, err := net.Dial("tcp", "localhost:23234")
		assert.Equal(t, nil, err)
	})
}
