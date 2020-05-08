package database

import (
	"github.com/golang/protobuf/ptypes"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertToPostgresSnapshot(t *testing.T) {
	correctTime, _ := ptypes.TimestampProto(time.Now())
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToPostgresSnapshot(&apiPb.SchedulerResponse{})
		assert.Error(t, err)
	})
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToPostgresSnapshot(&apiPb.SchedulerResponse{
			SchedulerId: "id",
			Snapshot:    &apiPb.Snapshot{},
		})
		assert.Error(t, err)
	})
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToPostgresSnapshot(&apiPb.SchedulerResponse{
			SchedulerId: "id",
			Snapshot: &apiPb.Snapshot{
				Code: 0,
				Type: 0,
				Meta: &apiPb.Snapshot_MetaData{
					StartTime: nil,
				},
			},
		})
		assert.Error(t, err)
	})
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToPostgresSnapshot(&apiPb.SchedulerResponse{
			SchedulerId: "id",
			Snapshot: &apiPb.Snapshot{
				Code: 0,
				Type: 0,
				Meta: &apiPb.Snapshot_MetaData{
					StartTime: correctTime,
					EndTime:   nil,
				},
			},
		})
		assert.Error(t, err)
	})
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToPostgresSnapshot(&apiPb.SchedulerResponse{
			SchedulerId: "id",
			Snapshot: &apiPb.Snapshot{
				Code:  0,
				Type:  0,
				Error: &apiPb.Snapshot_SnapshotError{Message: ""},
				Meta: &apiPb.Snapshot_MetaData{
					StartTime: correctTime,
					EndTime:   correctTime,
				},
			},
		})
		assert.NoError(t, err)
	})
}

func TestConvertFromPostgresSnapshots(t *testing.T) {
	wrongTime := time.Unix(-62135596888, -100000000) //Protobuf validate this seconds aas error
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertFromPostgresSnapshots([]*Snapshot{{}})
		assert.NotEmpty(t, err)
	})
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertFromPostgresSnapshots([]*Snapshot{
			{
				Meta: &MetaData{
					StartTime: wrongTime,
					EndTime:   time.Time{},
				},
			},
		})
		assert.NotEmpty(t, err)
	})
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertFromPostgresSnapshots([]*Snapshot{
			{
				Meta: &MetaData{
					StartTime: time.Time{},
					EndTime:   wrongTime,
				},
			},
		})
		assert.NotEmpty(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertFromPostgresSnapshots([]*Snapshot{
			{
				Meta: &MetaData{
					StartTime: time.Time{},
					EndTime:   time.Time{},
				},
			},
		})
		assert.Empty(t, err)
	})
}

func TestConvertToPostgressStatRequest(t *testing.T) {
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToPostgressStatRequest(&apiPb.SendMetricsRequest{
			Time: nil,
		})
		assert.Error(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertToPostgressStatRequest(&apiPb.SendMetricsRequest{
			Time: ptypes.TimestampNow(),
		})
		assert.NoError(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertToPostgressStatRequest(&apiPb.SendMetricsRequest{
			CpuInfo: &apiPb.CpuInfo{},
			MemoryInfo: &apiPb.MemoryInfo{
				Mem:  &apiPb.MemoryInfo_Memory{},
				Swap: &apiPb.MemoryInfo_Memory{},
			},
			DiskInfo: &apiPb.DiskInfo{},
			NetInfo:  &apiPb.NetInfo{},
			Time:     ptypes.TimestampNow(),
		})
		assert.NoError(t, err)
	})
}
