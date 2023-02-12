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

func (c configSecondMock) Retry() bool {
	//TODO implement me
	return true
}

func (c configSecondMock) RetryCount() int32 {
	//TODO implement me
	return 0
}

func (c configSecondMock) GetInterval() time.Duration {
	return time.Second * 5
}

func (c configSecondMock) GetAgentName() string {
	return ""
}

type serverSuccess struct {
	apiPb.UnimplementedAgentServerServer
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
			return rq.SendAndClose(&empty.Empty{})
		}
		if res != nil {
			fmt.Println(res.GetMetric())
			s.mutex.Lock()
			s.count += 1
			s.mutex.Unlock()
			s.ch <- res
		}
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

func (c configMock) Retry() bool {
	//TODO implement me
	panic("implement me")
}

func (c configMock) RetryCount() int32 {
	//TODO implement me
	panic("implement me")
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

func (c configErrorMock) Retry() bool {
	//TODO implement me
	return false
}

func (c configErrorMock) RetryCount() int32 {
	//TODO implement me
	panic("implement me")
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

type mockRegisterErrorClient struct {
}

func (m mockRegisterErrorClient) Register(ctx context.Context, in *apiPb.RegisterRequest, opts ...grpc.CallOption) (*apiPb.RegisterResponse, error) {
	//TODO implement me
	return nil, errors.New("hahah")
}

func (m mockRegisterErrorClient) GetByAgentName(ctx context.Context, in *apiPb.GetByAgentNameRequest, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockRegisterErrorClient) GetAgentById(ctx context.Context, in *apiPb.GetAgentByIdRequest, opts ...grpc.CallOption) (*apiPb.AgentItem, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockRegisterErrorClient) UnRegister(ctx context.Context, in *apiPb.UnRegisterRequest, opts ...grpc.CallOption) (*apiPb.UnRegisterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockRegisterErrorClient) GetAgentList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockRegisterErrorClient) SendMetrics(ctx context.Context, opts ...grpc.CallOption) (apiPb.AgentServer_SendMetricsClient, error) {
	//TODO implement me
	panic("implement me")
}

type configErrorRetryNegativeMock struct {
}

func (c configErrorRetryNegativeMock) GetAgentServer() string {
	//TODO implement me
	return ""
}

func (c configErrorRetryNegativeMock) GetInterval() time.Duration {
	//TODO implement me
	return time.Second * 10
}

func (c configErrorRetryNegativeMock) GetAgentName() string {
	//TODO implement me
	return ""
}

func (c configErrorRetryNegativeMock) Retry() bool {
	//TODO implement me
	return true
}

func (c configErrorRetryNegativeMock) RetryCount() int32 {
	//TODO implement me
	return 0
}

type configErrorRetryFiveMock struct {
}

func (c configErrorRetryFiveMock) GetAgentServer() string {
	//TODO implement me
	return ""
}

func (c configErrorRetryFiveMock) GetInterval() time.Duration {
	//TODO implement me
	return time.Second * 109
}

func (c configErrorRetryFiveMock) GetAgentName() string {
	//TODO implement me
	return ""
}

func (c configErrorRetryFiveMock) Retry() bool {
	//TODO implement me
	return true
}

func (c configErrorRetryFiveMock) RetryCount() int32 {
	//TODO implement me
	return 5
}

type configErrorRetryTenMock struct {
}

func (c configErrorRetryTenMock) GetAgentServer() string {
	//TODO implement me
	return ""
}

func (c configErrorRetryTenMock) GetInterval() time.Duration {
	//TODO implement me
	return time.Second * 10
}

func (c configErrorRetryTenMock) GetAgentName() string {
	//TODO implement me
	return "zsd"
}

func (c configErrorRetryTenMock) Retry() bool {
	//TODO implement me
	return true
}

func (c configErrorRetryTenMock) RetryCount() int32 {
	//TODO implement me
	return 1
}

type mockRegisterErrorThreeTimesClient struct {
	count int32
	mutex sync.Mutex
}

func (m *mockRegisterErrorThreeTimesClient) Register(ctx context.Context, in *apiPb.RegisterRequest, opts ...grpc.CallOption) (*apiPb.RegisterResponse, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.count += 1
	if m.count == 3 {
		return &apiPb.RegisterResponse{}, nil
	}
	return nil, errors.New("no")
}

func (m mockRegisterErrorThreeTimesClient) GetByAgentName(ctx context.Context, in *apiPb.GetByAgentNameRequest, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockRegisterErrorThreeTimesClient) GetAgentById(ctx context.Context, in *apiPb.GetAgentByIdRequest, opts ...grpc.CallOption) (*apiPb.AgentItem, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockRegisterErrorThreeTimesClient) UnRegister(ctx context.Context, in *apiPb.UnRegisterRequest, opts ...grpc.CallOption) (*apiPb.UnRegisterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockRegisterErrorThreeTimesClient) GetAgentList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*apiPb.GetAgentListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockRegisterErrorThreeTimesClient) SendMetrics(ctx context.Context, opts ...grpc.CallOption) (apiPb.AgentServer_SendMetricsClient, error) {
	//TODO implement me
	panic("implement me")
}

func TestNew(t *testing.T) {
	t.Run("Should: create application", func(t *testing.T) {
		a := New(context.Background(), &executorMock{}, nil, &configMock{}, func() (stat *host.InfoStat, err error) {
			return nil, nil
		}, func(agent apiPb.AgentServerClient) (statClient apiPb.AgentServer_SendMetricsClient, err error) {
			return nil, nil
		}, make(chan os.Signal, 1))
		assert.Implements(t, (*Application)(nil), a)
	})
}

func TestApplication_Run(t *testing.T) {
	t.Run("Should: throw error if cant get host information", func(t *testing.T) {
		a := New(context.Background(), &executorMock{}, nil, &configErrorMock{}, func() (stat *host.InfoStat, err error) {
			return nil, errors.New("asfasff")
		}, func(agent apiPb.AgentServerClient) (statClient apiPb.AgentServer_SendMetricsClient, err error) {
			return nil, nil
		}, make(chan os.Signal, 1))
		assert.NotEqual(t, nil, a.Run())
	})

	t.Run("Should: throw error because cant register and no retry", func(t *testing.T) {
		a := New(context.Background(), &executorMock{}, &mockRegisterErrorClient{}, &configErrorMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{}, nil
		}, func(agent apiPb.AgentServerClient) (statClient apiPb.AgentServer_SendMetricsClient, err error) {
			return nil, nil
		}, make(chan os.Signal, 1))
		assert.NotEqual(t, nil, a.Run())
	})

	t.Run("Should: throw error because cant register and retry 1 times", func(t *testing.T) {
		a := New(context.Background(), &executorMock{}, &mockRegisterErrorClient{}, &configErrorRetryTenMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{}, nil
		}, func(agent apiPb.AgentServerClient) (statClient apiPb.AgentServer_SendMetricsClient, err error) {
			return nil, nil
		}, make(chan os.Signal, 1))
		assert.NotEqual(t, nil, a.Run())
	})

	t.Run("Should: throw error because cant get stream with retry 5", func(t *testing.T) {
		a := New(context.Background(), &executorMock{}, &mockRegisterErrorThreeTimesClient{}, &configErrorRetryFiveMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{}, nil
		}, func(agent apiPb.AgentServerClient) (statClient apiPb.AgentServer_SendMetricsClient, err error) {
			return nil, errors.New("haha")
		}, make(chan os.Signal, 1))
		assert.NotEqual(t, nil, a.Run())
	})

	t.Run("Should: throw error because cant get stream", func(t *testing.T) {
		a := New(context.Background(), &executorMock{}, &mockRegisterErrorThreeTimesClient{}, &configErrorRetryNegativeMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{}, nil
		}, func(agent apiPb.AgentServerClient) (statClient apiPb.AgentServer_SendMetricsClient, err error) {
			return nil, errors.New("haha")
		}, make(chan os.Signal, 1))
		assert.NotEqual(t, nil, a.Run())
	})

	t.Run("Should: not throw error if all works like expected", func(t *testing.T) {
		lis, err := net.Listen("tcp", ":0")
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

		//time.Sleep(time.Second * 2)
		port := fmt.Sprintf(":%d", lis.Addr().(*net.TCPAddr).Port)
		conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())

		assert.NoError(t, err)

		client := apiPb.NewAgentServerClient(conn)

		ch := make(chan *apiPb.Metric)
		inter := make(chan os.Signal, 1)
		a := New(context.Background(), &executorMock{
			ch: ch,
		}, client, &configSecondMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{}, nil
		}, (&mockGetSteam{}).NewStream, inter)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			err = a.Run()
			assert.Equal(t, nil, err)
			wg.Done()
		}()
		ch <- &apiPb.Metric{
			CpuInfo: &apiPb.CpuInfo{Cpus: []*apiPb.CpuInfo_CPU{{
				Load: 1,
			}}},
		}
		value := <-msgChan

		assert.EqualValues(t, 1, (value.Msg).(*apiPb.SendMetricsRequest_Metric).Metric.CpuInfo.Cpus[0].Load)

		grpcServer.Stop()
		time.Sleep(time.Second * 2)
		ch <- &apiPb.Metric{
			CpuInfo: &apiPb.CpuInfo{Cpus: []*apiPb.CpuInfo_CPU{{
				Load: 2,
			}}},
		}

		ch <- &apiPb.Metric{
			CpuInfo: &apiPb.CpuInfo{Cpus: []*apiPb.CpuInfo_CPU{{
				Load: 3,
			}}},
		}

		time.Sleep(time.Second * 10)
		lis2, err := net.Listen("tcp", port)
		assert.Equal(t, nil, err)
		grpcServer = grpc.NewServer()
		apiPb.RegisterAgentServerServer(grpcServer, s)

		go func() {
			err = grpcServer.Serve(lis2)
			assert.Equal(t, nil, err)
		}()
		time.Sleep(time.Second * 2)
		value = <-msgChan
		assert.EqualValues(t, 2, (value.Msg).(*apiPb.SendMetricsRequest_Metric).Metric.CpuInfo.Cpus[0].Load)

		// Again break
		fmt.Println("Server broke again")
		grpcServer.Stop()
		time.Sleep(time.Second * 1)
		ch <- &apiPb.Metric{
			CpuInfo: &apiPb.CpuInfo{Cpus: []*apiPb.CpuInfo_CPU{{
				Load: 4,
			}}},
		}
		time.Sleep(time.Second * 10)
		lis3, err := net.Listen("tcp", port)
		assert.Equal(t, nil, err)
		grpcServer = grpc.NewServer()
		apiPb.RegisterAgentServerServer(grpcServer, s)

		go func() {
			err = grpcServer.Serve(lis3)
			assert.Equal(t, nil, err)
		}()
		time.Sleep(time.Second * 2)
		value = <-msgChan
		assert.EqualValues(t, 3, (value.Msg).(*apiPb.SendMetricsRequest_Metric).Metric.CpuInfo.Cpus[0].Load)

		value = <-msgChan
		assert.EqualValues(t, 4, (value.Msg).(*apiPb.SendMetricsRequest_Metric).Metric.CpuInfo.Cpus[0].Load)

		ch <- &apiPb.Metric{
			CpuInfo: &apiPb.CpuInfo{Cpus: []*apiPb.CpuInfo_CPU{{
				Load: 5,
			}}},
		}
		value = <-msgChan
		assert.EqualValues(t, 5, (value.Msg).(*apiPb.SendMetricsRequest_Metric).Metric.CpuInfo.Cpus[0].Load)

		inter <- syscall.SIGTERM
		wg.Wait()
		// 5 msg + disconnect + unregister
		grpcServer.Stop()
		time.Sleep(time.Second * 5)
		assert.Equal(t, 7, s.count)
	})
	t.Run("Should: not throw error if all works like expected", func(t *testing.T) {
		lis, err := net.Listen("tcp", ":0")
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

		//time.Sleep(time.Second * 2)
		port := fmt.Sprintf(":%d", lis.Addr().(*net.TCPAddr).Port)
		conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())

		assert.NoError(t, err)

		client := apiPb.NewAgentServerClient(conn)
		ctx, cancel := context.WithCancel(context.Background())

		ch := make(chan *apiPb.Metric)
		inter := make(chan os.Signal, 1)
		a := New(ctx, &executorMock{
			ch: ch,
		}, client, &configSecondMock{}, func() (stat *host.InfoStat, err error) {
			return &host.InfoStat{}, nil
		}, (&mockGetSteam{}).NewStream, inter)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			err = a.Run()
			assert.Equal(t, nil, err)
			wg.Done()
		}()
		time.Sleep(time.Second * 2)
		cancel()
		time.Sleep(time.Second * 2)
		inter <- syscall.SIGTERM
		wg.Wait()
		grpcServer.Stop()
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

type mockGetSteam struct {
	c int32
}

func (m *mockGetSteam) NewStream(agent apiPb.AgentServerClient) (apiPb.AgentServer_SendMetricsClient, error) {
	m.c += 1
	if m.c >= 2 && m.c < 10 {
		return nil, errors.New("hjahaha")
	}
	return agent.SendMetrics(context.Background(), grpc.WaitForReady(true))
}
