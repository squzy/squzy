package application

import (
	"context"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/squzy/squzy/apps/agent_client/config"
	agent_executor "github.com/squzy/squzy/internal/agent-executor"
	"github.com/squzy/squzy/internal/helpers"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"google.golang.org/grpc"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type getSteamFn func(agent apiPb.AgentServerClient) (apiPb.AgentServer_SendMetricsClient, error)

type application struct {
	executor      agent_executor.AgentExecutor
	config        config.Config
	hostStatFn    func() (*host.InfoStat, error)
	getStreamFn   getSteamFn
	isStreamAvail bool
	buffer        []*apiPb.SendMetricsRequest
	client        apiPb.AgentServerClient
	interrupt     chan os.Signal
	mutex         sync.Mutex
}

func (a *application) getClient(opts ...grpc.DialOption) apiPb.AgentServerClient {
	for {
		ctx, cancel := helpers.TimeoutContext(context.Background(), 0)
		defer cancel()
		conn, err := grpc.DialContext(ctx, a.config.GetAgentServer(), opts...)
		if err == nil {
			return apiPb.NewAgentServerClient(conn)
		}
	}
}

func (a *application) register(hostStat *host.InfoStat) string {
	for {
		ctx, cancel := helpers.TimeoutContext(context.Background(), 0)
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

		if err == nil {
			return res.Id
		}
	}
}

func (a *application) getStream() apiPb.AgentServer_SendMetricsClient {
	for {
		s, err := a.getStreamFn(a.client)
		if err == nil {
			return s
		}
	}
}

func (a *application) Run() error {
	hostStat, err := a.hostStatFn()

	if err != nil {
		return err
	}

	a.client = a.getClient(grpc.WithInsecure(), grpc.WithBlock())

	agentID := a.register(hostStat)

	signal.Notify(a.interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(a.interrupt)

	logger.Infof("Registred with ID=%s", agentID)

	st := a.getStream()

	go func() {
		a.isStreamAvail = true

		for stat := range a.executor.Execute() {
			stat.AgentId = agentID
			stat.AgentName = a.config.GetAgentName()
			// what we should do if squzy server cant get msg

			metric := &apiPb.SendMetricsRequest{
				Msg: &apiPb.SendMetricsRequest_Metric{
					Metric: stat,
				},
			}
			err = st.Send(metric)

			if err == io.EOF {
				a.buffer = append(a.buffer, metric)
				if a.isStreamAvail {
					a.isStreamAvail = false
					go func() {
						st = a.getStream()
						a.isStreamAvail = true
						a.mutex.Lock()
						defer a.mutex.Unlock()
						for _, v := range a.buffer {
							_ = st.Send(v)
						}
						a.buffer = []*apiPb.SendMetricsRequest{}
					}()
				}
			}
		}
	}()

	<-a.interrupt

	_ = st.Send(&apiPb.SendMetricsRequest{
		Msg: &apiPb.SendMetricsRequest_Disconnect_{
			Disconnect: &apiPb.SendMetricsRequest_Disconnect{
				AgentId: agentID,
				Time:    timestamp.Now(),
			},
		},
	})

	_ = st.CloseSend()

	ctxClose, cancelClose := helpers.TimeoutContext(context.Background(), 0)
	defer cancelClose()

	_, _ = a.client.UnRegister(ctxClose, &apiPb.UnRegisterRequest{
		Id:   agentID,
		Time: timestamp.Now(),
	})

	return nil
}

type Application interface {
	Run() error
}

func New(
	executor agent_executor.AgentExecutor,
	config config.Config,
	hostStatFn func() (*host.InfoStat, error),
	getStreamFn getSteamFn,
	interrupt chan os.Signal,
) Application {
	return &application{
		config:      config,
		executor:    executor,
		hostStatFn:  hostStatFn,
		getStreamFn: getStreamFn,
		interrupt:   interrupt,
	}
}

func NewStream(agent apiPb.AgentServerClient) (apiPb.AgentServer_SendMetricsClient, error) {
	return agent.SendMetrics(context.Background(), grpc.WaitForReady(true))
}
