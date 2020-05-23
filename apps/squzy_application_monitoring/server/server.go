package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

type server struct {

}

func (s *server) InitializeApplication(ctx context.Context, info *apiPb.ApplicationInfo) (*apiPb.InitializeApplicationResponse, error) {
	panic("implement me")
}

func (s *server) SaveTransaction(ctx context.Context, info *apiPb.TransactionInfo) (*empty.Empty, error) {
	panic("implement me")
}

func (s *server) SaveSegment(ctx context.Context, info *apiPb.SegmentInfo) (*empty.Empty, error) {
	panic("implement me")
}

func New() apiPb.ApplicationMonitoringServer {
	return &server{}
}