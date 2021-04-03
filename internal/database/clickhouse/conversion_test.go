package clickhouse

import (
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_convertFromIncidentHistory(t *testing.T) {
	t.Run("Should: return nil", func(t *testing.T) {
		res := convertToIncident(nil)
		assert.Nil(t, res)
	})
}

func Test_convertToIncidentHistory(t *testing.T) {
	t.Run("Should: return empty res", func(t *testing.T) {
		res, _, _ := convertToIncidentHistories(nil)
		assert.Equal(t, 0, len(res))
	})
	t.Run("Should: return empty res", func(t *testing.T) {
		maxValidSeconds := 253402300800
		res, _, _ := convertToIncidentHistories([]*apiPb.Incident_HistoryItem{{
			Timestamp: &tspb.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		}})
		assert.Equal(t, 0, len(res))
	})
}
