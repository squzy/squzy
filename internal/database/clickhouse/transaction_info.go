package clickhouse

//
import (
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/golang/protobuf/ptypes/wrappers"
	uuid "github.com/satori/go.uuid"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"squzy/internal/logger"
	"strings"
	"time"
)

type TransactionInfo struct {
	Model
	TransactionId     string
	ApplicationId     string
	ParentId          string
	MetaHost          string
	MetaPath          string
	MetaMethod        string
	Name              string
	StartTime         int64
	EndTime           int64
	TransactionStatus int32
	TransactionType   int32
	Error             string
}

type GroupResult struct {
	Name         string
	Count        int64
	SuccessCount int64
	Latency      string
	MinTime      string
	MaxTime      string
	LowTime      string
}

const (
	transactionInfoFields   = "id, created_at, updated_at, transaction_id, application_id, parent_id, meta_host, meta_path, meta_method, name, start_time, end_time, transactionStatus, transaction_type, error"
	transNameStr            = "name"
	transMetaHostStr        = "metaHost"
	transMetaMethodStr      = "metaMethod"
	transMetaPathStr        = "metaPath"
	transTransactionTypeStr = "transactionType"
)

var (
	applicationIdFilterString        = fmt.Sprintf(`"applicationId" = ?`)
	applicationStartTimeFilterString = fmt.Sprintf(`"startTime" BETWEEN ? and ?`)

	transOrderMap = map[apiPb.SortTransactionList]string{
		apiPb.SortTransactionList_SORT_TRANSACTION_LIST_UNSPECIFIED: fmt.Sprintf(`"%s"."startTime"`, dbTransactionInfoCollection),
		apiPb.SortTransactionList_DURATION:                          fmt.Sprintf(`"%s"."endTime" - "%s"."startTime"`, dbTransactionInfoCollection, dbTransactionInfoCollection),
		apiPb.SortTransactionList_BY_TRANSACTION_START_TIME:         fmt.Sprintf(`"%s"."startTime"`, dbTransactionInfoCollection),
		apiPb.SortTransactionList_BY_TRANSACTION_END_TIME:           fmt.Sprintf(`"%s"."endTime"`, dbTransactionInfoCollection),
	}
	groupMap = map[apiPb.GroupTransaction]string{
		apiPb.GroupTransaction_GROUP_TRANSACTION_UNSPECIFIED: transTransactionTypeStr,
		apiPb.GroupTransaction_BY_TYPE:                       transTransactionTypeStr,
		apiPb.GroupTransaction_BY_NAME:                       transNameStr,
		apiPb.GroupTransaction_BY_METHOD:                     transMetaMethodStr,
		apiPb.GroupTransaction_BY_HOST:                       transMetaHostStr,
		apiPb.GroupTransaction_BY_PATH:                       transMetaPathStr,
	}
)

func (c *Clickhouse) InsertTransactionInfo(data *apiPb.TransactionInfo) error {
	now := time.Now()

	info, err := convertToTransactionInfo(data)
	if err != nil {
		return err
	}

	tx, err := c.Db.Begin()
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`INSERT INTO transaction_info (%s) VALUES ($0, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`, transactionInfoFields)
	_, err = tx.Exec(q,
		clickhouse.UUID(uuid.NewV4().String()),
		now,
		now,
		info.TransactionId,
		info.ApplicationId,
		info.ParentId,
		info.MetaHost,
		info.MetaPath,
		info.MetaHost,
		info.Name,
		info.StartTime,
		info.EndTime,
		info.TransactionStatus,
		info.TransactionType,
		info.Error,
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

func (c *Clickhouse) GetTransactionInfo(request *apiPb.GetTransactionsRequest) ([]*apiPb.TransactionInfo, int64, error) {
	timeFrom, timeTo, err := getTimeInt64(request.TimeRange)
	if err != nil {
		return nil, -1, err
	}

	var count int64
	count, err = c.countTransactions(request, timeFrom, timeTo)
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, request.GetPagination())

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM transaction_info WHERE (%s AND %s AND %s AND %s AND %s AND %s AND %s, AND %s) ORDER BY %s LIMIT %d OFFSET %d`,
		transactionInfoFields,
		applicationIdFilterString,
		applicationStartTimeFilterString,
		getTransactionsByString(transMetaHostStr, request.GetHost()),
		getTransactionsByString(transNameStr, request.GetName()),
		getTransactionsByString(transMetaPathStr, request.GetPath()),
		getTransactionsByString(transMetaMethodStr, request.GetMethod()),
		getTransactionTypeWhere(request.GetType()),
		getTransactionStatusWhere(request.GetStatus()),
		getTransactionOrder(request.GetSort())+getTransactionDirection(request.GetSort()), // todo
		limit,
		offset),
		request.ApplicationId,
		timeFrom,
		timeTo,
	)

	if err != nil {
		logger.Error(err.Error())
		return nil, -1, errorDataBase
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	var infos []*TransactionInfo
	for rows.Next() {
		inf := &TransactionInfo{}
		if err := rows.Scan(&inf.Model.ID, &inf.Model.CreatedAt, &inf.Model.UpdatedAt,
			&inf.ApplicationId, &inf.ParentId, &inf.MetaHost, &inf.MetaPath,
			&inf.MetaMethod, &inf.Name, &inf.StartTime, &inf.EndTime,
			&inf.TransactionStatus, &inf.TransactionType, &inf.Error); err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}

		infos = append(infos, inf)
	}

	if err := rows.Err(); err != nil {
		logger.Error(err.Error())
		return nil, -1, errorDataBase
	}

	return convertFromTransactions(infos), count, nil
}

func (c *Clickhouse) countTransactions(request *apiPb.GetTransactionsRequest, timeFrom int64, timeTo int64) (int64, error) {
	var count int64
	rows, err := c.Db.Query(fmt.Sprintf(`SELECT count(*) FROM transaction_info WHERE %s AND (%s) AND %s AND %s AND %s AND %s AND %s AND %s LIMIT 1`,
		applicationIdFilterString,
		applicationStartTimeFilterString,
		getTransactionsByString(transMetaHostStr, request.GetHost()),
		getTransactionsByString(transNameStr, request.GetName()),
		getTransactionsByString(transMetaPathStr, request.GetPath()),
		getTransactionsByString(transMetaMethodStr, request.GetMethod()),
		getTransactionTypeWhere(request.GetType()),
		getTransactionStatusWhere(request.GetStatus())),
		request.ApplicationId,
		timeFrom,
		timeTo)

	if err != nil {
		logger.Error(err.Error())
		return -1, errorDataBase
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	if ok := rows.Next(); !ok {
		return -1, errorDataBase
	}

	if err := rows.Scan(&count); err != nil {
		logger.Error(err.Error())
		return -1, errorDataBase
	}

	if err := rows.Err(); err != nil {
		logger.Error(err.Error())
		return -1, errorDataBase

	}
	return count, nil
}

func (c *Clickhouse) GetTransactionByID(request *apiPb.GetTransactionByIdRequest) (*apiPb.TransactionInfo, []*apiPb.TransactionInfo, error) {
	transaction, err := c.getTransaction(request.TransactionId)
	if err != nil || &transaction == nil {
		return nil, nil, err
	}

	children, err := c.GetTransactionChildren(transaction.TransactionId, "")
	if err != nil {
		return nil, nil, err
	}

	return convertFromTransaction(transaction), convertFromTransactions(children), nil
}

func (c *Clickhouse) getTransaction(id string) (*TransactionInfo, error) {
	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM transaction_info WHERE "transactionId" = ? LIMIT 1`, transactionInfoFields), id)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	inf := &TransactionInfo{}

	if ok := rows.Next(); !ok {
		return nil, err
	}

	if err := rows.Scan(&inf.Model.ID, &inf.Model.CreatedAt, &inf.Model.UpdatedAt,
		&inf.ApplicationId, &inf.ParentId, &inf.MetaHost, &inf.MetaPath,
		&inf.MetaMethod, &inf.Name, &inf.StartTime, &inf.EndTime,
		&inf.TransactionStatus, &inf.TransactionType, &inf.Error); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if err := rows.Err(); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return inf, nil
}

func (c *Clickhouse) GetTransactionChildren(transactionId, cyclicalLoopCheck string) ([]*TransactionInfo, error) {
	if strings.Contains(cyclicalLoopCheck, transactionId) {
		return nil, nil
	}

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM transaction_info WHERE "parentId" = ?`, transactionInfoFields), transactionId)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	var childTransactions []*TransactionInfo
	for rows.Next() {
		child := &TransactionInfo{}
		if err := rows.Scan(&child.Model.ID, &child.Model.CreatedAt, &child.Model.UpdatedAt,
			&child.ApplicationId, &child.ParentId, &child.MetaHost, &child.MetaPath,
			&child.MetaMethod, &child.Name, &child.StartTime, &child.EndTime,
			&child.TransactionStatus, &child.TransactionType, &child.Error); err != nil {
			logger.Error(err.Error())
			return nil, err
		}

		childTransactions = append(childTransactions, child)
	}

	for _, v := range childTransactions {
		grandChildren, err := c.GetTransactionChildren(v.TransactionId, cyclicalLoopCheck+" "+v.ParentId)
		if err != nil {
			return nil, errorDataBase
		}
		childTransactions = append(grandChildren, grandChildren...)
	}

	if err := rows.Err(); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return childTransactions, nil
}

// todo order
func (c *Clickhouse) GetTransactionGroup(request *apiPb.GetTransactionGroupRequest) (map[string]*apiPb.TransactionGroup, error) {
	timeFrom, timeTo, err := getTimeInt64(request.TimeRange)
	if err != nil {
		return nil, err
	}

	selection := fmt.Sprintf(
		`%s as "groupName", COUNT(%s) as "count", COUNT(CASE WHEN "%s"."transactionStatus" = '1' THEN 1 ELSE NULL END) as "successCount", AVG("%s"."endTime"-"%s"."startTime") as "latency", min("%s"."endTime"-"%s"."startTime") as "minTime", max("%s"."endTime"-"%s"."startTime") as "maxTime", min("%s"."endTime") as "lowTime"`,
		getTransactionsGroupBy(request.GetGroupType()),
		getTransactionsGroupBy(request.GetGroupType()),
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
		dbTransactionInfoCollection,
	)

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s transaction_info WHERE %s AND %s AND %s AND %s ORDER BY %s`,
		selection,
		applicationIdFilterString,
		applicationStartTimeFilterString,
		getTransactionTypeWhere(request.GetType()),
		getTransactionStatusWhere(request.GetStatus()),
		getTransactionsGroupBy(request.GetGroupType())),
		request.ApplicationId,
		timeFrom,
		timeTo,
	)

	if err != nil {
		logger.Error(err.Error())
		return nil, errorDataBase
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	var groupResults []*GroupResult
	for rows.Next() {
		res := &GroupResult{}
		if err := rows.Scan(&res.Name, &res.Count, &res.SuccessCount, &res.Latency, &res.MinTime, &res.MaxTime, &res.LowTime); err != nil {
			logger.Error(err.Error())
			return nil, err
		}

		groupResults = append(groupResults, res)
	}

	if err := rows.Err(); err != nil {
		logger.Error(err.Error())
		return nil, errorDataBase
	}

	return convertFromGroupResult(groupResults, timeTo), nil
}

func getTransactionOrder(request *apiPb.SortingTransactionList) string {
	if request == nil {
		return fmt.Sprintf(`"%s"."startTime"`, dbTransactionInfoCollection)
	}
	if res, ok := transOrderMap[request.GetSortBy()]; ok {
		return res
	}
	return fmt.Sprintf(`"%s"."startTime"`, dbTransactionInfoCollection)
}

func getTransactionDirection(request *apiPb.SortingTransactionList) string {
	if request == nil {
		return ` desc`
	}
	if res, ok := directionMap[request.GetDirection()]; ok {
		return res
	}
	return ` desc`
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
	return fmt.Sprintf(`"%s"."%s" = '%d'`, dbTransactionInfoCollection, transTransactionTypeStr, transType)
}

func getTransactionStatusWhere(transType apiPb.TransactionStatus) string {
	if transType == apiPb.TransactionStatus_TRANSACTION_CODE_UNSPECIFIED {
		return ""
	}
	return fmt.Sprintf(`"%s"."transactionStatus" = '%d'`, dbTransactionInfoCollection, transType)
}

func getTransactionsGroupBy(group apiPb.GroupTransaction) string {
	if val, ok := groupMap[group]; ok {
		return fmt.Sprintf(`"%s"."%s"`, dbTransactionInfoCollection, val)
	}
	return fmt.Sprintf(`"%s"."%s"`, dbTransactionInfoCollection, transTransactionTypeStr)
}
