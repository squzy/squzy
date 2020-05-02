package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/shirou/gopsutil/host"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"net"
	"squzy/internal/grpcTools"
	"syscall"
	"testing"
	"time"
)

type executorMock struct {
	ch chan *apiPb.SendMetricsRequest
}

func (e *executorMock) Execute() chan *apiPb.SendMetricsRequest {
	return e.ch
}

type serverError struct {
}

func (s serverError) GetByAgentUniqName(context.Context, *apiPb.GetByAgentUniqNameRequest) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (s serverError) GetAgentList(context.Context, *empty.Empty) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (s serverError) SendMetrics(apiPb.AgentServer_SendMetricsServer) error {
	panic("implement me")
}

func (s serverError) Register(context.Context, *apiPb.RegisterRequest) (*apiPb.RegisterResponse, error) {
	return nil, errors.New("traarar")
}

func (s serverError) UnRegister(context.Context, *apiPb.UnRegisterRequest) (*apiPb.UnRegisterResponse, error) {
	return &apiPb.UnRegisterResponse{}, nil
}

type configSecondMock struct {
}

func (c configSecondMock) GetAgentName() string {
	return ""
}

type serverSteamError struct {
}

func (s serverSteamError) GetByAgentUniqName(context.Context, *apiPb.GetByAgentUniqNameRequest) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (s serverSteamError) GetAgentList(context.Context, *empty.Empty) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (s serverSteamError) SendMetrics(apiPb.AgentServer_SendMetricsServer) error {
	panic("implement me")
}

func (s serverSteamError) Register(context.Context, *apiPb.RegisterRequest) (*apiPb.RegisterResponse, error) {
	return &apiPb.RegisterResponse{
		Id: "",
	}, nil
}

func (s serverSteamError) UnRegister(context.Context, *apiPb.UnRegisterRequest) (*apiPb.UnRegisterResponse, error) {
	panic("implement me")
}

func (s serverSteamError) SendStat(req apiPb.AgentServer_SendMetricsClient) error {
	return nil
}

func (c configSecondMock) GetExecutionTimeout() time.Duration {
	return time.Second * 5
}

type serverSuccess struct {
	ch    chan *apiPb.SendMetricsRequest
	count int
}

func (s *serverSuccess) GetByAgentUniqName(context.Context, *apiPb.GetByAgentUniqNameRequest) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (s *serverSuccess) GetAgentList(context.Context, *empty.Empty) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (s *serverSuccess) SendMetrics(rq apiPb.AgentServer_SendMetricsServer) error {
	res, _ := rq.Recv()
	s.ch <- res
	return nil
}

func (s *serverSuccess) Register(context.Context, *apiPb.RegisterRequest) (*apiPb.RegisterResponse, error) {
	return &apiPb.RegisterResponse{
		Id: "asf",
	}, nil
}

func (s *serverSuccess) UnRegister(context.Context, *apiPb.UnRegisterRequest) (*apiPb.UnRegisterResponse, error) {
	s.count += 1
	return &apiPb.UnRegisterResponse{
		Id: "asf",
	}, nil
}

func (c configSecondMock) GetSquzyServer() string {
	return "localhost:14556"
}

func (c configSecondMock) GetSquzyServerTimeout() time.Duration {
	return time.Second
}

type grpcToolsMock struct {
}

func (g grpcToolsMock) GetConnection(address string, timeout time.Duration, option ...grpc.DialOption) (*grpc.ClientConn, error) {
	return &grpc.ClientConn{}, nil
}

type configMock struct {
}

func (c configMock) GetAgentName() string {
	return ""
}

func (c configMock) GetExecutionTimeout() time.Duration {
	return time.Second
}

func (c configMock) GetSquzyServer() string {
	return "localhost:14555"
}

func (c configMock) GetSquzyServerTimeout() time.Duration {
	return time.Second
}

type configErrorMock struct {
}

func (c configErrorMock) GetAgentName() string {
	return ""
}

func (c configErrorMock) GetExecutionTimeout() time.Duration {
	return time.Second
}

func (c configErrorMock) GetSquzyServer() string {
	return "safafasfafsf:12424"
}

func (c configErrorMock) GetSquzyServerTimeout() time.Duration {
	return time.Second
}

type configSuccessMock struct {
}

func (c configSuccessMock) GetAgentName() string {
	panic("implement me")
}

func (c configSuccessMock) GetExecutionTimeout() time.Duration {
	return time.Second
}

func (c configSuccessMock) GetSquzyServer() string {
	return "localhost:13453"
}

func (c configSuccessMock) GetSquzyServerTimeout() time.Duration {
	return time.Second
}

func (c configSuccessMock) GetPort() int32 {
	return 19944
}

type configSuccessSteamMock struct {
}

func (c configSuccessSteamMock) GetAgentName() string {
	return ""
}

func (c configSuccessSteamMock) GetSquzyServer() string {
	return "localhost:13454"
}

func (c configSuccessSteamMock) GetExecutionTimeout() time.Duration {
	return time.Second
}

func (c configSuccessSteamMock) GetSquzyServerTimeout() time.Duration {
	return time.Second
}

func TestNew(t *testing.T) {
	t.Run("Should: create application", func(t *testing.T) {
		a := New(&executorMock{}, &grpcToolsMock{}, &configMock{}, func() (stat *host.InfoStat, err error) {
			return nil, nil
		}, func(agent apiPb.AgentServerClient) (statClient apiPb.AgentServer_SendMetricsClient, err error) {
			return nil, nil
		})
		assert.Implements(t, (*Application)(nil), a)
	})
}

func TestApplication_Run(t *testing.T) {
	t.Run("Should: throw error if cant get host information", func(t *testing.T) {
		a := New(&executorMock{}, grpcTools.New(), &configErrorMock{}, func() (stat *host.InfoStat, err error) {
			return nil, errors.New("asfasff")
		}, func(agent apiPb.AgentServerClient) (statClient apiPb.AgentServer_SendMetricsClient, err error) {
			return nil, nil
		})
		assert.NotEqual(t, nil, a.Run())
	})
	t.Run("Should: throw error if cant connect to squzy server", func(t *testing.T) {
		a := New(&executorMock{}, grpcTools.New(), &configErrorMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{}, nil
		}, func(agent apiPb.AgentServerClient) (statClient apiPb.AgentServer_SendMetricsClient, err error) {
			return nil, nil
		})
		assert.NotEqual(t, nil, a.Run())
	})
	t.Run("Should: throw error if register now working on server side", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 13453))
		grpcServer := grpc.NewServer()
		apiPb.RegisterAgentServerServer(grpcServer, &serverError{})
		go func() {
			_ = grpcServer.Serve(lis)
		}()
		a := New(&executorMock{}, grpcTools.New(), &configSuccessMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{}, nil
		}, func(agent apiPb.AgentServerClient) (statClient apiPb.AgentServer_SendMetricsClient, err error) {
			return nil, nil
		})
		assert.NotEqual(t, nil, a.Run())
	})
	t.Run("Should: throw error if stream now working on server side", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 13454))
		grpcServer := grpc.NewServer()
		apiPb.RegisterAgentServerServer(grpcServer, &serverSteamError{})
		go func() {
			_ = grpcServer.Serve(lis)
		}()
		a := New(&executorMock{}, grpcTools.New(), &configSuccessSteamMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{}, nil
		}, func(agent apiPb.AgentServerClient) (statClient apiPb.AgentServer_SendMetricsClient, err error) {
			return nil, errors.New("asfasf")
		})
		assert.NotEqual(t, nil, a.Run())
	})

	t.Run("Should: not throw error if all works like expected", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 14556))
		msgChan := make(chan *apiPb.SendMetricsRequest)
		grpcServer := grpc.NewServer()
		s := &serverSuccess{
			ch:    msgChan,
			count: 0,
		}
		apiPb.RegisterAgentServerServer(grpcServer, s)
		go func() {
			_ = grpcServer.Serve(lis)
		}()
		time.Sleep(time.Second)
		ch := make(chan *apiPb.SendMetricsRequest)

		a := New(&executorMock{
			ch: ch,
		}, grpcTools.New(), &configSecondMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{}, nil
		}, NewStream)
		go func() {
			_ = a.Run()
		}()
		time.Sleep(time.Second)
		ch <- &apiPb.SendMetricsRequest{
			CpuInfo: &apiPb.CpuInfo{Cpus: []*apiPb.CpuInfo_CPU{{
				Load: 5,
			}}},
		}
		value := <-msgChan
		assert.EqualValues(t, []*apiPb.CpuInfo_CPU{{
			Load: 5,
		}}, value.CpuInfo.Cpus)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		time.Sleep(time.Second * 2)
		assert.Equal(t, 1, s.count)
	})
}

type client struct {
}

func (c client) GetByAgentUniqName(ctx context.Context, in *apiPb.GetByAgentUniqNameRequest, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (c client) GetAgentList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (c client) SendMetrics(ctx context.Context, opts ...grpc.CallOption) (apiPb.AgentServer_SendMetricsClient, error) {
	return nil, nil
}

func (c client) Register(ctx context.Context, in *apiPb.RegisterRequest, opts ...grpc.CallOption) (*apiPb.RegisterResponse, error) {
	panic("implement me")
}

func (c client) UnRegister(ctx context.Context, in *apiPb.UnRegisterRequest, opts ...grpc.CallOption) (*apiPb.UnRegisterResponse, error) {
	panic("implement me")
}

type clientError struct {
}

func (c clientError) GetByAgentUniqName(ctx context.Context, in *apiPb.GetByAgentUniqNameRequest, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (c clientError) GetAgentList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (c clientError) Register(ctx context.Context, in *apiPb.RegisterRequest, opts ...grpc.CallOption) (*apiPb.RegisterResponse, error) {
	panic("implement me")
}

func (c clientError) UnRegister(ctx context.Context, in *apiPb.UnRegisterRequest, opts ...grpc.CallOption) (*apiPb.UnRegisterResponse, error) {
	panic("implement me")
}

func (c clientError) SendMetrics(ctx context.Context, opts ...grpc.CallOption) (apiPb.AgentServer_SendMetricsClient, error) {
	return nil, errors.New("asf")
}

func TestNewStream(t *testing.T) {
	t.Run("Should: not throw error", func(t *testing.T) {
		_, err := NewStream(&client{})
		assert.Equal(t, nil, err)
	})
	t.Run("Should: throw error", func(t *testing.T) {
		_, err := NewStream(&clientError{})
		assert.NotEqual(t, nil, err)
	})
}
