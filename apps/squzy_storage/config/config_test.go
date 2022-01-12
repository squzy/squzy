package config

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
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
		assert.Equal(t, s.WithIncident(), false)
		assert.Equal(t, s.WithDbLogs(), false)
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
		assert.Equal(t, s.WithIncident(), true)
	})
}

func TestCfg_WithDbLogs(t *testing.T) {
	t.Run("Should: return from env", func(t *testing.T) {
		err := os.Setenv(ENV_DB_LOGS, "true")
		if err != nil {
			assert.NotNil(t, nil)
		}
		s := New()
		assert.Equal(t, s.WithDbLogs(), true)
	})
}

func TestReadFromMap(t *testing.T) {

	testCases := []struct {
		name          string
		input         map[string]interface{}
		expectedError error
	}{
		{
			name:          "Empty db host",
			input:         map[string]interface{}{},
			expectedError: errors.New(fmt.Sprintf("wrong %s", ENV_DB_HOST)),
		},
		{
			name: "Wrong db host",
			input: map[string]interface{}{
				strings.ToLower(ENV_DB_HOST): struct{}{},
			},
			expectedError: errors.New(fmt.Sprintf("wrong %s", ENV_DB_HOST)),
		},
		{
			name: "Empty db port",
			input: map[string]interface{}{
				strings.ToLower(ENV_DB_HOST): "",
			},
			expectedError: errors.New(fmt.Sprintf("wrong %s", ENV_DB_PORT)),
		},
		{
			name: "Empty db name",
			input: map[string]interface{}{
				strings.ToLower(ENV_DB_HOST): "",
				strings.ToLower(ENV_DB_PORT): "",
			},
			expectedError: errors.New(fmt.Sprintf("wrong %s", ENV_DB_NAME)),
		},
		{
			name: "Empty db user",
			input: map[string]interface{}{
				strings.ToLower(ENV_DB_HOST): "",
				strings.ToLower(ENV_DB_PORT): "",
				strings.ToLower(ENV_DB_NAME): "",
			},
			expectedError: errors.New(fmt.Sprintf("wrong %s", ENV_DB_USER)),
		},
		{
			name: "Empty db password",
			input: map[string]interface{}{
				strings.ToLower(ENV_DB_HOST): "",
				strings.ToLower(ENV_DB_PORT): "",
				strings.ToLower(ENV_DB_NAME): "",
				strings.ToLower(ENV_DB_USER): "",
			},
			expectedError: errors.New(fmt.Sprintf("wrong %s", ENV_DB_PASSWORD)),
		},
		{
			name: "Empty port",
			input: map[string]interface{}{
				strings.ToLower(ENV_DB_HOST):     "",
				strings.ToLower(ENV_DB_PORT):     "",
				strings.ToLower(ENV_DB_NAME):     "",
				strings.ToLower(ENV_DB_USER):     "",
				strings.ToLower(ENV_DB_PASSWORD): "",
			},
			expectedError: nil,
		},
		{
			name: "Port provided",
			input: map[string]interface{}{
				strings.ToLower(ENV_DB_HOST):     "",
				strings.ToLower(ENV_DB_PORT):     "",
				strings.ToLower(ENV_DB_NAME):     "",
				strings.ToLower(ENV_DB_USER):     "",
				strings.ToLower(ENV_DB_PASSWORD): "",
				strings.ToLower(ENV_PORT):        2020,
			},
			expectedError: nil,
		},
		{
			name: "Port has incorrect format",
			input: map[string]interface{}{
				strings.ToLower(ENV_DB_HOST):     "",
				strings.ToLower(ENV_DB_PORT):     "",
				strings.ToLower(ENV_DB_NAME):     "",
				strings.ToLower(ENV_DB_USER):     "",
				strings.ToLower(ENV_DB_PASSWORD): "",
				strings.ToLower(ENV_PORT):        struct{}{},
			},
			expectedError: nil,
		},
		{
			name: "With incident provided & empty incident host",
			input: map[string]interface{}{
				strings.ToLower(ENV_DB_HOST):         "",
				strings.ToLower(ENV_DB_PORT):         "",
				strings.ToLower(ENV_DB_NAME):         "",
				strings.ToLower(ENV_DB_USER):         "",
				strings.ToLower(ENV_DB_PASSWORD):     "",
				strings.ToLower(ENV_PORT):            2020,
				strings.ToLower(ENV_ENABLE_INCIDENT): true,
			},
			expectedError: errors.New(fmt.Sprintf("wrong %s", ENV_INCIDENT_SERVER_HOST)),
		},
		{
			name: "With incident provided & empty incident host",
			input: map[string]interface{}{
				strings.ToLower(ENV_DB_HOST):              "",
				strings.ToLower(ENV_DB_PORT):              "",
				strings.ToLower(ENV_DB_NAME):              "",
				strings.ToLower(ENV_DB_USER):              "",
				strings.ToLower(ENV_DB_PASSWORD):          "",
				strings.ToLower(ENV_PORT):                 2020,
				strings.ToLower(ENV_ENABLE_INCIDENT):      true,
				strings.ToLower(ENV_INCIDENT_SERVER_HOST): "",
			},
			expectedError: nil,
		},
		{
			name: "With db logs provided",
			input: map[string]interface{}{
				strings.ToLower(ENV_DB_HOST):              "",
				strings.ToLower(ENV_DB_PORT):              "",
				strings.ToLower(ENV_DB_NAME):              "",
				strings.ToLower(ENV_DB_USER):              "",
				strings.ToLower(ENV_DB_PASSWORD):          "",
				strings.ToLower(ENV_PORT):                 2020,
				strings.ToLower(ENV_ENABLE_INCIDENT):      true,
				strings.ToLower(ENV_INCIDENT_SERVER_HOST): "",
				strings.ToLower(ENV_DB_LOGS):              true,
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(ENV_DB_HOST, "")
			os.Setenv(ENV_DB_PORT, "")
			os.Setenv(ENV_DB_NAME, "")
			os.Setenv(ENV_DB_USER, "")
			os.Setenv(ENV_DB_PASSWORD, "")
			os.Setenv(ENV_PORT, "")
			os.Setenv(ENV_ENABLE_INCIDENT, "")
			os.Setenv(ENV_INCIDENT_SERVER_HOST, "")
			os.Setenv(ENV_DB_LOGS, "")
			_, err := ReadFromMap(tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}

	t.Run("Test when all env read from os", func(t *testing.T) {
		os.Setenv(ENV_DB_HOST, "ENV_DB_HOST")
		os.Setenv(ENV_DB_PORT, "ENV_DB_PORT")
		os.Setenv(ENV_DB_NAME, "ENV_DB_NAME")
		os.Setenv(ENV_DB_USER, "ENV_DB_USER")
		os.Setenv(ENV_DB_PASSWORD, "ENV_DB_PASSWORD")
		os.Setenv(ENV_PORT, "2020")
		os.Setenv(ENV_ENABLE_INCIDENT, "true")
		os.Setenv(ENV_INCIDENT_SERVER_HOST, "ENV_INCIDENT_SERVER_HOST")
		os.Setenv(ENV_DB_LOGS, "true")
		_, err := ReadFromMap(map[string]interface{}{})
		assert.Equal(t, nil, err)
	})
}
