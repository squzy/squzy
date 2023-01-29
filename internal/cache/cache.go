package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Redis struct {
	Client *redis.Client
}

type Cache interface {
	InsertSchedule(data *apiPb.InsertScheduleWithIdRequest) error
	GetScheduleById(data *apiPb.GetScheduleWithIdRequest) (*apiPb.GetScheduleWithIdResponse, error)
}

func New(ca interface{}) (Cache, error) {
	client, ok := ca.(*redis.Client)
	if !ok {
		return nil, errors.New("cannot convert to redis db connection")
	}
	return &Redis{
		Client: client,
	}, nil
}

func (c *Redis) InsertSchedule(data *apiPb.InsertScheduleWithIdRequest) error {
	return c.Client.Set(context.Background(), data.GetId(), time.Now().Unix(), time.Duration(0)).Err()
}

func (c *Redis) GetScheduleById(data *apiPb.GetScheduleWithIdRequest) (*apiPb.GetScheduleWithIdResponse, error) {
	res := c.Client.Get(context.Background(), data.GetId())

	if err := res.Err(); err != nil {
		return nil, err
	}

	uTime, err := res.Int64()
	if err != nil {
		return nil, err
	}

	return &apiPb.GetScheduleWithIdResponse{
		ScheduleTime: &timestamppb.Timestamp{
			Seconds: uTime,
			Nanos:   0,
		},
	}, nil
}
