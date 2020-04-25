package database

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
)

//docker run -d --rm --name postgres -e POSTGRES_USER="user" -e POSTGRES_PASSWORD="password" -e POSTGRES_DB="database" -p 5432:5432 postgres
var (
	postgr = &postgres{
		host:     "localhost",
		port:     "5432",
		user:     "user",
		password: "password",
		dbname:   "database",
	}
	postgrWrong = &postgres{
		host:     "localhost",
		port:     "5432",
		user:     "wrongUser",
		password: "wrongPassword",
		dbname:   "database",
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


func TestPostgres_NewClient(t *testing.T) {
	t.Run("wrongPostgress", func(t *testing.T) {
		err := postgrWrong.newClient()
		assert.Error(t, err)
	})
	/*t.Run("correctPostgress", func(t *testing.T) {
		err := postgr.newClient()
		assert.NoError(t, err)
	})*/
}

func (s *Suite) Test_InsertMetaData() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbMetaDataCollection)).
		WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	err := postgr.InsertMetaData(&MetaData{})
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetMetaData() {
	var (
	    id = "1"
	)
	query := fmt.Sprintf(`SELECT * FROM "%s" WHERE "%s"."deleted_at" IS NULL`, dbMetaDataCollection, dbMetaDataCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnRows(rows)

	_, err := postgr.GetMetaData(id)
	require.NoError(s.T(), err)
}

func (s *Suite) Test_InsertStatRequest() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCollection)).
		WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	err := postgr.InsertStatRequest(&StatRequest{})
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetStatRequest() {
	var (
		id = "1"
	)
	query := fmt.Sprintf(`SELECT * FROM "%s" WHERE "%s"."deleted_at" IS NULL`, dbStatRequestCollection, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnRows(rows)

	_, err := postgr.GetStatRequest(id)
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

func TestPostgres_InsertMetaData(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		err := postgrWrong.InsertMetaData(&MetaData{})
		assert.Error(t, err)
	})
}

func TestPostgres_GetMetaData(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := postgrWrong.GetMetaData("")
		assert.Error(t, err)
	})
}

func TestPostgres_InsertStatRequest(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		err := postgrWrong.InsertStatRequest(&StatRequest{})
		assert.Error(t, err)
	})
}

func TestPostgres_GetStatRequest(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := postgrWrong.GetStatRequest("")
		assert.Error(t, err)
	})
}
