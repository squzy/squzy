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

func (*dbErrorMock) GetSnapshots(*apiPb.GetSchedulerInformationRequest) ([]*apiPb.SchedulerSnapshot, int32, error) {
	return nil, -1, errors.New("error")
}

func (*dbErrorMock) GetSnapshotsUptime(request *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	return nil, errors.New("error")
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

func (mock *dbErrorMock) GetNetInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return nil, -1, errors.New("error")
}

func (*dbErrorMock) InsertTransactionInfo(data *apiPb.TransactionInfo) error {
	return errors.New("error")
}

func (*dbErrorMock) GetTransactionInfo(request *apiPb.GetTransactionsRequest) ([]*apiPb.TransactionInfo, int64, error) {
	return nil, -1, errors.New("error")
}

func (*dbErrorMock) GetTransactionByID(request *apiPb.GetTransactionByIdRequest) (*apiPb.TransactionInfo, []*apiPb.TransactionInfo, error) {
	return nil, nil, errors.New("error")
}

func (*dbErrorMock) GetTransactionGroup(request *apiPb.GetTransactionGroupRequest) (map[string]*apiPb.TransactionGroup, error) {
	return nil, errors.New("error")
}

type dbMock struct {
}

func (*dbMock) Migrate() error {
	return nil
}

func (*dbMock) InsertSnapshot(data *apiPb.SchedulerResponse) error {
	return nil
}

func (*dbMock) GetSnapshots(*apiPb.GetSchedulerInformationRequest) ([]*apiPb.SchedulerSnapshot, int32, error) {
	return nil, -1, nil
}

func (*dbMock) GetSnapshotsUptime(request *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	return nil, nil
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

func (*dbMock) InsertTransactionInfo(data *apiPb.TransactionInfo) error {
	return nil
}

func (*dbMock) GetTransactionInfo(request *apiPb.GetTransactionsRequest) ([]*apiPb.TransactionInfo, int64, error) {
	return nil, -1, nil
}

func (*dbMock) GetTransactionByID(request *apiPb.GetTransactionByIdRequest) (*apiPb.TransactionInfo, []*apiPb.TransactionInfo, error) {
	return nil, nil, nil
}

func (*dbMock) GetTransactionGroup(request *apiPb.GetTransactionGroupRequest) (map[string]*apiPb.TransactionGroup, error) {
	return nil, nil
}

func TestNewService(t *testing.T) {
	t.Run("Should: return no nil", func(t *testing.T) {
		assert.NotNil(t, NewService(nil))
	})
}

func TestService_SaveResponseFromScheduler(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := service{
			database: &dbErrorMock{},
		}
		_, err := s.SaveResponseFromScheduler(context.Background(), nil)
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.SaveResponseFromScheduler(context.Background(), nil)
		assert.NoError(t, err)
	})
}

func TestService_SaveResponseFromAgent(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := service{
			database: &dbErrorMock{},
		}
		_, err := s.SaveResponseFromAgent(context.Background(), nil)
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.SaveResponseFromAgent(context.Background(), nil)
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

func TestService_GetSchedulerUptime(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := service{
			database: &dbErrorMock{},
		}
		_, err := s.GetSchedulerUptime(context.Background(), &apiPb.GetSchedulerUptimeRequest{})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.GetSchedulerUptime(context.Background(), &apiPb.GetSchedulerUptimeRequest{})
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

func TestService_SaveTransaction(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := service{
			database: &dbErrorMock{},
		}
		_, err := s.SaveTransaction(context.Background(), &apiPb.TransactionInfo{})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.SaveTransaction(context.Background(), &apiPb.TransactionInfo{})
		assert.NoError(t, err)
	})
}

func TestService_GetTransactions(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := service{
			database: &dbErrorMock{},
		}
		_, err := s.GetTransactions(context.Background(), &apiPb.GetTransactionsRequest{})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.GetTransactions(context.Background(), &apiPb.GetTransactionsRequest{})
		assert.NoError(t, err)
	})
}

func TestService_GetTransactionById(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := service{
			database: &dbErrorMock{},
		}
		_, err := s.GetTransactionById(context.Background(), &apiPb.GetTransactionByIdRequest{})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.GetTransactionById(context.Background(), &apiPb.GetTransactionByIdRequest{})
		assert.NoError(t, err)
	})
}

func TestService_GetTransactionsGroup(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := service{
			database: &dbErrorMock{},
		}
		_, err := s.GetTransactionsGroup(context.Background(), &apiPb.GetTransactionGroupRequest{})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := service{
			database: &dbMock{},
		}
		_, err := s.GetTransactionsGroup(context.Background(), &apiPb.GetTransactionGroupRequest{})
		assert.NoError(t, err)
	})
}
