package main

import (
	"github.com/royiro10/cogo/services"
	"github.com/royiro10/cogo/util"
)

var logger = util.DefaultLogger

type CommandArgs string

const (
	RUN_DAEMON  CommandArgs = "daemon"
	STOP_DAEMON CommandArgs = "stop"
)

func main() {
	commandService := services.CreateCommandService(logger)
	lockService := services.CreateLockFileService(logger)
	daemon := services.CreateCogoDaemon(logger, commandService)

	cli := services.CreateCLI(services.CogoCLIDeps{
		Logger:      logger,
		LockService: lockService,
		CogoDaemon:  daemon,
	})

	cli.Handle()
}
