package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: create new without error", func(t *testing.T) {
		s := New()
		assert.NotNil(t, s)
	})
}

func TestCfg_GetStorageHost(t *testing.T) {
	t.Run("Should: get storage host", func(t *testing.T) {
		err := os.Setenv(ENV_STORAGE_HOST, "localhost:9090")
		assert.Nil(t, err)
		s := New()
		assert.Equal(t, s.GetStorageHost(), "localhost:9090")
	})
}

func TestCfg_GetMongoURI(t *testing.T) {
	t.Run("Should: get mongo host", func(t *testing.T) {
		err := os.Setenv(ENV_MONGO_URI, "localhost:9090")
		assert.Nil(t, err)
		s := New()
		assert.Equal(t, s.GetMongoURI(), "localhost:9090")
	})
}

func TestCfg_GetPort(t *testing.T) {
	t.Run("Should: get port", func(t *testing.T) {
		err := os.Setenv(ENV_PORT, "9090")
		assert.Nil(t, err)
		s := New()
		assert.Equal(t, s.GetPort(), int32(9090))
	})
}

func TestCfg_GetMongoDB(t *testing.T) {
	t.Run("Should: get mongo db", func(t *testing.T) {
		err := os.Setenv(ENV_MONGO_DB, "localhost:9090")
		assert.Nil(t, err)
		s := New()
		assert.Equal(t, s.GetMongoDB(), "localhost:9090")
	})
}

func TestCfg_GetNotificationListCollection(t *testing.T) {
	t.Run("Should: get mongo collection", func(t *testing.T) {
		err := os.Setenv(ENV_MONGO_LIST_COLLECTION, "localhost:9090")
		assert.Nil(t, err)
		s := New()
		assert.Equal(t, s.GetNotificationListCollection(), "localhost:9090")
	})
}

func TestCfg_GetNotificationMethodCollection(t *testing.T) {
	t.Run("Should: get mongo collection", func(t *testing.T) {
		err := os.Setenv(ENV_MONGO_METHOD_COLLECTION, "localhost:9090")
		assert.Nil(t, err)
		s := New()
		assert.Equal(t, s.GetNotificationMethodCollection(), "localhost:9090")
	})
}