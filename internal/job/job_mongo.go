package job

import (
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var (
	timeout = time.Second * 5
)

type mongoJob struct {
	schedulerID string
	dbConfig    *scheduler_config_storage.DbConfig
	mongo       MongoConnector
}

type mongoError struct {
	schedulerID string
	startTime   *timestamp.Timestamp
	endTime     *timestamp.Timestamp
	code        apiPb.SchedulerCode
	description string
	location    string
}

func newMongoError(schedulerID string, startTime, endTime *timestamp.Timestamp, code apiPb.SchedulerCode, description, location string) CheckError {
	return &mongoError{
		schedulerID: schedulerID,
		startTime:   startTime,
		endTime:     endTime,
		code:        code,
		description: description,
		location:    location,
	}
}

func (s *mongoError) GetLogData() *apiPb.SchedulerResponse {
	var err *apiPb.SchedulerSnapshot_Error
	if s.code == apiPb.SchedulerCode_ERROR {
		err = &apiPb.SchedulerSnapshot_Error{
			Message: s.description,
		}
	}
	return &apiPb.SchedulerResponse{
		SchedulerId: s.schedulerID,
		Snapshot: &apiPb.SchedulerSnapshot{
			Code:  s.code,
			Error: err,
			Type:  apiPb.SchedulerType_MONGO,
			Meta: &apiPb.SchedulerSnapshot_MetaData{
				StartTime: s.startTime,
				EndTime:   s.endTime,
			},
		},
	}
}

func ExecMongo(schedulerId string, config *scheduler_config_storage.DbConfig, mongo MongoConnector) CheckError {
	startTime := timestamppb.Now()

	clientOptions := options.Client().ApplyURI(config.Host)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return newMongoError(schedulerId, startTime, timestamppb.Now(), apiPb.SchedulerCode_ERROR, mongoConnectionError.Error(), config.Host)
	}

	err = mongo.Ping(context.TODO(), nil)
	if err != nil {
		return newMongoError(schedulerId, startTime, timestamppb.Now(), apiPb.SchedulerCode_ERROR, mongoPingError.Error(), config.Host)
	}

	return newMongoError(schedulerId, startTime, timestamppb.Now(), apiPb.SchedulerCode_OK, "", config.Host)
}

type MongoConnector interface {
	Connect(ctx context.Context, opts ...*options.ClientOptions) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error
}

type MongoConnection struct {
	Client   *mongo.Client
	Connect_ func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error)
}

func NewMongoConnection() MongoConnection {
	return MongoConnection{
		Connect_: mongo.Connect,
	}
}

func (m MongoConnection) Connect(ctx context.Context, opts ...*options.ClientOptions) error {
	if client, err := m.Connect_(ctx, opts...); err == nil {
		m.Client = client
		return err
	} else {
		return err
	}
}

func (m MongoConnection) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return m.Client.Ping(ctx, rp)
}
