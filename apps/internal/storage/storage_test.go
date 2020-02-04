package storage

import (
	"github.com/golang/protobuf/ptypes"
	storagePb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockStartTimeErrorMock struct {
	
}

func (m mockStartTimeErrorMock) GetLogData() *storagePb.Log {
	return &storagePb.Log{
		Code:                 0,
		Description:          "",
		Meta:                 &storagePb.MetaData{
			Id:                   "",
			Location:             "",
			Port:                 0,
			StartTime:            nil,
			EndTime:              nil,
			Type:                 0,
		},
	}
}

type mockEndTimeErrorMock struct {

}

func (m mockEndTimeErrorMock) GetLogData() *storagePb.Log {
	return &storagePb.Log{
		Code:                 0,
		Description:          "",
		Meta:                 &storagePb.MetaData{
			Id:                   "",
			Location:             "",
			Port:                 0,
			StartTime:            ptypes.TimestampNow(),
			EndTime:              nil,
			Type:                 0,
		},
	}
}

type mockError struct {
	
}

func (m mockError) GetLogData() *storagePb.Log {
	return &storagePb.Log{
		Code:                 storagePb.StatusCode_Error,
		Description:          "",
		Meta:                 &storagePb.MetaData{
			Id:                   "",
			Location:             "",
			Port:                 0,
			StartTime:            ptypes.TimestampNow(),
			EndTime:              ptypes.TimestampNow(),
			Type:                 0,
		},
	}
}

type mockOk struct {
	
}

func (m mockOk) GetLogData() *storagePb.Log {
	return &storagePb.Log{
		Code:                 storagePb.StatusCode_OK,
		Description:          "",
		Meta:                 &storagePb.MetaData{
			Id:                   "",
			Location:             "",
			Port:                 0,
			StartTime:            ptypes.TimestampNow(),
			EndTime:              ptypes.TimestampNow(),
			Type:                 0,
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
		assert.NotEqual(t, nil, s.Write("", &mockStartTimeErrorMock{}))
	})
	t.Run("Should: throw error because endTime", func(t *testing.T) {
		s := GetInMemoryStorage()
		assert.NotEqual(t, nil, s.Write("", &mockEndTimeErrorMock{}))
	})
	t.Run("Should: write to error log", func(t *testing.T) {
		s := GetInMemoryStorage()
		assert.Equal(t, nil, s.Write("", &mockError{}))
	})
	t.Run("Should: write to info log", func(t *testing.T) {
		s := GetInMemoryStorage()
		assert.Equal(t, nil, s.Write("", &mockOk{}))
	})
}