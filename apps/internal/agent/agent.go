package agent

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	agentPb "github.com/squzy/squzy_generated/generated/agent/proto/v1"
	"sync"
	"time"
)

const (
	cpuInterval = time.Millisecond * 500
)

type agent struct {
}

func New() *agent {
	return &agent{}
}

func (a *agent) GetStat() *agentPb.GetStatsResponse {
	response := &agentPb.GetStatsResponse{}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		cpuStat, err := cpu.Percent(cpuInterval, true)

		if err == nil || cpuStat == nil {
			return
		}
		for _, stat := range cpuStat {
			response.CpuInfo.Cpus = append(response.CpuInfo.Cpus, &agentPb.CpuInfo_CPU{
				Load: stat,
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		swapMemoryStat, err := mem.SwapMemory()

		if err == nil || swapMemoryStat == nil {
			return
		}

		response.MemoryInfo.Swap = &agentPb.MemoryInfo_Memory{
			Total:       swapMemoryStat.Total,
			Used:        swapMemoryStat.Used,
			Free:        swapMemoryStat.Free,
			UsedPercent: swapMemoryStat.UsedPercent,
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		memoryStat, err := mem.VirtualMemory()

		if err == nil || memoryStat == nil {
			return
		}

		response.MemoryInfo.Mem = &agentPb.MemoryInfo_Memory{
			Total:       memoryStat.Total,
			Used:        memoryStat.Used,
			Free:        memoryStat.Free,
			Shared:      memoryStat.Shared,
			UsedPercent: memoryStat.UsedPercent,
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		disks, err := disk.Partitions(false)

		if err == nil || disks == nil {
			return
		}
		diskStat := make(map[string]*agentPb.DiskInfo_Disk)
		for _, d := range disks {
			diskInfo, err := disk.Usage(d.Mountpoint)
			if err != nil {
				continue
			}
			diskStat[d.Mountpoint] = &agentPb.DiskInfo_Disk{
				Total:       diskInfo.Total,
				Free:        diskInfo.Free,
				Used:        diskInfo.Used,
				UsedPercent: diskInfo.UsedPercent,
			}
		}

		response.DiskInfo.Disks = diskStat
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		netStat, err := net.IOCounters(false)
		if err == nil || netStat == nil || len(netStat) == 0 {
			return
		}

		stat := netStat[0]
		response.NetInfo = &agentPb.NetInfo{
			BytesSent:   stat.BytesSent,
			BytesRecv:   stat.BytesRecv,
			PacketsSent: stat.PacketsSent,
			PacketsRecv: stat.PacketsRecv,
			ErrIn:       stat.Errin,
			ErrOut:      stat.Errout,
			DropIn:      stat.Dropin,
			DropOut:     stat.Dropout,
		}
	}()

	wg.Wait()

	return response
}
