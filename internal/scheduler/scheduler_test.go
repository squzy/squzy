package scheduler

import (
	"errors"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

type cacheMock struct {
	startTime     time.Time
	scheduledNext time.Duration
}

func (c cacheMock) InsertSchedule(data *apiPb.InsertScheduleWithIdRequest) error {
	return nil
}

func (c cacheMock) GetScheduleById(data *apiPb.GetScheduleWithIdRequest) (*apiPb.GetScheduleWithIdResponse, error) {
	t := c.startTime.Add(c.scheduledNext)
	return &apiPb.GetScheduleWithIdResponse{
		ScheduledNext: timestamppb.New(t),
	}, nil
}

func (c cacheMock) DeleteScheduleById(data *apiPb.DeleteScheduleWithIdRequest) error {
	return nil
}

type cacheMockErr struct {
}

func (c cacheMockErr) InsertSchedule(data *apiPb.InsertScheduleWithIdRequest) error {
	return errors.New("InsertSchedule")
}

func (c cacheMockErr) GetScheduleById(data *apiPb.GetScheduleWithIdRequest) (*apiPb.GetScheduleWithIdResponse, error) {
	return nil, errors.New("GetScheduleById")
}

func (c cacheMockErr) DeleteScheduleById(data *apiPb.DeleteScheduleWithIdRequest) error {
	return errors.New("DeleteScheduleById")
}

type cacheMockErrGet struct {
}

func (c cacheMockErrGet) InsertSchedule(data *apiPb.InsertScheduleWithIdRequest) error {
	return nil
}

func (c cacheMockErrGet) GetScheduleById(data *apiPb.GetScheduleWithIdRequest) (*apiPb.GetScheduleWithIdResponse, error) {
	return nil, errors.New("GetScheduleById")
}

func (c cacheMockErrGet) DeleteScheduleById(data *apiPb.DeleteScheduleWithIdRequest) error {
	return nil
}

type jobExecutor struct {
	count int
}

func (j *jobExecutor) Execute(schedulerId primitive.ObjectID) {
	j.count += 1
}

func TestNew(t *testing.T) {
	t.Run("Tests: Scheduler.New()", func(t *testing.T) {
		t.Run("Should: create new app without error", func(t *testing.T) {
			_, err := New(primitive.NewObjectID(), time.Second, nil, &cacheMock{})
			assert.Equal(t, nil, err)
		})
		t.Run("Should: create new app with 'intervalLessHalfSecondError' error", func(t *testing.T) {
			_, err := New(primitive.NewObjectID(), time.Millisecond, nil, &cacheMock{})
			assert.Equal(t, errIntervalLessHalfSecondError, err)
		})
	})
}

func TestSchl_Run(t *testing.T) {
	t.Run("Tests: Scheduler.Run()", func(t *testing.T) {
		t.Run("Should: run without error ", func(t *testing.T) {
			i, _ := New(primitive.NewObjectID(), time.Second, &jobExecutor{}, &cacheMock{})
			i.Run()
			i.Run()
			i.Stop()
		})
		t.Run("Should: run job every second ", func(t *testing.T) {
			store := &jobExecutor{}
			i, err := New(primitive.NewObjectID(), time.Second, store, &cacheMock{
				time.Now(), 900 * time.Millisecond,
			})
			assert.Equal(t, nil, err)
			err = i.Run()
			assert.Equal(t, nil, err)
			ch := make(chan bool)
			time.AfterFunc(time.Millisecond*1100, func() {
				assert.Equal(t, 1, store.count)
				ch <- true
			})
			<-ch
			time.AfterFunc(time.Millisecond*1100, func() {
				assert.Equal(t, 3, store.count)
				ch <- true
			})
			<-ch
			i.Stop()
		})
		t.Run("Should: observer should receive quit ch ", func(t *testing.T) {
			store := &jobExecutor{}
			i, err := New(primitive.NewObjectID(), time.Second, store, &cacheMock{
				time.Now(), 900 * time.Millisecond,
			})
			assert.Equal(t, nil, err)
			err = i.Run()
			assert.Equal(t, nil, err)
			i.Stop()
		})
		t.Run("Should: return err ", func(t *testing.T) {
			store := &jobExecutor{}
			i, err := New(primitive.NewObjectID(), time.Second, store, &cacheMockErr{})
			assert.Equal(t, nil, err)
			err = i.Run()
			assert.ErrorContains(t, err, "InsertSchedule")
		})
		t.Run("Should: observer return err ", func(t *testing.T) {
			store := &jobExecutor{}
			i, err := New(primitive.NewObjectID(), time.Second, store, &cacheMockErrGet{})
			assert.Equal(t, nil, err)
			err = i.Run()
			assert.Nil(t, err)
		})
	})
}

func TestSchl_Stop(t *testing.T) {
	t.Run("Tests: Scheduler.Stop()", func(t *testing.T) {
		t.Run("Should: stop without error ", func(t *testing.T) {
			i, _ := New(primitive.NewObjectID(), time.Second, &jobExecutor{}, &cacheMock{})
			i.Run()
			i.Stop()
			i.Stop()
		})
	})
}

func TestSchl_IsRun(t *testing.T) {
	t.Run("Tests: Scheduler.IsRun()", func(t *testing.T) {
		t.Run("Should: return true ", func(t *testing.T) {
			i, _ := New(primitive.NewObjectID(), time.Second, &jobExecutor{}, &cacheMock{})
			i.Run()
			assert.Equal(t, true, i.IsRun())
			i.Stop()

		})
		t.Run("Should: return false", func(t *testing.T) {
			t.Run("Suite: after creation", func(t *testing.T) {
				i, _ := New(primitive.NewObjectID(), time.Second, &jobExecutor{}, &cacheMock{})
				assert.Equal(t, false, i.IsRun())
			})
			t.Run("Suite: after stop", func(t *testing.T) {
				i, _ := New(primitive.NewObjectID(), time.Second, &jobExecutor{}, &cacheMock{})
				i.Run()
				i.Stop()
				assert.Equal(t, false, i.IsRun())
			})
		})
	})
}

func TestSchl_GetId(t *testing.T) {
	t.Run("Should: return id as string", func(t *testing.T) {
		id := primitive.NewObjectID()
		s, err := New(id, time.Second, &jobExecutor{}, &cacheMock{})
		assert.Equal(t, id.Hex(), s.GetID())
		assert.IsType(t, "", s.GetID())
		assert.Equal(t, nil, err)
	})
}

func TestSchl_GetIdBson(t *testing.T) {
	t.Run("Should: return id as bson", func(t *testing.T) {
		id := primitive.NewObjectID()
		s, err := New(id, time.Second, &jobExecutor{}, &cacheMock{})
		assert.Equal(t, id, s.GetIDBson())
		assert.IsType(t, primitive.ObjectID{}, s.GetIDBson())
		assert.Equal(t, nil, err)
	})
}
