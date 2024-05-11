package cmd

import (
	"fmt"
	"strings"

	"github.com/royiro10/cogo/client"
	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/models"
	"github.com/royiro10/cogo/services"
)

func makeRunCommand(lockService common.LockService, logger *common.Logger) models.CogoCLICommand {
	return func(cmdInfo models.CogoCLIInfo) error {
		if !lockService.IsAcquired(GetLockFile()) {
			// TODO: automaticcly start daemon
			return fmt.Errorf("cogo must be start before running commands")
		}

		client := client.CreateCogoClient(logger)
		defer client.Close()

		session := cmdInfo.Flags.Session
		if session == "" {
			session = services.DefaultSessionKey
		}

		return client.Run(models.NewExecuteRequest(session, strings.Join(cmdInfo.Args[:], " ")))
	}
}
