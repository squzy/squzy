package handlers

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"squzy/internal/helpers"
)

type Handlers interface {
	GetAgentList(ctx context.Context) ([]*apiPb.AgentItem, error)
	GetAgentByID(ctx context.Context, id string) (*apiPb.AgentItem, error)
}

type handlers struct {
	agentClient      apiPb.AgentServerClient
	monitoringClient apiPb.SchedulersExecutorClient
}

func (h *handlers) GetAgentList(ctx context.Context) ([]*apiPb.AgentItem, error) {
	c, cancel := helpers.TimeoutContext(ctx, 0)
	defer cancel()
	list, err := h.agentClient.GetAgentList(c, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return list.Agents, nil
}

func (h *handlers) GetAgentByID(ctx context.Context, id string) (*apiPb.AgentItem, error) {
	c, cancel := helpers.TimeoutContext(ctx, 0)
	defer cancel()
	agent, err := h.agentClient.GetAgentById(c, &apiPb.GetAgentByIdRequest{
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
		agentClient:      agentClient,
		monitoringClient: monitoringClient,
	}
}
