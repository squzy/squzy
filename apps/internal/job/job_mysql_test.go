package job

import (
	"context"
	"database/sql"
	"database/sql/driver"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mysqlDriverConn struct {}

func (mysqlDriverConn) Prepare(query string) (driver.Stmt, error) {
	return nil, nil
}

func (mysqlDriverConn) Close() error {
	return nil
}

func (mysqlDriverConn) Begin() (driver.Tx, error) {
	return nil, nil
}

type mysqlDbConnector struct {}

func (mysqlDbConnector) Connect(context.Context) (driver.Conn, error) {
	return &mysqlDriverConn{}, nil
}
func (mysqlDbConnector) Driver() driver.Driver {
	return nil
}

type sqlMockConnectOk struct{}

func (sqlMockConnectOk) open(name, info string) (*sql.DB, error) {
	return sql.OpenDB(&mysqlDbConnector{}), nil
}

func TestMysqlJob_Do(t *testing.T) {
	t.Run("Test: mysqlJob", func(t *testing.T) {
		t.Run("Should: return error connecting", func(t *testing.T) {
			j := NewMysqlJob("", 0, "", "", "")
			err := j.Do()
			expected := clientPb.StatusCode_Error
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return no error", func(t *testing.T) {
			j := mysqlJob{
				mysql:    &sqlMockConnectOk{},
			}
			err := j.Do()
			expected := clientPb.StatusCode_OK
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
	})
}
