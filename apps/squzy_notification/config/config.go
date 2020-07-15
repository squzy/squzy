package config

type Config interface {
	GetPort() int32
	GetMongoURI() string
	GetMongoDB() string
	GetNotificationMethodCollection() string
	GetNotificationListCollection() string
	GetStorageHost() string
}

type cfg struct {
}

func (c *cfg) GetStorageHost() string {
	panic("implement me")
}

func (c *cfg) GetPort() int32 {
	panic("implement me")
}

func (c *cfg) GetMongoURI() string {
	panic("implement me")
}

func (c *cfg) GetMongoDB() string {
	panic("implement me")
}

func (c *cfg) GetNotificationMethodCollection() string {
	panic("implement me")
}

func (c *cfg) GetNotificationListCollection() string {
	panic("implement me")
}

func New() Config {
	return &cfg{}
}
