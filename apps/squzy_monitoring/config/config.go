package config

import (
	"github.com/squzy/squzy/internal/helpers"
	"os"
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
	ENV_CACHE_ADDR       = "CACHE_ADDR"
	ENV_CACHE_PASSWORD   = "CACHE_PASSWORD"
	ENV_CACHE_DB         = "CACHE_DB"

	defaultCacheDb        int32 = 0
	defaultPort           int32 = 9094
	defaultStorageTimeout       = time.Second * 5
	defaultMongoDb              = "squzy_monitoring"
	defaultCollection           = "schedulers"
)

type cfg struct {
	port            int32
	timeout         time.Duration
	clientAddress   string
	mongoURI        string
	mongoDb         string
	mongoCollection string
	cacheAddr       string
	cachePassword   string
	cacheDB         int32
}

func (c *cfg) GetPort() int32 {
	return c.port
}

func (c *cfg) GetClientAddress() string {
	return c.clientAddress
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

func (c *cfg) GetCacheAddr() string {
	return c.cacheAddr
}

func (c *cfg) GetCachePassword() string {
	return c.cachePassword
}

func (c *cfg) GetCacheDB() int32 {
	return c.cacheDB
}

type Config interface {
	GetPort() int32
	GetClientAddress() string
	GetStorageTimeout() time.Duration
	GetMongoURI() string
	GetMongoDb() string
	GetMongoCollection() string
	GetCacheAddr() string
	GetCachePassword() string
	GetCacheDB() int32
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

	cacheDBValue := os.Getenv(ENV_CACHE_DB)
	cacheDB := defaultCacheDb
	if cacheDBValue != "" {
		i, err := strconv.ParseInt(cacheDBValue, 10, 32)
		if err == nil {
			cacheDB = int32(i)
		}
	}
	return &cfg{
		clientAddress:   os.Getenv(ENV_STORAGE_HOST),
		timeout:         timeoutStorage,
		port:            port,
		mongoURI:        os.Getenv(ENV_MONGO_URI),
		mongoDb:         mongoDb,
		mongoCollection: collection,
		cacheAddr:       os.Getenv(ENV_CACHE_ADDR),
		cachePassword:   os.Getenv(ENV_CACHE_PASSWORD),
		cacheDB:         cacheDB,
	}
}
