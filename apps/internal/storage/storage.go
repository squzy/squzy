package storage

import (
	logger "log"
	"squzy/apps/internal/job"
)

type Storage interface {
	Write(id string, log job.CheckError) error
}

type memory struct {
}

func (m memory) Write(id string, log job.CheckError) error {
	logger.Println(id, log.GetLogData())
	return nil
}

func GetInMemoryStorage() Storage {
	return &memory{}
}
