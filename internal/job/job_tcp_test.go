package job

import (
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestNewTcpJob(t *testing.T) {
	t.Run("Should: Should implement interface Job", func(t *testing.T) {
		job := NewTcpJob("localhost", 9090, 0)
		assert.Implements(t, (*Job)(nil), job)
	})
}

func TestJobTcp_Do(t *testing.T) {
	t.Run("Test: Testing tcp health_check:", func(t *testing.T) {
		t.Run("Should: return wrongConnectConfigError", func(t *testing.T) {
			job := NewTcpJob("localhost", 10002, 0)
			server, _ := net.Listen("tcp", "localhost:10003")
			defer server.Close()
			assert.Equal(t, wrongConnectConfigError.Error(), job.Do("").GetLogData().Error.Message)
		})
		t.Run("Should: return nil", func(t *testing.T) {
			job := NewTcpJob("localhost", 10003, 0)
			server, _ := net.Listen("tcp", "localhost:10003")
			go func() {
				_, _ = server.Accept()
			}()
			defer server.Close()
			assert.Equal(t, apiPb.SchedulerResponseCode_OK, job.Do("").GetLogData().Code)
		})
		t.Run("Should: return error because timeout", func(t *testing.T) {
			job := NewTcpJob("localhost", 10004, 1)
			go func() {
				time.Sleep(time.Second * 5)
				server, _ := net.Listen("tcp", "localhost:10004")
				_, _ = server.Accept()
				defer server.Close()
			}()
			assert.Equal(t, apiPb.SchedulerResponseCode_Error, job.Do("").GetLogData().Code)
		})
	})
}
