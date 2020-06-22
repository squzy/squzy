package config

import (
	"os"
	"strconv"
)

const (
	ENV_PORT          = "PORT"
	ENV_STRORAGE_HOST = "ENV_STRORAGE_HOST"

	defaultPort int32 = 9090
)

type cfg struct {
	port        int32
	storageHost string
}

func (c *cfg) GetPort() int32 {
	return c.port
}

func (c *cfg) GetStorageHost() string {
	return c.storageHost
}

type Config interface {
	GetPort() int32
	GetStorageHost() string
}

func New() Config {
	// Read port
	portValue := os.Getenv(ENV_PORT)
	port := defaultPort
	if portValue != "" {
		i, err := strconv.ParseInt(portValue, 10, 32)
		if err == nil {
			port = int32(i)
		}
	}
	return &cfg{
		port:        port,
		storageHost: os.Getenv(ENV_STRORAGE_HOST),
	}
}
