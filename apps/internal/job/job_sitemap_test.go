package job

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"squzy/apps/internal/parsers"
	"errors"
	"testing"
)
type mockHttpTools struct {
	
}

func (m mockHttpTools) SendRequest(req *http.Request) (int, []byte, error) {
	return 200, nil, nil
}

func (m mockHttpTools) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	return 200, nil, nil
}

type mockHttpToolsWithError struct {

}

func (m mockHttpToolsWithError) SendRequest(req *http.Request) (int, []byte, error) {
	return 500, nil, errors.New("Wrong code")
}

func (m mockHttpToolsWithError) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	return 500, nil, errors.New("Wrong code")
}

func TestNewSiteMapJob(t *testing.T) {
	t.Run("Should: Should implement interface Job", func(t *testing.T) {
		job := NewSiteMapJob(nil, &mockHttpTools{})
		assert.Implements(t, (*Job)(nil), job)
	})
}

func TestSiteMapJob_Do(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		t.Run("Because mock with 200", func(t *testing.T) {
			job := NewSiteMapJob(&parsers.SiteMap{
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
			}, &mockHttpTools{})
			assert.Equal(t, nil, job.Do())
		})
		t.Run("Because ignore url", func(t *testing.T) {
			job := NewSiteMapJob(&parsers.SiteMap{
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
			}, &mockHttpToolsWithError{})
			assert.Equal(t, nil, job.Do())
		})
	})
	t.Run("Should: return error", func(t *testing.T) {
		t.Run("Because return 500", func(t *testing.T) {
			job := NewSiteMapJob(&parsers.SiteMap{
				UrlSet: []parsers.SiteMapUrl{
					{
						Location: "localhost",
						Ignore:false,
					},
					{
						Location: "localhost",
						Ignore:true,
					},
				},
			}, &mockHttpToolsWithError{})
			assert.IsType(t, &siteMapError{}, job.Do())
		})
	})
}