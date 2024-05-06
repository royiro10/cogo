package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/royiro10/cogo/cmd"
	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/server"
	"github.com/royiro10/cogo/services"
)

type CommandArgs string

func main() {
	flags := cmd.InitFlags()
	args := flag.Args()

	logger := common.EmptyLogger
	if flags.IsLogging {
		logger = common.CreateLogger("./logs", fmt.Sprintf("cogo_%d.log", os.Getpid()))
	}

	commandService := services.CreateCommandService(logger)
	lockService := services.CreateLockFileService(logger)
	daemon := server.CreateCogoDaemon(logger, commandService)

	cli := cmd.CreateCLI(cmd.CogoCLIDeps{
		Logger:      logger,
		LockService: lockService,
		CogoDaemon:  daemon,
	})

	cli.Handle(args, flags)
}
