package job

import (
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCassandraJob_Do(t *testing.T) {
	t.Run("Test: cassandraJob", func(t *testing.T) {
		t.Run("Should: return error connecting", func(t *testing.T) {
			j := NewCassandraJob("", "", "")
			err := j.Do()
			expected := clientPb.StatusCode_Error
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
	})
}
