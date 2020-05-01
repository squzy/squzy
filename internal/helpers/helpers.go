package helpers

import (
	"context"
	"strings"
	"time"
)

const (
	httpPort               = int32(80)
	httpsPort              = int32(443)
	defaultTimeout         = 10
	defaultTimeoutDuration = time.Second * defaultTimeout
)

func GetPortByUrl(url string) int32 {
	if strings.HasPrefix(url, "https") {
		return httpsPort
	}
	return httpPort
}

func DurationFromSecond(seconds int32) time.Duration {
	return time.Duration(seconds) * time.Second
}

func TimeoutContext(parentCtx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout.Seconds() <= 0 {
		timeout = defaultTimeoutDuration
	}
	return context.WithTimeout(parentCtx, timeout)
}
