package config

import (
	"os"
	"squzy/internal/helpers"
	"strconv"
	"time"
)

const (
	ENV_TRACING_HEADER   = "TRACING_HEADER"
	ENV_PORT             = "PORT"
	ENV_MONGO_DB         = "MONGO_DB"
	ENV_MONGO_URI        = "MONGO_URI"
	ENV_MONGO_COLLECTION = "MONGO_COLLECTION"
	ENV_STORAGE_HOST     = "SQUZY_STORAGE_HOST"
	ENV_STORAGE_TIMEOUT  = "SQUZY_STORAGE_TIMEOUT"

	defaultTracingHeader        = "Squzy_transaction"
	defaultPort           int32 = 9095
	defaultMongoDb              = "applications_monitoring"
	defaultCollection           = "application"
	defaultStorageTimeout       = time.Second * 5
)

type Config interface {
	GetTracingHeader() string
	GetPort() int32
	GetMongoURI() string
	GetMongoDb() string
	GetMongoCollection() string
	GetStorageTimeout() time.Duration
	GetStorageHost() string
}

type cfg struct {
	tracingHeader   string
	port            int32
	timeout         time.Duration
	storageAddress  string
	mongoURI        string
	mongoDb         string
	mongoCollection string
}

func (c *cfg) GetStorageHost() string {
	return c.storageAddress
}

func (c *cfg) GetPort() int32 {
	return c.port
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

func (c *cfg) GetStorageTimeout() time.Duration {
	return c.timeout
}

func (c *cfg) GetTracingHeader() string {
	return c.tracingHeader
}

func New() Config {
	header := os.Getenv(ENV_TRACING_HEADER)
	if header == "" {
		header = defaultTracingHeader
	}

	portValue := os.Getenv(ENV_PORT)
	port := defaultPort
	if portValue != "" {
		i, err := strconv.ParseInt(portValue, 10, 32)
		if err == nil {
			port = int32(i)
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

	timeoutValue := os.Getenv(ENV_STORAGE_TIMEOUT)
	timeoutStorage := defaultStorageTimeout
	if timeoutValue != "" {
		i, err := strconv.ParseInt(timeoutValue, 10, 32)
		if err == nil {
			timeoutStorage = helpers.DurationFromSecond(int32(i))
		}
	}

	return &cfg{
		tracingHeader:   header,
		port:            port,
		mongoURI:        os.Getenv(ENV_MONGO_URI),
		mongoDb:         mongoDb,
		mongoCollection: collection,
		timeout:         timeoutStorage,
		storageAddress: os.Getenv(ENV_STORAGE_HOST),
	}
}
