package database

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/protobuf/ptypes"
	"github.com/jinzhu/gorm"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

//docker run -d --rm --name postgres -e POSTGRES_USER="user" -e POSTGRES_PASSWORD="password" -e POSTGRES_DB="database" -p 5432:5432 postgres
var (
	postgr      = &postgres{}
	postgrWrong = &postgres{}
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	postgr.db = s.DB

	s.DB.LogMode(true)
}

func TestPostgres_NewClient(t *testing.T) {
	t.Run("wrongPostgress", func(t *testing.T) {
		err := postgrWrong.newClient(func() (db *gorm.DB, e error) {
			return gorm.Open(
				"postgres",
				fmt.Sprintf("host=lkl port=00 user=us dbname=dbn password=ps connect_timeout=10 sslmode=disable"))
		})
		assert.Error(t, err)
	})
}

func (s *Suite) Test_InsertMetaData() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbSnapshotCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "meta_data")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	correctTime, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		panic("Time convertion error")
	}
	err = postgr.InsertSnapshot(&apiPb.SchedulerResponse{
		SchedulerId:          "schId",
		Snapshot:             &apiPb.Snapshot{
			Code:                 0,
			Type:                 0,
			Error:                &apiPb.Snapshot_SnapshotError{
				Message:              "message",
			},
			Meta:                 &apiPb.Snapshot_MetaData{
				StartTime:            correctTime,
				EndTime:              correctTime,
				Value:                nil,
			},
		},
	})
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetMetaData() {
	var (
		id = "1"
	)
	query := fmt.Sprintf(`SELECT * FROM "%s" WHERE "%s"."deleted_at" IS NULL`, dbSnapshotCollection, dbSnapshotCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnRows(rows)

	_, err := postgr.GetSnapshots(id)
	require.NoError(s.T(), err)
}

func (s *Suite) Test_InsertStatRequest() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "cpu_infos")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "memory_infos")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "memories")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "memories")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "disk_infos")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, "net_infos")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	statReq, _ := ConvertFromPostgressStatRequest(
		&StatRequest{
			CpuInfo: []*CpuInfo{
				{
					Load: 0,
				},
			},
			MemoryInfo: &MemoryInfo{
				StatRequestID: 0,
				Mem: &Memory{
					Total:       0,
					Used:        0,
					Free:        0,
					Shared:      0,
					UsedPercent: 0,
				},
				Swap: &Memory{
					Total:       0,
					Used:        0,
					Free:        0,
					Shared:      0,
					UsedPercent: 0,
				},
			},
			DiskInfo: []*DiskInfo{
				{
					Name:        "",
					Total:       0,
					Free:        0,
					Used:        0,
					UsedPercent: 0,
				},
			},
			NetInfo:  []*NetInfo{
				{
					Name:          "",
					BytesSent:     0,
					BytesRecv:     0,
					PacketsSent:   0,
					PacketsRecv:   0,
					ErrIn:         0,
					ErrOut:        0,
					DropIn:        0,
					DropOut:       0,
				},
			},
			Time:     time.Now(),
		})
	err := postgr.InsertStatRequest(statReq)
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetStatRequest() {
	var (
		id = "1"
	)
	query := fmt.Sprintf(`SELECT * FROM "%s" WHERE "%s"."deleted_at" IS NULL`, dbStatRequestCollection, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnRows(rows)

	_, err := postgr.GetStatRequest(id)
	require.NoError(s.T(), err)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func TestPostgres_Migrate(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		err := postgrWrong.Migrate()
		assert.Error(t, err)
	})
}

func TestPostgres_InsertMetaData(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		err := postgrWrong.InsertSnapshot(&apiPb.SchedulerResponse{})
		assert.Error(t, err)
	})
}

func TestPostgres_GetMetaData(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := postgrWrong.GetSnapshots("")
		assert.Error(t, err)
	})
}

func TestPostgres_InsertStatRequest(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		statReq, _ := ConvertFromPostgressStatRequest(
			&StatRequest{
				CpuInfo: []*CpuInfo{
					{
						Load: 0,
					},
				},
				MemoryInfo: &MemoryInfo{
					StatRequestID: 0,
					Mem: &Memory{
						Total:       0,
						Used:        0,
						Free:        0,
						Shared:      0,
						UsedPercent: 0,
					},
					Swap: &Memory{
						Total:       0,
						Used:        0,
						Free:        0,
						Shared:      0,
						UsedPercent: 0,
					},
				},
				DiskInfo: []*DiskInfo{
					{
						Name:        "",
						Total:       0,
						Free:        0,
						Used:        0,
						UsedPercent: 0,
					},
				},
				NetInfo:  []*NetInfo{
					{
						Name:          "",
						BytesSent:     0,
						BytesRecv:     0,
						PacketsSent:   0,
						PacketsRecv:   0,
						ErrIn:         0,
						ErrOut:        0,
						DropIn:        0,
						DropOut:       0,
					},
				},
				Time:     time.Now(),
			})
		err := postgrWrong.InsertStatRequest(statReq)
		assert.Error(t, err)
	})
}

func TestPostgres_GetStatRequest(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := postgrWrong.GetStatRequest("")
		assert.Error(t, err)
	})
}
