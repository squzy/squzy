package expression

import (
	"context"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
	"time"
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
	for _, filter := range filters {
		req = filter(req)
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
		"Index": func(index int32, filters ...FilterSnapshot) *apiPb.SchedulerSnapshot {
			return e.GetSnapshots(
				schedulerId,
				apiPb.SortDirection_ASC,
				&apiPb.Pagination{
					Page:  index,
					Limit: 1,
				},
				filters...)[0]
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
		"UnixNanoNow": func() int64 {
			return time.Now().UnixNano()
		},
		"timeDiff": func(t1, t2 time.Time) time.Duration {
			return t1.Sub(t2)
		},
		"durationLess": func(d1, d2 time.Duration) bool {
			return d1 < d2
		},
		"durationMore": func(d1, d2 time.Duration) bool {
			return d1 > d2
		},
		"durationEqual": func(d1, d2 time.Duration) bool {
			return d1 == d2
		},
		"durationToSecond": func(d time.Duration) int64 {
			return int64(d.Seconds())
		},
		"NowTime": func() time.Time {
			return time.Now()
		},
		"float64ToInt64": func(v float64) int64 {
			return int64(v)
		},
		"getValue": func(snapshot *apiPb.SchedulerSnapshot) *structpb.Value {
			return snapshot.GetMeta().GetValue()
		},
		"unixToTime": func(unix int64) time.Time {
			return time.Unix(unix, 0)
		},
		"unixNanoToTime": func(unixNano int64) time.Time {
			return time.Unix(0, unixNano)
		},
		"null": nil,
		"mulDuration": func(f int, t time.Duration) time.Duration {
			return time.Duration(f) * t
		},
		"Week":   time.Hour * 24 * 7,
		"Day":    time.Hour * 24,
		"Hour":   time.Hour,
		"Minute": time.Minute,
		"Second": time.Second,
		//Transaction status keys
		"Ok":    apiPb.SchedulerCode_OK,
		"Error": apiPb.SchedulerCode_ERROR,
	}
}
