package database

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
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

func (s *Suite) Test_InsertMetaData() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbSnapshotCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "meta_data")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

func (s *Suite) Test_GetMetaData() {
	var (
		id = "1"
	)
	query := fmt.Sprintf(`SELECT * FROM "%s" WHERE "%s"."deleted_at" IS NULL`, dbSnapshotCollection, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnRows(rows)

	_, err := postgr.GetSnapshots(id)
	require.NoError(s.T(), err)
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
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "memories")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "memories")).
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
		WithArgs(id).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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
		WithArgs(id).
		WillReturnRows(rows)

	_, _, err := postgr.GetStatRequest(id, &apiPb.Pagination{
		Page:  1, //random value
		Limit: 2, //random value
	}, nil)
	require.Error(s.T(), err)
}

func (s *Suite) Test_GetCpuInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT cpuInfo, time FROM "%s"`, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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
		WithArgs(id).
		WillReturnRows(rows)

	_, _, err := postgr.GetCPUInfo(id, &apiPb.Pagination{
		Page:  1, //random value
		Limit: 2, //random value
	}, nil)
	require.Error(s.T(), err)
}

func (s *Suite) Test_GetMemoryInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT memoryInfo, time FROM "%s"`, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetMemoryInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetDiskInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT diskInfo, time FROM "%s"`, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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
		WithArgs(id).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT netInfo, time FROM "%s"`, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgr.GetNetInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func TestPostgres_Migrate(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		err := postgrWrong.Migrate()
		assert.Error(t, err)
	})
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

func TestPostgres_GetMetaData(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := postgrWrong.GetSnapshots("")
		assert.Error(t, err)
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
