package clickhouse

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"os"
	"sort"
	"testing"
	"time"
)

var (
	pool        *dockertest.Pool
	resource    *dockertest.Resource
	wdb, _      = sql.Open("clickhouse", "tcp://user:password@lkl:00/debug=true&clicks?read_timeout=10&write_timeout=10")
	clickhWrong = &Clickhouse{
		wdb,
	}
	clickh *Clickhouse
)

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		logger.Fatalf("could not start test: %s", err)
	}
	code := m.Run()
	err = shutdown()
	if err != nil {
		logger.Fatalf("could not stop test: %s", err)
	}
	os.Exit(code)
}

func shutdown() error {
	if err := pool.Purge(resource); err != nil {
		return fmt.Errorf("could not purge resource: %w", err)
	}
	return nil
}

func setup() error {
	var db *sql.DB
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	resource, err = pool.RunWithOptions(&dockertest.RunOptions{Repository: "yandex/clickhouse-server",
		Tag:          "20.3.11.97",
		Cmd:          []string{"start-single-node", "--insecure"},
		ExposedPorts: []string{"9000/tcp", "8123/tcp"},
		//PortBindings: map[docker.Port][]docker.PortBinding{
		//	"9000/tcp": {{HostIP: "", HostPort: "9000"}},
		//	"8123/tcp": {{HostIP: "", HostPort: "8123"}},
		//},
	},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		db, err = sql.Open("clickhouse", fmt.Sprintf("tcp://%s:%s?debug=true", "localhost", resource.GetPort("9000/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to clickhouse container: %s", err)
	}

	err = clickh.Migrate()
	if err != nil {
		return err
	}
	return nil
}

func TestInsertIncident(t *testing.T) {
	lo := &apiPb.Incident{
		Id:     "insert",
		Status: 0,
		RuleId: "433",
		Histories: []*apiPb.Incident_HistoryItem{&apiPb.Incident_HistoryItem{
			Status: 0,
			Timestamp: &timestamp.Timestamp{
				Seconds: 3324,
				Nanos:   0,
			},
		}},
	}
	err := clickh.InsertIncident(lo)
	if err != nil {
		assert.Fail(t, err.Error())
	}
}

func TestGetIncidentById(t *testing.T) {
	lo := &apiPb.Incident{
		Id:     "select",
		Status: 1,
		RuleId: "433",
		Histories: []*apiPb.Incident_HistoryItem{
			&apiPb.Incident_HistoryItem{
				Status: 0,
				Timestamp: &timestamp.Timestamp{
					Seconds: 3324,
					Nanos:   0,
				},
			},
		},
	}
	err := clickh.InsertIncident(lo)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	inc, err := clickh.GetIncidentById(lo.Id)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	sort.Slice(lo.Histories, func(i, j int) bool {
		return lo.Histories[i].Status < lo.Histories[j].Status
	})
	assert.NotNil(t, inc)
	assert.Equal(t, lo.Id, inc.Id)
	assert.Equal(t, lo.Status, inc.Status)
	assert.Equal(t, lo.RuleId, inc.RuleId)
	assert.Equal(t, lo.Histories[0].Status, inc.Histories[0].Status)
	assert.Equal(t, lo.Histories[0].Timestamp, inc.Histories[0].Timestamp)
	assert.NotNil(t, inc)
}

func TestUpdateIncidentStatus(t *testing.T) {
	lo := &apiPb.Incident{
		Id:     "update",
		Status: 1,
		RuleId: "433",
		Histories: []*apiPb.Incident_HistoryItem{&apiPb.Incident_HistoryItem{
			Status: 0,
			Timestamp: &timestamp.Timestamp{
				Seconds: 3324,
				Nanos:   0,
			},
		}},
	}
	err := clickh.InsertIncident(lo)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	inc, err := clickh.UpdateIncidentStatus(lo.Id, apiPb.IncidentStatus(2))
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.NotNil(t, inc)
	assert.Equal(t, lo.Id, inc.Id)
	assert.Equal(t, apiPb.IncidentStatus(2), inc.Status)
	assert.Equal(t, 2, len(inc.Histories))

}

func TestGetActiveIncidentByRuleId(t *testing.T) {
	lo := &apiPb.Incident{
		Id:     "active",
		Status: 1,
		RuleId: "some rule",
		Histories: []*apiPb.Incident_HistoryItem{&apiPb.Incident_HistoryItem{
			Status: 0,
			Timestamp: &timestamp.Timestamp{
				Seconds: 3324,
				Nanos:   0,
			},
		}},
	}
	err := clickh.InsertIncident(lo)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	inc, err := clickh.GetActiveIncidentByRuleId(lo.RuleId)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.NotNil(t, inc)
	assert.Equal(t, lo.Id, inc.Id)
	assert.Equal(t, lo.Status, inc.Status)
	assert.Equal(t, lo.RuleId, inc.RuleId)
	assert.Equal(t, lo.Histories[0].Status, inc.Histories[0].Status)
	assert.Equal(t, lo.Histories[0].Timestamp, inc.Histories[0].Timestamp)
}

func TestGetIncidents(t *testing.T) {
	lo := &apiPb.Incident{
		Id:     "incidents",
		Status: 1,
		RuleId: "999",
		Histories: []*apiPb.Incident_HistoryItem{&apiPb.Incident_HistoryItem{
			Status: 1,
			Timestamp: &timestamp.Timestamp{
				Seconds: 3324,
				Nanos:   0,
			},
		}},
	}
	lo2 := &apiPb.Incident{
		Id:     "incidents2",
		Status: 1,
		RuleId: "999",
		Histories: []*apiPb.Incident_HistoryItem{&apiPb.Incident_HistoryItem{
			Status: 1,
			Timestamp: &timestamp.Timestamp{
				Seconds: 3424,
				Nanos:   0,
			},
		}},
	}
	lo3 := &apiPb.Incident{
		Id:     "incidents3",
		Status: 3,
		RuleId: "999",
	}

	err := clickh.InsertIncident(lo)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	err = clickh.InsertIncident(lo2)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	err = clickh.InsertIncident(lo3)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	timeTo, err := ptypes.TimestampProto(time.Now().Add(time.Second * 5))
	if err != nil {
		assert.Fail(t, err.Error())
	}
	incs, count, err := clickh.GetIncidents(&apiPb.GetIncidentsListRequest{
		Status: 1,
		RuleId: &wrappers.StringValue{Value: "999"},
		Pagination: &apiPb.Pagination{
			Page:  1,
			Limit: 10,
		},
		TimeRange: &apiPb.TimeFilter{
			From: lo.Histories[0].Timestamp,
			To:   timeTo,
		},
		Sort: &apiPb.SortingIncidentList{
			SortBy:    apiPb.SortIncidentList_INCIDENT_LIST_BY_START_TIME,
			Direction: apiPb.SortDirection_ASC,
		},
	})
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, 2, int(count))
	assert.NotNil(t, incs)
	assert.Equal(t, lo.Id, incs[0].Id)
	assert.Equal(t, 1, int(incs[0].Status))
	assert.Equal(t, lo.Histories[0].Status, incs[0].Histories[0].Status)
	assert.Equal(t, lo.Histories[0].Timestamp, incs[0].Histories[0].Timestamp)
	assert.Equal(t, lo2.RuleId, incs[0].RuleId)
	assert.Equal(t, lo2.Id, incs[1].Id)
	assert.Equal(t, 1, int(incs[1].Status))
	assert.Equal(t, lo2.RuleId, incs[1].RuleId)
	assert.Equal(t, lo2.Histories[0].Status, incs[1].Histories[0].Status)
	assert.Equal(t, lo2.Histories[0].Timestamp, incs[1].Histories[0].Timestamp)
}

func TestInsertSnapshot(t *testing.T) {
	sn := &apiPb.SchedulerResponse{
		SchedulerId: "insert",
		Snapshot: &apiPb.SchedulerSnapshot{
			Code:  apiPb.SchedulerCode_OK,
			Type:  apiPb.SchedulerType_TCP,
			Error: nil,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: &timestamp.Timestamp{
					Seconds: 3324,
					Nanos:   0,
				},
				EndTime: &timestamp.Timestamp{
					Seconds: 3324,
					Nanos:   0,
				},
				Value: &structpb.Value{
					Kind: &structpb.Value_StringValue{
						StringValue: "Value",
					},
				},
			},
		},
	}
	err := clickh.InsertSnapshot(sn)
	if err != nil {
		assert.Fail(t, err.Error())
	}
}

func TestGetSnapshots(t *testing.T) {
	sn := &apiPb.SchedulerResponse{
		SchedulerId: "GetSnapshots",
		Snapshot: &apiPb.SchedulerSnapshot{
			Code:  apiPb.SchedulerCode_OK,
			Type:  apiPb.SchedulerType_TCP,
			Error: nil,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: &timestamp.Timestamp{
					Seconds: 3324,
					Nanos:   0,
				},
				EndTime: &timestamp.Timestamp{
					Seconds: 3324,
					Nanos:   0,
				},
				Value: &structpb.Value{
					Kind: &structpb.Value_StringValue{
						StringValue: "Value",
					},
				},
			},
		},
	}

	sn2 := &apiPb.SchedulerResponse{
		SchedulerId: "GetSnapshots",
		Snapshot: &apiPb.SchedulerSnapshot{
			Code:  apiPb.SchedulerCode_OK,
			Type:  apiPb.SchedulerType_TCP,
			Error: nil,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: &timestamp.Timestamp{
					Seconds: 3324,
					Nanos:   0,
				},
				EndTime: &timestamp.Timestamp{
					Seconds: 3324,
					Nanos:   0,
				},
				Value: &structpb.Value{
					Kind: &structpb.Value_StringValue{
						StringValue: "Value",
					},
				},
			},
		},
	}

	sn3 := &apiPb.SchedulerResponse{
		SchedulerId: "GetSnapshots",
		Snapshot: &apiPb.SchedulerSnapshot{
			Code: apiPb.SchedulerCode_ERROR,
			Type: apiPb.SchedulerType_TCP,
			Error: &apiPb.SchedulerSnapshot_Error{
				Message: "Error",
			},
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: &timestamp.Timestamp{
					Seconds: 3324,
					Nanos:   0,
				},
				EndTime: &timestamp.Timestamp{
					Seconds: 3324,
					Nanos:   0,
				},
				Value: &structpb.Value{
					Kind: &structpb.Value_StringValue{
						StringValue: "Value",
					},
				},
			},
		},
	}

	err := clickh.InsertSnapshot(sn)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	err = clickh.InsertSnapshot(sn2)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	err = clickh.InsertSnapshot(sn3)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	timeTo, err := ptypes.TimestampProto(time.Now().Add(time.Second * 5))
	if err != nil {
		assert.Fail(t, err.Error())
	}
	snaps, count, err := clickh.GetSnapshots(&apiPb.GetSchedulerInformationRequest{
		SchedulerId: "GetSnapshots",
		Pagination: &apiPb.Pagination{
			Page:  1,
			Limit: 10,
		},
		TimeRange: &apiPb.TimeFilter{
			From: sn.Snapshot.Meta.StartTime,
			To:   timeTo,
		},
		Sort: &apiPb.SortingSchedulerList{
			SortBy:    apiPb.SortSchedulerList_BY_START_TIME,
			Direction: apiPb.SortDirection_ASC,
		},
		Status: apiPb.SchedulerCode_OK,
	})

	assert.Equal(t, 2, int(count))
	assert.NotNil(t, snaps)
	assert.Equal(t, sn.Snapshot.Code, snaps[0].Code)
	assert.Equal(t, sn.Snapshot.Type, snaps[0].Type)
	assert.Equal(t, sn.Snapshot.Error, snaps[0].Error)
	assert.Equal(t, sn.Snapshot.Meta, snaps[0].Meta)
	assert.Equal(t, sn.Snapshot.Code, snaps[1].Code)
	assert.Equal(t, sn.Snapshot.Type, snaps[1].Type)
	assert.Equal(t, sn.Snapshot.Error, snaps[1].Error)
	assert.Equal(t, sn.Snapshot.Meta, snaps[1].Meta)

}

func TestGetSnapshotsUptime(t *testing.T) {
	sn := &apiPb.SchedulerResponse{
		SchedulerId: "GetSnapshotsUptime",
		Snapshot: &apiPb.SchedulerSnapshot{
			Code:  apiPb.SchedulerCode_OK,
			Type:  apiPb.SchedulerType_TCP,
			Error: nil,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: &timestamp.Timestamp{
					Seconds: 10000,
					Nanos:   0,
				},
				EndTime: &timestamp.Timestamp{
					Seconds: 10000,
					Nanos:   10,
				},
				Value: &structpb.Value{
					Kind: &structpb.Value_StringValue{
						StringValue: "Value",
					},
				},
			},
		},
	}

	sn2 := &apiPb.SchedulerResponse{
		SchedulerId: "GetSnapshotsUptime",
		Snapshot: &apiPb.SchedulerSnapshot{
			Code:  apiPb.SchedulerCode_OK,
			Type:  apiPb.SchedulerType_TCP,
			Error: nil,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: &timestamp.Timestamp{
					Seconds: 10000,
					Nanos:   0,
				},
				EndTime: &timestamp.Timestamp{
					Seconds: 10000,
					Nanos:   20,
				},
				Value: &structpb.Value{
					Kind: &structpb.Value_StringValue{
						StringValue: "Value",
					},
				},
			},
		},
	}

	sn3 := &apiPb.SchedulerResponse{
		SchedulerId: "GetSnapshots",
		Snapshot: &apiPb.SchedulerSnapshot{
			Code: apiPb.SchedulerCode_ERROR,
			Type: apiPb.SchedulerType_TCP,
			Error: &apiPb.SchedulerSnapshot_Error{
				Message: "Error",
			},
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: &timestamp.Timestamp{
					Seconds: 1000,
					Nanos:   0,
				},
				EndTime: &timestamp.Timestamp{
					Seconds: 5000,
					Nanos:   0,
				},
				Value: &structpb.Value{
					Kind: &structpb.Value_StringValue{
						StringValue: "Value",
					},
				},
			},
		},
	}

	err := clickh.InsertSnapshot(sn)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	err = clickh.InsertSnapshot(sn2)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	err = clickh.InsertSnapshot(sn3)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	timeTo, err := ptypes.TimestampProto(time.Now().Add(time.Second * 5))
	if err != nil {
		assert.Fail(t, err.Error())
	}
	resp, err := clickh.GetSnapshotsUptime(&apiPb.GetSchedulerUptimeRequest{
		SchedulerId: "GetSnapshotsUptime",
		TimeRange: &apiPb.TimeFilter{
			From: nil,
			To:   timeTo,
		},
	})

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, float64(15), resp.Latency)
	assert.Equal(t, float64(1), resp.Uptime)

}

func TestInsertStatRequest(t *testing.T) {
	di := map[string]*apiPb.DiskInfo_Disk{"1": {
		Total:       100,
		Used:        100,
		Free:        100,
		UsedPercent: 100,
	}}

	ni := map[string]*apiPb.NetInfo_Interface{"1": {
		BytesSent:   100,
		BytesRecv:   100,
		PacketsSent: 100,
		PacketsRecv: 100,
		ErrIn:       100,
		ErrOut:      100,
		DropIn:      100,
		DropOut:     100,
	}}

	sr := &apiPb.Metric{
		AgentId:   "insertStatRequest",
		AgentName: "AgentName",
		CpuInfo: &apiPb.CpuInfo{
			Cpus: []*apiPb.CpuInfo_CPU{&apiPb.CpuInfo_CPU{
				Load: 100.0,
			}},
		},
		MemoryInfo: &apiPb.MemoryInfo{
			Mem: &apiPb.MemoryInfo_Memory{
				Total:       100,
				Used:        100,
				Free:        100,
				Shared:      100,
				UsedPercent: 100,
			},
			Swap: &apiPb.MemoryInfo_Memory{
				Total:       100,
				Used:        100,
				Free:        100,
				Shared:      100,
				UsedPercent: 100,
			},
		},
		DiskInfo: &apiPb.DiskInfo{
			Disks: di,
		},
		NetInfo: &apiPb.NetInfo{
			Interfaces: ni,
		},
		Time: &timestamp.Timestamp{
			Seconds: 1789,
			Nanos:   0,
		},
	}
	err := clickh.InsertStatRequest(sr)
	if err != nil {
		assert.Fail(t, err.Error())
	}
}

func TestGetStatRequest(t *testing.T) {
	di := map[string]*apiPb.DiskInfo_Disk{"1": {
		Total:       100,
		Used:        100,
		Free:        100,
		UsedPercent: 100,
	}}

	ni := map[string]*apiPb.NetInfo_Interface{"1": {
		BytesSent:   100,
		BytesRecv:   100,
		PacketsSent: 100,
		PacketsRecv: 100,
		ErrIn:       100,
		ErrOut:      100,
		DropIn:      100,
		DropOut:     100,
	}}

	sr := &apiPb.Metric{
		AgentId:   "getStatRequest",
		AgentName: "AgentName",
		CpuInfo: &apiPb.CpuInfo{
			Cpus: []*apiPb.CpuInfo_CPU{&apiPb.CpuInfo_CPU{
				Load: 100.0,
			}},
		},
		MemoryInfo: &apiPb.MemoryInfo{
			Mem: &apiPb.MemoryInfo_Memory{
				Total:       100,
				Used:        100,
				Free:        100,
				Shared:      100,
				UsedPercent: 100,
			},
			Swap: &apiPb.MemoryInfo_Memory{
				Total:       100,
				Used:        100,
				Free:        100,
				Shared:      100,
				UsedPercent: 100,
			},
		},
		DiskInfo: &apiPb.DiskInfo{
			Disks: di,
		},
		NetInfo: &apiPb.NetInfo{
			Interfaces: ni,
		},
		Time: &timestamp.Timestamp{
			Seconds: 1788,
			Nanos:   0,
		},
	}

	sr2 := &apiPb.Metric{
		AgentId:   "getStatRequest",
		AgentName: "AgentName",
		CpuInfo: &apiPb.CpuInfo{
			Cpus: []*apiPb.CpuInfo_CPU{&apiPb.CpuInfo_CPU{
				Load: 100.0,
			}},
		},
		MemoryInfo: &apiPb.MemoryInfo{
			Mem: &apiPb.MemoryInfo_Memory{
				Total:       100,
				Used:        100,
				Free:        100,
				Shared:      100,
				UsedPercent: 100,
			},
			Swap: &apiPb.MemoryInfo_Memory{
				Total:       100,
				Used:        100,
				Free:        100,
				Shared:      100,
				UsedPercent: 100,
			},
		},
		DiskInfo: &apiPb.DiskInfo{
			Disks: di,
		},
		NetInfo: &apiPb.NetInfo{
			Interfaces: ni,
		},
		Time: &timestamp.Timestamp{
			Seconds: 1789,
			Nanos:   0,
		},
	}

	err := clickh.InsertStatRequest(sr)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	err = clickh.InsertStatRequest(sr2)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	srs, count, err := clickh.GetStatRequest(sr.AgentId, &apiPb.Pagination{
		Page:  1,
		Limit: 10,
	}, &apiPb.TimeFilter{
		From: &timestamp.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
		To: &timestamp.Timestamp{
			Seconds: 1789,
			Nanos:   0,
		},
	})
	if err != nil {
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, float64(100), srs[0].CpuInfo.Cpus[0].Load)
	assert.Equal(t, float64(100), srs[1].CpuInfo.Cpus[0].Load)
	assert.Equal(t, uint64(100), srs[0].NetInfo.Interfaces["1"].BytesRecv)
	assert.Equal(t, uint64(100), srs[1].NetInfo.Interfaces["1"].BytesRecv)
	assert.Equal(t, uint64(100), srs[0].DiskInfo.Disks["1"].Free)
	assert.Equal(t, uint64(100), srs[1].DiskInfo.Disks["1"].Free)
	assert.Equal(t, uint64(100), srs[0].MemoryInfo.Mem.Free)
	assert.Equal(t, uint64(100), srs[0].MemoryInfo.Swap.Free)
	assert.Equal(t, uint64(100), srs[1].MemoryInfo.Mem.Free)
	assert.Equal(t, uint64(100), srs[1].MemoryInfo.Swap.Free)
	assert.Equal(t, count, int64(2))

	//todo: move these tests
	srsCPU, count, err := clickh.GetCPUInfo(sr.AgentId, &apiPb.Pagination{
		Page:  1,
		Limit: 10,
	}, &apiPb.TimeFilter{
		From: &timestamp.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
		To: &timestamp.Timestamp{
			Seconds: 1789,
			Nanos:   0,
		},
	})
	if err != nil {
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, float64(100), srsCPU[0].CpuInfo.Cpus[0].Load)
	assert.Equal(t, float64(100), srsCPU[1].CpuInfo.Cpus[0].Load)
	assert.Equal(t, count, int64(2))

	srsMem, count, err := clickh.GetMemoryInfo(sr.AgentId, &apiPb.Pagination{
		Page:  1,
		Limit: 10,
	}, &apiPb.TimeFilter{
		From: &timestamp.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
		To: &timestamp.Timestamp{
			Seconds: 1789,
			Nanos:   0,
		},
	})
	if err != nil {
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, uint64(100), srsMem[0].MemoryInfo.Mem.Free)
	assert.Equal(t, uint64(100), srsMem[0].MemoryInfo.Swap.Free)
	assert.Equal(t, uint64(100), srsMem[1].MemoryInfo.Mem.Free)
	assert.Equal(t, uint64(100), srsMem[1].MemoryInfo.Swap.Free)
	assert.Equal(t, count, int64(2))

	srsDisk, count, err := clickh.GetDiskInfo(sr.AgentId, &apiPb.Pagination{
		Page:  1,
		Limit: 10,
	}, &apiPb.TimeFilter{
		From: &timestamp.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
		To: &timestamp.Timestamp{
			Seconds: 1789,
			Nanos:   0,
		},
	})
	if err != nil {
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, uint64(100), srsDisk[0].DiskInfo.Disks["1"].Free)
	assert.Equal(t, uint64(100), srsDisk[1].DiskInfo.Disks["1"].Free)
	assert.Equal(t, count, int64(2))

	srsNet, count, err := clickh.GetNetInfo(sr.AgentId, &apiPb.Pagination{
		Page:  1,
		Limit: 10,
	}, &apiPb.TimeFilter{
		From: &timestamp.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
		To: &timestamp.Timestamp{
			Seconds: 1789,
			Nanos:   0,
		},
	})
	if err != nil {
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, uint64(100), srsNet[0].NetInfo.Interfaces["1"].BytesRecv)
	assert.Equal(t, uint64(100), srsNet[1].NetInfo.Interfaces["1"].BytesRecv)
	assert.Equal(t, count, int64(2))
}

func TestInsertTransactionInfo(t *testing.T) {
	err := clickh.InsertTransactionInfo(&apiPb.TransactionInfo{
		Id:            "InsertTransactionInfo",
		ApplicationId: "ApplicationId",
		ParentId:      "ParentId",
		Meta: &apiPb.TransactionInfo_Meta{
			Host:   "Host",
			Path:   "Path",
			Method: "Method",
		},
		Name: "Name",
		StartTime: &timestamp.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
		EndTime: &timestamp.Timestamp{
			Seconds: 6789,
			Nanos:   0,
		},
		Status: 1,
		Type:   1,
		Error: &apiPb.TransactionInfo_Error{
			Message: "TransactionInfo_Error",
		},
	})
	if err != nil {
		assert.Fail(t, err.Error())
	}
}

func TestGetTransactionInfo(t *testing.T) {
	tr1 := &apiPb.TransactionInfo{
		Id:            "InsertTransactionInfo",
		ApplicationId: "ApplicationId",
		ParentId:      "ParentId",
		Meta: &apiPb.TransactionInfo_Meta{
			Host:   "Host",
			Path:   "Path",
			Method: "Method",
		},
		Name: "Name",
		StartTime: &timestamp.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
		EndTime: &timestamp.Timestamp{
			Seconds: 6789,
			Nanos:   0,
		},
		Status: 1,
		Type:   1,
		Error: &apiPb.TransactionInfo_Error{
			Message: "TransactionInfo_Error",
		},
	}

	tr2 := &apiPb.TransactionInfo{
		Id:            "InsertTransactionInfo2",
		ApplicationId: "ApplicationId",
		ParentId:      "ParentId",
		Meta: &apiPb.TransactionInfo_Meta{
			Host:   "Host",
			Path:   "Path",
			Method: "Method",
		},
		Name: "Name",
		StartTime: &timestamp.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
		EndTime: &timestamp.Timestamp{
			Seconds: 6789,
			Nanos:   0,
		},
		Status: 1,
		Type:   1,
		Error: &apiPb.TransactionInfo_Error{
			Message: "TransactionInfo_Error",
		},
	}

	err := clickh.InsertTransactionInfo(tr1)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	err = clickh.InsertTransactionInfo(tr2)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	trs, count, err := clickh.GetTransactionInfo(&apiPb.GetTransactionsRequest{
		ApplicationId: "ApplicationId",
		Pagination:    nil,
		TimeRange:     nil,
		Type:          0,
		Status:        0,
		Host:          nil,
		Name:          nil,
		Path:          nil,
		Method:        nil,
		Sort:          nil,
	})

	assert.Equal(t, "ApplicationId", trs[0].ApplicationId)
	assert.Equal(t, "ApplicationId", trs[1].ApplicationId)
	assert.Equal(t, count, int64(2))
}

func TestGetTransactionChildren(t *testing.T) {
	tr1 := &apiPb.TransactionInfo{
		Id:            "GetTransactionChildren",
		ApplicationId: "ApplicationId",
		ParentId:      "ParentId",
		Meta: &apiPb.TransactionInfo_Meta{
			Host:   "Host",
			Path:   "Path",
			Method: "Method",
		},
		Name: "Name",
		StartTime: &timestamp.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
		EndTime: &timestamp.Timestamp{
			Seconds: 6789,
			Nanos:   0,
		},
		Status: 1,
		Type:   1,
		Error: &apiPb.TransactionInfo_Error{
			Message: "TransactionInfo_Error",
		},
	}

	tr2 := &apiPb.TransactionInfo{
		Id:            "GetTransactionChildren2",
		ApplicationId: "ApplicationId",
		ParentId:      "GetTransactionChildren",
		Meta: &apiPb.TransactionInfo_Meta{
			Host:   "Host",
			Path:   "Path",
			Method: "Method",
		},
		Name: "Name",
		StartTime: &timestamp.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
		EndTime: &timestamp.Timestamp{
			Seconds: 6789,
			Nanos:   0,
		},
		Status: 1,
		Type:   1,
		Error: &apiPb.TransactionInfo_Error{
			Message: "TransactionInfo_Error",
		},
	}

	tr3 := &apiPb.TransactionInfo{
		Id:            "GetTransactionChildren3",
		ApplicationId: "ApplicationId",
		ParentId:      "GetTransactionChildren",
		Meta: &apiPb.TransactionInfo_Meta{
			Host:   "Host",
			Path:   "Path",
			Method: "Method",
		},
		Name: "Name",
		StartTime: &timestamp.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
		EndTime: &timestamp.Timestamp{
			Seconds: 6789,
			Nanos:   0,
		},
		Status: 1,
		Type:   1,
		Error: &apiPb.TransactionInfo_Error{
			Message: "TransactionInfo_Error",
		},
	}

	err := clickh.InsertTransactionInfo(tr1)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	err = clickh.InsertTransactionInfo(tr2)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	err = clickh.InsertTransactionInfo(tr3)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	trChildren, err := clickh.GetTransactionChildren(tr1.Id, "")

	assert.Equal(t, "GetTransactionChildren2", trChildren[0].TransactionId)
	assert.Equal(t, "GetTransactionChildren3", trChildren[1].TransactionId)
}

func TestGetTransactionGroup(t *testing.T) {
	tr1 := &apiPb.TransactionInfo{
		Id:            "GetTransactionGroup",
		ApplicationId: "ApplicationId",
		ParentId:      "ParentId",
		Meta: &apiPb.TransactionInfo_Meta{
			Host:   "Host",
			Path:   "Path",
			Method: "Method",
		},
		Name: "Name",
		StartTime: &timestamp.Timestamp{
			Seconds: 6789,
			Nanos:   0,
		},
		EndTime: &timestamp.Timestamp{
			Seconds: 6899,
			Nanos:   0,
		},
		Status: 1,
		Type:   5,
		Error: &apiPb.TransactionInfo_Error{
			Message: "TransactionInfo_Error",
		},
	}

	tr2 := &apiPb.TransactionInfo{
		Id:            "GetTransactionGroup",
		ApplicationId: "ApplicationId",
		ParentId:      "ParentId",
		Meta: &apiPb.TransactionInfo_Meta{
			Host:   "Host",
			Path:   "Path",
			Method: "Method",
		},
		Name: "Name",
		StartTime: &timestamp.Timestamp{
			Seconds: 6789,
			Nanos:   0,
		},
		EndTime: &timestamp.Timestamp{
			Seconds: 6899,
			Nanos:   0,
		},
		Status: 1,
		Type:   5,
		Error: &apiPb.TransactionInfo_Error{
			Message: "TransactionInfo_Error",
		},
	}

	tr3 := &apiPb.TransactionInfo{
		Id:            "GetTransactionGroup",
		ApplicationId: "ApplicationId",
		ParentId:      "ParentId",
		Meta: &apiPb.TransactionInfo_Meta{
			Host:   "Host",
			Path:   "Path",
			Method: "Method",
		},
		Name: "Name",
		StartTime: &timestamp.Timestamp{
			Seconds: 6789,
			Nanos:   0,
		},
		EndTime: &timestamp.Timestamp{
			Seconds: 6899,
			Nanos:   0,
		},
		Status: 1,
		Type:   5,
		Error: &apiPb.TransactionInfo_Error{
			Message: "TransactionInfo_Error",
		},
	}

	err := clickh.InsertTransactionInfo(tr1)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	err = clickh.InsertTransactionInfo(tr2)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	err = clickh.InsertTransactionInfo(tr3)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	trGroup, err := clickh.GetTransactionGroup(&apiPb.GetTransactionGroupRequest{
		ApplicationId: "ApplicationId",
		TimeRange: &apiPb.TimeFilter{
			From: &timestamp.Timestamp{
				Seconds: 6700,
				Nanos:   0,
			},
			To: &timestamp.Timestamp{
				Seconds: 6900,
				Nanos:   0,
			},
		},
		GroupType: apiPb.GroupTransaction_BY_PATH,
		Type:      apiPb.TransactionType_TRANSACTION_TYPE_GRPC,
		Status:    1,
	})

	assert.Equal(t, 3, trGroup["Path"].Count)
}

func TestClickhouse_Migrate_error(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		err := clickhWrong.Migrate()
		assert.Error(t, err)
	})
}

type CustomConverter struct{}

func (s CustomConverter) ConvertValue(v interface{}) (driver.Value, error) {
	switch v.(type) {
	case clickhouse.UUID:
		return v.(clickhouse.UUID), nil
	case string:
		return v.(string), nil
	case []uint32:
		return v.([]uint32), nil
	case []int64:
		return v.([]int64), nil
	case int:
		return v.(int), nil
	case int32:
		return v.(int32), nil
	case int64:
		return v.(int64), nil
	case time.Time:
		return v.(time.Time), nil
	default:
		return nil, errors.New(fmt.Sprintf("cannot convert %T with value %v", v, v))
	}
}
