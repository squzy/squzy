package clickhouse

import (
	"github.com/golang/protobuf/ptypes"
	_struct "github.com/golang/protobuf/ptypes/struct"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_convertFromIncidentHistory(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		res := convertToIncident(nil, time.Time{})
		assert.Nil(t, res)
	})
}

func Test_convertToIncidentHistory(t *testing.T) {
	t.Run("Should: return empty res", func(t *testing.T) {
		res, _, _ := convertToIncidentHistories(nil)
		assert.Equal(t, 0, len(res))
	})
	t.Run("Should: return empty res", func(t *testing.T) {
		maxValidSeconds := 253402300800
		res, _, _ := convertToIncidentHistories([]*apiPb.Incident_HistoryItem{{
			Timestamp: &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		}})
		assert.Equal(t, 0, len(res))
	})
}

func TestConvertToSnapshot(t *testing.T) {
	correctTime, _ := ptypes.TimestampProto(time.Now())
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToSnapshot(&apiPb.SchedulerResponse{})
		assert.Error(t, err)
	})
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToSnapshot(&apiPb.SchedulerResponse{
			SchedulerId: "id",
			Snapshot:    &apiPb.SchedulerSnapshot{},
		})
		assert.Error(t, err)
	})
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToSnapshot(&apiPb.SchedulerResponse{
			SchedulerId: "id",
			Snapshot: &apiPb.SchedulerSnapshot{
				Code: 0,
				Type: 0,
				Meta: &apiPb.SchedulerSnapshot_MetaData{
					StartTime: nil,
				},
			},
		})
		assert.Error(t, err)
	})
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToSnapshot(&apiPb.SchedulerResponse{
			SchedulerId: "id",
			Snapshot: &apiPb.SchedulerSnapshot{
				Code: 0,
				Type: 0,
				Meta: &apiPb.SchedulerSnapshot_MetaData{
					StartTime: correctTime,
					EndTime:   nil,
				},
			},
		})
		assert.Error(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertToSnapshot(&apiPb.SchedulerResponse{
			SchedulerId: "id",
			Snapshot: &apiPb.SchedulerSnapshot{
				Code:  0,
				Type:  0,
				Error: &apiPb.SchedulerSnapshot_Error{Message: ""},
				Meta: &apiPb.SchedulerSnapshot_MetaData{
					StartTime: correctTime,
					EndTime:   correctTime,
					Value: &_struct.Value{
						Kind: &_struct.Value_StringValue{
							StringValue: "HUY",
						},
					},
				},
			},
		})
		assert.NoError(t, err)
	})
}

func TestConvertFromClickhouseSnapshots(t *testing.T) {
	wrongTime := time.Unix(-62135596888, -100000000) //Protobuf validate this seconds aas error
	t.Run("Test: error in convertation", func(t *testing.T) {
		res := ConvertFromSnapshots([]*Snapshot{{}})
		assert.NotNil(t, res)
	})
	t.Run("Test: error in convertation", func(t *testing.T) {
		res := ConvertFromSnapshots([]*Snapshot{
			{
				MetaStartTime: wrongTime.UnixNano(),
			},
		})
		assert.NotNil(t, res)
	})
	t.Run("Test: error", func(t *testing.T) {
		res := ConvertFromSnapshots([]*Snapshot{
			{
				Error: "error",
			},
		})
		assert.NotNil(t, res)
	})
	t.Run("Test: no error", func(t *testing.T) {
		res := ConvertFromSnapshots([]*Snapshot{
			{
				MetaValue: nil,
			},
		})
		assert.NotNil(t, res)
	})
	t.Run("Test: no error", func(t *testing.T) {
		res := ConvertFromSnapshots([]*Snapshot{
			{
				MetaValue: []byte(`{"stringValue":"HUY"}`),
			},
		})
		assert.NotNil(t, res)
	})
}

func TestConvertFromUptimeResult(t *testing.T) {
	t.Run("Test: error", func(t *testing.T) {
		res := convertFromUptimeResult(&UptimeResult{
			Count:   10,
			Latency: "qwe",
		}, 10)
		assert.NotNil(t, res)
	})
	t.Run("Test: error", func(t *testing.T) {
		res := convertFromUptimeResult(&UptimeResult{
			Count:   10,
			Latency: "10000",
		}, 10)
		assert.NotNil(t, res)
	})
}