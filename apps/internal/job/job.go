package job

import (
	"errors"
	"time"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
)

const (
	connTimeout = time.Second * 5
)

var (
	grpcNotServing          = errors.New("STATUS_NOT_SERVING")
	connTimeoutError        = errors.New("CONNECTION_TIMEOUT")
	wrongConnectConfigError = errors.New("WRONG_CONNECTION_CONFIGURATION")
)

type CheckError interface {
	GetLogData() *clientPb.Log
}

type Job interface {
	Do() CheckError
}
