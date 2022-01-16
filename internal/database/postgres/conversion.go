package postgres

import (
	"errors"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"google.golang.org/protobuf/types/known/structpb"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"strings"
	"time"
)

func ConvertToPostgresSnapshot(request *apiPb.SchedulerResponse) (*Snapshot, error) {
	return convertToSnapshot(request.GetSnapshot(), request.GetSchedulerId())
}

func ConvertFromPostgresSnapshots(snapshots []*Snapshot) []*apiPb.SchedulerSnapshot {
	var res []*apiPb.SchedulerSnapshot
	for _, v := range snapshots {
		snap, err := convertFromSnapshot(v)
		if err == nil {
			res = append(res, snap)
		}
		//TODO: log if error
	}
	return res
}

func ConvertToPostgressStatRequest(request *apiPb.Metric) (*StatRequest, error) {
	t := request.GetTime().AsTime()
	err := request.GetTime().CheckValid()
	if err != nil {
		return nil, err
	}
	return &StatRequest{
		AgentID:    request.GetAgentId(),
		AgentName:  request.GetAgentName(),
		CPUInfo:    convertToCPUInfo(request.GetCpuInfo()),
		MemoryInfo: convertToMemoryInfo(request.GetMemoryInfo()),
		DiskInfo:   convertToDiskInfo(request.GetDiskInfo()),
		NetInfo:    convertToNetInfo(request.GetNetInfo()),
		Time:       t,
	}, nil
}

func ConvertFromPostgressStatRequests(data []*StatRequest) []*apiPb.GetAgentInformationResponse_Statistic {
	var res []*apiPb.GetAgentInformationResponse_Statistic
	for _, request := range data {
		stat, err := ConvertFromPostgressStatRequest(request)
		if err == nil {
			res = append(res, stat)
		}
		//TODO: log if error
	}
	return res
}

func ConvertFromPostgressStatRequest(data *StatRequest) (*apiPb.GetAgentInformationResponse_Statistic, error) {
	t := timestamp.New(data.Time)
	err := t.CheckValid()
	if err != nil {
		return nil, err
	}
	return &apiPb.GetAgentInformationResponse_Statistic{
		CpuInfo:    convertFromCPUInfo(data.CPUInfo),
		MemoryInfo: convertFromMemoryInfo(data.MemoryInfo),
		DiskInfo:   convertFromDiskInfo(data.DiskInfo),
		NetInfo:    convertFromNetInfo(data.NetInfo),
		Time:       t,
	}, nil
}

func convertToSnapshot(request *apiPb.SchedulerSnapshot, schedulerID string) (*Snapshot, error) {
	if request == nil {
		return nil, errors.New("ERROR_SNAPSHOT_IS_EMPTY")
	}
	if request.GetMeta() == nil {
		return nil, errors.New("EMPTY_META_DATA")
	}
	startTime := request.GetMeta().GetStartTime().AsTime()
	err := request.GetMeta().GetStartTime().CheckValid()
	if err != nil {
		return nil, err
	}
	endTime := request.GetMeta().GetEndTime().AsTime()
	err = request.GetMeta().GetEndTime().CheckValid()
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

	bValue, err := request.GetMeta().GetValue().MarshalJSON()
	if err != nil {
		return res, nil
	}
	res.MetaValue = bValue
	return res, nil
}

func convertFromSnapshot(snapshot *Snapshot) (*apiPb.SchedulerSnapshot, error) {
	//Skip error, because this convertion is always correct (snapshot.MetaStartTime < maximum possible value)
	startTime := timestamp.New(time.Unix(0, snapshot.MetaStartTime))
	endTime := timestamp.New(time.Unix(0, snapshot.MetaEndTime))

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

	str := &structpb.Value{}

	if err := str.UnmarshalJSON(snapshot.MetaValue); err != nil {
		return res, nil
	}

	res.Meta.Value = str
	return res, nil
}

func convertToCPUInfo(request *apiPb.CpuInfo) []*CPUInfo {
	var res []*CPUInfo
	if request == nil {
		return res
	}
	for _, v := range request.Cpus {
		res = append(res, &CPUInfo{Load: v.GetLoad()})
	}
	return res
}

func convertToMemoryInfo(reqest *apiPb.MemoryInfo) *MemoryInfo {
	if reqest == nil {
		return nil
	}
	res := &MemoryInfo{}
	if reqest.GetMem() != nil {
		res.Mem = &MemoryMem{
			Total:       reqest.GetMem().GetTotal(),
			Used:        reqest.GetMem().GetUsed(),
			Free:        reqest.GetMem().GetFree(),
			Shared:      reqest.GetMem().GetShared(),
			UsedPercent: reqest.GetMem().GetUsedPercent(),
		}
	}
	if reqest.GetSwap() != nil {
		res.Swap = &MemorySwap{
			Total:       reqest.GetSwap().GetTotal(),
			Used:        reqest.GetSwap().GetUsed(),
			Free:        reqest.GetSwap().GetFree(),
			Shared:      reqest.GetSwap().GetShared(),
			UsedPercent: reqest.GetSwap().GetUsedPercent(),
		}
	}
	return res
}

func convertToDiskInfo(request *apiPb.DiskInfo) []*DiskInfo {
	var res []*DiskInfo
	if request == nil {
		return res
	}
	for name, v := range request.GetDisks() {
		res = append(res, &DiskInfo{
			Name:        name,
			Total:       v.GetTotal(),
			Free:        v.GetFree(),
			Used:        v.GetUsed(),
			UsedPercent: v.GetUsedPercent(),
		})
	}
	return res
}

func convertToNetInfo(request *apiPb.NetInfo) []*NetInfo {
	var res []*NetInfo
	if request == nil {
		return res
	}
	for name, v := range request.GetInterfaces() {
		res = append(res, &NetInfo{
			Name:        name,
			BytesSent:   v.GetBytesSent(),
			BytesRecv:   v.GetBytesRecv(),
			PacketsSent: v.GetPacketsSent(),
			PacketsRecv: v.GetPacketsRecv(),
			ErrIn:       v.GetErrIn(),
			ErrOut:      v.GetErrOut(),
			DropIn:      v.GetDropIn(),
			DropOut:     v.GetDropOut(),
		})
	}
	return res
}

func convertFromCPUInfo(data []*CPUInfo) *apiPb.CpuInfo {
	var cpus []*apiPb.CpuInfo_CPU
	for _, v := range data {
		cpus = append(cpus, &apiPb.CpuInfo_CPU{
			Load: v.Load,
		})
	}
	if len(cpus) == 0 {
		return nil
	}
	return &apiPb.CpuInfo{Cpus: cpus}
}

func convertFromMemoryInfo(data *MemoryInfo) *apiPb.MemoryInfo {
	if data == nil {
		return nil
	}
	res := &apiPb.MemoryInfo{
		Mem:  nil,
		Swap: nil,
	}
	if data.Mem != nil {
		res.Mem = &apiPb.MemoryInfo_Memory{
			Total:       data.Mem.Total,
			Used:        data.Mem.Used,
			Free:        data.Mem.Free,
			Shared:      data.Mem.Shared,
			UsedPercent: data.Mem.UsedPercent,
		}
	}
	if data.Swap != nil {
		res.Swap = &apiPb.MemoryInfo_Memory{
			Total:       data.Swap.Total,
			Used:        data.Swap.Used,
			Free:        data.Swap.Free,
			Shared:      data.Swap.Shared,
			UsedPercent: data.Swap.UsedPercent,
		}
	}
	if res.Mem == nil && res.Swap == nil {
		return nil
	}
	return res
}

func convertFromDiskInfo(data []*DiskInfo) *apiPb.DiskInfo {
	disks := map[string]*apiPb.DiskInfo_Disk{}
	for _, v := range data {
		disks[v.Name] = &apiPb.DiskInfo_Disk{
			Total:       v.Total,
			Free:        v.Free,
			Used:        v.Used,
			UsedPercent: v.UsedPercent,
		}
	}
	if len(disks) == 0 {
		return nil
	}
	return &apiPb.DiskInfo{Disks: disks}
}

func convertFromNetInfo(data []*NetInfo) *apiPb.NetInfo {
	interfaces := map[string]*apiPb.NetInfo_Interface{}
	for _, v := range data {
		interfaces[v.Name] = &apiPb.NetInfo_Interface{
			BytesSent:   v.BytesSent,
			BytesRecv:   v.BytesRecv,
			PacketsSent: v.PacketsSent,
			PacketsRecv: v.PacketsRecv,
			ErrIn:       v.ErrIn,
			ErrOut:      v.ErrOut,
			DropIn:      v.DropIn,
			DropOut:     v.DropOut,
		}
	}
	if len(interfaces) == 0 {
		return nil
	}
	return &apiPb.NetInfo{Interfaces: interfaces}
}

func convertToTransactionInfo(data *apiPb.TransactionInfo) (*TransactionInfo, error) {
	startTime := data.GetStartTime().AsTime()
	err := data.GetStartTime().CheckValid()
	if err != nil {
		return nil, err
	}
	endTime := data.GetEndTime().AsTime()
	err = data.GetEndTime().CheckValid()
	if err != nil {
		return nil, err
	}
	if data.GetMeta() == nil {
		data.Meta = &apiPb.TransactionInfo_Meta{
			Host:   "",
			Path:   "",
			Method: "",
		}
	}
	if data.Error == nil {
		data.Error = &apiPb.TransactionInfo_Error{
			Message: "",
		}
	}
	return &TransactionInfo{
		TransactionId:     data.GetId(),
		ApplicationId:     data.GetApplicationId(),
		ParentId:          data.GetParentId(),
		MetaHost:          data.GetMeta().GetHost(),
		MetaPath:          data.GetMeta().GetPath(),
		MetaMethod:        data.GetMeta().GetMethod(),
		Name:              data.GetName(),
		StartTime:         startTime.UnixNano(),
		EndTime:           endTime.UnixNano(),
		TransactionStatus: int32(data.GetStatus()),
		TransactionType:   int32(data.GetType()),
		Error:             data.GetError().GetMessage(),
	}, nil
}

func convertFromTransaction(data *TransactionInfo) *apiPb.TransactionInfo {
	//Skip error, because this convertion is always correct (data.StartTime < maximum possible value)
	startTime := timestamp.New(time.Unix(0, data.StartTime))
	endTime := timestamp.New(time.Unix(0, data.EndTime))

	transactionMeta := &apiPb.TransactionInfo_Meta{
		Host:   data.MetaHost,
		Path:   data.MetaPath,
		Method: data.MetaMethod,
	}
	if data.MetaPath == "" && data.MetaHost == "" && data.MetaMethod == "" {
		transactionMeta = nil
	}
	transactionError := &apiPb.TransactionInfo_Error{
		Message: data.Error,
	}
	if data.Error == "" {
		transactionError = nil
	}
	return &apiPb.TransactionInfo{

		Id:            data.TransactionId,
		ApplicationId: data.ApplicationId,
		ParentId:      data.ParentId,
		Meta:          transactionMeta,
		Name:          data.Name,
		StartTime:     startTime,
		EndTime:       endTime,
		Status:        apiPb.TransactionStatus(data.TransactionStatus),
		Type:          apiPb.TransactionType(data.TransactionType),
		Error:         transactionError,
	}
}

func convertFromTransactions(data []*TransactionInfo) []*apiPb.TransactionInfo {
	var res []*apiPb.TransactionInfo
	for _, request := range data {
		stat := convertFromTransaction(request)
		if stat != nil {
			res = append(res, stat)
		}
	}
	return res
}

func convertToIncident(data *apiPb.Incident) *Incident {
	histories, startTime, endTime := convertToIncidentHistories(data.GetHistories())
	if startTime == 0 || endTime == 0 {
		startTime = time.Now().UnixNano()
		endTime = time.Now().UnixNano()
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
	timeParsed := data.GetTimestamp().AsTime()
	err := data.GetTimestamp().CheckValid()
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return &IncidentHistory{
		Status:    int32(data.GetStatus()),
		Timestamp: timeParsed.UnixNano(),
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
	histories := []*apiPb.Incident_HistoryItem{}
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
	parsedTime := timestamp.New(time.Unix(0, data.Timestamp))
	return &apiPb.Incident_HistoryItem{
		Status:    apiPb.IncidentStatus(data.Status),
		Timestamp: parsedTime,
	}
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

func convertFromGroupResult(group []*GroupResult, upTime int64) map[string]*apiPb.TransactionGroup {
	res := map[string]*apiPb.TransactionGroup{}
	for _, v := range group {
		latency, err := strconv.ParseFloat(strings.Split(v.Latency, ".")[0], 64)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		minTime, err := strconv.ParseFloat(strings.Split(v.MinTime, ".")[0], 64)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		maxTime, err := strconv.ParseFloat(strings.Split(v.MaxTime, ".")[0], 64)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		lowTime, err := strconv.ParseFloat(strings.Split(v.LowTime, ".")[0], 64)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		res[v.Name] = &apiPb.TransactionGroup{
			Count:        v.Count,
			SuccessRatio: float64(v.SuccessCount) / float64(v.Count),
			AverageTime:  latency / 1000000, //div 1000 in order to get milliseconds
			MinTime:      minTime / 1000000,
			MaxTime:      maxTime / 1000000,
			Throughput:   getThroughput(v.Count, lowTime, upTime), //Take minutes and div by count
		}
	}
	return res
}

func getThroughput(count int64, lowTime float64, upTime int64) float64 {
	timeDiapasonMinutes := (float64(upTime) - lowTime) / 60000000000
	if timeDiapasonMinutes == 0 {
		return 0
	}
	return float64(count) / timeDiapasonMinutes
}
