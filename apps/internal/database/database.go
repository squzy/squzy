package database

type Database interface {
}

type db struct {

}

func New() Database {
	return &db{}
}