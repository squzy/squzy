package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Shoud: return default value", func(t *testing.T) {
		s := New()
		assert.Equal(t, s.GetPort(), int32(8080))
	})
}

func TestCfg_GetPort(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv("PORT", "11124")
		s := New()
		assert.Equal(t, s.GetPort(), int32(11124))
	})
}

func TestConfig_GetDbHost(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		expected := "localhost"
		os.Setenv("DB_HOST", expected)
		s := New()
		assert.Equal(t, s.GetDbHost(), expected)
	})
}

func TestConfig_GetDbPort(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		expected := "localhost"
		os.Setenv("DB_PORT", expected)
		s := New()
		assert.Equal(t, s.GetDbPort(), expected)
	})
}

func TestConfig_GetDbName(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		expected := "localhost"
		os.Setenv("DB_NAME", expected)
		s := New()
		assert.Equal(t, s.GetDbName(), expected)
	})
}

func TestConfig_GetDbUser(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		expected := "localhost"
		os.Setenv("DB_USER", expected)
		s := New()
		assert.Equal(t, s.GetDbUser(), expected)
	})
}

func TestConfig_GetDbPassword(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		expected := "localhost"
		os.Setenv("DB_PASSWORD", expected)
		s := New()
		assert.Equal(t, s.GetDbPassword(), expected)
	})
}