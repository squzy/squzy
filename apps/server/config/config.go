package config

import (
	"os"
	"strconv"
)

type Cfg interface {
	GetPort() int32
	GetDbHost() string
	GetDbPort() string
	GetDbUser() string
	GetDbPassword() string
	GetDbName() string
}

type config struct {
	port       int32
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
}

func (c *config) GetPort() int32 {
	return c.port
}

func (c *config) GetDbHost() string {
	return c.dbHost
}

func (c *config) GetDbPort() string {
	return c.dbPort
}

func (c *config) GetDbUser() string {
	return c.dbUser
}

func (c *config) GetDbPassword() string {
	return c.dbPassword
}

func (c *config) GetDbName() string {
	return c.dbName
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
		port:       port,
		dbHost:     getStringValue("DB_HOST", "localhost"),
		dbPort:     getStringValue("DB_PORT", "5432"),
		dbUser:     getStringValue("DB_USER", "user"),
		dbPassword: getStringValue("DB_PASSWORD", "password"),
		dbName:     getStringValue("DB_NAME", "database"),
	}
}

func getStringValue(name string, defaultValue string) string {
	res := os.Getenv(name)
	if res != "" {
		return res
	}
	return defaultValue
}
