package version

import "fmt"

var (
	Version = "version"
)

func init() {
	fmt.Println(Version)
}

func GetVersion() string {
	return Version
}