package postgres

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"time"
)

type Postgres struct {
	Db *gorm.DB
}

const (
	dbSnapshotCollection        = "snapshots"
	dbTransactionInfoCollection = "transaction_infos"
	dbStatRequestCollection     = "stat_requests"
)

var (
	errorDataBase = errors.New("ERROR_DATABASE_OPERATION")

	directionMap = map[apiPb.SortDirection]string{
		apiPb.SortDirection_SORT_DIRECTION_UNSPECIFIED: ``,
		apiPb.SortDirection_ASC:                        ` asc`,
		apiPb.SortDirection_DESC:                       ` desc`,
	}
)

func (p *Postgres) Migrate() error {
	models := []interface{}{
		&Snapshot{},
		&StatRequest{},
		&CPUInfo{},
		&MemoryInfo{},
		&MemoryMem{},
		&MemorySwap{},
		&DiskInfo{},
		&NetInfo{},
		&TransactionInfo{},
		&Incident{},
		&IncidentHistory{},
	}

	var err error
	for _, model := range models {
		err = p.Db.AutoMigrate(model).Error // migrate models one-by-one
	}

	return err
}

func getTime(filter *apiPb.TimeFilter) (time.Time, time.Time, error) {
	timeFrom := time.Unix(0, 0)
	timeTo := time.Now()
	var err error
	if filter != nil {
		if filter.GetFrom() != nil {
			timeFrom = filter.From.AsTime()
			err = filter.From.CheckValid()
		}
		if filter.GetTo() != nil {
			timeTo = filter.To.AsTime()
			err = filter.To.CheckValid()
		}
	}
	return timeFrom, timeTo, err
}

//Return time unixNanos
func getTimeInt64(filter *apiPb.TimeFilter) (int64, int64, error) {
	timeFrom := time.Unix(0, 0)
	timeTo := time.Now()
	var err error
	if filter != nil {
		if filter.GetFrom() != nil {
			timeFrom = filter.From.AsTime()
			err = filter.From.CheckValid()
		}
		if filter.GetTo() != nil {
			timeTo = filter.To.AsTime()
			err = filter.To.CheckValid()
		}
	}
	if err != nil {
		return 0, 0, err
	}
	return timeFrom.UnixNano(), timeTo.UnixNano(), nil
}

//Return offset and limit
func getOffsetAndLimit(count int64, pagination *apiPb.Pagination) (int, int) {
	if pagination != nil {
		if pagination.Page == -1 {
			return int(count) - int(pagination.Limit), int(pagination.Limit)
		}
		return int(pagination.GetLimit() * (pagination.GetPage() - 1)), int(pagination.GetLimit())
	}
	return 0, int(count)
}
