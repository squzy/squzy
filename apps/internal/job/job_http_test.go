package job

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	clientPb "github.com/squzy/squzy_generated/generated/logger"
)

func TestJobHTTP_Do(t *testing.T) {
	t.Run("Test: JobHTTP.Do()", func(t *testing.T) {
		expectStatus := http.StatusOK

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
			except := NewHttpError(
				ptypes.TimestampNow(),
				clientPb.StatusCode_Error,
				wrongStatusError.Error(),
				j.url,
			)
			assert.Equal(t, except.GetLogData().Code, err.GetLogData().Code)
			assert.Equal(t, except.GetLogData().Description, err.GetLogData().Description)
		})


		t.Run("Should: throw error incorrect ", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
			defer ts.Close()
			m := make(map[string]string)
			m["Accept"] = "application/json; charset=utf-8"
			j := NewJob("GET", ts.URL, m, 200)
			err := j.Do()
			except := NewHttpError(
				ptypes.TimestampNow(),
				clientPb.StatusCode_OK,
				"",
				j.url,
			)
			assert.Equal(t, except.GetLogData().Code, err.GetLogData().Code)
		})
	})
}