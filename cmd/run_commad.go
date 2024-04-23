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
		if !lockService.IsAquired(LOCK_FILE) {
			// TODO: automaticcly start daemon
			return fmt.Errorf("cogo must be start before running commands")
		}

		client := client.CreateCogoClient(logger)
		defer client.Close()

		commandRequest := &services.CommandParameters{
			SessionId: services.DefaultSessionKey,
			Command:   strings.Join(cmdInfo.Args[:], " "),
		}

		client.Run(commandRequest)
		return nil
	}
}
