package job

import (
	"context"
	"errors"
	"fmt"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	"squzy/internal/parsers"
	"squzy/internal/semaphore"
	"testing"
	"time"
)

type mockHttpTools struct {
}

func (m mockHttpTools) SendRequestTimeoutStatusCode(req *http.Request, timeout time.Duration, expectedCode int, ) (int, []byte, error) {
	return 200, nil, nil
}

func (m mockHttpTools) SendRequestTimeout(req *http.Request, timeout time.Duration) (int, []byte, error) {
	return 200, nil, nil
}

func (m mockHttpTools) CreateRequest(method string, url string, headers *map[string]string, log string) *http.Request {
	rq, _ := http.NewRequest(method, url, nil)
	return rq
}

func (m mockHttpTools) SendRequest(req *http.Request) (int, []byte, error) {
	return 200, nil, nil
}

func (m mockHttpTools) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	return 200, nil, nil
}

type siteMapStorage struct {
}

type siteMapStorageIgnore struct {
}

func (s siteMapStorageIgnore) Get(url string) (*parsers.SiteMap, error) {
	return &parsers.SiteMap{
		UrlSet: []parsers.SiteMapUrl{
			{
				Location: "localhost",
				Ignore:   true,
			},
			{
				Location: "localhost",
				Ignore:   true,
			},
		},
	}, nil
}

func (s siteMapStorage) Get(url string) (*parsers.SiteMap, error) {
	return &parsers.SiteMap{
		UrlSet: []parsers.SiteMapUrl{
			{
				Location: "localhost",
				Ignore:   false,
			},
			{
				Location: "localhost",
				Ignore:   false,
			},
		},
	}, nil
}

type siteMapStorageError struct {
}

func (s siteMapStorageError) Get(url string) (*parsers.SiteMap, error) {
	return nil, errors.New("SAFafs")
}

type siteMapStorageEmptyIgnore struct {
}

func (s siteMapStorageEmptyIgnore) Get(url string) (*parsers.SiteMap, error) {
	return &parsers.SiteMap{
		UrlSet: []parsers.SiteMapUrl{
		},
	}, nil
}

type mockHttpToolsWithError struct {
}

func (m mockHttpToolsWithError) SendRequestTimeoutStatusCode(req *http.Request, timeout time.Duration, expectedCode int, ) (int, []byte, error) {
	return 500, nil, errors.New("Wrong code")
}

func (m mockHttpToolsWithError) SendRequestTimeout(req *http.Request, timeout time.Duration) (int, []byte, error) {
	panic("implement me")
}

func (m mockHttpToolsWithError) GetWithRedirectsWithStatusCode(url string, expectedCode int) (int, []byte, error) {
	return 500, nil, errors.New("Wrong code")
}

func (m mockHttpToolsWithError) GetWithRedirects(url string) (int, []byte, error) {
	return 500, nil, errors.New("Wrong code")
}

func (m mockHttpToolsWithError) CreateRequest(method string, url string, headers *map[string]string, log string) *http.Request {
	rq, _ := http.NewRequest(method, url, nil)
	return rq
}

func (m mockHttpToolsWithError) SendRequest(req *http.Request) (int, []byte, error) {
	return 500, nil, errors.New("Wrong code")
}

func (m mockHttpToolsWithError) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	return 500, nil, errors.New("Wrong code")
}

func TestNewSiteMapJob(t *testing.T) {
	t.Run("Should: Should implement interface Job", func(t *testing.T) {
		job := NewSiteMapJob("", 0, &siteMapStorage{}, &mockHttpTools{}, semaphore.NewSemaphore, 5)
		assert.Implements(t, (*Job)(nil), job)
	})
}

type mockErrorSemaphore struct {
}

func (*mockErrorSemaphore) Acquire(ctx context.Context) error {
	fmt.Println("saf")
	return errors.New("Acquire error")
}
func (*mockErrorSemaphore) Release() {}

func errorFactory(i int) semaphore.Semaphore {
	return &mockErrorSemaphore{}
}

func successFactory(i int) semaphore.Semaphore {
	return semaphore.NewSemaphore(i)
}

func TestSiteMapJob_Do(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		t.Run("Because mock with 200", func(t *testing.T) {
			job := NewSiteMapJob("", 0, &siteMapStorage{}, &mockHttpTools{}, successFactory, -1)
			assert.Equal(t, clientPb.StatusCode_OK, job.Do().GetLogData().Code)
		})
		t.Run("Because ignore url", func(t *testing.T) {
			job := NewSiteMapJob("", 0, &siteMapStorageIgnore{}, &mockHttpToolsWithError{}, successFactory, 5)
			assert.Equal(t, clientPb.StatusCode_OK, job.Do().GetLogData().Code)
		})
		t.Run("Because: empty sitemap", func(t *testing.T) {
			job := NewSiteMapJob("", 0, &siteMapStorageEmptyIgnore{}, &mockHttpToolsWithError{}, successFactory, 5)
			assert.Equal(t, clientPb.StatusCode_OK, job.Do().GetLogData().Code)
		})
	})
	t.Run("Should: return error", func(t *testing.T) {
		t.Run("Because Acquire error", func(t *testing.T) {
			job := NewSiteMapJob("", 0, &siteMapStorage{}, &mockHttpToolsWithError{}, errorFactory, 5)
			assert.IsType(t, clientPb.StatusCode_Error, job.Do().GetLogData().Code)
		})
		t.Run("Because return 500", func(t *testing.T) {
			job := NewSiteMapJob("", 0, &siteMapStorage{}, &mockHttpToolsWithError{}, successFactory, 5)
			assert.IsType(t, clientPb.StatusCode_Error, job.Do().GetLogData().Code)
		})
		t.Run("Because sitemapError", func(t *testing.T) {
			job := NewSiteMapJob("", 0, &siteMapStorageError{}, &mockHttpTools{}, successFactory, 5)
			assert.IsType(t, clientPb.StatusCode_Error, job.Do().GetLogData().Code)
		})
	})
}
