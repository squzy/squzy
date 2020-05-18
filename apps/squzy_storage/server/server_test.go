package server

import (
	"errors"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"os"
	"squzy/apps/squzy_storage/config"
	"testing"
)

func TestNewServer(t *testing.T) {
	t.Run("Should: work", func(t *testing.T) {
		s := NewServer(nil, nil)
		assert.NotNil(t, s)
	})
}

func TestServer_Run(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := &server{
			config: nil,
			newServer: func(config.Config) (apiPb.StorageServer, error) {
				return nil, errors.New("Expected error")
			},
		}
		assert.Error(t, s.Run())
	})
	t.Run("Should: return error", func(t *testing.T) {
		err := os.Setenv("PORT", "1000000")
		if err != nil {
			panic("Error writing os env")
		}
		s := &server{
			config:    config.New(),
			newServer: func(config.Config) (apiPb.StorageServer, error) { return nil, nil },
		}
		assert.Error(t, s.Run())
	})
	t.Run("Should: return error", func(t *testing.T) {
		err := os.Setenv("PORT", "1000000")
		if err != nil {
			panic("Error writing os env")
		}
		s := &server{
			config:    config.New(),
			newServer: func(config.Config) (apiPb.StorageServer, error) { return nil, nil },
		}
		assert.Error(t, s.Run())
	})
}
