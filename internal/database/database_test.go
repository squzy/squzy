package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: not return nil", func(t *testing.T) {
		s := New(nil)
		assert.NotNil(t, s)
	})
}
