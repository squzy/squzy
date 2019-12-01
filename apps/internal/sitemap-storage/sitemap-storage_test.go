package sitemap_storage

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"squzy/apps/internal/parsers"
	"errors"
	"time"
	"testing"
)

type mockHttp struct {

}

func (m mockHttp) SendRequest(req *http.Request) (int, []byte, error) {
	return 200, nil, nil
}

func (m mockHttp) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	return 200, nil, nil
}

type mockSiteMapParser struct {

}

type mockSiteMapParserError struct {

}

func (m mockSiteMapParserError) Parse(xmlBytes []byte) (*parsers.SiteMap, error) {
	return nil, errors.New("asdsad")
}

func (m mockSiteMapParser) Parse(xmlBytes []byte) (*parsers.SiteMap, error) {
	return &parsers.SiteMap{}, nil
}

type mockHttpError struct {

}

func (m mockHttpError) SendRequest(req *http.Request) (int, []byte, error) {
	return 0, nil, errors.New("ascss")
}

func (m mockHttpError) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	return 0, nil, errors.New("ascss")
}

func TestNew(t *testing.T) {
	t.Run("Shoudle implement interface", func(t *testing.T) {
		s := New(time.Second, &mockHttp{}, &mockSiteMapParser{})
		assert.Implements(t, (*SiteMapStorage)(nil), s)
	})
}

func TestStorage_Get(t *testing.T) {
	t.Run("Should: return error because httpError", func(t *testing.T) {
		s := New(time.Second, &mockHttpError{}, &mockSiteMapParser{})
		_, err := s.Get("evrerver")
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because parseError", func(t *testing.T) {
		s := New(time.Second, &mockHttp{}, &mockSiteMapParserError{})
		_, err := s.Get("evrerver")
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return sitemap", func(t *testing.T) {
		s := New(time.Second, &mockHttp{}, &mockSiteMapParser{})
		sm, err := s.Get("evrerver")
		assert.Equal(t, nil, err)
		assert.NotEqual(t, sm, err)
	})
	t.Run("Should: return from cache", func(t *testing.T) {
		s := New(time.Minute, &mockHttp{}, &mockSiteMapParser{})
		sm, err := s.Get("evrerver")
		assert.Equal(t, nil, err)
		assert.NotEqual(t, sm, err)
		sm2, _ := s.Get("evrerver")
		assert.Equal(t, sm, sm2)
	})
}