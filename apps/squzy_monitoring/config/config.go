package config

import (
	"os"
	"squzy/internal/helpers"
	"strconv"
	"time"
)

type cfg struct {
	port          int32
	timeout       time.Duration
	clientAddress string
}

func (c cfg) GetPort() int32 {
	return c.port
}

func (c cfg) GetClientAddress() string {
	return c.clientAddress
}

func (c cfg) GetStorageTimeout() time.Duration {
	return c.timeout
}

type Config interface {
	GetPort() int32
	GetClientAddress() string
	GetStorageTimeout() time.Duration
}

func New() Config {
	// Read port
	portValue := os.Getenv("PORT")
	port := int32(8080)
	if portValue != "" {
		i, err := strconv.ParseInt(portValue, 10, 32)
		if err == nil {
			port = int32(i)
		}
	}
	// Read storage
	timeoutValue := os.Getenv("STORAGE_TIMEOUT")
	timeoutStorage := time.Second * 5
	if timeoutValue != "" {
		i, err := strconv.ParseInt(timeoutValue, 10, 32)
		if err == nil {
			timeoutStorage = helpers.DurationFromSecond(int32(i))
		}
	}
	// Client Adress
	clientAddress := os.Getenv("STORAGE_HOST")
	return &cfg{
		clientAddress: clientAddress,
		timeout:       timeoutStorage,
		port:          port,
	}
}
