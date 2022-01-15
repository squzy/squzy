package expression

import (
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpressionStruct_GetAgents(t *testing.T) {
	t.Run("Should: panic", func(t *testing.T) {
		assert.Panics(t, func() { exprErr.GetAgents("id", nil) }, "The code did not panic")
	})
	t.Run("Should: not panic", func(t *testing.T) {
		panicFunc := func() {
			exprCorr.GetAgents(
				"id",
				nil,
				func(req *apiPb.GetAgentInformationRequest) *apiPb.GetAgentInformationRequest {
					return req
				})
		}
		assert.NotPanics(t, panicFunc, "The code did not panic")
	})
}

func TestExpressionStruct_getAgentEnv(t *testing.T) {
	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_AGENT,
			"12345",
			`one(Last(10, UseTimeTo("3/1/2021 12:30"), UseType(CPU)), {one(.CpuInfo.Cpus, {.Load <= 10})})`)
		assert.True(t, res)
		assert.Nil(t, err)
	})
	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_AGENT,
			"12345",
			`one(Last(10, UseTimeFrom("3/1/2020"), UseTimeTo("3/1/2021"), UseType(All)), {one(.CpuInfo.Cpus, {.Load <= 10})})`)
		assert.True(t, res)
		assert.Nil(t, err)
	})
}
