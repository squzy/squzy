package job

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/logger"
	"golang.org/x/sync/errgroup"
	"net/http"
	"squzy/apps/internal/httpTools"
	sitemap_storage "squzy/apps/internal/sitemap-storage"
	"strings"
)

type siteMapJob struct {
	url            string
	siteMapStorage sitemap_storage.SiteMapStorage
	httpTools      httpTools.HttpTool
}

type siteMapError struct {
	time        *timestamp.Timestamp
	code        clientPb.StatusCode
	description string
	location    string
	port        int32
}

func (s *siteMapError) GetLogData() *clientPb.Log {
	return &clientPb.Log{
		Code:        s.code,
		Description: s.description,
		Meta: &clientPb.MetaData{
			Id:       uuid.New().String(),
			Location: s.location,
			Port:     s.port,
			Time:     s.time,
			Type:     clientPb.Type_SiteMap,
		},
	}
}

func newSiteMapError(time *timestamp.Timestamp, code clientPb.StatusCode, description string, location string, port int32) CheckError {
	return &siteMapError{
		time:        time,
		code:        code,
		description: description,
		location:    location,
		port:        port,
	}
}

type siteMapErr struct {
	location   string
	statusCode int
}

func (sme *siteMapErr) Error() string {
	return fmt.Sprintf("%s -  was return %d", sme.location, sme.statusCode)
}

func NewSiteMapJob(url string, siteMapStorage sitemap_storage.SiteMapStorage, httpTools httpTools.HttpTool) Job {
	return &siteMapJob{
		url:            url,
		siteMapStorage: siteMapStorage,
		httpTools:      httpTools,
	}
}

func (j *siteMapJob) Do() CheckError {
	ctx, cancel := context.WithCancel(context.Background())
	group, _ := errgroup.WithContext(ctx)
	siteMap, err := j.siteMapStorage.Get(j.url)
	if err != nil {
		return newSiteMapError(ptypes.TimestampNow(), clientPb.StatusCode_Error, err.Error(), j.url, 0)
	}
	for _, v := range siteMap.UrlSet {
		if v.Ignore {
			continue
		}
		location := v.Location
		group.Go(func() error {
			req, _ := http.NewRequest(http.MethodGet, location, nil)
			code, _, err := j.httpTools.SendRequestWithStatusCode(req, http.StatusOK)
			if err != nil {
				cancel()
				return newSiteMapErr(location, code)
			}
			return nil
		})
	}
	err = group.Wait()
	if err != nil {
		location := strings.Split(err.Error(), " - ")
		return newSiteMapError(ptypes.TimestampNow(), clientPb.StatusCode_Error, err.Error(), location[0], 80)
	}
	cancel()
	return newSiteMapError(ptypes.TimestampNow(), clientPb.StatusCode_OK, "", "", 0)
}

func newSiteMapErr(location string, code int) error {
	return &siteMapErr{
		location:   location,
		statusCode: code,
	}
}
