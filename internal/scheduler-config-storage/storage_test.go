package scheduler_config_storage

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

var (
	basicError = errors.New("")
)

type mockOk struct {
}

func (m mockOk) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, nil
}

func (m mockOk) FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error {
	return nil
}

func (m mockOk) FindAll(ctx context.Context, predicate bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error {
	return nil
}

func (m mockOk) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

type mockError struct {
}

func (m mockError) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, basicError
}

func (m mockError) FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error {
	return basicError
}

func (m mockError) FindAll(ctx context.Context, predicate bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error {
	return basicError
}

func (m mockError) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, basicError
}

func TestNew(t *testing.T) {
	t.Run("Should: implement interface", func(t *testing.T) {
		s := New(nil)
		assert.Implements(t, (*Storage)(nil), s)
	})
}

func TestStorage_Add(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&mockOk{})
		assert.Equal(t, nil, s.Add(context.Background(), nil))
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		assert.NotEqual(t, nil, s.Add(context.Background(), nil))
	})
}

func TestStorage_Get(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&mockOk{})
		_, err := s.Get(context.Background(), primitive.NewObjectID())
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		_, err := s.Get(context.Background(), primitive.NewObjectID())
		assert.NotEqual(t, nil, err)
	})
}

func TestStorage_GetAll(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&mockOk{})
		_, err := s.GetAll(context.Background())
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		_, err := s.GetAll(context.Background())
		assert.NotEqual(t, nil, err)
	})
}

func TestStorage_GetAllForSync(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&mockOk{})
		_, err := s.GetAllForSync(context.Background())
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		_, err := s.GetAllForSync(context.Background())
		assert.NotEqual(t, nil, err)
	})
}

func TestStorage_Remove(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&mockOk{})
		err := s.Remove(context.Background(), primitive.NewObjectID())
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		err := s.Remove(context.Background(), primitive.NewObjectID())
		assert.NotEqual(t, nil, err)
	})
}

func TestStorage_Run(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&mockOk{})
		err := s.Run(context.Background(), primitive.NewObjectID())
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		err := s.Run(context.Background(), primitive.NewObjectID())
		assert.NotEqual(t, nil, err)
	})
}

func TestStorage_Stop(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&mockOk{})
		err := s.Stop(context.Background(), primitive.NewObjectID())
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		err := s.Stop(context.Background(), primitive.NewObjectID())
		assert.NotEqual(t, nil, err)
	})
}
