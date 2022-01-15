package config

import (
	"os"
	"github.com/squzy/squzy/internal/helpers"
	"strconv"
	"time"
)

const (
	ENV_PORT             = "PORT"
	ENV_STORAGE_TIMEOUT  = "SQUZY_STORAGE_TIMEOUT"
	ENV_MONGO_DB         = "MONGO_DB"
	ENV_MONGO_URI        = "MONGO_URI"
	ENV_MONGO_COLLECTION = "MONGO_COLLECTION"
	ENV_STORAGE_HOST     = "SQUZY_STORAGE_HOST"

	defaultPort           int32 = 9091
	defaultStorageTimeout       = time.Second * 5
	defaultMongoDb              = "squzy_agent"
	defaultCollection           = "agents"
)

type cfg struct {
	port            int32
	timeout         time.Duration
	storageAddress  string
	mongoURI        string
	mongoDb         string
	mongoCollection string
}

func (c *cfg) GetPort() int32 {
	return c.port
}

func (c *cfg) GetStorageAddress() string {
	return c.storageAddress
}

func (c *cfg) GetStorageTimeout() time.Duration {
	return c.timeout
}

func (c *cfg) GetMongoURI() string {
	return c.mongoURI
}

func (c *cfg) GetMongoDb() string {
	return c.mongoDb
}

func (c *cfg) GetMongoCollection() string {
	return c.mongoCollection
}

type Config interface {
	GetPort() int32
	GetStorageAddress() string
	GetStorageTimeout() time.Duration
	GetMongoURI() string
	GetMongoDb() string
	GetMongoCollection() string
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
	// Read storage
	timeoutValue := os.Getenv(ENV_STORAGE_TIMEOUT)
	timeoutStorage := defaultStorageTimeout
	if timeoutValue != "" {
		i, err := strconv.ParseInt(timeoutValue, 10, 32)
		if err == nil {
			timeoutStorage = helpers.DurationFromSecond(int32(i))
		}
	}
	mongoDb := os.Getenv(ENV_MONGO_DB)
	if mongoDb == "" {
		mongoDb = defaultMongoDb
	}
	collection := os.Getenv(ENV_MONGO_COLLECTION)
	if collection == "" {
		collection = defaultCollection
	}
	return &cfg{
		storageAddress:  os.Getenv(ENV_STORAGE_HOST),
		timeout:         timeoutStorage,
		port:            port,
		mongoURI:        os.Getenv(ENV_MONGO_URI),
		mongoDb:         mongoDb,
		mongoCollection: collection,
	}
}
