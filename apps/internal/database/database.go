package database

type Database interface {
	InsertMetaData(data *MetaData) error
	GetMetaData(id string) (*MetaData, error)
	InsertStatRequest(data *StatRequest) error
	GetStatRequest(id string) (*StatRequest, error)
}

func New() Database {
	return &postgres{}
}