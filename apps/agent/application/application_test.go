package application

import (
	"context"
	"errors"
	"fmt"
	agentPb "github.com/squzy/squzy_generated/generated/agent/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"net"
	"squzy/apps/internal/grpcTools"
	"testing"
	"time"
)

type executorMock struct {
	ch chan *agentPb.SendStat
}

func (e *executorMock) Execute() chan *agentPb.SendStat {
	return e.ch
}

type serverError struct {
}

func (s *serverError) Register(r agentPb.AgentServer_RegisterServer) error {
	return errors.New("awf")
}

func (s *serverError) GetList(context.Context, *agentPb.GetListRequest) (*agentPb.GetListResponse, error) {
	panic("implement me")
}

type configSecondMock struct {
}

func (c configSecondMock) GetExecutionTimeout() time.Duration {
	return time.Second * 5
}

type serverSuccess struct {
	ch chan *agentPb.SendStat
}

func (s serverSuccess) Register(rq agentPb.AgentServer_RegisterServer) error {
	rq.Send(&agentPb.RegisterResponse{
		Id: "asfsaf",
	})
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

type grpcToolsErrorMock struct {
}

func (g grpcToolsErrorMock) GetConnection(address string, timeout time.Duration, option ...grpc.DialOption) (*grpc.ClientConn, error) {
	return nil, errors.New("asfaf")
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

func TestNew(t *testing.T) {
	t.Run("Should: create application", func(t *testing.T) {
		a := New(&executorMock{}, &grpcToolsMock{}, &configMock{}, )
		assert.Implements(t, (*Application)(nil), a)
	})
}

func TestApplication_Run(t *testing.T) {
	t.Run("Should: throw error if cant connect to squzy server", func(t *testing.T) {
		a := New(&executorMock{}, grpcTools.New(), &configErrorMock{})
		assert.NotEqual(t, nil, a.Run())
	})
	t.Run("Should: throw error if register now working on server side", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 13453))
		grpcServer := grpc.NewServer()
		agentPb.RegisterAgentServerServer(grpcServer, &serverError{})
		go func() {
			_ = grpcServer.Serve(lis)
		}()
		a := New(&executorMock{}, grpcTools.New(), &configSuccessMock{})
		assert.NotEqual(t, nil, a.Run())
	})
	t.Run("Should: not throw error if all works like expected", func(t *testing.T) {
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", 14556))
		msgChan := make(chan *agentPb.SendStat)
		grpcServer := grpc.NewServer()
		agentPb.RegisterAgentServerServer(grpcServer, &serverSuccess{
			ch: msgChan,
		})
		go func() {
			grpcServer.Serve(lis)
		}()
		ch := make(chan *agentPb.SendStat)

		a := New(&executorMock{
			ch: ch,
		}, grpcTools.New(), &configSecondMock{})
		assert.Equal(t, nil, a.Run())
		ch <- &agentPb.SendStat{
			CpuInfo: &agentPb.CpuInfo{Cpus: []*agentPb.CpuInfo_CPU{{
				Load: 5,
			}}},
		}
		value := <-msgChan
		assert.EqualValues(t, []*agentPb.CpuInfo_CPU{{
			Load: 5,
		}}, value.CpuInfo.Cpus)
	})
}
