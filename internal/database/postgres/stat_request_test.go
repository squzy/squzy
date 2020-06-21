package postgres

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
)

var (
	postgrStatRequest = &Postgres{}
	dbStatRequest, _  = gorm.Open(
		"postgres",
		fmt.Sprintf("host=lkl port=00 user=us dbname=dbn password=ps connect_timeout=10 sslmode=disable"))
	postgrWrongStatRequest = &Postgres{
		dbStatRequest,
	}
)

type SuiteStatRequest struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *SuiteStatRequest) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	postgrStatRequest.Db = s.DB

	s.DB.LogMode(true)
}

func TestPostgres_InsertStatRequest(t *testing.T) {
	t.Run("Should: return conv error", func(t *testing.T) {
		err := postgrStatRequest.InsertStatRequest(&apiPb.Metric{})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		err := postgrWrongStatRequest.InsertStatRequest(&apiPb.Metric{
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

func (s *SuiteStatRequest) Test_InsertStatRequest() {
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

	err := postgrStatRequest.InsertStatRequest(&apiPb.Metric{
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

func (s *SuiteStatRequest) Test_GetStatRequest() {
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

	_, _, err := postgrStatRequest.GetStatRequest(id, nil, nil)
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteStatRequest) Test_GetStatRequest_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgrStatRequest.GetStatRequest(id, &apiPb.Pagination{
		Page:  1, //random value
		Limit: 2, //random value
	}, nil)
	require.Error(s.T(), err)
}

func TestPostgres_GetStatRequest(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := postgrWrongStatRequest.GetStatRequest("", nil, &apiPb.TimeFilter{
			From: &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := postgrWrongStatRequest.GetStatRequest("", nil, nil)
		assert.Error(t, err)
	})
}

func (s *SuiteStatRequest) Test_GetCpuInfo() {
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

	_, _, err := postgrStatRequest.GetCPUInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
//Is used for getSpecialRecords test
func (s *SuiteStatRequest) Test_GetCpuInfo_Count_Error() {
	var (
		id = "1"
	)

	_, _, err := postgrStatRequest.GetCPUInfo(id, nil, nil)
	require.Error(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
//Is used for getSpecialRecords test
func (s *SuiteStatRequest) Test_GetCpuInfo_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgrStatRequest.GetCPUInfo(id, &apiPb.Pagination{
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
		_, _, err := postgrWrongStatRequest.GetCPUInfo("", nil, &apiPb.TimeFilter{
			From: &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		})
		assert.Error(t, err)
	})
}

func (s *SuiteStatRequest) Test_GetMemoryInfo() {
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

	_, _, err := postgrStatRequest.GetMemoryInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteStatRequest) Test_GetMemoryInfo_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgrStatRequest.GetMemoryInfo(id, &apiPb.Pagination{
		Page:  1, //random value
		Limit: 2, //random value
	}, nil)
	require.Error(s.T(), err)
}

func TestPostgres_GetMemoryInfo(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := postgrWrongStatRequest.GetMemoryInfo("", nil, &apiPb.TimeFilter{
			From: &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := postgrWrongStatRequest.GetMemoryInfo("", nil, nil)
		assert.Error(t, err)
	})
}

func (s *SuiteStatRequest) Test_GetDiskInfo() {
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

	_, _, err := postgrStatRequest.GetDiskInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

func (s *SuiteStatRequest) Test_GetNetInfo() {
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

	_, _, err := postgrStatRequest.GetNetInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

func (s *SuiteStatRequest) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInitStatRequest(t *testing.T) {
	suite.Run(t, new(SuiteStatRequest))
}