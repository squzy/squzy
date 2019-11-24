package parsers

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestNewSiteMapParser(t *testing.T) {
	t.Run("Test: SiteMapParse create", func(t *testing.T) {
		parser := NewSiteMapParser()
		assert.IsType(t, &SiteMapParser{}, parser)
		assert.NotEqual(t, nil, parser)
	})
}

func TestSiteMapParser_Parse(t *testing.T) {
	t.Run("Test: Parse", func(t *testing.T) {
		t.Run("Should: parse without error", func(t *testing.T) {
			parser := NewSiteMapParser()
			f, _ := ioutil.ReadFile("valid.xml")
			res, err := parser.Parse(f)
			assert.Equal(t, nil, err)
			assert.Equal(t, 5, len(res.UrlSet))
			assert.Equal(t, true, res.UrlSet[0].Ignore)
		})
		t.Run("Should: parse with error", func(t *testing.T) {
			parser := NewSiteMapParser()
			f, _ := ioutil.ReadFile("invalid.xml")
			_, err := parser.Parse(f)
			assert.NotEqual(t, nil, err)
		})
	})
}