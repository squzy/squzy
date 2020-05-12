package agent_executor

import (
	"errors"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"squzy/internal/agent"
	"time"
)

type AgentExecutor interface {
	Execute() chan *apiPb.Metric
}

type executor struct {
	agent       agent.Agent
	interval    time.Duration
	statChan    chan *apiPb.Metric
	executeChan chan bool
}

const (
	minIntervalExecute = time.Millisecond * 500
)

var (
	errIntervalLessHalfSecondError = errors.New("INTERVAL_LESS_THAN_HALF_SECOND")
)

func (e *executor) Execute() chan *apiPb.Metric {
	go func() {
		for range e.executeChan {
			c := make(chan *apiPb.Metric, 1)
			go func() {
				c <- e.agent.GetStat()
			}()
			select {
			case res := <-c:
				close(c)
				e.statChan <- res
				time.Sleep(e.interval)
				e.executeChan <- true
			case <-time.After(e.interval):
				close(c)
				e.executeChan <- true
			}
		}
	}()
	e.executeChan <- true
	return e.statChan
}

func New(agent agent.Agent, interval time.Duration) (AgentExecutor, error) {
	if interval < minIntervalExecute {
		return nil, errIntervalLessHalfSecondError
	}
	return &executor{
		agent:       agent,
		interval:    interval,
		statChan:    make(chan *apiPb.Metric, 1),
		executeChan: make(chan bool, 1),
	}, nil
}
