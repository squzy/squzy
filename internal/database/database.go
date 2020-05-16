package database

import (
	"github.com/jinzhu/gorm"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

type Database interface {
	InsertSnapshot(data *apiPb.SchedulerResponse) error //TODO: fix
	GetSnapshots(id string) ([]*apiPb.SchedulerSnapshot, error) //TODO: fix
	InsertStatRequest(data *apiPb.Metric) error
	GetStatRequest(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error)
	GetCpuInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error)
	GetMemoryInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error)
	GetDiskInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error)
	GetNetInfo(id string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error)
}

func New(getDB func() (*gorm.DB, error)) (Database, error) {
	db := &postgres{
	}
	err := db.newClient(getDB)
	return db, err
}