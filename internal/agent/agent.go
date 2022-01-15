package agent

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"sync"
	"time"
)

type Agent interface {
	GetStat() *apiPb.Metric
}

type agent struct {
	cpuStatFn           func(time.Duration, bool) ([]float64, error)
	swapMemoryStatFn    func() (*mem.SwapMemoryStat, error)
	virtualMemoryStatFn func() (*mem.VirtualMemoryStat, error)
	diskStatFn          func(bool) ([]disk.PartitionStat, error)
	diskUsageFn         func(string) (*disk.UsageStat, error)
	netStatFn           func(bool) ([]net.IOCountersStat, error)
	timeFn              func() *timestamp.Timestamp
}

func New(
	cpuStatFn func(time.Duration, bool) ([]float64, error),
	swapMemoryStatFn func() (*mem.SwapMemoryStat, error),
	virtualMemoryStatFn func() (*mem.VirtualMemoryStat, error),
	diskStatFn func(bool) ([]disk.PartitionStat, error),
	diskUsageFn func(string) (*disk.UsageStat, error),
	netStatFn func(bool) ([]net.IOCountersStat, error),
	timeFn func() *timestamp.Timestamp,
) *agent {
	return &agent{
		cpuStatFn:           cpuStatFn,
		swapMemoryStatFn:    swapMemoryStatFn,
		virtualMemoryStatFn: virtualMemoryStatFn,
		diskStatFn:          diskStatFn,
		diskUsageFn:         diskUsageFn,
		netStatFn:           netStatFn,
		timeFn:              timeFn,
	}
}

func (a *agent) GetStat() *apiPb.Metric {
	response := &apiPb.Metric{
		CpuInfo:    &apiPb.CpuInfo{},
		MemoryInfo: &apiPb.MemoryInfo{},
		DiskInfo:   &apiPb.DiskInfo{},
		NetInfo:    &apiPb.NetInfo{},
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		cpuStat, err := a.cpuStatFn(0, true)

		if err != nil || cpuStat == nil {
			return
		}

		for _, stat := range cpuStat {
			response.CpuInfo.Cpus = append(response.CpuInfo.Cpus, &apiPb.CpuInfo_CPU{
				Load: stat,
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		swapMemoryStat, err := a.swapMemoryStatFn()

		if err != nil || swapMemoryStat == nil {
			return
		}
		response.MemoryInfo.Swap = &apiPb.MemoryInfo_Memory{
			Total:       swapMemoryStat.Total,
			Used:        swapMemoryStat.Used,
			Free:        swapMemoryStat.Free,
			UsedPercent: swapMemoryStat.UsedPercent,
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		memoryStat, err := a.virtualMemoryStatFn()

		if err != nil || memoryStat == nil {
			return
		}

		response.MemoryInfo.Mem = &apiPb.MemoryInfo_Memory{
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

		disks, err := a.diskStatFn(false)

		if err != nil || disks == nil {
			return
		}
		diskStat := make(map[string]*apiPb.DiskInfo_Disk)
		for _, d := range disks {
			diskInfo, err := a.diskUsageFn(d.Mountpoint)
			if err != nil {
				continue
			}
			diskStat[d.Mountpoint] = &apiPb.DiskInfo_Disk{
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
		// take stat separate
		nets, err := a.netStatFn(true)

		if err != nil || nets == nil || len(nets) == 0 {
			return
		}

		netStat := make(map[string]*apiPb.NetInfo_Interface)

		for _, netInterface := range nets {
			netStat[netInterface.Name] = &apiPb.NetInfo_Interface{
				BytesSent:   netInterface.BytesSent,
				BytesRecv:   netInterface.BytesRecv,
				PacketsSent: netInterface.PacketsSent,
				PacketsRecv: netInterface.PacketsRecv,
				ErrIn:       netInterface.Errin,
				ErrOut:      netInterface.Errout,
				DropIn:      netInterface.Dropin,
				DropOut:     netInterface.Dropout,
			}
		}

		response.NetInfo.Interfaces = netStat
	}()

	wg.Wait()

	response.Time = a.timeFn()

	return response
}
