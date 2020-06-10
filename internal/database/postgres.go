package database

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"strings"
	"time"
)

type postgres struct {
	db *gorm.DB
}

type Snapshot struct {
	gorm.Model
	SchedulerID   string `gorm:"column:schedulerId"`
	Code          string `gorm:"column:code"`
	Type          string `gorm:"column:type"`
	Error         string `gorm:"column:error"`
	MetaStartTime int64  `gorm:"column:metaStartTime"`
	MetaEndTime   int64  `gorm:"column:metaEndTime"`
	MetaValue     []byte `gorm:"column:metaValue"`
}

//Agent gorm description
type StatRequest struct {
	gorm.Model
	AgentID    string `gorm:"column:agentID"`
	AgentName  string `gorm:"column:agentName"`
	CPUInfo    []*CPUInfo
	MemoryInfo *MemoryInfo `gorm:"column:memoryInfo"`
	DiskInfo   []*DiskInfo `gorm:"column:diskInfo"`
	NetInfo    []*NetInfo  `gorm:"column:netInfo"`
	Time       time.Time   `gorm:"column:time"`
}

const (
	cpuInfoKey  = "CPUInfo"
	diskInfoKey = "DiskInfo"
	netInfoKey  = "NetInfo"
)

type CPUInfo struct {
	gorm.Model
	StatRequestID uint    `gorm:"column:statRequestId"`
	Load          float64 `gorm:"column:load"`
}

type MemoryInfo struct {
	gorm.Model
	StatRequestID uint        `gorm:"column:statRequestId"`
	Mem           *MemoryMem  `gorm:"column:mem"`
	Swap          *MemorySwap `gorm:"column:swap"`
}

type MemoryMem struct {
	gorm.Model
	MemoryInfoID uint    `gorm:"column:memoryInfoId"`
	Total        uint64  `gorm:"column:total"`
	Used         uint64  `gorm:"column:used"`
	Free         uint64  `gorm:"column:free"`
	Shared       uint64  `gorm:"column:shared"`
	UsedPercent  float64 `gorm:"column:usedPercent"`
}

type MemorySwap struct {
	gorm.Model
	MemoryInfoID uint    `gorm:"column:memoryInfoId"`
	Total        uint64  `gorm:"column:total"`
	Used         uint64  `gorm:"column:used"`
	Free         uint64  `gorm:"column:free"`
	Shared       uint64  `gorm:"column:shared"`
	UsedPercent  float64 `gorm:"column:usedPercent"`
}

type DiskInfo struct {
	gorm.Model
	StatRequestID uint    `gorm:"column:statRequestId"`
	Name          string  `gorm:"column:name"`
	Total         uint64  `gorm:"column:total"`
	Free          uint64  `gorm:"column:free"`
	Used          uint64  `gorm:"column:used"`
	UsedPercent   float64 `gorm:"column:usedPercent"`
}

type NetInfo struct {
	gorm.Model
	StatRequestID uint   `gorm:"column:statRequestId"`
	Name          string `gorm:"column:name"`
	BytesSent     uint64 `gorm:"column:bytesSent"`
	BytesRecv     uint64 `gorm:"column:bytesRecv"`
	PacketsSent   uint64 `gorm:"column:packetsSent"`
	PacketsRecv   uint64 `gorm:"column:packetsRecv"`
	ErrIn         uint64 `gorm:"column:errIn"`
	ErrOut        uint64 `gorm:"column:errOut"`
	DropIn        uint64 `gorm:"column:dropIn"`
	DropOut       uint64 `gorm:"column:dropOut"`
}

type TransactionInfo struct {
	gorm.Model
	TransactionId     string `gorm:"column:transactionId"`
	ApplicationId     string `gorm:"column:applicationId"`
	ParentId          string `gorm:"column:parentId"`
	MetaHost          string `gorm:"column:metaHost"`
	MetaPath          string `gorm:"column:metaPath"`
	MetaMethod        string `gorm:"column:metaMethod"`
	Name              string `gorm:"column:name"`
	StartTime         int64  `gorm:"column:startTime"`
	EndTime           int64  `gorm:"column:endTime"`
	TransactionStatus string `gorm:"column:transactionStatus"`
	TransactionType   string `gorm:"column:transactionType"`
	Error             string `gorm:"column:error"`
}

const (
	dbSnapshotCollection        = "snapshots"
	dbTransactionInfoCollection = "transaction_infos" //TODO: check
	dbStatRequestCollection     = "stat_requests"
)

var (
	errorDataBase = errors.New("ERROR_DATABASE_OPERATION")
)

func (p *postgres) Migrate() error {
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
	}

	var err error
	for _, model := range models {
		err = p.db.AutoMigrate(model).Error // migrate models one-by-one
	}

	return err
}

func (p *postgres) InsertSnapshot(data *apiPb.SchedulerResponse) error {
	snapshot, err := ConvertToPostgresSnapshot(data)
	if err != nil {
		return err
	}
	if err := p.db.Table(dbSnapshotCollection).Create(snapshot).Error; err != nil {
		return errorDataBase
	}
	return nil
}

func (p *postgres) GetSnapshots(request *apiPb.GetSchedulerInformationRequest) ([]*apiPb.SchedulerSnapshot, int32, error) {
	timeFrom, timeTo, err := getTime(request.GetTimeRange())
	if err != nil {
		return nil, -1, err
	}

	var count int64
	if request.GetStatus() == apiPb.SchedulerCode_SCHEDULER_CODE_UNSPECIFIED {
		err = p.db.Table(dbSnapshotCollection).
			Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), request.GetSchedulerId()).
			Where(fmt.Sprintf(`"%s"."metaStartTime" BETWEEN ? and ?`, dbSnapshotCollection), timeFrom, timeTo).
			Count(&count).Error
	} else {
		err = p.db.Table(dbSnapshotCollection).
			Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), request.GetSchedulerId()).
			Where(fmt.Sprintf(`"%s"."metaStartTime" BETWEEN ? and ?`, dbSnapshotCollection), timeFrom, timeTo).
			Where(fmt.Sprintf(`"%s"."code" = ?`, dbSnapshotCollection), request.GetStatus().String()).
			Count(&count).Error
	}

	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, request.GetPagination())

	var dbSnapshots []*Snapshot
	if request.GetStatus() == apiPb.SchedulerCode_SCHEDULER_CODE_UNSPECIFIED {
		err = p.db.
			Table(dbSnapshotCollection).
			Set("gorm:auto_preload", true).
			Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), request.GetSchedulerId()).
			Where(fmt.Sprintf(`"%s"."metaStartTime" BETWEEN ? and ?`, dbSnapshotCollection), timeFrom, timeTo).
			Order(getSnapshotOrder(request.GetSort()) + " " + getSnapshotDirection(request.GetSort())).
			Offset(offset).
			Limit(limit).
			Find(&dbSnapshots).Error
	} else {
		err = p.db.
			Table(dbSnapshotCollection).
			Set("gorm:auto_preload", true).
			Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), request.GetSchedulerId()).
			Where(fmt.Sprintf(`"%s"."metaStartTime" BETWEEN ? and ?`, dbSnapshotCollection), timeFrom, timeTo).
			Where(fmt.Sprintf(`"%s"."code" = ?`, dbSnapshotCollection), request.GetStatus().String()).
			Order(getSnapshotOrder(request.GetSort()) + getSnapshotDirection(request.GetSort())).
			Offset(offset).
			Limit(limit).
			Find(&dbSnapshots).Error
	}

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgresSnapshots(dbSnapshots), int32(count), nil
}

type UptimeResult struct {
	Count        int64  `gorm:"column:count"`
	Latency      string `gorm:"column:latency"`
}

func (p *postgres) GetSnapshotsUptime(request *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	timeFrom, timeTo, err := getTime(request.GetTimeRange())
	if err != nil {
		return nil, err
	}
	var countAll int64
	err = p.db.Table(dbSnapshotCollection).
		Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), request.GetSchedulerId()).
		Where(fmt.Sprintf(`"%s"."metaStartTime" BETWEEN ? and ?`, dbSnapshotCollection), timeFrom, timeTo).
		Count(&countAll).Error

	selectString := fmt.Sprintf(
		`COUNT(*) as "count", AVG("%s"."metaEndTime"-"%s"."metaStartTime") as "latency"`,
		dbSnapshotCollection,
		dbSnapshotCollection,
	)

	var uptimeResult UptimeResult
	err = p.db.Table(dbSnapshotCollection).
		Select(selectString).
		Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), request.GetSchedulerId()).
		Where(fmt.Sprintf(`"%s"."metaStartTime" BETWEEN ? and ?`, dbSnapshotCollection), timeFrom, timeTo).
		Where(fmt.Sprintf(`"%s"."code" = ?`, dbSnapshotCollection), "OK").
		Find(&uptimeResult).Error
	if err != nil {
		return nil, err
	}

	return convertFromUptimeResult(&uptimeResult, countAll), nil
}

func getSnapshotOrder(request *apiPb.SortingSchedulerList) string {
	if request == nil {
		return fmt.Sprintf(`"%s"."metaStartTime"`, dbSnapshotCollection)
	}
	orderMap := map[apiPb.SortSchedulerList]string{
		apiPb.SortSchedulerList_SORT_SCHEDULER_LIST_UNSPECIFIED: fmt.Sprintf(`"%s"."metaStartTime"`, dbSnapshotCollection),
		apiPb.SortSchedulerList_BY_START_TIME:                   fmt.Sprintf(`"%s"."metaStartTime"`, dbSnapshotCollection),
		apiPb.SortSchedulerList_BY_END_TIME:                     fmt.Sprintf(`"%s"."metaEndTime"`, dbSnapshotCollection),
		apiPb.SortSchedulerList_BY_LATENCY:                      fmt.Sprintf(`"%s"."metaEndTime" - "%s"."metaStartTime"`, dbSnapshotCollection, dbSnapshotCollection),
	}
	if res, ok := orderMap[request.GetSortBy()]; ok {
		return res
	}
	return fmt.Sprintf(`"%s"."metaStartTime"`, dbSnapshotCollection)
}

func getSnapshotDirection(request *apiPb.SortingSchedulerList) string {
	if request == nil {
		return ``
	}
	directionMap := map[apiPb.SortDirection]string{
		apiPb.SortDirection_SORT_DIRECTION_UNSPECIFIED: ``,
		apiPb.SortDirection_ASC:                        ` asc`,
		apiPb.SortDirection_DESC:                       ` desc`,
	}
	if res, ok := directionMap[request.GetDirection()]; ok {
		return res
	}
	return ``
}

//TODO: remake
func getUptimeAndLatency(dbSnapshots []*Snapshot, countAll, countOk int64) (float64, float64, error) {
	if countAll == 0 || countOk == 0 {
		return 0, 0, nil
	}
	latency := int64(0)
	/*for _, snapshot := range dbSnapshots {
		//Recieveing time difference in mellisecinds
		if snapshot.Meta != nil {
			latency += (snapshot.Meta.EndTime.UnixNano() - snapshot.Meta.StartTime.UnixNano()) / int64(time.Millisecond)
		}
	}*/
	return float64(countOk) / float64(countAll), float64(latency) / float64(countOk), nil
}

func (p *postgres) InsertStatRequest(data *apiPb.Metric) error {
	pgData, err := ConvertToPostgressStatRequest(data)
	if err != nil {
		return err
	}
	if err := p.db.Table(dbStatRequestCollection).Create(pgData).Error; err != nil {
		//TODO: log?
		return errorDataBase
	}
	return nil
}

func (p *postgres) GetStatRequest(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int64
	err = p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	//TODO: test if it works
	var statRequests []*StatRequest
	err = p.db.
		Set("gorm:auto_preload", true).
		//Preload("disk_infos").
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Order("time").
		Offset(offset).
		Limit(limit).
		Find(&statRequests).Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}

func (p *postgres) GetCPUInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, cpuInfoKey)
}

func (p *postgres) GetMemoryInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int64
	err = p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	//TODO: test if it works
	var statRequests []*StatRequest
	err = p.db.
		Preload("MemoryInfo").
		Preload("MemoryInfo.Mem").
		Preload("MemoryInfo.Swap").
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Order("time").
		Offset(offset).
		Limit(limit).
		Find(&statRequests).
		Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}

func (p *postgres) GetDiskInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, diskInfoKey)
}

func (p *postgres) GetNetInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, netInfoKey)
}

func (p *postgres) getSpecialRecords(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter, key string) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int64
	err = p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	//TODO: test if it works
	var statRequests []*StatRequest
	err = p.db.
		Preload(key).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Order("time").
		Offset(offset).
		Limit(limit).
		Find(&statRequests).
		Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}

func (p *postgres) InsertTransactionInfo(data *apiPb.TransactionInfo) error {
	info, err := convertToTransactionInfo(data)
	if err != nil {
		return err
	}
	if err := p.db.Table(dbTransactionInfoCollection).Create(info).Error; err != nil {
		return errorDataBase
	}
	return nil
}

func (p *postgres) GetTransactionInfo(request *apiPb.GetTransactionsRequest) ([]*apiPb.TransactionInfo, int64, error) {
	timeFrom, timeTo, err := getTimeInt64(request.TimeRange)
	if err != nil {
		return nil, -1, err
	}

	var count int64
	err = p.db.Table(dbTransactionInfoCollection).
		Where(fmt.Sprintf(`"%s"."applicationId" = ?`, dbTransactionInfoCollection), request.GetApplicationId()).
		Where(fmt.Sprintf(`"%s"."startTime" BETWEEN ? and ?`, dbTransactionInfoCollection), timeFrom, timeTo).
		Where(getTransactionsByString("metaHost", request.GetHost())).
		Where(getTransactionsByString("name", request.GetName())).
		Where(getTransactionsByString("metaPath", request.GetPath())).
		Where(getTransactionsByString("metaMethod", request.GetMethod())).
		Where(getTransactionTypeWhere(request.GetType())).
		Where(getTransactionStatusWhere(request.GetStatus())).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, request.GetPagination())

	//TODO: order
	var statRequests []*TransactionInfo
	err = p.db.Table(dbTransactionInfoCollection).
		Where(fmt.Sprintf(`"%s"."applicationId" = ?`, dbTransactionInfoCollection), request.GetApplicationId()).
		Where(fmt.Sprintf(`"%s"."startTime" BETWEEN ? and ?`, dbTransactionInfoCollection), timeFrom, timeTo).
		Where(getTransactionsByString("metaHost", request.GetHost())).
		Where(getTransactionsByString("name", request.GetName())).
		Where(getTransactionsByString("metaPath", request.GetPath())).
		Where(getTransactionsByString("metaMethod", request.GetMethod())).
		Where(getTransactionTypeWhere(request.GetType())).
		Where(getTransactionStatusWhere(request.GetStatus())).
		Order(getTransactionOrder(request.GetSort()) + getTransactionDirection(request.GetSort())). //TODO
		Offset(offset).
		Limit(limit).
		Find(&statRequests).
		Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return convertFromTransactions(statRequests), count, nil
}

func (p *postgres) GetTransactionByID(request *apiPb.GetTransactionByIdRequest) (*apiPb.TransactionInfo, []*apiPb.TransactionInfo, error) {
	var transaction TransactionInfo
	err := p.db.Table(dbTransactionInfoCollection).
		Where(fmt.Sprintf(`"%s"."transactionId" = ?`, dbTransactionInfoCollection), request.GetTransactionId()).
		First(&transaction).
		Error
	if err != nil || &transaction == nil {
		return nil, nil, errorDataBase
	}

	children, err := p.GetTransactionChildren(transaction.TransactionId, "")
	if err != nil {
		return nil, nil, err
	}

	return convertFromTransaction(&transaction), convertFromTransactions(children), nil
}

//passedString is used in order to prevent cycles
func (p *postgres) GetTransactionChildren(transactionId, passedString string) ([]*TransactionInfo, error) {
	if strings.Contains(passedString, transactionId) {
		return nil, nil
	}

	var childrenTransactions []*TransactionInfo
	err := p.db.Table(dbTransactionInfoCollection).
		Where(fmt.Sprintf(`"%s"."parentId" = ?`, dbTransactionInfoCollection), transactionId).
		Find(&childrenTransactions).
		Error
	if err != nil {
		return nil, errorDataBase
	}

	res := childrenTransactions
	for _, v := range childrenTransactions {
		subchildren, err := p.GetTransactionChildren(v.TransactionId, passedString+" "+v.ParentId)
		if err != nil {
			return nil, errorDataBase
		}
		for _, v := range subchildren {
			res = append(res, v)
		}
	}

	return res, nil
}

type GroupResult struct {
	Name         string `gorm:"column:groupName"`
	Count        int64  `gorm:"column:count"`
	SuccessCount int64  `gorm:"column:successCount"`
	Latency      string `gorm:"column:latency"`
	MinTime      string `gorm:"column:minTime"`
	MaxTime      string `gorm:"column:maxTime"`
}

func (p *postgres) GetTransactionGroup(request *apiPb.GetTransactionGroupRequest) (map[string]*apiPb.TransactionGroup, error) {
	timeFrom, timeTo, err := getTimeInt64(request.TimeRange)
	if err != nil {
		return nil, err
	}

	selectString := fmt.Sprintf(
		`%s as "groupName", COUNT(%s) as "count", COUNT(CASE WHEN "%s"."transactionStatus" = 'TRANSACTION_SUCCESSFUL' THEN 1 ELSE NULL END) as "successCount", AVG("%s"."endTime"-"%s"."startTime") as "latency", min("%s"."endTime"-"%s"."startTime") as "minTime", max("%s"."endTime"-"%s"."startTime") as "maxTime"`,
		getTransactionsGroupBy(request.GetGroupType()),
		getTransactionsGroupBy(request.GetGroupType()),
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
	)

	//TODO: order
	var groupResult []*GroupResult
	err = p.db.Table(dbTransactionInfoCollection).
		Select(selectString).
		Where(fmt.Sprintf(`"%s"."applicationId" = ?`, dbTransactionInfoCollection), request.GetApplicationId()).
		Where(fmt.Sprintf(`"%s"."startTime" BETWEEN ? and ?`, dbTransactionInfoCollection), timeFrom, timeTo).
		Where(getTransactionTypeWhere(request.GetType())).
		Where(getTransactionStatusWhere(request.GetStatus())).
		Group(getTransactionsGroupBy(request.GetGroupType())).
		Find(&groupResult).
		Error
	if err != nil {
		return nil, errorDataBase
	}

	return convertFromGroupResult(groupResult), nil
}

func getTransactionOrder(request *apiPb.SortingTransactionList) string {
	if request == nil {
		return fmt.Sprintf(`"%s"."startTime"`, dbTransactionInfoCollection)
	}
	orderMap := map[apiPb.SortTransactionList]string{
		apiPb.SortTransactionList_SORT_TRANSACTION_LIST_UNSPECIFIED: fmt.Sprintf(`"%s"."startTime"`, dbTransactionInfoCollection),
		apiPb.SortTransactionList_DURATION:                          fmt.Sprintf(`"%s"."endTime" - "%s"."startTime"`, dbTransactionInfoCollection, dbTransactionInfoCollection),
	}
	if res, ok := orderMap[request.GetSortBy()]; ok {
		return res
	}
	return fmt.Sprintf(`"%s"."startTime"`, dbTransactionInfoCollection)
}

func getTransactionDirection(request *apiPb.SortingTransactionList) string {
	if request == nil {
		return ``
	}
	directionMap := map[apiPb.SortDirection]string{
		apiPb.SortDirection_SORT_DIRECTION_UNSPECIFIED: ``,
		apiPb.SortDirection_ASC:                        ` asc`,
		apiPb.SortDirection_DESC:                       ` desc`,
	}
	if res, ok := directionMap[request.GetDirection()]; ok {
		return res
	}
	return ``
}

func getTransactionsByString(key string, value *wrappers.StringValue) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf(`"%s"."%s" = '%s'`, dbTransactionInfoCollection, key, value.GetValue())
}

func getTransactionTypeWhere(transType apiPb.TransactionType) string {
	if transType == apiPb.TransactionType_TRANSACTION_TYPE_UNSPECIFIED {
		return ""
	}
	return fmt.Sprintf(`"%s"."transactionType" = '%s'`, dbTransactionInfoCollection, transType.String())
}

func getTransactionStatusWhere(transType apiPb.TransactionStatus) string {
	if transType == apiPb.TransactionStatus_TRANSACTION_CODE_UNSPECIFIED {
		return ""
	}
	return fmt.Sprintf(`"%s"."transactionStatus" = '%s'`, dbTransactionInfoCollection, transType.String())
}

var (
	groupMap = map[apiPb.GroupTransaction]string{
		apiPb.GroupTransaction_GROUP_TRANSACTION_UNSPECIFIED: "transactionType",
		apiPb.GroupTransaction_BY_TYPE:                       "transactionType",
		apiPb.GroupTransaction_BY_NAME:                       "name",
		apiPb.GroupTransaction_BY_METHOD:                     "metaMethod",
		apiPb.GroupTransaction_BY_HOST:                       "metaHost",
		apiPb.GroupTransaction_BY_PATH:                       "metaPath",
	}
)

func getTransactionsGroupBy(group apiPb.GroupTransaction) string {
	if val, ok := groupMap[group]; ok {
		return fmt.Sprintf(`"%s"."%s"`, dbTransactionInfoCollection, val)
	}
	return fmt.Sprintf(`"%s"."transactionType"`, dbTransactionInfoCollection)
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
