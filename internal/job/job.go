package job

import (
	"errors"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

var (
	errGrpcNotServing          = errors.New("STATUS_NOT_SERVING")
	errConnTimeoutError        = errors.New("CONNECTION_TIMEOUT")
	errWrongConnectConfigError = errors.New("WRONG_CONNECTION_CONFIGURATION")
)

type CheckError interface {
	GetLogData() *apiPb.SchedulerResponse
}
