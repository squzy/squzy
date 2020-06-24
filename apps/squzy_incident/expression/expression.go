package expression

import (
	"github.com/araddon/dateparse"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

type Expression interface {
	IsValidTransaction(applicationId string, rule string) bool
}

type expressionStruct struct {
	storageClient apiPb.StorageClient
}

func NewExpr(storage apiPb.StorageClient) Expression {
	return &expressionStruct{
		storageClient: storage,
	}
}

func convertToTimestamp(strTime string) *timestamp.Timestamp {
	t, err := dateparse.ParseAny(strTime)
	if err != nil {
		panic(err)
	}
	res, err := ptypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}
	return res
}
