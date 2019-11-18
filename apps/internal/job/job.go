package job

import (
	"errors"
	"time"
)

const (
	connTimeout = time.Second * 5
)

var (
	grpcNotServing          = errors.New("STATUS_NOT_SERVING")
	connTimeoutError        = errors.New("CONNECTION_TIMEOUT")
	wrongConnectConfigError = errors.New("WRONG_CONNECTION_CONFIGURATION")
	connectionNotExistError = errors.New("CONNECTION_NOT_EXIST")
)


type Job interface {
	Do() error
}