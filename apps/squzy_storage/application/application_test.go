package application

import (
	"context"
	"errors"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

type dbErrorMock struct {
}

func (*dbErrorMock) Migrate() error {
	return nil
}

func (*dbErrorMock) InsertSnapshot(data *apiPb.SchedulerResponse) error {
	return errors.New("error")
}

func (*dbErrorMock) GetSnapshots(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.SchedulerSnapshot, int32, error) {
	return nil, -1, errors.New("error")
}

func (*dbErrorMock) InsertStatRequest(data *apiPb.Metric) error {
	return errors.New("error")
}

func (*dbErrorMock) GetStatRequest(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return nil, -1, errors.New("error")
}

func (*dbErrorMock) GetCPUInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return nil, -1, errors.New("error")
}

func (*dbErrorMock) GetMemoryInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return nil, -1, errors.New("error")
}

func (*dbErrorMock) GetDiskInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return nil, -1, errors.New("error")
}

func (*dbErrorMock) GetNetInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return nil, -1, errors.New("error")
}

type dbMock struct {
}

func (*dbMock) Migrate() error {
	return nil
}

func (*dbMock) InsertSnapshot(data *apiPb.SchedulerResponse) error {
	return nil
}

func (*dbMock) GetSnapshots(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.SchedulerSnapshot, int32, error) {
	return nil, -1, nil
}

func (*dbMock) InsertStatRequest(data *apiPb.Metric) error {
	return nil
}

func (*dbMock) GetStatRequest(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return nil, -1, nil
}

func (*dbMock) GetCPUInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return nil, -1, nil
}

func (*dbMock) GetMemoryInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return nil, -1, nil
}

func (*dbMock) GetDiskInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return nil, -1, nil
}

func (*dbMock) GetNetInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return nil, -1, nil
}

func TestNewService(t *testing.T) {
	t.Run("Should: return no nil", func(t *testing.T) {
		assert.NotNil(t, NewService(nil))
	})
}

func TestService_SendResponseFromScheduler(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := service{
			database: &dbErrorMock{},
		}
		_, err := s.SendResponseFromScheduler(context.Background(), nil)
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.SendResponseFromScheduler(context.Background(), nil)
		assert.NoError(t, err)
	})
}

func TestService_SendResponseFromAgent(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := service{
			database: &dbErrorMock{},
		}
		_, err := s.SendResponseFromAgent(context.Background(), nil)
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.SendResponseFromAgent(context.Background(), nil)
		assert.NoError(t, err)
	})
}

func TestService_GetSchedulerInformation(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := service{
			database: &dbErrorMock{},
		}
		_, err := s.GetSchedulerInformation(context.Background(), &apiPb.GetSchedulerInformationRequest{})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.GetSchedulerInformation(context.Background(), &apiPb.GetSchedulerInformationRequest{})
		assert.NoError(t, err)
	})
}

func TestService_GetAgentInformation(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := service{
			database: &dbErrorMock{},
		}
		_, err := s.GetAgentInformation(context.Background(), &apiPb.GetAgentInformationRequest{Type: -1})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.GetAgentInformation(context.Background(), &apiPb.GetAgentInformationRequest{Type: apiPb.TypeAgentStat_ALL})
		assert.NoError(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.GetAgentInformation(context.Background(), &apiPb.GetAgentInformationRequest{Type: apiPb.TypeAgentStat_CPU})
		assert.NoError(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.GetAgentInformation(context.Background(), &apiPb.GetAgentInformationRequest{Type: apiPb.TypeAgentStat_MEMORY})
		assert.NoError(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.GetAgentInformation(context.Background(), &apiPb.GetAgentInformationRequest{Type: apiPb.TypeAgentStat_DISK})
		assert.NoError(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.GetAgentInformation(context.Background(), &apiPb.GetAgentInformationRequest{Type: apiPb.TypeAgentStat_NET})
		assert.NoError(t, err)
	})
}
