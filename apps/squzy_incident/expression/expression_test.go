package expression

import (
	"context"
	"errors"
	"github.com/araddon/dateparse"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
)

type mockStorage struct {
}

func (m mockStorage) SaveResponseFromScheduler(ctx context.Context, in *apiPb.SchedulerResponse, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, nil
}

func (m mockStorage) SaveResponseFromAgent(ctx context.Context, in *apiPb.Metric, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, nil
}

func (m mockStorage) SaveTransaction(ctx context.Context, in *apiPb.TransactionInfo, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, nil
}

func (m mockStorage) GetSchedulerInformation(ctx context.Context, in *apiPb.GetSchedulerInformationRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerInformationResponse, error) {
	return &apiPb.GetSchedulerInformationResponse{
		Snapshots: []*apiPb.SchedulerSnapshot{
			{
				Code: apiPb.SchedulerCode_OK,
				Type: apiPb.SchedulerType_GRPC,
				Meta: &apiPb.SchedulerSnapshot_MetaData{
					StartTime: ptypes.TimestampNow(),
					EndTime:   ptypes.TimestampNow(),
				},
			},
		},
	}, nil
}

func (m mockStorage) GetSchedulerUptime(ctx context.Context, in *apiPb.GetSchedulerUptimeRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerUptimeResponse, error) {
	return nil, nil
}

func (m mockStorage) GetAgentInformation(ctx context.Context, in *apiPb.GetAgentInformationRequest, opts ...grpc.CallOption) (*apiPb.GetAgentInformationResponse, error) {
	return &apiPb.GetAgentInformationResponse{
		Stats: []*apiPb.GetAgentInformationResponse_Statistic{
			{
				Time: ptypes.TimestampNow(),
				CpuInfo: &apiPb.CpuInfo{
					Cpus: []*apiPb.CpuInfo_CPU{
						{
							Load: 10,
						},
					},
				},
			},
		},
	}, nil
}

func (m mockStorage) GetTransactionsGroup(ctx context.Context, in *apiPb.GetTransactionGroupRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionGroupResponse, error) {
	return nil, nil
}

func (m mockStorage) GetTransactions(ctx context.Context, in *apiPb.GetTransactionsRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionsResponse, error) {
	return &apiPb.GetTransactionsResponse{
		Transactions: []*apiPb.TransactionInfo{
			{
				Meta: &apiPb.TransactionInfo_Meta{
					Host:   "host",
					Path:   "path",
					Method: "method",
				},
				Name:      "name",
				StartTime: ptypes.TimestampNow(),
				EndTime:   ptypes.TimestampNow(),
				Status:    apiPb.TransactionStatus_TRANSACTION_SUCCESSFUL,
				Type:      apiPb.TransactionType_TRANSACTION_TYPE_DB,
			},
		},
	}, nil
}

func (m mockStorage) GetTransactionById(ctx context.Context, in *apiPb.GetTransactionByIdRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionByIdResponse, error) {
	return nil, nil
}

func (m mockStorage) SaveIncident(ctx context.Context, in *apiPb.Incident, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, nil
}

func (m mockStorage) UpdateIncidentStatus(ctx context.Context, in *apiPb.UpdateIncidentStatusRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	return nil, nil
}

func (m mockStorage) GetIncidentById(ctx context.Context, in *apiPb.IncidentIdRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	return nil, nil
}

func (m mockStorage) GetIncidentByRuleId(ctx context.Context, in *apiPb.RuleIdRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	return nil, nil
}

func (m mockStorage) GetIncidentsList(ctx context.Context, in *apiPb.GetIncidentsListRequest, opts ...grpc.CallOption) (*apiPb.GetIncidentsListResponse, error) {
	return nil, nil
}

type mockErrorStorage struct {
}

func (m mockErrorStorage) SaveResponseFromScheduler(ctx context.Context, in *apiPb.SchedulerResponse, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) SaveResponseFromAgent(ctx context.Context, in *apiPb.Metric, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) SaveTransaction(ctx context.Context, in *apiPb.TransactionInfo, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetSchedulerInformation(ctx context.Context, in *apiPb.GetSchedulerInformationRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerInformationResponse, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetSchedulerUptime(ctx context.Context, in *apiPb.GetSchedulerUptimeRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerUptimeResponse, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetAgentInformation(ctx context.Context, in *apiPb.GetAgentInformationRequest, opts ...grpc.CallOption) (*apiPb.GetAgentInformationResponse, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetTransactionsGroup(ctx context.Context, in *apiPb.GetTransactionGroupRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionGroupResponse, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetTransactions(ctx context.Context, in *apiPb.GetTransactionsRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionsResponse, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetTransactionById(ctx context.Context, in *apiPb.GetTransactionByIdRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionByIdResponse, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) SaveIncident(ctx context.Context, in *apiPb.Incident, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) UpdateIncidentStatus(ctx context.Context, in *apiPb.UpdateIncidentStatusRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetIncidentById(ctx context.Context, in *apiPb.IncidentIdRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetIncidentByRuleId(ctx context.Context, in *apiPb.RuleIdRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetIncidentsList(ctx context.Context, in *apiPb.GetIncidentsListRequest, opts ...grpc.CallOption) (*apiPb.GetIncidentsListResponse, error) {
	return nil, errors.New("ERROR")
}

var (
	storage  = &mockStorage{}
	exprCorr = &expressionStruct{storageClient: storage}
	exprErr  = &expressionStruct{storageClient: &mockErrorStorage{}}
)

func TestNewExpression(t *testing.T) {
	t.Run("Should: return not nil", func(t *testing.T) {
		assert.NotNil(t, NewExpression(storage))
	})
}

func TestExpressionStruct_ProcessRule(t *testing.T) {
	t.Run("Should: panic", func(t *testing.T) {
		panicFunc := func() {
			_ = exprCorr.ProcessRule(
				10,
				"",
				"")
		}
		assert.Equal(t, true, assert.Panics(t, panicFunc, "The code did not panic"))
	})
	t.Run("Should: panic", func(t *testing.T) {
		panicFunc := func() {
			_ = exprCorr.ProcessRule(
				apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT,
				"12345",
				"wrongString")
		}
		assert.Equal(t, true, assert.Panics(t, panicFunc, "The code did not panic"))
	})
	t.Run("Should: panic because storage error", func(t *testing.T) {
		panicFunc := func() {
			_ = exprErr.ProcessRule(
				apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT,
				"12345",
				"one(Last(10), {one(.CpuInfo.Cpus, {.Load <= 10})})")
		}
		assert.Equal(t, true, assert.Panics(t, panicFunc, "The code did not panic"))
	})
	t.Run("Should: no panic", func(t *testing.T) {
		res := exprCorr.ProcessRule(
			apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT,
			"12345",
			"one(Last(10), {one(.CpuInfo.Cpus, {.Load <= 10})})")
		assert.True(t, res)
	})
	t.Run("Should: panic because not bool", func(t *testing.T) {
		panicFunc :=
			func() {
				_ = exprCorr.ProcessRule(
					apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT,
					"12345",
					"filter(Last(10), {one(.CpuInfo.Cpus, {.Load <= 10})})")
			}
		assert.Equal(t, true, assert.Panics(t, panicFunc, "The code did not panic"))
	})
	t.Run("Should: no panic", func(t *testing.T) {
		res := exprCorr.ProcessRule(
			apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT,
			"12345",
			"one(Last(10), {one(.CpuInfo.Cpus, {.Load < 5})})")
		assert.False(t, res)
	})
}

func TestExpressionStruct_IsValid(t *testing.T) {
	t.Run("Should: panic", func(t *testing.T) {
		err := exprCorr.IsValid(
			apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT,
			"wrongString")
		assert.Error(t, err)
	})
	t.Run("Should: no panic", func(t *testing.T) {
		err := exprCorr.IsValid(
			apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT,
			"one(Last(10), {one(.CpuInfo.Cpus, {.Load <= 10})})")
		assert.NoError(t, err)
	})
}

func Test_convertToTimestamp(t *testing.T) {
	t.Run("Should: panic", func(t *testing.T) {
		panicFunc := func() {
			_ = convertToTimestamp("01-10-2020")
		}
		assert.Equal(t, true, assert.Panics(t, panicFunc, "The code did not panic"))
	})
	t.Run("Should: no panic", func(t *testing.T) {
		str := "3/1/2014"
		v, err := dateparse.ParseAny(str)
		assert.Nil(t, err)
		res := convertToTimestamp(str)
		value, err := ptypes.TimestampProto(v)
		assert.Nil(t, err)
		res2, err := ptypes.Timestamp(value)
		assert.Nil(t, err)
		assert.EqualValues(t, res, value)
		assert.EqualValues(t, v, res2)
	})
	t.Run("Should: panic", func(t *testing.T) {
		panicFunc := func() {
			_ = convertToTimestamp("0000-01-01T00:00:00.899Z")
		}
		assert.Equal(t, true, assert.Panics(t, panicFunc, "The code did not panic"))
	})
}

func Test_getTimeRange(t *testing.T) {
	t.Run("Should: panic", func(t *testing.T) {
		panicFunc := func() {
			_ = getTimeRange(nil, nil)
		}
		assert.Equal(t, true, assert.Panics(t, panicFunc, "The code did not panic"))
	})
	t.Run("Should: panic", func(t *testing.T) {
		panicFunc := func() {
			_ = getTimeRange(ptypes.TimestampNow(), nil)
		}
		assert.Equal(t, true, assert.Panics(t, panicFunc, "The code did not panic"))
	})
}
