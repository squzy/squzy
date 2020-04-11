package server

import (
	"context"
	storagePb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
)

type server struct {

}

func (s server) SendLogMessage(ctz context.Context,rq  *storagePb.SendLogMessageRequest) (*storagePb.SendLogMessageResponse, error) {
	panic("implement me")
}

func New() storagePb.LoggerServer {
	return &server{}
}