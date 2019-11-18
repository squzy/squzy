package job

import (
	"net"
	"fmt"
)

type jobTcp struct {
	port int16
	host string
}

func NewTcpJob(host string, port int16) Job {
	return &jobTcp{
		port:    port,
		host: host,
	}
}

func (j *jobTcp) Do() error {
 	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", j.host, j.port), connTimeout)
 	if err != nil {
 		return wrongConnectConfigError
	}
	if conn != nil {
		defer func() {
			_ = conn.Close()
		}()
	}
	return nil
}