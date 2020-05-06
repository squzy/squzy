package config

import (
	"os"
	"squzy/internal/helpers"
	"strconv"
	"time"
)

const (
	ENV_SQUZY_AGENT_INTERVAL    = "SQUZY_AGENT_INTERVAL"
	ENV_SQUZY_AGENT_SERVER_HOST = "SQUZY_AGENT_SERVER_HOST"
	ENV_SQUZY_AGENT_NAME        = "SQUZY_AGENT_NAME"

	defaultTimeout = time.Second * 5
)

type Config interface {
	GetAgentServer() string
	GetInterval() time.Duration
	GetAgentName() string
}

type cfg struct {
	host             string
	executionTimeout time.Duration
	agentName        string
}

func (c *cfg) GetInterval() time.Duration {
	return c.executionTimeout
}

func (c *cfg) GetAgentServer() string {
	return c.host
}

func (c *cfg) GetAgentName() string {
	return c.agentName
}

func New() Config {
	// Server timeout execution
	timeoutExecutionValue := os.Getenv(ENV_SQUZY_AGENT_INTERVAL)
	timeoutExecution := defaultTimeout
	if timeoutExecutionValue != "" {
		i, err := strconv.ParseInt(timeoutExecutionValue, 10, 32)
		if err == nil {
			timeoutExecution = helpers.DurationFromSecond(int32(i))
		}
	}

	return &cfg{
		host:             os.Getenv(ENV_SQUZY_AGENT_SERVER_HOST),
		executionTimeout: timeoutExecution,
		agentName:        os.Getenv(ENV_SQUZY_AGENT_NAME),
	}
}
