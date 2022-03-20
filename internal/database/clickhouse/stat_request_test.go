package clickhouse

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"regexp"
	"testing"
	"time"
)

var (
	wdbStatRequest, _      = sql.Open("clickhouse", "tcp://user:password@lkl:00/debug=true&clicks?read_timeout=10&write_timeout=10")
	clickhWrongStatRequest = &Clickhouse{
		wdbStatRequest,
	}
	clickStatRequest = &Clickhouse{}
)

type SuiteStatRequest struct {
	suite.Suite
	DB   *sql.DB
	mock sqlmock.Sqlmock
}

func (s *SuiteStatRequest) SetupSuite() {
	var err error

	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	clickStatRequest.Db = s.DB
}

func TestPostgres_InsertStatRequest(t *testing.T) {
	t.Run("Should: return conv error", func(t *testing.T) {
		err := clickStatRequest.InsertStatRequest(&apiPb.Metric{})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		err := clickhWrongStatRequest.InsertStatRequest(&apiPb.Metric{
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
			Time: timestamp.Now(),
		})
		assert.Error(t, err)
	})
}

func (s *SuiteStatRequest) Test_InsertStatRequest() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCpuInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestMemoryInfoMemCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestMemoryInfoSwapCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestDiskInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestNetInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := clickStatRequest.InsertStatRequest(&apiPb.Metric{
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
		Time: timestamp.Now(),
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

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsCpuInfoFields, dbStatRequestCpuInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "load"}).
		AddRow("1", time.Now(), "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoMemCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoSwapCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsDiskInfoFields, dbStatRequestDiskInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "name", "total", "used", "free", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsNetInfoFields, dbStatRequestNetInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "name", "bytes_sent", "bytes_recv", "packets_sent", "packets_recv", "err_in", "err_out", "drop_in", "drop_out"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetStatRequest(id, nil, nil)
	require.NoError(s.T(), err)
}

func (s *SuiteStatRequest) Test_GetStatRequest_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetStatRequest(id, &apiPb.Pagination{
		Page:  1, //random value
		Limit: 2, //random value
	}, nil)
	require.Error(s.T(), err)
}

func TestPostgres_GetStatRequest(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := clickhWrongStatRequest.GetStatRequest("", nil, &apiPb.TimeFilter{
			From: &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := clickhWrongStatRequest.GetStatRequest("", nil, nil)
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

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsCpuInfoFields, dbStatRequestCpuInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "load"}).
		AddRow("1", time.Now(), "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetCPUInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

//Is used for getSpecialRecords test
func (s *SuiteStatRequest) Test_GetCpuInfo_Count_Error() {
	var (
		id = "1"
	)

	_, _, err := clickStatRequest.GetCPUInfo(id, nil, nil)
	require.Error(s.T(), err)
}

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

	_, _, err := clickStatRequest.GetCPUInfo(id, &apiPb.Pagination{
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
		_, _, err := clickhWrongStatRequest.GetCPUInfo("", nil, &apiPb.TimeFilter{
			From: &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
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

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoMemCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoSwapCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetMemoryInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

func (s *SuiteStatRequest) Test_GetMemoryInfo_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetMemoryInfo(id, &apiPb.Pagination{
		Page:  1, //random value
		Limit: 2, //random value
	}, nil)
	require.Error(s.T(), err)
}

func TestPostgres_GetMemoryInfo(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := clickhWrongStatRequest.GetMemoryInfo("", nil, &apiPb.TimeFilter{
			From: &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := clickhWrongStatRequest.GetMemoryInfo("", nil, nil)
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

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsDiskInfoFields, dbStatRequestDiskInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "name", "total", "used", "free", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetDiskInfo(id, nil, nil)
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

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsNetInfoFields, dbStatRequestNetInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "name", "bytes_sent", "bytes_recv", "packets_sent", "packets_recv", "err_in", "err_out", "drop_in", "drop_out"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetNetInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

func (s *SuiteStatRequest) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInitStatRequest(t *testing.T) {
	suite.Run(t, new(SuiteStatRequest))
}
