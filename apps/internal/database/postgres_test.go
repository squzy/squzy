package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//docker run -d --rm --name postgres -e POSTGRES_USER="user" -e POSTGRES_PASSWORD="password" -e POSTGRES_DB="database" -p 5432:5432 postgres
var (
	postgr = &postgres{
		host:     "localhost",
		port:     "5432",
		user:     "user",
		password: "password",
		dbname:   "database",
	}
)

func TestPostgres_Migrate(t *testing.T) {
	err := postgr.InsertMetaData(&MetaData{
		Location:  "loc",
		Port:      20,
		StartTime: &time.Time{},
		EndTime:   &time.Time{},
		Type:      20,
	})
	assert.NoError(t, err)
	now := time.Now()
	err = postgr.InsertStatRequest(&StatRequest{
		CpuInfo: []*CpuInfo{{Load: 10,}, {Load: 20,}},
		MemoryInfo: &MemoryInfo{
			Mem: &Memory{
				Total:       100,
				Used:        5,
				Free:        90,
				Shared:      5,
				UsedPercent: 5,
			},
			Swap: &Memory{
				Total:       200,
				Used:        10,
				Free:        180,
				Shared:      10,
				UsedPercent: 10,
			},
		},
		DiskInfo: []*DiskInfo{{
			Name:        "/disk",
			Total:       500,
			Free:        400,
			Used:        100,
			UsedPercent: 20,
		}},
		NetInfo: []*NetInfo{{
			Name:        "localhost",
			BytesSent:   500,
			BytesRecv:   200,
			PacketsSent: 10,
			PacketsRecv: 20,
			ErrIn:       0,
			ErrOut:      0,
			DropIn:      0,
			DropOut:     0,
		}},
		Time: &now,
	})
	assert.NoError(t, err)
}
