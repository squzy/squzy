package clickhouse

import (
	"database/sql"
	"errors"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

type Clickhouse struct {
	Db *sql.DB
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

func (c *Clickhouse) Migrate() error {
	_, err := c.Db.Exec(`
			CREATE TABLE IF NOT EXISTS snapshot (
				SchedulerId String,
				Code Int32,
				Type        Int32,
				Error   String,
				MetaStartTime   Int64,
				MetaEndTime   Int64,
				MetaValue  Array(Int8)
			)
		`)
	if err != nil {
		return err
	}

	_, err = c.Db.Exec(`
			CREATE TABLE IF NOT EXISTS stat_request (
				AgentId String,
				AgentName String,
				CPUInfo Nested
					(
						StatRequestID UInt32,
						Load Float64
					),
				MemoryInfoMem Nested
					(
						MemoryInfoID UInt32,
						Total UInt64,
						Used UInt64,
						Free UInt64,
						Shared UInt64,
						UsedPercent Float64
					),
				MemoryInfoSwap Nested
					(
						MemoryInfoID UInt32,
						Total UInt64,
						Used UInt64,
						Free UInt64,
						Shared UInt64,
						UsedPercent Float64
					),		
				DiskInfo Nested
					(
						MemoryInfoID UInt32,
						Total UInt64,
						Used UInt64,
						Free UInt64,
						Shared UInt64,
						UsedPercent Float64
					),		
				NetInfo Nested
					(
						StatRequestID UInt32,
						Name String,
						BytesSent UInt64,
						BytesRecv UInt64,
						PacketsSent UInt64,
						PacketsRecv UInt64,
						ErrIn UInt64,
						ErrOut UInt64,
						DropIn UInt64,
						DropOut UInt64
					),		
				Time DateTime
			)
		`)
	if err != nil {
		return err
	}

	_, err = c.Db.Exec(`
		CREATE TABLE IF NOT EXISTS transaction_info (
					TransactionId String,
					ApplicationId String,
					ParentId String,
					MetaHost String,
					MetaPath String,
					MetaMethod String,
					Name String,
					StartTime Int64,
					EndTime Int64,
					TransactionStatus Int32,
					TransactionType Int32,
					Error String
				)
			`)
	if err != nil {
		return err
	}

	_, err = c.Db.Exec(`
			CREATE TABLE IF NOT EXISTS incident (
				IncidentId String,
				Status Int32,
				RuleId String,
				StartTime Int64,
				EndTime Int64,
				Histories Nested
					(
						IncidentID UInt32,
						Status Int64,
						Timestamp Int64
					),		
				Time DateTime
			)
		`)
	if err != nil {
		return err
	}
	return nil
}

