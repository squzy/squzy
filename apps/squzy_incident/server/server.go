package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)


type server struct {

}

func NewIncidentServer() apiPb.IncidentServerServer {
	return &server{}
}

func (s server) CreateRule(context.Context, *apiPb.CreateRuleRequest) (*apiPb.Rule, error) {
	panic("implement me")
}

func (s server) GetRuleById(context.Context, *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	panic("implement me")
}

func (s server) GetRulesByOwnerId(context.Context, *apiPb.GetRulesByOwnerIdRequest) (*apiPb.Rules, error) {
	panic("implement me")
}

func (s server) RemoveRule(context.Context, *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	panic("implement me")
}

func (s server) ProcessRecordFromStorage(context.Context, *apiPb.StorageRecord) (*empty.Empty, error) {
	panic("implement me")
}

func (s server) CloseIncident(context.Context, *apiPb.IncidentIdRequest) (*apiPb.Incident, error) {
	panic("implement me")
}

func (s server) StudyIncident(context.Context, *apiPb.IncidentIdRequest) (*apiPb.Incident, error) {
	panic("implement me")
}
