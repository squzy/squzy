package agent

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	agentPb "github.com/squzy/squzy_generated/generated/agent/proto/v1"
	"sync"
	"time"
)

const (
	cpuInterval = time.Millisecond * 500
)

type Agent interface {
	GetStat() *agentPb.SendStat
}

type agent struct {
	cpuStatFn           func(time.Duration, bool) ([]float64, error)
	swapMemoryStatFn    func() (*mem.SwapMemoryStat, error)
	virtualMemoryStatFn func() (*mem.VirtualMemoryStat, error)
	diskStatFn          func(bool) ([]disk.PartitionStat, error)
	diskUsageFn         func(string) (*disk.UsageStat, error)
	netStatFn           func(bool) ([]net.IOCountersStat, error)
	hostStatFn          func() (*host.InfoStat, error)
	timeFn              func() *timestamp.Timestamp
}

func New(
	cpuStatFn func(time.Duration, bool) ([]float64, error),
	swapMemoryStatFn func() (*mem.SwapMemoryStat, error),
	virtualMemoryStatFn func() (*mem.VirtualMemoryStat, error),
	diskStatFn func(bool) ([]disk.PartitionStat, error),
	diskUsageFn func(string) (*disk.UsageStat, error),
	netStatFn func(bool) ([]net.IOCountersStat, error),
	hostStatFn func() (*host.InfoStat, error),
	timeFn func() *timestamp.Timestamp,
) *agent {
	return &agent{
		cpuStatFn:           cpuStatFn,
		swapMemoryStatFn:    swapMemoryStatFn,
		virtualMemoryStatFn: virtualMemoryStatFn,
		diskStatFn:          diskStatFn,
		diskUsageFn:         diskUsageFn,
		netStatFn:           netStatFn,
		hostStatFn:          hostStatFn,
		timeFn:              timeFn,
	}
}

func (a *agent) GetStat() *agentPb.SendStat {
	response := &agentPb.SendStat{
		CpuInfo:    &agentPb.CpuInfo{},
		MemoryInfo: &agentPb.MemoryInfo{},
		DiskInfo:   &agentPb.DiskInfo{},
		HostInfo: &agentPb.HostInfo{
			PlatformInfo: &agentPb.PlatformInfo{},
		},
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		cpuStat, err := a.cpuStatFn(cpuInterval, true)

		if err != nil || cpuStat == nil {
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
		swapMemoryStat, err := a.swapMemoryStatFn()

		if err != nil || swapMemoryStat == nil {
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
		memoryStat, err := a.virtualMemoryStatFn()

		if err != nil || memoryStat == nil {
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

		disks, err := a.diskStatFn(false)

		if err != nil || disks == nil {
			return
		}
		diskStat := make(map[string]*agentPb.DiskInfo_Disk)
		for _, d := range disks {
			diskInfo, err := a.diskUsageFn(d.Mountpoint)
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

		netStat, err := a.netStatFn(false)
		if err != nil || netStat == nil || len(netStat) == 0 {
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

	wg.Add(1)

	go func() {
		defer wg.Done()
		hostStat, err := a.hostStatFn()
		if err != nil || hostStat == nil {
			return
		}

		response.HostInfo = &agentPb.HostInfo{
			HostName: hostStat.Hostname,
			Os:       hostStat.OS,
			PlatformInfo: &agentPb.PlatformInfo{
				Name:    hostStat.Platform,
				Family:  hostStat.PlatformFamily,
				Version: hostStat.PlatformVersion,
			},
		}
	}()

	wg.Wait()

	response.Time = a.timeFn()

	return response
}
