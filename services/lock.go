package services

import (
	"github.com/royiro10/cogo/models"
	"github.com/royiro10/cogo/util"
)

type LockService interface {
	Acquire(lockName string) (util.IDisposable, error)
	Release(lockName string) error
	IsAquired(lockName string) bool
	GetLockCommit(lockName string) (*models.LockCommit, error)
}
