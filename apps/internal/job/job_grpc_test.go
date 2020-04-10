package job

import (
	"context"
	"errors"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	health_check "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"net"
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

func TestNewGrpcJob(t *testing.T) {
	t.Run("Should: Should implement interface Job", func(t *testing.T) {
		job := NewGrpcJob("test", "localhost", 9090, 0, []grpc.DialOption{}, []grpc.CallOption{})
		assert.Implements(t, (*Job)(nil), job)
	})
}

func TestGrpcJob_Do(t *testing.T) {
	t.Run("Test: Testing grpc health_check:", func(t *testing.T) {
		t.Run("Should: Return nil", func(t *testing.T) {
			lis, _ := net.Listen("tcp", ":9090")
			grpcServer := grpc.NewServer()
			health_check.RegisterHealthServer(grpcServer, &server{})
			go func() {
				_ = grpcServer.Serve(lis)
			}()
			job := NewGrpcJob("test", "localhost", 9090, 1, []grpc.DialOption{grpc.WithInsecure()}, []grpc.CallOption{})
			assert.Equal(t, clientPb.StatusCode_OK, job.Do().GetLogData().Code)
			grpcServer.Stop()
		})

		t.Run("Should: Return error because server response more than 10 second (default timeout)", func(t *testing.T) {
			lis, _ := net.Listen("tcp", ":9090")
			grpcServer := grpc.NewServer()
			health_check.RegisterHealthServer(grpcServer, &serverLong{})
			go func() {
				_ = grpcServer.Serve(lis)
			}()
			job := NewGrpcJob("test", "localhost", 9090, 2, []grpc.DialOption{grpc.WithInsecure()}, []grpc.CallOption{})
			assert.Equal(t, clientPb.StatusCode_Error, job.Do().GetLogData().Code)
			grpcServer.Stop()
		})

		t.Run("Should: Return grpcNotServing error", func(t *testing.T) {
			lis, _ := net.Listen("tcp", ":9090")
			grpcServer := grpc.NewServer()
			health_check.RegisterHealthServer(grpcServer, &errorServer{})
			go func() {
				_ = grpcServer.Serve(lis)
			}()
			job := NewGrpcJob("test", "localhost", 9090, 0, []grpc.DialOption{grpc.WithInsecure()}, []grpc.CallOption{})
			assert.Equal(t, grpcNotServing.Error(), job.Do().GetLogData().Description)
			grpcServer.Stop()
		})

		t.Run("Should: Return connTimeoutError error", func(t *testing.T) {
			job := NewGrpcJob("test", "localhost", 9091, 0, []grpc.DialOption{grpc.WithInsecure()}, []grpc.CallOption{})
			assert.Equal(t, connTimeoutError.Error(), job.Do().GetLogData().Description)
		})

		t.Run("Should: Return wrongConnectConfigError error", func(t *testing.T) {
			job := NewGrpcJob("test", "localhost", 9091, 0, []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}, []grpc.CallOption{})
			assert.Equal(t, wrongConnectConfigError.Error(), job.Do().GetLogData().Description)
		})
	})
}
