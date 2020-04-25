package database

type Database interface {
	InsertMetaData(data *MetaData) error
	GetMetaData(id string) (*MetaData, error)
	InsertStatRequest(data *StatRequest) error
	GetStatRequest(id string) (*StatRequest, error)
}

func New() (Database, error) {
	db := &postgres{
		host:     "",
		port:     "",
		user:     "",
		password: "",
		dbname:   "",
	}
	err := db.newClient()
	return db, err
}