package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/v3/host"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

type executorMock struct {
	ch chan *apiPb.Metric
}

func (e *executorMock) Execute() chan *apiPb.Metric {
	return e.ch
}

type configSecondMock struct {
}

func (c configSecondMock) GetInterval() time.Duration {
	return time.Second * 5
}

func (c configSecondMock) GetAgentName() string {
	return ""
}

type serverSuccess struct {
	ch    chan *apiPb.SendMetricsRequest
	count int
	mutex sync.Mutex
}

func (s *serverSuccess) GetAgentById(ctx context.Context, request *apiPb.GetAgentByIdRequest) (*apiPb.AgentItem, error) {
	panic("implement me")
}

func (s *serverSuccess) GetByAgentName(context.Context, *apiPb.GetByAgentNameRequest) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (s *serverSuccess) GetAgentList(context.Context, *empty.Empty) (*apiPb.GetAgentListResponse, error) {
	panic("implement me")
}

func (s *serverSuccess) SendMetrics(rq apiPb.AgentServer_SendMetricsServer) error {
	for {
		res, err := rq.Recv()
		if err == io.EOF {
			return rq.SendAndClose(&empty.Empty{})
		}
		if err != nil {
			continue
		}
		s.count += 1
		s.ch <- res
	}
}

func (s *serverSuccess) Register(context.Context, *apiPb.RegisterRequest) (*apiPb.RegisterResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count += 1
	return &apiPb.RegisterResponse{
		Id: "asf",
	}, nil
}

func (s *serverSuccess) UnRegister(context.Context, *apiPb.UnRegisterRequest) (*apiPb.UnRegisterResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count += 1
	return &apiPb.UnRegisterResponse{
		Id: "asf",
	}, nil
}

func (c configSecondMock) GetAgentServer() string {
	return "localhost:14556"
}

func (c configSecondMock) GetAgentServerTimeout() time.Duration {
	return time.Second
}

type configMock struct {
}

func (c configMock) GetAgentName() string {
	return ""
}

func (c configMock) GetInterval() time.Duration {
	return time.Second
}

func (c configMock) GetAgentServer() string {
	return "localhost:14555"
}

func (c configMock) GetAgentServerTimeout() time.Duration {
	return time.Second
}

type configErrorMock struct {
}

func (c configErrorMock) GetAgentName() string {
	return ""
}

func (c configErrorMock) GetInterval() time.Duration {
	return time.Second
}

func (c configErrorMock) GetAgentServer() string {
	return "safafasfafsf:12424"
}

func (c configErrorMock) GetAgentServerTimeout() time.Duration {
	return time.Second
}

func TestNew(t *testing.T) {
	t.Run("Should: create application", func(t *testing.T) {
		a := New(&executorMock{}, &configMock{}, func() (stat *host.InfoStat, err error) {
			return nil, nil
		}, func(agent apiPb.AgentServerClient) (statClient apiPb.AgentServer_SendMetricsClient, err error) {
			return nil, nil
		}, make(chan os.Signal, 1))
		assert.Implements(t, (*Application)(nil), a)
	})
}

func TestApplication_Run(t *testing.T) {
	t.Run("Should: throw error if cant get host information", func(t *testing.T) {
		a := New(&executorMock{}, &configErrorMock{}, func() (stat *host.InfoStat, err error) {
			return nil, errors.New("asfasff")
		}, func(agent apiPb.AgentServerClient) (statClient apiPb.AgentServer_SendMetricsClient, err error) {
			return nil, nil
		}, make(chan os.Signal, 1))
		assert.NotEqual(t, nil, a.Run())
	})

	t.Run("Should: not throw error if all works like expected", func(t *testing.T) {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 14556))
		assert.Equal(t, nil, err)
		msgChan := make(chan *apiPb.SendMetricsRequest)
		grpcServer := grpc.NewServer()
		s := &serverSuccess{
			ch:    msgChan,
			count: 0,
		}
		apiPb.RegisterAgentServerServer(grpcServer, s)
		go func() {
			err := grpcServer.Serve(lis)
			assert.Equal(t, nil, err)
		}()

		ch := make(chan *apiPb.Metric)
		inter := make(chan os.Signal, 1)
		a := New(&executorMock{
			ch: ch,
		}, &configSecondMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{}, nil
		}, NewStream, inter)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {

			err = a.Run()
			assert.Equal(t, nil, err)
			wg.Done()
		}()
		ch <- &apiPb.Metric{
			CpuInfo: &apiPb.CpuInfo{Cpus: []*apiPb.CpuInfo_CPU{{
				Load: 5,
			}}},
		}
		value := <-msgChan

		assert.EqualValues(t, []*apiPb.CpuInfo_CPU{{
			Load: 5,
		}}, (value.Msg).(*apiPb.SendMetricsRequest_Metric).Metric.CpuInfo.Cpus)

		grpcServer.Stop()
		ch <- &apiPb.Metric{
			CpuInfo: &apiPb.CpuInfo{Cpus: []*apiPb.CpuInfo_CPU{{
				Load: 5,
			}}},
		}

		lis2, err := net.Listen("tcp", fmt.Sprintf(":%d", 14556))
		assert.Equal(t, nil, err)
		grpcServer = grpc.NewServer()
		apiPb.RegisterAgentServerServer(grpcServer, s)

		go func() {
			err = grpcServer.Serve(lis2)
			assert.Equal(t, nil, err)
		}()

		ch <- &apiPb.Metric{
			CpuInfo: &apiPb.CpuInfo{Cpus: []*apiPb.CpuInfo_CPU{{
				Load: 5,
			}}},
		}
		value = <-msgChan
		assert.EqualValues(t, []*apiPb.CpuInfo_CPU{{
			Load: 5,
		}}, (value.Msg).(*apiPb.SendMetricsRequest_Metric).Metric.CpuInfo.Cpus)

		inter <- syscall.SIGTERM
		wg.Wait()
		assert.Equal(t, 5, s.count)
	})
}

type client struct {
}

func (c client) GetAgentById(ctx context.Context, in *apiPb.GetAgentByIdRequest, opts ...grpc.CallOption) (*apiPb.AgentItem, error) {
	panic("implement me")
}

func (c client) GetByAgentName(ctx context.Context, in *apiPb.GetByAgentNameRequest, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
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

func (c clientError) GetAgentById(ctx context.Context, in *apiPb.GetAgentByIdRequest, opts ...grpc.CallOption) (*apiPb.AgentItem, error) {
	panic("implement me")
}

func (c clientError) GetByAgentName(ctx context.Context, in *apiPb.GetByAgentNameRequest, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
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
