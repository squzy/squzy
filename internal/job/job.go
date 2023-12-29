package job

import (
	"errors"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
)

var (
	errGrpcNotServing          = errors.New("STATUS_NOT_SERVING")
	errConnTimeoutError        = errors.New("CONNECTION_TIMEOUT")
	errWrongConnectConfigError = errors.New("WRONG_CONNECTION_CONFIGURATION")

	mongoConnectionError     = errors.New("UNABLE_TO_CONNECT_MONGO")
	mongoPingError           = errors.New("NO_PING_MONGO")
	postgresConnectionError  = errors.New("UNABLE_TO_CONNECT_POSTGRES")
	postgresPingError        = errors.New("NO_PING_POSTGRES")
	cassandraConnectionError = errors.New("UNABLE_TO_CONNECT_CASSANDRA")
	cassandraPingError       = errors.New("NO_PING_CASSANDRA")
	mysqlConnectionError     = errors.New("UNABLE_TO_CONNECT_MYSQL")
	mysqlPingError           = errors.New("NO_PING_MYSQL")
)

type CheckError interface {
	GetLogData() *apiPb.SchedulerResponse
}
