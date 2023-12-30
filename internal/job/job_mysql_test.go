package job

import (
	"errors"
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
				db:       &DBConnection{},
			}
			err := ExecMysql(j.dbConfig, j.db)
			expected := apiPb.SchedulerCode_ERROR
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return error ping", func(t *testing.T) {
			j := mysqlJob{
				dbConfig: &scheduler_config_storage.DbConfig{},
				db:       &dbMockErr{},
			}
			err := ExecMysql(j.dbConfig, j.db)
			expected := apiPb.SchedulerCode_ERROR
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return no error", func(t *testing.T) {
			j := mysqlJob{
				dbConfig: &scheduler_config_storage.DbConfig{},
				db:       &dbMockOk{},
			}
			err := ExecMysql(j.dbConfig, j.db)
			expected := apiPb.SchedulerCode_OK
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
	})
}
