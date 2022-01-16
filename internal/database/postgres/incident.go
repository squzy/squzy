package postgres

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

type Incident struct {
	gorm.Model
	IncidentId string             `gorm:"column:incidentId"`
	Status     int32              `gorm:"column:status"`
	RuleId     string             `gorm:"column:ruleId"`
	StartTime  int64              `gorm:"column:startTime"`
	EndTime    int64              `gorm:"column:endTime"`
	Histories  []*IncidentHistory `gorm:"column:history"`
}

type IncidentHistory struct {
	gorm.Model
	IncidentID uint  `gorm:"column:incidentId"`
	Status     int32 `gorm:"column:status"`
	Timestamp  int64 `gorm:"column:time"`
}

const (
	dbIncidentCollection        = "incidents"
	dbIncidentHistoryCollection = "incident_histories"
	descPrefix                  = " desc"

	NoIncident = "NO_INCIDENT"
)

var (
	incidentIdFilterString        = fmt.Sprintf(`"%s"."incidentId" = ?`, dbIncidentCollection)
	incidentRuleIdFilterString    = fmt.Sprintf(`"%s"."ruleId" = ?`, dbIncidentCollection)
	incidentStatusString          = fmt.Sprintf(`"%s"."status"`, dbIncidentCollection)
	incidentStartTimeFilterString = fmt.Sprintf(`"%s"."startTime" BETWEEN ? and ?`, dbIncidentCollection)
	statusString                  = "status"
	endTimeString                 = "endTime"

	incidentOrderMap = map[apiPb.SortIncidentList]string{
		apiPb.SortIncidentList_SORT_INCIDENT_LIST_UNSPECIFIED: fmt.Sprintf(`"%s"."startTime"`, dbIncidentCollection),
		apiPb.SortIncidentList_INCIDENT_LIST_BY_START_TIME:    fmt.Sprintf(`"%s"."startTime"`, dbIncidentCollection),
		apiPb.SortIncidentList_INCIDENT_LIST_BY_END_TIME:      fmt.Sprintf(`"%s"."endTime"`, dbIncidentCollection),
	}

	errorNoIncident = errors.New(NoIncident)
)

func (p *Postgres) InsertIncident(data *apiPb.Incident) error {
	incident := convertToIncident(data)
	if err := p.Db.Table(dbIncidentCollection).Create(incident).Error; err != nil {
		return errorDataBase
	}
	return nil
}

func (p *Postgres) UpdateIncidentStatus(id string, status apiPb.IncidentStatus) (*apiPb.Incident, error) {
	var incident Incident
	if err := p.Db.Table(dbIncidentCollection).
		Set("gorm:auto_preload", true).
		Where(incidentIdFilterString, id).First(&incident).Error; err != nil {
		return nil, errorDataBase
	}

	tNow := time.Now().UnixNano()

	if err := p.Db.Table(dbIncidentCollection).Where(incidentIdFilterString, id).
		Updates(
			map[string]interface{}{
				statusString:  int32(status),
				endTimeString: tNow,
			}).Error; err != nil {

		return nil, errorDataBase
	}

	history := &IncidentHistory{
		IncidentID: incident.ID,
		Status:     int32(status),
		Timestamp:  tNow,
	}
	if err := p.Db.Table(dbIncidentHistoryCollection).Create(history).Error; err != nil {

		return nil, errorDataBase
	}

	incident.Histories = append(incident.Histories, history)
	return convertFromIncident(&incident), nil
}

func (p *Postgres) GetIncidentById(id string) (*apiPb.Incident, error) {
	var incident Incident
	if err := p.Db.Table(dbIncidentCollection).
		Set("gorm:auto_preload", true).
		Where(incidentIdFilterString, id).First(&incident).Error; err != nil {
		return nil, errorDataBase
	}
	return convertFromIncident(&incident), nil
}

func (p *Postgres) GetActiveIncidentByRuleId(ruleId string) (*apiPb.Incident, error) {
	var incident Incident
	if err := p.Db.Table(dbIncidentCollection).
		Set("gorm:auto_preload", true).
		Where(incidentRuleIdFilterString, ruleId).
		Where(fmt.Sprintf(`%s OR %s OR %s`,
			getIncidentStatusString(apiPb.IncidentStatus_INCIDENT_STATUS_OPENED),
			getIncidentStatusString(apiPb.IncidentStatus_INCIDENT_STATUS_CAN_BE_CLOSED),
			getIncidentStatusString(apiPb.IncidentStatus_INCIDENT_STATUS_STUDIED))).
		First(&incident).Error; err != nil {

		return checkNoFoundError(err)
	}
	return convertFromIncident(&incident), nil
}

func (p *Postgres) GetIncidents(request *apiPb.GetIncidentsListRequest) ([]*apiPb.Incident, int64, error) {
	timeFrom, timeTo, err := getTimeInt64(request.GetTimeRange())
	if err != nil {
		return nil, -1, err
	}

	var count int64
	err = p.Db.Table(dbIncidentCollection).
		Where(incidentStartTimeFilterString, timeFrom, timeTo).
		Where(getIncidentStatusString(request.GetStatus())).
		Where(getIncidentRuleString(request.GetRuleId())).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, request.GetPagination())

	var incidents []*Incident
	err = p.Db.
		Table(dbIncidentCollection).
		Set("gorm:auto_preload", true).
		Where(incidentStartTimeFilterString, timeFrom, timeTo).
		Where(getIncidentStatusString(request.GetStatus())).
		Where(getIncidentRuleString(request.GetRuleId())).
		Order(getIncidentOrder(request.GetSort()) + getIncidentDirection(request.GetSort())).
		Offset(offset).
		Limit(limit).
		Find(&incidents).Error
	if err != nil {
		return nil, -1, errorDataBase
	}

	return convertFromIncidents(incidents), count, nil
}

func getIncidentStatusString(code apiPb.IncidentStatus) string {
	if code == apiPb.IncidentStatus_INCIDENT_STATUS_UNSPECIFIED {
		return ""
	}
	return fmt.Sprintf(`"%s"."status" = '%d'`, dbIncidentCollection, code)
}

func getIncidentRuleString(ruleId *wrappers.StringValue) string {
	if ruleId == nil {
		return ""
	}
	return fmt.Sprintf(`"%s"."ruleId" = '%s'`, dbIncidentCollection, ruleId.Value)
}

func getIncidentOrder(request *apiPb.SortingIncidentList) string {
	if request == nil {
		return fmt.Sprintf(`"%s"."startTime"`, dbIncidentCollection)
	}
	if res, ok := incidentOrderMap[request.GetSortBy()]; ok {
		return res
	}
	return fmt.Sprintf(`"%s"."startTime"`, dbIncidentCollection)
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

//Return empty incident
func checkNoFoundError(err error) (*apiPb.Incident, error) {
	if gorm.IsRecordNotFoundError(err) {
		return &apiPb.Incident{}, nil
	}
	return nil, errorDataBase
}
