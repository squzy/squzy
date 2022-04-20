package postgres

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"regexp"
	"testing"
)

var (
	postgrIncident = &Postgres{}
	dbIncident, _  = gorm.Open(
		"postgres",
		fmt.Sprintf("host=lkl port=00 user=us dbname=dbn password=ps connect_timeout=10 sslmode=disable"))
	postgrWrongIncident = &Postgres{
		dbIncident,
	}
)

type SuiteIncident struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *SuiteIncident) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	postgrIncident.Db = s.DB

	s.DB.LogMode(true)
}

func (s *SuiteIncident) Test_InsertIncident() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbIncidentCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbIncidentHistoryCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	err := postgrIncident.InsertIncident(&apiPb.Incident{
		Status: apiPb.IncidentStatus_INCIDENT_STATUS_CAN_BE_CLOSED,
		RuleId: "12345",
		Histories: []*apiPb.Incident_HistoryItem{
			{
				Status:    apiPb.IncidentStatus_INCIDENT_STATUS_OPENED,
				Timestamp: timestamp.Now(),
			},
		},
	})
	require.NoError(s.T(), err)
}

func TestPostgres_InsertIncident(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		err := postgrWrongIncident.InsertIncident(&apiPb.Incident{})
		assert.Error(t, err)
	})
}

func (s *SuiteIncident) Test_UpdateIncidentStatus() {
	query := fmt.Sprintf(`SELECT * FROM "%s"`, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbIncidentHistoryCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	s.mock.ExpectBegin()
	query = fmt.Sprintf(`UPDATE`)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	s.mock.ExpectBegin()
	query = fmt.Sprintf(`INSERT INTO "%s"`, dbIncidentHistoryCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)
	s.mock.ExpectCommit()

	_, err := postgrIncident.UpdateIncidentStatus("", apiPb.IncidentStatus_INCIDENT_STATUS_OPENED)
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteIncident) Test_UpdateIncidentStatus_SelectError() {
	_, err := postgrIncident.UpdateIncidentStatus("", apiPb.IncidentStatus_INCIDENT_STATUS_OPENED)
	require.Error(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteIncident) Test_UpdateIncidentStatus_UpdateError() {
	query := fmt.Sprintf(`SELECT * FROM "%s"`, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbIncidentHistoryCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := postgrIncident.UpdateIncidentStatus("", apiPb.IncidentStatus_INCIDENT_STATUS_OPENED)
	require.Error(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteIncident) Test_UpdateIncidentStatus_InsertError() {
	query := fmt.Sprintf(`SELECT * FROM "%s"`, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbIncidentHistoryCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	s.mock.ExpectBegin()
	query = fmt.Sprintf(`UPDATE`)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	_, err := postgrIncident.UpdateIncidentStatus("", apiPb.IncidentStatus_INCIDENT_STATUS_OPENED)
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_GetIncidentById() {
	query := fmt.Sprintf(`SELECT * FROM "%s"`, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbIncidentHistoryCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := postgrIncident.GetIncidentById("")
	require.NoError(s.T(), err)
}

func (s *SuiteIncident) Test_GetIncidentById_Error() {
	_, err := postgrIncident.GetIncidentById("")
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_GetActiveIncidentByRuleId() {
	query := fmt.Sprintf(`SELECT * FROM "%s"`, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbIncidentHistoryCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := postgrIncident.GetActiveIncidentByRuleId("10")
	require.NoError(s.T(), err)
}

func (s *SuiteIncident) Test_GetActiveIncidentByRuleId_error() {
	_, err := postgrIncident.GetActiveIncidentByRuleId("1")
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_GetIncidents() {
	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbIncidentCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbIncidentHistoryCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgrIncident.GetIncidents(&apiPb.GetIncidentsListRequest{})
	require.NoError(s.T(), err)
}

func (s *SuiteIncident) Test_GetIncidents_timeError() {
	maxValidSeconds := 253402300800
	_, _, err := postgrIncident.GetIncidents(&apiPb.GetIncidentsListRequest{
		TimeRange: &apiPb.TimeFilter{
			From: &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		},
	})
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_GetIncidents_countError() {
	_, _, err := postgrIncident.GetIncidents(&apiPb.GetIncidentsListRequest{})
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_GetIncidents_selectIncidentError() {
	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgrIncident.GetIncidents(&apiPb.GetIncidentsListRequest{})
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_GetIncidents_selectIncidentHistoryError() {
	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbIncidentCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgrIncident.GetIncidents(&apiPb.GetIncidentsListRequest{})
	require.Error(s.T(), err)
}

func Test_getIncidentRuleString(t *testing.T) {
	t.Run("Should: return string", func(t *testing.T) {
		res := getIncidentRuleString(&wrappers.StringValue{Value: "id"})
		assert.NotEqual(t, "", res)
	})
}

func Test_getIncidentOrder(t *testing.T) {
	t.Run("Should: return string", func(t *testing.T) {
		res := getIncidentOrder(&apiPb.SortingIncidentList{
			SortBy: apiPb.SortIncidentList_INCIDENT_LIST_BY_END_TIME,
		})
		assert.NotEqual(t, "", res)
	})
	t.Run("Should: return string", func(t *testing.T) {
		res := getIncidentOrder(&apiPb.SortingIncidentList{
			SortBy: 10,
		})
		assert.NotEqual(t, "", res)
	})
}

func Test_getIncidentDirection(t *testing.T) {
	t.Run("Should: return value from map", func(t *testing.T) {
		res := getIncidentDirection(&apiPb.SortingIncidentList{
			Direction: apiPb.SortDirection_DESC,
		})
		assert.NotEqual(t, "", res)
	})
	t.Run("Should: return default value", func(t *testing.T) {
		res := getIncidentDirection(&apiPb.SortingIncidentList{
			Direction: 10,
		})
		assert.NotEqual(t, "", res)
	})
}

func Test_checkNoFoundError(t *testing.T) {
	t.Run("Should: return value from map", func(t *testing.T) {
		_, err := checkNoFoundError(gorm.ErrRecordNotFound)
		assert.Nil(t, err)
	})
}

func (s *SuiteIncident) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInitIncident(t *testing.T) {
	suite.Run(t, new(SuiteIncident))
}
