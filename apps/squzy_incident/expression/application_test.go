package expression

import (
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpressionStruct_GetTransactions(t *testing.T) {
	t.Run("Should: panic", func(t *testing.T) {
		assert.Panics(t, func() { exprErr.GetTransactions("id", 0, nil) }, "The code did not panic")
	})
	t.Run("Should: not panic", func(t *testing.T) {
		panicFunc := func() {
			exprCorr.GetTransactions(
				"id",
				apiPb.SortDirection_ASC,
				nil,
				func(req *apiPb.GetTransactionsRequest) *apiPb.GetTransactionsRequest {
					return req
				})
		}
		assert.NotPanics(t, panicFunc, "The code did not panic")
	})
}

func TestExpressionStruct_getTransactionEnv(t *testing.T) {
	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_APPLICATION,
			"12345",
			`count(Last(10, UseTimeFrom("3/1/2020"), UseTimeTo("3/1/2021")), {.Meta.Host == "host"}) == 1`)
		assert.True(t, res)
		assert.Nil(t, err)
	})
	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_APPLICATION,
			"12345",
			`count(First(10, UseTimeTo("3/1/2021")), {.Name == "name"}) == 1`)
		assert.True(t, res)
		assert.Nil(t, err)
	})
	//Duration
	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_APPLICATION,
			"12345",
			`count(First(10, UseTimeTo("3/1/2021")), {Duration(#) < 10}) == 1`)
		assert.True(t, res)
		assert.Nil(t, err)
	})
	t.Run("Should: no panic", func(t *testing.T) {
		res, err := exprCorr.ProcessRule(
			apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_APPLICATION,
			"12345",
			`len(Index(1, UseType(DB), UseStatus(Success), UseHost("host"), UseName("name"), UsePath("path"), UseMethod("method"))) == 1`)
		assert.True(t, res)
		assert.Nil(t, err)
	})
}
