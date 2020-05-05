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

func (s *server) SendResponseFromScheduler(context.Context, *apiPb.SchedulerResponse) (*apiPb.SendResponseFromSchedulerRequest, error) {
	return nil, nil  //TODO
}

func (s *server) SendResponseFromAgent(ctx context.Context, request *apiPb.SendMetricsRequest) (*apiPb.SendResponseFromAgentResponse, error) {
	return nil, s.database.InsertStatRequest(request) //TODO
}

func (s *server)  GetSchedulerInformation(context.Context, *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error)  {
	return nil, nil  //TODO
}

func (s *server)  GetAgentInformation(context.Context, *apiPb.GetAgentInformationRequest) (*apiPb.SendMetricsRequest, error) {
	return nil, nil  //TODO
}
