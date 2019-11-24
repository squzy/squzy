package job

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestJobHTTP_Do(t *testing.T) {
	t.Run("Test: JobHTTP.Do()", func(t *testing.T) {
		expectStatus := http.StatusOK

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(expectStatus)
			}
		})
		http.HandleFunc("/nilResp", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
			}
		})
		go http.ListenAndServe(":8080", nil)

		t.Run("Should: error client.Do incorrect ", func(t *testing.T) {
			j := NewJob("", "", nil, 0)
			err := j.Do()
			assert.NotEqual(t, nil, err)
		})

		t.Run("Should: throw error incorrect ", func(t *testing.T) {
			j := NewJob("GET", "http://localhost:8080/nilResp", nil, 0)
			err := j.Do()
			assert.Equal(t, errorNoResponse, err.Error())
		})

		t.Run("Should: throw error incorrect ", func(t *testing.T) {
			j := NewJob("GET", "http://localhost:8080/", nil, 0)
			err := j.Do()
			assert.Equal(t, errorWrongStatus, err.Error())
		})

		t.Run("Should: work without error ", func(t *testing.T) {
			j := NewJob("GET", "http://localhost:8080/", nil, expectStatus)
			err := j.Do()
			assert.Equal(t, nil, err)
		})
	})
}