package storage_client

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"time"
)

type Storage interface {
	SetStatus(ctx context.Context, id string, status apiPb.IncidentStatus) (*apiPb.Incident, error)
	GetIncident(ctx context.Context, ruleID string) (*apiPb.Incident, error)
	SaveIncident(ctx context.Context, incident *apiPb.Incident) error
	GetLastTransactions(ctx context.Context, count int, filter *FilterTransaction) []*apiPb.TransactionInfo
	GetFirstTransactions(ctx context.Context, count int, filter *FilterTransaction) []*apiPb.TransactionInfo
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

type FilterTransaction struct {
	TimeFrom time.Time
	TimeTo   time.Time
	Type     apiPb.TransactionType
	Status   apiPb.TransactionStatus
}

//As this used in expr, need to throw panic
func (s *storage) GetLastTransactions(ctx context.Context, count int, filter *FilterTransaction) []*apiPb.TransactionInfo {
	transactions, err := s.client.GetTransactions(ctx, &apiPb.GetTransactionsRequest{
		ApplicationId: "", //TODO:
		TimeRange:     &apiPb.TimeFilter{
			From:                 getTimestamp(filter.TimeFrom),
			To:                   getTimestamp(filter.TimeTo),
		},
		Type:          filter.Type,
		Status:        filter.Status,
	})
	if err != nil {
		panic(err)
	}
	return transactions.GetTransactions()
}

//As this used in expr, need to throw panic
func (s *storage) GetFirstTransactions(ctx context.Context, count int, filter *FilterTransaction) []*apiPb.TransactionInfo {
	transactions, err := s.client.GetTransactions(ctx, &apiPb.GetTransactionsRequest{
		ApplicationId: "", //TODO:
		TimeRange:     &apiPb.TimeFilter{
			From:                 getTimestamp(filter.TimeFrom),
			To:                   getTimestamp(filter.TimeTo),
		},
		Type:          filter.Type,
		Status:        filter.Status,
		Sort: &apiPb.SortingTransactionList{
			Direction: apiPb.SortDirection_ASC,
		},
	})
	if err != nil {
		panic(err)
	}
	return transactions.GetTransactions()
}

//As this used in expr, need to throw panic
func getTimestamp(t time.Time) *timestamp.Timestamp {
	res, err := ptypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}
	return res
}