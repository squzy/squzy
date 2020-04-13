package postgres_connection

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"squzy/apps/server/structures"
	"sync"
)

type Postgres interface {
	Migrate() []error
	InsertMetaData(data *structures.MetaData)
	GetMetaData(id string) *structures.MetaData
	InsertStatRequest(data *structures.StatRequest)
	GetStatRequest(id string) *structures.StatRequest
}

type postgres struct {
	initClient sync.Once
	db         *gorm.DB

	host     string
	port     string
	user     string
	password string
	dbname   string
}

const (
	dbMetaDataCollection    = "metaData"    //TODO: check
	dbStatRequestCollection = "statRequest" //TODO: check
)

func NewPostgressConnection(host, port, user, dbname, password string) Postgres {
	return &postgres{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

func (p *postgres) newClient() {
	var err error
	p.db, err = gorm.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s connect_timeout=10 sslmode=disable",
			p.host,
			p.port,
			p.user,
			p.dbname,
			p.password,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	p.db.LogMode(true)
}

func (p *postgres) getClient() *gorm.DB {
	p.initClient.Do(p.newClient)
	return p.db
}

func (p *postgres) Migrate() []error {
	var errs []error
	err := p.db.DB().Ping() // ping the pg
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	models := [2]interface{}{
		structures.MetaData{},
		structures.StatRequest{},
	}

	for _, model := range models {
		err = p.db.AutoMigrate(model).Error // migrate models one-by-one
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (p *postgres) InsertMetaData(data *structures.MetaData) {
	db := p.getClient()
	db.Table(dbMetaDataCollection).Create(data)
}

func (p *postgres) GetMetaData(id string) *structures.MetaData {
	db := p.getClient()
	metaData := &structures.MetaData{}
	db.Table(dbMetaDataCollection).Where(fmt.Sprintf(`"%s"."id" = ?`,dbMetaDataCollection), id).First(metaData)
	return metaData
}

func (p *postgres) InsertStatRequest(data *structures.StatRequest) {
	db := p.getClient()
	db.Table(dbStatRequestCollection).Create(data)
}

func (p *postgres) GetStatRequest(id string) *structures.StatRequest {
	db := p.getClient()
	statRequest := &structures.StatRequest{}
	db.Table(dbStatRequestCollection).Where(fmt.Sprintf(`"%s"."id" = ?`,dbStatRequestCollection), id).First(statRequest)
	return statRequest
}
