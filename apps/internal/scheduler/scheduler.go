package scheduler

import (
	"errors"
	"squzy/apps/internal/config"
	"squzy/apps/internal/job"
	"time"
)

var (
	alreadyRunError             = errors.New("SCHEDULER_ALREADY_RUN")
	alreadyStopError            = errors.New("SCHEDULER_ALREADY_STOP")
	intervalLessHalfSecondError = errors.New("INTERVAL_LESS_THEN_HALF_SECOND")
)

type Scheduler interface {
	// Should run Scheduler every tick
	Run() error
	// Should stop Scheduler
	Stop() error
	// Return true/false depends from current state
	IsRun() bool
}

type schl struct {
	cfg       config.Config
	ticker    *time.Ticker
	isStopped bool
	quitCh    chan bool
	interval  time.Duration
	job       job.Job
}

func New(cfg config.Config, interval time.Duration, job job.Job) (Scheduler, error) {
	if interval < time.Millisecond*500 {
		return nil, intervalLessHalfSecondError
	}
	return &schl{
		cfg:       cfg,
		interval:  interval,
		isStopped: true,
		job:       job,
	}, nil
}

func (s *schl) Run() error {
	if !s.isStopped {
		return alreadyRunError
	}
	s.ticker = time.NewTicker(s.interval)
	s.isStopped = false
	s.quitCh = make(chan bool, 1)
	s.observer()
	return nil
}

func (s *schl) observer() {
	go func() {
	loop:
		for {
			select {
			case <-s.ticker.C:
				_ = s.job.Do()
			case <-s.quitCh:
				break loop
			}
		}
	}()
}

func (s *schl) IsRun() bool {
	return !s.isStopped
}

func (s *schl) Stop() error {
	if s.isStopped {
		return alreadyStopError
	}
	s.ticker.Stop()
	s.quitCh <- true
	close(s.quitCh)
	s.isStopped = true
	return nil
}
