package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
)

type Snapshot struct {
	gorm.Model
	SchedulerID   string `gorm:"column:schedulerId"`
	Code          int32  `gorm:"column:code"`
	Type          int32  `gorm:"column:type"`
	Error         string `gorm:"column:error"`
	MetaStartTime int64  `gorm:"column:metaStartTime"`
	MetaEndTime   int64  `gorm:"column:metaEndTime"`
	MetaValue     []byte `gorm:"column:metaValue"`
}

type UptimeResult struct {
	Count   int64  `gorm:"column:count"`
	Latency string `gorm:"column:latency"`
}

var (
	schedulerIdFilterString   = fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection)
	metaStartTimeFilterString = fmt.Sprintf(`"%s"."metaStartTime" BETWEEN ? and ?`, dbSnapshotCollection)

	snapOrderMap = map[apiPb.SortSchedulerList]string{
		apiPb.SortSchedulerList_SORT_SCHEDULER_LIST_UNSPECIFIED: fmt.Sprintf(`"%s"."metaStartTime"`, dbSnapshotCollection),
		apiPb.SortSchedulerList_BY_START_TIME:                   fmt.Sprintf(`"%s"."metaStartTime"`, dbSnapshotCollection),
		apiPb.SortSchedulerList_BY_END_TIME:                     fmt.Sprintf(`"%s"."metaEndTime"`, dbSnapshotCollection),
		apiPb.SortSchedulerList_BY_LATENCY:                      fmt.Sprintf(`"%s"."metaEndTime" - "%s"."metaStartTime"`, dbSnapshotCollection, dbSnapshotCollection),
	}
)

func (p *Postgres) InsertSnapshot(data *apiPb.SchedulerResponse) error {
	snapshot, err := ConvertToPostgresSnapshot(data)
	if err != nil {
		return err
	}
	if err := p.Db.Table(dbSnapshotCollection).Create(snapshot).Error; err != nil {
		return errorDataBase
	}
	return nil
}

func (p *Postgres) GetSnapshots(request *apiPb.GetSchedulerInformationRequest) ([]*apiPb.SchedulerSnapshot, int32, error) {
	timeFrom, timeTo, err := getTimeInt64(request.GetTimeRange())
	if err != nil {
		return nil, -1, err
	}

	var count int64
	err = p.Db.Table(dbSnapshotCollection).
		Where(schedulerIdFilterString, request.GetSchedulerId()).
		Where(metaStartTimeFilterString, timeFrom, timeTo).
		Where(getCodeString(request.GetStatus())).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, request.GetPagination())

	var dbSnapshots []*Snapshot
	err = p.Db.
		Table(dbSnapshotCollection).
		Set("gorm:auto_preload", true).
		Where(schedulerIdFilterString, request.GetSchedulerId()).
		Where(metaStartTimeFilterString, timeFrom, timeTo).
		Where(getCodeString(request.GetStatus())).
		Order(getSnapshotOrder(request.GetSort()) + getSnapshotDirection(request.GetSort())).
		Offset(offset).
		Limit(limit).
		Find(&dbSnapshots).Error
	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgresSnapshots(dbSnapshots), int32(count), nil
}

func (p *Postgres) GetSnapshotsUptime(request *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	timeFrom, timeTo, err := getTimeInt64(request.GetTimeRange())
	if err != nil {
		return nil, err
	}
	var countAll int64
	err = p.Db.Table(dbSnapshotCollection).
		Where(schedulerIdFilterString, request.GetSchedulerId()).
		Where(metaStartTimeFilterString, timeFrom, timeTo).
		Count(&countAll).Error

	if err != nil {
		return nil, err
	}

	selectString := fmt.Sprintf(
		`COUNT(*) as "count", AVG("%s"."metaEndTime"-"%s"."metaStartTime") as "latency"`,
		dbSnapshotCollection,
		dbSnapshotCollection,
	)

	var uptimeResult UptimeResult
	err = p.Db.Table(dbSnapshotCollection).
		Select(selectString).
		Where(schedulerIdFilterString, request.GetSchedulerId()).
		Where(metaStartTimeFilterString, timeFrom, timeTo).
		Where(getCodeString(apiPb.SchedulerCode_OK)).
		Find(&uptimeResult).Error
	if err != nil {
		return nil, err
	}
	return convertFromUptimeResult(&uptimeResult, countAll), nil
}

func getCodeString(code apiPb.SchedulerCode) string {
	if code == apiPb.SchedulerCode_SCHEDULER_CODE_UNSPECIFIED {
		return ""
	}
	return fmt.Sprintf(`"%s"."code" = '%d'`, dbSnapshotCollection, code)
}

func getSnapshotOrder(request *apiPb.SortingSchedulerList) string {
	if request == nil {
		return fmt.Sprintf(`"%s"."metaStartTime"`, dbSnapshotCollection)
	}
	if res, ok := snapOrderMap[request.GetSortBy()]; ok {
		return res
	}
	return fmt.Sprintf(`"%s"."metaStartTime"`, dbSnapshotCollection)
}

func getSnapshotDirection(request *apiPb.SortingSchedulerList) string {
	if request == nil {
		return ` desc`
	}
	if res, ok := directionMap[request.GetDirection()]; ok {
		return res
	}
	return ` desc`
}
