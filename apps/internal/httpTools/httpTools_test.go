package httpTools

import (
	"github.com/stretchr/testify/assert"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("Test: Create new", func(t *testing.T) {
		j := New("veriosn")
		assert.IsType(t, &httpTool{}, j)
		assert.NotEqual(t, nil, j)
	})
}

func newRequest(method string, url string, body io.Reader) *http.Request {
	rq, _ := http.NewRequest(method, url, nil)
	return rq
}

func TestHttpTool_SendRequest(t *testing.T) {
	t.Run("Test: Should not return error", func(t *testing.T) {
		bytes := []byte("Hello, client")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)

			_, _ = w.Write(bytes)
		}))
		defer ts.Close()
		j := New("")
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
		j := New("veriosn")
		req := newRequest(http.MethodGet, ts.URL, nil)
		_, _, err := j.SendRequest(req)
		assert.NotEqual(t, nil, err)
	})
	t.Run("Test: Should return error", func(t *testing.T) {
		j := New("veriosn")
		req := newRequest(http.MethodGet, "ts.URL", nil)
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
		j := New("")
		req := newRequest(http.MethodGet, ts.URL, nil)
		_, body, _ := j.SendRequestWithStatusCode(req, http.StatusOK)
		assert.Equal(t, body, bytes)
	})
	t.Run("Test: Should return error", func(t *testing.T) {
		j := New("")
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
		j := New("")
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
		j := New("veriosn")
		req := newRequest(http.MethodGet, ts.URL, nil)
		_, _, err := j.SendRequestWithStatusCode(req, 200)
		assert.Equal(t, notExpectedStatusCodeFn(ts.URL, 201, 200), err)
	})
}

func TestHttpTool_CreateRequest(t *testing.T) {
	t.Run("Should: create request with header, url and method", func(t *testing.T) {
		h := New("veriosn")
		m := map[string]string{
			"trata": "trata",
		}
		rq := h.CreateRequest(http.MethodGet, "http://test.ru", &m, "23")
		assert.Equal(t, "http://test.ru", rq.URL.String())
		assert.Equal(t, http.MethodGet, rq.Method)
	})
	t.Run("Should: create request without headers", func(t *testing.T) {
		h := New("veriosn")
		rq := h.CreateRequest(http.MethodGet, "http://test.ru", nil, "")
		assert.Equal(t, "http://test.ru", rq.URL.String())
		assert.Equal(t, http.MethodGet, rq.Method)
	})
	t.Run("Should: create log header if present", func(t *testing.T) {
		h := New("veriosn")
		rq := h.CreateRequest(http.MethodGet, "http://test.ru", nil, "12")
		assert.Equal(t, "http://test.ru", rq.URL.String())
		assert.Equal(t, rq.Header.Get(logHeader), "12")
	})
	t.Run("Should: create with user agent", func(t *testing.T) {
		h := New("version")
		rq := h.CreateRequest(http.MethodGet, "http://test.ru", nil, "12")
		assert.Equal(t, "http://test.ru", rq.URL.String())
		assert.Equal(t, rq.Header.Get(userAgentHeaderKey), getUserAgent("version"))
	})
}

func TestHttpTool_SendRequestTimeout(t *testing.T) {
	t.Run("Should: return error because more then timeout", func(t *testing.T) {
		bytes := []byte("Hello, client")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Second * 2)
			w.WriteHeader(201)
			_, _ = w.Write(bytes)
		}))
		j := New("veriosn")
		req := newRequest(http.MethodGet, ts.URL, nil)
		s, _, err := j.SendRequestTimeout(req, time.Second)
		assert.NotEqual(t, nil, err)
		assert.Equal(t, 0, s)
	})
	t.Run("Should: return success because lest then timeout", func(t *testing.T) {
		bytes := []byte("Hello, client")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Second * 2)
			w.WriteHeader(201)
			_, _ = w.Write(bytes)
		}))
		j := New("veriosn")
		req := newRequest(http.MethodGet, ts.URL, nil)
		s, _, err := j.SendRequestTimeout(req, time.Second*4)
		assert.Equal(t, nil, err)
		assert.Equal(t, 201, s)
	})
	t.Run("Should: return success because use default timeout", func(t *testing.T) {
		bytes := []byte("Hello, client")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Second * 2)
			w.WriteHeader(201)
			_, _ = w.Write(bytes)
		}))
		j := New("veriosn")
		req := newRequest(http.MethodGet, ts.URL, nil)
		s, _, err := j.SendRequestTimeout(req, time.Second*(-4))
		assert.Equal(t, nil, err)
		assert.Equal(t, 201, s)
	})
}

func TestHttpTool_SendRequestTimeoutStatusCode(t *testing.T) {
	t.Run("Test: Should not return error", func(t *testing.T) {
		bytes := []byte("Hello, client")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Second)
			w.WriteHeader(200)
			_, _ = w.Write(bytes)
		}))
		defer ts.Close()
		j := New("")
		req := newRequest(http.MethodGet, ts.URL, nil)
		_, body, _ := j.SendRequestTimeoutStatusCode(req, time.Second*2, http.StatusOK)
		assert.Equal(t, body, bytes)
	})

	t.Run("Test: Should return error because timeout", func(t *testing.T) {
		bytes := []byte("Hello, client")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Second * 3)
			w.WriteHeader(200)
			_, _ = w.Write(bytes)
		}))
		defer ts.Close()
		j := New("")
		req := newRequest(http.MethodGet, ts.URL, nil)
		_, body, _ := j.SendRequestTimeoutStatusCode(req, time.Second*1, http.StatusOK)
		assert.Equal(t, []uint8([]byte(nil)), body)
	})
}
