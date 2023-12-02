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
	apiPb.UnimplementedApplicationMonitoringServer
}

func (m mock) GetApplicationListByAgentId(ctx context.Context, request *apiPb.AgentIdRequest) (*apiPb.GetApplicationListResponse, error) {
	panic("implement me")
}

func (m mock) ArchiveApplicationById(ctx context.Context, reuqest *apiPb.ApplicationByIdReuqest) (*apiPb.Application, error) {
	panic("implement me")
}

func (m mock) EnableApplicationById(ctx context.Context, reuqest *apiPb.ApplicationByIdReuqest) (*apiPb.Application, error) {
	panic("implement me")
}

func (m mock) DisableApplicationById(ctx context.Context, reuqest *apiPb.ApplicationByIdReuqest) (*apiPb.Application, error) {
	panic("implement me")
}

func (m mock) InitializeApplication(ctx context.Context, info *apiPb.ApplicationInfo) (*apiPb.InitializeApplicationResponse, error) {
	panic("implement me")
}

func (m mock) SaveTransaction(ctx context.Context, info *apiPb.TransactionInfo) (*empty.Empty, error) {
	panic("implement me")
}

func (m mock) GetApplicationById(ctx context.Context, request *apiPb.ApplicationByIdReuqest) (*apiPb.Application, error) {
	panic("implement me")
}

func (m mock) GetApplicationList(ctx context.Context, empty *empty.Empty) (*apiPb.GetApplicationListResponse, error) {
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
			_ = app.Run(11102)
		}()
		time.Sleep(time.Second)
		_, err := net.Dial("tcp", "localhost:11102")
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		app := New(nil)
		assert.NotEqual(t, nil, app.Run(1231323))
	})
}
