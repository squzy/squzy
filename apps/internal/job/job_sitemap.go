package job

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"golang.org/x/sync/errgroup"
	"net/http"
	"squzy/apps/internal/helpers"
	"squzy/apps/internal/httpTools"
	"squzy/apps/internal/semaphore"
	sitemap_storage "squzy/apps/internal/sitemap-storage"
)

type siteMapJob struct {
	url              string
	concurrency      int32
	timeout          int32
	siteMapStorage   sitemap_storage.SiteMapStorage
	httpTools        httpTools.HttpTool
	semaphoreFactory func(n int) semaphore.Semaphore
}

type siteMapError struct {
	logId       string
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
			Id:        s.logId,
			Location:  s.location,
			Port:      s.port,
			StartTime: s.startTime,
			EndTime:   s.endTime,
			Type:      clientPb.Type_SiteMap,
		},
	}
}

func newSiteMapError(logId string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code clientPb.StatusCode, description string, location string, port int32) CheckError {
	return &siteMapError{
		logId:       logId,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		location:    location,
		port:        port,
	}
}

func NewSiteMapJob(url string, timeout int32, siteMapStorage sitemap_storage.SiteMapStorage, httpTools httpTools.HttpTool, semaphoreFactoryFn func(n int) semaphore.Semaphore, concurrency int32) Job {
	return &siteMapJob{
		url:              url,
		concurrency:      concurrency,
		timeout:          timeout,
		siteMapStorage:   siteMapStorage,
		httpTools:        httpTools,
		semaphoreFactory: semaphoreFactoryFn,
	}
}

func (j *siteMapJob) Do() CheckError {
	logId := uuid.New().String()
	startTime := ptypes.TimestampNow()
	siteMap, err := j.siteMapStorage.Get(j.url)
	if err != nil {
		return newSiteMapError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_Error, err.Error(), j.url, 0)
	}

	count := len(siteMap.UrlSet)

	if count == 0 {
		return newSiteMapError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_OK, "", "", 0)
	}

	concurrency := int(j.concurrency)

	if concurrency <= 0 || concurrency > count {
		concurrency = len(siteMap.UrlSet)
	}

	sem := j.semaphoreFactory(concurrency)

	group, ctx := errgroup.WithContext(context.Background())
	for _, v := range siteMap.UrlSet {

		if v.Ignore {
			continue
		}
		location := v.Location

		group.Go(func() error {
			err := sem.Acquire(ctx)

			if err != nil {
				return err
			}

			defer sem.Release()

			rq := j.httpTools.CreateRequest(http.MethodGet, location, nil, logId)
			_, _, err = j.httpTools.SendRequestTimeoutStatusCode(rq, helpers.DurationFromSecond(j.timeout), http.StatusOK)

			if err != nil {
				return err
			}

			return nil
		})
	}
	err = group.Wait()
	if err != nil {
		return newSiteMapError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_Error, err.Error(), j.url, helpers.GetPortByUrl(j.url)) //nolint
	}
	return newSiteMapError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_OK, "", "", 0)
}
