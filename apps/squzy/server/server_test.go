package server

import (
	"context"
	"errors"
	serverPb "github.com/squzy/squzy_generated/generated/server/proto/v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	"squzy/apps/internal/job"
	"squzy/apps/internal/parsers"
	"squzy/apps/internal/scheduler"
	"squzy/apps/internal/semaphore"
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

func (m mockExternalStorage) Write(id string, log job.CheckError) error {
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
		assert.Implements(t, (*serverPb.ServerServer)(nil), s)
	})
}

func TestServer_AddScheduler(t *testing.T) {
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 0,
				Check:    nil,
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 1,
				Check: &serverPb.AddSchedulerRequest_TcpCheck{
					TcpCheck: &serverPb.TcpCheck{
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 1,
				Check: &serverPb.AddSchedulerRequest_HttpCheck{
					HttpCheck: &serverPb.HttpCheck{},
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 1,
				Check: &serverPb.AddSchedulerRequest_HttpJsonValue{
					HttpJsonValue: &serverPb.HttpJsonValueCheck{},
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 1,
				Check: &serverPb.AddSchedulerRequest_GrpcCheck{
					GrpcCheck: &serverPb.GrpcCheck{
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 1,
				Check: &serverPb.AddSchedulerRequest_SitemapCheck{
					SitemapCheck: &serverPb.SiteMapCheck{
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 0,
				Check: &serverPb.AddSchedulerRequest_TcpCheck{
					TcpCheck: &serverPb.TcpCheck{
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 0,
				Check: &serverPb.AddSchedulerRequest_HttpCheck{
					HttpCheck: &serverPb.HttpCheck{},
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 0,
				Check: &serverPb.AddSchedulerRequest_GrpcCheck{
					GrpcCheck: &serverPb.GrpcCheck{
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 0,
				Check: &serverPb.AddSchedulerRequest_SitemapCheck{
					SitemapCheck: &serverPb.SiteMapCheck{
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 0,
				Check: &serverPb.AddSchedulerRequest_HttpJsonValue{
					HttpJsonValue: &serverPb.HttpJsonValueCheck{},
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 1,
				Check: &serverPb.AddSchedulerRequest_TcpCheck{
					TcpCheck: &serverPb.TcpCheck{
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 1,
				Check: &serverPb.AddSchedulerRequest_HttpCheck{
					HttpCheck: &serverPb.HttpCheck{
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 1,
				Check: &serverPb.AddSchedulerRequest_HttpJsonValue{
					HttpJsonValue: &serverPb.HttpJsonValueCheck{
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 1,
				Check: &serverPb.AddSchedulerRequest_SitemapCheck{
					SitemapCheck: &serverPb.SiteMapCheck{
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
			_, err := s.AddScheduler(context.Background(), &serverPb.AddSchedulerRequest{
				Interval: 1,
				Check: &serverPb.AddSchedulerRequest_GrpcCheck{
					GrpcCheck: &serverPb.GrpcCheck{
						Port: 8080,
					}},
			})
			assert.NotEqual(t, nil, err)
		})
	})
}

func TestServer_RemoveScheduler(t *testing.T) {
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
			_, err := s.RemoveScheduler(context.Background(), &serverPb.RemoveSchedulerRequest{
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
			_, err := s.RemoveScheduler(context.Background(), &serverPb.RemoveSchedulerRequest{
				Id: "",
			})
			assert.NotEqual(t, nil, err)
		})
	})
}

func TestServer_RunScheduler(t *testing.T) {
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
			_, err := s.RunScheduler(context.Background(), &serverPb.RunSchedulerRequest{
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
			_, err := s.RunScheduler(context.Background(), &serverPb.RunSchedulerRequest{
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
			_, err := s.RunScheduler(context.Background(), &serverPb.RunSchedulerRequest{
				Id: "",
			})
			assert.NotEqual(t, nil, err)
		})
	})
}

func TestServer_StopScheduler(t *testing.T) {
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
			_, err := s.StopScheduler(context.Background(), &serverPb.StopSchedulerRequest{
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
			_, err := s.StopScheduler(context.Background(), &serverPb.StopSchedulerRequest{
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
			_, err := s.StopScheduler(context.Background(), &serverPb.StopSchedulerRequest{
				Id: "",
			})
			assert.NotEqual(t, nil, err)
		})
	})
}

func TestServer_GetList(t *testing.T) {
	t.Run("Should: return list with stop", func(t *testing.T) {
		s := New(
			&mockSchedulerStorage{},
			&mockExternalStorage{},
			&mockSiteMapStorage{},
			&mockHttpTools{},
			func(i int) semaphore.Semaphore {
				return semaphore.NewSemaphore(i)
			},
		)
		resp, err := s.GetList(context.Background(), &serverPb.GetListRequest{}, )
		assert.Equal(t, nil, err)
		assert.EqualValues(t, []*serverPb.SchedulerListItem{
			{
				Id:     "1",
				Status: serverPb.Status_STOPPED,
			},
		}, resp.List)
	})
	t.Run("Should: return list with run", func(t *testing.T) {
		s := New(
			&mockSchedulerStorageRunned{},
			&mockExternalStorage{},
			&mockSiteMapStorage{},
			&mockHttpTools{},
			func(i int) semaphore.Semaphore {
				return semaphore.NewSemaphore(i)
			},
		)
		resp, err := s.GetList(context.Background(), &serverPb.GetListRequest{}, )
		assert.Equal(t, nil, err)
		assert.EqualValues(t, []*serverPb.SchedulerListItem{
			{
				Id:     "1",
				Status: serverPb.Status_RUNNED,
			},
		}, resp.List)
	})
}
