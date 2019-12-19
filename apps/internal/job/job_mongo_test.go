package job

import (
	"context"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

type mongoMockConnectOk struct{}

func (mongoMockConnectOk) connect(context.Context, *options.ClientOptions) (*mongo.Client, error) {
	client, _ := mongo.NewClient()
	return client, nil
}

type mongoMockPingOk struct{}

func (mongoMockPingOk) ping(*mongo.Client, context.Context, *readpref.ReadPref) error {
	return nil
}

func TestMongoJob_Do(t *testing.T) {
	t.Run("Test: mongoJob", func(t *testing.T) {
		t.Run("Should: return error connecting", func(t *testing.T) {
			j := NewMongoJob("")
			err := j.Do()
			expected := clientPb.StatusCode_Error
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return error ping", func(t *testing.T) {
			j := mongoJob{
				url:   "",
				mongo: &mongoMockConnectOk{},
				ping:  &mongoPing{},
			}
			err := j.Do()
			expected := clientPb.StatusCode_Error
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
		t.Run("Should: return no error", func(t *testing.T) {
			j := mongoJob{
				url:   "",
				mongo: &mongoMockConnectOk{},
				ping:  &mongoMockPingOk{},
			}
			err := j.Do()
			expected := clientPb.StatusCode_OK
			actual := err.GetLogData().Code
			assert.EqualValues(t, expected, actual)
		})
	})
}
