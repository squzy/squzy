package clickhouse

import (
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	uuid "github.com/google/uuid"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

type Incident struct {
	Model      Model
	IncidentId string
	Status     int32
	RuleId     string
	StartTime  int64
	EndTime    int64
	Histories  []*IncidentHistory
}

type IncidentHistory struct {
	Model      Model
	IncidentID string
	Status     int32
	Timestamp  int64
}

const (
	descPrefix = " DESC"
	noSep      = ""
	orSep      = " OR"
	andSep     = " AND"
)

var (
	incidentFields          = "id, created_at, updated_at, incident_id, status, rule_id, start_time, end_time"
	incidentHistoriesFields = "id, created_at, incident_id, status, timestamp"
	incidentIdString        = fmt.Sprintf(`"incident_id" = ?`)
	incidentRuleIdString    = fmt.Sprintf(`"rule_id" = ?`)
	incidentStatusString    = fmt.Sprintf(`"status" = ?`)
	startTimeFilterString   = `start_time >= ? AND start_time <= ?`

	incidentOrderMap = map[apiPb.SortIncidentList]string{
		apiPb.SortIncidentList_SORT_INCIDENT_LIST_UNSPECIFIED: "start_time",
		apiPb.SortIncidentList_INCIDENT_LIST_BY_START_TIME:    "start_time",
		apiPb.SortIncidentList_INCIDENT_LIST_BY_END_TIME:      "end_time",
	}
)

func (c *Clickhouse) InsertIncident(data *apiPb.Incident) error {
	now := time.Now()

	incident := convertToIncident(data, now)

	err := c.insertIncident(now, incident)
	if err != nil {
		logger.Error(err.Error())
		return errorDataBase
	}

	for _, history := range incident.Histories {
		err = c.insertIncidentHistory(incident.IncidentId, now, history.Status, history.Timestamp)
		if err != nil {
			logger.Error(err.Error())
			return errorDataBase
		}
	}

	return nil
}

func (c *Clickhouse) insertIncident(now time.Time, incident *Incident) error {
	tx, err := c.Db.Begin()
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`INSERT INTO incidents (%s) VALUES ($0, $1, $2, $3, $4, $5, $6, $7)`, incidentFields)
	_, err = tx.Exec(q,
		clickhouse.UUID(uuid.New().String()),
		now,
		now,
		incident.IncidentId,
		incident.Status,
		incident.RuleId,
		incident.StartTime,
		incident.EndTime,
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

func (c *Clickhouse) insertIncidentHistory(incidentId string, now time.Time, status int32, timestamp int64) error {
	tx, err := c.Db.Begin()
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`INSERT INTO incidents_history (%s) VALUES ($0, $1, $2, $3, $4)`, incidentHistoriesFields)
	_, err = tx.Exec(q,
		clickhouse.UUID(uuid.New().String()),
		now,
		incidentId,
		status,
		timestamp,
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

func (c *Clickhouse) UpdateIncidentStatus(id string, status apiPb.IncidentStatus) (*apiPb.Incident, error) {
	now := time.Now()
	uNow := now.UnixNano()

	oldInc, err := c.getIncidentById(id)
	if err != nil {
		logger.Error(err.Error())
		return nil, errorDataBase
	}

	err = c.updateIncident(now, status, oldInc)
	if err != nil {
		logger.Error(err.Error())
		return nil, errorDataBase
	}

	err = c.insertIncidentHistory(id, now, int32(status), uNow)
	if err != nil {
		logger.Error(err.Error())
		return nil, errorDataBase
	}

	newInc, err := c.GetIncidentById(id)
	if err != nil {
		logger.Error(err.Error())
		return nil, errorDataBase
	}
	return newInc, nil
}

func (c *Clickhouse) updateIncident(now time.Time, status apiPb.IncidentStatus, incident *Incident) error {
	tx, err := c.Db.Begin()
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`INSERT INTO incidents (%s) VALUES ($0, $1, $2, $3, $4, $5, $6, $7)`, incidentFields)
	_, err = tx.Exec(q,
		clickhouse.UUID(uuid.New().String()),
		incident.Model.CreatedAt,
		now,
		incident.IncidentId,
		int32(status),
		incident.RuleId,
		incident.StartTime,
		incident.EndTime,
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

func (c *Clickhouse) GetIncidentById(id string) (*apiPb.Incident, error) {
	inc, err := c.getIncidentById(id)
	if err != nil {
		logger.Error(err.Error())
		return nil, errorDataBase
	}

	return convertFromIncident(inc), nil
}

func (c *Clickhouse) getIncidentById(id string) (*Incident, error) {
	inc, err := c.getIncident(id)
	if err != nil {
		return nil, err
	}

	histories, err := c.getIncidentHistories(id)
	if err != nil {
		return nil, err
	}

	if inc != nil {
		inc.Histories = histories
	}

	return inc, nil
}

func (c *Clickhouse) getIncident(id string) (*Incident, error) {
	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM incidents FINAL WHERE %s LIMIT 1`, incidentFields, incidentIdString), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if ok := rows.Next(); !ok {
		return &Incident{}, nil
	}

	inc := &Incident{}
	if err := rows.Scan(&inc.Model.ID, &inc.Model.CreatedAt, &inc.Model.UpdatedAt,
		&inc.IncidentId, &inc.Status, &inc.RuleId, &inc.StartTime, &inc.EndTime); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return inc, nil
}

func (c *Clickhouse) getIncidentHistories(id string) ([]*IncidentHistory, error) {
	var incs []*IncidentHistory

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM incidents_history WHERE %s`,
		incidentHistoriesFields, incidentIdString), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		inc := &IncidentHistory{}
		if err := rows.Scan(&inc.Model.ID, &inc.Model.CreatedAt, &inc.IncidentID, &inc.Status, &inc.Timestamp); err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		incs = append(incs, inc)
	}

	return incs, nil
}

func (c *Clickhouse) GetActiveIncidentByRuleId(ruleId string) (*apiPb.Incident, error) {
	inc, err := c.getActiveIncident(ruleId)
	if err != nil {
		logger.Error(err.Error())
		return nil, errorDataBase
	}

	inc.Histories, err = c.getIncidentHistories(inc.IncidentId)
	if err != nil {
		logger.Error(err.Error())
		return nil, errorDataBase
	}

	return convertFromIncident(inc), nil
}

func (c *Clickhouse) getActiveIncident(ruleId string) (*Incident, error) {
	inc := &Incident{}

	rIDQuery, rIDValue := getIncidentRuleString(ruleId, andSep)
	queryParams := createParamsWithVal(rIDValue)
	s1Query, s1Value := getIncidentStatusString(apiPb.IncidentStatus_INCIDENT_STATUS_OPENED, orSep)
	queryParams = addParamsWithIntVal(s1Value, queryParams)
	s2Query, s2Value := getIncidentStatusString(apiPb.IncidentStatus_INCIDENT_STATUS_CAN_BE_CLOSED, orSep)
	queryParams = addParamsWithIntVal(s2Value, queryParams)
	s3Query, s3Value := getIncidentStatusString(apiPb.IncidentStatus_INCIDENT_STATUS_STUDIED, noSep)
	queryParams = addParamsWithIntVal(s3Value, queryParams)

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM incidents WHERE %s (%s %s %s) LIMIT 1`,
		incidentFields,
		rIDQuery,
		s1Query,
		s2Query,
		s3Query,
	), queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if ok := rows.Next(); !ok {
		return inc, nil
	}

	if err := rows.Scan(&inc.Model.ID, &inc.Model.CreatedAt, &inc.Model.UpdatedAt,
		&inc.IncidentId, &inc.Status, &inc.RuleId, &inc.StartTime, &inc.EndTime); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return inc, nil
}

func (c *Clickhouse) GetIncidents(request *apiPb.GetIncidentsListRequest) ([]*apiPb.Incident, int64, error) {
	timeFrom, timeTo, err := getTimeInt64(request.GetTimeRange())
	if err != nil {
		return nil, -1, err
	}

	count, err := c.countIncidents(request, timeFrom, timeTo)
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, request.GetPagination())

	sQuery, sValue := getIncidentStatusString(request.Status, andSep)
	rIDQuery, rIDValue := getIncidentRuleString(unwrapRuleString(request.RuleId), andSep)
	queryParams := createParamsWithVal(rIDValue)
	queryParams = addParamsWithIntVal(sValue, queryParams)

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM incidents WHERE (%s %s %s) ORDER BY %s LIMIT %d OFFSET %d`,
		incidentFields,
		rIDQuery,
		sQuery,
		startTimeFilterString,
		getIncidentOrder(request.GetSort())+getIncidentDirection(request.GetSort()),
		limit,
		offset),
		append(queryParams, timeFrom, timeTo)...,
	)

	if err != nil {
		logger.Error(err.Error())
		return nil, -1, errorDataBase
	}
	defer rows.Close()

	var incs []*Incident
	for rows.Next() {
		inc := &Incident{}
		if err := rows.Scan(&inc.Model.ID, &inc.Model.CreatedAt, &inc.Model.UpdatedAt,
			&inc.IncidentId, &inc.Status, &inc.RuleId, &inc.StartTime, &inc.EndTime); err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}

		histories, err := c.getIncidentHistories(inc.IncidentId)
		if err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}

		inc.Histories = histories
		incs = append(incs, inc)
	}

	return convertFromIncidents(incs), count, nil
}

func createParamsWithVal(val string) []interface{} {
	queryParams := make([]interface{}, 0)
	if val != "" {
		queryParams = append(queryParams, val)
	}
	return queryParams
}

func addParamsWithIntVal(val int32, queryParams []interface{}) []interface{} {
	if val != 0 {
		queryParams = append(queryParams, val)
	}
	return queryParams
}

func (c *Clickhouse) countIncidents(request *apiPb.GetIncidentsListRequest, timeFrom int64, timeTo int64) (int64, error) {
	var count int64

	rIDQuery, rIDValue := getIncidentRuleString(unwrapRuleString(request.RuleId), andSep)
	queryParams := createParamsWithVal(rIDValue)
	sQuery, sValue := getIncidentStatusString(request.Status, andSep)
	queryParams = addParamsWithIntVal(sValue, queryParams)

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT count(*) FROM incidents WHERE %s %s %s`,
		rIDQuery,
		sQuery,
		startTimeFilterString),
		append(queryParams, timeFrom, timeTo)...)

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

func getIncidentStatusString(code apiPb.IncidentStatus, separator string) (string, int32) {
	if code == apiPb.IncidentStatus_INCIDENT_STATUS_UNSPECIFIED {
		return "", 0
	}
	return fmt.Sprintf(`%s %s`, incidentStatusString, separator), int32(code)
}

func unwrapRuleString(ruleId *wrapperspb.StringValue) string {
	if ruleId == nil {
		return ""
	}
	return ruleId.Value
}

func getIncidentRuleString(ruleId string, separator string) (string, string) {
	if ruleId == "" {
		return "", ""
	}
	return fmt.Sprintf(`%s %s`, incidentRuleIdString, separator), ruleId
}

func getIncidentOrder(request *apiPb.SortingIncidentList) string {
	if request == nil {
		return incidentOrderMap[apiPb.SortIncidentList_SORT_INCIDENT_LIST_UNSPECIFIED]
	}
	if res, ok := incidentOrderMap[request.GetSortBy()]; ok {
		return res
	}
	return incidentOrderMap[apiPb.SortIncidentList_SORT_INCIDENT_LIST_UNSPECIFIED]
}

func getIncidentDirection(request *apiPb.SortingIncidentList) string {
	if request == nil {
		return descPrefix
	}
	if res, ok := directionMap[request.GetDirection()]; ok {
		return res
	}
	return descPrefix
}
