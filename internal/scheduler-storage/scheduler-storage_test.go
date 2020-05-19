package scheduler_storage

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type schedulerMock struct {
}

func (s schedulerMock) GetIDBson() primitive.ObjectID {
	panic("implement me")
}

func (s schedulerMock) GetID() string {
	return "1"
}

func (s schedulerMock) Run() {
	panic("implement me")
}

func (s schedulerMock) Stop() {
}

func (s schedulerMock) IsRun() bool {
	return true
}

func TestNew(t *testing.T) {
	t.Run("Shoudle: create storage", func(t *testing.T) {
		s := New()
		assert.Implements(t, (*SchedulerStorage)(nil), s)
	})
}

func TestStorage_Get(t *testing.T) {
	t.Run("Should: return error errNotExistError", func(t *testing.T) {
		s := New()
		err := s.Set(&schedulerMock{})
		assert.Equal(t, nil, err)
		_, err = s.Get("0")
		assert.Equal(t, errNotExistError, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New()
		err := s.Set(&schedulerMock{})
		assert.Equal(t, nil, err)
		_, err = s.Get("1")
		assert.Equal(t, nil, err)
	})
}

func TestStorage_Set(t *testing.T) {
	t.Run("Should: return error errStorageKeyAlreadyExistError", func(t *testing.T) {
		s := New()
		err := s.Set(&schedulerMock{})
		assert.Equal(t, nil, err)
		err = s.Set(&schedulerMock{})
		assert.Equal(t, errStorageKeyAlreadyExistError, err)
	})
	t.Run("Should: return error errNotExistError", func(t *testing.T) {
		s := New()
		err := s.Set(&schedulerMock{})
		assert.Equal(t, nil, err)
		_, err = s.Get("0")
		assert.Equal(t, errNotExistError, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New()
		err := s.Set(&schedulerMock{})
		assert.Equal(t, nil, err)
		_, err = s.Get("1")
		assert.Equal(t, nil, err)
	})
}

func TestStorage_Remove(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New()
		err := s.Set(&schedulerMock{})
		assert.Equal(t, nil, err)
		err = s.Remove("0")
		assert.Equal(t, nil, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New()
		err := s.Set(&schedulerMock{})
		assert.Equal(t, nil, err)
		err = s.Remove("1")
		assert.Equal(t, nil, err)
	})
}
