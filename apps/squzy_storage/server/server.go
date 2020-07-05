package server

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
	"squzy/apps/squzy_storage/config"
	"squzy/internal/database"
	"squzy/internal/helpers"
	"time"
)

type server struct {
	database       database.Database
	incidentClient apiPb.IncidentServerClient
	cfg            config.Config
}

func NewServer(db database.Database, incidentClient apiPb.IncidentServerClient, cfg config.Config) apiPb.StorageServer {
	return &server{
		database:       db,
		incidentClient: incidentClient,
		cfg:            cfg,
	}
}

func (s *server) SaveResponseFromScheduler(ctx context.Context, request *apiPb.SchedulerResponse) (*empty.Empty, error) {
	err := s.database.InsertSnapshot(request)
	defer func() {
		if request == nil {
			return
		}
		s.SendRecordToIncident(&apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_Snapshot{
				Snapshot: &apiPb.SchedulerSnapshotWithId{
					Snapshot: request.GetSnapshot(),
					Id:       request.GetSchedulerId(),
				},
			},
		})
	}()
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func (s *server) SaveResponseFromAgent(ctx context.Context, request *apiPb.Metric) (*empty.Empty, error) {
	err := s.database.InsertStatRequest(request)
	defer func() {
		if request == nil {
			return
		}
		s.SendRecordToIncident(&apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_AgentMetric{
				AgentMetric: request,
			},
		})
	}()
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func (s *server) GetSchedulerInformation(ctx context.Context, request *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error) {
	snapshots, count, err := s.database.GetSnapshots(request)
	return &apiPb.GetSchedulerInformationResponse{
		Snapshots: snapshots,
		Count:     count,
	}, wrapError(err)
}

func (s *server) GetSchedulerUptime(ctx context.Context, request *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	response, err := s.database.GetSnapshotsUptime(request)
	return response, wrapError(err)
}

func (s *server) GetAgentInformation(ctx context.Context, request *apiPb.GetAgentInformationRequest) (*apiPb.GetAgentInformationResponse, error) {
	var res []*apiPb.GetAgentInformationResponse_Statistic
	var count int32
	var err error
	switch request.GetType() {
	case apiPb.TypeAgentStat_ALL:
		res, count, err = s.database.GetStatRequest(request.GetAgentId(), request.GetPagination(), request.GetTimeRange())
	case apiPb.TypeAgentStat_CPU:
		res, count, err = s.database.GetCPUInfo(request.GetAgentId(), request.GetPagination(), request.GetTimeRange())
	case apiPb.TypeAgentStat_MEMORY:
		res, count, err = s.database.GetMemoryInfo(request.GetAgentId(), request.GetPagination(), request.GetTimeRange())
	case apiPb.TypeAgentStat_DISK:
		res, count, err = s.database.GetDiskInfo(request.GetAgentId(), request.GetPagination(), request.GetTimeRange())
	case apiPb.TypeAgentStat_NET:
		res, count, err = s.database.GetNetInfo(request.GetAgentId(), request.GetPagination(), request.GetTimeRange())
	default:
		err = errors.New("invalid type")
	}
	return &apiPb.GetAgentInformationResponse{
		Stats: res,
		Count: count,
	}, wrapError(err)
}

func (s *server) SaveTransaction(ctx context.Context, req *apiPb.TransactionInfo) (*empty.Empty, error) {
	defer func() {
		if req == nil {
			return
		}
		s.SendRecordToIncident(&apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_Transaction{
				Transaction: req,
			},
		})
	}()
	return &empty.Empty{}, wrapError(s.database.InsertTransactionInfo(req))
}

func (s *server) GetTransactions(ctx context.Context, request *apiPb.GetTransactionsRequest) (*apiPb.GetTransactionsResponse, error) {
	transactions, count, err := s.database.GetTransactionInfo(request)
	return &apiPb.GetTransactionsResponse{
		Count:        count,
		Transactions: transactions,
	}, wrapError(err)
}

func (s *server) GetTransactionById(ctx context.Context, request *apiPb.GetTransactionByIdRequest) (*apiPb.GetTransactionByIdResponse, error) {
	transaction, children, err := s.database.GetTransactionByID(request)
	return &apiPb.GetTransactionByIdResponse{
		Transaction: transaction,
		Children:    children,
	}, wrapError(err)
}

func (s *server) GetTransactionsGroup(ctx context.Context, request *apiPb.GetTransactionGroupRequest) (*apiPb.GetTransactionGroupResponse, error) {
	res, err := s.database.GetTransactionGroup(request)
	return &apiPb.GetTransactionGroupResponse{
		Transactions: res,
	}, wrapError(err)
}

func wrapError(err error) error {
	if err != nil {
		return grpcStatus.Errorf(codes.Internal, err.Error())
	}
	return nil
}

func (s *server) SaveIncident(ctx context.Context, request *apiPb.Incident) (*empty.Empty, error) {
	return &empty.Empty{}, s.database.InsertIncident(request)
}

func (s *server) UpdateIncidentStatus(ctx context.Context, request *apiPb.UpdateIncidentStatusRequest) (*apiPb.Incident, error) {
	return s.database.UpdateIncidentStatus(request.GetIncidentId(), request.GetStatus())
}

func (s *server) GetIncidentById(ctx context.Context, request *apiPb.IncidentIdRequest) (*apiPb.Incident, error) {
	return s.database.GetIncidentById(request.GetIncidentId())
}

func (s *server) GetIncidentByRuleId(ctx context.Context, request *apiPb.RuleIdRequest) (*apiPb.Incident, error) {
	return s.database.GetActiveIncidentByRuleId(request.GetRuleId())
}

func (s *server) GetIncidentsList(ctx context.Context, request *apiPb.GetIncidentsListRequest) (*apiPb.GetIncidentsListResponse, error) {
	incidents, count, err := s.database.GetIncidents(request)
	if err != nil {
		return nil, err
	}
	return &apiPb.GetIncidentsListResponse{
		Count:     count,
		Incidents: incidents,
	}, nil
}

func (s *server) SendRecordToIncident(rq *apiPb.StorageRecord) {
	if !s.cfg.WithIncident() {
		return
	}
	if s.incidentClient == nil {
		return
	}
	go func() {
		ctx, cancel := helpers.TimeoutContext(context.Background(), time.Second*5)
		defer cancel()
		_, _ = s.incidentClient.ProcessRecordFromStorage(ctx, rq)
	}()
}
