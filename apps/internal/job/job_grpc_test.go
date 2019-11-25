package job

import (
	"context"
	clientPb "github.com/squzy/squzy_generated/generated/logger"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	health_check "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"testing"
)

type server struct {

}

type errorServer struct {

}

func (e errorServer) Check(context.Context, *health_check.HealthCheckRequest) (*health_check.HealthCheckResponse, error) {
	return &health_check.HealthCheckResponse{
		Status: health_check.HealthCheckResponse_NOT_SERVING,
	}, nil
}

func (e errorServer) Watch(*health_check.HealthCheckRequest, health_check.Health_WatchServer) error {
	panic("implement me")
}

func (s server) Check(context.Context, *health_check.HealthCheckRequest) (*health_check.HealthCheckResponse, error) {
	return &health_check.HealthCheckResponse{
		Status: health_check.HealthCheckResponse_SERVING,
	}, nil
}

func (s server) Watch(*health_check.HealthCheckRequest, health_check.Health_WatchServer) error {
	panic("implement me")
}

func TestNewGrpcJob(t *testing.T) {
	t.Run("Should: Should implement interface Job", func(t *testing.T) {
		job := NewGrpcJob("test", "localhost", 9090, []grpc.DialOption{}, []grpc.CallOption{})
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
			job := NewGrpcJob("test", "localhost",9090, []grpc.DialOption{grpc.WithInsecure()}, []grpc.CallOption{})
			assert.Equal(t, clientPb.StatusCode_OK, job.Do().GetLogData().Code)
			grpcServer.Stop()
		})

		t.Run("Should: Return grpcNotServing error", func(t *testing.T) {
			lis, _ := net.Listen("tcp", ":9090")
			grpcServer := grpc.NewServer()
			health_check.RegisterHealthServer(grpcServer, &errorServer{})
			go func() {
				_ = grpcServer.Serve(lis)
			}()
			job := NewGrpcJob("test", "localhost",9090, []grpc.DialOption{grpc.WithInsecure()}, []grpc.CallOption{})
			assert.Equal(t, grpcNotServing.Error(), job.Do().GetLogData().Description)
			grpcServer.Stop()
		})

		t.Run("Should: Return connTimeoutError error", func(t *testing.T) {
			job := NewGrpcJob("test", "localhost",9091, []grpc.DialOption{grpc.WithInsecure()}, []grpc.CallOption{})
			assert.Equal(t, connTimeoutError.Error(), job.Do().GetLogData().Description)
		})

		t.Run("Should: Return wrongConnectConfigError error", func(t *testing.T) {
			job := NewGrpcJob("test", "localhost",9091, []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}, []grpc.CallOption{})
			assert.Equal(t, wrongConnectConfigError.Error(), job.Do().GetLogData().Description)
		})
	})
}