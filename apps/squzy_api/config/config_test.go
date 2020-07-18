package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New()
		assert.NotEqual(t, nil, s)
	})
}

func TestCfg_GetAgentServerAddress(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv(ENV_AGENT_SERVER, "11124")
		s := New()
		assert.Equal(t, s.GetAgentServerAddress(), "11124")
	})
}

func TestCfg_GetMonitoringServerAddress(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv(ENV_MONITORING_SERVER, "11124")
		s := New()
		assert.Equal(t, s.GetMonitoringServerAddress(), "11124")
	})
}

func TestCfg_GetApplicationMonitoringAddress(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv(ENV_APPLICATION_MONITORING_SERVER, "11124")
		s := New()
		assert.Equal(t, s.GetApplicationMonitoringAddress(), "11124")
	})
}

func TestCfg_GetPort(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv(ENV_PORT, "11124")
		s := New()
		assert.Equal(t, s.GetPort(), int32(11124))
	})
}

func TestCfg_GetStorageServerAddress(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv(ENV_STORAGE_SERVER, "11124")
		s := New()
		assert.Equal(t, s.GetStorageServerAddress(), "11124")
	})
}

func TestCfg_GetIncidentServerAddress(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv(ENV_INCIDENT_SERVER, "11124")
		s := New()
		assert.Equal(t, s.GetIncidentServerAddress(), "11124")
	})
}

func TestCfg_GetNotificationServerAddress(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv(ENV_NOTIFICATION_SERVER, "11124")
		s := New()
		assert.Equal(t, s.GetNotificationServerAddress(), "11124")
	})
}
