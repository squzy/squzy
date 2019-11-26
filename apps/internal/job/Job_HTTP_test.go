package job

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJobHTTP_Do(t *testing.T) {
	t.Run("Test: JobHTTP.Do()", func(t *testing.T) {
		expectStatus := http.StatusOK

		/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
		s := &http.Server{
			ReadTimeout: 1 * time.Second,
			WriteTimeout: 10 * time.Second,
			Addr:":8080",
		}
		go s.ListenAndServe()*/

		t.Run("Should: error client.Do incorrect ", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(expectStatus)
			}))
			defer ts.Close()
			j := NewJob("", "", nil, 0)
			err := j.Do()
			assert.NotEqual(t, nil, err)
		})

		t.Run("Should: throw error incorrect ", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
			defer ts.Close()
			j := NewJob("GET", ts.URL, nil, 0)
			err := j.Do()
			assert.Equal(t, wrongStatusError, err)
		})
	})
}