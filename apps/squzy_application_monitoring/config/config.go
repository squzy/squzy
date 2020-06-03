package config

import (
	"os"
	"strconv"
)

const (
	ENV_TRACING_HEADER         = "TRACING_HEADER"
	ENV_PORT                   = "PORT"
	defaultTracingHeader       = "Squzy_transaction"
	defaultPort          int32 = 9095
)

type Config interface {
	GetTracingHeader() string
}

type cfg struct {
	tracingHeader string
	port          int32
}

func (c *cfg) GetTracingHeader() string {
	return c.tracingHeader
}

func New() Config {
	header := os.Getenv(ENV_TRACING_HEADER)
	if header == "" {
		header = defaultTracingHeader
	}

	portValue := os.Getenv(ENV_PORT)
	port := defaultPort
	if portValue != "" {
		i, err := strconv.ParseInt(portValue, 10, 32)
		if err == nil {
			port = int32(i)
		}
	}

	return &cfg{
		tracingHeader: header,
		port:          port,
	}
}
