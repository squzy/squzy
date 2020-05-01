package mongo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: implement interface", func(t *testing.T) {
		conn := New(nil)
		assert.Implements(t, (*Connector)(nil), conn)
	})
}