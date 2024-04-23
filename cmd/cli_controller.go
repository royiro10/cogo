package cmd

import (
	"fmt"
	"strings"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/models"
	"github.com/royiro10/cogo/server"
)

const LOCK_FILE = "./cogo.lock"
const (
	RUN_DAEMON   = "daemon"
	START_DAEMON = "start"
	STOP_DAEMON  = "stop"
	RUN_COMMAND  = "run"
	UNKOWN       = "unkown"
)

type CogoCLIDeps struct {
	Logger      *common.Logger
	LockService common.LockService
	CogoDaemon  server.Daemon
}

type CogoCLI struct {
	logger   *common.Logger
	commands map[string]models.CogoCLICommand
}

func CreateCLI(deps CogoCLIDeps) *CogoCLI {
	commands := map[string]models.CogoCLICommand{
		RUN_DAEMON:   makeHandleRunAsDaemon(deps.LockService, deps.Logger, deps.CogoDaemon),
		START_DAEMON: makeHandleStartDaemon(deps.LockService, deps.Logger),
		STOP_DAEMON:  makeHandleStopDaemon(deps.LockService, deps.Logger),
		RUN_COMMAND:  makeRunCommand(deps.LockService, deps.Logger),
	}

	service := &CogoCLI{
		commands: commands,
		logger:   deps.Logger,
	}

	return service
}

func (cli *CogoCLI) Handle(args []string) {
	if len(args) < 1 {
		cli.UsageMsg()
		return
	}

	cmdInfo := models.CogoCLIInfo{
		Commad: args[0],
		Args:   args[1:],
	}

	if command := cli.commands[cmdInfo.Commad]; command != nil {
		err := command(cmdInfo)
		if err != nil {
			cli.logger.Fatal(err)
		}

		return
	}

	cli.logger.Info("unkown command has been run", "command", cmdInfo.Commad)

	fmt.Println("Unkown commnad: ", cmdInfo.Commad)
	cli.UsageMsg()
}

func (cli *CogoCLI) UsageMsg() {
	avilableCommands := make([]string, len(cli.commands))

	i := 0
	for k := range cli.commands {
		avilableCommands[i] = k
		i++
	}

	fmt.Printf("Usage: cogo <%s>\n", strings.Join(avilableCommands, " | "))
}
