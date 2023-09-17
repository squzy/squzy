package application

import (
	"context"
	"errors"
	"github.com/squzy/squzy/internal/scheduler"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"testing"
	"time"
)

type mockExecuter struct {
}

func (m mockExecuter) Execute(schedulerId primitive.ObjectID) {
}

type mockStorageError struct {
}

func (m mockStorageError) Get(string) (scheduler.Scheduler, error) {
	panic("implement me")
}

func (m mockStorageError) Set(scheduler.Scheduler) error {
	return errors.New("")
}

func (m mockStorageError) Remove(string) error {
	panic("implement me")
}

type mockConfigStorageOk struct {
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
	panic("implement me")
}

func (m mockConfigStorageError) GetAllForSync(ctx context.Context) ([]*scheduler_config_storage.SchedulerConfig, error) {
	return nil, errors.New("asf")
}

func (m mockConfigStorageOk) Get(ctx context.Context, schedulerId primitive.ObjectID) (*scheduler_config_storage.SchedulerConfig, error) {
	panic("implement me")
}

func (m mockConfigStorageOk) Add(ctx context.Context, config *scheduler_config_storage.SchedulerConfig) error {
	panic("implement me")
}

func (m mockConfigStorageOk) Remove(ctx context.Context, schedulerId primitive.ObjectID) error {
	panic("implement me")
}

func (m mockConfigStorageOk) Run(ctx context.Context, schedulerId primitive.ObjectID) error {
	panic("implement me")
}

func (m mockConfigStorageOk) Stop(ctx context.Context, schedulerId primitive.ObjectID) error {
	panic("implement me")
}

func (m mockConfigStorageOk) GetAll(ctx context.Context) ([]*scheduler_config_storage.SchedulerConfig, error) {
	panic("implement me")
}

func (m mockConfigStorageOk) GetAllForSync(ctx context.Context) ([]*scheduler_config_storage.SchedulerConfig, error) {
	return []*scheduler_config_storage.SchedulerConfig{
		{
			ID:       primitive.ObjectID{},
			Type:     0,
			Status:   apiPb.SchedulerStatus_STOPPED,
			Interval: 1,
			Timeout:  1,
		},
		{
			ID:       primitive.ObjectID{},
			Type:     0,
			Status:   apiPb.SchedulerStatus_RUNNED,
			Interval: 1,
			Timeout:  1,
		},
	}, nil
}

type mockStorageOk struct {
}

func (m mockStorageOk) Get(string) (scheduler.Scheduler, error) {
	panic("implement me")
}

func (m mockStorageOk) Set(scheduler.Scheduler) error {
	return nil
}

func (m mockStorageOk) Remove(string) error {
	panic("implement me")
}

type mockCacheOk struct {
}

func (m mockCacheOk) InsertSchedule(data *apiPb.InsertScheduleWithIdRequest) error {
	return nil
}

func (m mockCacheOk) GetScheduleById(data *apiPb.GetScheduleWithIdRequest) (*apiPb.GetScheduleWithIdResponse, error) {
	return &apiPb.GetScheduleWithIdResponse{
		ScheduledNext: &timestamppb.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
	}, nil
}

func (m mockCacheOk) DeleteScheduleById(data *apiPb.DeleteScheduleWithIdRequest) error {
	return nil
}

type mockCacheErr struct {
}

func (m mockCacheErr) InsertSchedule(data *apiPb.InsertScheduleWithIdRequest) error {
	return errors.New("InsertSchedule")
}

func (m mockCacheErr) GetScheduleById(data *apiPb.GetScheduleWithIdRequest) (*apiPb.GetScheduleWithIdResponse, error) {
	return nil, errors.New("GetScheduleById")
}

func (m mockCacheErr) DeleteScheduleById(data *apiPb.DeleteScheduleWithIdRequest) error {
	return errors.New("DeleteScheduleById")
}

func TestNew(t *testing.T) {
	t.Run("Should: Create new application", func(t *testing.T) {
		app := New(nil, nil, nil, mockCacheOk{})
		assert.NotEqual(t, nil, app)
	})
}

func TestApp_Run(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		app := New(&mockStorageOk{}, &mockExecuter{}, &mockConfigStorageOk{}, mockCacheOk{})
		go func() {
			_ = app.Run(11111)
		}()
		time.Sleep(time.Second)
		_, err := net.Dial("tcp", "localhost:11111")
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return error because port is wrong", func(t *testing.T) {
		app := New(&mockStorageOk{}, &mockExecuter{}, &mockConfigStorageOk{}, mockCacheOk{})
		assert.NotEqual(t, nil, app.Run(1244214))
	})
	t.Run("Should: return err because cant sync with DB", func(t *testing.T) {
		app := New(&mockStorageOk{}, &mockExecuter{}, &mockConfigStorageError{}, mockCacheOk{})
		go func() {
			_ = app.Run(11111)
		}()
		time.Sleep(time.Second)
		_, err := net.Dial("tcp", "localhost:11111")
		assert.Equal(t, nil, err)
	})
}

func TestApp_SyncOne(t *testing.T) {
	t.Run("Should: return error because config wrong", func(t *testing.T) {
		app := New(&mockStorageOk{}, &mockExecuter{}, nil, nil)
		err := app.SyncOne(&scheduler_config_storage.SchedulerConfig{
			ID:       primitive.ObjectID{},
			Type:     0,
			Status:   0,
			Interval: 0,
			Timeout:  0,
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error because cant set in storage", func(t *testing.T) {
		app := New(&mockStorageError{}, &mockExecuter{}, nil, nil)
		err := app.SyncOne(&scheduler_config_storage.SchedulerConfig{
			ID:       primitive.ObjectID{},
			Type:     0,
			Status:   0,
			Interval: 1,
			Timeout:  1,
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return nil because status stopped", func(t *testing.T) {
		app := New(&mockStorageOk{}, &mockExecuter{}, nil, nil)
		err := app.SyncOne(&scheduler_config_storage.SchedulerConfig{
			ID:       primitive.ObjectID{},
			Type:     0,
			Status:   apiPb.SchedulerStatus_STOPPED,
			Interval: 1,
			Timeout:  1,
		})
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return nil because status runned", func(t *testing.T) {
		app := New(&mockStorageOk{}, &mockExecuter{}, nil, mockCacheOk{})
		err := app.SyncOne(&scheduler_config_storage.SchedulerConfig{
			ID:       primitive.ObjectID{},
			Type:     0,
			Status:   apiPb.SchedulerStatus_RUNNED,
			Interval: 1,
			Timeout:  1,
		})
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return err because cache returns error", func(t *testing.T) {
		app := New(&mockStorageOk{}, &mockExecuter{}, nil, mockCacheErr{})
		err := app.SyncOne(&scheduler_config_storage.SchedulerConfig{
			ID:       primitive.ObjectID{},
			Type:     0,
			Status:   apiPb.SchedulerStatus_RUNNED,
			Interval: 1,
			Timeout:  1,
		})
		assert.ErrorContains(t, err, "InsertSchedule")
	})
}
