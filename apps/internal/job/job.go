package job

import (
	"errors"
	"strings"
	"time"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
)

const (
	connTimeout = time.Second * 5
	httpPort  = int32(80)
	httpsPort = int32(443)
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

func GetPortByUrl(url string) int32 {
	if strings.HasPrefix(url, "https") {
		return httpsPort
	}
	return httpPort
}
