package version

import (
	"log"
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