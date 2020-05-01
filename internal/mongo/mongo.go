package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type connector struct {
	collection *mongo.Collection
}

type Connector interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error
	FindAll(ctx context.Context, predicate bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

func (c *connector) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.collection.UpdateOne(ctx, filter, update, opts...)
}

func (c *connector) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.collection.InsertOne(ctx, document, opts...)
}

func (c *connector) FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error {
	raw, err := c.collection.FindOne(ctx, filter, opts...).DecodeBytes()
	if err != nil {
		return err
	}
	return bson.UnmarshalWithContext(bsoncodec.DecodeContext{
		Registry: bson.DefaultRegistry,
		Truncate: true,
	}, raw, structToDeserialize)
}

func (c *connector) FindAll(ctx context.Context, filter bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error {
	cursor, err := c.collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	return cursor.All(ctx, structToDeserialize)
}

func New(collection *mongo.Collection) Connector {
	return &connector{
		collection: collection,
	}
}
