package clickhouse

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"regexp"
	"testing"
	"time"
)

var (
	wdbSnapshot, _      = sql.Open("clickhouse", "tcp://user:password@lkl:00/debug=true&clicks?read_timeout=10&write_timeout=10")
	clickhWrongSnapshot = &Clickhouse{
		wdbSnapshot,
	}
	clickSnapshot = &Clickhouse{}
)

type SuiteSnapshot struct {
	suite.Suite
	DB   *sql.DB
	mock sqlmock.Sqlmock
}

func (s *SuiteSnapshot) SetupSuite() {
	var err error

	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	clickSnapshot.Db = s.DB
}

func (s *SuiteSnapshot) Test_Snapshots() {
	s.mock.ExpectBegin()
	query := fmt.Sprintf(`INSERT INTO "%s" (%s)`, dbSnapshotCollection, snapshotFields)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	correctTime := timestamp.Now()

	err := clickSnapshot.InsertSnapshot(&apiPb.SchedulerResponse{
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
				Value:     structpb.NewStringValue(""),
			},
		},
	})
	require.NoError(s.T(), err)
}

func TestClickhouse_InsertSnapshots(t *testing.T) {
	t.Run("Should: return conv error", func(t *testing.T) {
		err := clickSnapshot.InsertSnapshot(&apiPb.SchedulerResponse{})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		correctTime := timestamp.Now()
		err := clickhWrongSnapshot.InsertSnapshot(&apiPb.SchedulerResponse{
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

func (s *SuiteSnapshot) Test_GetSnapshots() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM %s`, snapshotFields, dbSnapshotCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "updated_at", "scheduler_id", "code", "type", "error", "meta_start_time", "meta_end_time", "meta_value"}).
		AddRow("1", time.Now(), time.Now(), "1", "1", "1", "error", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, count, err := clickSnapshot.GetSnapshots(&apiPb.GetSchedulerInformationRequest{
		SchedulerId: id,
		Sort: &apiPb.SortingSchedulerList{
			SortBy:    -1,
			Direction: -1,
		},
	})
	require.NotEqual(s.T(), int32(0), count)
	require.NoError(s.T(), err)
}

func (s *SuiteSnapshot) Test_GetSnapshots_WithStatus() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM %s`, snapshotFields, dbSnapshotCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "updated_at", "scheduler_id", "code", "type", "error", "meta_start_time", "meta_end_time", "meta_value"}).
		AddRow("1", time.Now(), time.Now(), "1", "1", "1", "error", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, count, err := clickSnapshot.GetSnapshots(&apiPb.GetSchedulerInformationRequest{
		SchedulerId: id,
		Sort: &apiPb.SortingSchedulerList{
			SortBy:    apiPb.SortSchedulerList_BY_LATENCY,
			Direction: apiPb.SortDirection_ASC,
		},
		Status: apiPb.SchedulerCode_OK,
	})
	require.NotEqual(s.T(), int32(2), count)
	require.NoError(s.T(), err)
}

func (s *SuiteSnapshot) Test_GetSnapshots_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickSnapshot.GetSnapshots(
		&apiPb.GetSchedulerInformationRequest{
			SchedulerId: id,
			Pagination: &apiPb.Pagination{
				Page:  -1, //random value
				Limit: 2,  //random value
			},
		})
	require.Error(s.T(), err)
}

func TestClickhouse_GetSnapshots(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := clickhWrongSnapshot.GetSnapshots(
			&apiPb.GetSchedulerInformationRequest{
				SchedulerId: "",
				Pagination:  nil,
				TimeRange: &apiPb.TimeFilter{
					From: &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
					To:   &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
				},
			})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := clickhWrongSnapshot.GetSnapshots(&apiPb.GetSchedulerInformationRequest{})
		assert.Error(t, err)
	})
}

func (s *SuiteSnapshot) Test_GetSnapshotsUptime() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT count(*) as "count", avg(meta_end_time-meta_start_time) as "latency" FROM "%s"`, dbSnapshotCollection)
	rows = sqlmock.NewRows([]string{"count", "latency"}).AddRow("1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := clickSnapshot.GetSnapshotsUptime(&apiPb.GetSchedulerUptimeRequest{
		SchedulerId: id,
	})
	require.NoError(s.T(), err)
}

func (s *SuiteSnapshot) Test_GetSnapshotsUptime_FirstCountError() {
	var (
		id = "1"
	)

	_, err := clickSnapshot.GetSnapshotsUptime(&apiPb.GetSchedulerUptimeRequest{
		SchedulerId: id,
	})
	require.Error(s.T(), err)
}

func (s *SuiteSnapshot) Test_GetSnapshotsUptime_SelectError() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := clickSnapshot.GetSnapshotsUptime(&apiPb.GetSchedulerUptimeRequest{
		SchedulerId: id,
	})
	require.Error(s.T(), err)
}

func TestClickhouse_GetSnapshotsUptime(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, err := clickhWrongSnapshot.GetSnapshotsUptime(
			&apiPb.GetSchedulerUptimeRequest{
				SchedulerId: "",
				TimeRange: &apiPb.TimeFilter{
					From: &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
					To:   &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
				},
			})
		assert.Error(t, err)
	})
}

func (s *SuiteSnapshot) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInitSnapshot(t *testing.T) {
	suite.Run(t, new(SuiteSnapshot))
}
