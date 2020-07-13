package database

import (
	"context"
	"github.com/squzy/mongo_helper"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	Id primitive.ObjectID `bson:"_id"`
	OwnerId primitive.ObjectID `bson:"ownerId"`
	Type apiPb.NotificationMethodType `bson:"type"`
	NotificationMethodId primitive.ObjectID `bson:"notificationMethodId"`
}

type NotificationList interface {
	Add(ctx context.Context, notification *Notification) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetList(ctx context.Context, OwnerId primitive.ObjectID, Type apiPb.NotificationMethodType) ([]*Notification, error)
}

type notificationList struct {
	mongo mongo_helper.Connector
}

func (n *notificationList) Add(ctx context.Context, notification *Notification) error {
	_, err := n.mongo.InsertOne(ctx, notification)
	return err
}

func (n *notificationList) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := n.mongo.Delete(ctx, bson.M{
		"_id": id,
	})
	return err
}

func (n *notificationList) GetList(ctx context.Context, ownerId primitive.ObjectID, methodType apiPb.NotificationMethodType) ([]*Notification, error) {
	list := []*Notification{}
	err := n.mongo.FindAll(ctx, bson.M{
		"ownerId": ownerId,
		"type": methodType,
	}, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func NewList(mongo mongo_helper.Connector ) NotificationList {
	return &notificationList{
		mongo: mongo,
	}
}

