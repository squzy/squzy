package postgres

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

var (
	postgrSnapshot = &Postgres{}
	dbSnapshot, _  = gorm.Open(
		postgres.Open(fmt.Sprintf("host=lkl port=00 user=us dbname=dbn password=ps connect_timeout=10 sslmode=disable")),
		&gorm.Config{},
	)
	postgrWrongSnapshot = &Postgres{
		dbSnapshot,
	}
)

type SuiteSnapshot struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *SuiteSnapshot) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	require.NoError(s.T(), err)
	postgrSnapshot.Db = s.DB
}

func (s *SuiteSnapshot) Test_Snapshots() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbSnapshotCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	correctTime := timestamp.Now()

	err := postgrSnapshot.InsertSnapshot(&apiPb.SchedulerResponse{
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
		err := postgrSnapshot.InsertSnapshot(&apiPb.SchedulerResponse{})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		correctTime := timestamp.Now()
		err := postgrWrongSnapshot.InsertSnapshot(&apiPb.SchedulerResponse{
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

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbSnapshotCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgrSnapshot.GetSnapshots(&apiPb.GetSchedulerInformationRequest{
		SchedulerId: id,
		Sort: &apiPb.SortingSchedulerList{
			SortBy:    -1,
			Direction: -1,
		},
	})
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

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbSnapshotCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgrSnapshot.GetSnapshots(&apiPb.GetSchedulerInformationRequest{
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
func (s *SuiteSnapshot) Test_GetSnapshots_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgrSnapshot.GetSnapshots(
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
		_, _, err := postgrWrongSnapshot.GetSnapshots(
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
		_, _, err := postgrWrongSnapshot.GetSnapshots(&apiPb.GetSchedulerInformationRequest{})
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

	query = fmt.Sprintf(`COUNT(*) as "count", AVG("%s"."metaEndTime"-"%s"."metaStartTime") as "latency"`, dbSnapshotCollection, dbSnapshotCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := postgrSnapshot.GetSnapshotsUptime(&apiPb.GetSchedulerUptimeRequest{
		SchedulerId: id,
	})
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteSnapshot) Test_GetSnapshotsUptime_FirstCountError() {
	var (
		id = "1"
	)

	_, err := postgrSnapshot.GetSnapshotsUptime(&apiPb.GetSchedulerUptimeRequest{
		SchedulerId: id,
	})
	require.Error(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteSnapshot) Test_GetSnapshotsUptime_SelectError() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := postgrSnapshot.GetSnapshotsUptime(&apiPb.GetSchedulerUptimeRequest{
		SchedulerId: id,
	})
	require.Error(s.T(), err)
}

func TestPostgres_GetSnapshotsUptime(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, err := postgrWrongSnapshot.GetSnapshotsUptime(
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
