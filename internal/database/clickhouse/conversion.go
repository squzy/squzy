package clickhouse

import (
	"bytes"
	"errors"
	"github.com/golang/protobuf/jsonpb"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"strconv"
	"strings"

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

func ConvertToSnapshot(request *apiPb.SchedulerResponse) (*Snapshot, error) {
	return convertToSnapshot(request.GetSnapshot(), request.GetSchedulerId())
}

func convertToSnapshot(request *apiPb.SchedulerSnapshot, schedulerID string) (*Snapshot, error) {
	if request == nil {
		return nil, errors.New("ERROR_SNAPSHOT_IS_EMPTY")
	}
	if request.GetMeta() == nil {
		return nil, errors.New("EMPTY_META_DATA")
	}
	startTime, err := ptypes.Timestamp(request.GetMeta().GetStartTime())
	if err != nil {
		return nil, err
	}
	endTime, err := ptypes.Timestamp(request.GetMeta().GetEndTime())
	if err != nil {
		return nil, err
	}

	res := &Snapshot{
		SchedulerID:   schedulerID,
		Code:          int32(request.GetCode()),
		Type:          int32(request.GetType()),
		MetaStartTime: startTime.UnixNano(),
		MetaEndTime:   endTime.UnixNano(),
	}
	if request.GetError() != nil {
		res.Error = request.GetError().GetMessage()
	}

	var b bytes.Buffer
	err = (&jsonpb.Marshaler{}).Marshal(&b, request.GetMeta().GetValue())
	if err != nil {
		return res, nil
	}
	res.MetaValue = b.Bytes()
	return res, nil
}

func ConvertFromSnapshots(snapshots []*Snapshot) []*apiPb.SchedulerSnapshot {
	var res []*apiPb.SchedulerSnapshot
	for _, v := range snapshots {
		snap, err := convertFromSnapshot(v)
		if err == nil {
			res = append(res, snap)
		}
		if err != nil {
			logger.Error(err.Error())
		}
	}
	return res
}

func convertFromSnapshot(snapshot *Snapshot) (*apiPb.SchedulerSnapshot, error) {
	//Skip error, because this convertion is always correct (snapshot.MetaStartTime < maximum possible value)
	startTime, _ := ptypes.TimestampProto(time.Unix(0, snapshot.MetaStartTime))
	endTime, _ := ptypes.TimestampProto(time.Unix(0, snapshot.MetaEndTime))

	res := &apiPb.SchedulerSnapshot{
		Code: apiPb.SchedulerCode(snapshot.Code),
		Type: apiPb.SchedulerType(snapshot.Type),
		Meta: &apiPb.SchedulerSnapshot_MetaData{
			StartTime: startTime,
			EndTime:   endTime,
		},
	}
	if snapshot.Error != "" {
		res.Error = &apiPb.SchedulerSnapshot_Error{
			Message: snapshot.Error,
		}
	}

	str := &_struct.Value{}
	if err := jsonpb.Unmarshal(bytes.NewReader(snapshot.MetaValue), str); err != nil {
		return res, nil
	}

	res.Meta.Value = str
	return res, nil
}

func convertFromUptimeResult(uptimeResult *UptimeResult, countAll int64) *apiPb.GetSchedulerUptimeResponse {
	latency, err := strconv.ParseFloat(strings.Split(uptimeResult.Latency, ".")[0], 64)
	if err != nil {
		return &apiPb.GetSchedulerUptimeResponse{
			Uptime:  0,
			Latency: 0,
		}
	}
	return &apiPb.GetSchedulerUptimeResponse{
		Uptime:  float64(uptimeResult.Count) / float64(countAll),
		Latency: latency,
	}
}