package grpctools

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("Should: return new grpc tools", func(t *testing.T) {
		tools := New()
		assert.Implements(t, (*GrpcTool)(nil), tools)
	})
}

func TestGrpcTool_GetConnection(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		tools := New()
		_, err := tools.GetConnection("localhost", time.Second, grpc.WithInsecure())
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return error because timeout", func(t *testing.T) {
		tools := New()
		_, err := tools.GetConnection("localhost", time.Second, grpc.WithInsecure(), grpc.WithBlock())
		assert.NotEqual(t, nil, err)
	})
}
