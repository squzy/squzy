package clickhouse

import (
	_ "github.com/ClickHouse/clickhouse-go"

	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

//docker run -d --rm --name clickhouse-server -p 8123:8123 --ulimit nofile=262144:262144 yandex/clickhouse-server
var (
	db, _   = sql.Open("clickhouse", "http://user:password@lkl:00/clicks?read_timeout=10&write_timeout=10")
	chWrong = &Clickhouse{
		db,
	}
)

func TestPostgres_Migrate_error(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		err := chWrong.Migrate()
		assert.Error(t, err)
	})
}
