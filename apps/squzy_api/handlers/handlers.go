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
	GetSchedulerList(ctx context.Context) ([]*apiPb.Scheduler, error)
	GetSchedulerByID(ctx context.Context, id string) (*apiPb.Scheduler, error)
	GetSchedulerHistoryByID(ctx context.Context, rq *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error)
	GetAgentHistoryByID(ctx context.Context, rq *apiPb.GetAgentInformationRequest) (*apiPb.GetAgentInformationResponse, error)
	RunScheduler(ctx context.Context, id string) error
	StopScheduler(ctx context.Context, id string) error
	RemoveScheduler(ctx context.Context, id string) error
	AddScheduler(ctx context.Context, scheduler *apiPb.AddRequest) error
}

type handlers struct {
	agentClient      apiPb.AgentServerClient
	monitoringClient apiPb.SchedulersExecutorClient
	storageClient    apiPb.StorageClient
}

func (h *handlers) GetSchedulerHistoryByID(ctx context.Context, rq *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, 0)
	defer cancel()
	r, err := h.storageClient.GetSchedulerInformation(c, rq)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (h *handlers) GetAgentHistoryByID(ctx context.Context, rq *apiPb.GetAgentInformationRequest) (*apiPb.GetAgentInformationResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, 0)
	defer cancel()
	r, err := h.storageClient.GetAgentInformation(c, rq)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (h *handlers) RunScheduler(ctx context.Context, id string) error {
	c, cancel := helpers.TimeoutContext(ctx, 0)
	defer cancel()
	_, err := h.monitoringClient.Run(c, &apiPb.RunRequest{
		Id: id,
	})
	return err
}

func (h *handlers) AddScheduler(ctx context.Context, scheduler *apiPb.AddRequest) error {
	c, cancel := helpers.TimeoutContext(ctx, 0)
	defer cancel()
	_, err := h.monitoringClient.Add(c, scheduler)
	return err
}

func (h *handlers) StopScheduler(ctx context.Context, id string) error {
	c, cancel := helpers.TimeoutContext(ctx, 0)
	defer cancel()
	_, err := h.monitoringClient.Stop(c, &apiPb.StopRequest{
		Id: id,
	})
	return err
}

func (h *handlers) RemoveScheduler(ctx context.Context, id string) error {
	c, cancel := helpers.TimeoutContext(ctx, 0)
	defer cancel()
	_, err := h.monitoringClient.Remove(c, &apiPb.RemoveRequest{
		Id: id,
	})
	return err
}

func (h *handlers) GetSchedulerList(ctx context.Context) ([]*apiPb.Scheduler, error) {
	c, cancel := helpers.TimeoutContext(ctx, 0)
	defer cancel()
	list, err := h.monitoringClient.GetSchedulerList(c, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return list.List, nil
}

func (h *handlers) GetSchedulerByID(ctx context.Context, id string) (*apiPb.Scheduler, error) {
	c, cancel := helpers.TimeoutContext(ctx, 0)
	defer cancel()
	scheduler, err := h.monitoringClient.GetSchedulerById(c, &apiPb.GetSchedulerByIdRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return scheduler, nil
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
