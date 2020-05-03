package server

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"io"
	"time"
)

type server struct {
	
}

func (s *server) Register(ctx context.Context,rq *apiPb.RegisterRequest) (*apiPb.RegisterResponse, error) {
	fmt.Println(rq)
	return &apiPb.RegisterResponse{
		Id:                   "1",
	}, nil
}

func (s server) GetByAgentUniqName(context.Context, *apiPb.GetByAgentUniqNameRequest) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (s *server) UnRegister(ctx context.Context,rq *apiPb.UnRegisterRequest) (*apiPb.UnRegisterResponse, error) {
	return &apiPb.UnRegisterResponse{
		Id:                   rq.Id,
	}, nil
}

func (s *server) GetAgentList(context.Context, *empty.Empty) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (s *server) SendMetrics(stream apiPb.AgentServer_SendMetricsServer) error {
	for {
		stat, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&empty.Empty{})
		}

		if err != nil {
			return err
		}
		fmt.Println(time.Now().Format(time.RFC3339))
		fmt.Println(stat)
	}
}

func New() apiPb.AgentServerServer {
	return &server{}
}