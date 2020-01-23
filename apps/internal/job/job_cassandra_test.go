package job

import (
	"errors"
	"github.com/gocql/gocql"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockCassandraToolsCreateError struct {}

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

type mockCassandraToolsExecuteError struct {}

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

type mockCassandraToolsOk struct {}

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

func TestNewCassandraJob(t *testing.T) {
	t.Run("Should: create job", func(t *testing.T) {
		assert.NotNil(t, NewCassandraJob("", "", ""))
	})
}

func TestCassandraJob_Do(t *testing.T) {
	t.Run("Test: cassandra job", func(t *testing.T) {
		t.Run("Should: return error create session", func(t *testing.T) {
			j := cassandraJob{
				cassandraTools: &mockCassandraToolsCreateError{},
			}
			err := j.Do()
			expected := clientPb.StatusCode_Error
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return error execute batch", func(t *testing.T) {
			j := cassandraJob{
				cassandraTools: &mockCassandraToolsExecuteError{},
			}
			err := j.Do()
			expected := clientPb.StatusCode_Error
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return no error", func(t *testing.T) {
			j := cassandraJob{
				cassandraTools: &mockCassandraToolsOk{},
			}
			err := j.Do()
			expected := clientPb.StatusCode_OK
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
	})
}
