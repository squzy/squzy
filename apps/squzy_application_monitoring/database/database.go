package database

import (
	"context"
	"github.com/squzy/mongo_helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Id   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
	Host string             `bson:"host,omitempty"`
}

type Database interface {
	FindOrCreate(ctx context.Context, name string, host string) (*Application, error)
}

type db struct {
	connector mongo_helper.Connector
}

func (d *db) FindOrCreate(ctx context.Context, name string, host string) (*Application, error) {
	app := &Application{}
	err := d.connector.FindOne(ctx, bson.M{
		"name": bson.M{
			"$eq": name,
		},
	}, app)

	if err == nil {
		return app, nil
	}

	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	app = &Application{
		Id: primitive.NewObjectID(),
		Host: host,
		Name: name,
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
