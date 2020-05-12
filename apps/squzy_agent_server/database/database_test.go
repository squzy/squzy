package database

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"testing"
)

type mockError struct {
}

func (m mockError) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, errors.New("")
}

func (m mockError) FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error {
	return errors.New("")
}

func (m mockError) FindAll(ctx context.Context, predicate bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error {
	return errors.New("")
}

func (m mockError) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, errors.New("")
}

type mockOk struct {
}

func (m mockOk) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, nil
}

func (m mockOk) FindOne(ctx context.Context, filter interface{}, structToDeserialize interface{}, opts ...*options.FindOneOptions) error {
	return nil
}

func (m mockOk) FindAll(ctx context.Context, predicate bson.M, structToDeserialize interface{}, opts ...*options.FindOptions) error {
	agents := []*AgentDao{
		{
			Id:        primitive.NewObjectID(),
			AgentName: "test",
			Status:    0,
			HostInfo:  &HostInfo{
				HostName:     "",
				Os:           "",
				PlatFormInfo: &PlatFormInfo{
					Name:    "",
					Family:  "",
					Version: "",
				},
			},
			History:   nil,
		},
		{
			Id:        primitive.NewObjectID(),
			AgentName: "",
			Status:    0,
			HostInfo:  nil,
			History:   nil,
		},
	}
	val := reflect.ValueOf(structToDeserialize)
	val.Elem().Set(reflect.ValueOf(agents))
	return nil
}

func (m mockOk) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

func TestNew(t *testing.T) {
	t.Run("Should: implement interface", func(t *testing.T) {
		s := New(nil)
		assert.Implements(t, (*Database)(nil), s)
	})
}

func TestDb_Add(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&mockOk{})
		_, err := s.Add(context.Background(), &apiPb.RegisterRequest{
			AgentName: "",
			Time: ptypes.TimestampNow(),
			HostInfo: &apiPb.HostInfo{
				HostName: "",
				Os:       "",
				PlatformInfo: &apiPb.PlatformInfo{
					Name:                 "",
					Family:               "",
					Version:              "",
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
			},
		})
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return error because time is nil", func(t *testing.T) {
		s := New(&mockOk{})
		_, err := s.Add(context.Background(), &apiPb.RegisterRequest{
			AgentName: "",
			Time: nil,
			HostInfo: &apiPb.HostInfo{
				HostName: "",
				Os:       "",
				PlatformInfo: &apiPb.PlatformInfo{
					Name:                 "",
					Family:               "",
					Version:              "",
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
			},
		})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		_, err := s.Add(context.Background(), &apiPb.RegisterRequest{
			AgentName: "",
			Time: ptypes.TimestampNow(),
			HostInfo: &apiPb.HostInfo{
				HostName: "",
				Os:       "",
				PlatformInfo: &apiPb.PlatformInfo{
					Name:                 "",
					Family:               "",
					Version:              "",
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
			},
		})
		assert.NotEqual(t, nil, err)
	})
}

func TestDb_UpdateStatus(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&mockOk{})
		err := s.UpdateStatus(context.Background(), primitive.NewObjectID(), 0, ptypes.TimestampNow())
		assert.Equal(t, nil, err)
	})
	t.Run("Should: not return error", func(t *testing.T) {
		s := New(&mockOk{})
		err := s.UpdateStatus(context.Background(), primitive.NewObjectID(), apiPb.AgentStatus_DISCONNECTED, ptypes.TimestampNow())
		assert.Equal(t, nil, err)
	})
	t.Run("Should: return because time error", func(t *testing.T) {
		s := New(&mockOk{})
		err := s.UpdateStatus(context.Background(), primitive.NewObjectID(), 0, nil)
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		err := s.UpdateStatus(context.Background(), primitive.NewObjectID(), 0, ptypes.TimestampNow())
		assert.NotEqual(t, nil, err)
	})
}

func TestDb_GetAll(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		_, err := s.GetAll(context.Background(), bson.M{})
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return agent", func(t *testing.T) {
		s := New(&mockOk{})
		res, err := s.GetAll(context.Background(), bson.M{})
		assert.Equal(t, nil, err)
		assert.Equal(t, 2, len(res))
	})
}

func TestDb_GetById(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		s := New(&mockError{})
		_, err := s.GetById(context.Background(), primitive.NewObjectID())
		assert.NotEqual(t, nil, err)
	})
	t.Run("Should: return agent", func(t *testing.T) {
		s := New(&mockOk{})
		_, err := s.GetById(context.Background(), primitive.NewObjectID())
		assert.Equal(t, nil, err)
	})
}
