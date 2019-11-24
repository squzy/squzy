package parsers

import (
	"encoding/xml"
	"net/http"
	"squzy/apps/internal/httpTools"
)

type SiteMap struct {
	UrlSet []SiteMapUrl `xml:"urlset"`
}

type SiteMapUrl struct {
	Location string `xml:"loc"`
}

type SiteMapParser struct {
	httpTool httpTools.HttpTool
}

func NewSiteMapParser(httpTool httpTools.HttpTool) *SiteMapParser {
	return &SiteMapParser{
		httpTool: httpTool,
	}
}

func (parser *SiteMapParser) Parse(url string) (*SiteMap, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	xmlBytes, err := parser.httpTool.SendRequestWithStatusCode(req, http.StatusOK)
	if err != nil {
		return nil, err
	}
	siteMap := &SiteMap{}
	err = xml.Unmarshal(xmlBytes, siteMap)
	if err != nil {
		return nil, err
	}
	return siteMap, nil
}
