package application

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("Should: not be nil", func(t *testing.T) {
		s := New(nil)
		assert.NotEqual(t, nil, s)
	})
}

func TestApp_Run(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		app := New(nil)
		go func() {
			_ = app.Run(11101)
		}()
		time.Sleep(time.Second)
		_, err := net.Dial("tcp", "localhost:11101")
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		app := New(nil)
		assert.NotEqual(t, nil, app.Run(1231323))
	})
}