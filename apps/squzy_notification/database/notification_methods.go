package database

import (
	"context"
	"github.com/squzy/mongo_helper"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationMethodDb interface {
	Create(ctx context.Context, nm *NotificationMethod) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	Activate(ctx context.Context, id primitive.ObjectID) error
	Deactivate(ctx context.Context, id primitive.ObjectID) error
	Get(ctx context.Context, id primitive.ObjectID) (*NotificationMethod, error)
	GetAll(ctx context.Context) ([]*NotificationMethod, error)
}

type NotificationMethod struct {
	Id      primitive.ObjectID             `bson:"_id"`
	Name    string                         `bson:"name"`
	Status  apiPb.NotificationMethodStatus `bson:"status"`
	Type    apiPb.NotificationMethodType   `bson:"type"`
	Slack   *SlackConfig                   `bson:"slackConfig,omitempty"`
	WebHook *WebHookConfig                 `bson:"webhookConfig,omitempty"`
}

type SlackConfig struct {
	Url string `bson:"string"`
}

type WebHookConfig struct {
	Url string `bson:"string"`
}

type notificationMethodDb struct {
	mongo mongo_helper.Connector
}

var (
	activeStatus = []apiPb.NotificationMethodStatus{
		apiPb.NotificationMethodStatus_NOTIFICATION_STATUS_ACTIVE,
		apiPb.NotificationMethodStatus_NOTIFICATION_STATUS_INACTIVE,
	}
)

func (n *notificationMethodDb) GetAll(ctx context.Context) ([]*NotificationMethod, error) {
	items := []*NotificationMethod{}
	err := n.mongo.FindAll(ctx, bson.M{}, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (n *notificationMethodDb) Create(ctx context.Context, nm *NotificationMethod) error {
	_, err := n.mongo.InsertOne(ctx, nm)
	return err
}

func (n *notificationMethodDb) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := n.mongo.UpdateOne(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"status": apiPb.NotificationMethodStatus_NOTIFICATION_STATUS_REMOVED,
		},
	})
	return err
}

func (n *notificationMethodDb) Activate(ctx context.Context, id primitive.ObjectID) error {
	_, err := n.mongo.UpdateOne(ctx, bson.M{
		"_id": id,
		"status": bson.M{
			"$in": activeStatus,
		},
	}, bson.M{
		"$set": bson.M{
			"status": apiPb.NotificationMethodStatus_NOTIFICATION_STATUS_ACTIVE,
		},
	})
	return err
}

func (n *notificationMethodDb) Deactivate(ctx context.Context, id primitive.ObjectID) error {
	_, err := n.mongo.UpdateOne(ctx, bson.M{
		"_id": id,
		"status": bson.M{
			"$in": activeStatus,
		},
	}, bson.M{
		"$set": bson.M{
			"status": apiPb.NotificationMethodStatus_NOTIFICATION_STATUS_INACTIVE,
		},
	})
	return err
}

func (n *notificationMethodDb) Get(ctx context.Context, id primitive.ObjectID) (*NotificationMethod, error) {
	item := &NotificationMethod{}
	err := n.mongo.FindOne(ctx, bson.M{
		"_id": id,
	}, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func NewMethods(mongo mongo_helper.Connector) NotificationMethodDb {
	return &notificationMethodDb{
		mongo: mongo,
	}
}
