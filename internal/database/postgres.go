package database

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"squzy/internal/database/convertion"
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

type MetaData struct {
	gorm.Model
	Location  string     `gorm:"location"`
	Port      int32      `gorm:"port"`
	StartTime *time.Time `gorm:"startTime"`
	EndTime   *time.Time `gorm:"endTime"`
	Type      int32      `gorm:"column:type"`
}

type StatRequest struct {
	gorm.Model
	CpuInfo    []*CpuInfo  `gorm:"cpuInfo"`
	MemoryInfo *MemoryInfo `gorm:"memoryInfo"`
	DiskInfo   []*DiskInfo `gorm:"diskInfo"`
	NetInfo    []*NetInfo  `gorm:"netInfo"`
	Time       *time.Time  `gorm:"time"`
}

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
	dbMetaDataCollection    = "meta_data"    //TODO: check
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

func (p *postgres) InsertMetaData(data *MetaData) error {
	if err := p.db.Table(dbMetaDataCollection).Create(data).Error; err != nil {
		return errorDataBase
	}
	return nil
}

func (p *postgres) GetMetaData(id string) (*MetaData, error) {
	metaData := &MetaData{}
	if err := p.db.Table(dbMetaDataCollection).Where(fmt.Sprintf(`"%s"."id" = ?`, dbMetaDataCollection), id).First(metaData).Error; err != nil {
		fmt.Println(err.Error()) //TODO: log?
		return nil, errorDataBase
	}
	return metaData, nil
}

func (p *postgres) InsertStatRequest(data *apiPb.SendMetricsRequest) error {
	pgData := convertion.ConvertToPostgressStatRequest(data)
	if err := p.db.Table(dbStatRequestCollection).Create(pgData).Error; err != nil {
		fmt.Println(err.Error()) //TODO: log?
		return errorDataBase
	}
	return nil
}

func (p *postgres) GetStatRequest(id string) (*apiPb.SendMetricsRequest, error) {
	statRequest := &StatRequest{}
	if err := p.db.Table(dbStatRequestCollection).Where(fmt.Sprintf(`"%s"."id" = ?`, dbStatRequestCollection), id).First(statRequest).Error; err != nil {
		fmt.Println(err.Error()) //TODO: log?
		return nil, errorDataBase
	}
	return convertion.ConvertFromPostgressStatRequest(statRequest), nil
}
