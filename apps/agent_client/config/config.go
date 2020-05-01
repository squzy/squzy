package config

import (
	"os"
	"strconv"
	"time"
)

type Config interface {
	GetSquzyServer() string
	GetExecutionTimeout() time.Duration
	GetSquzyServerTimeout() time.Duration
	GetAgentName() string
}

type cfg struct {
	timeout          time.Duration
	host             string
	executionTimeout time.Duration
	agentName        string
}

func (c *cfg) GetExecutionTimeout() time.Duration {
	return c.executionTimeout
}

func (c *cfg) GetSquzyServer() string {
	return c.host
}

func (c *cfg) GetSquzyServerTimeout() time.Duration {
	return c.timeout
}

func (c *cfg) GetAgentName() string {
	return c.agentName
}

const defaultTimeout = time.Second * 5

func New() Config {
	// Server timeout connection
	timeoutValue := os.Getenv("SQUZY_SERVER_TIMEOUT")
	timeoutServer := defaultTimeout
	if timeoutValue != "" {
		i, err := strconv.ParseInt(timeoutValue, 10, 32)
		if err == nil {
			timeoutServer = time.Duration(int32(i)) * time.Second
		}
	}
	// Server timeout execution
	timeoutExecutionValue := os.Getenv("SQUZY_EXECUTION_TIMEOUT")
	timeoutExecution := defaultTimeout
	if timeoutExecutionValue != "" {
		i, err := strconv.ParseInt(timeoutExecutionValue, 10, 32)
		if err == nil {
			timeoutExecution = time.Duration(int32(i)) * time.Second
		}
	}
	// Squzy server address
	clientAddress := os.Getenv("SQUZY_SERVER_HOST")
	agentName := os.Getenv("SQUZY_AGENT_NAME")
	return &cfg{
		timeout:          timeoutServer,
		host:             clientAddress,
		executionTimeout: timeoutExecution,
		agentName:        agentName,
	}
}
