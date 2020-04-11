package agents

import (
	"context"
	agentPb "github.com/squzy/squzy_generated/generated/agent/proto/v1"
	"squzy/apps/internal/database"
)

type service struct {
	db database.Database
}

func (s *service) Register(context.Context, *agentPb.RegisterRequest) (*agentPb.RegisterResponse, error) {
	return nil, nil
}

func (s *service) UnRegister(context.Context, *agentPb.UnRegisterRequest) (*agentPb.UnRegisterResponse, error) {
	return nil, nil
}

func (s *service) SendStat(agentPb.AgentServer_SendStatServer) error {
	return nil
}

func (s *service) GetList(context.Context, *agentPb.GetListRequest) (*agentPb.GetListResponse, error) {
	return nil, nil
}

func New(db database.Database) agentPb.AgentServerServer {
	return &service{
		db:db,
	}
}