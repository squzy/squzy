package config

import (
	"os"
	"strconv"
)

const (
	ENV_PORT             = "PORT"
	ENV_STRORAGE_HOST    = "ENV_STRORAGE_HOST"
	ENV_MONGO_DB         = "MONGO_DB"
	ENV_MONGO_URI        = "MONGO_URI"
	ENV_MONGO_COLLECTION = "MONGO_COLLECTION"

	defaultPort       int32 = 9097
	defaultMongoDb          = "incident_manager"
	defaultCollection       = "rules"
)

type cfg struct {
	port            int32
	storageHost     string
	mongoURI        string
	mongoDb         string
	mongoCollection string
}

func (c *cfg) GetPort() int32 {
	return c.port
}

func (c *cfg) GetStorageHost() string {
	return c.storageHost
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
	GetStorageHost() string
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

	mongoDb := os.Getenv(ENV_MONGO_DB)
	if mongoDb == "" {
		mongoDb = defaultMongoDb
	}
	collection := os.Getenv(ENV_MONGO_COLLECTION)
	if collection == "" {
		collection = defaultCollection
	}

	return &cfg{
		port:            port,
		storageHost:     os.Getenv(ENV_STRORAGE_HOST),
		mongoURI:        os.Getenv(ENV_MONGO_URI),
		mongoDb:         mongoDb,
		mongoCollection: collection,
	}
}
