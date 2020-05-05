package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"squzy/apps/squzy_agent_server/database"
)

type server struct {
	db     database.Database
	client apiPb.StorageClient
}

func (s *server) Register(ctx context.Context, rq *apiPb.RegisterRequest) (*apiPb.RegisterResponse, error) {
	id, err := s.db.Add(ctx, rq)
	if err != nil {
		return nil, err
	}
	return &apiPb.RegisterResponse{
		Id: id,
	}, nil
}

func (s server) GetByAgentName(ctx context.Context, name *apiPb.GetByAgentNameRequest) (*apiPb.GetAgentListResponse, error) {
	agents, err := s.db.GetAll(ctx, bson.M{
		"agentName": bson.M{
			"$eq": name,
		},
	})
	if err != nil {
		return nil, err
	}

	return &apiPb.GetAgentListResponse{
		Agents: agents,
	}, nil
}

func (s *server) UnRegister(ctx context.Context, rq *apiPb.UnRegisterRequest) (*apiPb.UnRegisterResponse, error) {
	agentID, err := primitive.ObjectIDFromHex(rq.Id)
	if err != nil {
		return nil, err
	}
	err = s.db.UpdateStatus(context.Background(), agentID, apiPb.AgentStatus_UNREGISTRED)
	if err != nil {
		return nil, err
	}
	return &apiPb.UnRegisterResponse{
		Id: rq.Id,
	}, nil
}

func (s *server) GetAgentList(ctx context.Context, e *empty.Empty) (*apiPb.GetAgentListResponse, error) {
	agents, err := s.db.GetAll(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	return &apiPb.GetAgentListResponse{
		Agents: agents,
	}, nil
}

func (s *server) SendMetrics(stream apiPb.AgentServer_SendMetricsServer) error {
	var agentId primitive.ObjectID

	for {

		stat, err := stream.Recv()

		if err != nil {
			_ = s.db.UpdateStatus(context.Background(), agentId, apiPb.AgentStatus_DISCONNECTED)
			return stream.SendAndClose(&empty.Empty{})
		}

		id, err := primitive.ObjectIDFromHex(stat.AgentId)

		if err != nil {
			return err
		}

		agentId = id

		_ = s.db.UpdateStatus(context.Background(), agentId, apiPb.AgentStatus_RUNNED)
		// _, _ = s.client.SendResponseFromAgent(context.Background(), stat)
	}
}

func New(db database.Database, client apiPb.StorageClient) apiPb.AgentServerServer {
	return &server{
		db:     db,
		client: client,
	}
}
