package helpers

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetPortByUrl(t *testing.T) {
	t.Run("Should: return int32(80)", func(t *testing.T) {
		assert.Equal(t, int32(80), GetPortByUrl("http://google.com"))
	})
	t.Run("Should: return int32(443)", func(t *testing.T) {
		assert.Equal(t, int32(443), GetPortByUrl("https://google.com"))
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
