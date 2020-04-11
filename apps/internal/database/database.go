package database

type Database interface {
	Save() error
}

type db struct {

}

func (d db) Save() error {
	return nil
}

func New() Database {
	return &db{}
}