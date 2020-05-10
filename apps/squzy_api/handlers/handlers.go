package handlers

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

type Handlers interface {
	GetAgentList(ctx context.Context) ([]*apiPb.AgentItem, error)
	GetAgentById(ctx context.Context, id string) (*apiPb.AgentItem, error)
}

type handlers struct {
	agentClient apiPb.AgentServerClient
	monitoringClient apiPb.SchedulersExecutorClient
}

func (h *handlers) GetAgentList(ctx context.Context) ([]*apiPb.AgentItem, error) {
	list, err := h.agentClient.GetAgentList(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return list.Agents, nil
}

func (h *handlers) GetAgentById(ctx context.Context, id string) (*apiPb.AgentItem, error) {
	agent, err := h.agentClient.GetAgentById(ctx, &apiPb.GetAgentByIdRequest{
		AgentId: id,
	})
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func New(
	agentClient apiPb.AgentServerClient,
	monitoringClient apiPb.SchedulersExecutorClient,
	) Handlers {
	return &handlers{
		agentClient: agentClient,
		monitoringClient: monitoringClient,
	}
}
