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
	"squzy/internal/httptools"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
	"squzy/internal/semaphore"
	sitemap_storage "squzy/internal/sitemap-storage"
)

type siteMapError struct {
	schedulerID string
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
			Message: fmt.Sprintf("Error: %s, URL: %s", s.description, s.location),
		}
	}
	return &apiPb.SchedulerResponse{
		SchedulerId: s.schedulerID,
		Code:        s.code,
		Error:       err,
		Type:        apiPb.SchedulerType_SiteMap,
		Meta: &apiPb.SchedulerResponse_MetaData{
			StartTime: s.startTime,
			EndTime:   s.endTime,
		},
	}
}

func newSiteMapError(schedulerID string, startTime *timestamp.Timestamp, endTime *timestamp.Timestamp, code apiPb.SchedulerResponseCode, description string, location string) CheckError {
	return &siteMapError{
		schedulerID: schedulerID,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		location:    location,
	}
}

func ExecSiteMap(schedulerID string, timeout int32, config *scheduler_config_storage.SiteMapConfig, siteMapStorage sitemap_storage.SiteMapStorage, httpTools httptools.HTTPTool, semaphoreFactoryFn func(n int) semaphore.Semaphore) CheckError {
	startTime := ptypes.TimestampNow()
	siteMap, err := siteMapStorage.Get(config.URL)
	if err != nil {
		return newSiteMapError(schedulerID, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_Error, err.Error(), config.URL)
	}

	count := len(siteMap.URLSet)

	if count == 0 {
		return newSiteMapError(schedulerID, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_OK, "", "")
	}

	concurrency := int(config.Concurrency)

	if concurrency <= 0 || concurrency > count {
		concurrency = len(siteMap.URLSet)
	}

	sem := semaphoreFactoryFn(concurrency)

	group, ctx := errgroup.WithContext(context.Background())
	for _, v := range siteMap.URLSet {
		if v.Ignore {
			continue
		}
		location := v.Location

		group.Go(func() error {
			errSem := sem.Acquire(ctx)

			if errSem != nil {
				return err
			}

			defer sem.Release()

			rq := httpTools.CreateRequest(http.MethodGet, location, nil, schedulerID)
			_, _, errSem = httpTools.SendRequestTimeoutStatusCode(rq, helpers.DurationFromSecond(timeout), http.StatusOK)

			if errSem != nil {
				return err
			}

			return nil
		})
	}
	err = group.Wait()
	if err != nil {
		return newSiteMapError(schedulerID, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_Error, err.Error(), config.URL)
	}
	return newSiteMapError(schedulerID, startTime, ptypes.TimestampNow(), apiPb.SchedulerResponseCode_OK, "", "")
}
