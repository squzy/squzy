package version

import (
	"squzy/internal/logger"
)

var (
	Version = "local"
)

func init() {
	log.Println("Version: " + GetVersion())
}

func GetVersion() string {
	return Version
}
