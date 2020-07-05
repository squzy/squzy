package expression

import (
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpressionStruct_GetSnapshots(t *testing.T) {
	t.Run("Should: panic", func(t *testing.T) {
		assert.Panics(t, func() { exprErr.GetSnapshots("id", 0, nil) }, "The code did not panic")
	})
	t.Run("Should: not panic", func(t *testing.T) {
		panicFunc := func() {
			exprCorr.GetSnapshots(
				"id",
				apiPb.SortDirection_ASC,
				nil,
				func(req *apiPb.GetSchedulerInformationRequest) *apiPb.GetSchedulerInformationRequest {
					return req
				})
		}
		assert.NotPanics(t, panicFunc, "The code did not panic")
	})
}

func TestExpressionStruct_getSnapshotEnv(t *testing.T) {
	t.Run("Should: no panic", func(t *testing.T) {
		res := exprCorr.ProcessRule(
			apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_SCHEDULER,
			"12345",
			`count(Last(10, UseTimeFrom("3/1/2020"), UseTimeTo("3/1/2021"), UseCode(Ok)), {.Type == GRPC}) == 1`)
		assert.True(t, res)
	})
	t.Run("Should: no panic", func(t *testing.T) {
		res := exprCorr.ProcessRule(
			apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_SCHEDULER,
			"12345",
			`count(First(10, UseTimeTo("3/1/2021"), UseCode(Ok)), {.Type == GRPC}) == 1`)
		assert.True(t, res)
	})
	t.Run("Should: no panic", func(t *testing.T) {
		res := exprCorr.ProcessRule(
			apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_SCHEDULER,
			"12345",
			`count(Index(1, UseTimeTo("3/1/2021"), UseCode(Ok)), {.Type == GRPC}) == 1`)
		assert.True(t, res)
	})
}
