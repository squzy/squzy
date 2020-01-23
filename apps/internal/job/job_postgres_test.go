package job

import (
	"database/sql"
	"errors"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockPostgresSqlError(string, string) (*sql.DB, error) {
	return nil, errors.New("ERROR")
}

func mockPostgresSqlOk(string, string) (*sql.DB, error) {
	return sql.OpenDB(nil), nil
}

func mockPostgresPingError(*sql.DB) error {
	return errors.New("ERROR")
}

func mockPostgresPingOk(*sql.DB) error {
	return nil
}

func TestNewPosgresDbJob(t *testing.T) {
	t.Run("Should: create job", func(t *testing.T) {
		assert.NotNil(t, NewPosgresDbJob("", 0, "", "", "", func(db *sql.DB) error {return nil}))
	})
}

func TestPostgresJob_Do(t *testing.T) {
	t.Run("Test: postgresDbJob", func(t *testing.T) {
		t.Run("Should: return error connecting", func(t *testing.T) {
			j := postgresJob{
				postgresOpen: mockPostgresSqlError,
			}
			err := j.Do()
			expected := clientPb.StatusCode_Error
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return error ping", func(t *testing.T) {
			j := postgresJob{
				postgresOpen: mockPostgresSqlOk,
				postgresPing: mockPostgresPingError,
			}
			err := j.Do()
			expected := clientPb.StatusCode_Error
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return no error", func(t *testing.T) {
			j := postgresJob{
				postgresOpen: mockPostgresSqlOk,
				postgresPing: mockPostgresPingOk,
			}
			err := j.Do()
			expected := clientPb.StatusCode_OK
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
	})
}
