package cache

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("Should: create new cache", func(t *testing.T) {
		c, err := New(&redis.Client{})
		assert.Nil(t, err)
		assert.NotNil(t, c)
	})
	t.Run("Should: return err", func(t *testing.T) {
		_, err := New("")
		assert.ErrorContains(t, err, "cannot convert to redis db connection")
	})
}

func TestRedis_InsertSchedule(t *testing.T) {
	db, mock := redismock.NewClientMock()
	c := &Redis{
		Client: db,
	}

	data := &apiPb.InsertScheduleWithIdRequest{
		Id:            "id",
		ScheduledNext: nil,
	}

	mock.Regexp().ExpectSet(data.Id, "^[0-9]{10}$", time.Duration(0)).SetVal("")
	err := c.InsertSchedule(data)
	assert.Nil(t, err)

	data2 := &apiPb.GetScheduleWithIdRequest{
		Id: "id",
	}
	mock.ExpectGet(data2.Id).SetVal("1294960918")
	res, err := c.GetScheduleById(data2)
	assert.Nil(t, err)
	assert.Equal(t, res, &apiPb.GetScheduleWithIdResponse{
		ScheduledNext: &timestamp.Timestamp{
			Seconds: 1294960918,
			Nanos:   0,
		},
	})

	data3 := &apiPb.DeleteScheduleWithIdRequest{
		Id: "id",
	}
	mock.ExpectDel(data3.Id).SetVal(0)
	err = c.DeleteScheduleById(data3)
	assert.Nil(t, err)
}

func TestRedis_InsertScheduleError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	c := &Redis{
		Client: db,
	}

	data := &apiPb.InsertScheduleWithIdRequest{
		Id:            "id",
		ScheduledNext: nil,
	}

	mock.Regexp().ExpectSet(
		data.Id, "^[0-9]{10}$",
		time.Duration(0)).SetErr(errors.New("InsertScheduleWithIdRequest"))
	err := c.InsertSchedule(data)
	assert.ErrorContains(t, err, "InsertScheduleWithIdRequest")

	data2 := &apiPb.GetScheduleWithIdRequest{
		Id: "id",
	}
	mock.ExpectGet(data2.Id).SetErr(errors.New("GetScheduleWithIdRequest"))
	res, err := c.GetScheduleById(data2)
	assert.ErrorContains(t, err, "GetScheduleWithIdRequest")
	assert.Nil(t, res, "GetScheduleWithIdRequest")

	mock.ExpectGet(data2.Id).SetVal("string")
	res, err = c.GetScheduleById(data2)
	assert.ErrorContains(t, err, "strconv.ParseInt: parsing \"string\"")
	assert.Nil(t, res, "GetScheduleWithIdRequest")

	data3 := &apiPb.DeleteScheduleWithIdRequest{
		Id: "id",
	}
	mock.ExpectDel(data3.Id).SetErr(errors.New("DeleteScheduleWithIdRequest"))
	err = c.DeleteScheduleById(data3)
	assert.ErrorContains(t, err, "DeleteScheduleWithIdRequest")
}
