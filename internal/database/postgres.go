package database

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"time"
)

type postgres struct {
	db *gorm.DB
}

type Model struct {
	CreatedAt time.Time  `json:"createdAt" gorm:"index"`
	UpdatedAt time.Time  `json:"createdAt,omitEmpty" gorm:"index"`
	DeletedAt *time.Time `json:"createdAt" gorm:"index"`
}

type Snapshot struct {
	gorm.Model
	SchedulerId string    `gorm:"schedulerId"`
	Code        string    `gorm:"code"`
	Type        string    `gorm:"column:type"`
	Error       string    `gorm:"error"`
	Meta        *MetaData `gorm:"meta"`
}

type MetaData struct {
	gorm.Model
	SnapshotId uint           `gorm:"snapshotId"`
	StartTime  time.Time      `gorm:"startTime"`
	EndTime    time.Time      `gorm:"endTime"`
	Value      *_struct.Value `gorm:"value"` //TODO: google
}

//Agent gorm description
type StatRequest struct {
	gorm.Model
	CpuInfo    []*CpuInfo  `gorm:"cpuInfo"`
	MemoryInfo *MemoryInfo `gorm:"memoryInfo"`
	DiskInfo   []*DiskInfo `gorm:"diskInfo"`
	NetInfo    []*NetInfo  `gorm:"netInfo"`
	Time       time.Time   `gorm:"time"`
}


const (
	cpuInfoKey = "cpuInfo"
	memoryInfoKey = "memoryInfo"
	diskInfoKey = "diskInfo"
	netInfoKey = "netInfo"
)


type CpuInfo struct {
	gorm.Model
	StatRequestID uint    `gorm:"statRequestId"`
	Load          float64 `gorm:"load"`
}

type MemoryInfo struct {
	gorm.Model
	StatRequestID uint    `gorm:"statRequestId"`
	Mem           *Memory `gorm:"mem"`
	Swap          *Memory `gorm:"swap"`
}

type Memory struct { ///Will we need check for Total = used + free + shared?
	gorm.Model
	MemoryInfoID uint    `gorm:"memoryInfoId"`
	Total        uint64  `gorm:"total"`
	Used         uint64  `gorm:"used"`
	Free         uint64  `gorm:"free"`
	Shared       uint64  `gorm:"shared"`
	UsedPercent  float64 `gorm:"usedPercent"`
}

type DiskInfo struct { ///Will we need check for Total = used + free?
	gorm.Model
	StatRequestID uint    `gorm:"statRequestId"`
	Name          string  `gorm:"name"`
	Total         uint64  `gorm:"total"`
	Free          uint64  `gorm:"free"`
	Used          uint64  `gorm:"used"`
	UsedPercent   float64 `gorm:"usedPercent"`
}

type NetInfo struct {
	gorm.Model
	StatRequestID uint   `gorm:"statRequestId"`
	Name          string `gorm:"name"`
	BytesSent     uint64 `gorm:"bytesSent"`
	BytesRecv     uint64 `gorm:"bytesRecv"`
	PacketsSent   uint64 `gorm:"packetsSent"`
	PacketsRecv   uint64 `gorm:"packetsRecv"`
	ErrIn         uint64 `gorm:"errIn"`
	ErrOut        uint64 `gorm:"errOut"`
	DropIn        uint64 `gorm:"dropIn"`
	DropOut       uint64 `gorm:"dropOut"`
}

const (
	dbSnapshotCollection    = "snapshots"     //TODO: check
	dbStatRequestCollection = "stat_requests" //TODO: check
)

var (
	errorConnection = errors.New("ERROR_CONNECTING_TO_POSTGRES")
	errorDataBase   = errors.New("ERROR_DATABASE_OPERATION")
)

func (p *postgres) newClient(getDB func() (*gorm.DB, error)) error {
	var err error
	p.db, err = getDB()
	if err != nil {
		fmt.Println(err.Error()) //TODO: log?
		return errorConnection
	}
	p.db.LogMode(true)
	return p.Migrate()
}

func (p *postgres) Migrate() (resErr error) {
	resErr = nil
	models := []interface{}{
		&Snapshot{},
		&MetaData{},
		&StatRequest{},
		&CpuInfo{},
		&MemoryInfo{},
		&Memory{},
		&DiskInfo{},
		&NetInfo{},
	}

	for _, model := range models {
		err := p.db.AutoMigrate(model).Error // migrate models one-by-one
		if err != nil {
			fmt.Println(err.Error()) //TODO: log?
			resErr = errorDataBase
		}
	}
	return resErr
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

func (p *postgres) GetSnapshots(id string) ([]*apiPb.SchedulerSnapshot, error) {
	var snpashots []*Snapshot
	if err := p.db.Table(dbSnapshotCollection).Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), id).Find(&snpashots).Error; err != nil {
		fmt.Println(err.Error()) //TODO: log?
		return nil, errorDataBase
	}
	snapshots, errs := ConvertFromPostgresSnapshots(snpashots)
	if len(errs) != 0 {
		//TODO: log
	}
	return snapshots, nil
}

func (p *postgres) InsertStatRequest(data *apiPb.Metric) error {
	pgData, err := ConvertToPostgressStatRequest(data)
	if err != nil {
		return err
	}
	if err := p.db.Table(dbStatRequestCollection).Create(pgData).Error; err != nil {
		fmt.Println(err.Error()) //TODO: log?
		return errorDataBase
	}
	return nil
}

func (p *postgres) GetStatRequest(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	var statRequests []*StatRequest
	if err := p.db.Table(dbStatRequestCollection).Where(fmt.Sprintf(`"%s"."id" = ?`, dbStatRequestCollection), id).Find(&statRequests).Error; err != nil {
		fmt.Println(err.Error()) //TODO: log?
		return nil, -1, errorDataBase
	}
	count := int32(-1)
	p.db.Table(dbStatRequestCollection).Where(fmt.Sprintf(`"%s"."id" = ?`, dbStatRequestCollection), id).Count(count)
	res, err := ConvertFromPostgressStatRequests(statRequests)
	return res, count, err
}

func (p *postgres) GetCpuInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(id, pagination, filter, cpuInfoKey)
}

func (p *postgres) GetMemoryInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(id, pagination, filter, memoryInfoKey)
}

func (p *postgres) GetDiskInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(id, pagination, filter, diskInfoKey)
}

func (p *postgres) GetNetInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(id, pagination, filter, netInfoKey)
}

func (p *postgres) getSpecialRecords(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter, key string) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	count := int32(-1)
	if err := p.db.Table(dbStatRequestCollection).Where(fmt.Sprintf(`"%s"."id" = ?`, dbStatRequestCollection), id).Select(key).Count(count).Error; err != nil {
		return nil, -1, err
	}

	offset := int32(0)
	limit := count
	if pagination != nil {
		offset = pagination.GetLimit() * pagination.GetPage()
		limit = pagination.GetLimit()
	}

	//TODO: test if it works
	var statRequests []*StatRequest
	if err := p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."id" = ?`, dbStatRequestCollection), id).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Select("%s, time", key).
		Order("time").
		Offset(offset).
		Limit(limit).
		Find(&statRequests).Error;
		err != nil {

		fmt.Println(err.Error()) //TODO: log?
		return nil, -1, errorDataBase
	}
	res, err := ConvertFromPostgressStatRequests(statRequests)
	return res, count, err
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
			timeFrom, err = ptypes.Timestamp(filter.From)
		}
	}
	return timeFrom, timeTo, err
}
