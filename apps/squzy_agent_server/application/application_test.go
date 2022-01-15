package application

import (
	"context"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"net"
	"testing"
	"time"
)

type mock struct {
}

func (m mock) Register(ctx context.Context, request *apiPb.RegisterRequest) (*apiPb.RegisterResponse, error) {
	panic("implement me")
}

func (m mock) GetByAgentName(ctx context.Context, request *apiPb.GetByAgentNameRequest) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (m mock) GetAgentById(ctx context.Context, request *apiPb.GetAgentByIdRequest) (*apiPb.AgentItem, error) {
	panic("implement me")
}

func (m mock) UnRegister(ctx context.Context, request *apiPb.UnRegisterRequest) (*apiPb.UnRegisterResponse, error) {
	panic("implement me")
}

func (m mock) GetAgentList(ctx context.Context, empty *empty.Empty) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (m mock) SendMetrics(server apiPb.AgentServer_SendMetricsServer) error {
	panic("implement me")
}

func TestNew(t *testing.T) {
	t.Run("Should: not be nil", func(t *testing.T) {
		s := New(nil)
		assert.NotEqual(t, nil, s)
	})
}

func TestApp_Run(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		app := New(&mock{})
		go func() {
			_ = app.Run(11101)
		}()
		time.Sleep(time.Second)
		_, err := net.Dial("tcp", "localhost:11101")
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		app := New(nil)
		assert.NotEqual(t, nil, app.Run(1231323))
	})
}
