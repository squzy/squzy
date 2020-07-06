package expression

import (
	"context"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

type FilterSnapshot func(req *apiPb.GetSchedulerInformationRequest) *apiPb.GetSchedulerInformationRequest

func (e *expressionStruct) GetSnapshots(
	schedulerId string,
	direction apiPb.SortDirection,
	pagination *apiPb.Pagination,
	filters ...FilterSnapshot) []*apiPb.SchedulerSnapshot {

	req := &apiPb.GetSchedulerInformationRequest{
		SchedulerId: schedulerId,
		Pagination:  pagination,
		Sort: &apiPb.SortingSchedulerList{
			Direction: direction,
		},
	}
	if filters != nil {
		for _, filter := range filters {
			req = filter(req)
		}
	}
	list, err := e.storageClient.GetSchedulerInformation(context.Background(), req)
	if err != nil {
		panic(err)
	}
	return list.GetSnapshots()
}

func (e *expressionStruct) getSnapshotEnv(schedulerId string) map[string]interface{} {
	return map[string]interface{}{
		"Last": func(count int32, filters ...FilterSnapshot) []*apiPb.SchedulerSnapshot {
			return e.GetSnapshots(
				schedulerId,
				apiPb.SortDirection_ASC,
				&apiPb.Pagination{
					Page:  -1,
					Limit: count,
				},
				filters...)
		},
		"First": func(count int32, filters ...FilterSnapshot) []*apiPb.SchedulerSnapshot {
			return e.GetSnapshots(
				schedulerId,
				apiPb.SortDirection_ASC,
				&apiPb.Pagination{
					Page:  0,
					Limit: count,
				},
				filters...)
		},
		"Index": func(index int32, filters ...FilterSnapshot) []*apiPb.SchedulerSnapshot {
			return e.GetSnapshots(
				schedulerId,
				apiPb.SortDirection_ASC,
				&apiPb.Pagination{
					Page:  index,
					Limit: 1,
				},
				filters...)
		},
		"UseCode": func(status apiPb.SchedulerCode) FilterSnapshot {
			return func(req *apiPb.GetSchedulerInformationRequest) *apiPb.GetSchedulerInformationRequest {
				req.Status = status
				return req
			}
		},
		"UseTimeFrom": func(timeStr string) FilterSnapshot {
			return func(req *apiPb.GetSchedulerInformationRequest) *apiPb.GetSchedulerInformationRequest {
				if req.TimeRange == nil {
					req.TimeRange = &apiPb.TimeFilter{}
				}
				req.TimeRange.From = convertToTimestamp(timeStr)
				return req
			}
		},
		"UseTimeTo": func(timeStr string) FilterSnapshot {
			return func(req *apiPb.GetSchedulerInformationRequest) *apiPb.GetSchedulerInformationRequest {
				if req.TimeRange == nil {
					req.TimeRange = &apiPb.TimeFilter{}
				}
				req.TimeRange.To = convertToTimestamp(timeStr)
				return req
			}
		},
		"Duration": func(snapshot *apiPb.SchedulerSnapshot) int64 {
			return getTimeRange(snapshot.GetMeta().GetStartTime(), snapshot.GetMeta().GetEndTime())
		},
		//Transaction status keys
		"Ok":            apiPb.SchedulerCode_OK,
		"Error":         apiPb.SchedulerCode_ERROR,
		"TCP":           apiPb.SchedulerType_TCP,
		"GRPC":          apiPb.SchedulerType_GRPC,
		"HTTP":          apiPb.SchedulerType_HTTP,
		"SiteMap":       apiPb.SchedulerType_SITE_MAP,
		"HTTPJSONValue": apiPb.SchedulerType_HTTP_JSON_VALUE,
	}
}
