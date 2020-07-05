package database

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

type mockMongo struct {
}

func (m mockMongo) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, nil
}

func (m mockMongo) FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error {
	return nil
}

func (m mockMongo) FindAll(ctx context.Context, predicate bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error {
	return nil
}

func (m mockMongo) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

var (
	db    = New(&mockMongo{})
	dbErr = New(&mockErrorMongo{})
)

type mockErrorMongo struct {
}

func (m mockErrorMongo) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorMongo) FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error {
	return errors.New("ERROR")
}

func (m mockErrorMongo) FindAll(ctx context.Context, predicate bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error {
	return errors.New("ERROR")
}

func (m mockErrorMongo) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, errors.New("ERROR")
}

func TestNew(t *testing.T) {
	t.Run("Should: create db", func(t *testing.T) {
		assert.NotNil(t, New(&mockMongo{}))
	})
}

func TestDatabase_SaveRule(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		assert.Nil(t, db.SaveRule(context.Background(), &Rule{}))
	})
}

func TestDatabase_FindRuleById(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		_, err := db.FindRuleById(context.Background(), primitive.NewObjectID())
		assert.Nil(t, err)
	})
}

func TestDatabase_FindRulesByOwnerId(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		_, err := db.FindRulesByOwnerId(context.Background(), 0, primitive.NewObjectID())
		assert.Nil(t, err)
	})
}

func TestDatabase_DeactivateRule(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		_, err := db.DeactivateRule(context.Background(), primitive.NewObjectID())
		assert.Nil(t, err)
	})
	t.Run("Should: return err", func(t *testing.T) {
		_, err := dbErr.DeactivateRule(context.Background(), primitive.NewObjectID())
		assert.Error(t, err)
	})
}

func TestDatabase_RemoveRule(t *testing.T) {
	t.Run("Should: return err", func(t *testing.T) {
		_, err := db.RemoveRule(context.Background(), primitive.NewObjectID())
		assert.Nil(t, err)
	})
	t.Run("Should: return nil", func(t *testing.T) {
		_, err := dbErr.RemoveRule(context.Background(), primitive.NewObjectID())
		assert.Error(t, err)
	})
}

func TestDatabase_ActivateRule(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		_, err := db.ActivateRule(context.Background(), primitive.NewObjectID())
		assert.Nil(t, err)
	})
	t.Run("Should: return err", func(t *testing.T) {
		_, err := dbErr.ActivateRule(context.Background(), primitive.NewObjectID())
		assert.Error(t, err)
	})
}
