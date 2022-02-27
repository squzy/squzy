package clickhouse

import (
	"database/sql"
	"errors"
	"github.com/golang/protobuf/ptypes"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"time"
)

type Clickhouse struct {
	Db *sql.DB
}

const (
	dbSnapshotCollection                  = "snapshots"
	dbTransactionInfoCollection           = "transaction_info"
	dbStatRequestCollection               = "stat_requests"
	dbStatRequestCpuInfoCollection        = "stat_requests_cpu_info"
	dbStatRequestMemoryInfoMemCollection  = "stat_requests_memory_info_mem"
	dbStatRequestMemoryInfoSwapCollection = "stat_requests_memory_info_swap"
	dbStatRequestDiskInfoCollection       = "stat_requests_disk_info"
	dbStatRequestNetInfoCollection        = "stat_requests_net_info"
)

var (
	errorDataBase = errors.New("ERROR_DATABASE_OPERATION")

	directionMap = map[apiPb.SortDirection]string{
		apiPb.SortDirection_SORT_DIRECTION_UNSPECIFIED: ``,
		apiPb.SortDirection_ASC:                        ` ASC`,
		apiPb.SortDirection_DESC:                       ` DESC`,
	}
)

type Model struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Clickhouse) Migrate() error {
	_, err := c.Db.Exec(`
			CREATE TABLE IF NOT EXISTS snapshots (
				id UUID,
				created_at DateTime,
				updated_at DateTime,
				scheduler_id String,
				code Int32,
				type        Int32,
				error   String,
				meta_start_time   Int64,
				meta_end_time   Int64,
				meta_value  Array(UInt8)
			) ENGINE = MergeTree ORDER BY tuple()
		`)
	if err != nil {
		return err
	}

	_, err = c.Db.Exec(`
			CREATE TABLE IF NOT EXISTS stat_requests (
				id UUID,
				created_at DateTime,
				updated_at DateTime,
				agent_id String,
				agent_name String,
				time DateTime
			) ENGINE = MergeTree ORDER BY tuple()
		`)
	if err != nil {
		return err
	}

	_, err = c.Db.Exec(`
			CREATE TABLE IF NOT EXISTS stat_requests_cpu_info (
				id UUID,
				created_at DateTime,
				updated_at DateTime,
				stat_request_id UUID,
				load Float64
			) ENGINE = MergeTree ORDER BY tuple()
		`)
	if err != nil {
		return err
	}

	_, err = c.Db.Exec(`
			CREATE TABLE IF NOT EXISTS stat_requests_memory_info_mem (
				id UUID,
				created_at DateTime,
				updated_at DateTime,
				stat_request_id UUID,
				memory_info_id UUID,
				total UInt64,
				used UInt64,
				free UInt64,
				shared UInt64,
				used_percent Float64
			) ENGINE = MergeTree ORDER BY tuple()
		`)
	if err != nil {
		return err
	}
	_, err = c.Db.Exec(`
			CREATE TABLE IF NOT EXISTS stat_requests_memory_info_swap (
				id UUID,
				created_at DateTime,
				updated_at DateTime,
				stat_request_id UUID,
				memory_info_id UUID,
				total UInt64,
				used UInt64,
				free UInt64,
				shared UInt64,
				used_percent Float64
			) ENGINE = MergeTree ORDER BY tuple()
		`)
	if err != nil {
		return err
	}

	_, err = c.Db.Exec(`
			CREATE TABLE IF NOT EXISTS stat_requests_disk_info (
				id UUID,
				created_at DateTime,
				updated_at DateTime,
				stat_request_id UUID,
				memory_info_id UUID,
				name String,
				total UInt64,
				used UInt64,
				free UInt64,
				used_percent Float64
			) ENGINE = MergeTree ORDER BY tuple()
		`)
	if err != nil {
		return err
	}

	_, err = c.Db.Exec(`
			CREATE TABLE IF NOT EXISTS stat_requests_net_info (
				id UUID,
				created_at DateTime,
				updated_at DateTime,
				stat_request_id UUID,
				name String,
				bytes_sent UInt64,
				bytes_recv UInt64,
				packets_sent UInt64,
				packets_recv UInt64,
				err_in UInt64,
				err_out UInt64,
				drop_in UInt64,
				drop_out UInt64
			) ENGINE = MergeTree ORDER BY tuple();
		`)
	if err != nil {
		return err
	}

	_, err = c.Db.Exec(`
		CREATE TABLE IF NOT EXISTS transaction_info (
					id UUID,
			    	created_at DateTime,
			    	updated_at DateTime,
					transaction_id String,
					application_id String,
					parent_id String,
					meta_host String,
					meta_path String,
					meta_method String,
					name String,
					start_time Int64,
					end_time Int64,
					transaction_status Int32,
					transaction_type Int32,
					error String
				) ENGINE = MergeTree ORDER BY tuple()
			`)
	if err != nil {
		return err
	}

	_, err = c.Db.Exec(`
			CREATE TABLE IF NOT EXISTS incidents (
			    id UUID,
			    created_at DateTime,
			    updated_at DateTime,
				incident_id String,
				status Int32,
				rule_id String,
				start_time Int64,
				end_time Int64
			) ENGINE = ReplacingMergeTree(updated_at) ORDER BY tuple(incident_id)
		`)
	if err != nil {
		return err
	}

	_, err = c.Db.Exec(`
			CREATE TABLE IF NOT EXISTS incidents_history (
			    id UUID,
			    created_at DateTime,
			    updated_at DateTime,
				timestamp Int64,
				incident_id String,
				status Int32
			) ENGINE = StripeLog
		`)
	if err != nil {
		return err
	}
	return nil
}

func getTime(filter *apiPb.TimeFilter) (time.Time, time.Time, error) {
	timeFrom := time.Unix(0, 0)
	timeTo := time.Now()
	var err error
	if filter != nil {
		if filter.GetFrom() != nil {
			timeFrom, err = ptypes.Timestamp(filter.From)
		}
		if filter.GetTo() != nil {
			timeTo, err = ptypes.Timestamp(filter.To)
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
			timeFrom, err = ptypes.Timestamp(filter.From)
		}
		if filter.GetTo() != nil {
			timeTo, err = ptypes.Timestamp(filter.To)
		}
	}
	if err != nil {
		return 0, 0, err
	}
	return timeFrom.UnixNano(), timeTo.UnixNano(), nil
}

func getOffsetAndLimit(count int64, pagination *apiPb.Pagination) (int, int) {
	if pagination != nil {
		if pagination.Page == -1 {
			return int(count) - int(pagination.Limit), int(pagination.Limit)
		}
		return int(pagination.GetLimit() * (pagination.GetPage() - 1)), int(pagination.GetLimit())
	}
	return 0, int(count)
}
