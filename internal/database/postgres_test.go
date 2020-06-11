package database

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jinzhu/gorm"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

//docker run -d --rm --name postgres -e POSTGRES_USER="user" -e POSTGRES_PASSWORD="password" -e POSTGRES_DB="database" -p 5432:5432 postgres
var (
	postgr = &postgres{}
	db, _  = gorm.Open(
		"postgres",
		fmt.Sprintf("host=lkl port=00 user=us dbname=dbn password=ps connect_timeout=10 sslmode=disable"))
	postgrWrong = &postgres{
		db,
	}
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	postgr.db = s.DB

	s.DB.LogMode(true)
}

func TestPostgres_Migrate_error(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		err := postgrWrong.Migrate()
		assert.Error(t, err)
	})
}

func (s *Suite) Test_Snapshots() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbSnapshotCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	correctTime, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		require.NotNil(s.T(), nil)
	}
	err = postgr.InsertSnapshot(&apiPb.SchedulerResponse{
		SchedulerId: "schId",
		Snapshot: &apiPb.SchedulerSnapshot{
			Code: 0,
			Type: 0,
			Error: &apiPb.SchedulerSnapshot_Error{
				Message: "message",
			},
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: correctTime,
				EndTime:   correctTime,
				Value:     nil,
			},
		},
	})
	require.NoError(s.T(), err)
}

func TestPostgres_InsertSnapshots(t *testing.T) {
	t.Run("Should: return conv error", func(t *testing.T) {
		err := postgr.InsertSnapshot(&apiPb.SchedulerResponse{})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		correctTime, err := ptypes.TimestampProto(time.Now())
		if err != nil {
			assert.NotNil(t, nil)
		}
		err = postgrWrong.InsertSnapshot(&apiPb.SchedulerResponse{
			SchedulerId: "",
			Snapshot: &apiPb.SchedulerSnapshot{
				Code: 0,
				Type: 0,
				Meta: &apiPb.SchedulerSnapshot_MetaData{
					StartTime: correctTime,
					EndTime:   correctTime,
				},
			},
		})
		assert.Error(t, err)
	})
}

func (s *Suite) Test_GetSnapshots() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbSnapshotCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)


	_, _, err := postgr.GetSnapshots(&apiPb.GetSchedulerInformationRequest{
		SchedulerId: id,
		Sort: &apiPb.SortingSchedulerList{
			SortBy:    -1,
			Direction: -1,
		},
	})
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetSnapshots_WithStatus() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbSnapshotCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)


	_, _, err := postgr.GetSnapshots(&apiPb.GetSchedulerInformationRequest{
		SchedulerId: id,
		Sort: &apiPb.SortingSchedulerList{
			SortBy:    apiPb.SortSchedulerList_BY_LATENCY,
			Direction: apiPb.SortDirection_ASC,
		},
		Status: apiPb.SchedulerCode_OK,
	})
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *Suite) Test_GetSnapshots_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetSnapshots(
		&apiPb.GetSchedulerInformationRequest{
			SchedulerId: id,
			Pagination: &apiPb.Pagination{
				Page:  -1, //random value
				Limit: 2,  //random value
			},
		})
	require.Error(s.T(), err)
}

func TestPostgres_GetSnapshots(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := postgrWrong.GetSnapshots(
			&apiPb.GetSchedulerInformationRequest{
				SchedulerId: "",
				Pagination:  nil,
				TimeRange: &apiPb.TimeFilter{
					From: &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
					To:   &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
				},
			})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := postgrWrong.GetSnapshots(&apiPb.GetSchedulerInformationRequest{})
		assert.Error(t, err)
	})
}

func (s *Suite) Test_GetSnapshotsUptime() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`COUNT(*) as "count", AVG("%s"."metaEndTime"-"%s"."metaStartTime") as "latency"`, dbSnapshotCollection, dbSnapshotCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)


	_, err := postgr.GetSnapshotsUptime(&apiPb.GetSchedulerUptimeRequest{
		SchedulerId: id,
	})
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *Suite) Test_GetSnapshotsUptime_FirstCountError() {
	var (
		id = "1"
	)

	_, err := postgr.GetSnapshotsUptime(&apiPb.GetSchedulerUptimeRequest{
		SchedulerId: id,
	})
	require.Error(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *Suite) Test_GetSnapshotsUptime_SelectError() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := postgr.GetSnapshotsUptime(&apiPb.GetSchedulerUptimeRequest{
		SchedulerId: id,
	})
	require.Error(s.T(), err)
}

func TestPostgres_GetSnapshotsUptime(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, err := postgrWrong.GetSnapshotsUptime(
			&apiPb.GetSchedulerUptimeRequest{
				SchedulerId: "",
				TimeRange: &apiPb.TimeFilter{
					From: &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
					To:   &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
				},
			})
		assert.Error(t, err)
	})
}

func Test_getUptimeAndLatency(t *testing.T) {
	t.Run("Should: return 0 and no error", func(t *testing.T) {
		var snapshots []*Snapshot
		uptime, latency, err := getUptimeAndLatency(snapshots, 0, 0)
		assert.Equal(t, float64(0), uptime)
		assert.Equal(t, float64(0), latency)
		assert.NoError(t, err)
	})
	t.Run("Should: return not 0 and no error", func(t *testing.T) {
		snapshots := []*Snapshot{
			{
				Code: "OK",
				MetaStartTime: time.Now().UnixNano(),
				MetaEndTime: time.Now().UnixNano(),
			},
		}
		uptime, latency, err := getUptimeAndLatency(snapshots, 1, 1)
		assert.Equal(t, float64(1), uptime)
		assert.Equal(t, float64(0), latency)
		assert.NoError(t, err)
	})
}

func TestPostgres_InsertStatRequest(t *testing.T) {
	t.Run("Should: return conv error", func(t *testing.T) {
		err := postgr.InsertStatRequest(&apiPb.Metric{})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		err := postgrWrong.InsertStatRequest(&apiPb.Metric{
			CpuInfo: &apiPb.CpuInfo{
				Cpus: []*apiPb.CpuInfo_CPU{{}},
			},
			MemoryInfo: &apiPb.MemoryInfo{
				Mem:  &apiPb.MemoryInfo_Memory{},
				Swap: &apiPb.MemoryInfo_Memory{},
			},
			DiskInfo: &apiPb.DiskInfo{
				Disks: map[string]*apiPb.DiskInfo_Disk{},
			},
			NetInfo: &apiPb.NetInfo{
				Interfaces: map[string]*apiPb.NetInfo_Interface{},
			},
			Time: ptypes.TimestampNow(),
		})
		assert.Error(t, err)
	})
}

func (s *Suite) Test_InsertStatRequest() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "cpu_infos")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "memory_infos")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "memory_mems")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "memory_swaps")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "disk_infos")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "net_infos")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	err := postgr.InsertStatRequest(&apiPb.Metric{
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
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetStatRequest() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cpu_infos"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "memory_infos"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "memory_mems"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "memory_swaps"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "disk_infos"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "net_infos"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetStatRequest(id, nil, nil)
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *Suite) Test_GetStatRequest_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetStatRequest(id, &apiPb.Pagination{
		Page:  1, //random value
		Limit: 2, //random value
	}, nil)
	require.Error(s.T(), err)
}

func TestPostgres_GetStatRequest(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := postgrWrong.GetStatRequest("", nil, &apiPb.TimeFilter{
			From: &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := postgrWrong.GetStatRequest("", nil, nil)
		assert.Error(t, err)
	})
}

func (s *Suite) Test_GetCpuInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cpu_infos"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetCPUInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
//Is used for getSpecialRecords test
func (s *Suite) Test_GetCpuInfo_Count_Error() {
	var (
		id = "1"
	)

	_, _, err := postgr.GetCPUInfo(id, nil, nil)
	require.Error(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
//Is used for getSpecialRecords test
func (s *Suite) Test_GetCpuInfo_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetCPUInfo(id, &apiPb.Pagination{
		Page:  1, //random value
		Limit: 2, //random value
	}, nil)
	require.Error(s.T(), err)
}

//Time errors in getSpecialRecords
func TestPostgres_GetCpuInfo(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := postgrWrong.GetCPUInfo("", nil, &apiPb.TimeFilter{
			From: &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		})
		assert.Error(t, err)
	})
}

func (s *Suite) Test_GetMemoryInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "memory_infos"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetMemoryInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *Suite) Test_GetMemoryInfo_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetMemoryInfo(id, &apiPb.Pagination{
		Page:  1, //random value
		Limit: 2, //random value
	}, nil)
	require.Error(s.T(), err)
}

func TestPostgres_GetMemoryInfo(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := postgrWrong.GetMemoryInfo("", nil, &apiPb.TimeFilter{
			From: &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := postgrWrong.GetMemoryInfo("", nil, nil)
		assert.Error(t, err)
	})
}

func (s *Suite) Test_GetDiskInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "disk_infos"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetDiskInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetNetInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "net_infos"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetNetInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

func (s *Suite) Test_InsertTransactionInfo() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbTransactionInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	correctTime, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		require.NotNil(s.T(), nil)
	}
	err = postgr.InsertTransactionInfo(&apiPb.TransactionInfo{
		StartTime: correctTime,
		EndTime:   correctTime,
	})
	require.NoError(s.T(), err)
}

func TestPostgres_InsertTransactionInfo(t *testing.T) {
	t.Run("Should: return conv error", func(t *testing.T) {
		err := postgr.InsertTransactionInfo(&apiPb.TransactionInfo{})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		correctTime, err := ptypes.TimestampProto(time.Now())
		if err != nil {
			assert.NotNil(t, nil)
		}
		err = postgrWrong.InsertTransactionInfo(&apiPb.TransactionInfo{
			StartTime: correctTime,
			EndTime:   correctTime,
		})
		assert.Error(t, err)
	})
}

func (s *Suite) Test_GetTransactionInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbTransactionInfoCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetTransactionInfo(
		&apiPb.GetTransactionsRequest{
			ApplicationId: id,
			Host:          &wrappers.StringValue{Value: "q"},
			Sort: &apiPb.SortingTransactionList{
				SortBy:    0,
				Direction: 0,
			},
		})
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *Suite) Test_GetTransactionInfo_CountError() {
	var (
		id = "1"
	)

	_, _, err := postgr.GetTransactionInfo(
		&apiPb.GetTransactionsRequest{
			ApplicationId: id,
			Type:          1,
			Status:        1,
			Host:          &wrappers.StringValue{Value: "q"},
		})
	require.Error(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *Suite) Test_GetTransactionInfo_SelectError() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetTransactionInfo(
		&apiPb.GetTransactionsRequest{
			ApplicationId: id,
			Type:          1,
			Status:        1,
			Host:          &wrappers.StringValue{Value: "q"},
		})
	require.Error(s.T(), err)
}

func TestPostgres_GetTransactionInfo(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := postgrWrong.GetTransactionInfo(
			&apiPb.GetTransactionsRequest{
				TimeRange: &apiPb.TimeFilter{
					From: &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
					To:   &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
				},
			})
		assert.Error(t, err)
	})
}

func (s *Suite) Test_GetTransactionByID() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT * FROM "%s"`, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"transactionId", "parentId"}).AddRow("1", "0")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbTransactionInfoCollection)
	rows = sqlmock.NewRows([]string{"transactionId", "parentId"})
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetTransactionByID(
		&apiPb.GetTransactionByIdRequest{
			TransactionId: id,
		})
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *Suite) Test_GetTransactionByID_ChildrenError() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT * FROM "%s"`, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"transactionId", "parentId"}).AddRow("1", "0")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetTransactionByID(
		&apiPb.GetTransactionByIdRequest{
			TransactionId: id,
		})
	require.Error(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *Suite) Test_GetTransactionByID_Error() {
	var (
		id = "1"
	)

	_, _, err := postgr.GetTransactionByID(
		&apiPb.GetTransactionByIdRequest{
			TransactionId: id,
		})
	require.Error(s.T(), err)
}

func (s *Suite) Test_GetTransactionChildren() {
	query := fmt.Sprintf(`SELECT * FROM "%s"`, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"transactionId", "parentId"}).AddRow("1", "0")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbTransactionInfoCollection)
	rows = sqlmock.NewRows([]string{"transactionId", "parentId"}).AddRow("2", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	//Check for sycle
	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbTransactionInfoCollection)
	rows = sqlmock.NewRows([]string{"transactionId", "parentId"}).AddRow("1", "0")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := postgr.GetTransactionChildren("0", "")
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *Suite) Test_GetTransactionChildren_Error() {
	var (
		id = "1"
	)

	_, err := postgr.GetTransactionChildren(id, "")
	require.Error(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *Suite) Test_GetTransactionChildren_SubchildrenError() {
	query := fmt.Sprintf(`SELECT * FROM "%s"`, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"transactionId", "parentId"}).AddRow("1", "0")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbTransactionInfoCollection)
	rows = sqlmock.NewRows([]string{"transactionId", "parentId"}).AddRow("2", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := postgr.GetTransactionChildren("0", "")
	require.Error(s.T(), err)
}

func (s *Suite) Test_GetTransactionGroup() {
	query := fmt.Sprintf(
		`SELECT "%s"."name" as "groupName", COUNT("%s"."name") as "count", COUNT(CASE WHEN "transaction_infos"."transactionStatus" = 'TRANSACTION_SUCCESSFUL' THEN 1 ELSE NULL END) as "successCount", AVG("%s"."endTime"-"%s"."startTime") as "latency", min("transaction_infos"."endTime"-"transaction_infos"."startTime") as "minTime", max("transaction_infos"."endTime"-"transaction_infos"."startTime") as "maxTime", min("%s"."endTime") as "lowTime"`,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"groupName", "count", "latency"})
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := postgr.GetTransactionGroup(&apiPb.GetTransactionGroupRequest{
		ApplicationId: "1",
		GroupType:     2,
	})
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *Suite) Test_GetTransactionGroup_Error() {
	_, err := postgr.GetTransactionGroup(&apiPb.GetTransactionGroupRequest{
		ApplicationId: "1",
		GroupType:     -2,
	})
	require.Error(s.T(), err)
}

func TestPostgres_GetTransactionGroup(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, err := postgrWrong.GetTransactionGroup(
			&apiPb.GetTransactionGroupRequest{
				TimeRange: &apiPb.TimeFilter{
					From: &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
					To:   &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
				},
			})
		assert.Error(t, err)
	})
}

func Test_getTransactionDirection(t *testing.T) {
	//Time for invalid timestamp
	t.Run("Should: return empty string", func(t *testing.T) {
		res := getTransactionDirection(&apiPb.SortingTransactionList{
			SortBy:    -1,
			Direction: -1,
		})
		assert.Equal(t, "", res)
	})
}

func Test_getTransactionOrder(t *testing.T) {
	//Time for invalid timestamp
	t.Run("Should: return empty string", func(t *testing.T) {
		res := getTransactionOrder(&apiPb.SortingTransactionList{
			SortBy:    -1,
			Direction: -1,
		})
		assert.Equal(t, fmt.Sprintf(`"%s"."startTime"`, dbTransactionInfoCollection), res)
	})
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}
