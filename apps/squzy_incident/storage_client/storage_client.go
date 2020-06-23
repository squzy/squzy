package storage_client

import (
	"context"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

type Storage interface {
	SetStatus(ctx context.Context, id string, status apiPb.IncidentStatus) (*apiPb.Incident, error)
	GetIncident(ctx context.Context, ruleID string) (*apiPb.Incident, error)
	SaveIncident(ctx context.Context, incident *apiPb.Incident) error
}

type storage struct {
	client apiPb.StorageClient
}

func New(client apiPb.StorageClient) Storage {
	return &storage{
		client: client,
	}
}

func (s *storage) SetStatus(ctx context.Context, id string, status apiPb.IncidentStatus) (*apiPb.Incident, error) {
	return s.client.UpdateIncidentStatus(ctx, &apiPb.UpdateIncidentStatusRequest{
		IncidentId: id,
		Status:     status,
	})
}

func (s *storage) GetIncident(ctx context.Context, ruleID string) (*apiPb.Incident, error) {
	return s.client.GetIncidentByRuleId(ctx, &apiPb.RuleIdRequest{
		RuleId: ruleID,
	})
}

func (s *storage) SaveIncident(ctx context.Context, incident *apiPb.Incident) error {
	_, err := s.client.SaveIncident(ctx, incident)
	return err
}
