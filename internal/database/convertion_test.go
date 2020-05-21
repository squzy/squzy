package database

import (
	"github.com/golang/protobuf/ptypes"
	_struct "github.com/golang/protobuf/ptypes/struct"
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
			Snapshot:    &apiPb.SchedulerSnapshot{},
		})
		assert.Error(t, err)
	})
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToPostgresSnapshot(&apiPb.SchedulerResponse{
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
		_, err := ConvertToPostgresSnapshot(&apiPb.SchedulerResponse{
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
		_, err := ConvertToPostgresSnapshot(&apiPb.SchedulerResponse{
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

func TestConvertFromPostgresSnapshots(t *testing.T) {
	wrongTime := time.Unix(-62135596888, -100000000) //Protobuf validate this seconds aas error
	t.Run("Test: error in convertation", func(t *testing.T) {
		res := ConvertFromPostgresSnapshots([]*Snapshot{{}})
		assert.Nil(t, res)
	})
	t.Run("Test: error in convertation", func(t *testing.T) {
		res := ConvertFromPostgresSnapshots([]*Snapshot{
			{
				Meta: &MetaData{
					StartTime: wrongTime,
					EndTime:   time.Time{},
				},
			},
		})
		assert.Nil(t, res)
	})
	t.Run("Test: error", func(t *testing.T) {
		res := ConvertFromPostgresSnapshots([]*Snapshot{
			{
				Meta: &MetaData{
					StartTime: time.Time{},
					EndTime:   wrongTime,
				},
			},
		})
		assert.Nil(t, res)
	})
	t.Run("Test: no error", func(t *testing.T) {
		res := ConvertFromPostgresSnapshots([]*Snapshot{
			{
				Meta: &MetaData{
					StartTime: time.Time{},
					EndTime:   time.Time{},
					Value:     nil,
				},
			},
		})
		assert.NotNil(t, res)
	})
	t.Run("Test: no error", func(t *testing.T) {
		res := ConvertFromPostgresSnapshots([]*Snapshot{
			{
				Error: "error",
				Meta: &MetaData{
					StartTime: time.Time{},
					EndTime:   time.Time{},
					Value:     []byte(`{"stringValue":"HUY"}`),
				},
			},
		})
		assert.NotNil(t, res)
	})
}

func TestConvertToPostgressStatRequest(t *testing.T) {
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToPostgressStatRequest(&apiPb.Metric{
			Time: nil,
		})
		assert.Error(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertToPostgressStatRequest(&apiPb.Metric{
			Time: ptypes.TimestampNow(),
		})
		assert.NoError(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertToPostgressStatRequest(&apiPb.Metric{
			CpuInfo: &apiPb.CpuInfo{
				Cpus: []*apiPb.CpuInfo_CPU{{}},
			},
			MemoryInfo: &apiPb.MemoryInfo{
				Mem:  &apiPb.MemoryInfo_Memory{},
				Swap: &apiPb.MemoryInfo_Memory{},
			},
			DiskInfo: &apiPb.DiskInfo{
				Disks: map[string]*apiPb.DiskInfo_Disk{
					"": {},
				},
			},
			NetInfo: &apiPb.NetInfo{
				Interfaces: map[string]*apiPb.NetInfo_Interface{
					"": {},
				},
			},
			Time: ptypes.TimestampNow(),
		})
		assert.NoError(t, err)
	})
}

func TestConvertFromPostgressStatRequest(t *testing.T) {
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertFromPostgressStatRequest(&StatRequest{
			Time: time.Unix(-62135596888, -100000000), //Protobuf validate this seconds aas error
		})
		assert.Error(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertFromPostgressStatRequest(&StatRequest{
			CPUInfo:    nil,
			MemoryInfo: nil,
			DiskInfo:   nil,
			NetInfo:    nil,
			Time:       time.Time{},
		})
		assert.NoError(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertFromPostgressStatRequest(&StatRequest{
			CPUInfo:    []*CPUInfo{{}},
			MemoryInfo: &MemoryInfo{},
			DiskInfo:   []*DiskInfo{{}},
			NetInfo:    []*NetInfo{{}},
			Time:       time.Time{},
		})
		assert.NoError(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertFromPostgressStatRequest(&StatRequest{
			CPUInfo: []*CPUInfo{{}},
			MemoryInfo: &MemoryInfo{
				Mem:  &MemoryMem{},
				Swap: &MemorySwap{},
			},
			DiskInfo: []*DiskInfo{{}},
			NetInfo:  []*NetInfo{{}},
			Time:     time.Time{},
		})
		assert.NoError(t, err)
	})
}
