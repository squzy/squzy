package sitemap_storage

import (
	"net/http"
	"squzy/apps/internal/httpTools"
	"squzy/apps/internal/parsers"
	"sync"
	"time"
)

type SiteMapStorage interface {
	Get(url string) (*parsers.SiteMap, error)
}

type storage struct {
	httpTools     httpTools.HttpTool
	duration      time.Duration
	kv            map[string]*StorageItem
	mutex         sync.RWMutex
	siteMapParser parsers.SiteMapParser
}

type StorageItem struct {
	deadline time.Time
	siteMap  *parsers.SiteMap
}

func (s storage) Get(url string) (*parsers.SiteMap, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	value, exist := s.kv[url]
	if exist && time.Now().Before(value.deadline) {
		return value.siteMap, nil
	}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	_, resp, err := s.httpTools.SendRequestWithStatusCode(req, http.StatusOK)
	if err != nil {
		return nil, err
	}
	siteMap, err := s.siteMapParser.Parse(resp)
	if err != nil {
		return nil, err
	}
	s.kv[url] = &StorageItem{
		deadline: time.Now().Add(s.duration),
		siteMap:  siteMap,
	}
	return siteMap, nil
}

func New(duration time.Duration, httpTools httpTools.HttpTool, siteMapParser parsers.SiteMapParser) SiteMapStorage {
	return &storage{
		duration:      duration,
		kv:            make(map[string]*StorageItem),
		httpTools:     httpTools,
		siteMapParser: siteMapParser,
	}
}
