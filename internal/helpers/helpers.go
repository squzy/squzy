package helpers

import (
	"context"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	"strings"
	"time"
)

const (
	httpPort               int32 = 80
	httpsPort              int32 = 443
	defaultTimeout               = 10
	defaultTimeoutDuration       = time.Second * defaultTimeout
)

func GetPortByURL(url string) int32 {
	if strings.HasPrefix(url, "https") {
		return httpsPort
	}
	return httpPort
}

func DurationNotNegative(seconds int32) time.Duration {
	if seconds <= 0 {
		return defaultTimeoutDuration
	}
	return DurationFromSecond(seconds)
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

func SelectorsToDb(selectors []*apiPb.HttpJsonValueConfig_Selectors) []*scheduler_config_storage.Selectors {
	arr := []*scheduler_config_storage.Selectors{}
	for _, v := range selectors {
		arr = append(arr, &scheduler_config_storage.Selectors{
			Type: v.Type,
			Path: v.Path,
		})
	}
	return arr
}

func SelectorsToProto(selectors []*scheduler_config_storage.Selectors) []*apiPb.HttpJsonValueConfig_Selectors {
	arr := []*apiPb.HttpJsonValueConfig_Selectors{}
	for _, v := range selectors {
		arr = append(arr, &apiPb.HttpJsonValueConfig_Selectors{
			Type: v.Type,
			Path: v.Path,
		})
	}
	return arr
}
