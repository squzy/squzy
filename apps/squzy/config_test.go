package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestReadConfig(t *testing.T) {
	t.Run("Should: read port from env", func(t *testing.T) {
		os.Setenv("PORT", "11124")
		ReadConfig()
		assert.Equal(t, port, int32(11124))
	})
	t.Run("Should: read timeout from env", func(t *testing.T) {
		os.Setenv("STORAGE_TIMEOUT", "11")
		ReadConfig()
		assert.Equal(t, timeoutStorage, time.Second * 11)
	})
	t.Run("Should: read client address from env", func(t *testing.T) {
		os.Setenv("STORAGE_HOST", "11")
		ReadConfig()
		assert.Equal(t, clientAddress, "11")
	})
}