package scheduler

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	job_executor "github.com/squzy/squzy/internal/job-executor"
	"time"
)

var (
	errIntervalLessHalfSecondError = errors.New("INTERVAL_LESS_THAN_HALF_SECOND")
)

type Scheduler interface {
	// Should return id
	GetID() string
	//Get ID bson
	GetIDBson() primitive.ObjectID
	// Should run Scheduler every tick
	Run()
	// Should stop Scheduler
	Stop()
	// Return true/false depends from current state
	IsRun() bool
}

type schl struct {
	ticker      *time.Ticker
	isStopped   bool
	quitCh      chan bool
	interval    time.Duration
	id          primitive.ObjectID
	jobExecutor job_executor.JobExecutor
}

func New(id primitive.ObjectID, interval time.Duration, jobExecutor job_executor.JobExecutor) (Scheduler, error) {
	if interval < time.Millisecond*500 {
		return nil, errIntervalLessHalfSecondError
	}
	return &schl{
		id:          id,
		interval:    interval,
		isStopped:   true,
		jobExecutor: jobExecutor,
	}, nil
}

func (s *schl) Run() {
	if !s.isStopped {
		return
	}
	s.ticker = time.NewTicker(s.interval)
	s.isStopped = false
	s.quitCh = make(chan bool, 1)
	s.observer()
}

func (s *schl) observer() {
	go func() {
	loop:
		for {
			select {
			case <-s.ticker.C:
				s.jobExecutor.Execute(s.id)
			case <-s.quitCh:
				break loop
			}
		}
	}()
}

func (s *schl) IsRun() bool {
	return !s.isStopped
}

func (s *schl) GetID() string {
	return s.id.Hex()
}

func (s *schl) GetIDBson() primitive.ObjectID {
	return s.id
}

func (s *schl) Stop() {
	if s.isStopped {
		return
	}
	s.ticker.Stop()
	s.quitCh <- true
	close(s.quitCh)
	s.isStopped = true
}
