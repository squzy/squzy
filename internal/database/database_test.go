package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: not return nil", func(t *testing.T) {
		s, err := New(nil)
		assert.NotNil(t, err)
		assert.Nil(t, s)
	})
}
