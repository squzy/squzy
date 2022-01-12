package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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

// ReadFromMap is used for receiving config from map[string]interface{}
// The last one provided in cli by viper package
func ReadFromMap(cfgMap map[string]interface{}) (Config, error) {
	config := &cfg{}
	// Checking that required fields are present and have correct format
	dbHost, err := readStringValueFromMap(cfgMap, ENV_DB_HOST)
	if err != nil {
		return nil, err
	}
	config.dbHost = dbHost

	dbPort, err := readStringValueFromMap(cfgMap, ENV_DB_PORT)
	if err != nil {
		// trying to read as int
		dbPortInt, err := readIntValueFromMap(cfgMap, ENV_DB_PORT)
		if err != nil {
			return nil, err
		}
		config.dbPort = strconv.Itoa(dbPortInt)
	} else {
		config.dbPort = dbPort
	}

	dbName, err := readStringValueFromMap(cfgMap, ENV_DB_NAME)
	if err != nil {
		return nil, err
	}
	config.dbName = dbName

	dbUser, err := readStringValueFromMap(cfgMap, ENV_DB_USER)
	if err != nil {
		return nil, err
	}
	config.dbUser = dbUser

	dbPassword, err := readStringValueFromMap(cfgMap, ENV_DB_PASSWORD)
	if err != nil {
		return nil, err
	}
	config.dbUser = dbPassword

	// Not required fields
	port, err := readIntValueFromMap(cfgMap, ENV_PORT)
	if err == nil {
		config.port = int32(port)
	} else {
		config.port = defaultPort
	}

	withIncident, err := readBoolValueFromMap(cfgMap, ENV_ENABLE_INCIDENT)
	if err == nil {
		config.withIncident = withIncident
	}
	if withIncident {
		incidentServer, err := readStringValueFromMap(cfgMap, ENV_INCIDENT_SERVER_HOST)
		if err != nil {
			return nil, err
		}
		config.incidentServer = incidentServer
	}
	withDbLogs, err := readBoolValueFromMap(cfgMap, ENV_DB_LOGS)
	if err == nil {
		config.withDbLogs = withDbLogs
	}
	return config, nil
}

// readStringValueFromMap reads the field with given key
// trying to read from env if field is not presented or has wrong format
// We need strings.ToLower(key), because viper lowercases all the keys
func readStringValueFromMap(cfgMap map[string]interface{}, key string) (string, error) {
	var err error
	if _, ok := cfgMap[strings.ToLower(key)]; !ok {
		err = errors.New(fmt.Sprintf("empty %s", key))
	}
	if value, ok := cfgMap[strings.ToLower(key)].(string); !ok {
		err = errors.New(fmt.Sprintf("wrong %s", key))
	} else {
		return value, nil
	}
	// trying to read from env if field is not presented or has wrong format
	if os.Getenv(key) == "" {
		return "", err
	}
	return os.Getenv(key), nil
}

// readStringValueFromMap reads the field with given key
// trying to read from env if field is not presented or has wrong format
// We need strings.ToLower(key), because viper lowercases all the keys
func readIntValueFromMap(cfgMap map[string]interface{}, key string) (int, error) {
	var err error
	if _, ok := cfgMap[strings.ToLower(key)]; !ok {
		err = errors.New(fmt.Sprintf("empty %s", key))
	}
	if value, ok := cfgMap[strings.ToLower(key)].(int); !ok {
		err = errors.New(fmt.Sprintf("wrong %s", key))
	} else {
		return value, nil
	}
	// trying to read from env if field is not presented or has wrong format
	if os.Getenv(key) == "" {
		return 0, err
	}
	return strconv.Atoi(os.Getenv(key))
}

// readStringValueFromMap reads the field with given key
// trying to read from env if field is not presented or has wrong format
// We need strings.ToLower(key), because viper lowercases all the keys
func readBoolValueFromMap(cfgMap map[string]interface{}, key string) (bool, error) {
	var err error
	if _, ok := cfgMap[strings.ToLower(key)]; !ok {
		err = errors.New(fmt.Sprintf("empty %s", key))
	}
	if value, ok := cfgMap[strings.ToLower(key)].(bool); !ok {
		err = errors.New(fmt.Sprintf("wrong %s", key))
	} else {
		return value, nil
	}
	// trying to read from env if field is not presented or has wrong format
	if os.Getenv(key) == "" {
		return false, err
	}
	return strconv.ParseBool(os.Getenv(key))
}
