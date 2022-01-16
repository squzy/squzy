package scheduler_storage

import (
	"errors"
	"github.com/squzy/squzy/internal/scheduler"
	"sync"
)

var (
	errNotExistError               = errors.New("SCHEDULER_NOT_EXIST")
	errStorageKeyAlreadyExistError = errors.New("STORAGE_KEY_ALREADY_EXIST")
)

type SchedulerStorage interface {
	Get(string) (scheduler.Scheduler, error)
	Set(scheduler.Scheduler) error
	Remove(string) error
}

type storage struct {
	kv    map[string]scheduler.Scheduler
	mutex sync.RWMutex
}

func (s *storage) Get(id string) (scheduler.Scheduler, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	value, exist := s.kv[id]
	if !exist {
		return nil, errNotExistError
	}
	return value, nil
}

func (s *storage) Set(schl scheduler.Scheduler) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	id := schl.GetID()
	_, exist := s.kv[id]

	if exist {
		return errStorageKeyAlreadyExistError
	}
	s.kv[id] = schl
	return nil
}

func (s *storage) Remove(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, exist := s.kv[id]
	if !exist {
		return nil
	}
	// make sure that observer stop before delete
	s.kv[id].Stop()
	delete(s.kv, id)
	return nil
}

func New() SchedulerStorage {
	return &storage{
		kv: make(map[string]scheduler.Scheduler),
	}
}
