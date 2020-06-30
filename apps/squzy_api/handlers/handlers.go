package handlers

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"squzy/internal/helpers"
	"time"
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
	AddScheduler(ctx context.Context, scheduler *apiPb.AddRequest) (*apiPb.AddResponse, error)
	RegisterApplication(ctx context.Context, rq *apiPb.ApplicationInfo) (*apiPb.InitializeApplicationResponse, error)
	SaveTransaction(ctx context.Context, rq *apiPb.TransactionInfo) (*empty.Empty, error)
	GetSchedulerUptime(ctx context.Context, rq *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error)
	GetTransactionGroups(ctx context.Context, req *apiPb.GetTransactionGroupRequest) (*apiPb.GetTransactionGroupResponse, error)
	GetTransactionsList(ctx context.Context, req *apiPb.GetTransactionsRequest) (*apiPb.GetTransactionsResponse, error)
	GetApplicationById(ctx context.Context, id string) (*apiPb.Application, error)
	ArchivedApplicationById(ctx context.Context, id string) (*apiPb.Application, error)
	EnabledApplicationById(ctx context.Context, id string) (*apiPb.Application, error)
	DisabledApplicationById(ctx context.Context, id string) (*apiPb.Application, error)
	GetApplicationList(ctx context.Context) ([]*apiPb.Application, error)
	GetTransactionById(ctx context.Context, id string) (*apiPb.GetTransactionByIdResponse, error)
	CreateRule(ctx context.Context, rule *apiPb.CreateRuleRequest) (*apiPb.Rule, error)
	ValidateRule(ctx context.Context, rule *apiPb.ValidateRuleRequest) (*apiPb.ValidateRuleResponse, error)
	GetRulesByOwnerId(ctx context.Context, req *apiPb.GetRulesByOwnerIdRequest) (*apiPb.Rules, error)
	GetRuleById(ctx context.Context, req *apiPb.RuleIdRequest) (*apiPb.Rule, error)
	ActivateRuleById(ctx context.Context, req *apiPb.RuleIdRequest) (*apiPb.Rule, error)
	DeactivateRuleById(ctx context.Context, req *apiPb.RuleIdRequest) (*apiPb.Rule, error)
	RemoveRuleById(ctx context.Context, req *apiPb.RuleIdRequest) (*apiPb.Rule, error)
}

const (
	defaultRequestTimeout = time.Second * 30
)

type handlers struct {
	agentClient                 apiPb.AgentServerClient
	monitoringClient            apiPb.SchedulersExecutorClient
	storageClient               apiPb.StorageClient
	applicationMonitoringClient apiPb.ApplicationMonitoringClient
	incidentClient              apiPb.IncidentServerClient
}

func (h *handlers) GetRuleById(ctx context.Context, req *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.incidentClient.GetRuleById(c, req)
}

func (h *handlers) ActivateRuleById(ctx context.Context, req *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.incidentClient.ActivateRule(c, req)
}

func (h *handlers) DeactivateRuleById(ctx context.Context, req *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.incidentClient.DeactivateRule(c, req)
}

func (h *handlers) RemoveRuleById(ctx context.Context, req *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.incidentClient.RemoveRule(c, req)
}

func (h *handlers) GetRulesByOwnerId(ctx context.Context, req *apiPb.GetRulesByOwnerIdRequest) (*apiPb.Rules, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.incidentClient.GetRulesByOwnerId(c, req)
}

func (h *handlers) ValidateRule(ctx context.Context, rule *apiPb.ValidateRuleRequest) (*apiPb.ValidateRuleResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.incidentClient.ValidateRule(c, rule)
}

func (h *handlers) CreateRule(ctx context.Context, rule *apiPb.CreateRuleRequest) (*apiPb.Rule, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.incidentClient.CreateRule(c, rule)
}

func (h *handlers) ArchivedApplicationById(ctx context.Context, id string) (*apiPb.Application, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.applicationMonitoringClient.ArchiveApplicationById(c, &apiPb.ApplicationByIdReuqest{
		ApplicationId: id,
	})
}

func (h *handlers) EnabledApplicationById(ctx context.Context, id string) (*apiPb.Application, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.applicationMonitoringClient.EnableApplicationById(c, &apiPb.ApplicationByIdReuqest{
		ApplicationId: id,
	})
}

func (h *handlers) DisabledApplicationById(ctx context.Context, id string) (*apiPb.Application, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.applicationMonitoringClient.DisableApplicationById(c, &apiPb.ApplicationByIdReuqest{
		ApplicationId: id,
	})
}

func (h *handlers) SaveTransaction(ctx context.Context, rq *apiPb.TransactionInfo) (*empty.Empty, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.applicationMonitoringClient.SaveTransaction(c, rq)
}

func (h *handlers) GetSchedulerHistoryByID(ctx context.Context, rq *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.storageClient.GetSchedulerInformation(c, rq)
}

func (h *handlers) GetAgentHistoryByID(ctx context.Context, rq *apiPb.GetAgentInformationRequest) (*apiPb.GetAgentInformationResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.storageClient.GetAgentInformation(c, rq)
}

func (h *handlers) RunScheduler(ctx context.Context, id string) error {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	_, err := h.monitoringClient.Run(c, &apiPb.RunRequest{
		Id: id,
	})
	return err
}

func (h *handlers) AddScheduler(ctx context.Context, scheduler *apiPb.AddRequest) (*apiPb.AddResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.monitoringClient.Add(c, scheduler)
}

func (h *handlers) StopScheduler(ctx context.Context, id string) error {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	_, err := h.monitoringClient.Stop(c, &apiPb.StopRequest{
		Id: id,
	})
	return err
}

func (h *handlers) RemoveScheduler(ctx context.Context, id string) error {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	_, err := h.monitoringClient.Remove(c, &apiPb.RemoveRequest{
		Id: id,
	})
	return err
}

func (h *handlers) GetSchedulerList(ctx context.Context) ([]*apiPb.Scheduler, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	list, err := h.monitoringClient.GetSchedulerList(c, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return list.Lists, nil
}

func (h *handlers) GetSchedulerByID(ctx context.Context, id string) (*apiPb.Scheduler, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.monitoringClient.GetSchedulerById(c, &apiPb.GetSchedulerByIdRequest{
		Id: id,
	})
}

func (h *handlers) GetAgentList(ctx context.Context) ([]*apiPb.AgentItem, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	list, err := h.agentClient.GetAgentList(c, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return list.Agents, nil
}

func (h *handlers) GetAgentByID(ctx context.Context, id string) (*apiPb.AgentItem, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.agentClient.GetAgentById(c, &apiPb.GetAgentByIdRequest{
		AgentId: id,
	})
}

func (h *handlers) RegisterApplication(ctx context.Context, rq *apiPb.ApplicationInfo) (*apiPb.InitializeApplicationResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.applicationMonitoringClient.InitializeApplication(c, rq)
}

func (h *handlers) GetSchedulerUptime(ctx context.Context, rq *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.storageClient.GetSchedulerUptime(c, rq)
}

func (h *handlers) GetTransactionGroups(ctx context.Context, req *apiPb.GetTransactionGroupRequest) (*apiPb.GetTransactionGroupResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.storageClient.GetTransactionsGroup(c, req)
}

func (h *handlers) GetTransactionsList(ctx context.Context, req *apiPb.GetTransactionsRequest) (*apiPb.GetTransactionsResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.storageClient.GetTransactions(c, req)
}

func (h *handlers) GetTransactionById(ctx context.Context, id string) (*apiPb.GetTransactionByIdResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.storageClient.GetTransactionById(c, &apiPb.GetTransactionByIdRequest{
		TransactionId: id,
	})
}

func (h *handlers) GetApplicationById(ctx context.Context, id string) (*apiPb.Application, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	return h.applicationMonitoringClient.GetApplicationById(c, &apiPb.ApplicationByIdReuqest{
		ApplicationId: id,
	})
}

func (h *handlers) GetApplicationList(ctx context.Context) ([]*apiPb.Application, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	res, err := h.applicationMonitoringClient.GetApplicationList(c, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return res.Applications, err
}

func New(
	agentClient apiPb.AgentServerClient,
	monitoringClient apiPb.SchedulersExecutorClient,
	storageClient apiPb.StorageClient,
	applicationMonitoringClient apiPb.ApplicationMonitoringClient,
	incidentClient apiPb.IncidentServerClient,
) Handlers {
	return &handlers{
		agentClient:                 agentClient,
		monitoringClient:            monitoringClient,
		storageClient:               storageClient,
		applicationMonitoringClient: applicationMonitoringClient,
		incidentClient:              incidentClient,
	}
}
