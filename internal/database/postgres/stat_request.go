package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"time"
)

//Agent gorm description
type StatRequest struct {
	gorm.Model
	AgentID    string `gorm:"column:agentID"`
	AgentName  string `gorm:"column:agentName"`
	CPUInfo    []*CPUInfo
	MemoryInfo *MemoryInfo `gorm:"column:memoryInfo"`
	DiskInfo   []*DiskInfo `gorm:"column:diskInfo"`
	NetInfo    []*NetInfo  `gorm:"column:netInfo"`
	Time       time.Time   `gorm:"column:time"`
}

const (
	cpuInfoKey  = "CPUInfo"
	diskInfoKey = "DiskInfo"
	netInfoKey  = "NetInfo"
)

type CPUInfo struct {
	gorm.Model
	StatRequestID uint    `gorm:"column:statRequestId"`
	Load          float64 `gorm:"column:load"`
}

type MemoryInfo struct {
	gorm.Model
	StatRequestID uint        `gorm:"column:statRequestId"`
	Mem           *MemoryMem  `gorm:"column:mem"`
	Swap          *MemorySwap `gorm:"column:swap"`
}

type MemoryMem struct {
	gorm.Model
	MemoryInfoID uint    `gorm:"column:memoryInfoId"`
	Total        uint64  `gorm:"column:total"`
	Used         uint64  `gorm:"column:used"`
	Free         uint64  `gorm:"column:free"`
	Shared       uint64  `gorm:"column:shared"`
	UsedPercent  float64 `gorm:"column:usedPercent"`
}

type MemorySwap struct {
	gorm.Model
	MemoryInfoID uint    `gorm:"column:memoryInfoId"`
	Total        uint64  `gorm:"column:total"`
	Used         uint64  `gorm:"column:used"`
	Free         uint64  `gorm:"column:free"`
	Shared       uint64  `gorm:"column:shared"`
	UsedPercent  float64 `gorm:"column:usedPercent"`
}

type DiskInfo struct {
	gorm.Model
	StatRequestID uint    `gorm:"column:statRequestId"`
	Name          string  `gorm:"column:name"`
	Total         uint64  `gorm:"column:total"`
	Free          uint64  `gorm:"column:free"`
	Used          uint64  `gorm:"column:used"`
	UsedPercent   float64 `gorm:"column:usedPercent"`
}

type NetInfo struct {
	gorm.Model
	StatRequestID uint   `gorm:"column:statRequestId"`
	Name          string `gorm:"column:name"`
	BytesSent     uint64 `gorm:"column:bytesSent"`
	BytesRecv     uint64 `gorm:"column:bytesRecv"`
	PacketsSent   uint64 `gorm:"column:packetsSent"`
	PacketsRecv   uint64 `gorm:"column:packetsRecv"`
	ErrIn         uint64 `gorm:"column:errIn"`
	ErrOut        uint64 `gorm:"column:errOut"`
	DropIn        uint64 `gorm:"column:dropIn"`
	DropOut       uint64 `gorm:"column:dropOut"`
}

var (
	agentIdFilterString         = fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection)
	statRequestTimeFilterString = fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection)
	statRequestTimeString       = "time"
)

func (p *Postgres) InsertStatRequest(data *apiPb.Metric) error {
	pgData, err := ConvertToPostgressStatRequest(data)
	if err != nil {
		return err
	}
	if err := p.Db.Table(dbStatRequestCollection).Create(pgData).Error; err != nil {
		//TODO: log?
		return errorDataBase
	}
	return nil
}

func (p *Postgres) GetStatRequest(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	fmt.Println(agentID)

	var count int64
	err = p.Db.Table(dbStatRequestCollection).
		Where(agentIdFilterString, agentID).
		Where(statRequestTimeFilterString, timeFrom, timeTo).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	var statRequests []*StatRequest
	err = p.Db.
		Set("gorm:auto_preload", true).
		Where(agentIdFilterString, agentID).
		Where(statRequestTimeFilterString, timeFrom, timeTo).
		Order(statRequestTimeString).
		Offset(offset).
		Limit(limit).
		Find(&statRequests).Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}

func (p *Postgres) GetCPUInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, cpuInfoKey)
}

func (p *Postgres) GetMemoryInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int64
	err = p.Db.Table(dbStatRequestCollection).
		Where(agentIdFilterString, agentID).
		Where(statRequestTimeFilterString, timeFrom, timeTo).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	var statRequests []*StatRequest
	err = p.Db.
		Preload("MemoryInfo").
		Preload("MemoryInfo.Mem").
		Preload("MemoryInfo.Swap").
		Where(agentIdFilterString, agentID).
		Where(statRequestTimeFilterString, timeFrom, timeTo).
		Order(statRequestTimeString).
		Offset(offset).
		Limit(limit).
		Find(&statRequests).
		Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}

func (p *Postgres) GetDiskInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, diskInfoKey)
}

func (p *Postgres) GetNetInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, netInfoKey)
}

func (p *Postgres) getSpecialRecords(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter, key string) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int64
	err = p.Db.Table(dbStatRequestCollection).
		Where(agentIdFilterString, agentID).
		Where(statRequestTimeFilterString, timeFrom, timeTo).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	var statRequests []*StatRequest
	err = p.Db.
		Preload(key).
		Where(agentIdFilterString, agentID).
		Where(statRequestTimeFilterString, timeFrom, timeTo).
		Order(statRequestTimeString).
		Offset(offset).
		Limit(limit).
		Find(&statRequests).
		Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}
