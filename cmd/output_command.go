package cmd

import (
	"fmt"
	"os"

	"github.com/royiro10/cogo/client"
	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/models"
	"github.com/royiro10/cogo/services"
)

func makeOutputCommand(
	lockService common.LockService,
	logger *common.Logger,
) models.CogoCLICommand {
	return func(cmdInfo models.CogoCLIInfo) error {
		if !lockService.IsAcquired(GetLockFile()) {
			return fmt.Errorf("cogo must be start before running commands")
		}
		workdir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		client := client.CreateCogoClient(logger, workdir)
		defer client.Close()

		session := cmdInfo.Flags.Session
		if session == "" {
			session = services.DefaultSessionKey
		}

		return client.Output(models.NewOutputRequest(session, cmdInfo.Flags.IsStream))
	}
}
