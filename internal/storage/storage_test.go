package storage

import (
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
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
				StartTime: timestamp.Now(),
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
			Code: apiPb.SchedulerCode_ERROR,
			Error: &apiPb.SchedulerSnapshot_Error{
				Message: "",
			},
			Type: 0,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: timestamp.Now(),
				EndTime:   timestamp.Now(),
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
				StartTime: timestamp.Now(),
				EndTime:   timestamp.Now(),
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
