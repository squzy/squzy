package helpers

import (
	"context"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
	"testing"
	"time"
)

func TestGetPortByUrl(t *testing.T) {
	t.Run("Should: return int32(80)", func(t *testing.T) {
		assert.Equal(t, int32(80), GetPortByURL("http://google.com"))
	})
	t.Run("Should: return int32(443)", func(t *testing.T) {
		assert.Equal(t, int32(443), GetPortByURL("https://google.com"))
	})
}

func TestDurationFromSecond(t *testing.T) {
	t.Run("Should: be equal", func(t *testing.T) {
		assert.Equal(t, time.Second*5, DurationFromSecond(5))
	})
}

func TestTimeoutContext(t *testing.T) {
	t.Run("Should: create context with timeout", func(t *testing.T) {
		ctx, cancel := TimeoutContext(context.Background(), time.Second)
		defer cancel()
		deadline, _ := ctx.Deadline()
		assert.Equal(t, time.Now().Add(time.Second).Unix(), deadline.Unix())
	})
	t.Run("Should: create context with default time is less then 0", func(t *testing.T) {
		ctx, cancel := TimeoutContext(context.Background(), -time.Second)
		defer cancel()
		deadline, _ := ctx.Deadline()
		assert.Equal(t, time.Now().Add(defaultTimeoutDuration).Unix(), deadline.Unix())
	})
}

func TestSelectorsToDb(t *testing.T) {
	t.Run("Should: convert correct", func(t *testing.T) {
		assert.EqualValues(t, []*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_STRING,
				Path: "select",
			},
		}, SelectorsToDb([]*apiPb.HttpJsonValueConfig_Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_STRING,
				Path: "select",
			},
		}))
	})
}

func TestSelectorsToProto(t *testing.T) {
	t.Run("Should: convert correct", func(t *testing.T) {
		assert.EqualValues(t, []*apiPb.HttpJsonValueConfig_Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_STRING,
				Path: "select",
			},
		}, SelectorsToProto([]*scheduler_config_storage.Selectors{
			{
				Type: apiPb.HttpJsonValueConfig_STRING,
				Path: "select",
			},
		}))
	})
}
