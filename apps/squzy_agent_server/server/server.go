package server

import (
	"context"
	"github.com/squzy/squzy/apps/squzy_agent_server/database"
	"github.com/squzy/squzy/internal/helpers"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	empty "google.golang.org/protobuf/types/known/emptypb"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	db     database.Database
	client apiPb.StorageClient
}

func (s *server) GetAgentById(ctx context.Context, rq *apiPb.GetAgentByIdRequest) (*apiPb.AgentItem, error) {
	agentID, err := primitive.ObjectIDFromHex(rq.AgentId)
	if err != nil {
		return nil, err
	}
	res, err := s.db.GetByID(ctx, agentID)
	if err != nil {
		return nil, err
	}
	return res, nil
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

func (s server) GetByAgentName(ctx context.Context, rq *apiPb.GetByAgentNameRequest) (*apiPb.GetAgentListResponse, error) {
	agents, err := s.db.GetAll(ctx, bson.M{
		"agentName": bson.M{
			"$eq": rq.AgentName,
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
	err = s.db.UpdateStatus(context.Background(), agentID, apiPb.AgentStatus_UNREGISTRED, rq.Time)
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
	stat, err := stream.Recv()

	if err != nil {
		return err
	}

	switch msg := stat.Msg.(type) {
	case *apiPb.SendMetricsRequest_Metric:
		id, err := primitive.ObjectIDFromHex(msg.Metric.AgentId)

		if err != nil {
			return err
		}

		_ = s.db.UpdateStatus(context.Background(), id, apiPb.AgentStatus_RUNNED, msg.Metric.Time)

		for {
			incomeMsg, err := stream.Recv()

			if err != nil {
				// If error happens try to disconnect
				ctx, cancel := helpers.TimeoutContext(context.Background(), 0)
				defer cancel()
				_ = s.db.UpdateStatus(ctx, id, apiPb.AgentStatus_DISCONNECTED, timestamp.Now())
				return stream.SendAndClose(&empty.Empty{})
			}

			switch newMsg := incomeMsg.Msg.(type) {
			case *apiPb.SendMetricsRequest_Metric:
				ctx, cancel := helpers.TimeoutContext(context.Background(), 0)
				defer cancel()
				_, _ = s.client.SaveResponseFromAgent(ctx, newMsg.Metric)
				continue
			case *apiPb.SendMetricsRequest_Disconnect_:
				ctx, cancel := helpers.TimeoutContext(context.Background(), 0)
				defer cancel()
				_ = s.db.UpdateStatus(ctx, id, apiPb.AgentStatus_DISCONNECTED, newMsg.Disconnect.Time)
				return stream.SendAndClose(&empty.Empty{})
			}
		}
	default:
		return stream.SendAndClose(&empty.Empty{})
	}
}

func New(db database.Database, client apiPb.StorageClient) apiPb.AgentServerServer {
	return &server{
		db:     db,
		client: client,
	}
}
