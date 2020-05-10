package config

import (
	"os"
	"strconv"
)

type Config interface {
	GetAgentServerAddress() string
	GetMonitoringServerAddress() string
	GetPort() int32
}

type cfg struct {
	port             int32
	agentServer      string
	monitoringServer string
}

func (c *cfg) GetAgentServerAddress() string {
	return c.agentServer
}

func (c *cfg) GetMonitoringServerAddress() string {
	return c.monitoringServer
}

func (c *cfg) GetPort() int32 {
	return c.port
}

const (
	ENV_PORT              = "PORT"
	ENV_AGENT_SERVER      = "AGENT_SERVER_HOST"
	ENV_MONITORING_SERVER = "MONITORING_SERVER_HOST"

	defaultPort    = int32(8080)
)

func New() Config {
	portValue := os.Getenv(ENV_PORT)
	port := defaultPort
	if portValue != "" {
		i, err := strconv.ParseInt(portValue, 10, 32)
		if err == nil {
			port = int32(i)
		}
	}
	return &cfg{
		port:             port,
		agentServer:      os.Getenv(ENV_AGENT_SERVER),
		monitoringServer: os.Getenv(ENV_MONITORING_SERVER),
	}
}
