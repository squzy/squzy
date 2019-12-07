package job

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type jobMongo struct {
	url string
}

func NewMongoJob(url string) Job {
	return &jobMongo{
		url: url,
	}
}

type mongoError struct {
	time        *timestamp.Timestamp
	code        clientPb.StatusCode
	description string
	location    string
}

func newMongoError(time *timestamp.Timestamp, code clientPb.StatusCode, description string, location string) CheckError  {
	return &mongoError{
		time:        time,
		code:        code,
		description: description,
		location:    location,
	}
}
func (m *mongoError) GetLogData() *clientPb.Log {
	return &clientPb.Log{
		Code:                 m.code,
		Description:          m.description,
		Meta:                 &clientPb.MetaData{
			Id:                   uuid.New().String(),
			Location:             m.location,
			Time:                 m.time,
		},
	}
}

func (j *jobMongo) Do() CheckError {
	clientOptions := options.Client().ApplyURI(j.url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return newMongoError(ptypes.TimestampNow(), clientPb.StatusCode_Error, mongoConnectionError.Error(), j.url)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return newMongoError(ptypes.TimestampNow(), clientPb.StatusCode_Error, mongoPingError.Error(), j.url)
	}

	return newMongoError(ptypes.TimestampNow(), clientPb.StatusCode_OK, "", j.url)
}
