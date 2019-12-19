package job

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
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
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
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
			Id:        uuid.New().String(),
			Location:  s.location,
			Port:      s.port,
			StartTime: s.startTime,
			EndTime:   s.endTime,
			Type:      clientPb.Type_SiteMap,
		},
	}
}

func newSiteMapError(startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code clientPb.StatusCode, description string, location string, port int32) CheckError {
	return &siteMapError{
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		location:    location,
		port:        port,
	}
}

type siteMapErr struct {
	location      string
	statusCode    int
	internalError error
}

func (sme *siteMapErr) Error() string {
	return fmt.Sprintf("StatusCode %d, fullError - %s", sme.statusCode, sme.internalError.Error())
}

func NewSiteMapJob(url string, siteMapStorage sitemap_storage.SiteMapStorage, httpTools httpTools.HttpTool) Job {
	return &siteMapJob{
		url:            url,
		siteMapStorage: siteMapStorage,
		httpTools:      httpTools,
	}
}

func (j *siteMapJob) Do() CheckError {
	startTime := ptypes.TimestampNow()
	siteMap, err := j.siteMapStorage.Get(j.url)
	if err != nil {
		return newSiteMapError(startTime, ptypes.TimestampNow(), clientPb.StatusCode_Error, err.Error(), j.url, 0)
	}
	ctx, cancel := context.WithCancel(context.Background()) //nolint
	group, _ := errgroup.WithContext(ctx)
	for _, v := range siteMap.UrlSet {
		if v.Ignore {
			continue
		}
		location := v.Location
		group.Go(func() error {
			rq, _ := http.NewRequest(http.MethodGet, location, nil)
			code, _, err := j.httpTools.SendRequestWithStatusCode(rq, http.StatusOK)
			if err != nil {
				cancel()
				return newSiteMapErr(location, code, err)
			}
			return nil
		})
	}
	err = group.Wait()
	if err != nil {
		location := strings.Split(err.Error(), " - ")
		var url string
		if len(location) > 0 {
			url = location[0]
		}
		return newSiteMapError(startTime, ptypes.TimestampNow(), clientPb.StatusCode_Error, err.Error(), url, GetPortByUrl(url)) //nolint
	}
	cancel()
	return newSiteMapError(startTime, ptypes.TimestampNow(), clientPb.StatusCode_OK, "", "", 0)
}

func newSiteMapErr(location string, code int, error error) error {
	return &siteMapErr{
		location:      location,
		statusCode:    code,
		internalError: error,
	}
}
