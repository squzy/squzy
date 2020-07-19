package config

import (
	"os"
	"strconv"
)

const (
	ENV_PORT             = "PORT"
	ENV_STORAGE_HOST     = "STORAGE_HOST"
	ENV_MONGO_DB         = "MONGO_DB"
	ENV_MONGO_URI        = "MONGO_URI"
	ENV_MONGO_LIST_COLLECTION = "MONGO_LIST_COLLECTION"
	ENV_MONGO_METHOD_COLLECTION = "MONGO_METHOD_COLLECTION"
	ENV_DASHBOARD_HOST = "DASHBOARD_HOST"

	defaultPort       int32 = 9098
	defaultMongoDb          = "notification_manager"
	defaultListCollection       = "list"
	defaultMethodCollection       = "methods"
)

type Config interface {
	GetPort() int32
	GetMongoURI() string
	GetMongoDB() string
	GetNotificationMethodCollection() string
	GetNotificationListCollection() string
	GetStorageHost() string
	GetDashboardHost() string
}

type cfg struct {
	storageHost string
	port int32
	mongoDB string
	mongoURI string
	notificationListCollection string
	notificationMethodCollection string
	dashboardHost string
}

func (c *cfg) GetDashboardHost() string {
	return c.dashboardHost
}

func (c *cfg) GetStorageHost() string {
	return c.storageHost
}

func (c *cfg) GetPort() int32 {
	return c.port
}

func (c *cfg) GetMongoURI() string {
	return c.mongoURI
}

func (c *cfg) GetMongoDB() string {
	return c.mongoDB
}

func (c *cfg) GetNotificationMethodCollection() string {
	return c.notificationMethodCollection
}

func (c *cfg) GetNotificationListCollection() string {
	return c.notificationListCollection
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
	listCollection := os.Getenv(ENV_MONGO_LIST_COLLECTION)
	if listCollection == "" {
		listCollection = defaultListCollection
	}
	methodCollection := os.Getenv(ENV_MONGO_METHOD_COLLECTION)
	if methodCollection == "" {
		methodCollection = defaultMethodCollection
	}

	return &cfg{
		port:            port,
		storageHost:     os.Getenv(ENV_STORAGE_HOST),
		mongoURI:        os.Getenv(ENV_MONGO_URI),
		mongoDB:         mongoDb,
		notificationListCollection: listCollection,
		notificationMethodCollection: methodCollection,
		dashboardHost: os.Getenv(ENV_DASHBOARD_HOST),
	}
}
