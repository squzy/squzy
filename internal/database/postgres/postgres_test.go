package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

//docker run -d --rm --name postgres -e POSTGRES_USER="user" -e POSTGRES_PASSWORD="password" -e POSTGRES_DB="database" -p 5432:5432 postgres
var (
	db, _  = gorm.Open(
		"postgres",
		fmt.Sprintf("host=lkl port=00 user=us dbname=dbn password=ps connect_timeout=10 sslmode=disable"))
	postgrWrong = &Postgres{
		db,
	}
)

func TestPostgres_Migrate_error(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		err := postgrWrong.Migrate()
		assert.Error(t, err)
	})
}
