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
		assert.Equal(t, s.GetPort(), int32(8080))
		assert.Equal(t, s.GetClientAddress(), "")
	})
}

func TestCfg_GetClientAddress(t *testing.T) {

}

func TestCfg_GetPort(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv("PORT", "11124")
		s := New()
		assert.Equal(t, s.GetPort(), int32(11124))
	})
}
func TestCfg_GetStorageTimeout(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		os.Setenv("STORAGE_TIMEOUT", "11")
		s := New()
		assert.Equal(t, s.GetStorageTimeout(), time.Second * 11)
	})
}