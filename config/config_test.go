package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: return default value", func(t *testing.T) {
		s := New()
		assert.Equal(t, s.GetLogLevel(), "")
	})
}

func TestCfg_GetAgentServerTimeout(t *testing.T) {
	t.Run("Should: return execution timeout from env", func(t *testing.T) {
		os.Setenv(ENV_SQUZY_LOG_LEVEL, "ERROR")
		s := New()
		assert.Equal(t, s.GetLogLevel(), "ERROR")
	})
}
