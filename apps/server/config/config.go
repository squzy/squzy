package config

import (
	"os"
	"strconv"
)

type Cfg interface {
	GetPort() int32
}

type config struct {
	port int32
}

func (c *config) GetPort() int32 {
	return c.port
}

func New() Cfg {
	// Read port
	portValue := os.Getenv("PORT")
	port := int32(8080)
	if portValue != "" {
		i, err := strconv.ParseInt(portValue, 10, 32)
		if err == nil {
			port = int32(i)
		}
	}
	return &config{
		port: port,
	}
}