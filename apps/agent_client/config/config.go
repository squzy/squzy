package config

import (
	"github.com/squzy/squzy/internal/helpers"
	"os"
	"strconv"
	"time"
)

const (
	ENV_SQUZY_AGENT_INTERVAL    = "SQUZY_AGENT_INTERVAL"
	ENV_SQUZY_AGENT_SERVER_HOST = "SQUZY_AGENT_SERVER_HOST"
	ENV_SQUZY_AGENT_NAME        = "SQUZY_AGENT_NAME"
	ENV_SQUZY_AGENT_RETRY       = "SQUZY_AGENT_RETRY"
	ENV_SQUZY_AGENT_RETRY_COUNT = "SQUZY_AGENT_RETRY_COUNT"

	defaultTimeout = time.Second * 5
)

type Config interface {
	GetAgentServer() string
	GetInterval() time.Duration
	GetAgentName() string
	Retry() bool
	RetryCount() int32
}

type cfg struct {
	host             string
	executionTimeout time.Duration
	agentName        string
	retry            bool
	retryCount       int32
}

func (c *cfg) Retry() bool {
	return c.retry
}

func (c *cfg) RetryCount() int32 {
	return c.retryCount
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

	retry := false
	retryCount := int32(0)
	retryValue := os.Getenv(ENV_SQUZY_AGENT_RETRY)
	retryCountValue := os.Getenv(ENV_SQUZY_AGENT_RETRY_COUNT)

	if retryValue != "" {
		needRetry, err := strconv.ParseBool(retryValue)
		if err == nil {
			retry = needRetry
		}
	}

	if retryCountValue != "" {
		countValue, err := strconv.ParseInt(retryCountValue, 10, 32)
		if err == nil {
			retryCount = int32(countValue)
		}
	}

	return &cfg{
		host:             os.Getenv(ENV_SQUZY_AGENT_SERVER_HOST),
		executionTimeout: timeoutExecution,
		agentName:        os.Getenv(ENV_SQUZY_AGENT_NAME),
		retry:            retry,
		retryCount:       retryCount,
	}
}
