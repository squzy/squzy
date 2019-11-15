package job

import "testing"

func TestJobHTTP_Do(t *testing.T) {
	t.Run("Test: JobHTTP.Do()", func(t *testing.T) {
		t.Run("Should: throw error incorrect ")
		j := jobHTTP{
			methodType: "",
			url:        "",
			headers:    nil,
			statusCode: 0,
		}
	})
}