package main

import (
	host_metric "squzy/apps/internal/host-metric"
	"squzy/apps/internal/job"
	"squzy/apps/internal/scheduler"
	"squzy/apps/internal/storage"
	"time"
	clientPb "github.com/squzy/squzy_generated/generated/storage/proto/v1"
)

type mockJob struct {
	metric host_metric.Metric
}

type mock struct {

}

func (m mock) GetLogData() *clientPb.Log {
	return &clientPb.Log{
		Code: clientPb.StatusCode_OK,
	}
}

func (m *mockJob) Do() job.CheckError {
	m.metric.GetStat()
	return &mock{}
}

func main() {
	schl, _ := scheduler.New(time.Second, &mockJob{
		metric: host_metric.NewCpuMetric(time.Second),
	}, storage.GetInMemoryStorage())
	_ = schl.Run()
}
