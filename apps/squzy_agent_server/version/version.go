package version

import (
	"github.com/squzy/squzy/internal/logger"
)

var (
	Version = "local"
)

func init() {
	logger.Info("Version: " + GetVersion())
}

func GetVersion() string {
	return Version
}
