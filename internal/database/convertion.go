package database

import (
	"errors"
	"github.com/golang/protobuf/ptypes"
	_struct "github.com/golang/protobuf/ptypes/struct"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"google.golang.org/protobuf/encoding/protojson"
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
	t, err := ptypes.Timestamp(request.GetTime())
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
	t, err := ptypes.TimestampProto(data.Time)
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
	metaData, err := convertToMetaData(request.GetMeta())
	if err != nil {
		return nil, err
	}
	res := &Snapshot{
		SchedulerID: schedulerID,
		Code:        request.GetCode().String(),
		Type:        request.GetType().String(),
		Meta:        metaData,
	}
	if request.GetError() != nil {
		res.Error = request.GetError().GetMessage()
	}
	return res, nil
}

func convertToMetaData(request *apiPb.SchedulerSnapshot_MetaData) (*MetaData, error) {
	if request == nil {
		return nil, errors.New("EMPTY_META_DATA")
	}
	startTime, err := ptypes.Timestamp(request.GetStartTime())
	if err != nil {
		return nil, err
	}
	endTime, err := ptypes.Timestamp(request.GetEndTime())
	if err != nil {
		return nil, err
	}

	buffer, err := protojson.Marshal(request.GetValue())
	if err != nil {
		return &MetaData{
			StartTime: startTime,
			EndTime:   endTime,
		}, nil
	}
	return &MetaData{
		StartTime: startTime,
		EndTime:   endTime,
		Value:     buffer,
	}, nil
}

func convertFromSnapshot(snapshot *Snapshot) (*apiPb.SchedulerSnapshot, error) {
	meta, err := convertFromMetaData(snapshot.Meta)
	if err != nil {
		return nil, err
	}
	return &apiPb.SchedulerSnapshot{
		Code: apiPb.SchedulerCode(apiPb.SchedulerCode_value[snapshot.Code]),
		Type: apiPb.SchedulerType(apiPb.SchedulerType_value[snapshot.Type]),
		Error: &apiPb.SchedulerSnapshot_Error{
			Message: snapshot.Error,
		},
		Meta: meta,
	}, nil
}

func convertFromMetaData(metaData *MetaData) (*apiPb.SchedulerSnapshot_MetaData, error) {
	if metaData == nil {
		return nil, errors.New("EMPTY_META_DATA")
	}
	startTime, err := ptypes.TimestampProto(metaData.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := ptypes.TimestampProto(metaData.EndTime)
	if err != nil {
		return nil, err
	}
	str := &_struct.Value{}

	if err := protojson.Unmarshal(metaData.Value, str); err != nil {
		return &apiPb.SchedulerSnapshot_MetaData{
			StartTime: startTime,
			EndTime:   endTime,
		}, nil
	}
	return &apiPb.SchedulerSnapshot_MetaData{
		StartTime: startTime,
		EndTime:   endTime,
		Value:     str,
	}, nil
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
