package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/royiro10/cogo/services"
	"github.com/royiro10/cogo/util"
)

type CommandArgs string

const (
	RUN_DAEMON  CommandArgs = "daemon"
	STOP_DAEMON CommandArgs = "stop"
)

func main() {
	isLogging := flag.Bool("logger", false, "shold log")

	flag.Parse()
	args := flag.Args()

	var logger = util.EmptyLogger
	if *isLogging {
		logger = util.CreateLogger(fmt.Sprintf("./logs/cogo_%d.log", os.Getpid()))
	}

	commandService := services.CreateCommandService(logger)
	lockService := services.CreateLockFileService(logger)
	daemon := services.CreateCogoDaemon(logger, commandService)

	cli := services.CreateCLI(services.CogoCLIDeps{
		Logger:      logger,
		LockService: lockService,
		CogoDaemon:  daemon,
	})

	cli.Handle(args)
}
