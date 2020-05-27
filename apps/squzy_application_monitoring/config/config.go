package config

type Config interface {
	GetTracingHeader() string
}

type cfg struct {
	tracingHeader string
}

func (c *cfg) GetTracingHeader() string {
	return c.tracingHeader
}

func New() Config {
	return &cfg{
		tracingHeader: "Squzy_tracing",
	}
}