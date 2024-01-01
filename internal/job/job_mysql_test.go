package job

import (
	"database/sql"
	"errors"
	"fmt"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

type dbMockOk struct {
}

func (m dbMockOk) Connect(string, string) error {
	return nil
}

func (m dbMockOk) Ping() error {
	return nil
}

func (m dbMockOk) Close() error {
	return nil
}

type dbMockErr struct {
}

func (m dbMockErr) Connect(string, string) error {
	return nil
}

func (m dbMockErr) Ping() error {
	return errors.New("Ping")
}

func (m dbMockErr) Close() error {
	return nil
}

func TestMysqlJob_Do(t *testing.T) {
	t.Run("Test: mysqlJob", func(t *testing.T) {
		t.Run("Should: return error connecting", func(t *testing.T) {
			j := mysqlJob{
				dbConfig: &scheduler_config_storage.DbConfig{},
				db:       NewDBConnection(),
			}
			err := ExecMysql(j.schedulerID, j.dbConfig, j.db)
			expected := apiPb.SchedulerCode_ERROR
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return error ping", func(t *testing.T) {
			j := mysqlJob{
				dbConfig: &scheduler_config_storage.DbConfig{},
				db:       &dbMockErr{},
			}
			err := ExecMysql(j.schedulerID, j.dbConfig, j.db)
			expected := apiPb.SchedulerCode_ERROR
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return no error", func(t *testing.T) {
			j := mysqlJob{
				dbConfig: &scheduler_config_storage.DbConfig{},
				db:       &dbMockOk{},
			}
			err := ExecMysql(j.schedulerID, j.dbConfig, j.db)
			expected := apiPb.SchedulerCode_OK
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
	})
}

type dbConnectionOk struct {
}

func (d dbConnectionOk) Ping() error {
	return nil
}

func (d dbConnectionOk) Close() error {
	return nil
}

var openMock = func(driverName, dataSourceName string) (*sql.DB, error) {
	return &sql.DB{}, nil
}

func TestDBConnection(t *testing.T) {
	t.Run("Test: DBConnection", func(t *testing.T) {
		t.Run("Should: get Connect", func(t *testing.T) {
			config := &scheduler_config_storage.DbConfig{}
			m := DBConnection{Open: openMock}
			args := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				config.Host, config.Port, config.User, config.Password, config.DbName)

			err := m.Connect("mysql", args)
			assert.Nil(t, err)
			assert.NotNil(t, t, m.Client)
		})
		t.Run("Should: Ping", func(t *testing.T) {
			m := DBConnection{
				Client: &dbConnectionOk{},
			}

			err := m.Ping()
			assert.Nil(t, err)
		})
		t.Run("Should: Close", func(t *testing.T) {
			m := DBConnection{
				Client: &dbConnectionOk{},
			}

			err := m.Close()
			assert.Nil(t, err)
		})
	})
}
