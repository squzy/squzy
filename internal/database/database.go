package database

import (
	"github.com/jinzhu/gorm"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

type Database interface {
	InsertSnapshot(data *apiPb.SchedulerResponse) error //TODO: fix
	GetSnapshots(id string) ([]*apiPb.Snapshot, error) //TODO: fix
	InsertStatRequest(data *apiPb.SendMetricsRequest) error
	GetStatRequest(id string) (*apiPb.SendMetricsRequest, error)
}

func New(getDB func() (*gorm.DB, error)) (Database, error) {
	db := &postgres{
	}
	err := db.newClient(getDB)
	return db, err
}