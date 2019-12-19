package job

import (
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgresJob_Do(t *testing.T) {
	t.Run("Test: postgresDbJob", func(t *testing.T) {
		t.Run("Should: return error connecting", func(t *testing.T) {
			j := NewPosgresDbJob("", 0, "", "", "")
			err := j.Do()
			expected := clientPb.StatusCode_Error
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return no error", func(t *testing.T) {
			j := postgresJob{
				mysql:    &sqlMockConnectOk{},
			}
			err := j.Do()
			expected := clientPb.StatusCode_OK
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
	})
}
