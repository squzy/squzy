package job

import (
	"errors"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	"squzy/apps/internal/parsers"
	"testing"
)

type mockHttpTools struct {
	
}

func (m mockHttpTools) CreateRequest(method string, url string, headers *map[string]string) *http.Request {
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
				Ignore:true,
			},
			{
				Location: "localhost",
				Ignore:true,
			},
		},
	}, nil
}

func (s siteMapStorage) Get(url string) (*parsers.SiteMap, error) {
	return &parsers.SiteMap{
		UrlSet: []parsers.SiteMapUrl{
			{
				Location: "localhost",
				Ignore:false,
			},
			{
				Location: "localhost",
				Ignore:false,
			},
		},
	}, nil
}

type siteMapStorageError struct {

}

func (s siteMapStorageError) Get(url string) (*parsers.SiteMap, error) {
	return nil, errors.New("SAFafs")
}

type mockHttpToolsWithError struct {

}

func (m mockHttpToolsWithError) GetWithRedirectsWithStatusCode(url string, expectedCode int) (int, []byte, error) {
	return 500, nil, errors.New("Wrong code")
}

func (m mockHttpToolsWithError) GetWithRedirects(url string) (int, []byte, error) {
	return 500, nil, errors.New("Wrong code")
}

func (m mockHttpToolsWithError) CreateRequest(method string, url string, headers *map[string]string) *http.Request {
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
		job := NewSiteMapJob("", &siteMapStorage{}, &mockHttpTools{}, 5)
		assert.Implements(t, (*Job)(nil), job)
	})
}

func TestSiteMapJob_Do(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		t.Run("Because mock with 200", func(t *testing.T) {
			job := NewSiteMapJob("", &siteMapStorage{}, &mockHttpTools{}, -1)
			assert.Equal(t, clientPb.StatusCode_OK, job.Do().GetLogData().Code)
		})
		t.Run("Because ignore url", func(t *testing.T) {
			job := NewSiteMapJob("", &siteMapStorageIgnore{}, &mockHttpToolsWithError{}, 5)
			assert.Equal(t, clientPb.StatusCode_OK, job.Do().GetLogData().Code)
		})
	})
	t.Run("Should: return error", func(t *testing.T) {
		t.Run("Because return 500", func(t *testing.T) {
			job := NewSiteMapJob("", &siteMapStorage{}, &mockHttpToolsWithError{}, 5)
			assert.IsType(t, clientPb.StatusCode_Error, job.Do().GetLogData().Code)
		})
		t.Run("Because sitemapError", func(t *testing.T) {
			job := NewSiteMapJob("", &siteMapStorageError{}, &mockHttpTools{}, 5)
			assert.IsType(t, clientPb.StatusCode_Error, job.Do().GetLogData().Code)
		})
	})
}