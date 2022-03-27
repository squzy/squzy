package clickhouse

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
)

var (
	click = &Clickhouse{}
)

type SuiteClickhouse struct {
	suite.Suite
	DB   *sql.DB
	mock sqlmock.Sqlmock
}

func (s *SuiteClickhouse) SetupSuite() {
	var err error

	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	click.Db = s.DB
}

func (s *SuiteClickhouse) Test_Migrate_1() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbSnapshotCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := click.Migrate()
	require.Error(s.T(), err)
}

func (s *SuiteClickhouse) Test_Migrate_2() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbSnapshotCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := click.Migrate()
	require.Error(s.T(), err)
}

func (s *SuiteClickhouse) Test_Migrate_3() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbSnapshotCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCpuInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := click.Migrate()
	require.Error(s.T(), err)
}

func (s *SuiteClickhouse) Test_Migrate_4() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbSnapshotCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCpuInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := click.Migrate()
	require.Error(s.T(), err)
}

func (s *SuiteClickhouse) Test_Migrate_5() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbSnapshotCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCpuInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoMemCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := click.Migrate()
	require.Error(s.T(), err)
}

func (s *SuiteClickhouse) Test_Migrate_6() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbSnapshotCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCpuInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoMemCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoSwapCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := click.Migrate()
	require.Error(s.T(), err)
}

func (s *SuiteClickhouse) Test_Migrate_7() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbSnapshotCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCpuInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoMemCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoSwapCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestDiskInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := click.Migrate()
	require.Error(s.T(), err)
}

func (s *SuiteClickhouse) Test_Migrate_8() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbSnapshotCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCpuInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoMemCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoSwapCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestDiskInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestNetInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := click.Migrate()
	require.Error(s.T(), err)
}

func (s *SuiteClickhouse) Test_Migrate_9() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbSnapshotCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCpuInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoMemCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoSwapCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestDiskInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestNetInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbTransactionInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := click.Migrate()
	require.Error(s.T(), err)
}

func (s *SuiteClickhouse) Test_Migrate_10() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbSnapshotCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCpuInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoMemCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoSwapCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestDiskInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestNetInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbTransactionInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbIncidentCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := click.Migrate()
	require.Error(s.T(), err)
}

func (s *SuiteClickhouse) Test_Migrate_11() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbSnapshotCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCpuInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoMemCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoSwapCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestDiskInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestNetInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbTransactionInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbIncidentCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := click.Migrate()
	require.Error(s.T(), err)
}

func (s *SuiteClickhouse) Test_Migrate_12() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbSnapshotCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestCpuInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoMemCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestMemoryInfoSwapCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestDiskInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbStatRequestNetInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbTransactionInfoCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbIncidentCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, dbIncidentHistoryCollection)
	s.mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := click.Migrate()
	require.Nil(s.T(), err)
}

//func (s *SuiteClickhouse) AfterTest(_, _ string) {
//	require.NoError(s.T(), s.mock.ExpectationsWereMet())
//}

func TestInitClickhouse(t *testing.T) {
	suite.Run(t, new(SuiteClickhouse))
}
