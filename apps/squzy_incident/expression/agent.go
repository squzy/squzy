package expression

import (
	"context"
	"fmt"
	"github.com/antonmedv/expr"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"strconv"
)

type FilterAgent func(req *apiPb.GetAgentInformationRequest) *apiPb.GetAgentInformationRequest

func (e *expressionStruct) IsValidAgent(agentId string, rule string) bool {
	env := map[string]interface{}{
		"last": func(count int32, filters ...FilterAgent) []*apiPb.GetAgentInformationResponse_Statistic {
			return e.GetAgents(
				agentId,
				&apiPb.Pagination{
					Page:                 0,
					Limit:                count,
				},
				filters...)
		},
		"SetTimeFrom": func(timeStr string) FilterAgent {
			return func(req *apiPb.GetAgentInformationRequest) *apiPb.GetAgentInformationRequest {
				if req.TimeRange == nil {
					req.TimeRange = &apiPb.TimeFilter{}
				}
				req.TimeRange.From = convertToTimestamp(timeStr)
				return req
			}
		},
		"SetTimeTo": func(timeStr string) FilterAgent {
			return func(req *apiPb.GetAgentInformationRequest) *apiPb.GetAgentInformationRequest {
				if req.TimeRange == nil {
					req.TimeRange = &apiPb.TimeFilter{}
				}
				req.TimeRange.To = convertToTimestamp(timeStr)
				return req
			}
		},
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

func (e *expressionStruct) GetAgents(
	agentId string,
	pagination *apiPb.Pagination,
	filters ...FilterAgent) []*apiPb.GetAgentInformationResponse_Statistic {

	req := &apiPb.GetAgentInformationRequest{
		AgentId: agentId,
		Pagination: pagination,
	}
	if filters != nil {
		for _, filter := range filters {
			req = filter(req)
		}
	}
	list, err := e.storageClient.GetAgentInformation(context.Background(), req)
	if err != nil {
		panic(err)
	}
	return list.GetStats()
}
