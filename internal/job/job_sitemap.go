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
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
	"squzy/internal/semaphore"
	sitemap_storage "squzy/internal/sitemap-storage"
)

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
		Code:        s.code,
		Error:       err,
		Type:        apiPb.SchedulerType_SiteMap,
		Meta: &apiPb.SchedulerResponse_MetaData{
			StartTime: s.startTime,
			EndTime:   s.endTime,
		},
	}
}

func newSiteMapError(schedulerId string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerResponseCode, description string, location string) CheckError {
	return &siteMapError{
		schedulerId: schedulerId,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		location:    location,
	}
}

func ExecSiteMap(schedulerId string, timeout int32, config *scheduler_config_storage.SiteMapConfig, siteMapStorage sitemap_storage.SiteMapStorage, httpTools httpTools.HttpTool, semaphoreFactoryFn func(n int) semaphore.Semaphore) CheckError {
	startTime := ptypes.TimestampNow()
	siteMap, err := siteMapStorage.Get(config.Url)
	if err != nil {
		return newSiteMapError(schedulerId, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_Error, err.Error(), config.Url)
	}

	count := len(siteMap.UrlSet)

	if count == 0 {
		return newSiteMapError(schedulerId, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_OK, "", "")
	}

	concurrency := int(config.Concurrency)

	if concurrency <= 0 || concurrency > count {
		concurrency = len(siteMap.UrlSet)
	}

	sem := semaphoreFactoryFn(concurrency)

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

			rq := httpTools.CreateRequest(http.MethodGet, location, nil, schedulerId)
			_, _, err = httpTools.SendRequestTimeoutStatusCode(rq, helpers.DurationFromSecond(timeout), http.StatusOK)

			if err != nil {
				return err
			}

			return nil
		})
	}
	err = group.Wait()
	if err != nil {
		return newSiteMapError(schedulerId, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_Error, err.Error(), config.Url)
	}
	return newSiteMapError(schedulerId, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_OK, "", "")
}
