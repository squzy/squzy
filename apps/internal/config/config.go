package config

type Config interface {
}

type cfg struct {
}

func NewConfig() Config {
	return &cfg{}
}