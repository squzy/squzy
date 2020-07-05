package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: return default value", func(t *testing.T) {
		s := New()
		assert.Equal(t, s.GetPort(), defaultPort)
		assert.Equal(t, s.GetStorageHost(), "")
	})
}

func TestCfg_GetPort(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_PORT, "11124")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetPort(), int32(11124))
	})
}

func TestCfg_GetStorageHost(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_STORAGE_HOST, "localhost:9090")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetStorageHost(), "localhost:9090")
	})
}

func TestCfg_GetMongoURI(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_MONGO_URI, "11124")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetMongoURI(), "11124")
	})
}

func TestCfg_GetMongoDb(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_MONGO_DB, "11124")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetMongoDb(), "11124")
	})
}

func TestCfg_GetMongoCollection(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_MONGO_COLLECTION, "11124")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetMongoCollection(), "11124")
	})
}
