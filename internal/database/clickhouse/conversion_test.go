package clickhouse

import (
	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"testing"
	"time"
)

func Test_convertToIncident(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		res := convertToIncident(nil, time.Time{})
		assert.Nil(t, res)
	})
}

func Test_convertFromIncidentHistory(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		res := convertFromIncidentHistory(nil)
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
	t.Run("Should: return empty res", func(t *testing.T) {
		res, _, _ := convertToIncidentHistories(nil)
		assert.Nil(t, res)
		assert.Equal(t, 0, len(res))
	})
}

func Test_convertToIncidentHistoryErr(t *testing.T) {
	t.Run("Should: return empty res", func(t *testing.T) {
		res := convertToIncidentHistory(nil)
		assert.Nil(t, res)
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
					Value: &structpb.Value{
						Kind: &structpb.Value_StringValue{
							StringValue: "hey",
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

func TestConvertToClickhouseStatRequest(t *testing.T) {
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToClickhouseStatRequest(&apiPb.Metric{
			Time: nil,
		})
		assert.Error(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertToClickhouseStatRequest(&apiPb.Metric{
			Time: ptypes.TimestampNow(),
		})
		assert.NoError(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertToClickhouseStatRequest(&apiPb.Metric{
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

func TestConvertFromClickhouseStatRequests(t *testing.T) {
	t.Run("Test: error", func(t *testing.T) {
		res := ConvertFromClickhouseStatRequests([]*StatRequest{
			{
				Time: time.Date(
					-1000, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
		})
		assert.Equal(t, 0, len(res))
	})
}

func TestConvertFromClickhouseStatRequest(t *testing.T) {
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertFromClickhouseStatRequest(&StatRequest{
			Time: time.Unix(-62135596888, -100000000), //Protobuf validate this as error
		})
		assert.Error(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertFromClickhouseStatRequest(&StatRequest{
			CPUInfo:    nil,
			MemoryInfo: nil,
			DiskInfo:   nil,
			NetInfo:    nil,
			Time:       time.Time{},
		})
		assert.NoError(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertFromClickhouseStatRequest(&StatRequest{
			CPUInfo:    []*CPUInfo{{}},
			MemoryInfo: &MemoryInfo{},
			DiskInfo:   []*DiskInfo{{}},
			NetInfo:    []*NetInfo{{}},
			Time:       time.Time{},
		})
		assert.NoError(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertFromClickhouseStatRequest(&StatRequest{
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

func TestConvertToClickhouseSnapshot(t *testing.T) {
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
					Value: &structpb.Value{
						Kind: &structpb.Value_StringValue{
							StringValue: "hey",
						},
					},
				},
			},
		})
		assert.NoError(t, err)
	})
}

func TestConvertFromSnapshots(t *testing.T) {
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

func TestConvertToStatRequest(t *testing.T) {
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertToClickhouseStatRequest(&apiPb.Metric{
			Time: nil,
		})
		assert.Error(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertToClickhouseStatRequest(&apiPb.Metric{
			Time: ptypes.TimestampNow(),
		})
		assert.NoError(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertToClickhouseStatRequest(&apiPb.Metric{
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

func TestConvertFromStatRequest(t *testing.T) {
	t.Run("Test: error", func(t *testing.T) {
		_, err := ConvertFromClickhouseStatRequest(&StatRequest{
			Time: time.Unix(-62135596888, -100000000), //Protobuf validate this seconds aas error
		})
		assert.Error(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertFromClickhouseStatRequest(&StatRequest{
			CPUInfo:    nil,
			MemoryInfo: nil,
			DiskInfo:   nil,
			NetInfo:    nil,
			Time:       time.Time{},
		})
		assert.NoError(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertFromClickhouseStatRequest(&StatRequest{
			CPUInfo:    []*CPUInfo{{}},
			MemoryInfo: &MemoryInfo{},
			DiskInfo:   []*DiskInfo{{}},
			NetInfo:    []*NetInfo{{}},
			Time:       time.Time{},
		})
		assert.NoError(t, err)
	})
	t.Run("Test: no error", func(t *testing.T) {
		_, err := ConvertFromClickhouseStatRequest(&StatRequest{
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

func TestConvertFromGroupResult(t *testing.T) {
	t.Run("Test: error", func(t *testing.T) {
		res := convertFromGroupResult([]*GroupResult{
			{
				Name:  "Name",
				Count: 0,
			},
		}, time.Now().UnixNano())
		assert.NotNil(t, res)
	})
	t.Run("Test: no error", func(t *testing.T) {
		res := convertFromGroupResult([]*GroupResult{
			{
				Name:    "Name",
				Count:   0,
				Latency: "10000.000",
			},
		}, time.Now().UnixNano())
		assert.NotNil(t, res)
	})
	t.Run("Test: no error", func(t *testing.T) {
		res := convertFromGroupResult([]*GroupResult{
			{
				Name:    "Name",
				Count:   0,
				Latency: "10000.000",
				MinTime: "10000.000",
			},
		}, time.Now().UnixNano())
		assert.NotNil(t, res)
	})
	t.Run("Test: no error", func(t *testing.T) {
		res := convertFromGroupResult([]*GroupResult{
			{
				Name:    "Name",
				Count:   0,
				Latency: "10000.000",
				MinTime: "10000.000",
				MaxTime: "10000.000",
			},
		}, time.Now().UnixNano())
		assert.NotNil(t, res)
	})
	t.Run("Test: no error", func(t *testing.T) {
		res := convertFromGroupResult([]*GroupResult{
			{
				Name:    "Name",
				Count:   0,
				Latency: "10000.000",
				MinTime: "10000.000",
				MaxTime: "10000.000",
				LowTime: "10000.000",
			},
		}, time.Now().UnixNano())
		assert.NotNil(t, res)
	})
}

func TestGetThroughput(t *testing.T) {
	t.Run("Test: error", func(t *testing.T) {
		res := getThroughput(0, 10, 10)
		assert.NotNil(t, res)
	})
}

func TestConvertFromTransaction(t *testing.T) {
	t.Run("Test: empty meta", func(t *testing.T) {
		res := convertFromTransaction(&TransactionInfo{})
		assert.Nil(t, res.Meta)
	})
	t.Run("Test: empty error", func(t *testing.T) {
		res := convertFromTransaction(&TransactionInfo{
			MetaHost:   "a",
			MetaPath:   "b",
			MetaMethod: "c",
			Error:      "",
		})
		assert.Nil(t, res.Error)
	})
}

func TestConvertToTransactionInfo(t *testing.T) {
	t.Run("Test: error", func(t *testing.T) {
		_, err := convertToTransactionInfo(&apiPb.TransactionInfo{
			StartTime: nil,
		})
		assert.Error(t, err)
	})
	t.Run("Test: error", func(t *testing.T) {
		_, err := convertToTransactionInfo(&apiPb.TransactionInfo{
			StartTime: ptypes.TimestampNow(),
			EndTime:   nil,
		})
		assert.Error(t, err)
	})
}

func Test_conversion(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		o, l := getOffsetAndLimit(0, &apiPb.Pagination{Page: 0, Limit: 5})
		assert.Equal(t, 0, o)
		assert.Equal(t, 5, l)
	})
}
