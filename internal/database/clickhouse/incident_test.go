package clickhouse

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"regexp"
	"testing"
	"time"
)

var (
	dbIncidentCollection        = "incidents"
	dbIncidentHistoryCollection = "incidents_history"

	wdbIncident, _      = sql.Open("clickhouse", "tcp://user:password@lkl:00/debug=true&clicks?read_timeout=10&write_timeout=10")
	clickhWrongIncident = &Clickhouse{
		wdbIncident,
	}
	clickIncident = &Clickhouse{}
)

type SuiteIncident struct {
	suite.Suite
	DB   *sql.DB
	mock sqlmock.Sqlmock
}

func (s *SuiteIncident) SetupSuite() {
	var err error

	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	clickIncident.Db = s.DB
}

type CustomStruct struct{}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteIncident) Test_UpdateIncidentStatus_InsertError_getIncidentById() {
	query := fmt.Sprintf(`SELECT %s FROM %s`, incidentFields, dbIncidentCollection)
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnError(errors.New("Test_UpdateIncidentStatus_InsertError_getIncidentById"))

	_, err := clickIncident.UpdateIncidentStatus("", apiPb.IncidentStatus_INCIDENT_STATUS_OPENED)
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_UpdateIncidentStatus_InsertError_updateIncident() {
	query := fmt.Sprintf(`SELECT %s FROM %s`, incidentFields, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "incident_id", "status", "rule_id", "start_time", "end_time"}).
		AddRow("1", time.Now(), time.Now(), "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnRows(rows)
	query = fmt.Sprintf(`SELECT %s FROM %s`, incidentHistoriesFields, dbIncidentHistoryCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "incident_id", "status", "timestamp"}).
		AddRow("1", time.Now(), time.Now(), "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnRows(rows)

	s.mock.ExpectBegin()
	query = fmt.Sprintf(`INSERT INTO %s (%s)`, dbIncidentCollection, incidentFields)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("Test_UpdateIncidentStatus_InsertError_updateIncident"))

	_, err := clickIncident.UpdateIncidentStatus("", apiPb.IncidentStatus_INCIDENT_STATUS_OPENED)
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_UpdateIncidentStatus_InsertError_insertIncidentHistory() {
	query := fmt.Sprintf(`SELECT %s FROM %s`, incidentFields, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "incident_id", "status", "rule_id", "start_time", "end_time"}).
		AddRow("1", time.Now(), time.Now(), "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnRows(rows)
	query = fmt.Sprintf(`SELECT %s FROM %s`, incidentHistoriesFields, dbIncidentHistoryCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "incident_id", "status", "timestamp"}).
		AddRow("1", time.Now(), time.Now(), "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnRows(rows)

	s.mock.ExpectBegin()
	query = fmt.Sprintf(`INSERT INTO %s (%s)`, dbIncidentCollection, incidentFields)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	s.mock.ExpectBegin()
	query = fmt.Sprintf(`INSERT INTO %s (%s)`, dbIncidentHistoryCollection, incidentHistoriesFields)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("Test_UpdateIncidentStatus_InsertError_updateIncident"))

	_, err := clickIncident.UpdateIncidentStatus("", apiPb.IncidentStatus_INCIDENT_STATUS_OPENED)
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_UpdateIncidentStatus_InsertError_GetIncidentById() {
	query := fmt.Sprintf(`SELECT %s FROM %s`, incidentFields, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "incident_id", "status", "rule_id", "start_time", "end_time"}).
		AddRow("1", time.Now(), time.Now(), "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnRows(rows)
	query = fmt.Sprintf(`SELECT %s FROM %s`, incidentHistoriesFields, dbIncidentHistoryCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "incident_id", "status", "timestamp"}).
		AddRow("1", time.Now(), time.Now(), "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnRows(rows)

	s.mock.ExpectBegin()
	query = fmt.Sprintf(`INSERT INTO %s (%s)`, dbIncidentCollection, incidentFields)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	s.mock.ExpectBegin()
	query = fmt.Sprintf(`INSERT INTO %s (%s)`, dbIncidentHistoryCollection, incidentHistoriesFields)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	query = fmt.Sprintf(`SELECT %s FROM %s`, incidentFields, dbIncidentCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "updated_at", "incident_id", "status", "rule_id", "start_time", "end_time"}).
		AddRow("1", time.Now(), time.Now(), "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnError(errors.New("Test_UpdateIncidentStatus_InsertError_GetIncidentById"))

	_, err := clickIncident.UpdateIncidentStatus("", apiPb.IncidentStatus_INCIDENT_STATUS_OPENED)
	require.Error(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteIncident) Test_getIncident() {
	query := fmt.Sprintf(`SELECT %s FROM %s`, incidentFields, dbIncidentCollection)

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New("Test_getIncident"))

	_, err := clickIncident.getIncident("")
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_getIncident_next() {
	query := fmt.Sprintf(`SELECT %s FROM %s`, incidentFields, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{})

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	res, err := clickIncident.getIncident("")
	require.Nil(s.T(), res)
	require.Nil(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteIncident) Test_UpdateIncidentStatus_Scan() {
	query := fmt.Sprintf(`SELECT %s FROM %s`, incidentHistoriesFields, dbIncidentHistoryCollection)
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New(""))

	_, err := clickIncident.getIncidentHistories("")
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_GetIncidentById_Error() {
	_, err := clickIncident.GetIncidentById("")
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_GetActiveIncidentByRuleId_error() {
	_, err := clickIncident.GetActiveIncidentByRuleId("1")
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_GetIncidents_timeError() {
	maxValidSeconds := 253402300800
	_, _, err := clickIncident.GetIncidents(&apiPb.GetIncidentsListRequest{
		RuleId: &wrappers.StringValue{Value: "id"},
		TimeRange: &apiPb.TimeFilter{
			From: &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		},
	})
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_GetIncidents_countError() {
	_, _, err := clickIncident.GetIncidents(&apiPb.GetIncidentsListRequest{
		RuleId: &wrappers.StringValue{Value: "id"},
	})
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_GetIncidents_selectIncidentError() {
	query := fmt.Sprintf(`SELECT count(*) FROM %s`, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickIncident.GetIncidents(&apiPb.GetIncidentsListRequest{
		RuleId: &wrappers.StringValue{Value: "id"},
	})
	require.Error(s.T(), err)
}

func (s *SuiteIncident) Test_GetIncidents_selectIncidentHistoryError() {
	query := fmt.Sprintf(`SELECT count(*) FROM %s`, dbIncidentCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM %s`, incidentFields, dbIncidentCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickIncident.GetIncidents(&apiPb.GetIncidentsListRequest{
		RuleId: &wrappers.StringValue{Value: "id"},
	})
	require.Error(s.T(), err)
}

func Test_getIncidentRuleString(t *testing.T) {
	t.Run("Should: return string", func(t *testing.T) {
		res := getIncidentRuleString("id")
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

func (s *SuiteIncident) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInitIncident(t *testing.T) {
	suite.Run(t, new(SuiteIncident))
}
