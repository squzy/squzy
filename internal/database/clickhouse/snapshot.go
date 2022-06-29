package clickhouse

import (
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	uuid "github.com/google/uuid"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"time"
)

type Snapshot struct {
	Model         Model
	SchedulerID   string
	Code          int32
	Type          int32
	Error         string
	MetaStartTime int64
	MetaEndTime   int64
	MetaValue     []byte
}

type UptimeResult struct {
	Count   int64
	Latency string
}

var (
	snapshotFields                    = "id, created_at, updated_at, scheduler_id, code, type, error, meta_start_time, meta_end_time, meta_value"
	snapshotSchedulerIdString         = fmt.Sprintf(`"scheduler_id" = ?`)
	snapshotMetaStartTimeFilterString = fmt.Sprintf(`"meta_start_time" BETWEEN ? and ?`)

	snapOrderMap = map[apiPb.SortSchedulerList]string{
		apiPb.SortSchedulerList_SORT_SCHEDULER_LIST_UNSPECIFIED: fmt.Sprintf(`"%s"."meta_start_time"`, dbSnapshotCollection),
		apiPb.SortSchedulerList_BY_START_TIME:                   fmt.Sprintf(`"%s"."meta_start_time"`, dbSnapshotCollection),
		apiPb.SortSchedulerList_BY_END_TIME:                     fmt.Sprintf(`"%s"."meta_end_time"`, dbSnapshotCollection),
		apiPb.SortSchedulerList_BY_LATENCY:                      fmt.Sprintf(`"%s"."meta_end_time" - "%s"."meta_start_time"`, dbSnapshotCollection, dbSnapshotCollection),
	}
)

func (c *Clickhouse) InsertSnapshot(data *apiPb.SchedulerResponse) error {
	now := time.Now()

	snapshot, err := ConvertToSnapshot(data)
	if err != nil {
		return err
	}

	err = c.insertSnapshot(now, snapshot)
	if err != nil {
		logger.Error(err.Error())
		return errorDataBase
	}
	return nil
}

func (c *Clickhouse) insertSnapshot(now time.Time, snapshot *Snapshot) error {
	tx, err := c.Db.Begin()
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES ($0, $1, $2, $3, $4, $5, $6, $7, $8, $9)`, dbSnapshotCollection, snapshotFields)
	_, err = tx.Exec(q,
		clickhouse.UUID(uuid.New().String()),
		now,
		now,
		snapshot.SchedulerID,
		snapshot.Code,
		snapshot.Type,
		snapshot.Error,
		snapshot.MetaStartTime,
		snapshot.MetaEndTime,
		snapshot.MetaValue,
	)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (c *Clickhouse) GetSnapshots(request *apiPb.GetSchedulerInformationRequest) ([]*apiPb.SchedulerSnapshot, int32, error) {
	timeFrom, timeTo, err := getTimeInt64(request.GetTimeRange())
	if err != nil {
		return nil, -1, err
	}

	var count int64
	count, err = c.countSnapshots(request, timeFrom, timeTo)
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, request.GetPagination())

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM snapshots WHERE (%s AND %s %s) ORDER BY %s LIMIT %d OFFSET %d`,
		snapshotFields,
		snapshotSchedulerIdString,
		getCodeString(request.GetStatus(), andSep),
		snapshotMetaStartTimeFilterString,
		getSnapshotOrder(request.GetSort())+getSnapshotDirection(request.GetSort()),
		limit,
		offset),
		request.SchedulerId,
		timeFrom,
		timeTo,
	)

	if err != nil {
		logger.Error(err.Error())
		return nil, -1, errorDataBase
	}
	defer rows.Close()

	var snapshots []*Snapshot
	for rows.Next() {
		snp := &Snapshot{}
		if err := rows.Scan(&snp.Model.ID, &snp.Model.CreatedAt, &snp.Model.UpdatedAt,
			&snp.SchedulerID, &snp.Code, &snp.Type, &snp.Error,
			&snp.MetaStartTime, &snp.MetaEndTime, &snp.MetaValue); err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}

		snapshots = append(snapshots, snp)
	}

	return ConvertFromSnapshots(snapshots), int32(count), nil
}

func (c *Clickhouse) countSnapshots(request *apiPb.GetSchedulerInformationRequest, timeFrom int64, timeTo int64) (int64, error) {
	var count int64
	rows, err := c.Db.Query(fmt.Sprintf(`SELECT count(*) FROM "%s" WHERE %s AND %s %s LIMIT 1`,
		dbSnapshotCollection,
		snapshotSchedulerIdString,
		getCodeString(request.Status, andSep),
		snapshotMetaStartTimeFilterString),
		request.SchedulerId,
		timeFrom,
		timeTo)

	if err != nil {
		logger.Error(err.Error())
		return -1, errorDataBase
	}

	defer rows.Close()

	if ok := rows.Next(); !ok {
		return 0, nil
	}

	if err := rows.Scan(&count); err != nil {
		logger.Error(err.Error())
		return -1, errorDataBase
	}
	return count, nil
}

func (c *Clickhouse) GetSnapshotsUptime(request *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	timeFrom, timeTo, err := getTimeInt64(request.GetTimeRange())
	if err != nil {
		return nil, err
	}
	countAll, err := c.countAllSnapshots(request, timeFrom, timeTo)
	if err != nil {
		return nil, err
	}

	var uptimeResult UptimeResult
	uptimeResult, err = c.countSnapshotsUptime(request, timeFrom, timeTo)
	if err != nil {
		return nil, err
	}
	return convertFromUptimeResult(&uptimeResult, countAll), nil
}

func (c *Clickhouse) countAllSnapshots(request *apiPb.GetSchedulerUptimeRequest, timeFrom int64, timeTo int64) (int64, error) {
	var count int64
	rows, err := c.Db.Query(fmt.Sprintf(`SELECT count(*) FROM "%s" WHERE %s AND (%s)`,
		dbSnapshotCollection,
		snapshotSchedulerIdString,
		snapshotMetaStartTimeFilterString),
		request.SchedulerId,
		timeFrom,
		timeTo)

	if err != nil {
		logger.Error(err.Error())
		return -1, errorDataBase
	}

	defer rows.Close()

	if ok := rows.Next(); !ok {
		return 0, nil
	}

	if err := rows.Scan(&count); err != nil {
		logger.Error(err.Error())
		return -1, errorDataBase
	}

	return count, nil
}

func (c *Clickhouse) countSnapshotsUptime(request *apiPb.GetSchedulerUptimeRequest, timeFrom int64, timeTo int64) (UptimeResult, error) {
	var result UptimeResult

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT count(*) as "count", avg(meta_end_time-meta_start_time) as "latency" FROM "%s" WHERE %s AND %s %s`,
		dbSnapshotCollection,
		snapshotSchedulerIdString,
		getCodeString(apiPb.SchedulerCode_OK, andSep),
		snapshotMetaStartTimeFilterString),
		request.SchedulerId,
		timeFrom,
		timeTo)

	if err != nil {
		logger.Error(err.Error())
		return UptimeResult{
			Count:   -1,
			Latency: "",
		}, errorDataBase
	}

	defer rows.Close()

	if ok := rows.Next(); !ok {
		return UptimeResult{
			Count:   0,
			Latency: "",
		}, nil
	}

	if err := rows.Scan(&result.Count, &result.Latency); err != nil {
		logger.Error(err.Error())
		return UptimeResult{
			Count:   -1,
			Latency: "",
		}, errorDataBase
	}

	return result, nil
}

func getCodeString(code apiPb.SchedulerCode, separator string) string {
	if code == apiPb.SchedulerCode_SCHEDULER_CODE_UNSPECIFIED {
		return ""
	}
	return fmt.Sprintf(`"code" = '%d' %s`, code, separator)
}

func getSnapshotOrder(request *apiPb.SortingSchedulerList) string {
	if request == nil {
		return fmt.Sprintf(`"meta_start_time"`)
	}
	if res, ok := snapOrderMap[request.GetSortBy()]; ok {
		return res
	}
	return fmt.Sprintf(`"meta_start_time"`)
}

func getSnapshotDirection(request *apiPb.SortingSchedulerList) string {
	if request == nil {
		return descPrefix
	}
	if res, ok := directionMap[request.GetDirection()]; ok {
		return res
	}
	return descPrefix
}
