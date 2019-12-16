package job

import (
	"errors"
	storagePb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"
)

type httpToolsMock struct {
}

func (h httpToolsMock) GetWithRedirectsWithStatusCode(url string, expectedCode int) (int, []byte, error) {
	panic("implement me")
}

func (h httpToolsMock) GetWithRedirects(url string) (int, []byte, error) {
	panic("implement me")
}

func (h httpToolsMock) CreateRequest(method string, url string, headers *map[string]string) *fasthttp.Request {
	return fasthttp.AcquireRequest()
}

type httpToolsMockError struct {
}

func (h httpToolsMockError) GetWithRedirectsWithStatusCode(url string, expectedCode int) (int, []byte, error) {
	panic("implement me")
}

func (h httpToolsMockError) GetWithRedirects(url string) (int, []byte, error) {
	panic("implement me")
}

func (h httpToolsMockError) CreateRequest(method string, url string, headers *map[string]string) *fasthttp.Request {
	return fasthttp.AcquireRequest()
}

func (h httpToolsMockError) SendRequest(req *fasthttp.Request) (int, []byte, error) {
	panic("implement me")
}

func (h httpToolsMockError) SendRequestWithStatusCode(req *fasthttp.Request, expectedCode int) (int, []byte, error) {
	return 0, nil, errors.New("safsaf")
}

func (h httpToolsMock) SendRequest(req *fasthttp.Request) (int, []byte, error) {
	panic("implement me")
}

func (h httpToolsMock) SendRequestWithStatusCode(req *fasthttp.Request, expectedCode int) (int, []byte, error) {
	return 0, nil, nil
}

func TestNewHttpJob(t *testing.T) {
	t.Run("Should: implement interface", func(t *testing.T) {
		s := NewHttpJob(http.MethodGet, "", map[string]string{}, http.StatusOK, &httpToolsMock{})
		assert.Implements(t, (*Job)(nil), s)
	})

}

func TestJobHTTP_Do(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := NewHttpJob(http.MethodGet, "", map[string]string{}, http.StatusOK, &httpToolsMock{})
		assert.Equal(t, storagePb.StatusCode_OK, s.Do().GetLogData().Code)
	})
	t.Run("Should: not return error with headers", func(t *testing.T) {
		s := NewHttpJob(http.MethodGet, "", map[string]string{
			"test": "asf",
		}, http.StatusOK, &httpToolsMock{})
		assert.Equal(t, storagePb.StatusCode_OK, s.Do().GetLogData().Code)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := NewHttpJob(http.MethodGet, "", map[string]string{}, http.StatusOK, &httpToolsMockError{})
		assert.Equal(t, storagePb.StatusCode_Error, s.Do().GetLogData().Code)
	})
	t.Run("Should: return error port 80", func(t *testing.T) {
		s := NewHttpJob(http.MethodGet, "http://google.ru", map[string]string{}, http.StatusOK, &httpToolsMockError{})
		assert.Equal(t, int32(80), s.Do().GetLogData().Meta.Port)
	})
	t.Run("Should: return error port 80", func(t *testing.T) {
		s := NewHttpJob(http.MethodGet, "https://google.ru", map[string]string{}, http.StatusOK, &httpToolsMockError{})
		assert.Equal(t, int32(443), s.Do().GetLogData().Meta.Port)
	})
}
