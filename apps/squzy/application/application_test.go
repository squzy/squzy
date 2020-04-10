package application

import (
	"github.com/stretchr/testify/assert"
	"net"
	"squzy/apps/internal/semaphore"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("Should: Create new application", func(t *testing.T) {
		app := New(nil, nil, nil, nil, func(i int) semaphore.Semaphore {
			return semaphore.NewSemaphore(i)
		})
		assert.NotEqual(t, nil, app)
	})
}

func TestApp_Run(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		app := New(nil, nil, nil, nil, func(i int) semaphore.Semaphore {
			return semaphore.NewSemaphore(i)
		})
		assert.NotEqual(t, nil, app.Run(1244214))
	})
	t.Run("Should: not return error", func(t *testing.T) {
		app := New(nil, nil, nil, nil, func(i int) semaphore.Semaphore {
			return semaphore.NewSemaphore(i)
		})
		go func() {
			_ = app.Run(11111)
		}()
		time.Sleep(time.Second)
		_, err := net.Dial("tcp", "localhost:11111")
		assert.Equal(t, nil, err)
	})
}
