package scheduler

import (
	"errors"
	"fmt"
	"github.com/squzy/squzy/apps/squzy_monitoring/config"
	"github.com/squzy/squzy/internal/cache"
	job_executor "github.com/squzy/squzy/internal/job-executor"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	Run() error
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
	cache       cache.Cache
}

func New(id primitive.ObjectID, interval time.Duration, jobExecutor job_executor.JobExecutor, cache cache.Cache) (Scheduler, error) {
	if interval < config.SmallestInterval {
		return nil, errIntervalLessHalfSecondError
	}
	return &schl{
		id:          id,
		interval:    interval,
		isStopped:   true,
		jobExecutor: jobExecutor,
		cache:       cache,
	}, nil
}

func (s *schl) Run() error {
	if !s.isStopped {
		return nil
	}
	err := s.cache.InsertSchedule(&apiPb.InsertScheduleWithIdRequest{
		Id:            s.id.Hex(),
		ScheduledNext: timestamppb.New(time.Now().Add(s.interval)),
	})
	if err != nil {
		return fmt.Errorf("could not insert during run: %w", err)
	}
	s.ticker = time.NewTicker(config.SmallestInterval / 2)
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
				res, err := s.cache.GetScheduleById(&apiPb.GetScheduleWithIdRequest{
					Id: s.id.Hex(),
				})
				if err != nil {
					logger.Error("could not get schedule" + err.Error())
				}

				if res == nil {
					return
				}
				next := res.ScheduledNext.AsTime()
				now := time.Now()

				if now.After(next) {
					s.jobExecutor.Execute(s.id)

					err := s.cache.InsertSchedule(&apiPb.InsertScheduleWithIdRequest{
						Id:            s.id.Hex(),
						ScheduledNext: timestamppb.New(now.Add(s.interval)),
					})
					if err != nil {
						logger.Error("could not update schedule" + err.Error())
					}
				}
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
	err := s.cache.DeleteScheduleById(&apiPb.DeleteScheduleWithIdRequest{
		Id: s.id.Hex(),
	})
	if err != nil {
		logger.Info("could not delete schedule by id: " + err.Error())
	}
	s.ticker.Stop()
	s.quitCh <- true
	close(s.quitCh)
	s.isStopped = true
}
