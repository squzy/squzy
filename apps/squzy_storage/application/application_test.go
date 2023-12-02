package application

import (
	"context"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"net"
	"testing"
	"time"
)

type configErrorMock struct {
}

func (*configErrorMock) GetPort() int32 {
	return 1000000
}

func (*configErrorMock) GetDbHost() string {
	panic("implement me!")
}

func (*configErrorMock) GetDbPort() string {
	panic("implement me!")
}

func (*configErrorMock) GetDbName() string {
	panic("implement me!")
}

func (*configErrorMock) GetDbUser() string {
	panic("implement me!")
}

func (*configErrorMock) GetDbPassword() string {
	panic("implement me!")
}

func (*configErrorMock) GetIncidentServerAddress() string {
	return ""
}

func (*configErrorMock) WithIncident() bool {
	return false
}

func (*configErrorMock) WithDbLogs() bool {
	return false
}

type configMock struct {
}

func (*configMock) GetPort() int32 {
	return 23233
}

func (*configMock) GetDbHost() string {
	panic("implement me!")
}

func (*configMock) GetIncidentServerAddress() string {
	return ""
}

func (*configMock) WithIncident() bool {
	return false
}

func (*configMock) WithDbLogs() bool {
	return false
}

func (*configMock) GetDbPort() string {
	panic("implement me!")
}

func (*configMock) GetDbName() string {
	panic("implement me!")
}

func (*configMock) GetDbUser() string {
	panic("implement me!")
}

func (*configMock) GetDbPassword() string {
	panic("implement me!")
}

type mockApiStorage struct {
	apiPb.UnimplementedStorageServer
}

func (s mockApiStorage) SaveIncident(context.Context, *apiPb.Incident) (*empty.Empty, error) {
	panic("implement me")
}

func (s mockApiStorage) UpdateIncidentStatus(context.Context, *apiPb.UpdateIncidentStatusRequest) (*apiPb.Incident, error) {
	panic("implement me")
}

func (s mockApiStorage) GetIncidentById(context.Context, *apiPb.IncidentIdRequest) (*apiPb.Incident, error) {
	panic("implement me")
}

func (s mockApiStorage) GetIncidentByRuleId(context.Context, *apiPb.RuleIdRequest) (*apiPb.Incident, error) {
	panic("implement me")
}

func (s mockApiStorage) GetIncidentsList(context.Context, *apiPb.GetIncidentsListRequest) (*apiPb.GetIncidentsListResponse, error) {
	panic("implement me")
}

func (s mockApiStorage) GetSchedulerUptime(ctx context.Context, in *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	panic("implement me")
}

func (m mockApiStorage) SaveResponseFromScheduler(ctx context.Context, response *apiPb.SchedulerResponse) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockApiStorage) SaveResponseFromAgent(ctx context.Context, metric *apiPb.Metric) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockApiStorage) SaveTransaction(ctx context.Context, info *apiPb.TransactionInfo) (*empty.Empty, error) {
	panic("implement me")
}

func (m mockApiStorage) GetSchedulerInformation(ctx context.Context, request *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error) {
	panic("implement me")
}

func (m mockApiStorage) GetAgentInformation(ctx context.Context, request *apiPb.GetAgentInformationRequest) (*apiPb.GetAgentInformationResponse, error) {
	panic("implement me")
}

func (m mockApiStorage) GetTransactionsGroup(ctx context.Context, request *apiPb.GetTransactionGroupRequest) (*apiPb.GetTransactionGroupResponse, error) {
	panic("implement me")
}

func (m mockApiStorage) GetTransactions(ctx context.Context, request *apiPb.GetTransactionsRequest) (*apiPb.GetTransactionsResponse, error) {
	panic("implement me")
}

func (m mockApiStorage) GetTransactionById(ctx context.Context, request *apiPb.GetTransactionByIdRequest) (*apiPb.GetTransactionByIdResponse, error) {
	panic("implement me")
}

func TestNewServer(t *testing.T) {
	t.Run("Should: work", func(t *testing.T) {
		s := NewApplication(nil, nil)
		assert.NotNil(t, s)
	})
}

func TestServer_Run(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := &application{
			config:  &configErrorMock{},
			apiServ: nil,
		}
		assert.Error(t, s.Run())
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := &application{
			config:  &configMock{},
			apiServ: &mockApiStorage{},
		}
		go func() {
			_ = s.Run()
		}()
		time.Sleep(time.Second)
		_, err := net.Dial("tcp", "localhost:23233")
		assert.Equal(t, nil, err)
	})
}
