package version

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion(t *testing.T) {
	t.Run("Should: return not nil", func(t *testing.T) {
		assert.NotNil(t, GetVersion())
	})
}
