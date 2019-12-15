package job

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"net"
)

type tcpError struct {
	time        *timestamp.Timestamp
	code        clientPb.StatusCode
	description string
	location    string
	port        int32
}

func (s *tcpError) GetLogData() *clientPb.Log {
	return &clientPb.Log{
		Code:        s.code,
		Description: s.description,
		Meta: &clientPb.MetaData{
			Id:       uuid.New().String(),
			Location: s.location,
			Port:     s.port,
				Time:     s.time,
				Type:     clientPb.Type_Tcp,
			},
		}
}

func newTcpError(time *timestamp.Timestamp, code clientPb.StatusCode, description string, location string, port int32) CheckError {
	return &tcpError{
		time:        time,
		code:        code,
		description: description,
		location:    location,
		port:        port,
	}
}

type jobTcp struct {
	port int32
	host string
}

func NewTcpJob(host string, port int32) Job {
	return &jobTcp{
		port: port,
		host: host,
	}
}

func (j *jobTcp) Do() CheckError {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", j.host, j.port), connTimeout)
	if err != nil {
		return newTcpError(ptypes.TimestampNow(), clientPb.StatusCode_Error, wrongConnectConfigError.Error(), j.host, j.port)
	}
	if conn != nil {
		defer func() {
			_ = conn.Close()
		}()
	}
	return newTcpError(ptypes.TimestampNow(), clientPb.StatusCode_OK, "", j.host, j.port)
}
