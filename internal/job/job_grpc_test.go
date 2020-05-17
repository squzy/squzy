package job

import (
	"context"
	"errors"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	health_check "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"net"
	scheduler_config_storage "squzy/internal/scheduler-config-storage"
	"testing"
	"time"
)

type serverLong struct {
}

func (s serverLong) Check(ctx context.Context, rq *health_check.HealthCheckRequest) (*health_check.HealthCheckResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	// metadata should present
	if ok != true {
		return nil, errors.New("Metadata not present")
	}

	if md.Get(logMetaData)[0] == "" {
		return nil, errors.New("Metadata not present")
	}
	time.Sleep(time.Second * 4)
	return &health_check.HealthCheckResponse{
		Status: health_check.HealthCheckResponse_SERVING,
	}, nil
}

func (s serverLong) Watch(*health_check.HealthCheckRequest, health_check.Health_WatchServer) error {
	panic("implement me")
}

type server struct {
}

type errorServer struct {
}

func (e errorServer) Check(ctx context.Context, r *health_check.HealthCheckRequest) (*health_check.HealthCheckResponse, error) {
	return &health_check.HealthCheckResponse{
		Status: health_check.HealthCheckResponse_NOT_SERVING,
	}, nil
}

func (e errorServer) Watch(*health_check.HealthCheckRequest, health_check.Health_WatchServer) error {
	panic("implement me")
}

func (s server) Check(ctx context.Context, e *health_check.HealthCheckRequest) (*health_check.HealthCheckResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	// metadata should present
	if ok != true {
		return nil, errors.New("Metadata not present")
	}

	if md.Get(logMetaData)[0] == "" {
		return nil, errors.New("Metadata not present")
	}
	return &health_check.HealthCheckResponse{
		Status: health_check.HealthCheckResponse_SERVING,
	}, nil
}

func (s server) Watch(*health_check.HealthCheckRequest, health_check.Health_WatchServer) error {
	panic("implement me")
}

func TestExecGrpc(t *testing.T) {
	t.Run("Test: Testing grpc health_check:", func(t *testing.T) {
		t.Run("Should: Return nil", func(t *testing.T) {
			lis, _ := net.Listen("tcp", ":9090")
			grpcServer := grpc.NewServer()
			health_check.RegisterHealthServer(grpcServer, &server{})
			go func() {
				_ = grpcServer.Serve(lis)
			}()
			job := ExecGrpc("data", 1, &scheduler_config_storage.GrpcConfig{
				Service: "",
				Host:    "localhost",
				Port:    9090,
			}, grpc.WithInsecure())
			assert.Equal(t, apiPb.SchedulerCode_OK, job.GetLogData().Snapshot.Code)
			grpcServer.Stop()
		})

		t.Run("Should: Return error because server response more than 10 second (default timeout)", func(t *testing.T) {
			lis, _ := net.Listen("tcp", ":9090")
			grpcServer := grpc.NewServer()
			health_check.RegisterHealthServer(grpcServer, &serverLong{})
			go func() {
				_ = grpcServer.Serve(lis)
			}()
			job := ExecGrpc("data", 2, &scheduler_config_storage.GrpcConfig{
				Service: "test",
				Host:    "localhost",
				Port:    9090,
			}, grpc.WithInsecure())
			assert.Equal(t, apiPb.SchedulerCode_Error, job.GetLogData().Snapshot.Code)
			grpcServer.Stop()
		})

		t.Run("Should: Return errGrpcNotServing error", func(t *testing.T) {
			lis, _ := net.Listen("tcp", ":9090")
			grpcServer := grpc.NewServer()
			health_check.RegisterHealthServer(grpcServer, &errorServer{})
			go func() {
				_ = grpcServer.Serve(lis)
			}()
			job := ExecGrpc("", 0, &scheduler_config_storage.GrpcConfig{
				Service: "test",
				Host:    "localhost",
				Port:    9090,
			}, grpc.WithInsecure())
			assert.Equal(t, errGrpcNotServing.Error(), job.GetLogData().Snapshot.Error.Message)
			grpcServer.Stop()
		})

		t.Run("Should: Return errConnTimeoutError error", func(t *testing.T) {
			job := ExecGrpc("", 0, &scheduler_config_storage.GrpcConfig{Service: "test", Host: "localhost", Port: 9091}, grpc.WithInsecure())
			assert.Equal(t, errConnTimeoutError.Error(), job.GetLogData().Snapshot.Error.Message)
		})

		t.Run("Should: Return errWrongConnectConfigError error", func(t *testing.T) {
			job := ExecGrpc("", 0, &scheduler_config_storage.GrpcConfig{Service: "test", Host: "localhost", Port: 9091}, grpc.WithInsecure(), grpc.WithBlock())
			assert.Equal(t, errWrongConnectConfigError.Error(), job.GetLogData().Snapshot.Error.Message)
		})
	})
}
