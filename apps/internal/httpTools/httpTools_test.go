package httpTools

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Test: Create new", func(t *testing.T) {
		j := New()
		assert.IsType(t, &httpTool{}, j)
		assert.NotEqual(t, nil, j)
	})
}

func TestHttpTool_SendRequest(t *testing.T) {
	t.Run("Test: Should not return error", func(t *testing.T) {
		bytes := []byte("Hello, client")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)

			_, _ = w.Write(bytes)
		}))
		defer ts.Close()
		j := New()
		req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
		code, body, _ := j.SendRequest(req)
		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, body, bytes)
	})
	t.Run("Test: Should return error", func(t *testing.T) {
		j := New()
		req, _ := http.NewRequest(http.MethodGet, "ts.URL", nil)
		_, _, err := j.SendRequest(req)
		assert.NotEqual(t, nil, err)
	})
}

func TestHttpTool_SendRequestWithStatusCode(t *testing.T) {
	t.Run("Test: Should not return error", func(t *testing.T) {
		bytes := []byte("Hello, client")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)

			_, _ = w.Write(bytes)
		}))
		defer ts.Close()
		j := New()
		req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
		_, body, _ := j.SendRequestWithStatusCode(req, http.StatusOK)
		assert.Equal(t, body, bytes)
	})
	t.Run("Test: Should return error", func(t *testing.T) {
		j := New()
		req, _ := http.NewRequest(http.MethodGet, "ts.URL", nil)
		_, _, err := j.SendRequestWithStatusCode(req, 200)
		assert.NotEqual(t, nil, err)
	})

	t.Run("Test: Should return notExpectedStatusCode", func(t *testing.T) {
		bytes := []byte("Hello, client")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(201)

			_, _ = w.Write(bytes)
		}))
		j := New()
		req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
		_, _, err := j.SendRequestWithStatusCode(req, 200)
		assert.Equal(t, notExpectedStatusCode, err)
	})
}