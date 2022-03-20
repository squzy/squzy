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
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"regexp"
	"testing"
	"time"
)

var (
	wdbTransactionInfo, _      = sql.Open("clickhouse", "tcp://user:password@lkl:00/debug=true&clicks?read_timeout=10&write_timeout=10")
	clickhWrongTransactionInfo = &Clickhouse{
		wdbTransactionInfo,
	}
	clickTransactionInfo = &Clickhouse{}
)

type SuiteTransInfo struct {
	suite.Suite
	DB   *sql.DB
	mock sqlmock.Sqlmock
}

func (s *SuiteTransInfo) SetupSuite() {
	var err error

	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	clickTransactionInfo.Db = s.DB
}

func (s *SuiteTransInfo) Test_InsertTransactionInfo() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbTransactionInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	
	s.mock.ExpectCommit()

	correctTime := timestamp.Now()
	err := clickTransactionInfo.InsertTransactionInfo(&apiPb.TransactionInfo{
		StartTime: correctTime,
		EndTime:   correctTime,
		Error: &apiPb.TransactionInfo_Error{
			Message: "d",
		},
	})
	require.NoError(s.T(), err)
}

func TestClickhouse_InsertTransactionInfo(t *testing.T) {
	t.Run("Should: return conv error", func(t *testing.T) {
		err := clickTransactionInfo.InsertTransactionInfo(&apiPb.TransactionInfo{})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		correctTime := timestamp.Now()
		err := clickhWrongTransactionInfo.InsertTransactionInfo(&apiPb.TransactionInfo{
			StartTime: correctTime,
			EndTime:   correctTime,
		})
		assert.Error(t, err)
	})
}

func (s *SuiteTransInfo) Test_GetTransactionInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, transactionInfoFields, dbTransactionInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "updated_at", "transaction_id", "application_id", "parent_id", "meta_host", "meta_path", "meta_method", "name", "start_time", "end_time", "transaction_status", "transaction_type", "error"}).
		AddRow("1", time.Now(), time.Now(), "1", "0", "1", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickTransactionInfo.GetTransactionInfo(
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

func (s *SuiteTransInfo) Test_GetTransactionInfo_CountError() {
	var (
		id = "1"
	)

	_, _, err := clickTransactionInfo.GetTransactionInfo(
		&apiPb.GetTransactionsRequest{
			ApplicationId: id,
			Type:          1,
			Status:        1,
			Host:          &wrappers.StringValue{Value: "q"},
		})
	require.Error(s.T(), err)
}

func (s *SuiteTransInfo) Test_GetTransactionInfo_SelectError() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickTransactionInfo.GetTransactionInfo(
		&apiPb.GetTransactionsRequest{
			ApplicationId: id,
			Type:          1,
			Status:        1,
			Host:          &wrappers.StringValue{Value: "q"},
		})
	require.Error(s.T(), err)
}

func TestClickhouse_GetTransactionInfo(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := clickhWrongTransactionInfo.GetTransactionInfo(
			&apiPb.GetTransactionsRequest{
				TimeRange: &apiPb.TimeFilter{
					From: &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
					To:   &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
				},
			})
		assert.Error(t, err)
	})
}

func (s *SuiteTransInfo) Test_GetTransactionByID() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT %s FROM "%s"`, transactionInfoFields, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "transaction_id", "application_id", "parent_id", "meta_host", "meta_path", "meta_method", "name", "start_time", "end_time", "transaction_status", "transaction_type", "error"}).
		AddRow("1", time.Now(), time.Now(), "1", "0", "1", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, transactionInfoFields, dbTransactionInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "updated_at", "transaction_id", "application_id", "parent_id", "meta_host", "meta_path", "meta_method", "name", "start_time", "end_time", "transaction_status", "transaction_type", "error"}).
		AddRow("1", time.Now(), time.Now(), "1", "0", "1", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickTransactionInfo.GetTransactionByID(
		&apiPb.GetTransactionByIdRequest{
			TransactionId: id,
		})
	require.NoError(s.T(), err)
}

func (s *SuiteTransInfo) Test_GetTransactionByID_ChildrenError() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT %s FROM "%s"`, transactionInfoFields, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "transaction_id", "application_id", "parent_id", "meta_host", "meta_path", "meta_method", "name", "start_time", "end_time", "transaction_status", "transaction_type", "error"}).
		AddRow("1", time.Now(), time.Now(), "1", "0", "1", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickTransactionInfo.GetTransactionByID(
		&apiPb.GetTransactionByIdRequest{
			TransactionId: id,
		})
	require.Error(s.T(), err)
}

func (s *SuiteTransInfo) Test_GetTransactionByID_Error() {
	var (
		id = "1"
	)

	_, _, err := clickTransactionInfo.GetTransactionByID(
		&apiPb.GetTransactionByIdRequest{
			TransactionId: id,
		})
	require.Error(s.T(), err)
}

func (s *SuiteTransInfo) Test_GetTransactionChildren() {
	query := fmt.Sprintf(`SELECT %s FROM "%s"`, transactionInfoFields, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "transaction_id", "application_id", "parent_id", "meta_host", "meta_path", "meta_method", "name", "start_time", "end_time", "transaction_status", "transaction_type", "error"}).
		AddRow("1", time.Now(), time.Now(), "1", "0", "100", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, transactionInfoFields, dbTransactionInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "updated_at", "transaction_id", "application_id", "parent_id", "meta_host", "meta_path", "meta_method", "name", "start_time", "end_time", "transaction_status", "transaction_type", "error"}).
		AddRow("2", time.Now(), time.Now(), "2", "0", "1", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, transactionInfoFields, dbTransactionInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "updated_at", "transaction_id", "application_id", "parent_id", "meta_host", "meta_path", "meta_method", "name", "start_time", "end_time", "transaction_status", "transaction_type", "error"}).
		AddRow("3", time.Now(), time.Now(), "1", "0", "100", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := clickTransactionInfo.GetTransactionChildren("1", "")
	require.NoError(s.T(), err)
}

func (s *SuiteTransInfo) Test_GetTransactionChildren_Error() {
	var (
		id = "1"
	)

	_, err := clickTransactionInfo.GetTransactionChildren(id, "")
	require.Error(s.T(), err)
}

func (s *SuiteTransInfo) Test_GetTransactionChildren_SubchildrenError() {
	query := fmt.Sprintf(`SELECT %s FROM "%s"`, transactionInfoFields, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "transaction_id", "application_id", "parent_id", "meta_host", "meta_path", "meta_method", "name", "start_time", "end_time", "transaction_status", "transaction_type", "error"}).
		AddRow("1", time.Now(), time.Now(), "1", "0", "100", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, transactionInfoFields, dbTransactionInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "updated_at", "transaction_id", "application_id", "parent_id", "meta_host", "meta_path", "meta_method", "name", "start_time", "end_time", "transaction_status", "transaction_type", "error"}).
		AddRow("1", time.Now(), time.Now(), "2", "0", "100", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := clickTransactionInfo.GetTransactionChildren("0", "")
	require.Error(s.T(), err)
}

func (s *SuiteTransInfo) Test_GetTransactionGroup() {
	query := fmt.Sprintf(
		`SELECT "%s"."name" as "groupName", COUNT("%s"."name") as "count", COUNT(CASE WHEN "transaction_info"."transaction_status" = '1' THEN 1 ELSE NULL END) as "successCount", AVG("%s"."end_time"-"%s"."start_time") as "latency", min("transaction_info"."end_time"-"transaction_info"."start_time") as "minTime", max("transaction_info"."end_time"-"transaction_info"."start_time") as "maxTime", min("%s"."end_time") as "lowTime"`,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"groupName", "count", "latency"})
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := clickTransactionInfo.GetTransactionGroup(&apiPb.GetTransactionGroupRequest{
		ApplicationId: "1",
		GroupType:     2,
	})
	require.NoError(s.T(), err)
}

func (s *SuiteTransInfo) Test_GetTransactionGroup_Error() {
	_, err := clickTransactionInfo.GetTransactionGroup(&apiPb.GetTransactionGroupRequest{
		ApplicationId: "1",
		GroupType:     -2,
	})
	require.Error(s.T(), err)
}

func TestClickhouse_GetTransactionGroup(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, err := clickhWrongTransactionInfo.GetTransactionGroup(
			&apiPb.GetTransactionGroupRequest{
				TimeRange: &apiPb.TimeFilter{
					From: &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
					To:   &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
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
		assert.Equal(t, " desc", res)
	})
}

func Test_getTransactionOrder(t *testing.T) {
	//Time for invalid timestamp
	t.Run("Should: return empty string", func(t *testing.T) {
		res := getTransactionOrder(&apiPb.SortingTransactionList{
			SortBy:    -1,
			Direction: -1,
		})
		assert.Equal(t, fmt.Sprintf(`"%s"."start_time"`, dbTransactionInfoCollection), res)
	})
}

func (s *SuiteTransInfo) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInitTransInfo(t *testing.T) {
	suite.Run(t, new(SuiteTransInfo))
}
