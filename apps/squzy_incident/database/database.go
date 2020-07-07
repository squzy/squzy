package database

import (
	"context"
	"github.com/squzy/mongo_helper"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Database interface {
	SaveRule(context.Context, *Rule) error
	FindRuleById(context.Context, primitive.ObjectID) (*Rule, error)
	FindRulesByOwnerId(ctx context.Context, ownerType apiPb.RuleOwnerType, ownerId primitive.ObjectID) ([]*Rule, error)
	RemoveRule(ctx context.Context, ruleId primitive.ObjectID) (*Rule, error)
	ActivateRule(ctx context.Context, ruleId primitive.ObjectID) (*Rule, error)
	DeactivateRule(ctx context.Context, ruleId primitive.ObjectID) (*Rule, error)
}

type Rule struct {
	Id        primitive.ObjectID  `bson:"_id"`
	Rule      string              `bson:"rule,omitempty"`
	Name      string              `bson:"name,omitempty"`
	AutoClose bool                `bson:"autoClose"`
	OwnerType apiPb.RuleOwnerType `bson:"ownerType"`
	OwnerId   primitive.ObjectID  `bson:"ownerId"`
	Status    apiPb.RuleStatus    `bson:"status"`
}

type database struct {
	mongo mongo_helper.Connector
}

var (
	activeStatus = []apiPb.RuleStatus{
		apiPb.RuleStatus_RULE_STATUS_ACTIVE, apiPb.RuleStatus_RULE_STATUS_INACTIVE,
	}
)

func New(mongo mongo_helper.Connector) Database {
	return &database{
		mongo: mongo,
	}
}

func (db *database) DeactivateRule(ctx context.Context, ruleId primitive.ObjectID) (*Rule, error) {
	return db.setStatus(ctx, ruleId, apiPb.RuleStatus_RULE_STATUS_INACTIVE)
}

func (db *database) setStatus(ctx context.Context, ruleId primitive.ObjectID, status apiPb.RuleStatus) (*Rule, error) {
	filter := bson.M{
		"_id": ruleId,
	}
	if status != apiPb.RuleStatus_RULE_STATUS_REMOVED {
		filter["status"] = bson.M{
			"$in": activeStatus,
		}
	}
	_, err := db.mongo.UpdateOne(ctx, filter, bson.M{
		"$set": bson.M{
			"status": status,
		},
	})
	if err != nil {
		return nil, err
	}
	return db.FindRuleById(ctx, ruleId)
}

func (db *database) ActivateRule(ctx context.Context, ruleId primitive.ObjectID) (*Rule, error) {
	return db.setStatus(ctx, ruleId, apiPb.RuleStatus_RULE_STATUS_ACTIVE)
}

func (db *database) SaveRule(ctx context.Context, rule *Rule) error {
	_, err := db.mongo.InsertOne(ctx, rule)
	return err
}

func (db *database) FindRuleById(ctx context.Context, id primitive.ObjectID) (*Rule, error) {
	rule := &Rule{}
	filter := bson.M{
		"_id": id,
	}
	err := db.mongo.FindOne(ctx, filter, rule)
	return rule, err
}

func (db *database) FindRulesByOwnerId(ctx context.Context, ownerType apiPb.RuleOwnerType, ownerId primitive.ObjectID) ([]*Rule, error) {
	var rules []*Rule
	filter := bson.M{
		"ownerType": ownerType,
		"ownerId":   ownerId,
	}
	err := db.mongo.FindAll(ctx, filter, &rules)
	return rules, err
}

func (db *database) RemoveRule(ctx context.Context, ruleId primitive.ObjectID) (*Rule, error) {
	return db.setStatus(ctx, ruleId, apiPb.RuleStatus_RULE_STATUS_REMOVED)
}
