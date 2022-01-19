package config

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Shoud: return default value", func(t *testing.T) {
		s := New()
		assert.Equal(t, s.GetPort(), defaultPort)
		assert.Equal(t, s.GetDbHost(), "")
		assert.Equal(t, s.GetDbPort(), "")
		assert.Equal(t, s.GetDbName(), "")
		assert.Equal(t, s.GetDbUser(), "")
		assert.Equal(t, s.GetDbPassword(), "")
		assert.Equal(t, s.GetIncidentServerAddress(), "")
		assert.Equal(t, s.GetWithIncident(), false)
		assert.Equal(t, s.GetWithDbLogs(), false)
	})
}

func TestCfg_GetClientAddress(t *testing.T) {

}

func TestCfg_GetPort(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_PORT, "11124")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetPort(), int32(11124))
	})
}

func TestCfg_GetDbHost(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_DB_HOST, "dbhost")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetDbHost(), "dbhost")
	})
}

func TestCfg_GetDbPort(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_DB_PORT, "dbport")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetDbPort(), "dbport")
	})
}

func TestCfg_GetDbName(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_DB_NAME, "dbname")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetDbName(), "dbname")
	})
}

func TestCfg_GetDbUser(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_DB_USER, "dbuser")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetDbUser(), "dbuser")
	})
}

func TestCfg_GetDbPassword(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_DB_PASSWORD, "dbpassword")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetDbPassword(), "dbpassword")
	})
}

func TestCfg_GetIncidentServerAddress(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_INCIDENT_SERVER_HOST, "dbpassword")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetIncidentServerAddress(), "dbpassword")
	})
}

func TestCfg_WithIncident(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_ENABLE_INCIDENT, "true")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetWithIncident(), true)
	})
}

func TestCfg_WithDbLogs(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_DB_LOGS, "true")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetWithDbLogs(), true)
	})
}

func TestCfg_DbType(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_DB_TYPE, DB_TYPE_POSTGRES)
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.GetWithDbLogs(), true)
	})
}

func TestNewConfigFromYaml(t *testing.T) {
	marshall := func(config *cfg) []byte {
		res, err := yaml.Marshal(config)
		if err != nil {
			return []byte("error")
		}
		return res
	}
	testCases := []struct {
		name          string
		cfgByte       []byte
		expectedError error
	}{
		{
			name:          "error unmarshalling file",
			cfgByte:       []byte("error"),
			expectedError: fmt.Errorf("error unmarshalling cfg file: %w", errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `error` into config.cfg")),
		},
		{
			name: "error empty dbHost",
			cfgByte: marshall(&cfg{
			}),
			expectedError: errors.New("empty config dbHost"),
		},
		{
			name: "error empty dbPort",
			cfgByte: marshall(&cfg{
				DbHost: "dbHost",
			}),
			expectedError: errors.New("empty config dbPort"),
		},
		{
			name: "error empty dbName",
			cfgByte: marshall(&cfg{
				DbHost: "dbHost",
				DbPort: "dbPort",
			}),
			expectedError: errors.New("empty config dbName"),
		},
		{
			name: "error empty dbUser",
			cfgByte: marshall(&cfg{
				DbHost: "dbHost",
				DbPort: "dbPort",
				DbName: "dbName",
			}),
			expectedError: errors.New("empty config dbUser"),
		},
		{
			name: "error empty dbPassword",
			cfgByte: marshall(&cfg{
				DbHost: "dbHost",
				DbPort: "dbPort",
				DbName: "dbName",
				DbUser: "dbUser",
			}),
			expectedError: errors.New("empty config dbPassword"),
		},
		{
			name: "error empty incidentServer when withIncident = true",
			cfgByte: marshall(&cfg{
				DbHost:       "dbHost",
				DbPort:       "dbPort",
				DbName:       "dbName",
				DbUser:       "dbUser",
				DbPassword:   "dbPassword",
				WithIncident: true,
			}),
			expectedError: errors.New("empty config incidentServer when withIncident true"),
		},
		{
			name: "no error, default port and database",
			cfgByte: marshall(&cfg{
				DbHost:         "dbHost",
				DbPort:         "dbPort",
				DbName:         "dbName",
				DbUser:         "dbUser",
				DbPassword:     "dbPassword",
				IncidentServer: "incidentServer",
				WithIncident:   true,
				WithDbLogs:     true,
			}),
		},
		{
			name: "no error",
			cfgByte: marshall(&cfg{
				Port:           3030,
				DbHost:         "dbHost",
				DbPort:         "dbPort",
				DbName:         "dbName",
				DbUser:         "dbUser",
				DbPassword:     "dbPassword",
				DbType:         DB_TYPE_CLICKHOUSE,
				IncidentServer: "incidentServer",
				WithIncident:   true,
				WithDbLogs:     true,
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewConfigFromYaml(tc.cfgByte)
			if err != nil && tc.expectedError != nil {
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, tc.expectedError)
				assert.Nil(t, err)
			}
		})
	}
}
