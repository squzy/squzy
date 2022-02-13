package clickhouse

//
import (
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/golang/protobuf/ptypes/wrappers"
	uuid "github.com/google/uuid"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
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
	transactionInfoFields   = "id, created_at, updated_at, transaction_id, application_id, parent_id, meta_host, meta_path, meta_method, name, start_time, end_time, transaction_status, transaction_type, error"
	transNameStr            = "name"
	transMetaHostStr        = "meta_host"
	transMetaMethodStr      = "meta_method"
	transMetaPathStr        = "meta_path"
	transTransactionTypeStr = "transaction_type"
)

var (
	applicationIdFilterString        = fmt.Sprintf(`"application_id" = ?`)
	applicationStartTimeFilterString = fmt.Sprintf(`"start_time" BETWEEN ? and ?`)

	transOrderMap = map[apiPb.SortTransactionList]string{
		apiPb.SortTransactionList_SORT_TRANSACTION_LIST_UNSPECIFIED: fmt.Sprintf(`"%s"."start_time"`, dbTransactionInfoCollection),
		apiPb.SortTransactionList_DURATION:                          fmt.Sprintf(`"%s"."end_time" - "%s"."start_time"`, dbTransactionInfoCollection, dbTransactionInfoCollection),
		apiPb.SortTransactionList_BY_TRANSACTION_START_TIME:         fmt.Sprintf(`"%s"."start_time"`, dbTransactionInfoCollection),
		apiPb.SortTransactionList_BY_TRANSACTION_END_TIME:           fmt.Sprintf(`"%s"."end_time"`, dbTransactionInfoCollection),
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
		clickhouse.UUID(uuid.New().String()),
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

	q := fmt.Sprintf(`SELECT %s FROM transaction_info WHERE ( %s AND %s AND %s AND %s AND %s AND %s AND %s AND %s ) ORDER BY %s LIMIT %d OFFSET %d`,
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
		offset)
	q = strings.ReplaceAll(q, "AND  ", "")
	rows, err := c.Db.Query(q,
		request.ApplicationId,
		timeFrom,
		timeTo,
	)

	if err != nil {
		logger.Error(err.Error())
		return nil, -1, errorDataBase
	}
	defer rows.Close()

	var infos []*TransactionInfo
	for rows.Next() {
		inf := &TransactionInfo{}
		if err := rows.Scan(&inf.Model.ID, &inf.Model.CreatedAt, &inf.Model.UpdatedAt,
			&inf.TransactionId, &inf.ApplicationId, &inf.ParentId, &inf.MetaHost, &inf.MetaPath,
			&inf.MetaMethod, &inf.Name, &inf.StartTime, &inf.EndTime,
			&inf.TransactionStatus, &inf.TransactionType, &inf.Error); err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}

		infos = append(infos, inf)
	}

	return convertFromTransactions(infos), count, nil
}

func (c *Clickhouse) countTransactions(request *apiPb.GetTransactionsRequest, timeFrom int64, timeTo int64) (int64, error) {
	var count int64
	q := fmt.Sprintf(`SELECT count(*) FROM transaction_info WHERE %s AND (%s) AND %s AND %s AND %s AND %s AND %s AND %s LIMIT 1`,
		applicationIdFilterString,
		applicationStartTimeFilterString,
		getTransactionsByString(transMetaHostStr, request.GetHost()),
		getTransactionsByString(transNameStr, request.GetName()),
		getTransactionsByString(transMetaPathStr, request.GetPath()),
		getTransactionsByString(transMetaMethodStr, request.GetMethod()),
		getTransactionTypeWhere(request.GetType()),
		getTransactionStatusWhere(request.GetStatus()))
	q = strings.ReplaceAll(q, "AND  ", "")
	rows, err := c.Db.Query(q,
		request.ApplicationId,
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
	defer rows.Close()

	inf := &TransactionInfo{}

	if ok := rows.Next(); !ok {
		return nil, nil
	}

	if err := rows.Scan(&inf.Model.ID, &inf.Model.CreatedAt, &inf.Model.UpdatedAt,
		&inf.ApplicationId, &inf.ParentId, &inf.MetaHost, &inf.MetaPath,
		&inf.MetaMethod, &inf.Name, &inf.StartTime, &inf.EndTime,
		&inf.TransactionStatus, &inf.TransactionType, &inf.Error); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return inf, nil
}

func (c *Clickhouse) GetTransactionChildren(transactionId, cyclicalLoopCheck string) ([]*TransactionInfo, error) {
	if strings.Contains(cyclicalLoopCheck, transactionId) {
		return nil, nil
	}

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM transaction_info WHERE "parent_id" = ?`, transactionInfoFields), transactionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var childTransactions []*TransactionInfo
	for rows.Next() {
		child := &TransactionInfo{}
		if err := rows.Scan(&child.Model.ID, &child.Model.CreatedAt, &child.Model.UpdatedAt,
			&child.TransactionId, &child.ApplicationId,
			&child.ParentId, &child.MetaHost, &child.MetaPath,
			&child.MetaMethod, &child.Name, &child.StartTime, &child.EndTime,
			&child.TransactionStatus, &child.TransactionType, &child.Error); err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		childTransactions = append(childTransactions, child)
	}

	var allChildTransactions []*TransactionInfo
	for _, v := range childTransactions {
		grandChildren, err := c.GetTransactionChildren(v.TransactionId, cyclicalLoopCheck+" "+v.ParentId)
		if err != nil {
			return nil, errorDataBase
		}
		allChildTransactions = append(allChildTransactions, grandChildren...)
	}

	allChildTransactions = append(allChildTransactions, childTransactions...)
	return allChildTransactions, nil
}

// todo order
func (c *Clickhouse) GetTransactionGroup(request *apiPb.GetTransactionGroupRequest) (map[string]*apiPb.TransactionGroup, error) {
	timeFrom, timeTo, err := getTimeInt64(request.TimeRange)
	if err != nil {
		return nil, err
	}

	selection := fmt.Sprintf(
		`%s as "groupName", COUNT(%s) as "count", COUNT(CASE WHEN "%s"."transaction_status" = '1' THEN 1 ELSE NULL END) as "successCount", AVG("%s"."end_time"-"%s"."start_time") as "latency", min("%s"."end_time"-"%s"."start_time") as "minTime", max("%s"."end_time"-"%s"."start_time") as "maxTime", min("%s"."end_time") as "lowTime"`,
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

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM transaction_info WHERE %s AND %s AND %s AND %s GROUP BY %s`,
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
	defer rows.Close()

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
		return fmt.Sprintf(`"%s"."start_time"`, dbTransactionInfoCollection)
	}
	if res, ok := transOrderMap[request.GetSortBy()]; ok {
		return res
	}
	return fmt.Sprintf(`"%s"."start_time"`, dbTransactionInfoCollection)
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
	return fmt.Sprintf(`"%s"."transaction_status" = '%d'`, dbTransactionInfoCollection, transType)
}

func getTransactionsGroupBy(group apiPb.GroupTransaction) string {
	if val, ok := groupMap[group]; ok {
		return fmt.Sprintf(`"%s"."%s"`, dbTransactionInfoCollection, val)
	}
	return fmt.Sprintf(`"%s"."%s"`, dbTransactionInfoCollection, transTransactionTypeStr)
}
