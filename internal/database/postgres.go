package database

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"time"
)

type postgres struct {
	db *gorm.DB
}

type Snapshot struct {
	gorm.Model
	SchedulerID string    `gorm:"column:schedulerId"`
	Code        string    `gorm:"column:code"`
	Type        string    `gorm:"column:type"`
	Error       string    `gorm:"column:error"`
	Meta        *MetaData `gorm:"column:meta"`
}

type MetaData struct {
	gorm.Model
	SnapshotID uint      `gorm:"column:snapshotId"`
	StartTime  time.Time `gorm:"column:startTime"`
	EndTime    time.Time `gorm:"column:endTime"`
	Value      []byte    `gorm:"column:value"`
}

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

const (
	dmMetaDataCollection    = "meta_data"
	dbSnapshotCollection    = "snapshots"     //TODO: check
	dbStatRequestCollection = "stat_requests" //TODO: check
)

var (
	errorDataBase = errors.New("ERROR_DATABASE_OPERATION")
)

func (p *postgres) Migrate() error {
	models := []interface{}{
		&Snapshot{},
		&MetaData{},
		&StatRequest{},
		&CPUInfo{},
		&MemoryInfo{},
		&MemoryMem{},
		&MemorySwap{},
		&DiskInfo{},
		&NetInfo{},
	}

	var err error
	for _, model := range models {
		err = p.db.AutoMigrate(model).Error // migrate models one-by-one
	}

	return err
}

func (p *postgres) InsertSnapshot(data *apiPb.SchedulerResponse) error {
	snapshot, err := ConvertToPostgresSnapshot(data)
	if err != nil {
		return err
	}
	if err := p.db.Table(dbSnapshotCollection).Create(snapshot).Error; err != nil {
		return errorDataBase
	}
	return nil
}

func (p *postgres) GetSnapshots(request *apiPb.GetSchedulerInformationRequest) ([]*apiPb.SchedulerSnapshot, int32, error) {
	timeFrom, timeTo, err := getTime(request.GetTimeRange())
	if err != nil {
		return nil, -1, err
	}

	var count int64
	if request.GetStatus() == apiPb.SchedulerCode_SCHEDULER_CODE_UNSPECIFIED {
		err = p.db.Table(dbSnapshotCollection).
			Joins(fmt.Sprintf(`JOIN "%s" ON "%s.snapshotId" = "%s"."ID"`, dmMetaDataCollection, dmMetaDataCollection, dbSnapshotCollection)).
			Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), request.GetSchedulerId()).
			Where(fmt.Sprintf(`"%s"."startTime" BETWEEN ? and ?`, dmMetaDataCollection), timeFrom, timeTo).
			Count(&count).Error
	} else {
		err = p.db.Table(dbSnapshotCollection).
			Joins(fmt.Sprintf(`JOIN "%s" ON "%s.snapshotId" = "%s"."ID"`, dmMetaDataCollection, dmMetaDataCollection, dbSnapshotCollection)).
			Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), request.GetSchedulerId()).
			Where(fmt.Sprintf(`"%s"."startTime" BETWEEN ? and ?`, dmMetaDataCollection), timeFrom, timeTo).
			Where(fmt.Sprintf(`"%s"."code" = ?`, dbSnapshotCollection), request.GetStatus().String()).
			Count(&count).Error
	}

	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, request.GetPagination())

	var dbSnapshots []*Snapshot
	if request.GetStatus() == apiPb.SchedulerCode_SCHEDULER_CODE_UNSPECIFIED {
		err = p.db.
			Table(dbSnapshotCollection).
			Set("gorm:auto_preload", true).
			Joins(fmt.Sprintf(`JOIN "%s" ON "%s.snapshotId" = "%s"."ID"`, dmMetaDataCollection, dmMetaDataCollection, dbSnapshotCollection)).
			Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), request.GetSchedulerId()).
			Where(fmt.Sprintf(`"%s"."startTime" BETWEEN ? and ?`, dmMetaDataCollection), timeFrom, timeTo).
			Order(getOrder(request.GetSort()) + " " + getDirection(request.GetSort())).
			Offset(offset).
			Limit(limit).
			Find(&dbSnapshots).Error
	} else {
		err = p.db.
			Table(dbSnapshotCollection).
			Set("gorm:auto_preload", true).
			Joins(fmt.Sprintf(`JOIN "%s" ON "%s.snapshotId" = "%s"."ID"`, dmMetaDataCollection, dmMetaDataCollection, dbSnapshotCollection)).
			Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), request.GetSchedulerId()).
			Where(fmt.Sprintf(`"%s"."startTime" BETWEEN ? and ?`, dmMetaDataCollection), timeFrom, timeTo).
			Where(fmt.Sprintf(`"%s"."code" = ?`, dbSnapshotCollection), request.GetStatus().String()).
			Order(getOrder(request.GetSort()) + " " + getDirection(request.GetSort())).
			Offset(offset).
			Limit(limit).
			Find(&dbSnapshots).Error

	}

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgresSnapshots(dbSnapshots), int32(count), nil
}

func (p *postgres) GetSnapshotsUptime(request *apiPb.GetSchedulerUptimeRequest) (float64, float64, error) {
	timeFrom, timeTo, err := getTime(request.GetTimeRange())
	if err != nil {
		return -1, -1, err
	}

	var countOk int64
	err = p.db.Table(dbSnapshotCollection).
		Joins(fmt.Sprintf(`JOIN "%s" ON "%s.snapshotId" = "%s"."ID"`, dmMetaDataCollection, dmMetaDataCollection, dbSnapshotCollection)).
		Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), request.GetSchedulerId()).
		Where(fmt.Sprintf(`"%s"."startTime" BETWEEN ? and ?`, dmMetaDataCollection), timeFrom, timeTo).
		Where(fmt.Sprintf(`"%s"."code" = ?`, dbSnapshotCollection), apiPb.SchedulerCode_OK.String()).
		Count(&countOk).Error
	if err != nil {
		return -1, -1, err
	}

	var countAll int
	err = p.db.Table(dbSnapshotCollection).
		Joins(fmt.Sprintf(`JOIN "%s" ON "%s.snapshotId" = "%s"."ID"`, dmMetaDataCollection, dmMetaDataCollection, dbSnapshotCollection)).
		Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), request.GetSchedulerId()).
		Where(fmt.Sprintf(`"%s"."startTime" BETWEEN ? and ?`, dmMetaDataCollection), timeFrom, timeTo).
		Count(&countOk).Error
	if err != nil {
		return -1, -1, err
	}

	var dbSnapshots []*Snapshot
	err = p.db.
		Table(dbSnapshotCollection).
		Set("gorm:auto_preload", true).
		Joins(fmt.Sprintf(`JOIN "%s" ON "%s.snapshotId" = "%s"."ID"`, dmMetaDataCollection, dmMetaDataCollection, dbSnapshotCollection)).
		Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), request.GetSchedulerId()).
		Where(fmt.Sprintf(`"%s"."startTime" BETWEEN ? and ?`, dmMetaDataCollection), timeFrom, timeTo).
		Where(fmt.Sprintf(`"%s"."code" = ?`, dbSnapshotCollection), apiPb.SchedulerCode_OK.String()).
		Find(&dbSnapshots).Error
	if err != nil {
		return -1, -1, errorDataBase
	}

	latency := int64(0)
	for _, snapshot := range dbSnapshots {
		//Recieveing time difference in mellisecinds
		if snapshot.Meta != nil {
			latency += (snapshot.Meta.EndTime.UnixNano() - snapshot.Meta.StartTime.UnixNano()) / int64(time.Millisecond)
		}
	}

	return float64(countOk) / float64(countAll), float64(latency) / float64(countOk), nil
}

func getOrder(request *apiPb.SortingSchedulerList) string {
	if request == nil {
		return fmt.Sprintf(`"%s"."startTime"`, dmMetaDataCollection)
	}
	orderMap := map[apiPb.SortSchedulerList]string{
		apiPb.SortSchedulerList_SORT_SCHEDULER_LIST_UNSPECIFIED: fmt.Sprintf(`"%s"."startTime"`, dmMetaDataCollection),
		apiPb.SortSchedulerList_BY_START_TIME:                   fmt.Sprintf(`"%s"."startTime"`, dmMetaDataCollection),
		apiPb.SortSchedulerList_BY_END_TIME:                     fmt.Sprintf(`"%s"."endTime"`, dmMetaDataCollection),
		apiPb.SortSchedulerList_BY_LATENCY:                      fmt.Sprintf(`"%s"."endTime" - "%s"."startTime"`, dmMetaDataCollection, dmMetaDataCollection),
	}
	if res, ok := orderMap[request.GetSortBy()]; ok {
		return res
	}
	return fmt.Sprintf(`"%s"."startTime"`, dmMetaDataCollection)
}

func getDirection(request *apiPb.SortingSchedulerList) string {
	if request == nil {
		return ``
	}
	directionMap := map[apiPb.SortDirection]string{
		apiPb.SortDirection_SORT_DIRECTION_UNSPECIFIED: ``,
		apiPb.SortDirection_ASC:                        `asc`,
		apiPb.SortDirection_DESC:                       `desc`,
	}
	if res, ok := directionMap[request.GetDirection()]; ok {
		return res
	}
	return ``
}

func (p *postgres) InsertStatRequest(data *apiPb.Metric) error {
	pgData, err := ConvertToPostgressStatRequest(data)
	if err != nil {
		return err
	}
	if err := p.db.Table(dbStatRequestCollection).Create(pgData).Error; err != nil {
		//TODO: log?
		return errorDataBase
	}
	return nil
}

func (p *postgres) GetStatRequest(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int64
	err = p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	//TODO: test if it works
	var statRequests []*StatRequest
	err = p.db.
		Set("gorm:auto_preload", true).
		//Preload("disk_infos").
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Order("time").
		Offset(offset).
		Limit(limit).
		Find(&statRequests).Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}

func (p *postgres) GetCPUInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, cpuInfoKey)
}

func (p *postgres) GetMemoryInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int64
	err = p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	//TODO: test if it works
	var statRequests []*StatRequest
	err = p.db.
		Preload("MemoryInfo").
		Preload("MemoryInfo.Mem").
		Preload("MemoryInfo.Swap").
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Order("time").
		Offset(offset).
		Limit(limit).
		Find(&statRequests).
		Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}

func (p *postgres) GetDiskInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, diskInfoKey)
}

func (p *postgres) GetNetInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, netInfoKey)
}

func (p *postgres) getSpecialRecords(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter, key string) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int64
	err = p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	//TODO: test if it works
	var statRequests []*StatRequest
	err = p.db.
		Preload(key).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Order("time").
		Offset(offset).
		Limit(limit).
		Find(&statRequests).
		Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}

func getTime(filter *apiPb.TimeFilter) (time.Time, time.Time, error) {
	timeFrom := time.Unix(0, 0)
	timeTo := time.Now()
	var err error
	if filter != nil {
		if filter.GetFrom() != nil {
			timeFrom, err = ptypes.Timestamp(filter.From)
		}
		if filter.GetTo() != nil {
			timeTo, err = ptypes.Timestamp(filter.To)
		}
	}
	return timeFrom, timeTo, err
}

//Return offset and limit
func getOffsetAndLimit(count int64, pagination *apiPb.Pagination) (int, int) {
	if pagination != nil {
		if pagination.Page == -1 {
			return int(count) - int(pagination.Limit), int(pagination.Limit)
		}
		return int(pagination.GetLimit() * pagination.GetPage()), int(pagination.GetLimit())
	}
	return 0, int(count)
}
