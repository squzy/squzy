package host_metric

type Metric interface {
	GetStat() interface{}
}