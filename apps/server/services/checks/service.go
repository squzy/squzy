package checks

import (
	"context"
	"fmt"
	storagePb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"squzy/apps/internal/database"
)

type server struct {
	db database.Database
}

func (s *server) SendLogMessage(ctx context.Context,rq  *storagePb.SendLogMessageRequest) (*storagePb.SendLogMessageResponse, error) {
	fmt.Println(rq.SchedulerId, rq.Log.Code)
	return &storagePb.SendLogMessageResponse{
		Success:true,
	}, nil
}

func New(db database.Database) storagePb.LoggerServer {
	return &server{
		db: db,
	}
}