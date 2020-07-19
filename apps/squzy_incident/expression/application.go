package expression

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

type FilterTransaction func(req *apiPb.GetTransactionsRequest) *apiPb.GetTransactionsRequest

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
	for _, filter := range filters {
		req = filter(req)
	}
	list, err := e.storageClient.GetTransactions(context.Background(), req)
	if err != nil {
		panic(err)
	}
	return list.GetTransactions()
}

func (e *expressionStruct) getTransactionEnv(applicationId string) map[string]interface{} {
	return map[string]interface{}{
		"Last": func(count int32, filters ...FilterTransaction) []*apiPb.TransactionInfo {
			return e.GetTransactions(
				applicationId,
				apiPb.SortDirection_ASC,
				&apiPb.Pagination{
					Page:  -1,
					Limit: count,
				},
				filters...)
		},
		"First": func(count int32, filters ...FilterTransaction) []*apiPb.TransactionInfo {
			return e.GetTransactions(
				applicationId,
				apiPb.SortDirection_ASC,
				&apiPb.Pagination{
					Page:  0,
					Limit: count,
				},
				filters...)
		},
		"Index": func(index int32, filters ...FilterTransaction) []*apiPb.TransactionInfo {
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
		"UseTimeFrom": func(timeStr string) FilterTransaction {
			return func(req *apiPb.GetTransactionsRequest) *apiPb.GetTransactionsRequest {
				if req.TimeRange == nil {
					req.TimeRange = &apiPb.TimeFilter{}
				}
				req.TimeRange.From = convertToTimestamp(timeStr)
				return req
			}
		},
		"UseTimeTo": func(timeStr string) FilterTransaction {
			return func(req *apiPb.GetTransactionsRequest) *apiPb.GetTransactionsRequest {
				if req.TimeRange == nil {
					req.TimeRange = &apiPb.TimeFilter{}
				}
				req.TimeRange.To = convertToTimestamp(timeStr)
				return req
			}
		},
		"UseHost": func(host string) FilterTransaction {
			return func(req *apiPb.GetTransactionsRequest) *apiPb.GetTransactionsRequest {
				req.Host = &wrappers.StringValue{
					Value: host,
				}
				return req
			}
		},
		"UseName": func(name string) FilterTransaction {
			return func(req *apiPb.GetTransactionsRequest) *apiPb.GetTransactionsRequest {
				req.Name = &wrappers.StringValue{
					Value: name,
				}
				return req
			}
		},
		"UsePath": func(path string) FilterTransaction {
			return func(req *apiPb.GetTransactionsRequest) *apiPb.GetTransactionsRequest {
				req.Path = &wrappers.StringValue{
					Value: path,
				}
				return req
			}
		},
		"UseMethod": func(method string) FilterTransaction {
			return func(req *apiPb.GetTransactionsRequest) *apiPb.GetTransactionsRequest {
				req.Method = &wrappers.StringValue{
					Value: method,
				}
				return req
			}
		},
		"Duration": func(snapshot *apiPb.TransactionInfo) int64 {
			return getTimeRange(snapshot.GetStartTime(), snapshot.GetEndTime())
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
}
