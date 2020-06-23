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
	db database.Database
	storage storage_client.Storage
}

func NewIncidentServer() apiPb.IncidentServerServer {
	return &server{}
}

func (s *server) CreateRule(ctx context.Context, request *apiPb.CreateRuleRequest) (*apiPb.Rule, error) {
	rule := &apiPb.Rule{
		Id:        primitive.NewObjectID().String(),
		Rule:      request.GetRule(),
		Name:      request.GetName(),
		AutoClose: request.GetAutoClose(),
		OwnerType: request.GetOwnerType(),
		OwnerId:   request.GetOwnerId(),
	}
	err := s.db.SaveRule(ctx, rule)
	return rule, err
}

func (s *server) GetRuleById(ctx context.Context, request *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	return s.db.FindRuleById(ctx, request.GetRuleId())
}

func (s *server) GetRulesByOwnerId(ctx context.Context, request *apiPb.GetRulesByOwnerIdRequest) (*apiPb.Rules, error) {
	rules, err := s.db.FindRulesByOwnerId(ctx, int32(request.GetOwnerType()), request.GetOwnerId())
	return &apiPb.Rules{
		Rules:                rules,
	}, err
}

func (s *server) RemoveRule(ctx context.Context, request *apiPb.RuleIdRequest) (*apiPb.Rule, error) {
	return nil, s.db.RemoveRule(ctx, request.RuleId)
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
		wasIncident := checkRule(rule.GetRule())
		incident, err := s.storage.GetIncident(ctx, rule.GetId())
		if err != nil {
			wasError = true
			continue
		}

		if isIncidentExist(incident) && isIncidentOpened(incident) && !wasIncident {
			if err := s.tryCloseIncident(ctx, rule, incident); err != nil {
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

func getOwnerTypeAndId(request *apiPb.StorageRecord) (int32, string, error) {
	if request.GetScheduler() != nil {
		return int32(apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT), request.GetScheduler().GetId(), nil
	}
	if request.GetAgent() != nil {
		return int32(apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT), request.GetAgent().GetAgentId(), nil
	}
	if request.GetTransaction() != nil {
		return int32(apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_APPLICATION), request.GetTransaction().GetId(), nil
	}
	return 0, "", errors.New("ERROR_NO_RECORD")
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

func (s *server) tryCloseIncident(ctx context.Context, rule *apiPb.Rule, incident *apiPb.Incident) error {
	if rule.AutoClose {
		_, err := s.storage.SetStatus(ctx, incident.GetId(), apiPb.IncidentStatus_INCIDENT_STATUS_CLOSED)
		return err
	}
	_, err := s.storage.SetStatus(ctx, incident.GetId(), apiPb.IncidentStatus_INCIDENT_STATUS_CAN_BE_CLOSED)
	return err
}