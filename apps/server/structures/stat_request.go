package structures

import "time"

type StatRequest struct {
	Id         string     `gorm:"index:id"`
	CpuInfo    float64    `gorm:"cpuInfo"`
	MemoryInfo int64      `gorm:"memoryInfo"`
	DiskInfo   int64      `gorm:"diskInfo"`
	NetInfo    int64      `gorm:"netInfo"`
	Time       *time.Time `gorm:"time"`
}
