package job

import (
	"errors"
	"github.com/gocql/gocql"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockCassandraToolsCreateError struct{}

func (*mockCassandraToolsCreateError) CreateSession() (*gocql.Session, error) {
	return nil, errors.New("ERROR")
}
func (*mockCassandraToolsCreateError) ExecuteBatch(session *gocql.Session, batch *gocql.Batch) error {
	return nil
}
func (*mockCassandraToolsCreateError) NewBatch(session *gocql.Session) *gocql.Batch {
	return nil
}
func (*mockCassandraToolsCreateError) Close(session *gocql.Session) {}

type mockCassandraToolsExecuteError struct{}

func (*mockCassandraToolsExecuteError) CreateSession() (*gocql.Session, error) {
	return nil, nil
}
func (*mockCassandraToolsExecuteError) ExecuteBatch(session *gocql.Session, batch *gocql.Batch) error {
	return errors.New("ERROR")
}
func (*mockCassandraToolsExecuteError) NewBatch(session *gocql.Session) *gocql.Batch {
	return nil
}
func (*mockCassandraToolsExecuteError) Close(session *gocql.Session) {}

type mockCassandraToolsOk struct{}

func (*mockCassandraToolsOk) CreateSession() (*gocql.Session, error) {
	return &gocql.Session{}, nil
}
func (*mockCassandraToolsOk) ExecuteBatch(session *gocql.Session, batch *gocql.Batch) error {
	return nil
}
func (*mockCassandraToolsOk) NewBatch(session *gocql.Session) *gocql.Batch {
	return nil
}
func (*mockCassandraToolsOk) Close(session *gocql.Session) {}

func TestCassandraJob_Exec(t *testing.T) {
	t.Run("Test: cassandra job exec", func(t *testing.T) {
		t.Run("Should: return error create session", func(t *testing.T) {
			cassandraTools := &mockCassandraToolsCreateError{}
			err := ExecCassandra("", &scheduler_config_storage.DbConfig{}, cassandraTools)
			expected := apiPb.SchedulerCode_ERROR
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return error execute batch", func(t *testing.T) {
			cassandraTools := &mockCassandraToolsExecuteError{}
			err := ExecCassandra("", &scheduler_config_storage.DbConfig{}, cassandraTools)
			expected := apiPb.SchedulerCode_ERROR
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return no error", func(t *testing.T) {
			cassandraTools := &mockCassandraToolsOk{}
			err := ExecCassandra("", &scheduler_config_storage.DbConfig{}, cassandraTools)
			expected := apiPb.SchedulerCode_OK
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
	})
}
