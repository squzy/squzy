package job

import (
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgresJob_Do(t *testing.T) {
	t.Run("Test: postgresDbJob", func(t *testing.T) {
		t.Run("Should: return error connecting", func(t *testing.T) {
			j := postgresJob{
				dbConfig: &scheduler_config_storage.DbConfig{},
				db:       NewDBConnection(),
			}
			err := ExecPostgres(j.schedulerID, j.dbConfig, j.db)
			expected := apiPb.SchedulerCode_ERROR
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return error ping", func(t *testing.T) {
			j := postgresJob{
				dbConfig: &scheduler_config_storage.DbConfig{},
				db:       &dbMockErr{},
			}
			err := ExecPostgres(j.schedulerID, j.dbConfig, j.db)
			expected := apiPb.SchedulerCode_ERROR
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return no error", func(t *testing.T) {
			j := postgresJob{
				dbConfig: &scheduler_config_storage.DbConfig{},
				db:       &dbMockOk{},
			}
			err := ExecPostgres(j.schedulerID, j.dbConfig, j.db)
			expected := apiPb.SchedulerCode_OK
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
	})
}
