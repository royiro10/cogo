package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/royiro10/cogo/client"
	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/models"
	"github.com/royiro10/cogo/services/commands"
)

func makeRunCommand(lockService common.LockService, logger *common.Logger) models.CogoCLICommand {
	return func(cmdInfo models.CogoCLIInfo) error {
		if !lockService.IsAcquired(GetLockFile()) {
			// TODO: automaticcly start daemon
			return fmt.Errorf("cogo must be start before running commands")
		}
		workdir, err := os.Getwd()
		if err != nil {
			logger.Fatal(err)
		}
		client := client.CreateCogoClient(logger, workdir)
		defer client.Close()

		session := cmdInfo.Flags.Session
		if session == "" {
			session = commands.DefaultSessionKey
		}

		return client.Run(models.NewExecuteRequest(session, strings.Join(cmdInfo.Args[:], " ")))
	}
}
