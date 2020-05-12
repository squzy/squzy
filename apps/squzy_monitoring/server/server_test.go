package server

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"squzy/internal/scheduler"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
	"testing"
)

var (
	successTcpConfig = &scheduler_config_storage.SchedulerConfig{
		ID:       primitive.ObjectID{},
		Type:     apiPb.SchedulerType_Tcp,
		Status:   0,
		Interval: 0,
		Timeout:  0,
		TCPConfig: &scheduler_config_storage.TCPConfig{
			Host: "",
			Port: 0,
		},
	}

	successHttpConfig = &scheduler_config_storage.SchedulerConfig{
		ID:       primitive.NewObjectID(),
		Type:     apiPb.SchedulerType_Http,
		Status:   0,
		Interval: 0,
		Timeout:  0,
		HTTPConfig: &scheduler_config_storage.HTTPConfig{
			Method:     "",
			URL:        "",
			Headers:    nil,
			StatusCode: 0,
		},
	}

	successHttpValueConfig = &scheduler_config_storage.SchedulerConfig{
		ID:       primitive.NewObjectID(),
		Type:     apiPb.SchedulerType_HttpJsonValue,
		Status:   0,
		Interval: 0,
		Timeout:  0,
		HTTPValueConfig: &scheduler_config_storage.HTTPValueConfig{
			Method:    "",
			URL:       "",
			Headers:   nil,
			Selectors: nil,
		},
	}

	successGrpcConfig = &scheduler_config_storage.SchedulerConfig{
		ID:       primitive.NewObjectID(),
		Type:     apiPb.SchedulerType_Grpc,
		Status:   0,
		Interval: 0,
		Timeout:  0,
		GrpcConfig: &scheduler_config_storage.GrpcConfig{
			Service: "",
			Host:    "",
			Port:    0,
		},
	}

	successSiteMapConfig = &scheduler_config_storage.SchedulerConfig{
		ID:       primitive.NewObjectID(),
		Type:     apiPb.SchedulerType_SiteMap,
		Status:   0,
		Interval: 0,
		Timeout:  0,
		SiteMapConfig: &scheduler_config_storage.SiteMapConfig{
			URL:         "",
			Concurrency: 0,
		},
	}

	errorConfig = &scheduler_config_storage.SchedulerConfig{
		ID:       primitive.NewObjectID(),
		Type:     11111,
		Status:   0,
		Interval: 0,
		Timeout:  0,
	}

	cfgMap = map[primitive.ObjectID]*scheduler_config_storage.SchedulerConfig{
		successTcpConfig.ID:       successTcpConfig,
		successGrpcConfig.ID:      successGrpcConfig,
		successHttpConfig.ID:      successHttpConfig,
		successHttpValueConfig.ID: successHttpValueConfig,
		successSiteMapConfig.ID:   successSiteMapConfig,
		errorConfig.ID:            errorConfig,
	}

	rqMap = map[apiPb.SchedulerType]*apiPb.AddRequest{
		apiPb.SchedulerType_Tcp: {
			Interval: 10,
			Timeout:  0,
			Config: &apiPb.AddRequest_Tcp{
				Tcp: &apiPb.TcpConfig{
					Host:                 "",
					Port:                 0,
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
			},
		},
		apiPb.SchedulerType_Http: {
			Interval: 10,
			Timeout:  0,
			Config: &apiPb.AddRequest_Http{
				Http: &apiPb.HttpConfig{
					Method:     "",
					Url:        "",
					Headers:    nil,
					StatusCode: 0,
				},
			},
		},
		apiPb.SchedulerType_HttpJsonValue: {
			Interval: 10,
			Timeout:  0,
			Config: &apiPb.AddRequest_HttpValue{
				HttpValue: &apiPb.HttpJsonValueConfig{
					Method:    "",
					Url:       "",
					Headers:   nil,
					Selectors: nil,
				},
			},
		},
		apiPb.SchedulerType_SiteMap: {
			Interval: 10,
			Timeout:  0,
			Config: &apiPb.AddRequest_Sitemap{
				Sitemap: &apiPb.SiteMapConfig{
					Url:         "",
					Concurrency: 0,
				},
			},
		},
		apiPb.SchedulerType_Grpc: {
			Interval: 10,
			Timeout:  0,
			Config: &apiPb.AddRequest_Grpc{
				Grpc: &apiPb.GrpcConfig{
					Service: "",
					Host:    "",
					Port:    0,
				},
			},
		},
		1000: {
			Interval: 10,
			Timeout:  0,
		},
	}
)

type schedulerMock struct {
}

func (s schedulerMock) GetID() string {
	panic("implement me")
}

func (s schedulerMock) GetIDBson() primitive.ObjectID {
	panic("implement me")
}

func (s schedulerMock) Run() {
}

func (s schedulerMock) Stop() {
}

func (s schedulerMock) IsRun() bool {
	panic("implement me")
}

type mockStorageOk struct {
}

func (m mockStorageOk) Get(string) (scheduler.Scheduler, error) {
	return &schedulerMock{}, nil
}

func (m mockStorageOk) Set(scheduler.Scheduler) error {
	return nil
}

func (m mockStorageOk) Remove(string) error {
	return nil
}

type mockStorageError struct {
}

func (m mockStorageError) Get(string) (scheduler.Scheduler, error) {
	return nil, errors.New("")
}

func (m mockStorageError) Set(scheduler.Scheduler) error {
	return errors.New("")
}

func (m mockStorageError) Remove(string) error {
	return errors.New("")
}

type mockConfigStorageOk struct {
}

func (m mockConfigStorageOk) Get(ctx context.Context, schedulerId primitive.ObjectID) (*scheduler_config_storage.SchedulerConfig, error) {
	return cfgMap[schedulerId], nil
}

func (m mockConfigStorageOk) Add(ctx context.Context, config *scheduler_config_storage.SchedulerConfig) error {
	return nil
}

func (m mockConfigStorageOk) Remove(ctx context.Context, schedulerId primitive.ObjectID) error {
	return nil
}

func (m mockConfigStorageOk) Run(ctx context.Context, schedulerId primitive.ObjectID) error {
	return nil
}

func (m mockConfigStorageOk) Stop(ctx context.Context, schedulerId primitive.ObjectID) error {
	return nil
}

func (m mockConfigStorageOk) GetAll(ctx context.Context) ([]*scheduler_config_storage.SchedulerConfig, error) {
	return []*scheduler_config_storage.SchedulerConfig{
		{
			ID: successGrpcConfig.ID,
		},
	}, nil
}

func (m mockConfigStorageOk) GetAllForSync(ctx context.Context) ([]*scheduler_config_storage.SchedulerConfig, error) {
	panic("implement me")
}

type mockConfigStorageErrorSingle struct {
}

func (m mockConfigStorageErrorSingle) Get(ctx context.Context, schedulerId primitive.ObjectID) (*scheduler_config_storage.SchedulerConfig, error) {
	return nil, errors.New("")
}

func (m mockConfigStorageErrorSingle) Add(ctx context.Context, config *scheduler_config_storage.SchedulerConfig) error {
	return errors.New("")
}

func (m mockConfigStorageErrorSingle) Remove(ctx context.Context, schedulerId primitive.ObjectID) error {
	return errors.New("")
}

func (m mockConfigStorageErrorSingle) Run(ctx context.Context, schedulerId primitive.ObjectID) error {
	return errors.New("")
}

func (m mockConfigStorageErrorSingle) Stop(ctx context.Context, schedulerId primitive.ObjectID) error {
	return errors.New("")
}

func (m mockConfigStorageErrorSingle) GetAll(ctx context.Context) ([]*scheduler_config_storage.SchedulerConfig, error) {
	return []*scheduler_config_storage.SchedulerConfig{
		{
			ID: primitive.NewObjectID(),
		},
	}, nil
}

func (m mockConfigStorageErrorSingle) GetAllForSync(ctx context.Context) ([]*scheduler_config_storage.SchedulerConfig, error) {
	panic("implement me")
}

type mockConfigStorageError struct {
}

func (m mockConfigStorageError) Get(ctx context.Context, schedulerId primitive.ObjectID) (*scheduler_config_storage.SchedulerConfig, error) {
	panic("implement me")
}

func (m mockConfigStorageError) Add(ctx context.Context, config *scheduler_config_storage.SchedulerConfig) error {
	panic("implement me")
}

func (m mockConfigStorageError) Remove(ctx context.Context, schedulerId primitive.ObjectID) error {
	panic("implement me")
}

func (m mockConfigStorageError) Run(ctx context.Context, schedulerId primitive.ObjectID) error {
	panic("implement me")
}

func (m mockConfigStorageError) Stop(ctx context.Context, schedulerId primitive.ObjectID) error {
	panic("implement me")
}

func (m mockConfigStorageError) GetAll(ctx context.Context) ([]*scheduler_config_storage.SchedulerConfig, error) {
	return nil, errors.New("")
}

func (m mockConfigStorageError) GetAllForSync(ctx context.Context) ([]*scheduler_config_storage.SchedulerConfig, error) {
	panic("implement me")
}

func TestNew(t *testing.T) {
	t.Run("Should: implement interface", func(t *testing.T) {
		s := New(nil, nil, nil)
		assert.Implements(t, (*apiPb.SchedulersExecutorServer)(nil), s)
	})
}

func TestServer_GetSchedulerList(t *testing.T) {
	t.Run("Should: return error because DB", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageError{})
		_, err := s.GetSchedulerList(context.Background(), &empty.Empty{})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because sinle DB error", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageErrorSingle{})
		_, err := s.GetSchedulerList(context.Background(), &empty.Empty{})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return without error", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageOk{})
		_, err := s.GetSchedulerList(context.Background(), &empty.Empty{})
		assert.Equal(t, nil, err)
	})
}

func TestServer_GetSchedulerById(t *testing.T) {
	t.Run("Should: return error because DB", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageErrorSingle{})
		_, err := s.GetSchedulerById(context.Background(), &apiPb.GetSchedulerByIdRequest{
			Id: "",
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return tcp config", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageOk{})
		_, err := s.GetSchedulerById(context.Background(), &apiPb.GetSchedulerByIdRequest{
			Id: successTcpConfig.ID.Hex(),
		})
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return grpc config", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageOk{})
		_, err := s.GetSchedulerById(context.Background(), &apiPb.GetSchedulerByIdRequest{
			Id: successGrpcConfig.ID.Hex(),
		})
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return http config", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageOk{})
		_, err := s.GetSchedulerById(context.Background(), &apiPb.GetSchedulerByIdRequest{
			Id: successHttpConfig.ID.Hex(),
		})
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return sitemap config", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageOk{})
		_, err := s.GetSchedulerById(context.Background(), &apiPb.GetSchedulerByIdRequest{
			Id: successSiteMapConfig.ID.Hex(),
		})
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return httpValue config", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageOk{})
		_, err := s.GetSchedulerById(context.Background(), &apiPb.GetSchedulerByIdRequest{
			Id: successHttpValueConfig.ID.Hex(),
		})
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return error because not correct typw", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageOk{})
		_, err := s.GetSchedulerById(context.Background(), &apiPb.GetSchedulerByIdRequest{
			Id: errorConfig.ID.Hex(),
		})
		assert.NotEqual(t, nil, err)
	})
}

func TestServer_Run(t *testing.T) {
	t.Run("Should: return error because id not bson", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageOk{})
		_, err := s.Run(context.Background(), &apiPb.RunRequest{
			Id: "sff",
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because id not found in DB", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageErrorSingle{})
		_, err := s.Run(context.Background(), &apiPb.RunRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because cant find in memory", func(t *testing.T) {
		s := New(&mockStorageError{}, nil, &mockConfigStorageOk{})
		_, err := s.Run(context.Background(), &apiPb.RunRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&mockStorageOk{}, nil, &mockConfigStorageOk{})
		_, err := s.Run(context.Background(), &apiPb.RunRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.Equal(t, nil, err)
	})
}

func TestServer_Stop(t *testing.T) {
	t.Run("Should: return error because id not bson", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageOk{})
		_, err := s.Stop(context.Background(), &apiPb.StopRequest{
			Id: "sff",
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because id not found in DB", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageErrorSingle{})
		_, err := s.Stop(context.Background(), &apiPb.StopRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because cant find in memory", func(t *testing.T) {
		s := New(&mockStorageError{}, nil, &mockConfigStorageOk{})
		_, err := s.Stop(context.Background(), &apiPb.StopRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&mockStorageOk{}, nil, &mockConfigStorageOk{})
		_, err := s.Stop(context.Background(), &apiPb.StopRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.Equal(t, nil, err)
	})
}

func TestServer_Remove(t *testing.T) {
	t.Run("Should: return error because id not bson", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageOk{})
		_, err := s.Remove(context.Background(), &apiPb.RemoveRequest{
			Id: "sff",
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because id not found in DB", func(t *testing.T) {
		s := New(nil, nil, &mockConfigStorageErrorSingle{})
		_, err := s.Remove(context.Background(), &apiPb.RemoveRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because cant find in memory", func(t *testing.T) {
		s := New(&mockStorageError{}, nil, &mockConfigStorageOk{})
		_, err := s.Remove(context.Background(), &apiPb.RemoveRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&mockStorageOk{}, nil, &mockConfigStorageOk{})
		_, err := s.Remove(context.Background(), &apiPb.RemoveRequest{
			Id: primitive.NewObjectID().Hex(),
		})
		assert.Equal(t, nil, err)
	})
}

func TestServer_Add(t *testing.T) {
	t.Run("Should: return error because wrong interval", func(t *testing.T) {
		s := New(nil, nil, nil)
		_, err := s.Add(context.Background(), &apiPb.AddRequest{
			Interval: 0,
			Timeout:  0,
			Config:   nil,
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because wrong type", func(t *testing.T) {
		s := New(nil, nil, nil)
		_, err := s.Add(context.Background(), rqMap[1000])
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because cant add to DB", func(t *testing.T) {
		s := New(&mockStorageOk{}, nil, &mockConfigStorageErrorSingle{})
		_, err := s.Add(context.Background(), rqMap[apiPb.SchedulerType_Tcp])
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because cant add to in memory", func(t *testing.T) {
		s := New(&mockStorageError{}, nil, &mockConfigStorageOk{})
		_, err := s.Add(context.Background(), rqMap[apiPb.SchedulerType_Tcp])
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: add tcp check without error", func(t *testing.T) {
		s := New(&mockStorageOk{}, nil, &mockConfigStorageOk{})
		_, err := s.Add(context.Background(), rqMap[apiPb.SchedulerType_Tcp])
		assert.Equal(t, nil, err)
	})
	t.Run("Should: add grcp check without error", func(t *testing.T) {
		s := New(&mockStorageOk{}, nil, &mockConfigStorageOk{})
		_, err := s.Add(context.Background(), rqMap[apiPb.SchedulerType_Grpc])
		assert.Equal(t, nil, err)
	})
	t.Run("Should: add sitemap check without error", func(t *testing.T) {
		s := New(&mockStorageOk{}, nil, &mockConfigStorageOk{})
		_, err := s.Add(context.Background(), rqMap[apiPb.SchedulerType_SiteMap])
		assert.Equal(t, nil, err)
	})
	t.Run("Should: add httpValue check without error", func(t *testing.T) {
		s := New(&mockStorageOk{}, nil, &mockConfigStorageOk{})
		_, err := s.Add(context.Background(), rqMap[apiPb.SchedulerType_HttpJsonValue])
		assert.Equal(t, nil, err)
	})
	t.Run("Should: add http check without error", func(t *testing.T) {
		s := New(&mockStorageOk{}, nil, &mockConfigStorageOk{})
		_, err := s.Add(context.Background(), rqMap[apiPb.SchedulerType_Http])
		assert.Equal(t, nil, err)
	})
}
