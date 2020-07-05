package server

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
)

type mockClient struct {

}

func (m mockClient) CreateRule(ctx context.Context, in *apiPb.CreateRuleRequest, opts ...grpc.CallOption) (*apiPb.Rule, error) {
	panic("implement me")
}

func (m mockClient) GetRuleById(ctx context.Context, in *apiPb.RuleIdRequest, opts ...grpc.CallOption) (*apiPb.Rule, error) {
	panic("implement me")
}

func (m mockClient) GetRulesByOwnerId(ctx context.Context, in *apiPb.GetRulesByOwnerIdRequest, opts ...grpc.CallOption) (*apiPb.Rules, error) {
	panic("implement me")
}

func (m mockClient) RemoveRule(ctx context.Context, in *apiPb.RuleIdRequest, opts ...grpc.CallOption) (*apiPb.Rule, error) {
	panic("implement me")
}

func (m mockClient) ValidateRule(ctx context.Context, in *apiPb.ValidateRuleRequest, opts ...grpc.CallOption) (*apiPb.ValidateRuleResponse, error) {
	panic("implement me")
}

func (m mockClient) ProcessRecordFromStorage(ctx context.Context, in *apiPb.StorageRecord, opts ...grpc.CallOption) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (m mockClient) CloseIncident(ctx context.Context, in *apiPb.IncidentIdRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	panic("implement me")
}

func (m mockClient) ActivateRule(ctx context.Context, in *apiPb.RuleIdRequest, opts ...grpc.CallOption) (*apiPb.Rule, error) {
	panic("implement me")
}

func (m mockClient) DeactivateRule(ctx context.Context, in *apiPb.RuleIdRequest, opts ...grpc.CallOption) (*apiPb.Rule, error) {
	panic("implement me")
}

func (m mockClient) StudyIncident(ctx context.Context, in *apiPb.IncidentIdRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	panic("implement me")
}

type mockConfigDisable struct {

}
func (m mockConfigDisable) GetPort() int32 {
	panic("implement me")
}

func (m mockConfigDisable) GetDbHost() string {
	panic("implement me")
}

func (m mockConfigDisable) GetDbPort() string {
	panic("implement me")
}

func (m mockConfigDisable) GetDbName() string {
	panic("implement me")
}

func (m mockConfigDisable) GetDbUser() string {
	panic("implement me")
}

func (m mockConfigDisable) GetDbPassword() string {
	panic("implement me")
}

func (m mockConfigDisable) GetIncidentServerAddress() string {
	panic("implement me")
}

func (m mockConfigDisable) WithIncident() bool {
	return false
}

type mockConfigEnable struct {

}

func (m mockConfigEnable) GetPort() int32 {
	panic("implement me")
}

func (m mockConfigEnable) GetDbHost() string {
	panic("implement me")
}

func (m mockConfigEnable) GetDbPort() string {
	panic("implement me")
}

func (m mockConfigEnable) GetDbName() string {
	panic("implement me")
}

func (m mockConfigEnable) GetDbUser() string {
	panic("implement me")
}

func (m mockConfigEnable) GetDbPassword() string {
	panic("implement me")
}

func (m mockConfigEnable) GetIncidentServerAddress() string {
	panic("implement me")
}

func (m mockConfigEnable) WithIncident() bool {
	return true
}

type dbErrorMock struct {
}

func (mock *dbErrorMock) InsertIncident(*apiPb.Incident) error {
	return errors.New("ERROR")
}

func (mock *dbErrorMock) GetIncidentById(id string) (*apiPb.Incident, error) {
	return nil, errors.New("ERROR")
}

func (mock *dbErrorMock) GetActiveIncidentByRuleId(ruleId string) (*apiPb.Incident, error) {
	return nil, errors.New("ERROR")
}

func (mock *dbErrorMock) UpdateIncidentStatus(id string, status apiPb.IncidentStatus) (*apiPb.Incident, error) {
	return nil, errors.New("ERROR")
}

func (mock *dbErrorMock) GetIncidents(request *apiPb.GetIncidentsListRequest) ([]*apiPb.Incident, int64, error) {
	return nil, 0, errors.New("ERROR")
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

func (mock *dbMock) InsertIncident(*apiPb.Incident) error {
	return nil
}

func (mock *dbMock) GetIncidentById(id string) (*apiPb.Incident, error) {
	return &apiPb.Incident{}, nil
}

func (mock *dbMock) GetActiveIncidentByRuleId(ruleId string) (*apiPb.Incident, error) {
	return &apiPb.Incident{}, nil
}

func (mock *dbMock) UpdateIncidentStatus(id string, status apiPb.IncidentStatus) (*apiPb.Incident, error) {
	return &apiPb.Incident{}, nil
}

func (mock *dbMock) GetIncidents(request *apiPb.GetIncidentsListRequest) ([]*apiPb.Incident, int64, error) {
	return nil, 0, nil
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
		assert.NotNil(t, NewServer(nil, nil, nil))
	})
}

func TestService_SaveResponseFromScheduler(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
			cfg: mockConfigEnable{},
		}
		_, err := s.SaveResponseFromScheduler(context.Background(), nil)
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
			cfg: mockConfigEnable{},
		}
		_, err := s.SaveResponseFromScheduler(context.Background(), nil)
		assert.NoError(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
			cfg: mockConfigEnable{},
		}
		_, err := s.SaveResponseFromScheduler(context.Background(), &apiPb.SchedulerResponse{
			Snapshot: &apiPb.SchedulerSnapshot{},
		})
		assert.NoError(t, err)
	})
}

func TestService_SaveResponseFromAgent(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
			cfg: mockConfigDisable{},
		}
		_, err := s.SaveResponseFromAgent(context.Background(), nil)
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
			cfg: mockConfigDisable{},
		}
		_, err := s.SaveResponseFromAgent(context.Background(), nil)
		assert.NoError(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
			cfg: mockConfigEnable{},
		}
		_, err := s.SaveResponseFromAgent(context.Background(), &apiPb.Metric{})
		assert.NoError(t, err)
	})
}

func TestService_GetSchedulerInformation(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
		}
		_, err := s.GetSchedulerInformation(context.Background(), &apiPb.GetSchedulerInformationRequest{})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
		}
		_, err := s.GetSchedulerInformation(context.Background(), &apiPb.GetSchedulerInformationRequest{})
		assert.NoError(t, err)
	})
}

func TestService_GetSchedulerUptime(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
		}
		_, err := s.GetSchedulerUptime(context.Background(), &apiPb.GetSchedulerUptimeRequest{})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
		}
		_, err := s.GetSchedulerUptime(context.Background(), &apiPb.GetSchedulerUptimeRequest{})
		assert.NoError(t, err)
	})
}

func TestService_GetAgentInformation(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
		}
		_, err := s.GetAgentInformation(context.Background(), &apiPb.GetAgentInformationRequest{Type: -1})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
		}
		_, err := s.GetAgentInformation(context.Background(), &apiPb.GetAgentInformationRequest{Type: apiPb.TypeAgentStat_ALL})
		assert.NoError(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
		}
		_, err := s.GetAgentInformation(context.Background(), &apiPb.GetAgentInformationRequest{Type: apiPb.TypeAgentStat_CPU})
		assert.NoError(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
		}
		_, err := s.GetAgentInformation(context.Background(), &apiPb.GetAgentInformationRequest{Type: apiPb.TypeAgentStat_MEMORY})
		assert.NoError(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
		}
		_, err := s.GetAgentInformation(context.Background(), &apiPb.GetAgentInformationRequest{Type: apiPb.TypeAgentStat_DISK})
		assert.NoError(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
		}
		_, err := s.GetAgentInformation(context.Background(), &apiPb.GetAgentInformationRequest{Type: apiPb.TypeAgentStat_NET})
		assert.NoError(t, err)
	})
}

func TestService_SaveTransaction(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
			cfg: mockConfigDisable{},
		}
		_, err := s.SaveTransaction(context.Background(), &apiPb.TransactionInfo{})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
			cfg: mockConfigDisable{},
		}
		_, err := s.SaveTransaction(context.Background(), &apiPb.TransactionInfo{})
		assert.NoError(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
			cfg: mockConfigEnable{},
		}
		_, err := s.SaveTransaction(context.Background(), &apiPb.TransactionInfo{})
		assert.NoError(t, err)
	})
}

func TestService_GetTransactions(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
		}
		_, err := s.GetTransactions(context.Background(), &apiPb.GetTransactionsRequest{})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
		}
		_, err := s.GetTransactions(context.Background(), &apiPb.GetTransactionsRequest{})
		assert.NoError(t, err)
	})
}

func TestService_GetTransactionById(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
		}
		_, err := s.GetTransactionById(context.Background(), &apiPb.GetTransactionByIdRequest{})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
		}
		_, err := s.GetTransactionById(context.Background(), &apiPb.GetTransactionByIdRequest{})
		assert.NoError(t, err)
	})
}

func TestService_GetTransactionsGroup(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
		}
		_, err := s.GetTransactionsGroup(context.Background(), &apiPb.GetTransactionGroupRequest{})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
		}
		_, err := s.GetTransactionsGroup(context.Background(), &apiPb.GetTransactionGroupRequest{})
		assert.NoError(t, err)
	})
}

func TestServer_SaveIncident(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
		}
		_, err := s.SaveIncident(context.Background(), nil)
		assert.Error(t, err)
	})
}

func TestServer_UpdateIncidentStatus(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
		}
		_, err := s.UpdateIncidentStatus(context.Background(), nil)
		assert.Error(t, err)
	})
}

func TestServer_GetIncidentById(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
		}
		_, err := s.GetIncidentById(context.Background(), nil)
		assert.Error(t, err)
	})
}

func TestServer_GetIncidentByRuleId(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
		}
		_, err := s.GetIncidentByRuleId(context.Background(), nil)
		assert.Error(t, err)
	})
}

func TestServer_GetIncidentsList(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := server{
			database: &dbErrorMock{},
		}
		_, err := s.GetIncidentsList(context.Background(), nil)
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		s := server{
			database: &dbMock{},
		}
		_, err := s.GetIncidentsList(context.Background(), nil)
		assert.NoError(t, err)
	})
}

func TestServer_SendRecordToIncident(t *testing.T) {
	t.Run("Should: not throw panic and send incidnet", func(t *testing.T) {
		s := &server{
			database: nil,
			incidentClient: &mockClient{},
			cfg: &mockConfigEnable{},
		}
		panicFn := func() {
			s.SendRecordToIncident(nil)
		}
		assert.NotPanics(t, panicFn)
	})
	t.Run("Should: not send incident", func(t *testing.T) {
		s := &server{
			database: nil,
			cfg: &mockConfigDisable{},
		}
		panicFn := func() {
			s.SendRecordToIncident(nil)
		}
		assert.NotPanics(t, panicFn)
	})
	t.Run("Should: not send incident", func(t *testing.T) {
		s := &server{
			database: nil,
			cfg: &mockConfigEnable{},
		}
		panicFn := func() {
			s.SendRecordToIncident(nil)
		}
		assert.NotPanics(t, panicFn)
	})

}
