package parsers

import (
	"encoding/xml"
)

type SiteMap struct {
	XMLName xml.Name `xml:"urlset"`
	UrlSet []SiteMapUrl `xml:"url"`
}

type SiteMapUrl struct {
	XMLName  xml.Name `xml:"url"`
	Location string `xml:"loc"`
	Ignore bool `xml:"ignore"`
}

type SiteMapParser struct {
}

func NewSiteMapParser() *SiteMapParser {
	return &SiteMapParser{}
}

func (parser *SiteMapParser) Parse(xmlBytes []byte) (*SiteMap, error) {
	siteMap := &SiteMap{}
	err := xml.Unmarshal(xmlBytes, siteMap)
	if err != nil {
		return nil, err
	}
	return siteMap, nil
}
