package common

import (
	"github.com/royiro10/cogo/models"
)

type LockService interface {
	Aquire(lockName string) (IDisposable, error)
	Release(lockName string) error
	IsAquired(lockName string) bool
	GetLockCommit(lockName string) (*models.LockCommit, error)
}
