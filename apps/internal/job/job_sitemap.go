package job

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"squzy/apps/internal/httpTools"
	"squzy/apps/internal/parsers"
)

type siteMapJob struct {
	siteMap   *parsers.SiteMap
	httpTools httpTools.HttpTool
}

type siteMapError struct {
	location   string
	statusCode int
}

func (sme *siteMapError) Error() string {
	return fmt.Sprintf("Location: %s was return %d", sme.location, sme.statusCode)
}

func NewSiteMapJob(siteMap *parsers.SiteMap, httpTools httpTools.HttpTool) Job {
	return &siteMapJob{
		siteMap:   siteMap,
		httpTools: httpTools,
	}
}

func (j *siteMapJob) Do() error {
	ctx, cancel := context.WithCancel(context.Background())
	group, _ := errgroup.WithContext(ctx)

	for _, v := range j.siteMap.UrlSet {
		if v.Ignore {
			continue
		}
		group.Go(func() error {
			location := v.Location
			req, err := http.NewRequest(http.MethodGet, v.Location, nil)
			if err != nil {
				cancel()
				return err
			}
			code, _, err := j.httpTools.SendRequestWithStatusCode(req, http.StatusOK)
			if err != nil {
				cancel()
				return newSiteMapError(location, code)
			}
			return err
		})
	}
	err := group.Wait()
	if err != nil {
		return err
	}
	cancel()
	return nil
}

func newSiteMapError(location string, code int) error {
	return &siteMapError{
		location:   location,
		statusCode: code,
	}
}
