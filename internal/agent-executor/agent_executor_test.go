package agent_executor

import (
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mock struct {
}

type mockMoreSecond struct {
}

type mockLessSecond struct {
}

func (m mockLessSecond) GetStat() *apiPb.Metric {
	time.Sleep(time.Millisecond * 500)
	return &apiPb.Metric{}
}

func (m mockMoreSecond) GetStat() *apiPb.Metric {
	time.Sleep(time.Second * 2)
	return &apiPb.Metric{}
}

func (m mock) GetStat() *apiPb.Metric {
	return &apiPb.Metric{}
}

func TestNew(t *testing.T) {
	t.Run("Should: create new executor", func(t *testing.T) {
		a, _ := New(&mock{}, time.Second)
		assert.Implements(t, (*AgentExecutor)(nil), a)
	})
	t.Run("Should: return error if interval less then 500 millisecond", func(t *testing.T) {
		a, err := New(&mock{}, time.Millisecond)
		assert.EqualValues(t, nil, a)
		assert.Equal(t, intervalLessHalfSecondError, err)
	})
}

func TestExecutor_Execute(t *testing.T) {
	t.Run("Should: get value immediately", func(t *testing.T) {
		a, _ := New(&mock{}, time.Second)
		channel := a.Execute()
		value := <-channel
		assert.EqualValues(t, &apiPb.Metric{}, value)
		value2 := <-channel
		assert.EqualValues(t, &apiPb.Metric{}, value2)
	})
	t.Run("Should: not get value, if job execute more that inteval", func(t *testing.T) {
		a, _ := New(&mockMoreSecond{}, time.Second)
		channel := a.Execute()
		select {
		case <-channel:
			assert.FailNow(t, "Job execute less than second")
		case <-time.After(time.Second):
			assert.Equal(t, true, true)
		}
	})
	t.Run("Should: get value, if job execute less than inteval", func(t *testing.T) {
		a, _ := New(&mockLessSecond{}, time.Second)
		channel := a.Execute()
		select {
		case value := <-channel:
			assert.EqualValues(t, &apiPb.Metric{}, value)
		case <-time.After(time.Second):
			assert.FailNow(t, "Job execute less than second")
		}
	})
}
