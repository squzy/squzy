package config

import (
	"os"
)

const (
	ENV_SQUZY_LOG_LEVEL = "SQUZY_AGENT_INTERVAL"
)

type Config interface {
	GetLogLevel() string
}

type cfg struct {
	logLevel string
}

func (c *cfg) GetLogLevel() string {
	return c.logLevel
}

func New() Config {
	return &cfg{
		logLevel: os.Getenv(ENV_SQUZY_LOG_LEVEL),
	}
}
