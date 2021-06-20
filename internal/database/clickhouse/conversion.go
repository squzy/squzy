package clickhouse

import (
	"bytes"
	"errors"
	"github.com/golang/protobuf/jsonpb"
	_struct "github.com/golang/protobuf/ptypes/struct"

	//nolint:staticcheck
	"github.com/golang/protobuf/ptypes"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"squzy/internal/logger"
	"time"
)

func convertToIncident(data *apiPb.Incident, t time.Time) *Incident {
	if data == nil {
		return nil
	}
	histories, startTime, endTime := convertToIncidentHistories(data.GetHistories())

	if startTime == 0 || endTime == 0 {
		startTime = t.UnixNano()
		endTime = t.UnixNano()
	}
	return &Incident{
		IncidentId: data.GetId(),
		Status:     int32(data.GetStatus()),
		RuleId:     data.GetRuleId(),
		StartTime:  startTime,
		EndTime:    endTime,
		Histories:  histories,
	}
}
func convertToIncidentHistories(data []*apiPb.Incident_HistoryItem) ([]*IncidentHistory, int64, int64) {
	if data == nil {
		return nil, 0, 0
	}
	var histories []*IncidentHistory
	minTime := int64(0)
	maxTime := int64(0)
	for _, v := range data {
		history := convertToIncidentHistory(v)
		if history != nil {
			histories = append(histories, history)
			if minTime < history.Timestamp || minTime == 0 {
				minTime = history.Timestamp
			}
			if maxTime > history.Timestamp || maxTime == 0 {
				maxTime = history.Timestamp
			}
		}
	}
	return histories, minTime, maxTime
}

func convertToIncidentHistory(data *apiPb.Incident_HistoryItem) *IncidentHistory {
	if data == nil {
		return nil
	}
	t, err := ptypes.Timestamp(data.GetTimestamp())
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return &IncidentHistory{
		Status:    int32(data.GetStatus()),
		Timestamp: t.UnixNano(),
	}
}

func convertFromIncidents(data []*Incident) []*apiPb.Incident {
	var incidents []*apiPb.Incident
	for _, v := range data {
		incidents = append(incidents, convertFromIncident(v))
	}
	return incidents
}

func convertFromIncident(data *Incident) *apiPb.Incident {
	return &apiPb.Incident{
		Id:        data.IncidentId,
		Status:    apiPb.IncidentStatus(data.Status),
		RuleId:    data.RuleId,
		Histories: convertFromIncidentHistories(data.Histories),
	}
}

func convertFromIncidentHistories(data []*IncidentHistory) []*apiPb.Incident_HistoryItem {
	var histories []*apiPb.Incident_HistoryItem
	for _, v := range data {
		history := convertFromIncidentHistory(v)
		if history != nil {
			histories = append(histories, history)
		}
	}
	return histories
}

func convertFromIncidentHistory(data *IncidentHistory) *apiPb.Incident_HistoryItem {
	if data == nil {
		return nil
	}
	parsedTime, _ := ptypes.TimestampProto(time.Unix(0, data.Timestamp))
	return &apiPb.Incident_HistoryItem{
		Status:    apiPb.IncidentStatus(data.Status),
		Timestamp: parsedTime,
	}
}
