package version

import (
	"squzy/internal/logger"
)

var (
	Version = "local"
)

func init() {
	logger.Infof("Version: " + GetVersion())
}

func GetVersion() string {
	return Version
}
