package expression

import (
	"fmt"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_SCHEDULER,
			"12345",
			`len(Last(10, UseTimeFrom("3/1/2020"), UseTimeTo("3/1/2021"), UseCode(Ok))) == 1`)
		assert.True(t, res)
		assert.Nil(t, err)
	})
	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_SCHEDULER,
			"12345",
			`len(First(10, UseTimeTo("3/1/2021"), UseCode(Ok))) == 1`)
		assert.True(t, res)
		assert.Nil(t, err)
	})
	//Duration
	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_SCHEDULER,
			"12345",
			`count(First(10, UseTimeTo("3/1/2021"), UseCode(Ok)), {Duration(#) < 10}) == 1`)
		assert.True(t, res)
		assert.Nil(t, err)
	})
	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_SCHEDULER,
			"12345",
			`Index(1, UseTimeTo("3/1/2021"), UseCode(Ok)) != null`)
		assert.True(t, res)
		assert.Nil(t, err)
	})
	//UnixNanoNow
	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_SCHEDULER,
			"12345",
			`UnixNanoNow() - 1 > 0`)
		assert.True(t, res)
		assert.Nil(t, err)
	})
	// NowTime
	t.Run("Should: no panic", func(t *testing.T) {
		year := time.Now().Year()
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_SCHEDULER,
			"12345",
			fmt.Sprintf("NowTime().Year() == %d", year))
		assert.True(t, res)
		assert.Nil(t, err)
	})

	//timeDiff
	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_SCHEDULER,
			"12345",
			"durationToSecond(timeDiff(unixToTime(1611413752308 * 1000), unixToTime(1611413752308 * 1000))) == 0")
		assert.True(t, res)
		assert.Nil(t, err)
	})
	//timeDiff
	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_SCHEDULER,
			"12345",
			"durationToSecond(timeDiff(unixNanoToTime(1611413752308 * 1000), unixNanoToTime(1611413752308 * 1000))) == 0")
		assert.True(t, res)
		assert.Nil(t, err)
	})
	// float64ToInt64
	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_SCHEDULER,
			"12345",
			"float64ToInt64(0.0) == 0")
		assert.True(t, res)
		assert.Nil(t, err)
	})

	// getValue

	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_SCHEDULER,
			"12345",
			"getValue(Index(1)) != null")
		assert.True(t, res)
		assert.Nil(t, err)
	})

	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_SCHEDULER,
			"12345",
			"durationToSecond(Week) > 0 && durationToSecond(Day) > 0 && durationToSecond(Hour) > 0 && durationToSecond(Second) > 0 && durationToSecond(Minute) > 0")
		assert.True(t, res)
		assert.Nil(t, err)
	})


}
