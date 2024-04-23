package cmd

import (
	"fmt"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/models"
	"github.com/royiro10/cogo/util"
)

func makeHandleStopDaemon(lockService common.LockService, logger *common.Logger) models.CogoCLICommand {
	handleStopDaemon := func(cmdInfo models.CogoCLIInfo) error {
		lockCommit, err := lockService.GetLockCommit(LOCK_FILE)
		if err != nil {
			return fmt.Errorf("can not get commited lock: %w", err)
		}

		if err = util.SendInterrupt(lockCommit.Pid); err != nil {
			logger.Error("could not send interrupt to process", "err", err, "pid", lockCommit.Pid)
			logger.Warn("hard kill", "pid", lockCommit.Pid)
			if err := util.KillCmd(lockCommit.Pid).Start(); err != nil {
				return fmt.Errorf("hard kill failed: %w", err)
			}
		}

		return lockService.Release(LOCK_FILE)
	}

	return handleStopDaemon
}
