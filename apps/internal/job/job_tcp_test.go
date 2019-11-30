package job

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
)

func TestNewTcpJob(t *testing.T) {
	t.Run("Should: Should implement interface Job", func(t *testing.T) {
		job := NewTcpJob("localhost", 9090)
		assert.Implements(t, (*Job)(nil), job)
	})
}

func TestJobTcp_Do(t *testing.T) {
	t.Run("Test: Testing tcp health_check:", func(t *testing.T) {
		t.Run("Should: return wrongConnectConfigError", func(t *testing.T) {
			job := NewTcpJob("localhost", 10002)
			server, _ := net.Listen("tcp", "localhost:10003")
			defer server.Close()
			assert.Equal(t, wrongConnectConfigError.Error(), job.Do().GetLogData().Description)
		})
		t.Run("Should: return nil", func(t *testing.T) {
			job := NewTcpJob("localhost", 10003)
			server, _ := net.Listen("tcp", "localhost:10003")
			go func() {
				server.Accept()
			}()
			defer server.Close()
			assert.Equal(t, clientPb.StatusCode_OK, job.Do().GetLogData().Code)
		})
	})
}