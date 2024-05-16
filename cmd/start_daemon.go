package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/models"
)

func makeHandleStartDaemon(
	lockService common.LockService,
	logger *common.Logger,
) models.CogoCLICommand {
	return func(cci models.CogoCLIInfo) error {
		_, err := startDaemon(lockService, logger)
		return err
	}
}

func startDaemon(lockService common.LockService, logger *common.Logger) (bool, error) {
	if lockService.IsAcquired(GetLockFile()) {
		logger.Info("Daemon is already running.")
		return false, nil
	}

	logger.Info("Starting daemon...", "arg0", os.Args[0])

	cmd := exec.Command(os.Args[0], "--logger", RUN_DAEMON)
	err := cmd.Start()
	if err != nil {
		return false, fmt.Errorf("cmd.Start failed: %w ", err)
	}

	err = cmd.Process.Release()
	if err != nil {
		return false, fmt.Errorf("cmd.Process.Release failed: %w ", err)
	}

	return true, nil
}
