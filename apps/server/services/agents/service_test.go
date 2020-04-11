package agents

import (
	"context"
	agentPb "github.com/squzy/squzy_generated/generated/agent/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: implement interface of agents logger", func(t *testing.T) {
		t.Run("Should: implement interface of the LoggerServer", func(t *testing.T) {
			s := New(nil)
			assert.Implements(t, (*agentPb.AgentServerServer)(nil), s)
		})
	})
}

func TestService_GetList(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		t.Run("Should: implement interface of the LoggerServer", func(t *testing.T) {
			s := New(nil)
			_, err := s.GetList(context.Background(), &agentPb.GetListRequest{})
			assert.Equal(t, nil, err)
		})
	})
}

func TestService_Register(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		t.Run("Should: implement interface of the LoggerServer", func(t *testing.T) {
			s := New(nil)
			_, err := s.Register(context.Background(), &agentPb.RegisterRequest{})
			assert.Equal(t, nil, err)
		})
	})
}

func TestService_UnRegister(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		t.Run("Should: implement interface of the LoggerServer", func(t *testing.T) {
			s := New(nil)
			_, err := s.UnRegister(context.Background(), &agentPb.UnRegisterRequest{})
			assert.Equal(t, nil, err)
		})
	})
}

func TestService_SendStat(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		t.Run("Should: implement interface of the LoggerServer", func(t *testing.T) {
			s := New(nil)
			err := s.SendStat(nil)
			assert.Equal(t, nil, err)
		})
	})
}