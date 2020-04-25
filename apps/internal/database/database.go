package database

import "github.com/jinzhu/gorm"

type Database interface {
	InsertMetaData(data *MetaData) error
	GetMetaData(id string) (*MetaData, error)
	InsertStatRequest(data *StatRequest) error
	GetStatRequest(id string) (*StatRequest, error)
}

func New(getDB func() (*gorm.DB, error)) (Database, error) {
	db := &postgres{
	}
	err := db.newClient(getDB)
	return db, err
}