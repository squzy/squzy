package database

import (
	"github.com/golang/protobuf/ptypes"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertToPostgressScheduler(t *testing.T) {
	correctTime, _ := ptypes.TimestampProto(time.Now())
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToPostgresScheduler(&apiPb.SchedulerResponse{})
		assert.Error(t, err)
	})
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToPostgresScheduler(&apiPb.SchedulerResponse{
			SchedulerId: "id",
			Snapshot:    &apiPb.Snapshot{},
		})
		assert.Error(t, err)
	})
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToPostgresScheduler(&apiPb.SchedulerResponse{
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
		_, err := ConvertToPostgresScheduler(&apiPb.SchedulerResponse{
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
		_, err := ConvertToPostgresScheduler(&apiPb.SchedulerResponse{
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
