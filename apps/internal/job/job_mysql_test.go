package job

import (
	"database/sql"
	"errors"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockMySqlError(string, string) (*sql.DB, error) {
	return nil, errors.New("ERROR")
}

func mockMySqlOk(string, string) (*sql.DB, error) {
	return sql.OpenDB(nil), nil
}

func mockMySqlPingError(*sql.DB) error {
	return errors.New("ERROR")
}

func mockMySqlPingOk(*sql.DB) error {
	return nil
}

func TestNewMysqlJob(t *testing.T) {
	t.Run("Should: create job", func(t *testing.T) {
		assert.NotNil(t, NewMysqlJob("", 0, "", "", "", func(db *sql.DB) error {return nil}))
	})
}

func TestMysqlJob_Do(t *testing.T) {
	t.Run("Test: mysqlJob", func(t *testing.T) {
		t.Run("Should: return error connecting", func(t *testing.T) {
			j := mysqlJob{
				mySqlOpen: mockMySqlError,
			}
			err := j.Do()
			expected := clientPb.StatusCode_Error
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return error ping", func(t *testing.T) {
			j := mysqlJob{
				mySqlOpen: mockMySqlOk,
				mySqlPing: mockMySqlPingError,
			}
			err := j.Do()
			expected := clientPb.StatusCode_Error
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return no error", func(t *testing.T) {
			j := mysqlJob{
				mySqlOpen: mockMySqlOk,
				mySqlPing: mockMySqlPingOk,
			}
			err := j.Do()
			expected := clientPb.StatusCode_OK
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
	})
}
