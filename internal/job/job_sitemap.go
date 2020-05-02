package job

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"golang.org/x/sync/errgroup"
	"net/http"
	"squzy/internal/helpers"
	"squzy/internal/httpTools"
	"squzy/internal/semaphore"
	sitemap_storage "squzy/internal/sitemap-storage"
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
	schedulerId string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerResponseCode
	description string
	location    string
}

func (s *siteMapError) GetLogData() *apiPb.SchedulerResponse {
	var err *apiPb.SchedulerResponse_Error
	if s.code == apiPb.SchedulerResponseCode_Error {
		err = &apiPb.SchedulerResponse_Error{
			Message: fmt.Sprintf("Error: %s, Url: %s", s.description, s.location),
		}
	}
	return &apiPb.SchedulerResponse{
		SchedulerId: s.schedulerId,
		Code:  s.code,
		Error: err,
		Type:  apiPb.SchedulerType_SiteMap,
		Meta: &apiPb.SchedulerResponse_MetaData{
			StartTime: s.startTime,
			EndTime:   s.endTime,
		},
	}
}

func newSiteMapError(schedulerId string,startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerResponseCode, description string, location string) CheckError {
	return &siteMapError{
		schedulerId:schedulerId,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		location:    location,
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

func (j *siteMapJob) Do(schedulerId string) CheckError {
	startTime := ptypes.TimestampNow()
	siteMap, err := j.siteMapStorage.Get(j.url)
	if err != nil {
		return newSiteMapError(schedulerId,startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_Error, err.Error(), j.url)
	}

	count := len(siteMap.UrlSet)

	if count == 0 {
		return newSiteMapError(schedulerId, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_OK, "", "")
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

			rq := j.httpTools.CreateRequest(http.MethodGet, location, nil, schedulerId)
			_, _, err = j.httpTools.SendRequestTimeoutStatusCode(rq, helpers.DurationFromSecond(j.timeout), http.StatusOK)

			if err != nil {
				return err
			}

			return nil
		})
	}
	err = group.Wait()
	if err != nil {
		return newSiteMapError(schedulerId, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_Error, err.Error(), j.url) //nolint
	}
	return newSiteMapError(schedulerId, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_OK, "", "")
}
