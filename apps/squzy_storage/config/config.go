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
	ENV_DB_LOGS     = "DB_LOGS"

	ENV_INCIDENT_SERVER_HOST = "INCIDENT_SERVER_HOST"
	ENV_ENABLE_INCIDENT      = "ENABLE_INCIDENT"

	defaultPort int32 = 9090
)

type cfg struct {
	port           int32
	dbHost         string
	dbPort         string
	dbName         string
	dbUser         string
	dbPassword     string
	incidentServer string
	withIncident   bool
	withDbLogs     bool
}

func (c *cfg) GetPort() int32 {
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

func (c *cfg) GetIncidentServerAddress() string {
	return c.incidentServer
}

func (c *cfg) WithIncident() bool {
	return c.withIncident
}

func (c *cfg) WithDbLogs() bool {
	return c.withDbLogs
}

type Config interface {
	GetPort() int32
	GetDbHost() string
	GetDbPort() string
	GetDbName() string
	GetDbUser() string
	GetDbPassword() string
	GetIncidentServerAddress() string
	WithIncident() bool
	WithDbLogs() bool
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
	withIncident := false
	incidentValue := os.Getenv(ENV_ENABLE_INCIDENT)
	if incidentValue != "" {
		value, err := strconv.ParseBool(incidentValue)
		if err == nil {
			withIncident = value
		}
	}

	withDbLog := false
	dbLogsValue := os.Getenv(ENV_DB_LOGS)
	if dbLogsValue != "" {
		value, err := strconv.ParseBool(dbLogsValue)
		if err == nil {
			withDbLog = value
		}
	}
	return &cfg{
		port:           port,
		dbHost:         os.Getenv(ENV_DB_HOST),
		dbPort:         os.Getenv(ENV_DB_PORT),
		dbName:         os.Getenv(ENV_DB_NAME),
		dbUser:         os.Getenv(ENV_DB_USER),
		dbPassword:     os.Getenv(ENV_DB_PASSWORD),
		incidentServer: os.Getenv(ENV_INCIDENT_SERVER_HOST),
		withIncident:   withIncident,
		withDbLogs:     withDbLog,
	}
}
