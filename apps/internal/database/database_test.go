package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		_, s := New()
		assert.NotEqual(t, nil, s)
	})
}