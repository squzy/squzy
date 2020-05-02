package scheduler_storage

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type schedulerStopErrorMock struct {
}

func (s schedulerStopErrorMock) GetIdBson() primitive.ObjectID {
	panic("implement me")
}

func (s schedulerStopErrorMock) GetId() string {
	return "1"
}

func (s schedulerStopErrorMock) Run() {
}

func (s schedulerStopErrorMock) Stop() {
	return
}

func (s schedulerStopErrorMock) IsRun() bool {
	panic("implement me")
}

type schedulerMock struct {
}

func (s schedulerMock) GetIdBson() primitive.ObjectID {
	panic("implement me")
}

func (s schedulerMock) GetId() string {
	return "1"
}

func (s schedulerMock) Run() {
	panic("implement me")
}

func (s schedulerMock) Stop() {
	return
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
	t.Run("Should: return error notExistError", func(t *testing.T) {
		s := New()
		err := s.Set(&schedulerMock{})
		assert.Equal(t, nil, err)
		_, err = s.Get("0")
		assert.Equal(t, notExistError, err)
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
	t.Run("Should: return error storageKeyAlreadyExistError", func(t *testing.T) {
		s := New()
		err := s.Set(&schedulerMock{})
		assert.Equal(t, nil, err)
		err = s.Set(&schedulerMock{})
		assert.Equal(t, storageKeyAlreadyExistError, err)
	})
	t.Run("Should: return error notExistError", func(t *testing.T) {
		s := New()
		err := s.Set(&schedulerMock{})
		assert.Equal(t, nil, err)
		_, err = s.Get("0")
		assert.Equal(t, notExistError, err)
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
	t.Run("Should: return error storageKeyNotExistError", func(t *testing.T) {
		s := New()
		err := s.Set(&schedulerMock{})
		assert.Equal(t, nil, err)
		err = s.Remove("0")
		assert.Equal(t, storageKeyNotExistError, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New()
		err := s.Set(&schedulerMock{})
		assert.Equal(t, nil, err)
		err = s.Remove("1")
		assert.Equal(t, nil, err)
	})
}
