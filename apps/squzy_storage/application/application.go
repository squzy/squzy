package application

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
	"squzy/internal/database"
)

type service struct {
	database database.Database
}

func (s *service) GetSchedulerUptime(ctx context.Context, request *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	panic("implement me")
}

func (s *service) SaveResponseFromScheduler(ctx context.Context, response *apiPb.SchedulerResponse) (*empty.Empty, error) {
	panic("implement me")
}

func (s *service) SaveResponseFromAgent(ctx context.Context, metric *apiPb.Metric) (*empty.Empty, error) {
	panic("implement me")
}

func (s *service) SaveTransaction(ctx context.Context, info *apiPb.TransactionInfo) (*empty.Empty, error) {
	panic("implement me")
}

func (s *service) GetTransactionsGroup(ctx context.Context, request *apiPb.GetTransactionGroupRequest) (*apiPb.GetTransactionGroupResponse, error) {
	panic("implement me")
}

func (s *service) GetTransactions(ctx context.Context, request *apiPb.GetTransactionsRequest) (*apiPb.GetTransactionsResponse, error) {
	panic("implement me")
}

func (s *service) GetTransactionById(ctx context.Context, request *apiPb.GetTransactionByIdRequest) (*apiPb.GetTransactionByIdResponse, error) {
	panic("implement me")
}

func NewService(db database.Database) apiPb.StorageServer {
	return &service{
		database: db,
	}
}

func (s *service) SendResponseFromScheduler(ctx context.Context, request *apiPb.SchedulerResponse) (*empty.Empty, error) {
	err := s.database.InsertSnapshot(request)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func (s *service) SendResponseFromAgent(ctx context.Context, request *apiPb.Metric) (*empty.Empty, error) {
	err := s.database.InsertStatRequest(request)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func (s *service) GetSchedulerInformation(ctx context.Context, request *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error) {
	snapshots, count, err := s.database.GetSnapshots(request.SchedulerId, request.GetPagination(), request.GetTimeRange())
	return &apiPb.GetSchedulerInformationResponse{
		Snapshots: snapshots,
		Count:     count,
	}, wrapError(err)
}

func (s *service) GetAgentInformation(ctx context.Context, request *apiPb.GetAgentInformationRequest) (*apiPb.GetAgentInformationResponse, error) {
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

func wrapError(err error) error {
	if err != nil {
		return grpcStatus.Errorf(codes.Internal, err.Error())
	}
	return nil
}
