package application

import (
	"context"
	"errors"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/squzy/squzy/apps/agent_client/config"
	agent_executor "github.com/squzy/squzy/internal/agent-executor"
	"github.com/squzy/squzy/internal/helpers"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"go.uber.org/atomic"
	"google.golang.org/grpc"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	RetryInterval = time.Second * 5
)

var (
	errCantRegister  = errors.New("cant register client")
	errCantGetStream = errors.New("cant get stream")
)

type getSteamFn func(agent apiPb.AgentServerClient) (apiPb.AgentServer_SendMetricsClient, error)

type application struct {
	ctx           context.Context
	cancel        context.CancelFunc
	executor      agent_executor.AgentExecutor
	config        config.Config
	hostStatFn    func() (*host.InfoStat, error)
	getStreamFn   getSteamFn
	isStreamAvail bool
	buffer        []*apiPb.SendMetricsRequest
	client        apiPb.AgentServerClient
	interrupt     chan os.Signal
	mutex         sync.Mutex
	stream        apiPb.AgentServer_SendMetricsClient
	streamFailed  atomic.Bool
}

func (a *application) register(hostStat *host.InfoStat) (*apiPb.RegisterResponse, error) {
	ctx, cancel := helpers.TimeoutContext(a.ctx, RetryInterval)
	defer cancel()
	res, err := a.client.Register(ctx, &apiPb.RegisterRequest{
		AgentName: a.config.GetAgentName(),
		Interval:  int64(a.config.GetInterval()) / 1e9,
		HostInfo: &apiPb.HostInfo{
			HostName: hostStat.Hostname,
			Os:       hostStat.OS,
			PlatformInfo: &apiPb.PlatformInfo{
				Name:    hostStat.Platform,
				Family:  hostStat.PlatformFamily,
				Version: hostStat.PlatformVersion,
			},
		},
		Time: timestamp.Now(),
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *application) registerWithRetry(hostStat *host.InfoStat) (*apiPb.RegisterResponse, error) {
	res, err := a.register(hostStat)
	if err == nil {
		return res, nil
	}

	if !a.config.Retry() {
		return nil, errCantRegister
	}

	count := a.config.RetryCount()
	if count <= 0 {
		logger.Infof("Can't connect to server, will try connect in %s", RetryInterval.String())
		for {
			time.Sleep(RetryInterval)
			resRetry, errRetry := a.register(hostStat)
			if errRetry == nil {
				return resRetry, nil
			}
		}
	}
	logger.Infof("Can't connect to server, will try connect in %s times %d", RetryInterval.String(), count)
	for count > 0 {
		time.Sleep(RetryInterval)
		resRetry, errRetry := a.register(hostStat)
		if errRetry == nil {
			return resRetry, nil
		}
		count -= 1
	}
	return nil, errCantRegister
}

func (a *application) getStream() (apiPb.AgentServer_SendMetricsClient, error) {
	return a.getStreamFn(a.client)
}

func (a *application) flushBuffer(stream apiPb.AgentServer_SendMetricsClient) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for _, v := range a.buffer {
		// Because we use wait for ready
		_ = stream.Send(v)
	}
	a.buffer = []*apiPb.SendMetricsRequest{}
	return nil
}

func (a *application) sendToStream(clientId string) {
	input := a.executor.Execute()
	for {
		select {
		case stat := <-input:
			stat.AgentId = clientId
			stat.AgentName = a.config.GetAgentName()
			// what we should do if squzy server cant get msg

			metric := &apiPb.SendMetricsRequest{
				Msg: &apiPb.SendMetricsRequest_Metric{
					Metric: stat,
				},
			}
			if a.streamFailed.Load() {
				a.mutex.Lock()
				a.buffer = append(a.buffer, metric)
				a.mutex.Unlock()
				continue
			}
			err := a.stream.Send(metric)
			if err != nil {
				a.mutex.Lock()
				a.buffer = append(a.buffer, metric)
				a.mutex.Unlock()
				a.streamFailed.Store(true)
				go func() {
					for {
						newStream, errSteam := a.getStream()
						if errSteam != nil {
							continue
						}
						errFlush := a.flushBuffer(newStream)
						if errFlush == nil {
							a.stream = newStream
							a.streamFailed.Store(false)
							break
						}
					}
				}()
			}
		case <-a.ctx.Done():
			return
		}
	}
}

func (a *application) Run() error {
	hostStat, err := a.hostStatFn()

	if err != nil {
		return err
	}

	resp, err := a.registerWithRetry(hostStat)

	if err != nil {
		return err
	}

	stream, err := a.getStream()

	if err != nil {
		return err
	}

	a.stream = stream

	signal.Notify(a.interrupt, syscall.SIGTERM, syscall.SIGINT)

	defer signal.Stop(a.interrupt)

	logger.Infof("Registered with ID=%s", resp.Id)

	go a.sendToStream(resp.Id)

	<-a.interrupt

	a.cancel()

	_ = a.stream.Send(&apiPb.SendMetricsRequest{
		Msg: &apiPb.SendMetricsRequest_Disconnect_{
			Disconnect: &apiPb.SendMetricsRequest_Disconnect{
				AgentId: resp.Id,
				Time:    timestamp.Now(),
			},
		},
	})

	_ = a.stream.CloseSend()

	ctxClose, cancelClose := helpers.TimeoutContext(a.ctx, 0)
	defer cancelClose()

	_, _ = a.client.UnRegister(ctxClose, &apiPb.UnRegisterRequest{
		Id:   resp.Id,
		Time: timestamp.Now(),
	})

	return nil
}

type Application interface {
	Run() error
}

func New(
	ctx context.Context,
	executor agent_executor.AgentExecutor,
	client apiPb.AgentServerClient,
	config config.Config,
	hostStatFn func() (*host.InfoStat, error),
	getStreamFn getSteamFn,
	interrupt chan os.Signal,
) Application {
	appCtx, cancel := context.WithCancel(ctx)
	return &application{
		ctx:         appCtx,
		cancel:      cancel,
		config:      config,
		client:      client,
		executor:    executor,
		hostStatFn:  hostStatFn,
		getStreamFn: getStreamFn,
		interrupt:   interrupt,
	}
}

func NewStream(agent apiPb.AgentServerClient) (apiPb.AgentServer_SendMetricsClient, error) {
	return agent.SendMetrics(context.Background(), grpc.WaitForReady(true))
}
