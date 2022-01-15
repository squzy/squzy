package database

import (
	"context"
	"errors"
	"github.com/squzy/mongo_helper"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Id      primitive.ObjectID      `bson:"_id"`
	Name    string                  `bson:"name"`
	Host    string                  `bson:"host,omitempty"`
	Status  apiPb.ApplicationStatus `bson:"status"`
	AgentId string                  `bson:"agentId,omitempty"`
}

type Database interface {
	FindOrCreate(ctx context.Context, name string, host string, agentId string) (*Application, error)
	FindApplicationByName(ctx context.Context, name string) (*Application, error)
	FindApplicationById(ctx context.Context, id primitive.ObjectID) (*Application, error)
	FindAllApplication(ctx context.Context) ([]*Application, error)
	FindApplicationByAgentId(ctx context.Context, agentId primitive.ObjectID) ([]*Application, error)
	SetStatus(ctx context.Context, id primitive.ObjectID, status apiPb.ApplicationStatus) error
}

type db struct {
	connector mongo_helper.Connector
}

var (
	errArchived = errors.New("application already archived")
)

func (d *db) FindApplicationByAgentId(ctx context.Context, agentId primitive.ObjectID) ([]*Application, error) {
	return d.findListApplication(ctx, bson.M{
		"agentId": agentId,
	})
}

func (d *db) SetStatus(ctx context.Context, id primitive.ObjectID, status apiPb.ApplicationStatus) error {
	app, err := d.findApplication(ctx, bson.M{
		"_id": bson.M{
			"$eq": id,
		},
	})

	if err != nil {
		return err
	}

	if app.Status == status {
		return nil
	}

	if app.Status == apiPb.ApplicationStatus_APPLICATION_STATUS_ARCHIVED {
		return errArchived
	}

	_, err = d.connector.UpdateOne(ctx, bson.M{
		"_id":    id,
		"status": app.Status,
	}, bson.M{
		"$set": bson.M{
			"status": status,
		},
	})

	return err
}

func (d *db) FindApplicationById(ctx context.Context, id primitive.ObjectID) (*Application, error) {
	app, err := d.findApplication(ctx, bson.M{
		"_id": bson.M{
			"$eq": id,
		},
	})
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (d *db) findApplication(ctx context.Context, filter bson.M) (*Application, error) {
	app := &Application{}
	err := d.connector.FindOne(ctx, filter, app)

	if err != nil {
		return nil, err
	}
	return app, nil
}

func (d *db) findListApplication(ctx context.Context, filter bson.M) ([]*Application, error) {
	list := []*Application{}
	err := d.connector.FindAll(ctx, filter, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (d *db) FindAllApplication(ctx context.Context) ([]*Application, error) {
	return d.findListApplication(ctx, bson.M{})
}

func (d *db) FindApplicationByName(ctx context.Context, name string) (*Application, error) {
	app, err := d.findApplication(ctx, bson.M{
		"name": bson.M{
			"$eq": name,
		},
	})
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (d *db) FindOrCreate(ctx context.Context, name string, host string, agentID string) (*Application, error) {
	app, err := d.FindApplicationByName(ctx, name)

	if err == nil {
		return app, nil
	}

	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	app = &Application{
		Id:      primitive.NewObjectID(),
		Host:    host,
		Name:    name,
		Status:  apiPb.ApplicationStatus_APPLICATION_STATUS_ENABLED,
		AgentId: agentID,
	}

	_, err = d.connector.InsertOne(ctx, app)

	if err != nil {
		return nil, err
	}

	return app, nil
}

func New(connector mongo_helper.Connector) Database {
	return &db{
		connector: connector,
	}
}
