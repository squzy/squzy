package httpTools

import (
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Test: Create new", func(t *testing.T) {
		j := New()
		assert.IsType(t, &httpTool{}, j)
		assert.NotEqual(t, nil, j)
	})
}

func newRequest(method string, url string, body io.Reader) *fasthttp.Request {
	rq := fasthttp.AcquireRequest()
	rq.SetRequestURI(url)
	rq.Header.SetMethod(method)
	return rq
}

func TestHttpTool_GetWithRedirects(t *testing.T) {
	t.Run("Test: Should not return error", func(t *testing.T) {
		bytes := []byte("Hello, client")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)

			_, _ = w.Write(bytes)
		}))
		defer ts.Close()
		j := New()
		code, body, _ := j.GetWithRedirects(ts.URL)
		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, body, bytes)
	})
	t.Run("Test: Should return error because of body", func(t *testing.T) {
		bytes := []byte(strings.Repeat("hello", math.MaxInt8))
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1")
			w.WriteHeader(200)
			_, _ = w.Write(bytes)
		}))
		j := New()
		_, _, err := j.GetWithRedirects(ts.URL)
		assert.NotEqual(t, nil, err)
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
		req := newRequest(http.MethodGet, ts.URL, nil)
		code, body, _ := j.SendRequest(req)
		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, body, bytes)
	})
	t.Run("Test: Should return error because of body", func(t *testing.T) {
		bytes := []byte(strings.Repeat("hello", math.MaxInt8))
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1")
			w.WriteHeader(200)
			_, _ = w.Write(bytes)
		}))
		j := New()
		req := newRequest(http.MethodGet, ts.URL, nil)
		_, _, err := j.SendRequest(req)
		assert.NotEqual(t, nil, err)
	})
	t.Run("Test: Should return error", func(t *testing.T) {
		j := New()
		req := newRequest (http.MethodGet, "ts.URL", nil)
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
		req := newRequest(http.MethodGet, ts.URL, nil)
		_, body, _ := j.SendRequestWithStatusCode(req, http.StatusOK)
		assert.Equal(t, body, bytes)
	})
	t.Run("Test: Should return error", func(t *testing.T) {
		j := New()
		req := newRequest(http.MethodGet, "ts.URL", nil)
		_, _, err := j.SendRequestWithStatusCode(req, 200)
		assert.NotEqual(t, nil, err)
	})

	t.Run("Test: Should return because of body", func(t *testing.T) {
		bytes := []byte(strings.Repeat("hello", math.MaxInt8))
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1")
			w.WriteHeader(200)
			_, _ = w.Write(bytes)
		}))
		j := New()
		req := newRequest(http.MethodGet, ts.URL, nil)
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
		req  := newRequest(http.MethodGet, ts.URL, nil)
		_, _, err := j.SendRequestWithStatusCode(req, 200)
		assert.Equal(t, notExpectedStatusCode, err)
	})
}

func TestHttpTool_CreateRequest(t *testing.T) {
	t.Run("Should: create request with header, url and method", func(t *testing.T) {
		h := New()
		m := map[string]string{
			"trata": "trata",
		}
		rq := h.CreateRequest(http.MethodGet, "http://test.ru", &m)
		assert.Equal(t,"http://test.ru/", string(rq.URI().FullURI()))
		assert.Equal(t, http.MethodGet, string(rq.Header.Method()))
	})
	t.Run("Should: create request without headers", func(t *testing.T) {
		h := New()
		rq := h.CreateRequest(http.MethodGet, "http://test.ru", nil)
		assert.Equal(t,"http://test.ru/", string(rq.URI().FullURI()))
		assert.Equal(t, http.MethodGet, string(rq.Header.Method()))
	})
}