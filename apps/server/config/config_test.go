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