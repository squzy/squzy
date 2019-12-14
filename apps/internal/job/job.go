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
	postgresConnectionError = errors.New("UNABLE_TO_CONNECT_POSTGRES")
	postgresPingError       = errors.New("NO_PING_POSTGRES")
	mysqlConnectionError    = errors.New("UNABLE_TO_CONNECT_MYSQL")
	mysqlPingError          = errors.New("NO_PING_MYSQL")
)

type CheckError interface {
	GetLogData() *clientPb.Log
}

type Job interface {
	Do() CheckError
}
