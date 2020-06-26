package expression

import (
	"context"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

type FilterAgent func(req *apiPb.GetAgentInformationRequest) *apiPb.GetAgentInformationRequest

func (e *expressionStruct) GetAgents(
	agentId string,
	pagination *apiPb.Pagination,
	filters ...FilterAgent) []*apiPb.GetAgentInformationResponse_Statistic {

	req := &apiPb.GetAgentInformationRequest{
		AgentId:    agentId,
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

func (e *expressionStruct) getAgentEnv(agentId string) map[string]interface{} {
	return map[string]interface{}{
		"last": func(count int32, filters ...FilterAgent) []*apiPb.GetAgentInformationResponse_Statistic {
			return e.GetAgents(
				agentId,
				&apiPb.Pagination{
					Page:  0,
					Limit: count,
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
}
