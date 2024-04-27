package cmd

import (
	"fmt"

	"github.com/royiro10/cogo/client"
	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/models"
	"github.com/royiro10/cogo/services"
)

func makeKillCommand(lockService common.LockService, logger *common.Logger) models.CogoCLICommand {
	return func(cmdInfo models.CogoCLIInfo) error {
		if !lockService.IsAcquired(LOCK_FILE) {
			return fmt.Errorf("cogo must be start before running commands")
		}

		client := client.CreateCogoClient(logger)
		defer client.Close()

		session := cmdInfo.Flags.Session
		if session == "" {
			session = services.DefaultSessionKey
		}

		client.Kill(models.NewKillRequest(session))

		return nil
	}
}
