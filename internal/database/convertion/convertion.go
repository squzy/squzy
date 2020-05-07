package convertion

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"squzy/internal/database"
	"strconv"
)

func ConvertToPostgressScheduler(request *apiPb.SchedulerResponse) (*database.Snapshot, error) {
	id, err := strconv.ParseUint(request.GetSchedulerId(), 10, 32)
	if err != nil {
		return nil, err
	}
	return convertToSnapshot(request.GetSnapshot(), uint(id))
}

func ConvertFromPostgressSnapshots(snapshots []*database.Snapshot) ([]*apiPb.Snapshot, []error) {
	var errors []error
	var res []*apiPb.Snapshot
	for _, v := range snapshots {
		snap, err := convertFromSnapshot(v)
		if err != nil {
			errors = append(errors, err)
		} else {
			res = append(res, snap)
		}
	}
	return res, errors
}

func ConvertToPostgressStatRequest(request *apiPb.SendMetricsRequest) (*database.StatRequest, error) {
	t, err := ptypes.Timestamp(request.GetTime())
	if err != nil {
		return nil, err
	}
	return &database.StatRequest{
		CpuInfo:    convertToCpuInfo(request.GetCpuInfo()),
		MemoryInfo: convertToMemoryInfo(request.GetMemoryInfo()),
		DiskInfo:   convertToDiskInfo(request.GetDiskInfo()),
		NetInfo:    convertToNetInfo(request.GetNetInfo()),
		Time:       t,
	}, nil
}

func ConvertFromPostgressStatRequest(data *database.StatRequest) (*apiPb.SendMetricsRequest, error) {
	t, err := ptypes.TimestampProto(data.Time)
	if err != nil {
		return nil, err
	}
	return &apiPb.SendMetricsRequest{
		AgentId:       fmt.Sprint(data.ID),
		AgentUniqName: "", //TODO
		CpuInfo:       convertFromCpuInfo(data.CpuInfo),
		MemoryInfo:    convertFromMemoryInfo(data.MemoryInfo),
		DiskInfo:      convertFromDiskInfo(data.DiskInfo),
		NetInfo:       convertFromNetInfo(data.NetInfo),
		Time:          t,
	}, nil
}

func convertToSnapshot(request *apiPb.Snapshot, schedulerId uint) (*database.Snapshot, error) {
	if request == nil {
		return nil, errors.New("ERROR_SNAPSHOT_IS_EMPTY")
	}
	metaData, err := convertToMetaData(request.GetMeta())
	if err != nil {
		return nil, err
	}
	res := &database.Snapshot{
		SchedulerId: schedulerId,
		Code:        request.GetCode().String(),
		Type:        request.GetType().String(),
		Meta:        metaData,
	}
	if request.GetError() != nil {
		res.Error = request.GetError().GetMessage()
	}
	return res, nil
}

func convertToMetaData(request *apiPb.Snapshot_MetaData) (*database.MetaData, error) {
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
	return &database.MetaData{
		StartTime: startTime,
		EndTime:   endTime,
		Value:     request.GetValue(),
	}, nil
}

func convertFromSnapshot(snapshot *database.Snapshot) (*apiPb.Snapshot, error) {
	meta, err := convertFromMetaData(snapshot.Meta)
	if err != nil {
		return nil, err
	}
	return &apiPb.Snapshot{
		Code:                 apiPb.Snapshot_Code(apiPb.SchedulerResponseCode_value[snapshot.Code]),
		Type:                 apiPb.SchedulerType(apiPb.SchedulerType_value[snapshot.Type]),
		Error:                &apiPb.Snapshot_SnapshotError{
			Message: snapshot.Error,
		},
		Meta:                 meta,
	}, nil
}

func convertFromMetaData(metaData *database.MetaData) (*apiPb.Snapshot_MetaData, error) {
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
	return &apiPb.Snapshot_MetaData{
		StartTime:            startTime,
		EndTime:              endTime,
		Value:                metaData.Value,
	}, nil
}

func convertToCpuInfo(request *apiPb.CpuInfo) []*database.CpuInfo {
	var res []*database.CpuInfo
	if request == nil {
		return res
	}
	for _, v := range request.Cpus {
		res = append(res, &database.CpuInfo{Load: v.GetLoad()})
	}
	return res
}

func convertToMemoryInfo(reqest *apiPb.MemoryInfo) *database.MemoryInfo {
	if reqest == nil {
		return nil
	}
	res := &database.MemoryInfo{}
	if reqest.GetMem() != nil {
		res.Mem = &database.Memory{
			Total:       reqest.GetMem().GetTotal(),
			Used:        reqest.GetMem().GetUsed(),
			Free:        reqest.GetMem().GetFree(),
			Shared:      reqest.GetMem().GetShared(),
			UsedPercent: reqest.GetMem().GetUsedPercent(),
		}
	}
	if reqest.GetSwap() != nil {
		res.Swap = &database.Memory{
			Total:       reqest.GetSwap().GetTotal(),
			Used:        reqest.GetSwap().GetUsed(),
			Free:        reqest.GetSwap().GetFree(),
			Shared:      reqest.GetSwap().GetShared(),
			UsedPercent: reqest.GetSwap().GetUsedPercent(),
		}
	}
	return res
}

func convertToDiskInfo(request *apiPb.DiskInfo) []*database.DiskInfo {
	var res []*database.DiskInfo
	if request == nil {
		return res
	}
	for name, v := range request.GetDisks() {
		res = append(res, &database.DiskInfo{
			Name:        name,
			Total:       v.GetTotal(),
			Free:        v.GetFree(),
			Used:        v.GetUsed(),
			UsedPercent: v.GetUsedPercent(),
		})
	}
	return res
}

func convertToNetInfo(request *apiPb.NetInfo) []*database.NetInfo {
	var res []*database.NetInfo
	if request == nil {
		return res
	}
	for name, v := range request.GetInterfaces() {
		res = append(res, &database.NetInfo{
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

func convertFromCpuInfo(data []*database.CpuInfo) *apiPb.CpuInfo {
	var cpus []*apiPb.CpuInfo_CPU
	for _, v := range data {
		cpus = append(cpus, &apiPb.CpuInfo_CPU{
			Load: v.Load,
		})
	}
	return &apiPb.CpuInfo{Cpus: cpus,}
}

func convertFromMemoryInfo(data *database.MemoryInfo) *apiPb.MemoryInfo {
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
	return res
}

func convertFromDiskInfo(data []*database.DiskInfo) *apiPb.DiskInfo {
	disks := map[string]*apiPb.DiskInfo_Disk{}
	for _, v := range data {
		disks[v.Name] = &apiPb.DiskInfo_Disk{
			Total:       v.Total,
			Free:        v.Free,
			Used:        v.Used,
			UsedPercent: v.UsedPercent,
		}
	}
	return &apiPb.DiskInfo{Disks: disks,}
}

func convertFromNetInfo(data []*database.NetInfo) *apiPb.NetInfo {
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
	return &apiPb.NetInfo{Interfaces: interfaces,}
}
