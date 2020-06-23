package server

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"squzy/apps/squzy_incident/database"
	"squzy/apps/squzy_incident/storage_client"
)

type server struct {
	db      database.Database
	storage storage_client.Storage
}

func NewIncidentServer() apiPb.IncidentServerServer {
	return &server{}
}

func dbRuleToProto(rule *database.Rule) *apiPb.Rule {
	return &apiPb.Rule{
		Id:        rule.Id.Hex(),
		Rule:      rule.Rule,
		Name:      rule.Name,
		AutoClose: rule.AutoClose,
		OwnerType: rule.OwnerType,
		OwnerId:   rule.OwnerId.Hex(),
		Status:    rule.Status,
	}
}

func (s *server) ActivateRule(ctx context.Context, request *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	ruleId, err := primitive.ObjectIDFromHex(request.RuleId)
	if err != nil {
		return nil, err
	}

	rule, err := s.db.ActivateRule(ctx, ruleId)

	if err != nil {
		return nil, err
	}

	return dbRuleToProto(rule), nil
}

func (s *server) DeactivateRule(ctx context.Context, request *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	ruleId, err := primitive.ObjectIDFromHex(request.RuleId)
	if err != nil {
		return nil, err
	}

	rule, err := s.db.DeactivateRule(ctx, ruleId)

	if err != nil {
		return nil, err
	}

	return dbRuleToProto(rule), nil
}

func (s *server) CreateRule(ctx context.Context, request *apiPb.CreateRuleRequest) (*apiPb.Rule, error) {
	ownerId, err := primitive.ObjectIDFromHex(request.OwnerId)
	if err != nil {
		return nil, err
	}
	rule := &database.Rule{
		Id:        primitive.NewObjectID(),
		Rule:      request.GetRule(),
		Name:      request.GetName(),
		AutoClose: request.GetAutoClose(),
		OwnerType: request.GetOwnerType(),
		OwnerId:   ownerId,
		Status:    apiPb.RuleStatus_RULE_STATUS_ACTIVE,
	}
	err = s.db.SaveRule(ctx, rule)
	if err != nil {
		return nil, err
	}
	return dbRuleToProto(rule), err
}

func (s *server) GetRuleById(ctx context.Context, request *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	ruleId, err := primitive.ObjectIDFromHex(request.RuleId)
	if err != nil {
		return nil, err
	}
	rule, err := s.db.FindRuleById(ctx, ruleId)
	if err != nil {
		return nil, err
	}
	return dbRuleToProto(rule), nil
}

func (s *server) GetRulesByOwnerId(ctx context.Context, request *apiPb.GetRulesByOwnerIdRequest) (*apiPb.Rules, error) {
	ownerId, err := primitive.ObjectIDFromHex(request.OwnerId)
	if err != nil {
		return nil, err
	}
	dbRules, err := s.db.FindRulesByOwnerId(ctx, request.GetOwnerType(), ownerId)
	rules := []*apiPb.Rule{}
	for _, rule := range dbRules {
		rules = append(rules, dbRuleToProto(rule))
	}
	return &apiPb.Rules{
		Rules: rules,
	}, err
}

func (s *server) RemoveRule(ctx context.Context, request *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	ruleId, err := primitive.ObjectIDFromHex(request.RuleId)
	if err != nil {
		return nil, err
	}
	rule, err := s.db.RemoveRule(ctx, ruleId)
	if err != nil {
		return nil, err
	}
	return dbRuleToProto(rule), nil
}

func (s *server) ProcessRecordFromStorage(ctx context.Context, request *apiPb.StorageRecord) (*empty.Empty, error) {
	ownerType, ownerId, err := getOwnerTypeAndId(request)
	if err != nil {
		return nil, err
	}

	rules, err := s.db.FindRulesByOwnerId(ctx, ownerType, ownerId)
	if err != nil {
		return nil, err
	}

	wasError := false
	for _, rule := range rules {
		wasIncident := checkRule(rule.Rule)
		incident, err := s.storage.GetIncident(ctx, rule.Id.Hex())
		if err != nil {
			wasError = true
			continue
		}

		if isIncidentExist(incident) && isIncidentOpened(incident) && !wasIncident {
			if err := s.tryCloseIncident(ctx, rule.AutoClose, incident); err != nil {
				wasError = true
			}
			continue
		}

		if !isIncidentExist(incident) && wasIncident {
			if err := s.storage.SaveIncident(ctx, incident); err != nil {
				wasError = true
			}
			continue
		}
	}

	if wasError {
		return &empty.Empty{}, errors.New("WAS_ERROR_WHILE_RULE_PROCESSING")
	}
	return &empty.Empty{}, nil
}

func (s *server) CloseIncident(ctx context.Context, request *apiPb.IncidentIdRequest) (*apiPb.Incident, error) {
	return s.storage.SetStatus(ctx, request.GetIncidentId(), apiPb.IncidentStatus_INCIDENT_STATUS_CLOSED)
}

func (s *server) StudyIncident(context.Context, *apiPb.IncidentIdRequest) (*apiPb.Incident, error) {
	panic("implement me")
}

func getOwnerTypeAndId(request *apiPb.StorageRecord) (apiPb.RuleOwnerType, primitive.ObjectID, error) {
	if request.GetScheduler() != nil {
		ownerId, err := primitive.ObjectIDFromHex(request.GetScheduler().Id)
		if err != nil {
			return 0, primitive.ObjectID{}, errors.New("ERROR_NO_RECORD")
		}
		return apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT, ownerId, nil
	}
	if request.GetAgent() != nil {
		ownerId, err := primitive.ObjectIDFromHex(request.GetAgent().AgentId)
		if err != nil {
			return 0, primitive.ObjectID{}, errors.New("ERROR_NO_RECORD")
		}
		return apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT, ownerId, nil
	}
	if request.GetTransaction() != nil {
		ownerId, err := primitive.ObjectIDFromHex(request.GetTransaction().Id)
		if err != nil {
			return 0, primitive.ObjectID{}, errors.New("ERROR_NO_RECORD")
		}
		return apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_APPLICATION, ownerId, nil
	}
	return 0, primitive.ObjectID{}, errors.New("ERROR_NO_RECORD")
}

func checkRule(rule string) bool {
	//TODO
	return false
}

func isIncidentExist(incident *apiPb.Incident) bool {
	return incident != nil
}

func isIncidentOpened(incident *apiPb.Incident) bool {
	if incident == nil {
		return false
	}
	return incident.GetStatus() == apiPb.IncidentStatus_INCIDENT_STATUS_OPENED
}

func (s *server) tryCloseIncident(ctx context.Context, autoClose bool, incident *apiPb.Incident) error {
	if autoClose {
		_, err := s.storage.SetStatus(ctx, incident.GetId(), apiPb.IncidentStatus_INCIDENT_STATUS_CLOSED)
		return err
	}
	_, err := s.storage.SetStatus(ctx, incident.GetId(), apiPb.IncidentStatus_INCIDENT_STATUS_CAN_BE_CLOSED)
	return err
}
