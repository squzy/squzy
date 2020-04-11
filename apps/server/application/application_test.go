package application

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("Should: create new application", func(t *testing.T) {
		a := New(nil)
		assert.NotEqual(t, nil, a)
	})
}

func TestApplication_Run(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		app := New(nil)
		assert.NotEqual(t, nil, app.Run(1244214))
	})
	t.Run("Should: not return error", func(t *testing.T) {
		app := New(nil, )
		go func() {
			_ = app.Run(11122)
		}()
		time.Sleep(time.Second)
		_, err := net.Dial("tcp", "localhost:11122")
		assert.Equal(t, nil, err)
	})
}