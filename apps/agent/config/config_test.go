package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("Shoud: return default value", func(t *testing.T) {
		s := New()
		assert.Equal(t, s.GetSquzyServerTimeout(), defaultTimeout)
		assert.Equal(t, s.GetSquzyServer(), "")
	})
}

func TestCfg_GetSquzyServer(t *testing.T) {
	t.Run("Should: return server address from env", func(t *testing.T) {
		os.Setenv("SQUZY_SERVER_HOST", "11124")
		s := New()
		assert.Equal(t, s.GetSquzyServer(), "11124")
	})
}

func TestCfg_GetExecutionTimeout(t *testing.T) {
	t.Run("Should: return execution timeout from env", func(t *testing.T) {
		os.Setenv("SQUZY_EXECUTION_TIMEOUT", "12")
		s := New()
		assert.Equal(t, s.GetExecutionTimeout(), time.Second*12)
	})
}

func TestCfg_GetStorageTimeout(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv("SQUZY_SERVER_TIMEOUT", "11")
		s := New()
		assert.Equal(t, s.GetSquzyServerTimeout(), time.Second*11)
	})
}
