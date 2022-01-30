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
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

var (
	postgrTransInfo = &Postgres{}
	dbTransInfo, _  = gorm.Open(
		postgres.Open(fmt.Sprintf("host=lkl port=00 user=us dbname=dbn password=ps connect_timeout=10 sslmode=disable")),
		&gorm.Config{},
	)
	postgrWrongTransInfo = &Postgres{
		dbTransInfo,
	}
)

type SuiteTransInfo struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *SuiteTransInfo) SetupSuite() {
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
	postgrTransInfo.Db = s.DB
}

func (s *SuiteTransInfo) Test_InsertTransactionInfo() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbTransactionInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	correctTime := timestamp.Now()
	err := postgrTransInfo.InsertTransactionInfo(&apiPb.TransactionInfo{
		StartTime: correctTime,
		EndTime:   correctTime,
	})
	require.NoError(s.T(), err)
}

func TestPostgres_InsertTransactionInfo(t *testing.T) {
	t.Run("Should: return conv error", func(t *testing.T) {
		err := postgrTransInfo.InsertTransactionInfo(&apiPb.TransactionInfo{})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		correctTime := timestamp.Now()
		err := postgrWrongTransInfo.InsertTransactionInfo(&apiPb.TransactionInfo{
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

	query = fmt.Sprintf(`SELECT * FROM "%s"`, dbTransactionInfoCollection)
	rows = sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgrTransInfo.GetTransactionInfo(
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
func (s *SuiteTransInfo) Test_GetTransactionInfo_CountError() {
	var (
		id = "1"
	)

	_, _, err := postgrTransInfo.GetTransactionInfo(
		&apiPb.GetTransactionsRequest{
			ApplicationId: id,
			Type:          1,
			Status:        1,
			Host:          &wrappers.StringValue{Value: "q"},
		})
	require.Error(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteTransInfo) Test_GetTransactionInfo_SelectError() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgrTransInfo.GetTransactionInfo(
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
		_, _, err := postgrWrongTransInfo.GetTransactionInfo(
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

	_, _, err := postgrTransInfo.GetTransactionByID(
		&apiPb.GetTransactionByIdRequest{
			TransactionId: id,
		})
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteTransInfo) Test_GetTransactionByID_ChildrenError() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT * FROM "%s"`, dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"transactionId", "parentId"}).AddRow("1", "0")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := postgrTransInfo.GetTransactionByID(
		&apiPb.GetTransactionByIdRequest{
			TransactionId: id,
		})
	require.Error(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteTransInfo) Test_GetTransactionByID_Error() {
	var (
		id = "1"
	)

	_, _, err := postgrTransInfo.GetTransactionByID(
		&apiPb.GetTransactionByIdRequest{
			TransactionId: id,
		})
	require.Error(s.T(), err)
}

func (s *SuiteTransInfo) Test_GetTransactionChildren() {
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

	_, err := postgrTransInfo.GetTransactionChildren("0", "")
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteTransInfo) Test_GetTransactionChildren_Error() {
	var (
		id = "1"
	)

	_, err := postgrTransInfo.GetTransactionChildren(id, "")
	require.Error(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteTransInfo) Test_GetTransactionChildren_SubchildrenError() {
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

	_, err := postgrTransInfo.GetTransactionChildren("0", "")
	require.Error(s.T(), err)
}

func (s *SuiteTransInfo) Test_GetTransactionGroup() {
	query := fmt.Sprintf(
		`SELECT "%s"."name" as "groupName", COUNT("%s"."name") as "count", COUNT(CASE WHEN "transaction_infos"."transactionStatus" = '1' THEN 1 ELSE NULL END) as "successCount", AVG("%s"."endTime"-"%s"."startTime") as "latency", min("transaction_infos"."endTime"-"transaction_infos"."startTime") as "minTime", max("transaction_infos"."endTime"-"transaction_infos"."startTime") as "maxTime", min("%s"."endTime") as "lowTime"`,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection)
	rows := sqlmock.NewRows([]string{"groupName", "count", "latency"})
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := postgrTransInfo.GetTransactionGroup(&apiPb.GetTransactionGroupRequest{
		ApplicationId: "1",
		GroupType:     2,
	})
	require.NoError(s.T(), err)
}

//Based on fact, that if request is not mocked, it will return error
func (s *SuiteTransInfo) Test_GetTransactionGroup_Error() {
	_, err := postgrTransInfo.GetTransactionGroup(&apiPb.GetTransactionGroupRequest{
		ApplicationId: "1",
		GroupType:     -2,
	})
	require.Error(s.T(), err)
}

func TestPostgres_GetTransactionGroup(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, err := postgrWrongTransInfo.GetTransactionGroup(
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
		assert.Equal(t, fmt.Sprintf(`"%s"."startTime"`, dbTransactionInfoCollection), res)
	})
}

func (s *SuiteTransInfo) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInitTransInfo(t *testing.T) {
	suite.Run(t, new(SuiteTransInfo))
}
