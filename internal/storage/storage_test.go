package storage

import (
	"github.com/golang/protobuf/ptypes"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockStartTimeErrorMock struct {
}

func (m mockStartTimeErrorMock) GetLogData() *apiPb.SchedulerResponse {
	return &apiPb.SchedulerResponse{
		Snapshot: &apiPb.SchedulerSnapshot{
			Code: 0,
			Error: &apiPb.SchedulerSnapshot_Error{
				Message: "",
			},
			Type: 0,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: nil,
				EndTime:   nil,
			},
		},
	}
}

type mockEndTimeErrorMock struct {
}

func (m mockEndTimeErrorMock) GetLogData() *apiPb.SchedulerResponse {
	return &apiPb.SchedulerResponse{
		Snapshot: &apiPb.SchedulerSnapshot{
			Code: 0,
			Error: &apiPb.SchedulerSnapshot_Error{
				Message: "",
			},
			Type: 0,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: ptypes.TimestampNow(),
				EndTime:   nil,
			},
		},
	}
}

type mockError struct {
}

func (m mockError) GetLogData() *apiPb.SchedulerResponse {
	return &apiPb.SchedulerResponse{
		Snapshot: &apiPb.SchedulerSnapshot{
			Code: apiPb.SchedulerCode_Error,
			Error: &apiPb.SchedulerSnapshot_Error{
				Message: "",
			},
			Type: 0,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: ptypes.TimestampNow(),
				EndTime:   ptypes.TimestampNow(),
			},
		},
	}
}

type mockOk struct {
}

func (m mockOk) GetLogData() *apiPb.SchedulerResponse {
	return &apiPb.SchedulerResponse{
		Snapshot: &apiPb.SchedulerSnapshot{
			Code:  apiPb.SchedulerCode_OK,
			Error: nil,
			Type:  0,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: ptypes.TimestampNow(),
				EndTime:   ptypes.TimestampNow(),
			},
		},
	}
}

func TestGetInMemoryStorage(t *testing.T) {
	t.Run("Memory storage", func(t *testing.T) {
		s := GetInMemoryStorage()
		assert.Implements(t, (*Storage)(nil), s)
	})
}

func TestMemory_Write(t *testing.T) {
	t.Run("Should: throw error because startTime", func(t *testing.T) {
		s := GetInMemoryStorage()
		assert.NotEqual(t, nil, s.Write(&mockStartTimeErrorMock{}))
	})
	t.Run("Should: throw error because endTime", func(t *testing.T) {
		s := GetInMemoryStorage()
		assert.NotEqual(t, nil, s.Write(&mockEndTimeErrorMock{}))
	})
	t.Run("Should: write to error log", func(t *testing.T) {
		s := GetInMemoryStorage()
		assert.Equal(t, nil, s.Write(&mockError{}))
	})
	t.Run("Should: write to info log", func(t *testing.T) {
		s := GetInMemoryStorage()
		assert.Equal(t, nil, s.Write(&mockOk{}))
	})
}
