package expression

import (
	"context"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
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

	for _, filter := range filters {
		req = filter(req)
	}
	list, err := e.storageClient.GetAgentInformation(context.Background(), req)
	if err != nil {
		panic(err)
	}
	return list.GetStats()
}

func (e *expressionStruct) getAgentEnv(agentId string) map[string]interface{} {
	return map[string]interface{}{
		"Last": func(count int32, filters ...FilterAgent) []*apiPb.GetAgentInformationResponse_Statistic {
			return e.GetAgents(
				agentId,
				&apiPb.Pagination{
					Page:  -1,
					Limit: count,
				},
				filters...)
		},
		"UseTimeFrom": func(timeStr string) FilterAgent {
			return func(req *apiPb.GetAgentInformationRequest) *apiPb.GetAgentInformationRequest {
				if req.TimeRange == nil {
					req.TimeRange = &apiPb.TimeFilter{}
				}
				req.TimeRange.From = convertToTimestamp(timeStr)
				return req
			}
		},
		"UseTimeTo": func(timeStr string) FilterAgent {
			return func(req *apiPb.GetAgentInformationRequest) *apiPb.GetAgentInformationRequest {
				if req.TimeRange == nil {
					req.TimeRange = &apiPb.TimeFilter{}
				}
				req.TimeRange.To = convertToTimestamp(timeStr)
				return req
			}
		},
		"UseType": func(reqType apiPb.TypeAgentStat) FilterAgent {
			return func(req *apiPb.GetAgentInformationRequest) *apiPb.GetAgentInformationRequest {
				req.Type = reqType
				return req
			}
		},
		"All":    apiPb.TypeAgentStat_ALL,
		"CPU":    apiPb.TypeAgentStat_CPU,
		"Disk":   apiPb.TypeAgentStat_DISK,
		"Memory": apiPb.TypeAgentStat_MEMORY,
		"Net":    apiPb.TypeAgentStat_NET,
	}
}
