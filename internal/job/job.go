package job

import (
	"errors"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

var (
	grpcNotServing          = errors.New("STATUS_NOT_SERVING")
	connTimeoutError        = errors.New("CONNECTION_TIMEOUT")
	wrongConnectConfigError = errors.New("WRONG_CONNECTION_CONFIGURATION")
)

type CheckError interface {
	GetLogData() *apiPb.SchedulerResponse
}