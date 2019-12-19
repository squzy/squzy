package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPortByUrl(t *testing.T) {
	t.Run("Should: return int32(80)", func(t *testing.T) {
		assert.Equal(t, int32(80), GetPortByUrl("http://google.com"))
	})
	t.Run("Should: return int32(443)", func(t *testing.T) {
		assert.Equal(t, int32(443), GetPortByUrl("https://google.com"))
	})
}