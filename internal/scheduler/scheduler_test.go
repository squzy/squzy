package scheduler

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

type jobExecutor struct {
	count int
}

func (j *jobExecutor) Execute(schedulerId primitive.ObjectID) {
	j.count += 1
}

func TestNew(t *testing.T) {
	t.Run("Tests: Scheduler.New()", func(t *testing.T) {
		t.Run("Should: create new app without error", func(t *testing.T) {
			_, err := New(primitive.NewObjectID(), time.Second, nil)
			assert.Equal(t, nil, err)
		})
		t.Run("Should: create new app with 'intervalLessHalfSecondError' error", func(t *testing.T) {
			_, err := New(primitive.NewObjectID(), time.Millisecond, nil)
			assert.Equal(t, intervalLessHalfSecondError, err)
		})
	})
}

func TestSchl_Run(t *testing.T) {
	t.Run("Tests: Scheduler.Run()", func(t *testing.T) {
		t.Run("Should: run without error ", func(t *testing.T) {
			i, _ := New(primitive.NewObjectID(), time.Second, &jobExecutor{})
			i.Run()
			i.Run()
			i.Stop()
		})
		t.Run("Should: run job every second ", func(t *testing.T) {
			store := &jobExecutor{}
			i, err := New(primitive.NewObjectID(), time.Second, store)
			assert.Equal(t, nil, err)
			i.Run()
			assert.Equal(t, nil, err)
			ch := make(chan bool)
			time.AfterFunc(time.Millisecond*1100, func() {
				assert.Equal(t, 1, store.count)
				ch <- true
			})
			<-ch
			time.AfterFunc(time.Millisecond*1100, func() {
				assert.Equal(t, 2, store.count)
				ch <- true
			})
			<-ch
			i.Stop()
		})
	})
}

func TestSchl_Stop(t *testing.T) {
	t.Run("Tests: Scheduler.Stop()", func(t *testing.T) {
		t.Run("Should: stop without error ", func(t *testing.T) {
			i, _ := New(primitive.NewObjectID(), time.Second, &jobExecutor{})
			i.Run()
			i.Stop()
			i.Stop()
		})
	})
}

func TestSchl_IsRun(t *testing.T) {
	t.Run("Tests: Scheduler.IsRun()", func(t *testing.T) {
		t.Run("Should: return true ", func(t *testing.T) {
			i, _ := New(primitive.NewObjectID(), time.Second, &jobExecutor{})
			i.Run()
			assert.Equal(t, true, i.IsRun())
			i.Stop()

		})
		t.Run("Should: return false", func(t *testing.T) {
			t.Run("Suite: after creation", func(t *testing.T) {
				i, _ := New(primitive.NewObjectID(), time.Second, &jobExecutor{})
				assert.Equal(t, false, i.IsRun())
			})
			t.Run("Suite: after stop", func(t *testing.T) {
				i, _ := New(primitive.NewObjectID(), time.Second, &jobExecutor{})
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
		s, err := New(id, time.Second, &jobExecutor{})
		assert.Equal(t, id.Hex(), s.GetId())
		assert.IsType(t, "", s.GetId())
		assert.Equal(t, nil, err)
	})
}

func TestSchl_GetIdBson(t *testing.T) {
	t.Run("Should: return id as bson", func(t *testing.T) {
		id := primitive.NewObjectID()
		s, err := New(id, time.Second, &jobExecutor{})
		assert.Equal(t, id, s.GetIdBson())
		assert.IsType(t, primitive.ObjectID{}, s.GetIdBson())
		assert.Equal(t, nil, err)
	})
}
