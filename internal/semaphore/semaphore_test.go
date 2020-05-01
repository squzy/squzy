package semaphore

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSem(t *testing.T) {
	t.Run("Should: return nil error", func(t *testing.T) {
		s := NewSemaphore(10)
		actual := s.Acquire(context.Background())
		assert.NoError(t, actual)
		s.Release()
	})
}
