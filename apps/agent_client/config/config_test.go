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
		assert.Equal(t, s.GetAgentServer(), "")
		assert.Equal(t, s.GetInterval(), defaultTimeout)
		assert.Equal(t, s.GetAgentName(), "")
	})
}

func TestCfg_GetSquzyServer(t *testing.T) {
	t.Run("Should: return server address from env", func(t *testing.T) {
		os.Setenv(ENV_SQUZY_AGENT_SERVER_HOST, "11124")
		s := New()
		assert.Equal(t, s.GetAgentServer(), "11124")
	})
}

func TestCfg_GetAgentServerTimeout(t *testing.T) {
	t.Run("Should: return execution timeout from env", func(t *testing.T) {
		os.Setenv(ENV_SQUZY_AGENT_INTERVAL, "12")
		s := New()
		assert.Equal(t, s.GetInterval(), time.Second*12)
	})
}

func TestCfg_Retry(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv(ENV_SQUZY_AGENT_RETRY, "false")
		s := New()
		assert.Equal(t, s.Retry(), false)
	})
}

func TestCfg_RetryCount(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv(ENV_SQUZY_AGENT_RETRY_COUNT, "11")
		s := New()
		assert.Equal(t, s.RetryCount(), int32(11))
	})
}

func TestCfg_GetAgentName(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv(ENV_SQUZY_AGENT_NAME, "11124")
		s := New()
		assert.Equal(t, s.GetAgentName(), "11124")
	})
}
