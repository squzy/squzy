package config

import (
	"os"
)

const (
	ENV_PORT = "PORT"

	ENV_DB_HOST     = "DB_HOST"
	ENV_DB_PORT     = "DB_PORT"
	ENV_DB_NAME     = "DB_NAME"
	ENV_DB_USER     = "DB_USER"
	ENV_DB_PASSWORD = "DB_PASSWORD"

	defaultPort = "9090"
)

type cfg struct {
	port       string
	dbHost     string
	dbPort     string
	dbName     string
	dbUser     string
	dbPassword string
}

func (c *cfg) GetPort() string {
	return c.port
}

func (c *cfg) GetDbHost() string {
	return c.dbHost
}

func (c *cfg) GetDbPort() string {
	return c.dbPort
}

func (c *cfg) GetDbName() string {
	return c.dbName
}

func (c *cfg) GetDbUser() string {
	return c.dbUser
}

func (c *cfg) GetDbPassword() string {
	return c.dbPassword
}

type Config interface {
	GetPort() string
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
		port = os.Getenv(ENV_PORT)
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
