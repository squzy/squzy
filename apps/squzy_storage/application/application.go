package application

import (
	"context"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"squzy/internal/database"
)

type server struct {
	database database.Database
}

func NewStorage() *apiPb.StorageServer {
	return &server{}
}

func (s *server) SendResponseFromScheduler(ctx context.Context, request *apiPb.SchedulerResponse) (*apiPb.SendResponseFromSchedulerRequest, error) {
	err := s.database.InsertSnapshot(request)
	if err != nil {
		return &apiPb.SendResponseFromSchedulerRequest{
			Config: &apiPb.SendResponseFromSchedulerRequest_Error{
				Error: &apiPb.StrorageError{
					Description: err.Error(),
				},
			},
		}, nil
	}
	return &apiPb.SendResponseFromSchedulerRequest{
		Config: &apiPb.SendResponseFromSchedulerRequest_OkResult{
			OkResult: "OK",
		},
	}, nil
}

func (s *server) SendResponseFromAgent(ctx context.Context, request *apiPb.SendMetricsRequest) (*apiPb.SendResponseFromAgentResponse, error) {
	err := s.database.InsertStatRequest(request)
	if err != nil {
		return &apiPb.SendResponseFromAgentResponse{
			Config: &apiPb.SendResponseFromAgentResponse_Error{
				Error: &apiPb.StrorageError{
					Description: err.Error(),
				},
			},
		}, nil
	}
	return &apiPb.SendResponseFromAgentResponse{
		Config: &apiPb.SendResponseFromAgentResponse_OkResult{
			OkResult: "OK",
		},
	}, nil
}

func (s *server) GetSchedulerInformation(ctx context.Context, request *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error) {
	snapshots, err := s.database.GetSnapshots(request.SchedulerId)
	return &apiPb.GetSchedulerInformationResponse{
		Snapshots:            snapshots,
	}, err
}

func (s *server) GetAgentInformation(ctx context.Context, request *apiPb.GetAgentInformationRequest) (*apiPb.SendMetricsRequest, error) {
	return s.database.GetStatRequest(request.GetAgentId())
}
