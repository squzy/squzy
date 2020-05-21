package scheduler_config_storage

import (
	"context"
	"github.com/squzy/mongo_helper"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GrpcConfig struct {
	Service string `bson:"service"`
	Host    string `bson:"host"`
	Port    int32  `bson:"port"`
}

type HTTPConfig struct {
	Method     string            `bson:"string"`
	URL        string            `bson:"url"`
	Headers    map[string]string `bson:"headers"`
	StatusCode int32             `bson:"statusCode"`
}

type HTTPValueConfig struct {
	Method    string            `bson:"method"`
	URL       string            `bson:"url"`
	Headers   map[string]string `bson:"headers"`
	Selectors []*Selectors      `bson:"selectors"`
}

type Selectors struct {
	Type apiPb.HttpJsonValueConfig_JsonValueParseType `bson:"type"`
	Path string                                       `bson:"path"`
}

type TCPConfig struct {
	Host string `bson:"host"`
	Port int32  `bson:"port"`
}

type SiteMapConfig struct {
	URL         string `bson:"url"`
	Concurrency int32  `bson:"concurrency"`
}

type SchedulerConfig struct {
	ID              primitive.ObjectID    `bson:"_id"`
	Type            apiPb.SchedulerType   `bson:"type"`
	Status          apiPb.SchedulerStatus `bson:"status"`
	Interval        int32                 `bson:"interval"`
	Timeout         int32                 `bson:"timeout"`
	TCPConfig       *TCPConfig            `bson:"tcpConfig,omitempty"`
	SiteMapConfig   *SiteMapConfig        `bson:"siteMapConfig,omitempty"`
	GrpcConfig      *GrpcConfig           `bson:"grpcConfig,omitempty"`
	HTTPConfig      *HTTPConfig           `bson:"httpConfig,omitempty"`
	HTTPValueConfig *HTTPValueConfig      `bson:"httpValueConfig,omitempty"`
}

type Storage interface {
	Get(ctx context.Context, schedulerID primitive.ObjectID) (*SchedulerConfig, error)
	Add(ctx context.Context, config *SchedulerConfig) error
	Remove(ctx context.Context, schedulerID primitive.ObjectID) error
	Run(ctx context.Context, schedulerID primitive.ObjectID) error
	Stop(ctx context.Context, schedulerID primitive.ObjectID) error
	GetAll(ctx context.Context) ([]*SchedulerConfig, error)
	GetAllForSync(ctx context.Context) ([]*SchedulerConfig, error)
}

type storage struct {
	connector mongo_helper.Connector
}

var (
	statusForAction = []apiPb.SchedulerStatus{
		apiPb.SchedulerStatus_STOPPED,
		apiPb.SchedulerStatus_RUNNED,
	}
)

func (s *storage) GetAllForSync(ctx context.Context) ([]*SchedulerConfig, error) {
	configs := []*SchedulerConfig{}
	err := s.connector.FindAll(ctx, bson.M{
		"status": bson.M{
			"$in": statusForAction,
		},
	}, &configs)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (s *storage) GetAll(ctx context.Context) ([]*SchedulerConfig, error) {
	configs := []*SchedulerConfig{}
	err := s.connector.FindAll(ctx, bson.M{}, &configs)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (s *storage) Add(ctx context.Context, config *SchedulerConfig) error {
	_, err := s.connector.InsertOne(ctx, config)
	return err
}

func (s *storage) Remove(ctx context.Context, schedulerID primitive.ObjectID) error {
	_, err := s.connector.UpdateOne(ctx, bson.M{
		"_id": schedulerID,
	}, bson.M{
		"$set": bson.M{
			"status": apiPb.SchedulerStatus_REMOVED,
		},
	})
	return err
}

func (s *storage) Run(ctx context.Context, schedulerID primitive.ObjectID) error {
	_, err := s.connector.UpdateOne(ctx, bson.M{
		"_id": schedulerID,
		"status": bson.M{
			"$in": statusForAction,
		},
	}, bson.M{
		"$set": bson.M{
			"status": apiPb.SchedulerStatus_RUNNED,
		},
	})
	return err
}

func (s *storage) Stop(ctx context.Context, schedulerID primitive.ObjectID) error {
	_, err := s.connector.UpdateOne(ctx, bson.M{
		"_id": schedulerID,
		"status": bson.M{
			"$in": statusForAction,
		},
	}, bson.M{
		"$set": bson.M{
			"status": apiPb.SchedulerStatus_STOPPED,
		},
	})
	return err
}

func (s *storage) Get(ctx context.Context, schedulerID primitive.ObjectID) (*SchedulerConfig, error) {
	config := &SchedulerConfig{}
	err := s.connector.FindOne(ctx, bson.M{
		"_id": schedulerID,
	}, config)
	if err != nil {
		return nil, err
	}
	return config, err
}

func New(
	connector mongo_helper.Connector,
) Storage {
	return &storage{
		connector: connector,
	}
}
