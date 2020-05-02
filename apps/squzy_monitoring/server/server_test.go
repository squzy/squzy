package server

import (
	"context"
	"errors"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	"squzy/internal/job"
	"squzy/internal/parsers"
	"squzy/internal/scheduler"
	"squzy/internal/semaphore"
	"testing"
	"time"
)

type mockSchedulerStorageError struct {
}

type mockSchedulerError struct {
}

type mockSchedulerStorageGetError struct {
}

func (m mockSchedulerStorageGetError) Get(string) (scheduler.Scheduler, error) {
	return nil, errors.New("safsaf")
}

func (m mockSchedulerStorageGetError) Set(scheduler.Scheduler) error {
	panic("implement me")
}

func (m mockSchedulerStorageGetError) Remove(string) error {
	panic("implement me")
}

func (m mockSchedulerStorageGetError) GetList() map[string]bool {
	panic("implement me")
}

func (m mockSchedulerError) GetId() string {
	panic("implement me")
}

func (m mockSchedulerError) Run() error {
	return errors.New("e")
}

func (m mockSchedulerError) Stop() error {
	return errors.New("e")
}

func (m mockSchedulerError) IsRun() bool {
	return true
}

func (m mockSchedulerStorageError) Get(string) (scheduler.Scheduler, error) {
	return &mockSchedulerError{}, nil
}

func (m mockSchedulerStorageError) Set(scheduler.Scheduler) error {
	return errors.New("e")
}

func (m mockSchedulerStorageError) Remove(string) error {
	return errors.New("e")
}

func (m mockSchedulerStorageError) GetList() map[string]bool {
	panic("implement me")
}

type mockSchedulerStorage struct {
}

type schedulerMock struct {
}

func (s schedulerMock) GetId() string {
	return "1"
}

func (s schedulerMock) Run() error {
	return nil
}

func (s schedulerMock) Stop() error {
	return nil
}

func (s schedulerMock) IsRun() bool {
	return true
}

type mockSchedulerStorageRunned struct {
}

func (m mockSchedulerStorageRunned) Get(string) (scheduler.Scheduler, error) {
	panic("implement me")
}

func (m mockSchedulerStorageRunned) Set(scheduler.Scheduler) error {
	panic("implement me")
}

func (m mockSchedulerStorageRunned) Remove(string) error {
	panic("implement me")
}

func (m mockSchedulerStorageRunned) GetList() map[string]bool {
	return map[string]bool{
		"1": true,
	}
}

func (m mockSchedulerStorage) Get(string) (scheduler.Scheduler, error) {
	return &schedulerMock{}, nil
}

func (m mockSchedulerStorage) Set(scheduler.Scheduler) error {
	return nil
}

func (m mockSchedulerStorage) Remove(string) error {
	return nil
}

func (m mockSchedulerStorage) GetList() map[string]bool {
	return map[string]bool{
		"1": false,
	}
}

type mockHttpTools struct {
}

func (m mockHttpTools) SendRequestTimeoutStatusCode(req *http.Request, timeout time.Duration, expectedCode int, ) (int, []byte, error) {
	panic("implement me")
}

func (m mockHttpTools) SendRequestTimeout(req *http.Request, timeout time.Duration) (int, []byte, error) {
	panic("implement me")
}

func (m mockHttpTools) GetWithRedirectsWithStatusCode(url string, expectedCode int) (int, []byte, error) {
	panic("implement me")
}

func (m mockHttpTools) GetWithRedirects(url string) (int, []byte, error) {
	panic("implement me")
}

func (m mockHttpTools) CreateRequest(method string, url string, headers *map[string]string, log string) *http.Request {
	panic("implement me")
}

func (m mockHttpTools) SendRequest(req *http.Request) (int, []byte, error) {
	panic("implement me")
}

func (m mockHttpTools) SendRequestWithStatusCode(req *http.Request, expectedCode int) (int, []byte, error) {
	panic("implement me")
}

type mockSiteMapStorage struct {
}

func (m mockSiteMapStorage) Get(url string) (*parsers.SiteMap, error) {
	panic("implement me")
}

type mockExternalStorage struct {
}

func (m mockExternalStorage) Write(log job.CheckError) error {
	panic("implement me")
}

func TestNew(t *testing.T) {
	t.Run("Should: create new server", func(t *testing.T) {
		s := New(
			&mockSchedulerStorage{},
			&mockExternalStorage{},
			&mockSiteMapStorage{},
			&mockHttpTools{},
			func(i int) semaphore.Semaphore {
				return semaphore.NewSemaphore(i)
			},
		)
		assert.Implements(t, (*apiPb.SchedulersExecutorServer)(nil), s)
	})
}

func TestServer_Add(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		t.Run("Because default", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 0,
				Config:    nil,
			})
			assert.Equal(t, nil, err)
		})
		t.Run("Because: correct tcp", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 1,
				Config: &apiPb.AddRequest_Tcp{
					Tcp: &apiPb.TcpConfig{
						Host: "wefewf",
						Port: 23,
					}},
			})
			assert.Equal(t, nil, err)
		})
		t.Run("Because: correct http", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 1,
				Config: &apiPb.AddRequest_Http{
					Http: &apiPb.HttpConfig{},
				},
			})
			assert.Equal(t, nil, err)
		})
		t.Run("Because: correct httpJsonValue", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 1,
				Config: &apiPb.AddRequest_HttpValue{
					HttpValue: &apiPb.HttpJsonValueConfig{},
				},
			})
			assert.Equal(t, nil, err)
		})
		t.Run("Because: correct grpc", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 1,
				Config: &apiPb.AddRequest_Grpc{
					Grpc: &apiPb.GrpcConfig{
						Service: "",
						Host:    "wefewf",
						Port:    23,
					}},
			})
			assert.Equal(t, nil, err)
		})
		t.Run("Because: correct sitemap", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 1,
				Config: &apiPb.AddRequest_Sitemap{
					Sitemap: &apiPb.SiteMapConfig{
						Url: "",
					}},
			})
			assert.Equal(t, nil, err)
		})
	})
	t.Run("Should: return error", func(t *testing.T) {
		t.Run("Because: not correct timeout tcp", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 0,
				Config: &apiPb.AddRequest_Tcp{
					Tcp: &apiPb.TcpConfig{
						Host: "wefewf",
						Port: 23,
					}},
			})
			assert.NotEqual(t, nil, err)
		})
		t.Run("Because: not correct timeout http", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 0,
				Config: &apiPb.AddRequest_Http{
					Http: &apiPb.HttpConfig{},
				},
			})
			assert.NotEqual(t, nil, err)
		})
		t.Run("Because: not correct timeout grpc", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 0,
				Config: &apiPb.AddRequest_Grpc{
					Grpc: &apiPb.GrpcConfig{
						Service: "",
						Host:    "wefewf",
						Port:    23,
					}},
			})
			assert.NotEqual(t, nil, err)
		})
		t.Run("Because: not correct sitemap tcp", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 0,
				Config: &apiPb.AddRequest_Sitemap{
					Sitemap: &apiPb.SiteMapConfig{
						Url: "",
					}},
			})
			assert.NotEqual(t, nil, err)
		})
		t.Run("Because: not correct timeout httpJsonValue", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 0,
				Config: &apiPb.AddRequest_HttpValue{
					HttpValue: &apiPb.HttpJsonValueConfig{},
				},
			})
			assert.NotEqual(t, nil, err)
		})
		t.Run("Because: tcp alreadyExistId", func(t *testing.T) {
			s := New(
				&mockSchedulerStorageError{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 1,
				Config: &apiPb.AddRequest_Tcp{
					Tcp: &apiPb.TcpConfig{
						Host: "wefewf",
						Port: 23,
					}},
			})
			assert.NotEqual(t, nil, err)
		})
		t.Run("Because: http alreadyExistId", func(t *testing.T) {
			s := New(
				&mockSchedulerStorageError{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 1,
				Config: &apiPb.AddRequest_Http{
					Http: &apiPb.HttpConfig{
					}},
			})
			assert.NotEqual(t, nil, err)
		})
		t.Run("Because: httpJsonValue alreadyExistId", func(t *testing.T) {
			s := New(
				&mockSchedulerStorageError{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 1,
				Config: &apiPb.AddRequest_HttpValue{
					HttpValue: &apiPb.HttpJsonValueConfig{
					}},
			})
			assert.NotEqual(t, nil, err)
		})
		t.Run("Because: siteMap alreadyExistId", func(t *testing.T) {
			s := New(
				&mockSchedulerStorageError{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 1,
				Config: &apiPb.AddRequest_Sitemap{
					Sitemap: &apiPb.SiteMapConfig{
						Url: "",
					}},
			})
			assert.NotEqual(t, nil, err)
		})
		t.Run("Because: grpc alreadyExistId", func(t *testing.T) {
			s := New(
				&mockSchedulerStorageError{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Add(context.Background(), &apiPb.AddRequest{
				Interval: 1,
				Config: &apiPb.AddRequest_Grpc{
					Grpc: &apiPb.GrpcConfig{
						Port: 8080,
					}},
			})
			assert.NotEqual(t, nil, err)
		})
	})
}

func TestServer_Remove(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		t.Run("Because: correct setting", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Remove(context.Background(), &apiPb.RemoveRequest{
				Id: "",
			})
			assert.Equal(t, nil, err)
		})
	})
	t.Run("Should: return error", func(t *testing.T) {
		t.Run("Because: not exist", func(t *testing.T) {
			s := New(
				&mockSchedulerStorageError{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Remove(context.Background(), &apiPb.RemoveRequest{
				Id: "",
			})
			assert.NotEqual(t, nil, err)
		})
	})
}

func TestServer_Run(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		t.Run("Because: correct setting", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Run(context.Background(), &apiPb.RunRequest{
				Id: "",
			})
			assert.Equal(t, nil, err)
		})
	})
	t.Run("Should: return error", func(t *testing.T) {
		t.Run("Because: cant run", func(t *testing.T) {
			s := New(
				&mockSchedulerStorageError{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Run(context.Background(), &apiPb.RunRequest{
				Id: "",
			})
			assert.NotEqual(t, nil, err)
		})
		t.Run("Because: not exist", func(t *testing.T) {
			s := New(
				&mockSchedulerStorageGetError{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Run(context.Background(), &apiPb.RunRequest{
				Id: "",
			})
			assert.NotEqual(t, nil, err)
		})
	})
}

func TestServer_Stop(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		t.Run("Because: correct setting", func(t *testing.T) {
			s := New(
				&mockSchedulerStorage{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Stop(context.Background(), &apiPb.StopRequest{
				Id: "",
			})
			assert.Equal(t, nil, err)
		})
	})
	t.Run("Should: return error", func(t *testing.T) {
		t.Run("Because: not cant stop", func(t *testing.T) {
			s := New(
				&mockSchedulerStorageError{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Stop(context.Background(), &apiPb.StopRequest{
				Id: "",
			})
			assert.NotEqual(t, nil, err)
		})
		t.Run("Because: not exist", func(t *testing.T) {
			s := New(
				&mockSchedulerStorageGetError{},
				&mockExternalStorage{},
				&mockSiteMapStorage{},
				&mockHttpTools{},
				func(i int) semaphore.Semaphore {
					return semaphore.NewSemaphore(i)
				},
			)
			_, err := s.Stop(context.Background(), &apiPb.StopRequest{
				Id: "",
			})
			assert.NotEqual(t, nil, err)
		})
	})
}

func TestServer_GetSchedulerList(t *testing.T) {

}

func TestServer_GetSchedulerById(t *testing.T) {

}