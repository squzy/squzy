package expression

import (
	"context"
	"fmt"
	"github.com/antonmedv/expr"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"strconv"
)

type FilterTransaction func(req *apiPb.GetTransactionsRequest) *apiPb.GetTransactionsRequest

func (e *expressionStruct) IsValidTransaction(applicationId string, rule string) bool {
	env := map[string]interface{}{
		"last": func(count int32, filters ...FilterTransaction) []*apiPb.TransactionInfo {
			return e.GetTransactions(
				applicationId,
				apiPb.SortDirection_DESC,
				&apiPb.Pagination{
					Page:  0,
					Limit: count,
				},
				filters...)
		},
		"first": func(count int32, filters ...FilterTransaction) []*apiPb.TransactionInfo {
			return e.GetTransactions(
				applicationId,
				apiPb.SortDirection_ASC,
				&apiPb.Pagination{
					Page:  0,
					Limit: count,
				},
				filters...)
		},
		"index": func(index int32, filters ...FilterTransaction) []*apiPb.TransactionInfo {
			return e.GetTransactions(
				applicationId,
				apiPb.SortDirection_ASC,
				&apiPb.Pagination{
					Page:  index,
					Limit: 1,
				},
				filters...)
		},
		"UseType": func(trType apiPb.TransactionType) FilterTransaction {
			return func(req *apiPb.GetTransactionsRequest) *apiPb.GetTransactionsRequest {
				req.Type = trType
				return req
			}
		},
		"UseStatus": func(trStatus apiPb.TransactionStatus) FilterTransaction {
			return func(req *apiPb.GetTransactionsRequest) *apiPb.GetTransactionsRequest {
				req.Status = trStatus
				return req
			}
		},
		"SetTimeFrom": func(timeStr string) FilterTransaction {
			return func(req *apiPb.GetTransactionsRequest) *apiPb.GetTransactionsRequest {
				if req.TimeRange == nil {
					req.TimeRange = &apiPb.TimeFilter{}
				}
				req.TimeRange.From = convertToTimestamp(timeStr)
				return req
			}
		},
		"SetTimeTo": func(timeStr string) FilterTransaction {
			return func(req *apiPb.GetTransactionsRequest) *apiPb.GetTransactionsRequest {
				if req.TimeRange == nil {
					req.TimeRange = &apiPb.TimeFilter{}
				}
				req.TimeRange.To = convertToTimestamp(timeStr)
				return req
			}
		},
		//Transaction status keys
		"Success": apiPb.TransactionStatus_TRANSACTION_SUCCESSFUL,
		"Failed":  apiPb.TransactionStatus_TRANSACTION_FAILED,
		//Transaction type keys
		"Xhr":       apiPb.TransactionType_TRANSACTION_TYPE_XHR,
		"Fetch":     apiPb.TransactionType_TRANSACTION_TYPE_FETCH,
		"Websocket": apiPb.TransactionType_TRANSACTION_TYPE_WEBSOCKET,
		"HTTP":      apiPb.TransactionType_TRANSACTION_TYPE_HTTP,
		"GRPC":      apiPb.TransactionType_TRANSACTION_TYPE_GRPC,
		"DB":        apiPb.TransactionType_TRANSACTION_TYPE_DB,
		"Internal":  apiPb.TransactionType_TRANSACTION_TYPE_INTERNAL,
		"Router":    apiPb.TransactionType_TRANSACTION_TYPE_ROUTER,
	}

	program, err := expr.Compile(rule, expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	value, err := strconv.ParseBool(fmt.Sprintf("%v", output))
	if err == nil {
		return value
	}
	return false
}

func (e *expressionStruct) GetTransactions(
	applicationId string,
	direction apiPb.SortDirection,
	pagination *apiPb.Pagination,
	filters ...FilterTransaction) []*apiPb.TransactionInfo {

	req := &apiPb.GetTransactionsRequest{
		ApplicationId: applicationId,
		Pagination:    pagination,
		Sort: &apiPb.SortingTransactionList{
			Direction: direction,
		},
	}
	if filters != nil {
		for _, filter := range filters {
			req = filter(req)
		}
	}
	list, err := e.storageClient.GetTransactions(context.Background(), req)
	if err != nil {
		panic(err)
	}
	return list.GetTransactions()
}
