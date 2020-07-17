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

type mockSuccess struct {

}

func (m mockSuccess) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, nil
}

func (m mockSuccess) FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error {
	return nil
}

func (m mockSuccess) Delete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return nil, nil
}

func (m mockSuccess) FindAll(ctx context.Context, predicate bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error {
	return nil
}

func (m mockSuccess) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

type mockError struct {

}

func (m mockError) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, errors.New("")
}

func (m mockError) FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error {
	return errors.New("")
}

func (m mockError) Delete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return nil, errors.New("")
}

func (m mockError) FindAll(ctx context.Context, predicate bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error {
	return errors.New("")
}

func (m mockError) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, errors.New("")
}

func TestNewList(t *testing.T) {
	t.Run("Should: not be nil", func(t *testing.T) {
		s := NewList(nil)
		assert.NotNil(t, s)
	})
}

func TestNewMethods(t *testing.T) {
	t.Run("Should: not be nil", func(t *testing.T) {
		s := NewMethods(nil)
		assert.NotNil(t, s)
	})
}

func TestNotificationList_Add(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		s := NewList(&mockSuccess{})
		assert.Nil(t, s.Add(context.Background(), &Notification{}))
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := NewList(&mockError{})
		assert.NotNil(t, s.Add(context.Background(), &Notification{}))
	})
}

func TestNotificationList_Delete(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		s := NewList(&mockSuccess{})
		assert.Nil(t, s.Delete(context.Background(), primitive.NewObjectID()))
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := NewList(&mockError{})
		assert.NotNil(t, s.Delete(context.Background(), primitive.NewObjectID()))
	})
}

func TestNotificationList_GetList(t *testing.T) {
	t.Run("Should: return no error", func(t *testing.T) {
		s := NewList(&mockSuccess{})
		_, err := s.GetList(context.Background(),  primitive.NewObjectID(), 1)
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := NewList(&mockError{})
		_, err := s.GetList(context.Background(),  primitive.NewObjectID(), 1)
		assert.NotNil(t, err)
	})
}

func TestNotificationMethodDb_Activate(t *testing.T) {
	t.Run("Should: return no error", func(t *testing.T) {
		s := NewMethods(&mockSuccess{})
		assert.Nil(t, s.Activate(context.Background(), primitive.NewObjectID()))
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := NewMethods(&mockError{})
		assert.NotNil(t, s.Activate(context.Background(), primitive.NewObjectID()))
	})
}

func TestNotificationMethodDb_Deactivate(t *testing.T) {
	t.Run("Should: return no error", func(t *testing.T) {
		s := NewMethods(&mockSuccess{})
		assert.Nil(t, s.Deactivate(context.Background(), primitive.NewObjectID()))
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := NewMethods(&mockError{})
		assert.NotNil(t, s.Deactivate(context.Background(), primitive.NewObjectID()))
	})
}

func TestNotificationMethodDb_Create(t *testing.T) {
	t.Run("Should: return no error", func(t *testing.T) {
		s := NewMethods(&mockSuccess{})
		assert.Nil(t, s.Create(context.Background(), &NotificationMethod{}))
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := NewMethods(&mockError{})
		assert.NotNil(t, s.Create(context.Background(), &NotificationMethod{}))
	})
}

func TestNotificationMethodDb_Delete(t *testing.T) {
	t.Run("Should: return no error", func(t *testing.T) {
		s := NewMethods(&mockSuccess{})
		assert.Nil(t, s.Delete(context.Background(), primitive.NewObjectID()))
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := NewMethods(&mockError{})
		assert.NotNil(t, s.Delete(context.Background(), primitive.NewObjectID()))
	})
}

func TestNotificationMethodDb_Get(t *testing.T) {
	t.Run("Should: return no error", func(t *testing.T) {
		s := NewMethods(&mockSuccess{})
		_, err := s.Get(context.Background(), primitive.NewObjectID())
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := NewMethods(&mockError{})
		_, err := s.Get(context.Background(), primitive.NewObjectID())
		assert.NotNil(t, err)
	})
}