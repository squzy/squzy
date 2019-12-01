package storage

import (
	"squzy/apps/internal/job"
)

type Storage interface {
	Write(id string, log job.CheckError) error
}

type memory struct {
}

func (m memory) Write(id string, log job.CheckError) error {
	return nil
}

func GetInMemoryStorage() Storage {
	return &memory{}
}
