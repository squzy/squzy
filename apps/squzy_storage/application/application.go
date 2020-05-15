package application

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
	"os"
	"squzy/apps/squzy_storage/config"
	"squzy/internal/database"
)

type service struct {
	database database.Database
	config   config.Config
}

func NewService(cnfg config.Config) (apiPb.StorageServer, error) {
	db, err := database.New(func() (db *gorm.DB, e error) {
		return gorm.Open(
			"postgres",
			fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s connect_timeout=10 sslmode=disable",
				cnfg.GetDbHost(),
				cnfg.GetDbPort(),
				cnfg.GetDbName(),
				cnfg.GetDbUser(),
				cnfg.GetDbPassword(),
			))
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error connect to database"))
		//TODO: logger
		return nil, err
	}
	return &service{
		database: db,
		config:   cnfg,
	}, nil
}

func (s *service) SendResponseFromScheduler(ctx context.Context, request *apiPb.SchedulerResponse) (*apiPb.SendResponseFromSchedulerRequest, error) {
	err := s.database.InsertSnapshot(request)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, err.Error())
	}
	return &apiPb.SendResponseFromSchedulerRequest{
		Config: &apiPb.SendResponseFromSchedulerRequest_OkResult{
			OkResult: "OK",
		},
	}, nil
}

func (s *service) SendResponseFromAgent(ctx context.Context, request *apiPb.SendMetricsRequest) (*apiPb.SendResponseFromAgentResponse, error) {
	err := s.database.InsertStatRequest(request)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, err.Error())
	}
	return &apiPb.SendResponseFromAgentResponse{
		Config: &apiPb.SendResponseFromAgentResponse_OkResult{
			OkResult: "OK",
		},
	}, nil
}

func (s *service) GetSchedulerInformation(ctx context.Context, request *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error) {
	snapshots, err := s.database.GetSnapshots(request.SchedulerId)
	return &apiPb.GetSchedulerInformationResponse{
		Snapshots: snapshots,
	}, wrapError(err)
}

func (s *service) GetAgentInformation(ctx context.Context, request *apiPb.GetAgentInformationRequest) (*apiPb.SendMetricsRequest, error) {
	res, err := s.database.GetStatRequest(request.GetAgentId())
	return res, wrapError(err)
}

func wrapError(err error) error {
	if err != nil {
		return grpcStatus.Errorf(codes.Internal, err.Error())
	}
	return nil
}