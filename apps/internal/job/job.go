package job

type Job interface {
	Do() error
}