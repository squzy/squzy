package config

import (
	"os"
	"strconv"
)

const (
	ENV_PORT = "PORT"

	ENV_DB_HOST     = "DB_HOST"
	ENV_DB_PORT     = "DB_PORT"
	ENV_DB_NAME     = "DB_NAME"
	ENV_DB_USER     = "DB_USER"
	ENV_DB_PASSWORD = "DB_PASSWORD"

	defaultPort int32 = 9090
)

type cfg struct {
	port       int32
	dbHost     string
	dbPort     string
	dbName     string
	dbUser     string
	dbPassword string
}

func (c *cfg) GetPort() int32 {
	return c.port
}

func (c *cfg) GetDbHost() string {
	return "localhost"
	return c.dbHost
}

func (c *cfg) GetDbPort() string {
	return "5432"
	return c.dbPort
}

func (c *cfg) GetDbName() string {
	return "database"
	return c.dbName
}

func (c *cfg) GetDbUser() string {
	return "user"
	return c.dbUser
}

func (c *cfg) GetDbPassword() string {
	return "password"
	return c.dbPassword
}

type Config interface {
	GetPort() int32
	GetDbHost() string
	GetDbPort() string
	GetDbName() string
	GetDbUser() string
	GetDbPassword() string
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
		port:       port,
		dbHost:     os.Getenv(ENV_DB_HOST),
		dbPort:     os.Getenv(ENV_DB_PORT),
		dbName:     os.Getenv(ENV_DB_NAME),
		dbUser:     os.Getenv(ENV_DB_USER),
		dbPassword: os.Getenv(ENV_DB_PASSWORD),
	}
}
