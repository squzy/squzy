package clickhouse

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	uuid "github.com/google/uuid"
	"github.com/squzy/squzy/internal/logger"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	"time"
)

type StatRequest struct {
	Model      Model
	AgentID    string
	AgentName  string
	CPUInfo    []*CPUInfo
	MemoryInfo *MemoryInfo
	DiskInfo   []*DiskInfo
	NetInfo    []*NetInfo
	Time       time.Time
}

const (
	cpuInfoKey  = "CPUInfo"
	diskInfoKey = "DiskInfo"
	netInfoKey  = "NetInfo"
)

type CPUInfo struct {
	Model         Model
	StatRequestID string
	Load          float64
}

type MemoryInfo struct {
	Model         Model
	StatRequestID string
	Mem           *MemoryMem
	Swap          *MemorySwap
}

type MemoryMem struct {
	Model         Model
	StatRequestID string
	MemoryInfoID  string
	Total         uint64
	Used          uint64
	Free          uint64
	Shared        uint64
	UsedPercent   float64
}

type MemorySwap struct {
	Model         Model
	StatRequestID string
	MemoryInfoID  string
	Total         uint64
	Used          uint64
	Free          uint64
	Shared        uint64
	UsedPercent   float64
}

type DiskInfo struct {
	Model         Model
	StatRequestID string
	Name          string
	Total         uint64
	Free          uint64
	Used          uint64
	UsedPercent   float64
}

type NetInfo struct {
	Model         Model
	StatRequestID string
	Name          string
	BytesSent     uint64
	BytesRecv     uint64
	PacketsSent   uint64
	PacketsRecv   uint64
	ErrIn         uint64
	ErrOut        uint64
	DropIn        uint64
	DropOut       uint64
}

var (
	statRequestFields            = "id, created_at, agent_id, agent_name, time"
	statRequestsCpuInfoFields    = "id, created_at, stat_request_id, load"
	statRequestsMemoryInfoFields = "id, created_at, stat_request_id, memory_info_id, total, used, free, shared, used_percent"
	statRequestsDiskInfoFields   = "id, created_at, stat_request_id, name, total, used, free, used_percent"
	statRequestsNetInfoFields    = "id, created_at, stat_request_id, name, bytes_sent, bytes_recv, packets_sent, packets_recv, err_in, err_out, drop_in, drop_out"
	statRequestIdFilterString    = fmt.Sprintf(`"stat_request_id" = ?`)
	agentIdFilterString          = fmt.Sprintf(`"agent_id" = ?`)
	statRequestTimeFilterString  = fmt.Sprintf(`"time" BETWEEN ? and ?`)
	statRequestTimeString        = `"time"`
)

func (c *Clickhouse) InsertStatRequest(data *apiPb.Metric) error {
	now := time.Now()
	srId := clickhouse.UUID(uuid.New().String())
	srData, err := ConvertToClickhouseStatRequest(data)
	if err != nil {
		return err
	}
	err = c.insertStatRequest(now, srId, srData)
	if err != nil {
		logger.Error(err.Error())
		return errorDataBase
	}

	for _, cpuInfo := range srData.CPUInfo {
		err = c.insertStatRequestCPUInfo(now, srId, cpuInfo)
		if err != nil {
			logger.Error(err.Error())
			return errorDataBase
		}
	}

	err = c.insertStatRequestsMemoryInfoMem(now, srId, srData.MemoryInfo)
	if err != nil {
		logger.Error(err.Error())
		return errorDataBase
	}

	err = c.insertStatRequestsMemoryInfoSwap(now, srId, srData.MemoryInfo)
	if err != nil {
		logger.Error(err.Error())
		return errorDataBase
	}

	for _, diskInfo := range srData.DiskInfo {
		err = c.insertStatRequestsDiskInfo(now, srId, diskInfo)
		if err != nil {
			logger.Error(err.Error())
			return errorDataBase
		}
	}

	for _, netInfo := range srData.NetInfo {
		err = c.insertStatRequestsNetInfo(now, srId, netInfo)
		if err != nil {
			logger.Error(err.Error())
			return errorDataBase
		}
	}

	return nil
}

func (c *Clickhouse) insertStatRequest(now time.Time, sr_id clickhouse.UUID, statRequest *StatRequest) error {
	tx, err := c.Db.Begin()
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES ($0, $1, $2, $3, $4)`, dbStatRequestCollection, statRequestFields)
	_, err = tx.Exec(q,
		sr_id,
		now,
		statRequest.AgentID,
		statRequest.AgentName,
		statRequest.Time,
	)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (c *Clickhouse) insertStatRequestCPUInfo(now time.Time, sr_id clickhouse.UUID, cpuInfo *CPUInfo) error {
	tx, err := c.Db.Begin()
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES ($0, $1, $2, $3)`, dbStatRequestCpuInfoCollection, statRequestsCpuInfoFields)
	_, err = tx.Exec(q,
		clickhouse.UUID(uuid.New().String()),
		now,
		sr_id,
		cpuInfo.Load,
	)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (c *Clickhouse) insertStatRequestsMemoryInfoMem(now time.Time, sr_id clickhouse.UUID, memoryInfo *MemoryInfo) error {
	tx, err := c.Db.Begin()
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES ($0, $1, $2, $3, $4, $5, $6, $7)`, dbStatRequestMemoryInfoMemCollection, statRequestsMemoryInfoFields)
	_, err = tx.Exec(q,
		clickhouse.UUID(uuid.New().String()),
		now,
		sr_id,
		memoryInfo.Mem.MemoryInfoID,
		memoryInfo.Mem.Total,
		memoryInfo.Mem.Used,
		memoryInfo.Mem.Free,
		memoryInfo.Mem.Shared,
		memoryInfo.Mem.UsedPercent,
	)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (c *Clickhouse) insertStatRequestsMemoryInfoSwap(now time.Time, sr_id clickhouse.UUID, memoryInfo *MemoryInfo) error {
	tx, err := c.Db.Begin()
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES ($0, $1, $2, $3, $4, $5, $6, $7)`, dbStatRequestMemoryInfoSwapCollection, statRequestsMemoryInfoFields)
	_, err = tx.Exec(q,
		clickhouse.UUID(uuid.New().String()),
		now,
		sr_id,
		memoryInfo.Swap.MemoryInfoID,
		memoryInfo.Swap.Total,
		memoryInfo.Swap.Used,
		memoryInfo.Swap.Free,
		memoryInfo.Swap.Shared,
		memoryInfo.Swap.UsedPercent,
	)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (c *Clickhouse) insertStatRequestsDiskInfo(now time.Time, sr_id clickhouse.UUID, diskInfo *DiskInfo) error {
	tx, err := c.Db.Begin()
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES ($0, $1, $2, $3, $4, $5, $6, $7)`, dbStatRequestDiskInfoCollection, statRequestsDiskInfoFields)
	_, err = tx.Exec(q,
		clickhouse.UUID(uuid.New().String()),
		now,
		sr_id,
		diskInfo.Name,
		diskInfo.Total,
		diskInfo.Used,
		diskInfo.Free,
		diskInfo.UsedPercent,
	)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (c *Clickhouse) insertStatRequestsNetInfo(now time.Time, sr_id clickhouse.UUID, netInfo *NetInfo) error {
	tx, err := c.Db.Begin()
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES ($0, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, dbStatRequestNetInfoCollection, statRequestsNetInfoFields)
	_, err = tx.Exec(q,
		clickhouse.UUID(uuid.New().String()),
		now,
		sr_id,
		netInfo.Name,
		netInfo.BytesSent,
		netInfo.BytesRecv,
		netInfo.PacketsSent,
		netInfo.PacketsRecv,
		netInfo.ErrIn,
		netInfo.ErrOut,
		netInfo.DropIn,
		netInfo.DropOut,
	)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (c *Clickhouse) GetStatRequest(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	srs, count, err := c.getStatRequests(agentID, pagination, filter)
	if err != nil {
		return nil, -1, err
	}

	return ConvertFromClickhouseStatRequests(srs), count, nil
}

func (c *Clickhouse) getStatRequests(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*StatRequest, int32, error) {
	rows, count, err := c.getStatRequestsRows(agentID, pagination, filter)
	if err != nil {
		return nil, -1, err
	}

	defer rows.Close()

	var statRequests []*StatRequest
	for rows.Next() {
		sr := &StatRequest{}
		if err := rows.Scan(&sr.Model.ID, &sr.Model.CreatedAt,
			&sr.AgentID, &sr.AgentName, &sr.Time); err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}

		sr.CPUInfo, err = c.getStatRequestsCpuInfo(sr.Model.ID)
		if err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}
		mem, err := c.getStatRequestsMemoryInfoMem(sr.Model.ID)
		if err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}
		swp, err := c.getStatRequestsMemoryInfoSwap(sr.Model.ID)
		if err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}

		sr.MemoryInfo = &MemoryInfo{
			Model:         sr.Model,
			StatRequestID: sr.Model.ID,
			Mem:           mem,
			Swap:          swp,
		}

		sr.DiskInfo, err = c.getStatRequestsDiskInfo(sr.Model.ID)
		if err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}
		sr.NetInfo, err = c.getStatRequestsNetInfo(sr.Model.ID)
		if err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}

		statRequests = append(statRequests, sr)
	}

	return statRequests, count, nil
}

func (c *Clickhouse) getStatRequestsRows(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) (*sql.Rows, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	count, err := c.countStatRequests(agentID, timeFrom, timeTo)
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(int64(count), pagination)

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM "%s" WHERE (%s AND %s) ORDER BY %s LIMIT %d OFFSET %d`,
		statRequestFields,
		dbStatRequestCollection,
		agentIdFilterString,
		statRequestTimeFilterString,
		statRequestTimeString,
		limit,
		offset),
		agentID,
		timeFrom,
		timeTo,
	)

	if err != nil {
		logger.Error(err.Error())
		return nil, -1, errorDataBase
	}
	return rows, int32(count), nil
}

func (c *Clickhouse) countStatRequests(agentID string, timeFrom time.Time, timeTo time.Time) (int32, error) {
	var count int64
	rows, err := c.Db.Query(fmt.Sprintf(`SELECT count(*) FROM "%s" WHERE %s AND %s`,
		dbStatRequestCollection,
		agentIdFilterString,
		statRequestTimeFilterString),
		agentID,
		timeFrom,
		timeTo)

	if err != nil {
		logger.Error(err.Error())
		return -1, errorDataBase
	}

	defer rows.Close()

	if ok := rows.Next(); !ok {
		return 0, nil
	}

	if err := rows.Scan(&count); err != nil {
		logger.Error(err.Error())
		return -1, errorDataBase
	}

	return int32(count), nil
}
func (c *Clickhouse) getStatRequestsCpuInfo(id string) ([]*CPUInfo, error) {
	var cpuInfos []*CPUInfo

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM "%s" WHERE %s`, statRequestsCpuInfoFields, dbStatRequestCpuInfoCollection, statRequestIdFilterString), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cpu := &CPUInfo{}
		if err := rows.Scan(&cpu.Model.ID, &cpu.Model.CreatedAt, &cpu.StatRequestID, &cpu.Load); err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		cpuInfos = append(cpuInfos, cpu)
	}

	return cpuInfos, nil
}

func (c *Clickhouse) getStatRequestsMemoryInfoMem(id string) (*MemoryMem, error) {
	mem := &MemoryMem{}

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM "%s" WHERE %s`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoMemCollection, statRequestIdFilterString), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&mem.Model.ID, &mem.Model.CreatedAt, &mem.StatRequestID, &mem.MemoryInfoID, &mem.Total, &mem.Used, &mem.Free, &mem.Shared, &mem.UsedPercent); err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	}

	return mem, nil
}

func (c *Clickhouse) getStatRequestsMemoryInfoSwap(id string) (*MemorySwap, error) {
	mem := &MemorySwap{}

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM "%s" WHERE %s`, statRequestsMemoryInfoFields, dbStatRequestMemoryInfoSwapCollection, statRequestIdFilterString), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&mem.Model.ID, &mem.Model.CreatedAt, &mem.StatRequestID, &mem.MemoryInfoID, &mem.Total, &mem.Used, &mem.Free, &mem.Shared, &mem.UsedPercent); err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	}

	return mem, nil
}

func (c *Clickhouse) getStatRequestsDiskInfo(id string) ([]*DiskInfo, error) {
	var diskInfos []*DiskInfo

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM "%s" WHERE %s`, statRequestsDiskInfoFields, dbStatRequestDiskInfoCollection, statRequestIdFilterString), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		dis := &DiskInfo{}
		if err := rows.Scan(&dis.Model.ID, &dis.Model.CreatedAt, &dis.StatRequestID, &dis.Name, &dis.Total, &dis.Used, &dis.Free, &dis.UsedPercent); err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		diskInfos = append(diskInfos, dis)
	}

	return diskInfos, nil
}

func (c *Clickhouse) getStatRequestsNetInfo(id string) ([]*NetInfo, error) {
	var netInfos []*NetInfo

	rows, err := c.Db.Query(fmt.Sprintf(`SELECT %s FROM "%s" WHERE %s`, statRequestsNetInfoFields, dbStatRequestNetInfoCollection, statRequestIdFilterString), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		net := &NetInfo{}
		if err := rows.Scan(&net.Model.ID, &net.Model.CreatedAt, &net.StatRequestID, &net.Name, &net.BytesSent, &net.BytesRecv, &net.PacketsSent, &net.PacketsRecv, &net.ErrIn, &net.ErrOut, &net.DropIn, &net.DropOut); err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		netInfos = append(netInfos, net)
	}

	return netInfos, nil
}

func (c *Clickhouse) GetCPUInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	rows, count, err := c.getStatRequestsRows(agentID, pagination, filter)
	if err != nil {
		return nil, -1, err
	}

	defer rows.Close()

	var statRequests []*StatRequest
	for rows.Next() {
		sr := &StatRequest{}
		if err := rows.Scan(&sr.Model.ID, &sr.Model.CreatedAt,
			&sr.AgentID, &sr.AgentName, &sr.Time); err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}

		sr.CPUInfo, err = c.getStatRequestsCpuInfo(sr.Model.ID)
		if err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}
		statRequests = append(statRequests, sr)
	}

	return ConvertFromClickhouseStatRequests(statRequests), count, nil
}

func (c *Clickhouse) GetMemoryInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	rows, count, err := c.getStatRequestsRows(agentID, pagination, filter)
	if err != nil {
		return nil, -1, err
	}

	defer rows.Close()

	var statRequests []*StatRequest
	for rows.Next() {
		sr := &StatRequest{}
		if err := rows.Scan(&sr.Model.ID, &sr.Model.CreatedAt,
			&sr.AgentID, &sr.AgentName, &sr.Time); err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}

		mem, err := c.getStatRequestsMemoryInfoMem(sr.Model.ID)
		if err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}
		swp, err := c.getStatRequestsMemoryInfoSwap(sr.Model.ID)
		if err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}

		sr.MemoryInfo = &MemoryInfo{
			Model:         sr.Model,
			StatRequestID: sr.Model.ID,
			Mem:           mem,
			Swap:          swp,
		}

		statRequests = append(statRequests, sr)
	}

	return ConvertFromClickhouseStatRequests(statRequests), count, nil
}

func (c *Clickhouse) GetDiskInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	rows, count, err := c.getStatRequestsRows(agentID, pagination, filter)
	if err != nil {
		return nil, -1, err
	}

	defer rows.Close()

	var statRequests []*StatRequest
	for rows.Next() {
		sr := &StatRequest{}
		if err := rows.Scan(&sr.Model.ID, &sr.Model.CreatedAt,
			&sr.AgentID, &sr.AgentName, &sr.Time); err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}

		sr.DiskInfo, err = c.getStatRequestsDiskInfo(sr.Model.ID)
		if err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}
		statRequests = append(statRequests, sr)
	}

	return ConvertFromClickhouseStatRequests(statRequests), count, nil
}

func (c *Clickhouse) GetNetInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	rows, count, err := c.getStatRequestsRows(agentID, pagination, filter)
	if err != nil {
		return nil, -1, err
	}

	defer rows.Close()

	var statRequests []*StatRequest
	for rows.Next() {
		sr := &StatRequest{}
		if err := rows.Scan(&sr.Model.ID, &sr.Model.CreatedAt,
			&sr.AgentID, &sr.AgentName, &sr.Time); err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}

		sr.NetInfo, err = c.getStatRequestsNetInfo(sr.Model.ID)
		if err != nil {
			logger.Error(err.Error())
			return nil, -1, err
		}
		statRequests = append(statRequests, sr)
	}

	return ConvertFromClickhouseStatRequests(statRequests), count, nil
}
