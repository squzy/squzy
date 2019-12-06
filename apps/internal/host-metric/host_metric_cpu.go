package host_metric

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

type cpu_metric struct {
	interval time.Duration
}

func (c *cpu_metric) GetStat() interface{} {
	info, err := cpu.Percent(c.interval, true)
	if err !=nil {
		return nil
	}
	fmt.Println(info)
	return nil
}

func NewCpuMetric(interval time.Duration) Metric {
	return &cpu_metric{
		interval:interval,
	}
}