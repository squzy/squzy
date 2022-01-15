package agent

import (
	"errors"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("Should: create new agent", func(t *testing.T) {
		a := New(func(duration time.Duration, b bool) (float64s []float64, err error) {
			return nil, nil
		}, func() (stat *mem.SwapMemoryStat, err error) {
			return nil, nil
		}, func() (stat *mem.VirtualMemoryStat, err error) {
			return nil, nil
		}, func(b bool) (stats []disk.PartitionStat, err error) {
			return nil, nil
		}, func(s string) (stat *disk.UsageStat, err error) {
			return nil, nil
		}, func(b bool) (stat []net.IOCountersStat, err error) {
			return nil, nil
		}, func() *timestamp.Timestamp {
			return &timestamp.Timestamp{}
		})
		assert.IsType(t, &agent{}, a)
		assert.Implements(t, (*Agent)(nil), a)
	})
}

func TestAgent_GetStat(t *testing.T) {
	t.Run("Should: return stat about computer", func(t *testing.T) {
		a := New(func(duration time.Duration, b bool) (float64s []float64, err error) {
			return nil, nil
		}, func() (stat *mem.SwapMemoryStat, err error) {
			return nil, nil
		}, func() (stat *mem.VirtualMemoryStat, err error) {
			return nil, nil
		}, func(b bool) (stats []disk.PartitionStat, err error) {
			return nil, nil
		}, func(s string) (stat *disk.UsageStat, err error) {
			return nil, nil
		}, func(b bool) (stat []net.IOCountersStat, err error) {
			return nil, nil
		}, func() *timestamp.Timestamp {
			return &timestamp.Timestamp{}
		})
		assert.IsType(t, &apiPb.Metric{}, a.GetStat())
	})
	t.Run("Should: return cpu info", func(t *testing.T) {
		a := New(func(duration time.Duration, b bool) (float64s []float64, err error) {
			return []float64{6}, nil
		}, func() (stat *mem.SwapMemoryStat, err error) {
			return nil, nil
		}, func() (stat *mem.VirtualMemoryStat, err error) {
			return nil, nil
		}, func(b bool) (stats []disk.PartitionStat, err error) {
			return nil, nil
		}, func(s string) (stat *disk.UsageStat, err error) {
			return nil, nil
		}, func(b bool) (stat []net.IOCountersStat, err error) {
			return nil, nil
		}, func() *timestamp.Timestamp {
			return &timestamp.Timestamp{}
		})
		assert.EqualValues(t, []*apiPb.CpuInfo_CPU{
			{
				Load: 6,
			},
		}, a.GetStat().CpuInfo.Cpus)
	})
	t.Run("Should: return memory info", func(t *testing.T) {
		a := New(func(duration time.Duration, b bool) (float64s []float64, err error) {
			return nil, nil
		}, func() (stat *mem.SwapMemoryStat, err error) {
			return &mem.SwapMemoryStat{Used: 6}, nil
		}, func() (stat *mem.VirtualMemoryStat, err error) {
			return nil, nil
		}, func(b bool) (stats []disk.PartitionStat, err error) {
			return nil, nil
		}, func(s string) (stat *disk.UsageStat, err error) {
			return nil, nil
		}, func(b bool) (stat []net.IOCountersStat, err error) {
			return nil, nil
		}, func() *timestamp.Timestamp {
			return &timestamp.Timestamp{}
		})
		assert.EqualValues(t, &apiPb.MemoryInfo_Memory{
			Used: 6,
		}, a.GetStat().MemoryInfo.Swap)
	})
	t.Run("Should: return virtual memory info", func(t *testing.T) {
		a := New(func(duration time.Duration, b bool) (float64s []float64, err error) {
			return nil, nil
		}, func() (stat *mem.SwapMemoryStat, err error) {
			return nil, nil
		}, func() (stat *mem.VirtualMemoryStat, err error) {
			return &mem.VirtualMemoryStat{Used: 6}, nil
		}, func(b bool) (stats []disk.PartitionStat, err error) {
			return nil, nil
		}, func(s string) (stat *disk.UsageStat, err error) {
			return nil, nil
		}, func(b bool) (stat []net.IOCountersStat, err error) {
			return nil, nil
		}, func() *timestamp.Timestamp {
			return &timestamp.Timestamp{}
		})
		assert.EqualValues(t, &apiPb.MemoryInfo_Memory{
			Used: 6,
		}, a.GetStat().MemoryInfo.Mem)
	})
	t.Run("Should: return disk stat", func(t *testing.T) {
		a := New(func(duration time.Duration, b bool) (float64s []float64, err error) {
			return nil, nil
		}, func() (stat *mem.SwapMemoryStat, err error) {
			return nil, nil
		}, func() (stat *mem.VirtualMemoryStat, err error) {
			return nil, nil
		}, func(b bool) (stats []disk.PartitionStat, err error) {
			return []disk.PartitionStat{{Mountpoint: "/"}}, nil
		}, func(s string) (stat *disk.UsageStat, err error) {
			return &disk.UsageStat{
				Used: 6,
			}, nil
		}, func(b bool) (stat []net.IOCountersStat, err error) {
			return nil, nil
		}, func() *timestamp.Timestamp {
			return &timestamp.Timestamp{}
		})
		assert.EqualValues(t, &apiPb.DiskInfo{Disks: map[string]*apiPb.DiskInfo_Disk{
			"/": {Used: 6},
		}}, a.GetStat().DiskInfo)
	})
	t.Run("Should: return net stat", func(t *testing.T) {
		a := New(func(duration time.Duration, b bool) (float64s []float64, err error) {
			return nil, nil
		}, func() (stat *mem.SwapMemoryStat, err error) {
			return nil, nil
		}, func() (stat *mem.VirtualMemoryStat, err error) {
			return nil, nil
		}, func(b bool) (stats []disk.PartitionStat, err error) {
			return nil, nil
		}, func(s string) (stat *disk.UsageStat, err error) {
			return nil, nil
		}, func(b bool) (stat []net.IOCountersStat, err error) {
			return []net.IOCountersStat{
				{
					Name:      "test",
					BytesRecv: 5,
				},
			}, nil
		}, func() *timestamp.Timestamp {
			return &timestamp.Timestamp{}
		})
		assert.EqualValues(t, &apiPb.NetInfo{Interfaces: map[string]*apiPb.NetInfo_Interface{
			"test": {
				BytesRecv: 5,
			},
		}}, a.GetStat().NetInfo)
	})
	t.Run("Should: fill default value if throw error", func(t *testing.T) {
		errValue := errors.New("test")
		a := New(func(duration time.Duration, b bool) (float64s []float64, err error) {
			return nil, errValue
		}, func() (stat *mem.SwapMemoryStat, err error) {
			return nil, errValue
		}, func() (stat *mem.VirtualMemoryStat, err error) {
			return nil, errValue
		}, func(b bool) (stats []disk.PartitionStat, err error) {
			return nil, errValue
		}, func(s string) (stat *disk.UsageStat, err error) {
			return nil, errValue
		}, func(b bool) (stat []net.IOCountersStat, err error) {
			return nil, errValue
		}, func() *timestamp.Timestamp {
			return &timestamp.Timestamp{}
		})
		assert.EqualValues(t, &apiPb.Metric{
			CpuInfo:    &apiPb.CpuInfo{},
			MemoryInfo: &apiPb.MemoryInfo{},
			DiskInfo:   &apiPb.DiskInfo{},
			Time:       &timestamp.Timestamp{},
			NetInfo:    &apiPb.NetInfo{},
		}, a.GetStat())
	})
	t.Run("Should: fill default value if throw error, if disk usage not throw error", func(t *testing.T) {
		errValue := errors.New("test")
		a := New(func(duration time.Duration, b bool) (float64s []float64, err error) {
			return nil, errValue
		}, func() (stat *mem.SwapMemoryStat, err error) {
			return nil, errValue
		}, func() (stat *mem.VirtualMemoryStat, err error) {
			return nil, errValue
		}, func(b bool) (stats []disk.PartitionStat, err error) {
			return []disk.PartitionStat{{Mountpoint: "/"}}, nil
		}, func(s string) (stat *disk.UsageStat, err error) {
			return nil, errValue
		}, func(b bool) (stat []net.IOCountersStat, err error) {
			return nil, errValue
		}, func() *timestamp.Timestamp {
			return &timestamp.Timestamp{}
		})
		assert.Equal(t, &apiPb.Metric{
			CpuInfo:    &apiPb.CpuInfo{},
			MemoryInfo: &apiPb.MemoryInfo{},
			NetInfo:    &apiPb.NetInfo{},
			DiskInfo: &apiPb.DiskInfo{
				Disks: make(map[string]*apiPb.DiskInfo_Disk),
			},
			Time: &timestamp.Timestamp{},
		}, a.GetStat())
	})
}
