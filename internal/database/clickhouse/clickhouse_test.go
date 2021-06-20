package clickhouse

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/docker/go-connections/nat"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/protobuf/types/known/structpb"
	"os"
	"sort"
	"squzy/internal/logger"
	"testing"
	"time"
)

var (
	db, _       = sql.Open("clickhouse", "tcp://user:password@lkl:00/debug=true&clicks?read_timeout=10&write_timeout=10")
	clickhWrong = &Clickhouse{
		db,
	}
	clickh        *Clickhouse
	testContainer testcontainers.Container
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	err := setup(ctx)
	if err != nil {
		logger.Fatalf("could not start test: %s", err)
	}
	code := m.Run()
	err = shutdown(ctx)
	if err != nil {
		logger.Fatalf("could not stop test: %s", err)
	}
	os.Exit(code)
}

func shutdown(ctx context.Context) error {
	err := testContainer.Terminate(ctx)
	if err != nil {
		return err
	}
	return nil
}

func setup(ctx context.Context) error {
	var err error
	req := testcontainers.ContainerRequest{
		Image:        "yandex/clickhouse-server",
		ExposedPorts: []string{"9000/tcp"},
		WaitingFor:   wait.ForListeningPort(nat.Port("9000/tcp")),
	}
	testContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return err
	}

	ip, err := testContainer.Host(ctx)
	if err != nil {
		return err
	}
	port, err := testContainer.MappedPort(ctx, "9000")
	if err != nil {
		return err
	}
	db, err = sql.Open("clickhouse", fmt.Sprintf("tcp://%s:%s?debug=true", ip, port.Port()))
	if err != nil {
		return err
	}
	clickh = &Clickhouse{
		db,
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
		SchedulerId: "getSnapshots",
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
		SchedulerId: "getSnapshots",
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
		SchedulerId: "getSnapshots",
		Snapshot: &apiPb.SchedulerSnapshot{
			Code:  apiPb.SchedulerCode_ERROR,
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
		SchedulerId: "getSnapshots",
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
