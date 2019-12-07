package job

import (
	"errors"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"time"
)

const (
	connTimeout = time.Second * 5
)

var (
	grpcNotServing          = errors.New("STATUS_NOT_SERVING")
	connTimeoutError        = errors.New("CONNECTION_TIMEOUT")
	wrongConnectConfigError = errors.New("WRONG_CONNECTION_CONFIGURATION")
	mongoConnectionError    = errors.New("UNABLE_TO_CONNECT_MONGO")
	mongoPingError          = errors.New("NO_PING_MONGO")
)

type CheckError interface {
	GetLogData() *clientPb.Log
}

type Job interface {
	Do() CheckError
}
