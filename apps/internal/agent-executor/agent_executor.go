package agent_executor

import (
	"errors"
	agentPb "github.com/squzy/squzy_generated/generated/agent/proto/v1"
	"squzy/apps/internal/agent"
	"time"
)

type AgentExecutor interface {
	Execute() chan *agentPb.GetStatsResponse
}

type executor struct {
	agent       agent.Agent
	interval    time.Duration
	statChan    chan *agentPb.GetStatsResponse
	executeChan chan bool
}

const (
	minIntervalExecute = time.Millisecond * 500
)

var (
	intervalLessHalfSecondError = errors.New("INTERVAL_LESS_THAN_HALF_SECOND")
)

func (e *executor) Execute() chan *agentPb.GetStatsResponse {
	go func() {
		for range e.executeChan {
			c := make(chan *agentPb.GetStatsResponse, 1)
			go func() {
				c <- e.agent.GetStat()
			}()
			select {
			case res := <-c:
				close(c)
				e.statChan <- res
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
		return nil, intervalLessHalfSecondError
	}
	return &executor{
		agent:       agent,
		interval:    interval,
		statChan:    make(chan *agentPb.GetStatsResponse, 1),
		executeChan: make(chan bool, 1),
	}, nil
}
