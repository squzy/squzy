package storage

import "squzy/apps/internal/job"

type Storage interface {
	Write(id string, log job.CheckError) error
}
