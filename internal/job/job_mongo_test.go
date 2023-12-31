package job

import (
	"context"
	"errors"
	scheduler_config_storage "github.com/squzy/squzy/internal/scheduler-config-storage"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

type mongoMockOk struct {
}

func (m mongoMockOk) Connect(ctx context.Context, opts ...*options.ClientOptions) error {
	return nil
}

func (m mongoMockOk) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return nil
}

type mongoMockErr struct {
}

func (m mongoMockErr) Connect(ctx context.Context, opts ...*options.ClientOptions) error {
	return nil
}

func (m mongoMockErr) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return errors.New("Ping")
}

func TestMongoJob_Do(t *testing.T) {
	t.Run("Test: mongoJob", func(t *testing.T) {
		t.Run("Should: return error connecting", func(t *testing.T) {
			j := mongoJob{
				dbConfig: &scheduler_config_storage.DbConfig{},
				mongo:    NewMongoConnection(),
			}
			err := ExecMongo(j.schedulerID, j.dbConfig, j.mongo)
			expected := apiPb.SchedulerCode_ERROR
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return error ping", func(t *testing.T) {
			j := mongoJob{
				dbConfig: &scheduler_config_storage.DbConfig{},
				mongo:    &mongoMockErr{},
			}
			err := ExecMongo(j.schedulerID, j.dbConfig, j.mongo)
			expected := apiPb.SchedulerCode_ERROR
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return no error", func(t *testing.T) {
			j := mongoJob{
				dbConfig: &scheduler_config_storage.DbConfig{},
				mongo:    &mongoMockOk{},
			}
			err := ExecMongo(j.schedulerID, j.dbConfig, j.mongo)
			expected := apiPb.SchedulerCode_OK
			actual := err.GetLogData().Snapshot.Code
			assert.EqualValues(t, expected, actual)
		})
	})
}

var connectMock = func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
	return &mongo.Client{}, nil
}

func TestMongoConnection_Connect(t *testing.T) {
	t.Run("Should: get connection", func(t *testing.T) {
		config := &scheduler_config_storage.DbConfig{}
		m := MongoConnection{Connect_: connectMock}

		clientOptions := options.Client().ApplyURI(config.Host)
		err := m.Connect(context.Background(), clientOptions)
		assert.Nil(t, err)
		assert.NotNil(t, t, m.Client)
	})
}
