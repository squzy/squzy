package checks

import (
	"context"
	storagePb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: implement interface of the LoggerServer", func(t *testing.T) {
		s := New(nil)
		assert.Implements(t, (*storagePb.LoggerServer)(nil), s)
	})
}

func TestServer_SendLogMessage(t *testing.T) {
	t.Run("Should: not return error on message sending", func(t *testing.T) {
		s := New(nil)
		res, err := s.SendLogMessage(context.Background(), &storagePb.SendLogMessageRequest{
			Log: &storagePb.Log{
				Code: storagePb.StatusCode_OK,
			},
		})
		assert.Equal(t, nil, err)
		assert.Equal(t, true, res.Success)
	})
}
