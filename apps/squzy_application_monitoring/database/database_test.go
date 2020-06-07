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

type mockOk struct {
}

type mockError struct {
}

type mockNotFoundOk struct {
}

type mockNotFoundError struct {
}

func (m mockNotFoundError) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, errors.New("")
}

func (m mockNotFoundError) FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error {
	return mongo.ErrNoDocuments
}

func (m mockNotFoundError) FindAll(ctx context.Context, predicate bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error {
	panic("implement me")
}

func (m mockNotFoundError) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	panic("implement me")
}

func (m mockNotFoundOk) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, nil
}

func (m mockNotFoundOk) FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error {
	return mongo.ErrNoDocuments
}

func (m mockNotFoundOk) FindAll(ctx context.Context, predicate bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error {
	panic("implement me")
}

func (m mockNotFoundOk) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	panic("implement me")
}

func (m mockError) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	panic("implement me")
}

func (m mockError) FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error {
	return errors.New("")
}

func (m mockError) FindAll(ctx context.Context, predicate bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error {
	return errors.New("asf")
}

func (m mockError) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	panic("implement me")
}

func (m mockOk) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	panic("implement me")
}

func (m mockOk) FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error {
	return nil
}

func (m mockOk) FindAll(ctx context.Context, predicate bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error {
	return nil
}

func (m mockOk) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	panic("implement me")
}

func TestNew(t *testing.T) {
	t.Run("Should: not be nil", func(t *testing.T) {
		s := New(nil)
		assert.NotNil(t, s)
	})
}

func TestDb_FindAllApplication(t *testing.T) {
	t.Run("Should: return all application", func(t *testing.T) {
		s := New(&mockOk{})
		_, err := s.FindAllApplication(context.Background())
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		_, err := s.FindAllApplication(context.Background())
		assert.NotNil(t, err)
	})
}

func TestDb_FindApplicationById(t *testing.T) {
	t.Run("Should: return application", func(t *testing.T) {
		s := New(&mockOk{})
		_, err := s.FindApplicationById(context.Background(), primitive.NewObjectID())
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		_, err := s.FindApplicationById(context.Background(), primitive.NewObjectID())
		assert.NotNil(t, err)
	})
}

func TestDb_FindApplicationByName(t *testing.T) {
	t.Run("Should: return application", func(t *testing.T) {
		s := New(&mockOk{})
		_, err := s.FindApplicationByName(context.Background(), "")
		assert.Nil(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		_, err := s.FindApplicationByName(context.Background(), "primitive.NewObjectID()")
		assert.NotNil(t, err)
	})
}

func TestDb_FindOrCreate(t *testing.T) {
	t.Run("Should: return application immediately", func(t *testing.T) {
		s := New(&mockOk{})
		_, err := s.FindOrCreate(context.Background(), "", "")
		assert.Nil(t, err)
	})
	t.Run("Should: return error because internal error", func(t *testing.T) {
		s := New(&mockError{})
		_, err := s.FindOrCreate(context.Background(), "", "")
		assert.NotNil(t, err)
	})
	t.Run("Should: create new application", func(t *testing.T) {
		s := New(&mockNotFoundOk{})
		_, err := s.FindOrCreate(context.Background(), "", "")
		assert.Nil(t, err)
	})
	t.Run("Should: return error because internal error while creation", func(t *testing.T) {
		s := New(&mockNotFoundError{})
		_, err := s.FindOrCreate(context.Background(), "", "")
		assert.NotNil(t, err)
	})
}
