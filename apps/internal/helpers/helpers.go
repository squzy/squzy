package helpers

import "strings"

const (
	httpPort                      = int32(80)
	httpsPort                     = int32(443)
)

func GetPortByUrl(url string) int32 {
	if strings.HasPrefix(url, "https") {
		return httpsPort
	}
	return httpPort
}