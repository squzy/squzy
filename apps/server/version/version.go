package version

import (
	"fmt"
	"log"
)

var (
	Version = "local"
)

func init() {
	log.Println(fmt.Sprintf("Version: %s", GetVersion()))
}

func GetVersion() string {
	return Version
}
