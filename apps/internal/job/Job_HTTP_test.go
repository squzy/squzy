package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJobHTTP_Do(t *testing.T) {
	t.Run("Test: JobHTTP.Do()", func(t *testing.T) {
		t.Run("Should: throw error incorrect ", func(t *testing.T) {
			j := jobHTTP{
				methodType: "",
				url:        "",
				headers:    nil,
				statusCode: 0,
			}
			err := j.Do()
			assert.NotEqual(t, nil, err)
		})
	})
}