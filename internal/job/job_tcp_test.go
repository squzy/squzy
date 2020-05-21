package job

import (
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"net"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
	"testing"
	"time"
)

func TestExecTcp(t *testing.T) {
	t.Run("Test: Testing tcp health_check:", func(t *testing.T) {
		t.Run("Should: return errWrongConnectConfigError", func(t *testing.T) {
			server, _ := net.Listen("tcp", "localhost:10003")
			defer server.Close()
			job := ExecTCP("", 0, &scheduler_config_storage.TCPConfig{
				Host: "localhost",
				Port: 10002,
			})
			assert.Equal(t, errWrongConnectConfigError.Error(), job.GetLogData().Snapshot.Error.Message)
		})
		t.Run("Should: return nil", func(t *testing.T) {
			server, err := net.Listen("tcp", "localhost:10003")
			assert.Equal(t, nil, err)
			go func() {
				_, _ = server.Accept()
			}()
			defer server.Close()
			job := ExecTCP("", 0, &scheduler_config_storage.TCPConfig{
				Host: "localhost",
				Port: 10003,
			})
			assert.Equal(t, apiPb.SchedulerCode_OK, job.GetLogData().Snapshot.Code)
		})
		t.Run("Should: return error because timeout", func(t *testing.T) {
			go func() {
				time.Sleep(time.Second * 5)
				server, _ := net.Listen("tcp", "localhost:10004")
				_, _ = server.Accept()
				defer server.Close()
			}()
			job := ExecTCP("", 1, &scheduler_config_storage.TCPConfig{
				Host: "localhost",
				Port: 10004,
			})
			assert.Equal(t, apiPb.SchedulerCode_ERROR, job.GetLogData().Snapshot.Code)
		})
	})
}
