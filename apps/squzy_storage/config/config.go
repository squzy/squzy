package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
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
	Port           int32  `yaml:"port"`
	DbHost         string `yaml:"dbHost"`
	DbPort         string `yaml:"dbPort"`
	DbName         string `yaml:"dbName"`
	DbUser         string `yaml:"dbUser"`
	DbPassword     string `yaml:"dbPassword"`
	IncidentServer string `yaml:"incidentServer"`
	WithIncident   bool   `yaml:"withIncident"`
	WithDbLogs     bool   `yaml:"withDbLogs"`
}

func (c *cfg) GetPort() int32 {
	return c.Port
}

func (c *cfg) GetDbHost() string {
	return c.DbHost
}

func (c *cfg) GetDbPort() string {
	return c.DbPort
}

func (c *cfg) GetDbName() string {
	return c.DbName
}

func (c *cfg) GetDbUser() string {
	return c.DbUser
}

func (c *cfg) GetDbPassword() string {
	return c.DbPassword
}

func (c *cfg) GetIncidentServerAddress() string {
	return c.IncidentServer
}

func (c *cfg) GetWithIncident() bool {
	return c.WithIncident
}

func (c *cfg) GetWithDbLogs() bool {
	return c.WithDbLogs
}

type Config interface {
	GetPort() int32
	GetDbHost() string
	GetDbPort() string
	GetDbName() string
	GetDbUser() string
	GetDbPassword() string
	GetIncidentServerAddress() string
	GetWithIncident() bool
	GetWithDbLogs() bool
}

func New() Config {
	// Read Port
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
		Port:           port,
		DbHost:         os.Getenv(ENV_DB_HOST),
		DbPort:         os.Getenv(ENV_DB_PORT),
		DbName:         os.Getenv(ENV_DB_NAME),
		DbUser:         os.Getenv(ENV_DB_USER),
		DbPassword:     os.Getenv(ENV_DB_PASSWORD),
		IncidentServer: os.Getenv(ENV_INCIDENT_SERVER_HOST),
		WithIncident:   withIncident,
		WithDbLogs:     withDbLog,
	}
}

// Read config from .yaml file
// reader & unmarshall are provided in order to easy-mock
func NewConfigFromYaml(cfgByte []byte) (Config, error) {
	var config cfg
	err := yaml.Unmarshal(cfgByte, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling cfg file: %w", err)
	}

	// validate config
	if config.DbHost == "" {
		return nil, errors.New("empty config dbHost")
	}
	if config.DbPort == "" {
		return nil, errors.New("empty config dbPort")
	}
	if config.DbName == "" {
		return nil, errors.New("empty config dbName")
	}
	if config.DbUser == "" {
		return nil, errors.New("empty config dbUser")
	}
	if config.DbPassword == "" {
		return nil, errors.New("empty config dbPassword")
	}
	if config.WithIncident && config.IncidentServer == "" {
		return nil, errors.New("empty config incidentServer when withIncident true")
	}
	if config.Port == 0 {
		config.Port = defaultPort
	}
	return &config, nil
}
