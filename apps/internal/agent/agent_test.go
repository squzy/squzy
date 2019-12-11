package agent

import (
	agentPb "github.com/squzy/squzy_generated/generated/agent/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: create new agent", func(t *testing.T) {
		a := New()
		assert.IsType(t, &agent{}, a)
	})
}

func TestAgent_GetStat(t *testing.T) {
	t.Run("Should: return stat about computer", func(t *testing.T) {
		a := New()
		assert.IsType(t, &agentPb.GetStatsResponse{}, a.GetStat())
	})
}