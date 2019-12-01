package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInMemoryStorage(t *testing.T) {
	t.Run("Should: not throw error", func(t *testing.T) {
		s := GetInMemoryStorage()
		assert.Equal(t, nil, s.Write("", &mock{}))
	})
}