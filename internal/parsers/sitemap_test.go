package parsers

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestNewSiteMapParser(t *testing.T) {
	t.Run("Test: SiteMapParse create", func(t *testing.T) {
		parser := NewSiteMapParser()
		assert.IsType(t, &siteMapParser{}, parser)
		assert.NotEqual(t, nil, parser)
	})
}

func TestSiteMapParser_Parse(t *testing.T) {
	t.Run("Test: Parse", func(t *testing.T) {
		t.Run("Should: parse without error", func(t *testing.T) {
			parser := NewSiteMapParser()
			f, errRead := ioutil.ReadFile("valid.xml")
			assert.NoError(t, errRead)
			res, err := parser.Parse(f)
			assert.Equal(t, nil, err)
			assert.Equal(t, 5, len(res.URLSet))
			assert.Equal(t, true, res.URLSet[0].Ignore)
		})
		t.Run("Should: parse with error", func(t *testing.T) {
			parser := NewSiteMapParser()
			f, errRead := ioutil.ReadFile("invalid.xml")
			assert.NoError(t, errRead)
			_, err := parser.Parse(f)
			assert.NotEqual(t, nil, err)
		})
	})
}
