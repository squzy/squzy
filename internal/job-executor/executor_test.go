package job_executor

import (
	"context"
	"crypto/tls"
	"errors"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"github.com/squzy/squzy/internal/httptools"
	"github.com/squzy/squzy/internal/job"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	"github.com/squzy/squzy/internal/semaphore"
	sitemap_storage "github.com/squzy/squzy/internal/sitemap-storage"
	"testing"
)

type externalStorageMock struct {
}

func (e externalStorageMock) Write(log job.CheckError) error {
	return nil
}

type configStorageMockOk struct {
	typeOfChecker apiPb.SchedulerType
}

func (c configStorageMockOk) Get(ctx context.Context, schedulerId primitive.ObjectID) (*scheduler_config_storage.SchedulerConfig, error) {
	return &scheduler_config_storage.SchedulerConfig{
		ID:   primitive.NewObjectID(),
		Type: c.typeOfChecker,
	}, nil
}

func (c configStorageMockOk) Add(ctx context.Context, config *scheduler_config_storage.SchedulerConfig) error {
	panic("implement me")
}

func (c configStorageMockOk) Remove(ctx context.Context, schedulerId primitive.ObjectID) error {
	panic("implement me")
}

func (c configStorageMockOk) Run(ctx context.Context, schedulerId primitive.ObjectID) error {
	panic("implement me")
}

func (c configStorageMockOk) Stop(ctx context.Context, schedulerId primitive.ObjectID) error {
	panic("implement me")
}

func (c configStorageMockOk) GetAll(ctx context.Context) ([]*scheduler_config_storage.SchedulerConfig, error) {
	panic("implement me")
}

func (c configStorageMockOk) GetAllForSync(ctx context.Context) ([]*scheduler_config_storage.SchedulerConfig, error) {
	panic("implement me")
}

type configStorageMockError struct {
}

func (c configStorageMockError) Get(ctx context.Context, schedulerId primitive.ObjectID) (*scheduler_config_storage.SchedulerConfig, error) {
	return nil, errors.New("cant get config")
}

func (c configStorageMockError) Add(ctx context.Context, config *scheduler_config_storage.SchedulerConfig) error {
	panic("implement me")
}

func (c configStorageMockError) Remove(ctx context.Context, schedulerId primitive.ObjectID) error {
	panic("implement me")
}

func (c configStorageMockError) Run(ctx context.Context, schedulerId primitive.ObjectID) error {
	panic("implement me")
}

func (c configStorageMockError) Stop(ctx context.Context, schedulerId primitive.ObjectID) error {
	panic("implement me")
}

func (c configStorageMockError) GetAll(ctx context.Context) ([]*scheduler_config_storage.SchedulerConfig, error) {
	panic("implement me")
}

func (c configStorageMockError) GetAllForSync(ctx context.Context) ([]*scheduler_config_storage.SchedulerConfig, error) {
	panic("implement me")
}

type fnMock struct {
	executed bool
}

func (m *fnMock) TcpMock(schedulerId string, timeout int32, config *scheduler_config_storage.TCPConfig) job.CheckError {
	m.executed = true
	return nil
}

func (m *fnMock) SiteMapMock(schedulerId string, timeout int32, config *scheduler_config_storage.SiteMapConfig, siteMapStorage sitemap_storage.SiteMapStorage, httpTools httptools.HTTPTool, semaphoreFactoryFn func(n int) semaphore.Semaphore) job.CheckError {
	m.executed = true
	return nil
}

func (m *fnMock) GrpcMock(schedulerId string, timeout int32, config *scheduler_config_storage.GrpcConfig, opts ...grpc.DialOption) job.CheckError {
	m.executed = true
	return nil
}

func (m *fnMock) HttpMock(schedulerId string, timeout int32, config *scheduler_config_storage.HTTPConfig, httpTool httptools.HTTPTool) job.CheckError {
	m.executed = true
	return nil
}

func (m *fnMock) HttpValueMock(schedulerId string, timeout int32, config *scheduler_config_storage.HTTPValueConfig, httpTool httptools.HTTPTool) job.CheckError {
	m.executed = true
	return nil
}

func (m *fnMock) SSLExpirationMock(
	schedulerId string,
	timeout int32,
	config *scheduler_config_storage.SslExpirationConfig,
	cfg *tls.Config,
) job.CheckError {
	m.executed = true
	return nil
}

func TestNewExecutor(t *testing.T) {
	t.Run("Should: implement interface", func(t *testing.T) {
		s := NewExecutor(
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		assert.Implements(t, (*JobExecutor)(nil), s)
	})
}

func TestExecutor_Execute(t *testing.T) {
	t.Run("Should:  not execute because cant get config", func(t *testing.T) {
		fnMock := &fnMock{}
		s := NewExecutor(
			nil,
			nil,
			nil,
			nil,
			&configStorageMockError{},
			fnMock.TcpMock,
			fnMock.GrpcMock,
			fnMock.HttpMock,
			fnMock.SiteMapMock,
			fnMock.HttpValueMock,
			fnMock.SSLExpirationMock,
		)
		s.Execute(primitive.NewObjectID())
		assert.Equal(t, false, fnMock.executed)
	})
	t.Run("Should: execute tcp mock", func(t *testing.T) {
		fnMock := &fnMock{}
		s := NewExecutor(
			&externalStorageMock{},
			nil,
			nil,
			nil,
			&configStorageMockOk{
				apiPb.SchedulerType_TCP,
			},
			fnMock.TcpMock,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		s.Execute(primitive.NewObjectID())
		assert.Equal(t, true, fnMock.executed)
	})
	t.Run("Should: execute grpc mock", func(t *testing.T) {
		fnMock := &fnMock{}
		s := NewExecutor(
			&externalStorageMock{},
			nil,
			nil,
			nil,
			&configStorageMockOk{
				apiPb.SchedulerType_GRPC,
			},
			nil,
			fnMock.GrpcMock,
			nil,
			nil,
			nil,
			nil,
		)
		s.Execute(primitive.NewObjectID())
		assert.Equal(t, true, fnMock.executed)
	})
	t.Run("Should: execute http mock", func(t *testing.T) {
		fnMock := &fnMock{}
		s := NewExecutor(
			&externalStorageMock{},
			nil,
			nil,
			nil,
			&configStorageMockOk{
				apiPb.SchedulerType_HTTP,
			},
			nil,
			nil,
			fnMock.HttpMock,
			nil,
			nil,
			nil,
		)
		s.Execute(primitive.NewObjectID())
		assert.Equal(t, true, fnMock.executed)
	})
	t.Run("Should: execute sitemap mock", func(t *testing.T) {
		fnMock := &fnMock{}
		s := NewExecutor(
			&externalStorageMock{},
			nil,
			nil,
			nil,
			&configStorageMockOk{
				apiPb.SchedulerType_SITE_MAP,
			},
			nil,
			nil,
			nil,
			fnMock.SiteMapMock,
			nil,
			nil,
		)
		s.Execute(primitive.NewObjectID())
		assert.Equal(t, true, fnMock.executed)
	})
	t.Run("Should: execute sslexpirztion mock", func(t *testing.T) {
		fnMock := &fnMock{}
		s := NewExecutor(
			&externalStorageMock{},
			nil,
			nil,
			nil,
			&configStorageMockOk{
				apiPb.SchedulerType_SSL_EXPIRATION,
			},
			nil,
			nil,
			nil,
			nil,
			nil,
			fnMock.SSLExpirationMock,
		)
		s.Execute(primitive.NewObjectID())
		assert.Equal(t, true, fnMock.executed)
	})
	t.Run("Should: execute httpvalue mock", func(t *testing.T) {
		fnMock := &fnMock{}
		s := NewExecutor(
			&externalStorageMock{},
			nil,
			nil,
			nil,
			&configStorageMockOk{
				apiPb.SchedulerType_HTTP_JSON_VALUE,
			},
			nil,
			nil,
			nil,
			nil,
			fnMock.HttpValueMock,
			nil,
		)
		s.Execute(primitive.NewObjectID())
		assert.Equal(t, true, fnMock.executed)
	})
	t.Run("Should: nothing execute", func(t *testing.T) {
		fnMock := &fnMock{}
		s := NewExecutor(
			&externalStorageMock{},
			nil,
			nil,
			nil,
			&configStorageMockOk{
				11111,
			},
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		s.Execute(primitive.NewObjectID())
		assert.Equal(t, false, fnMock.executed)
	})
}
