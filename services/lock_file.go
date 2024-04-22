package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/royiro10/cogo/models"
	"github.com/royiro10/cogo/util"
)

type LockFileService struct {
	logger    *util.Logger
	isAquired bool
}

func CreateLockFileService(logger *util.Logger) LockService {
	l := logger
	if l == nil {
		l = util.DefaultLogger
	}

	service := &LockFileService{
		logger:    l,
		isAquired: false,
	}

	return service
}

func (s *LockFileService) Acquire(lockName string) (util.IDisposable, error) {
	if err := s.aquireLockFile(lockName); err != nil {
		return func() {}, err
	}

	return func() { s.releaseLockFile(lockName) }, nil
}

func (s *LockFileService) Release(lockName string) error {
	return s.releaseLockFile(lockName)
}

func (s *LockFileService) GetLockCommit(lockName string) (*models.LockCommit, error) {
	f, err := os.OpenFile(lockName, os.O_RDONLY, 0600)
	if err != nil {
		s.logger.Error("could not find daemon process", err)
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		s.logger.Error("could not read process of running daemon")
		return nil, err
	}

	var lockfileCommit models.LockCommit
	if err = json.Unmarshal(data, &lockfileCommit); err != nil {
		s.logger.Error("failed to parse LockFileCommit", "data", data)
		return nil, err
	}

	return &lockfileCommit, nil
}

func (s *LockFileService) aquireLockFile(lockFile string) error {
	if s.IsAquired(lockFile) {
		return fmt.Errorf("can not aquire, already locked")
	}

	f, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	lockCommit := &models.LockCommit{
		Time: time.Now(),
		Pid:  os.Getpid(),
	}

	jsonData, err := json.Marshal(lockCommit)
	if err != nil {
		s.logger.Error("could not marshel lockCommit", "err", err)
		return err
	}

	if _, err := f.Write(jsonData); err != nil {
		s.logger.Error("could write commit lock to file", "err", err, "file", lockFile)
		return err
	}

	return nil
}

func (s *LockFileService) releaseLockFile(lockFile string) error {
	if !s.IsAquired(lockFile) {
		return nil
	}

	lockCommit, err := s.GetLockCommit(lockFile)
	if err != nil {
		return fmt.Errorf("could not get the lock commit from lock file", "err", err)
	}

	err = os.Remove(lockFile)
	if err != nil {
		return fmt.Errorf("could not release lock file", err)
	}

	s.logger.Info("released lock file", "previus", lockCommit)
	return nil
}

func (s *LockFileService) IsAquired(lockFile string) bool {
	_, err := os.Stat(lockFile)
	return !os.IsNotExist(err)
}
