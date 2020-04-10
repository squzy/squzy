package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/host"
	agentPb "github.com/squzy/squzy_generated/generated/agent/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"net"
	"squzy/apps/internal/grpcTools"
	"syscall"
	"testing"
	"time"
)

type executorMock struct {
	ch chan *agentPb.SendStatRequest
}

func (e *executorMock) Execute() chan *agentPb.SendStatRequest {
	return e.ch
}

type serverError struct {
}

func (s serverError) Register(context.Context, *agentPb.RegisterRequest) (*agentPb.RegisterResponse, error) {
	return nil, errors.New("traarar")
}

func (s serverError) UnRegister(context.Context, *agentPb.UnRegisterRequest) (*agentPb.UnRegisterResponse, error) {
	return &agentPb.UnRegisterResponse{

	}, nil
}

func (s serverError) SendStat(agentPb.AgentServer_SendStatServer) error {
	panic("implement me")
}

func (s serverError) GetList(context.Context, *agentPb.GetListRequest) (*agentPb.GetListResponse, error) {
	panic("implement me")
}

type configSecondMock struct {
}

type serverSteamError struct {
}

func (s serverSteamError) Register(context.Context, *agentPb.RegisterRequest) (*agentPb.RegisterResponse, error) {
	return &agentPb.RegisterResponse{
		Id: "",
	}, nil
}

func (s serverSteamError) UnRegister(context.Context, *agentPb.UnRegisterRequest) (*agentPb.UnRegisterResponse, error) {
	panic("implement me")
}

func (s serverSteamError) SendStat(req agentPb.AgentServer_SendStatServer) error {
	return nil
}

func (s serverSteamError) GetList(context.Context, *agentPb.GetListRequest) (*agentPb.GetListResponse, error) {
	panic("implement me")
}

func (c configSecondMock) GetExecutionTimeout() time.Duration {
	return time.Second * 5
}

type serverSuccess struct {
	ch    chan *agentPb.SendStatRequest
	count int
}

func (s *serverSuccess) Register(context.Context, *agentPb.RegisterRequest) (*agentPb.RegisterResponse, error) {
	return &agentPb.RegisterResponse{
		Id: "asf",
	}, nil
}

func (s *serverSuccess) UnRegister(context.Context, *agentPb.UnRegisterRequest) (*agentPb.UnRegisterResponse, error) {
	s.count += 1
	return &agentPb.UnRegisterResponse{
		Id: "asf",
	}, nil
}

func (s serverSuccess) SendStat(rq agentPb.AgentServer_SendStatServer) error {
	res, _ := rq.Recv()
	s.ch <- res
	return nil
}

func (s serverSuccess) GetList(context.Context, *agentPb.GetListRequest) (*agentPb.GetListResponse, error) {
	panic("implement me")
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
		}, func(agent agentPb.AgentServerClient) (statClient agentPb.AgentServer_SendStatClient, err error) {
			return nil, nil
		})
		assert.Implements(t, (*Application)(nil), a)
	})
}

func TestApplication_Run(t *testing.T) {
	t.Run("Should: throw error if cant get host information", func(t *testing.T) {
		a := New(&executorMock{}, grpcTools.New(), &configErrorMock{}, func() (stat *host.InfoStat, err error) {
			return nil, errors.New("asfasff")
		}, func(agent agentPb.AgentServerClient) (statClient agentPb.AgentServer_SendStatClient, err error) {
			return nil, nil
		})
		assert.NotEqual(t, nil, a.Run())
	})
	t.Run("Should: throw error if cant connect to squzy server", func(t *testing.T) {
		a := New(&executorMock{}, grpcTools.New(), &configErrorMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{

			}, nil
		}, func(agent agentPb.AgentServerClient) (statClient agentPb.AgentServer_SendStatClient, err error) {
			return nil, nil
		})
		assert.NotEqual(t, nil, a.Run())
	})
	t.Run("Should: throw error if register now working on server side", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 13453))
		grpcServer := grpc.NewServer()
		agentPb.RegisterAgentServerServer(grpcServer, &serverError{})
		go func() {
			_ = grpcServer.Serve(lis)
		}()
		a := New(&executorMock{}, grpcTools.New(), &configSuccessMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{

			}, nil
		}, func(agent agentPb.AgentServerClient) (statClient agentPb.AgentServer_SendStatClient, err error) {
			return nil, nil
		})
		assert.NotEqual(t, nil, a.Run())
	})
	t.Run("Should: throw error if stream now working on server side", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 13454))
		grpcServer := grpc.NewServer()
		agentPb.RegisterAgentServerServer(grpcServer, &serverSteamError{})
		go func() {
			_ = grpcServer.Serve(lis)
		}()
		a := New(&executorMock{}, grpcTools.New(), &configSuccessSteamMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{

			}, nil
		}, func(agent agentPb.AgentServerClient) (statClient agentPb.AgentServer_SendStatClient, err error) {
			return nil, errors.New("asfasf")
		})
		assert.NotEqual(t, nil, a.Run())
	})

	t.Run("Should: not throw error if all works like expected", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 14556))
		msgChan := make(chan *agentPb.SendStatRequest)
		grpcServer := grpc.NewServer()
		s := &serverSuccess{
			ch:    msgChan,
			count: 0,
		}
		agentPb.RegisterAgentServerServer(grpcServer, s)
		go func() {
			_ = grpcServer.Serve(lis)
		}()
		time.Sleep(time.Second)
		ch := make(chan *agentPb.SendStatRequest)

		a := New(&executorMock{
			ch: ch,
		}, grpcTools.New(), &configSecondMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{}, nil
		}, NewStream)
		go func() {
			_ = a.Run()
		}()
		time.Sleep(time.Second)
		ch <- &agentPb.SendStatRequest{
			CpuInfo: &agentPb.CpuInfo{Cpus: []*agentPb.CpuInfo_CPU{{
				Load: 5,
			}}},
		}
		value := <-msgChan
		assert.EqualValues(t, []*agentPb.CpuInfo_CPU{{
			Load: 5,
		}}, value.CpuInfo.Cpus)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		time.Sleep(time.Second * 2)
		assert.Equal(t, 1, s.count)
	})
}

type client struct {
}

func (c client) Register(ctx context.Context, in *agentPb.RegisterRequest, opts ...grpc.CallOption) (*agentPb.RegisterResponse, error) {
	panic("implement me")
}

func (c client) UnRegister(ctx context.Context, in *agentPb.UnRegisterRequest, opts ...grpc.CallOption) (*agentPb.UnRegisterResponse, error) {
	panic("implement me")
}

func (c client) SendStat(ctx context.Context, opts ...grpc.CallOption) (agentPb.AgentServer_SendStatClient, error) {
	return nil, nil
}

func (c client) GetList(ctx context.Context, in *agentPb.GetListRequest, opts ...grpc.CallOption) (*agentPb.GetListResponse, error) {
	panic("implement me")
}

type clientError struct {
}

func (c clientError) Register(ctx context.Context, in *agentPb.RegisterRequest, opts ...grpc.CallOption) (*agentPb.RegisterResponse, error) {
	panic("implement me")
}

func (c clientError) UnRegister(ctx context.Context, in *agentPb.UnRegisterRequest, opts ...grpc.CallOption) (*agentPb.UnRegisterResponse, error) {
	panic("implement me")
}

func (c clientError) SendStat(ctx context.Context, opts ...grpc.CallOption) (agentPb.AgentServer_SendStatClient, error) {
	return nil, errors.New("asf")
}

func (c clientError) GetList(ctx context.Context, in *agentPb.GetListRequest, opts ...grpc.CallOption) (*agentPb.GetListResponse, error) {
	panic("implement me")
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
