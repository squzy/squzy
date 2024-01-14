package cassandra_tools

import (
	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCassandraTools(t *testing.T) {
	t.Run("Test: NewCassandraTools", func(t *testing.T) {
		assert.NotNil(t, NewCassandraTools("", "", "", 0))
	})
}

func TestCassandraTools_CreateSession(t *testing.T) {
	t.Run("Test: CassandraTools.CreateSession", func(t *testing.T) {
		s := NewCassandraTools("", "", "", 0)
		_, err := s.CreateSession()
		assert.Error(t, err)
	})
}

func TestCassandraTools_ExecuteBatch(t *testing.T) {
	t.Run("Test: CassandraTools.ExecuteBatch", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		s := NewCassandraTools("", "", "", 0)
		_ = s.ExecuteBatch(&gocql.Session{}, &gocql.Batch{})
	})
}

func TestCassandraTools_NewBatch(t *testing.T) {
	t.Run("Test: CassandraTools.NewBatch", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		s := NewCassandraTools("", "", "", 0)
		_ = s.NewBatch(nil)
	})
}

func TestCassandraTools_Close(t *testing.T) {
	t.Run("Test: CassandraTools.Close", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		s := NewCassandraTools("", "", "", 0)
		s.Close(nil)
	})
}
