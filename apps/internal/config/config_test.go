package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	t.Run("Tests:", func(t *testing.T) {
		t.Run("Should: Create new config", func(t *testing.T) {
			cfg := NewConfig()
			assert.Implements(t, (*Config)(nil), cfg)
		})
	})
}