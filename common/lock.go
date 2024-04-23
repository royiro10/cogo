package common

import (
	"github.com/royiro10/cogo/models"
)

type LockService interface {
	Acquire(lockName string) (IDisposable, error)
	Release(lockName string) error
	IsAcquired(lockName string) bool
	GetLockCommit(lockName string) (*models.LockCommit, error)
}
