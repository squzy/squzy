package database

import (
	"github.com/squzy/mongo_helper"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type Database interface {
	SaveRule(context.Context, *apiPb.Rule) error
	FindRuleById(context.Context, string) (*apiPb.Rule, error)
	FindRulesByOwnerId(context.Context, string) ([]*apiPb.Rule, error)
	RemoveRule(context.Context, string) error
}

type database struct {
	mongo mongo_helper.Connector
}

func New(mongo mongo_helper.Connector) Database {
	return &database{
		mongo: mongo,
	}
}

func (db *database) SaveRule(ctx context.Context, rule *apiPb.Rule) error {
	_, err := db.mongo.InsertOne(ctx, rule)
	return err
}

func (db *database) FindRuleById(ctx context.Context, id string) (*apiPb.Rule, error) {
	rule := &apiPb.Rule{}
	filter := bson.M{
		"_id":    id,
	}
	err := db.mongo.FindOne(ctx, filter, rule)
	return rule, err
}

func (db *database) FindRulesByOwnerId(ctx context.Context, id string) ([]*apiPb.Rule, error) {
	var rule []*apiPb.Rule
	filter := bson.M{
		"parent":    id,
	}
	err := db.mongo.FindAll(ctx, filter, rule)
	return rule, err
}

func (db *database) RemoveRule(context.Context, string) error {
	//TODO: implement removal in mongo
	return nil
}
