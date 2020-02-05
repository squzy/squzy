package job

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var (
	timeout = time.Second * 5
)

type mongoJob struct {
	url   string
	mongo mongoConnectorI
	ping  mongoPingI
}

func NewMongoJob(url string) Job {
	return &mongoJob{
		url:   url,
		mongo: &mongoConnector{},
		ping:  &mongoPing{},
	}
}

type mongoError struct {
	logId       string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        clientPb.StatusCode
	description string
	location    string
}

func newMongoError(logId string, startTime, endTime *timestamp.Timestamp, code clientPb.StatusCode, description string, location string) CheckError {
	return &mongoError{
		logId:       logId,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		location:    location,
	}
}

func (m *mongoError) GetLogData() *clientPb.Log {
	return &clientPb.Log{
		Code:        m.code,
		Description: m.description,
		Meta: &clientPb.MetaData{
			Id:       uuid.New().String(),
			Location: m.location,
		},
	}
}

func (j *mongoJob) Do() CheckError {
	logId := uuid.New().String()
	startTime := ptypes.TimestampNow()

	clientOptions := options.Client().ApplyURI(j.url)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client, err := j.mongo.connect(ctx, clientOptions)
	if err != nil {
		return newMongoError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_Error, mongoConnectionError.Error(), j.url)
	}

	fmt.Println(client)
	err = j.ping.ping(client, context.TODO(), nil)
	if err != nil {
		return newMongoError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_Error, mongoPingError.Error(), j.url)
	}

	return newMongoError(logId, startTime, ptypes.TimestampNow(), clientPb.StatusCode_OK, "", j.url)
}

type mongoConnectorI interface {
	connect(ctx context.Context, opts *options.ClientOptions) (*mongo.Client, error)
}

type mongoConnector struct {
}

func (mongoConnector) connect(ctx context.Context, opts *options.ClientOptions) (*mongo.Client, error) {
	return mongo.Connect(ctx, opts)
}

type mongoPingI interface {
	ping(client *mongo.Client, ctx context.Context, rp *readpref.ReadPref) error
}

type mongoPing struct{}

func (mongoPing) ping(client *mongo.Client, ctx context.Context, rp *readpref.ReadPref) error {
	return client.Ping(ctx, rp)
}
