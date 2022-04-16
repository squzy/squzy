package clickhouse

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"regexp"
	"testing"
	"time"
)

var (
	wdbStatRequest, _      = sql.Open("clickhouse", "tcp://user:password@lkl:00/debug=true&clicks?read_timeout=10&write_timeout=10")
	clickhWrongStatRequest = &Clickhouse{
		wdbStatRequest,
	}
	clickStatRequest = &Clickhouse{}
)

type SuiteStatRequest struct {
	suite.Suite
	DB   *sql.DB
	mock sqlmock.Sqlmock
}

func (s *SuiteStatRequest) SetupSuite() {
	var err error

	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	clickStatRequest.Db = s.DB
}

func TestClickhouse_InsertStatRequest(t *testing.T) {
	t.Run("Should: return conv error", func(t *testing.T) {
		err := clickStatRequest.InsertStatRequest(&apiPb.Metric{})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		err := clickhWrongStatRequest.InsertStatRequest(&apiPb.Metric{
			CpuInfo: &apiPb.CpuInfo{
				Cpus: []*apiPb.CpuInfo_CPU{{}},
			},
			MemoryInfo: &apiPb.MemoryInfo{
				Mem:  &apiPb.MemoryInfo_Memory{},
				Swap: &apiPb.MemoryInfo_Memory{},
			},
			DiskInfo: &apiPb.DiskInfo{
				Disks: map[string]*apiPb.DiskInfo_Disk{},
			},
			NetInfo: &apiPb.NetInfo{
				Interfaces: map[string]*apiPb.NetInfo_Interface{},
			},
			Time: timestamp.Now(),
		})
		assert.Error(t, err)
	})
}

func (s *SuiteStatRequest) Test_InsertStatRequest() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCpuInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestMemoryInfoMemCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestMemoryInfoSwapCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestDiskInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestNetInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := clickStatRequest.InsertStatRequest(&apiPb.Metric{
		CpuInfo: &apiPb.CpuInfo{
			Cpus: []*apiPb.CpuInfo_CPU{{}},
		},
		MemoryInfo: &apiPb.MemoryInfo{
			Mem:  &apiPb.MemoryInfo_Memory{},
			Swap: &apiPb.MemoryInfo_Memory{},
		},
		DiskInfo: &apiPb.DiskInfo{
			Disks: map[string]*apiPb.DiskInfo_Disk{
				"": {},
			},
		},
		NetInfo: &apiPb.NetInfo{
			Interfaces: map[string]*apiPb.NetInfo_Interface{
				"": {},
			},
		},
		Time: timestamp.Now(),
	})
	require.NoError(s.T(), err)
}

func (s *SuiteStatRequest) Test_InsertStatRequest_insertStatRequestCPUInfoError() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCpuInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("Test_InsertStatRequest_insertStatRequestCPUInfoError"))

	err := clickStatRequest.InsertStatRequest(&apiPb.Metric{
		CpuInfo: &apiPb.CpuInfo{
			Cpus: []*apiPb.CpuInfo_CPU{{}},
		},
		MemoryInfo: &apiPb.MemoryInfo{
			Mem:  &apiPb.MemoryInfo_Memory{},
			Swap: &apiPb.MemoryInfo_Memory{},
		},
		DiskInfo: &apiPb.DiskInfo{
			Disks: map[string]*apiPb.DiskInfo_Disk{
				"": {},
			},
		},
		NetInfo: &apiPb.NetInfo{
			Interfaces: map[string]*apiPb.NetInfo_Interface{
				"": {},
			},
		},
		Time: timestamp.Now(),
	})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_InsertStatRequest_insertStatRequestsMemoryInfoMemError() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCpuInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestMemoryInfoMemCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("Test_InsertStatRequest_insertStatRequestsMemoryInfoMemError"))

	err := clickStatRequest.InsertStatRequest(&apiPb.Metric{
		CpuInfo: &apiPb.CpuInfo{
			Cpus: []*apiPb.CpuInfo_CPU{{}},
		},
		MemoryInfo: &apiPb.MemoryInfo{
			Mem:  &apiPb.MemoryInfo_Memory{},
			Swap: &apiPb.MemoryInfo_Memory{},
		},
		DiskInfo: &apiPb.DiskInfo{
			Disks: map[string]*apiPb.DiskInfo_Disk{
				"": {},
			},
		},
		NetInfo: &apiPb.NetInfo{
			Interfaces: map[string]*apiPb.NetInfo_Interface{
				"": {},
			},
		},
		Time: timestamp.Now(),
	})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_InsertStatRequest_insertStatRequestsMemoryInfoSwapError() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCpuInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestMemoryInfoMemCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestMemoryInfoSwapCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("Test_InsertStatRequest_insertStatRequestsMemoryInfoMemError"))

	err := clickStatRequest.InsertStatRequest(&apiPb.Metric{
		CpuInfo: &apiPb.CpuInfo{
			Cpus: []*apiPb.CpuInfo_CPU{{}},
		},
		MemoryInfo: &apiPb.MemoryInfo{
			Mem:  &apiPb.MemoryInfo_Memory{},
			Swap: &apiPb.MemoryInfo_Memory{},
		},
		DiskInfo: &apiPb.DiskInfo{
			Disks: map[string]*apiPb.DiskInfo_Disk{
				"": {},
			},
		},
		NetInfo: &apiPb.NetInfo{
			Interfaces: map[string]*apiPb.NetInfo_Interface{
				"": {},
			},
		},
		Time: timestamp.Now(),
	})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_InsertStatRequest_insertStatRequestsDiskInfoError() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCpuInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestMemoryInfoMemCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestMemoryInfoSwapCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestDiskInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("Test_InsertStatRequest_insertStatRequestsDiskInfoError"))

	err := clickStatRequest.InsertStatRequest(&apiPb.Metric{
		CpuInfo: &apiPb.CpuInfo{
			Cpus: []*apiPb.CpuInfo_CPU{{}},
		},
		MemoryInfo: &apiPb.MemoryInfo{
			Mem:  &apiPb.MemoryInfo_Memory{},
			Swap: &apiPb.MemoryInfo_Memory{},
		},
		DiskInfo: &apiPb.DiskInfo{
			Disks: map[string]*apiPb.DiskInfo_Disk{
				"": {},
			},
		},
		NetInfo: &apiPb.NetInfo{
			Interfaces: map[string]*apiPb.NetInfo_Interface{
				"": {},
			},
		},
		Time: timestamp.Now(),
	})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_InsertStatRequest_insertStatRequestsNetInfoError() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCpuInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestMemoryInfoMemCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestMemoryInfoSwapCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestDiskInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestNetInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("Test_InsertStatRequest_insertStatRequestsNetInfoError"))

	err := clickStatRequest.InsertStatRequest(&apiPb.Metric{
		CpuInfo: &apiPb.CpuInfo{
			Cpus: []*apiPb.CpuInfo_CPU{{}},
		},
		MemoryInfo: &apiPb.MemoryInfo{
			Mem:  &apiPb.MemoryInfo_Memory{},
			Swap: &apiPb.MemoryInfo_Memory{},
		},
		DiskInfo: &apiPb.DiskInfo{
			Disks: map[string]*apiPb.DiskInfo_Disk{
				"": {},
			},
		},
		NetInfo: &apiPb.NetInfo{
			Interfaces: map[string]*apiPb.NetInfo_Interface{
				"": {},
			},
		},
		Time: timestamp.Now(),
	})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_insertStatRequest_error() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("Test_insertStatRequest_Error"))

	srId := clickhouse.UUID(uuid.New().String())
	err := clickStatRequest.insertStatRequest(time.Now(), srId, &StatRequest{})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_insertStatRequest_commitError() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit().WillReturnError(errors.New("Test_insertStatRequest_commitError"))

	srId := clickhouse.UUID(uuid.New().String())
	err := clickStatRequest.insertStatRequest(time.Now(), srId, &StatRequest{})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_insertStatRequestCPUInfo_error() {
	srId := clickhouse.UUID(uuid.New().String())
	err := clickStatRequest.insertStatRequestCPUInfo(time.Now(), srId, &CPUInfo{})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_insertStatRequestCPUInfo_commitError() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCpuInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit().WillReturnError(errors.New("Test_insertStatRequestCPUInfo_commitError"))

	srId := clickhouse.UUID(uuid.New().String())
	err := clickStatRequest.insertStatRequestCPUInfo(time.Now(), srId, &CPUInfo{})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_insertStatRequestsMemoryInfoMem_error() {
	srId := clickhouse.UUID(uuid.New().String())
	err := clickStatRequest.insertStatRequestsMemoryInfoMem(time.Now(), srId, &MemoryInfo{
		Model:         Model{},
		StatRequestID: "",
		Mem:           &MemoryMem{},
		Swap:          &MemorySwap{},
	})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_insertStatRequestsMemoryInfoMem_commitError() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestMemoryInfoMemCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit().WillReturnError(errors.New("Test_insertStatRequestsMemoryInfoMem_commitError"))

	srId := clickhouse.UUID(uuid.New().String())
	err := clickStatRequest.insertStatRequestsMemoryInfoMem(time.Now(), srId, &MemoryInfo{
		Model:         Model{},
		StatRequestID: "",
		Mem:           &MemoryMem{},
		Swap:          &MemorySwap{},
	})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_insertStatRequestsMemoryInfoSwap_error() {
	srId := clickhouse.UUID(uuid.New().String())
	err := clickStatRequest.insertStatRequestsMemoryInfoSwap(time.Now(), srId, &MemoryInfo{
		Model:         Model{},
		StatRequestID: "",
		Mem:           &MemoryMem{},
		Swap:          &MemorySwap{},
	})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_insertStatRequestsMemoryInfoSwap_commitError() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestMemoryInfoSwapCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit().WillReturnError(errors.New("Test_insertStatRequestsMemoryInfoSwap_commitError"))

	srId := clickhouse.UUID(uuid.New().String())
	err := clickStatRequest.insertStatRequestsMemoryInfoSwap(time.Now(), srId, &MemoryInfo{
		Model:         Model{},
		StatRequestID: "",
		Mem:           &MemoryMem{},
		Swap:          &MemorySwap{},
	})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_insertStatRequestsDiskInfo_error() {
	srId := clickhouse.UUID(uuid.New().String())
	err := clickStatRequest.insertStatRequestsDiskInfo(time.Now(), srId, &DiskInfo{})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_insertStatRequestsDiskInfo_commitError() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestDiskInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit().WillReturnError(errors.New("Test_insertStatRequestsDiskInfo_commitError"))

	srId := clickhouse.UUID(uuid.New().String())
	err := clickStatRequest.insertStatRequestsDiskInfo(time.Now(), srId, &DiskInfo{})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_insertStatRequestsNetInfo_error() {
	srId := clickhouse.UUID(uuid.New().String())
	err := clickStatRequest.insertStatRequestsNetInfo(time.Now(), srId, &NetInfo{
		Model:         Model{},
		StatRequestID: "",
		Name:          "",
		BytesSent:     0,
		BytesRecv:     0,
		PacketsSent:   0,
		PacketsRecv:   0,
		ErrIn:         0,
		ErrOut:        0,
		DropIn:        0,
		DropOut:       0,
	})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_insertStatRequestsNetInfo_commitError() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestNetInfoCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit().WillReturnError(errors.New("Test_insertStatRequestsDiskInfo_commitError"))

	srId := clickhouse.UUID(uuid.New().String())
	err := clickStatRequest.insertStatRequestsNetInfo(time.Now(), srId, &NetInfo{
		Model:         Model{},
		StatRequestID: "",
		Name:          "",
		BytesSent:     0,
		BytesRecv:     0,
		PacketsSent:   0,
		PacketsRecv:   0,
		ErrIn:         0,
		ErrOut:        0,
		DropIn:        0,
		DropOut:       0,
	})
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_GetStatRequest() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsCpuInfoFields, dbStatRequestCpuInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "load"}).
		AddRow("1", time.Now(), "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoMemCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoSwapCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsDiskInfoFields, dbStatRequestDiskInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "name", "total", "used", "free", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsNetInfoFields, dbStatRequestNetInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "name", "bytes_sent", "bytes_recv", "packets_sent", "packets_recv", "err_in", "err_out", "drop_in", "drop_out"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetStatRequest(id, nil, nil)
	require.NoError(s.T(), err)
}

func (s *SuiteStatRequest) Test_GetStatRequest_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetStatRequest(id, &apiPb.Pagination{
		Page:  1, //random value
		Limit: 2, //random value
	}, nil)
	require.Error(s.T(), err)
}

func TestClickhouse_GetStatRequest(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := clickhWrongStatRequest.GetStatRequest("", nil, &apiPb.TimeFilter{
			From: &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := clickhWrongStatRequest.GetStatRequest("", nil, nil)
		assert.Error(t, err)
	})
}

func (s *SuiteStatRequest) Test_GetStatRequest_StatRequestScanError() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time", "a"}).
		AddRow("1", time.Now(), "1", "1", time.Now(), "a")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.getStatRequests(id, nil, nil)
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_GetStatRequest_getStatRequestsCpuInfoError() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsCpuInfoFields, dbStatRequestCpuInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "load"}).
		AddRow("1", time.Now(), "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New("Test_GetStatRequest_getStatRequestsCpuInfoError"))

	_, _, err := clickStatRequest.getStatRequests(id, nil, nil)
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_GetStatRequest_getStatRequestsMemoryInfoMem() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsCpuInfoFields, dbStatRequestCpuInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "load"}).
		AddRow("1", time.Now(), "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoMemCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New("Test_GetStatRequest_getStatRequestsMemoryInfoMem"))

	_, _, err := clickStatRequest.getStatRequests(id, nil, nil)
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_GetStatRequest_getStatRequestsMemoryInfoSwap() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsCpuInfoFields, dbStatRequestCpuInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "load"}).
		AddRow("1", time.Now(), "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoMemCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoSwapCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnError(errors.New("Test_GetStatRequest_getStatRequestsMemoryInfoSwap"))

	_, _, err := clickStatRequest.getStatRequests(id, nil, nil)
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_GetStatRequest_getStatRequestsDiskInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsCpuInfoFields, dbStatRequestCpuInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "load"}).
		AddRow("1", time.Now(), "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoMemCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoSwapCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsDiskInfoFields, dbStatRequestDiskInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "name", "total", "used", "free", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnError(errors.New("Test_GetStatRequest_getStatRequestsDiskInfo"))

	_, _, err := clickStatRequest.getStatRequests(id, nil, nil)
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_GetStatRequest_getStatRequestsNetInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsCpuInfoFields, dbStatRequestCpuInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "load"}).
		AddRow("1", time.Now(), "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoMemCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoSwapCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsDiskInfoFields, dbStatRequestDiskInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "name", "total", "used", "free", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsNetInfoFields, dbStatRequestNetInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "name", "bytes_sent", "bytes_recv", "packets_sent", "packets_recv", "err_in", "err_out", "drop_in", "drop_out"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New("Test_GetStatRequest_getStatRequestsNetInfo"))

	_, _, err := clickStatRequest.getStatRequests(id, nil, nil)
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_countStatRequests_nextError() {
	var (
		id = "1"
	)
	_, err := clickStatRequest.countStatRequests(id, time.Now(), time.Now())
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_countStatRequests_selectError() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count", "a"}).AddRow("1", "a")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := clickStatRequest.countStatRequests(id, time.Now(), time.Now())
	require.Error(s.T(), err)
}

func (s *SuiteStatRequest) Test_GetCpuInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsCpuInfoFields, dbStatRequestCpuInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "load"}).
		AddRow("1", time.Now(), "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetCPUInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

//Is used for getSpecialRecords test
func (s *SuiteStatRequest) Test_GetCpuInfo_Count_Error() {
	var (
		id = "1"
	)

	_, _, err := clickStatRequest.GetCPUInfo(id, nil, nil)
	require.Error(s.T(), err)
}

//Is used for getSpecialRecords test
func (s *SuiteStatRequest) Test_GetCpuInfo_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetCPUInfo(id, &apiPb.Pagination{
		Page:  1, //random value
		Limit: 2, //random value
	}, nil)
	require.Error(s.T(), err)
}

//Time errors in getSpecialRecords
func TestClickhouse_GetCpuInfo(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := clickhWrongStatRequest.GetCPUInfo("", nil, &apiPb.TimeFilter{
			From: &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		})
		assert.Error(t, err)
	})
}

func (s *SuiteStatRequest) Test_GetMemoryInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoMemCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoSwapCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "memory_info_id", "total", "used", "free", "shared", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetMemoryInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

func (s *SuiteStatRequest) Test_GetMemoryInfo_Select_Error() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetMemoryInfo(id, &apiPb.Pagination{
		Page:  1, //random value
		Limit: 2, //random value
	}, nil)
	require.Error(s.T(), err)
}

func TestClickhouse_GetMemoryInfo(t *testing.T) {
	//Time for invalid timestamp
	maxValidSeconds := 253402300800
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := clickhWrongStatRequest.GetMemoryInfo("", nil, &apiPb.TimeFilter{
			From: &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
			To:   &timestamp.Timestamp{Seconds: int64(maxValidSeconds), Nanos: 0},
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, _, err := clickhWrongStatRequest.GetMemoryInfo("", nil, nil)
		assert.Error(t, err)
	})
}

func (s *SuiteStatRequest) Test_GetDiskInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsDiskInfoFields, dbStatRequestDiskInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "name", "total", "used", "free", "used_percent"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetDiskInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

func (s *SuiteStatRequest) Test_GetNetInfo() {
	var (
		id = "1"
	)

	query := fmt.Sprintf(`SELECT count(*) FROM "%s"`, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestFields, dbStatRequestCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "agent_id", "agent_name", "time"}).
		AddRow("1", time.Now(), "1", "1", time.Now())
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	query = fmt.Sprintf(`SELECT %s FROM "%s"`, statRequestsNetInfoFields, dbStatRequestNetInfoCollection)
	rows = sqlmock.NewRows([]string{"id", "created_at", "stat_request_id", "name", "bytes_sent", "bytes_recv", "packets_sent", "packets_recv", "err_in", "err_out", "drop_in", "drop_out"}).
		AddRow("1", time.Now(), "1", "1", "1", "1", "1", "1", "1", "1", "1", "1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, _, err := clickStatRequest.GetNetInfo(id, nil, nil)
	require.NoError(s.T(), err)
}

func (s *SuiteStatRequest) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInitStatRequest(t *testing.T) {
	suite.Run(t, new(SuiteStatRequest))
}
