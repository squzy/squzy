package semaphore

import (
	"context"
	"golang.org/x/sync/semaphore"
)

type SemaphoreFactory func (int) Semaphore

type Semaphore interface {
	Acquire(ctx context.Context) error
	Release()
}

type sem struct {
	s *semaphore.Weighted
}

func NewSemaphore(weight int) Semaphore {
	return &sem{
		s: semaphore.NewWeighted(int64(weight)),
	}
}

func (s *sem) Acquire(ctx context.Context) error {
	return s.s.Acquire(ctx, 1)
}

func (s *sem) Release() {
	s.s.Release(1)
}
