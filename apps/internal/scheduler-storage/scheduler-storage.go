package scheduler_storage

import (
	"errors"
	"squzy/apps/internal/scheduler"
	"sync"
)

var (
	notExistError               = errors.New("SCHEDULER_NOT_EXIST")
	storageKeyAlreadyExistError = errors.New("STORAGE_KEY_ALREADY_EXIST")
	storageKeyNotExistError     = errors.New("STORAGE_KEY_NOT_EXIST")
)

type SchedulerStorage interface {
	Get(string) (scheduler.Scheduler, error)
	Set(scheduler.Scheduler) error
	Remove(string) error
	GetList() map[string]bool
}

type storage struct {
	kv    map[string]scheduler.Scheduler
	mutex sync.RWMutex
}

func (s *storage) GetList() map[string]bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	statusMap := make(map[string]bool)
	for _, schl := range s.kv {
		statusMap[schl.GetId()] = schl.IsRun()
	}
	return statusMap
}

func (s *storage) Get(id string) (scheduler.Scheduler, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	value, exist := s.kv[id]
	if !exist {
		return nil, notExistError
	}
	return value, nil
}

func (s *storage) Set(schl scheduler.Scheduler) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	id := schl.GetId()
	_, exist := s.kv[id]

	if exist {
		return storageKeyAlreadyExistError
	}
	s.kv[id] = schl
	return nil
}

func (s *storage) Remove(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, exist := s.kv[id]
	if !exist {
		return storageKeyNotExistError
	}
	// make sure that observer stop before delete
	err := s.kv[id].Stop()
	if err != nil {
		// @TODO here we should add log to StdERR because that we should handle on monitoring side
		return err
	}
	delete(s.kv, id)
	return nil
}

func New() SchedulerStorage {
	return &storage{
		kv: make(map[string]scheduler.Scheduler),
	}
}
